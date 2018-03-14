# Validate

Validate is a CLI tool for validating OAS files. For proper validation all
references from original file are expanded. This process can take some time,
so if you don't want to wait on each start you can skip validation in your
code and perform it as a one-time action (for example on docker build)
using this tool.

Install

```sh
go get -u github.com/hypnoglow/oas2/cmd/oas-validate
```

Run to validate specification file

```sh
oas-expand petstore.yaml
```

You can easily see the difference

```go
now := time.Now()
oas.LoadSpec("./petstore.yaml")
log.Printf("Spec parsed and validated in %s\n", time.Since(now))

now = time.Now()
oas.LoadSpec("./petstore.yaml", Validation(false))
log.Printf("Spec parsed in %s\n", time.Since(now))
```

```
2018/02/13 18:47:08 Spec parsed and validated in 1.224684788s
2018/02/13 18:47:08 Spec parsed in 39.403587ms
```