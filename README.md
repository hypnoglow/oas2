# oas2

[![GoDoc](https://godoc.org/github.com/hypnoglow/oas2?status.svg)](https://godoc.org/github.com/hypnoglow/oas2)
[![CircleCI](https://circleci.com/gh/hypnoglow/oas2.svg?style=shield)](https://circleci.com/gh/hypnoglow/oas2)
[![codecov](https://codecov.io/gh/hypnoglow/oas2/branch/master/graph/badge.svg)](https://codecov.io/gh/hypnoglow/oas2)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypnoglow/oas2)](https://goreportcard.com/report/github.com/hypnoglow/oas2)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

**WIP. Not stable. API may change at any time.**

Package oas2 provides utilities to work with OpenAPI 2.0 specification
(aka Swagger).

The purpose of this package is to provide utilities for building APIs
from the OpenAPI specification in Go idiomatic way on top of `net/http`.

You don't need to learn any special framework or write `net/http`-incompatible
code - just delegate request validation, request parameters decoding
and other routines to this library - and focus on your application logic.

## Features

#### Decode query parameters to a struct

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

See [`examples/`](https://github.com/hypnoglow/oas2/tree/master/examples) directory for complete examples.

## License

[MIT](https://github.com/hypnoglow/oas2/blob/master/LICENSE).
