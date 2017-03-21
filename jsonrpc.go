package jsonrpc

import (
	"encoding/json"
	"errors"
)

const (

	// Invalid message type
	Invalid = "invalid"
	// NotificationType message type
	NotificationType = "notification"
	// RequestType message type
	RequestType = "request"
	// ErrorType message type
	ErrorType = "error"
	// SuccessType message type
	SuccessType = "success"

	jsonRPCVersion = "2.0"
)

var (
	// ErrEmptyMessage means empty jsonrpc message error
	errEmptyMessage = errors.New("Empty jsonrpc message")
	// ErrResultArgument means invalid 'result' argument
	errResultArgument = errors.New("Invalid 'result' argument that's required")
	// ErrJsonrpcVersion means invalid jsonrpc version
	errJsonrpcVersion = errors.New("invalid jsonrpc version")
	// ErrJsonrpcObject means invalid jsonrpc object
	errJsonrpcObject = errors.New("invalid jsonrpc object")
)

// ClientRequest represents a JSON-RPC data from client.
type ClientRequest struct {
	Type     string
	PlayLoad *PayloadReq
}

// PayloadReq represents a JSON-RPC request.
type PayloadReq struct {
	Version string `json:"jsonrpc"`
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params interface{} `json:"params,omitempty"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	ID interface{} `json:"id,omitempty"`
}

// ClientResponse represents a JSON-RPC response returned to a client.
type ClientResponse struct {
	Type     string
	PlayLoad *PayloadRes
}

// PayloadRes represents a JSON-RPC request.
type PayloadRes struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *ErrorObj   `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

// ErrorObj ...
type ErrorObj struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Request creates a JSON-RPC 2.0 request message structures.
// the id must be {String|Integer|nil} type
func Request(id interface{}, method string, args ...interface{}) (result string, err error) {
	if err = validateID(id); err != nil {
		return
	}
	p := &PayloadReq{
		Version: jsonRPCVersion,
		Method:  method,
		ID:      id,
	}
	if len(args) > 0 {
		p.Params = args[0]
	}
	val, err := json.Marshal(p)
	return string(val), err
}

// Notification Creates a JSON-RPC 2.0 notification message structures.
func Notification(method string, args ...interface{}) (string, error) {
	return Request(nil, method, args...)
}

//Batch Creates a JSON-RPC 2.0 batch message structures.
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

// Parse Parse message of from client request.
func Parse(msg string) (req *ClientRequest, err error) {
	p := make(map[string]interface{}, 0)
	if err = validateMsg(msg, &p); err != nil {
		return
	}
	req, err = parseReqMap(p)
	return
}

// ParseBatch Parse batch message of from client request.
func ParseBatch(msg string) (req []*ClientRequest, err error) {
	if msg == "" || len(msg) < 2 {
		err = errors.New("empty message")
		return
	}
	payloads := make([]map[string]interface{}, 0)
	if err = validateMsg(msg, &payloads); err == nil {
		req = make([]*ClientRequest, len(payloads))
		for i, val := range payloads {
			r, _ := parseReqMap(val)
			req[i] = r
		}
	}
	return
}

// ParseReply Parse message of from server reply.
func ParseReply(msg string) (res *ClientResponse, err error) {
	p := make(map[string]interface{}, 0)
	if err = validateMsg(msg, &p); err != nil {
		return
	}
	return parseResMap(p)
}

// ParseBatchReply Parse message of server reply batch request.
func ParseBatchReply(msg string) (res []*ClientResponse, err error) {
	if msg == "" || len(msg) < 2 {
		return res, errEmptyMessage
	}
	payloads := make([]map[string]interface{}, 0)
	if err = validateMsg(msg, &payloads); err == nil {
		res = make([]*ClientResponse, len(payloads))
		for i, val := range payloads {
			r, _ := parseResMap(val)
			res[i] = r
		}
	}
	return
}

// Success Creates a JSON-RPC 2.0 success response object, return JsonRpc json.
// The result parameter is required
func Success(id interface{}, result interface{}) (str string, err error) {
	if result == nil {
		return str, errResultArgument
	}
	if err = validateID(id); err != nil {
		return
	}
	p := &PayloadRes{
		Version: jsonRPCVersion,
		Result:  result,
		ID:      id,
	}
	data, err := json.Marshal(p)
	return string(data), err
}

//Error Creates a JSON-RPC 2.0 error response object, return JsonRpc json.
func Error(id interface{}, rpcerr *ErrorObj) (str string, err error) {
	if err = validateID(id); err != nil {
		return
	}
	p := &PayloadRes{
		Version: jsonRPCVersion,
		Error:   rpcerr,
		ID:      id,
	}
	data, err := json.Marshal(p)
	return string(data), err
}

// ErrorFrom return an ErrorObj by argument
func ErrorFrom(err error) (obj *ErrorObj) {
	return &ErrorObj{Code: -32701, Message: err.Error()}
}

// ErrorWith a JsonRpc error
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
func parseResMap(val map[string]interface{}) (r *ClientResponse, err error) {
	r = &ClientResponse{PlayLoad: &PayloadRes{}}
	if version, ok := val["jsonrpc"]; ok {
		r.PlayLoad.Version = version.(string)
	}
	if ID, ok := val["id"]; ok {
		r.PlayLoad.ID = ID
	}
	if errs, ok := val["error"]; ok {
		r.PlayLoad.Error = &ErrorObj{}
		r.Type = ErrorType
		if err2, ok := errs.(map[string]interface{}); ok {
			r.PlayLoad.Error.Code = int(err2["code"].(float64))
			r.PlayLoad.Error.Message = err2["message"].(string)
			r.PlayLoad.Error.Data = err2["data"]
		}
	} else if result, ok := val["result"]; ok {
		r.PlayLoad.Result = result
		r.Type = SuccessType
	} else {
		r.Type = Invalid
		err = errJsonrpcObject
	}
	if r.PlayLoad.Version != jsonRPCVersion {
		r.Type = Invalid
		err = errJsonrpcVersion
	}
	return
}
func parseReqMap(val map[string]interface{}) (r *ClientRequest, err error) {
	r = &ClientRequest{PlayLoad: &PayloadReq{}}
	if version, ok := val["jsonrpc"]; ok {
		r.PlayLoad.Version, _ = version.(string)
	}
	if ID, ok := val["id"]; ok {
		r.PlayLoad.ID = ID
	}
	if method, ok := val["method"]; ok {
		r.PlayLoad.Method, _ = method.(string)
	}
	if params, ok := val["params"]; ok {
		r.PlayLoad.Params = params
	}
	if r.PlayLoad.Version != jsonRPCVersion {
		r.Type = Invalid
		err = errJsonrpcVersion
	} else if r.PlayLoad.Method == "" {
		r.Type = Invalid
		err = errJsonrpcObject
	} else if r.PlayLoad.ID == nil {
		r.Type = NotificationType
	} else {
		r.Type = RequestType
	}
	return
}
func validateID(id interface{}) (err error) {
	if id != nil {
		switch id.(type) {
		case string:
		case int:
		default:
			err = errors.New("invalid id that MUST contain a String, Number, or NULL value")
		}
	}
	return
}
func validateMsg(msg string, p interface{}) (err error) {
	if msg == "" {
		return errEmptyMessage
	}
	err = json.Unmarshal([]byte(msg), p)
	if err != nil {
		err = errors.New("invalid jsonrpc message structures")
	}
	return
}
