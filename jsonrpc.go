package jsonrpc

import "encoding/json"
import "strings"

// MsgType ...
type MsgType string

const (
	// InvalidType Invalid message type
	InvalidType MsgType = "invalid"
	// NotificationType Notification message type
	NotificationType MsgType = "notification"
	// RequestType Request message type
	RequestType MsgType = "request"
	// ErrorType Error message type
	ErrorType MsgType = "error"
	// SuccessType Success message type
	SuccessType    MsgType = "success"
	jsonRPCVersion         = "2.0"
)

// RPC represents a JSON-RPC data object.
type RPC struct {
	// Type that should be <'request'|'notification'|'success'|'error'|'invalid'>
	Type    MsgType `json:"-"`
	Version string  `json:"jsonrpc"`
	// A String containing the name of the method to be invoked.
	Method string `json:"method,omitempty"`
	// Object to pass as request parameter to the method.
	Params interface{} `json:"params,omitempty"`
	Result interface{} `json:"result,omitempty"`
	Error  *ErrorObj   `json:"error,omitempty"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	ID interface{} `json:"id,omitempty"`
}

// ErrorObj ...
type ErrorObj struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Request return a JSON-RPC 2.0 request message structures.
// the id must be {String|Integer|nil} type
func Request(id interface{}, method string, args ...interface{}) (result string, err *ErrorObj) {
	if err = validateID(id); err != nil {
		return
	}
	p := &RPC{
		Version: jsonRPCVersion,
		Method:  method,
		ID:      id,
	}
	if len(args) > 0 {
		p.Params = args[0]
	}
	return marshal(p)
}

// Notification return a JSON-RPC 2.0 notification message structures without id.
func Notification(method string, args ...interface{}) (string, *ErrorObj) {
	return Request(nil, method, args...)
}

//Batch return a JSON-RPC 2.0 batch message structures.
func Batch(batch ...string) (arrstr string) {
	if len(batch) == 0 {
		return "[]"
	}
	arrstr = "["
	for index := 0; index < len(batch)-1; index++ {
		arrstr += batch[index]
		arrstr += ","
	}
	arrstr += batch[len(batch)-1]
	arrstr += "]"
	return
}

// Success return a JSON-RPC 2.0 success message structures.
// The result parameter is required
func Success(id interface{}, msg interface{}) (result string, err *ErrorObj) {
	if msg == nil {
		return result, InternalError()
	}
	if err = validateID(id); err != nil {
		return
	}
	p := &RPC{
		Version: jsonRPCVersion,
		Result:  msg,
		ID:      id,
	}
	return marshal(p)
}

//Error return a JSON-RPC 2.0 error message structures.
func Error(id interface{}, rpcerr *ErrorObj) (result string, err *ErrorObj) {
	if err = validateID(id); err != nil {
		return
	}
	p := &RPC{
		Version: jsonRPCVersion,
		Error:   rpcerr,
		ID:      id,
	}
	return marshal(p)
}

// ErrorFrom return an ErrorObj by native error
func ErrorFrom(err error) (obj *ErrorObj) {
	return &ErrorObj{Code: -32701, Message: err.Error()}
}

// ErrorWith return an ErrorObj by arguments
func ErrorWith(code int, msg string, data ...interface{}) (obj *ErrorObj) {
	obj = &ErrorObj{Code: code, Message: msg}
	if len(data) > 0 {
		obj.Data = data[0]
	}
	return obj
}

// ParseError invalid JSON was received by the server.An error occurred on the server while parsing the JSON text.
func ParseError(data ...interface{}) *ErrorObj {
	return ErrorWith(-32700, "Parse error", data...)
}

// InvalidRequest the request is not a valid Request object.
func InvalidRequest(data ...interface{}) *ErrorObj {
	return ErrorWith(-32600, "Invalid Request", data...)
}

// MethodNotFound the method does not exist or is not available.
func MethodNotFound(data ...interface{}) *ErrorObj {
	return ErrorWith(-32601, "Method not found", data...)
}

// InvalidParams Invalid method parameter(s).
func InvalidParams(data ...interface{}) *ErrorObj {
	return ErrorWith(-32602, "Invalid params", data...)
}

// InternalError  Internal JSON-RPC error.
func InternalError(data ...interface{}) *ErrorObj {
	return ErrorWith(-32603, "Internal error", data...)
}

// Parse return jsonrpc 2.0 message object, ignore the first return value if the msg is batch rpc.
func Parse(msg string) (req *RPC, batch []*RPC) {
	if strings.HasPrefix(msg, "[") && strings.HasSuffix(msg, "]") {
		batch = make([]*RPC, 1)
		if err := validateMsg(msg, &batch); err == nil {
			for _, val := range batch {
				parse(val)
			}
		} else {
			req := &RPC{}
			req.Error = err
			req.Type = InvalidType
			batch[0] = req
		}
		return nil, batch
	}
	req = &RPC{}
	if err := validateMsg(msg, req); err != nil {
		req.Error = err
		req.Type = InvalidType
	} else {
		parse(req)
	}
	return req, nil
}
func parse(r *RPC) {
	if r.Version != jsonRPCVersion {
		r.Type = InvalidType
		r.Error = InvalidRequest()
		return
	}
	if err := validateID(r.ID); err != nil {
		r.Error = err
		return
	}
	if r.Method != "" { //Request
		if r.ID == nil {
			r.Type = NotificationType
		} else {
			r.Type = RequestType
		}
	} else {
		if r.Error != nil { // Response
			r.Type = ErrorType
		} else if r.Result != nil {
			r.Type = SuccessType
		} else {
			r.Type = InvalidType
			r.Error = InvalidRequest()
		}
	}
	return
}
func validateID(id interface{}) (err *ErrorObj) {
	if id != nil {
		switch id.(type) {
		case string:
		case int:
		default:
			err = InternalError()
		}
	}
	return
}
func validateMsg(msg string, p interface{}) *ErrorObj {
	if msg == "" {
		return InvalidRequest()
	}
	err := json.Unmarshal([]byte(msg), p)
	if err != nil {
		return ParseError(err.Error())
	}
	return nil
}
func marshal(v interface{}) (string, *ErrorObj) {
	data, errs := json.Marshal(v)
	if errs != nil {
		return "", InternalError(errs.Error())
	}
	return string(data), nil
}
