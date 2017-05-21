# esrest 
Easy, elegant, fluent HTTP client API for Go

[![Travis branch](https://img.shields.io/travis/easonlin404/esrest/master.svg)](https://travis-ci.org/easonlin404/esrest)
[![Codecov branch](https://img.shields.io/codecov/c/github/easonlin404/esrest/master.svg)](https://codecov.io/gh/easonlin404/esrest)

## Features
* __100%__ code coverage
* Support __GET__/__POST__/__PUT__/__DELETE__ http methods
* Only use __Body__ chain method to send payload(__JSON__/__String__/__Slice__) 
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
Set header (Default ContentType is "application/json")
``` go
res, err := esrest.New().
		    Get("http://httpbin.org/get").
		    Header("MyHader", "headvalue").
		    Do()
```

Send __JSON__/__String__/__Slice__ payload only call same __Body__ chain method:
``` go
//JSON
json := struct {
		Message string `json:"message"`
	}{"ok"}

res, err := esrest.New().
		    Post("http://httpbin.org/post").
		    Body(json).
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

