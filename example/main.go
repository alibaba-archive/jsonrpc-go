package main

import "github.com/teambition/jsonrpc"

func main() {
	arr := "[{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": null},{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"1\"},{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32601, \"message\": \"Method not found\"}, \"id\": \"1\"}]"

	_, err := jsonrpc.ParseBatchReply(arr)
	if err != nil {
		println(err.Error())
	}
}
