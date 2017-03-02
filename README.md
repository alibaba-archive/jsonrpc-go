# jsonrpc

[![Build Status](https://travis-ci.org/teambition/jsonrpc.svg?branch=master)](https://travis-ci.org/teambition/jsonrpc)
[![Coverage Status](http://img.shields.io/coveralls/teambition/jsonrpc.svg?style=flat-square)](https://coveralls.io/r/teambition/jsonrpc)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/teambition/jsonrpc/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/teambition/jsonrpc)

jsonrpc is an golang implementation just for parse and serialize JSON-RPC2 messages, it's easy to integration with your any application.

Inspired by [https://github.com/teambition/jsonrpc-lite](https://github.com/teambition/jsonrpc-lite) and [JSON-RPC 2.0 specifications](http://jsonrpc.org/specification)
## Installation
```go
go get github.com/teambition/jsonrpc
```

## API
####  jsonrpc.Request(id, method[, params])
Creates a JSON-RPC 2.0 message structures 
- `id`: {String|Int|nil}
- `method`: {String}
- `params`:  {interface{}}, optional
```go
	val, err := jsonrpc.Request(123, "update")
	{
    "jsonrpc": "2.0",
    "method": "update",
    "id": 123
    }
```
####  jsonrpc.Request2(method[, params])
Creates a JSON-RPC 2.0 message structures, the id is automatic generation by strconv.FormatInt(rand.Int63(), 10)
- `method`: {String}
- `params`:  {interface{}}, optional
```go
	val, err := jsonrpc.Request("update")
	{
    "jsonrpc": "2.0",
    "method": "update",
    "id": randnum
    }
```
####  jsonrpc.Notification(method[, params])
Creates a JSON-RPC 2.0 notification message structures
- `method`: {String}
- `params`:  {interface{}}, optional
```go
    val, err = jsonrpc.Notification("update")
    {
    "jsonrpc": "2.0",
    "method": "update"
    }
```   