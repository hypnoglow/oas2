# Expand

Expand is a CLI tool for expanding OAS files. It makes new specification file
with all references expanded. Loading this new file is up to 100 times faster
than loading regular (non-expanded) spec. So this may be used as a one-time
action (for example on building docker image) to reduce time of all futher
application starts.

Install

```sh
go get -u github.com/hypnoglow/oas2/cmd/oas-expand
```

Run to make specification file expanded

```sh
oas-expand -target-dir=./cache spec.yaml
```

Cache directory can be passed to oas loader using `CacheDir` parameter

```go
doc, err := oas.LoadSpec(specPath, CacheDir("./cache"))
```

You can easily see the difference

```go
now := time.Now()
oas.LoadSpec("./petstore.yaml")
log.Printf("Spec parsed in %s\n", time.Since(now))

now = time.Now()
oas.LoadSpec("./petstore.yaml", oas.CacheDir("./cache"))
log.Printf("Expanded spec parsed in %s\n", time.Since(now))
```

```
2018/02/13 18:47:08 Spec parsed in 1.224684788s
2018/02/13 18:47:08 Expanded spec parsed in 39.403587ms
```