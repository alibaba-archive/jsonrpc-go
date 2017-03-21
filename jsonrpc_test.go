package jsonrpc

import "testing"
import "github.com/stretchr/testify/assert"

var (
	invalid      = "invalid"
	notification = "notification"
	request      = "request"
	errorType    = "error"
	success      = "success"
)

func TestProducer(t *testing.T) {

	t.Run("jsonrpc with request that should be", func(t *testing.T) {
		assert := assert.New(t)

		val, err := Request(123, "update")
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":123}", val)

		val, err = Request("123", "update")
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}", val)

		val, err = Request(true, "update")
		assert.NotNil(err)
		assert.Equal("invalid id that MUST contain a String, Number, or NULL value", err.Error())

		val, err = Notification("update")
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\"}", val)

		val, err = Notification("update", 0)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"params\":0}", val)

	})
	t.Run("jsonrpc with Batch func that should be", func(t *testing.T) {
		assert := assert.New(t)

		Request("1", "sum")

		val, err := Request(123, "update")
		assert.Nil(err)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":123}", val)

		val2, err := Request("123", "update")
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}", val2)

		val = Batch(val, val2)
		assert.Equal("[{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":123},{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}]", val)

		val = Batch()
		assert.Equal("[]", val)
	})
	t.Run("jsonrpc with response that should be", func(t *testing.T) {
		assert := assert.New(t)

		val, err := Success("123", nil)
		assert.NotNil(err)
		assert.Equal(errResultArgument, err)

		val, err = Success("123", "OK")
		assert.Nil(err)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"result\":\"OK\",\"id\":\"123\"}", val)

		val, err = Success(123, []string{})
		assert.Nil(err)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"result\":[],\"id\":123}", val)

		val, err = Success(true, "")
		assert.NotNil(err)
		assert.Equal("invalid id that MUST contain a String, Number, or NULL value", err.Error())

		rpcerr := ErrorWith(1, "test")
		val, err = Error(nil, rpcerr)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":1,\"message\":\"test\"}}", val)

		val, err = Error(true, rpcerr)
		assert.NotNil(err)
		assert.Equal("invalid id that MUST contain a String, Number, or NULL value", err.Error())

		objerr := ErrorFrom(errJsonrpcObject)
		assert.NotNil(objerr)
		assert.Equal(objerr.Message, errJsonrpcObject.Error())

		rpcerr = ErrorWith(1, "test", "xx")
		val, err = Error(nil, rpcerr)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":1,\"message\":\"test\",\"data\":\"xx\"}}", val)
	})

	t.Run("jsonrpc with ParseReq func that should be", func(t *testing.T) {

		assert := assert.New(t)

		val, err := Parse("")
		assert.Empty(val)
		assert.NotNil(err)
		assert.Equal(errEmptyMessage, err)

		val, err = Parse("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}")
		assert.Nil(err)
		assert.Equal("123", val.PlayLoad.ID)
		assert.Equal("update", val.PlayLoad.Method)
		assert.Equal(request, val.Type)

		val, err = Parse("{\"jsonrpc\":\"2.0,\"result\":\"OK\",\"id\":\"123\"}")
		assert.NotNil(err)
		assert.Equal("invalid jsonrpc message structures", err.Error())

		val, err = Parse("{\"jsonrpc\":\"3.0\",\"result\":\"OK\",\"id\":\"123\"}")
		assert.NotNil(err)
		assert.Equal("invalid jsonrpc version", err.Error())

		val, err = Parse("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"params\":0}")
		assert.Nil(err)
		assert.Equal(float64(0), val.PlayLoad.Params)
		assert.Equal("update", val.PlayLoad.Method)
		assert.Equal(notification, val.Type)

		val, err = Parse("{\"jsonrpc\":\"2.0\",\"params\":0,\"id\":\"123\"}")
		assert.NotNil(err)
		assert.Equal("invalid jsonrpc object", err.Error())

	})
	t.Run("jsonrpc with ParseReply fun that should be", func(t *testing.T) {
		assert := assert.New(t)

		val, err := ParseReply("")
		assert.Empty(val)
		assert.NotNil(err)
		assert.Equal(errEmptyMessage, err)

		val, err = ParseReply("{\"jsonrpc\":\"2.0\",\"result\":\"OK\",\"id\":\"123\"}")
		assert.Nil(err)
		assert.Equal("123", val.PlayLoad.ID)
		assert.Equal("OK", val.PlayLoad.Result)
		assert.Equal(success, val.Type)

		val, err = ParseReply("{\"jsonrpc\":\"2.0,\"result\":\"OK\",\"id\":\"123\"}")
		assert.NotNil(err)
		assert.Equal("invalid jsonrpc message structures", err.Error())

		val, err = ParseReply("{\"jsonrpc\":\"3.0\",\"result\":\"OK\",\"id\":\"123\"}")
		assert.NotNil(err)
		assert.Equal("invalid jsonrpc version", err.Error())

		val, err = ParseReply("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":1,\"message\":\"test\"}}")
		assert.Nil(err)
		assert.Equal("test", val.PlayLoad.Error.Message)
		assert.Equal(1, val.PlayLoad.Error.Code)
		assert.Equal(errorType, val.Type)

		val, err = ParseReply("{\"jsonrpc\":\"2.0\",\"id\":\"123\"}")
		if assert.NotNil(err) {
			assert.Equal("invalid jsonrpc object", err.Error())
		}

		val, err = ParseReply("{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"1\"}")
		assert.Nil(err)
		assert.Equal("1", val.PlayLoad.ID)
		assert.Equal(-32601, val.PlayLoad.Error.Code)
		assert.Equal(errorType, val.Type)

	})
	t.Run("jsonrpc with ParseBatchReply func that should be", func(t *testing.T) {
		arr := "[{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": null},{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"1\"},{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"2\"}]"
		assert := assert.New(t)

		val, err := ParseBatchReply(arr)

		if assert.Nil(err) {
			assert.Equal(3, len(val))
			assert.Equal(nil, val[0].PlayLoad.ID)
			assert.Equal(-32601, val[0].PlayLoad.Error.Code)
			assert.Equal("Method not found", val[0].PlayLoad.Error.Message)
			assert.Equal("1", val[1].PlayLoad.ID)
			assert.Equal(-32601, val[1].PlayLoad.Error.Code)
			assert.Equal("Method not found", val[1].PlayLoad.Error.Message)
			assert.Equal("2", val[2].PlayLoad.ID)
			assert.Equal(-32601, val[2].PlayLoad.Error.Code)
			assert.Equal("Method not found", val[2].PlayLoad.Error.Message)
		}

		val, err = ParseBatchReply("")
		assert.Equal(errEmptyMessage, err)

		str := `[
        {"jsonrpc": "2.0", "result": 7, "id": "1"},
        {"jsonrpc": "2.0", "result": 19, "id": "2"},
        {"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
        {"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": 5},
        {"jsonrpc": "2.0", "result": ["hello", 5], "id": "9"}
      ]`

		val, err = ParseBatchReply(str)
		assert.Equal("1", val[0].PlayLoad.ID)
		assert.Equal(float64(19), val[1].PlayLoad.Result)
		assert.Equal("Invalid Request", val[2].PlayLoad.Error.Message)
		assert.Equal(float64(5), val[3].PlayLoad.ID)
		assert.Equal([]interface{}{"hello", float64(5)}, val[4].PlayLoad.Result.([]interface{}))

	})
	t.Run("jsonrpc with ParseReqBatch func that should be", func(t *testing.T) {
		arr := `[
			{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
			{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]},
			{"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": "2"},
			{"foo": "boo"},
			{"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": "5"},
			{"jsonrpc": "1.0", "method": "get_data", "id": "9"} 
    	]`
		assert := assert.New(t)

		val, err := ParseBatch(arr)

		assert.Nil(err)

		if assert.Equal(6, len(val)) {
			assert.Equal("1", val[0].PlayLoad.ID)
			assert.Equal("notify_hello", val[1].PlayLoad.Method)
			assert.Equal([]interface{}{float64(42), float64(23)}, val[2].PlayLoad.Params.([]interface{}))
			assert.Equal(invalid, val[3].Type)
			assert.Equal("myself", val[4].PlayLoad.Params.(map[string]interface{})["name"])
			assert.Equal(invalid, val[5].Type)
		}
		val, err = ParseBatch("")
		assert.Equal("empty message", err.Error())
	})
	t.Run("jsonrpc with RandID func that should be", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal("Parse error", ParseError().Message)
		assert.Equal("Invalid Request", InvalidRequest().Message)
		assert.Equal("Method not found", MethodNotFound().Message)
		assert.Equal("Invalid params", InvalidParams().Message)
		assert.Equal("Internal error", InternalError().Message)
	})
}
