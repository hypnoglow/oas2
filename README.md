# oas2

[![GoDoc](https://godoc.org/github.com/hypnoglow/oas2?status.svg)](https://godoc.org/github.com/hypnoglow/oas2)
[![CircleCI](https://circleci.com/gh/hypnoglow/oas2.svg?style=shield)](https://circleci.com/gh/hypnoglow/oas2)
[![codecov](https://codecov.io/gh/hypnoglow/oas2/branch/master/graph/badge.svg)](https://codecov.io/gh/hypnoglow/oas2)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypnoglow/oas2)](https://goreportcard.com/report/github.com/hypnoglow/oas2)
[![GitHub release](https://img.shields.io/github/tag/hypnoglow/oas2.svg)](https://github.com/hypnoglow/oas2/releases)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

**Note that this is not stable yet. In accordance with semantic versioning, the API can change between any minor versions. Use a vendoring tool of your
preference to lock an exact [release](https://github.com/hypnoglow/oas2/releases) version.**

Package oas2 provides utilities for building APIs using the OpenAPI 2.0 
specification (aka Swagger) in Go idiomatic way on top of `net/http`.

You don't need to learn any special framework or write `net/http`-incompatible
code - just delegate request validation, request parameters decoding
and other routines to this library - and focus on your application logic.

This package is built on top of [OpenAPI Initiative golang toolkit](https://github.com/go-openapi).

### Should I have an OpenAPI specification for my API?

If you don't have a spec for your API yet - it's definitely worth it to create 
one. The specification itself provides many useful things, such as documentation,
usage examples, and others. [Learn more](https://www.openapis.org/) about OpenAPI
and its purposes. The great thing is that it is compatible with many tools for 
developers and consumers; [Swagger Toolkit](https://swagger.io/) is the most popular
set of utilities for OpenAPI.

This package offers an integration of the spec with your code. And tightly 
coupling your code with the spec is a good thing - you create a strong contract
for API consumers, and any changes to your API will be clearly reflected in the 
spec. You will see many benefits, such as distinctly recognize the situation when
you need to increase the major version of your API because of incompatible changes.

## Features

### Router from a spec

This package provides an easy way to automatically create a router supporting
all resources from your OpenAPI specification file. The underlying router is only
your choice - you can use [gorilla/mux](https://github.com/gorilla/mux), [chi](https://github.com/go-chi/chi)
or any other.

Let's dive into a simple example.

Given a spec: [petstore.yaml](_examples/petstore.yaml)

First of all, load your spec in your app (note that though package import path ends in `oas2`, the package namespace is actually `oas`):

```go
import "github.com/hypnoglow/oas2"

// ...

// specPath is a path to your spec file.
doc, _ := oas.LoadFile(specPath)
```

Next, create an [operation](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#operationObject) handler. 
Let's define a handler for `findPetsByStatus` operation:

```go
type FindPetsByStatusHandler struct {
	storage PetStorage
}

func (h FindPetsByStatusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	statuses := req.URL.Query()["status"]

	pets := h.storage.FindByStatus(statuses)

	_ = json.NewEncoder(w).Encode(pets)
}
```

```go
handlers := oas.OperationHandlers{
    "findPetsByStatus":    findPetsByStatus{},
}
```

Define what options (logger, middleware) you will use:

```go
logger := logrus.New()
logger.SetLevel(logrus.DebugLevel)
queryValidator := oas.QueryValidator(errHandler)
```

Create a router:

```go
router, _ := oas.NewRouter(
    doc, 
    handlers, 
    oas.DebugLog(logger.Debugf), 
    oas.Use(queryValidator)
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

See the full [example](_examples/router/main.go) for the complete code.

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
oas.DecodeQuery(req, &m)

fmt.Printf("%#v", m) // Member{Name:"John", Age:27, LovesApples:true}
```

Note that it works only with oas router, because it needs to extract operation
spec from the request. To use custom parameters spec, use `oas.DecodeQueryParams()`.
See [`godoc example`](https://godoc.org/github.com/hypnoglow/oas2#example-DecodeQueryParams) for details.

### Pluggable formats & validators

The specification [allows](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types) to have custom formats and to validate against them.

This package provides the following custom formats and validators:
- [`partial-time`](formats/partial_time.go)

You can also implement your custom format and validator for it, and then register it:
```go
validate.RegisterFormat("myformat", &MyCustomFormat{}, ValidateMyCustomFormat)
```

## License

[MIT](https://github.com/hypnoglow/oas2/blob/master/LICENSE).
