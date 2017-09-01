# oas2

**WIP. Not stable. API may change at any time.**

Package oas2 provides utilities to work with OpenAPI v2 specification
(aka Swagger).

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
oas.DecodeQuery(paramSpec, req.URL.Query(), &m)

fmt.Printf("%#v", m) // Member{Name:"John", Age:27, LovesApples:true}
```

See [`examples/`](https://github.com/hypnoglow/oas2/tree/master/examples) directory for complete examples.

## License

[MIT](https://github.com/hypnoglow/oas2/blob/master/LICENSE).
