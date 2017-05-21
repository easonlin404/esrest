# esrest 
Easy, elegant, fluent HTTP client API for Go

[![Travis branch](https://img.shields.io/travis/easonlin404/esrest/master.svg)](https://travis-ci.org/easonlin404/esrest)
[![Codecov branch](https://img.shields.io/codecov/c/github/easonlin404/esrest/master.svg)](https://codecov.io/gh/easonlin404/esrest)

## Features
* Support __GET__/__POST__/__PUT__/__DELETE__ http methods
* Only use `Body` chain method to send payload(JSON/string/slice) 
* Receive unmarshal JSON
* todo

## Installation
```sh
$ go get github.com/easonlin404/esrest
```
## Usage

__GET__/__POST__/__PUT__/__DELETE__
```go

res, err := esrest.New().Get("http://httpbin.org/get").Do()

```
Add header (Default ContentType is "application/json")
``` go
res, err := esrest.New().
		    Get("http://httpbin.org/get").
		    Header("MyHader", "headvalue").
		    Do()
```

Sending _JSON_ payload use call `Body` chain method smae as other:
``` go
//JSON
res, err := esrest.New().
		    Post("http://httpbin.org/post").
		    Body(struct {
                 		Message string `json:"message"`
                 	}{"ok"}).
		    Do()
```

``` go
//slice
res, err := esrest.New().
		    Post("http://httpbin.org/post").
		    Body([]byte(`{"message":"ok"}`)).
		    Do()
```
``` go
//string
res, err := esrest.New().
		    Post("http://httpbin.org/post").
		    Body(string(`{"message":"ok"}`)).
		    Do()
```
Add Query parameter:
``` go
res, err := esrest.New().
		    Get("http://httpbin.org/get").
		    Query("Param1", "value").
		    Do()
```

Receive unmarshal JSON:
``` go
json := struct {
		Message string `json:"message"`
	}{}
res, err := esrest.New().
		    Post("http://httpbin.org/post").
		    DoJson(json)
```

