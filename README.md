# oas2

[![GoDoc](https://godoc.org/github.com/hypnoglow/oas2?status.svg)](https://godoc.org/github.com/hypnoglow/oas2)
[![CircleCI](https://circleci.com/gh/hypnoglow/oas2.svg?style=shield)](https://circleci.com/gh/hypnoglow/oas2)
[![codecov](https://codecov.io/gh/hypnoglow/oas2/branch/master/graph/badge.svg)](https://codecov.io/gh/hypnoglow/oas2)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypnoglow/oas2)](https://goreportcard.com/report/github.com/hypnoglow/oas2)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

**WIP. Not stable yet. API may change at any time. Use a vendoring tool of your
preference to lock an exact [release](https://github.com/hypnoglow/oas2/releases) version.**

Package oas2 provides utilities to work with OpenAPI 2.0 specification
(aka Swagger).

The purpose of this package is to provide utilities for building APIs
around the OpenAPI specification in Go idiomatic way on top of `net/http`.

You don't need to learn any special framework or write `net/http`-incompatible
code - just delegate request validation, request parameters decoding
and other routines to this library - and focus on your application logic.

## Features

### Router from a spec

This package provides an easy way to automatically create a router supporting
all resources from your OpenAPI specification file. The underlying router is only
your choice - you can use [gorilla/mux](https://github.com/gorilla/mux), [chi](https://github.com/go-chi/chi)
or any other.

Let's dive into a simple example.

Given a spec: [petstore.yaml](examples/petstore.yaml)

First of all, load your spec in your app:

```go
// specPath is a path to your spec file.
doc, _ := oas2.LoadSpec(specPath)
```

Next, create an [operation](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#operationObject) handler. 
Let's define a handler for `loginUser` operation:

```go
func loginHandler(w http.ResponseWriter, req *http.Request) {
    username := req.URL.Query().Get("username")
    password := req.URL.Query().Get("password")
    
    if login(username, password) {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    w.WriteHeader(http.StatusBadRequest)
}
```

```go
handlers := oas2.OperationHandlers{
    "loginUser":    http.HandlerFunc(loginHandler),
}
```

Define what options (logger, middleware) you will use:

```go
logger := logrus.New()
logger.SetLevel(logrus.DebugLevel)
queryValidator := oas2.NewQueryValidator(errHandler)
```

Create a router:

```go
router, _ := oas2.NewRouter(
    doc.Spec(), 
    handlers, 
    oas2.DebugLog(logger.Debugf), 
    oas2.Use(queryValidator)
)
```

Then you can use your `router` as an argument for `http.ListenAndServe` 
or even as a subrouter for the other router.

```go
http.ListenAndServe(":8080", router)
``` 

Now the server handles requests based on the paths defined in the given spec.
It validates request query parameters against the spec and runs `errHandler` 
func if any error occured during validation. The router also sets the operation
identifier to each request's context, so it can be used in a handler or any custom
middleware.

See the full [example](examples/router/main.go) for the complete code.

### Decode query parameters to a struct

Given request query parameters: `?name=John&age=27`

Given OpenAPI v2 (swagger) spec:

```yaml
...
parameters:
- name: name
  type: string
- name: age
  type: integer
  format: int32
- name: loves_apples
  type: bool
  default: true
...
```

In your Go code you create a struct:

```go
type Member struct {
	Name        string `oas:"name"`
	Age         int32  `oas:"age"`
	LovesApples bool   `oas:"loves_apples"`
}
```

And populate it:

```go
var m Member 
oas2.DecodeQuery(paramSpec, req.URL.Query(), &m)

fmt.Printf("%#v", m) // Member{Name:"John", Age:27, LovesApples:true}
```

See [`godoc example`](https://godoc.org/github.com/hypnoglow/oas2#example-DecodeQuery) for complete example code.

## License

[MIT](https://github.com/hypnoglow/oas2/blob/master/LICENSE).
