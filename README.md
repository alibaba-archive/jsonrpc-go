# jsonrpc

[![Build Status](https://travis-ci.org/teambition/jsonrpc-go.svg?branch=master)](https://travis-ci.org/teambition/jsonrpc-go)
[![Coverage Status](http://img.shields.io/coveralls/teambition/jsonrpc-go.svg?style=flat-square)](https://coveralls.io/r/teambition/jsonrpc-go)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/teambition/jsonrpc-go/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/teambition/jsonrpc-go)

jsonrpc is an golang implementation just for parse and serialize JSON-RPC2 messages, it's easy to integration with your any application.

Inspired by [https://github.com/teambition/jsonrpc-lite](https://github.com/teambition/jsonrpc-lite) and [JSON-RPC 2.0 specifications](http://jsonrpc.org/specification)
## Installation
```go
go get github.com/teambition/jsonrpc-go
```

## API
### Generate jsonrpc 2.0 message.
```go
// Request creates a JSON-RPC 2.0 request message structures.
// the id must be {String|Integer|nil} type
func Request(id interface{}, method string, args ...interface{}) (result string, err *ErrorObj)
// Notification Creates a JSON-RPC 2.0 notification message structures without id.
func Notification(method string, args ...interface{}) (string, *ErrorObj) 
//Batch return a JSON-RPC 2.0 batch message structures.
func Batch(batch ...string) (arrstr string)
// Success return a JSON-RPC 2.0 success message structures.
// The result parameter is required
func Success(id interface{}, msg interface{}) (result string, err *ErrorObj)
//Error return a JSON-RPC 2.0 error message structures.
func Error(id interface{}, rpcerr *ErrorObj) (result string, err *ErrorObj) 
``` 

## Parse jsonrpc 2.0 message structures.
```go
// Parse return jsonrpc 2.0 message object.
func Parse(msg string) (req *RPC) 
```           