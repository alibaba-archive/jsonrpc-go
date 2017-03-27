package jsonrpc

import (
	"errors"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	request := func(id interface{}, method string, args ...interface{}) string {
		str, err := Request(id, method, args...)
		if err != nil {
			c, _ := json.Marshal(err)
			str = string(c)
		}
		return str
	}
	notification := func(method string, args ...interface{}) string {
		str, err := Notification(method, args...)
		if err != nil {
			c, _ := json.Marshal(err)
			str = string(c)
		}
		return str
	}
	success := func(id interface{}, result interface{}) string {
		str, err := Success(id, result)
		if err != nil {
			c, _ := json.Marshal(err)
			str = string(c)
		}
		return str
	}
	rpcerr := func(id interface{}, rpcerr *ErrorObj) (str string) {
		str, err := Error(id, rpcerr)
		if err != nil {
			c, _ := json.Marshal(err)
			str = string(c)
		}
		return str
	}
	cases := []struct {
		value string
		want  string
	}{

		{request(123, "update"), "{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":123}"},
		{request("123", "update"), "{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}"},
		{request(123, "update", make(chan bool)), "{\"code\":-32603,\"message\":\"Internal error\",\"data\":\"json: unsupported type: chan bool\"}"},
		{request(true, "update"), "{\"code\":-32603,\"message\":\"Internal error\"}"},
		{notification("update"), "{\"jsonrpc\":\"2.0\",\"method\":\"update\"}"},
		{notification("update", 0), "{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"params\":0}"},
		{Batch(request(123, "update"), request("123", "update")), "[{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":123},{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}]"},
		{Batch(), "[]"},
		{success("123", nil), "{\"code\":-32603,\"message\":\"Internal error\"}"},
		{success("123", "OK"), "{\"jsonrpc\":\"2.0\",\"result\":\"OK\",\"id\":\"123\"}"},
		{success(123, []string{}), "{\"jsonrpc\":\"2.0\",\"result\":[],\"id\":123}"},
		{success(true, ""), "{\"code\":-32603,\"message\":\"Internal error\"}"},
		{rpcerr(nil, ErrorWith(1, "test")), "{\"jsonrpc\":\"2.0\",\"error\":{\"code\":1,\"message\":\"test\"}}"},
		{rpcerr(true, ErrorWith(1, "test")), "{\"code\":-32603,\"message\":\"Internal error\"}"},
		{ErrorFrom(errors.New("invalid jsonrpc object")).Message, "invalid jsonrpc object"},
		{rpcerr(nil, ErrorWith(1, "test", "xx")), "{\"jsonrpc\":\"2.0\",\"error\":{\"code\":1,\"message\":\"test\",\"data\":\"xx\"}}"},
		{ParseError().Message, "Parse error"},
		{InvalidRequest().Message, "Invalid Request"},
		{MethodNotFound().Message, "Method not found"},
		{InvalidParams().Message, "Invalid params"},
		{InternalError().Message, "Internal error"},
	}
	for _, item := range cases {
		if item.value != item.want {
			t.Errorf("Encode get: %q ,want: %q", item.value, item.want)
		}
	}
}
func TestParse(t *testing.T) {
	assert := assert.New(t)
	cases := []string{
		"{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}",
		"{\"jsonrpc\":\"2.0\"result\":\"OK\",\"id\":\"123\"}",
		"{\"jsonrpc\":\"3.0\",\"result\":\"OK\",\"id\":\"123\"}",
		"{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"params\":0}",
		"{\"jsonrpc\":\"2.0\",\"params\":0,\"id\":\"123\"}",
		"{\"jsonrpc\":\"2.0\",\"result\":\"OK\",\"id\":123}",
		"{\"jsonrpc\":\"2.0\",\"result\":Null,\"id\":123}",
	}

	val, _ := Parse(cases[0])
	assert.Nil(val.Error)
	assert.Equal("123", val.ID)
	assert.Equal("update", val.Method)
	assert.Equal(RequestType, val.Type)

	val, _ = Parse(cases[1])
	assert.Equal("Parse error", val.Error.Message)

	val, _ = Parse(cases[2])
	assert.Equal("Invalid Request", val.Error.Message)

	val, _ = Parse(cases[3])
	assert.Nil(val.Error)
	assert.Equal(float64(0), val.Params)
	assert.Equal("update", val.Method)
	assert.Equal(NotificationType, val.Type)

	val, _ = Parse(cases[4])
	assert.Equal("Invalid Request", val.Error.Message)

	val, _ = Parse(cases[5])
	assert.Equal(float64(123), val.ID)
	assert.Equal("OK", val.Result)

	val, _ = Parse(cases[6])
	assert.Equal(InvalidType, val.Type)
	assert.Equal("Parse error", val.Error.Message)
	assert.NotNil(val.Error.Data)

	cases = []string{
		"",
		"{\"jsonrpc\":\"2.0\",\"result\":\"OK\",\"id\":\"123\"}",
		"{\"jsonrpc\":\"2.0,\"result\":\"OK\",\"id\":\"123\"}",
		"{\"jsonrpc\":\"3.0\",\"result\":\"OK\",\"id\":\"123\"}",
		"{\"jsonrpc\":\"2.0\",\"error\":{\"code\":1,\"message\":\"test\"}}",
		"{\"jsonrpc\":\"2.0\",\"id\":\"123\"}",
		"{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"1\"}",
	}

	val, _ = Parse(cases[0])

	assert.Equal("Invalid Request", val.Error.Message)

	val, _ = Parse(cases[1])
	assert.Nil(val.Error)
	assert.Equal("123", val.ID)
	assert.Equal("OK", val.Result)
	assert.Equal(SuccessType, val.Type)

	val, _ = Parse(cases[2])
	assert.NotNil(val.Error)
	assert.Equal("Parse error", val.Error.Message)

	val, _ = Parse(cases[3])
	assert.Equal("Invalid Request", val.Error.Message)

	val, _ = Parse(cases[4])
	assert.NotNil(val.Error)
	assert.Equal("test", val.Error.Message)
	assert.Equal(1, val.Error.Code)
	assert.Equal(ErrorType, val.Type)

	val, _ = Parse(cases[5])
	assert.Equal("Invalid Request", val.Error.Message)

	val, _ = Parse(cases[6])
	assert.NotNil(val.Error)
	assert.Equal("1", val.ID)
	assert.Equal(-32601, val.Error.Code)
	assert.Equal(ErrorType, val.Type)

}
func TestParseBatch(t *testing.T) {

	arr := "[{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": null},{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"1\"},{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"2\"}]"
	assert := assert.New(t)

	_, val := Parse(arr)

	assert.Equal(3, len(val))
	assert.Equal(nil, val[0].ID)
	assert.Equal(-32601, val[0].Error.Code)
	assert.Equal("Method not found", val[0].Error.Message)
	assert.Equal("1", val[1].ID)
	assert.Equal(-32601, val[1].Error.Code)
	assert.Equal("Method not found", val[1].Error.Message)
	assert.Equal("2", val[2].ID)
	assert.Equal(-32601, val[2].Error.Code)
	assert.Equal("Method not found", val[2].Error.Message)

	req, _ := Parse("")
	assert.Equal("Invalid Request", req.Error.Message)

	arr = `[
        {"jsonrpc": "2.0", "result": 7, "id": "1"},
        {"jsonrpc": "2.0", "result": 19, "id": "2"},
        {"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
        {"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": 5},
        {"jsonrpc": "2.0", "result": ["hello", 5], "id": "9"}
      ]`

	_, val = Parse(arr)
	assert.Equal("1", val[0].ID)
	assert.Equal(float64(19), val[1].Result)
	assert.Equal("Invalid Request", val[2].Error.Message)
	assert.Equal(float64(5), val[3].ID)
	assert.Equal([]interface{}{"hello", float64(5)}, val[4].Result.([]interface{}))

	arr = `[
			{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
			{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]},
			{"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": "2"},
			{"foo": "boo"},
			{"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": "5"},
			{"jsonrpc": "1.0", "method": "get_data", "id": "9"} 
    	]`

	_, val = Parse(arr)

	if assert.Equal(6, len(val)) {
		assert.Equal("1", val[0].ID)
		assert.Equal("notify_hello", val[1].Method)
		assert.Equal([]interface{}{float64(42), float64(23)}, val[2].Params.([]interface{}))
		assert.Equal(InvalidType, val[3].Type)
		assert.Equal("myself", val[4].Params.(map[string]interface{})["name"])
		assert.Equal(InvalidType, val[5].Type)
	}
	req, _ = Parse("")
	assert.Equal("Invalid Request", req.Error.Message)

	_, val = Parse(`[x:x]`)
	assert.Equal("Parse error", val[0].Error.Message)
}
