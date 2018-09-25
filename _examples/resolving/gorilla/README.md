# Example: Resolving basis

This example demonstrates the use of a resolving basis
with gorilla/mux as an underlying router implementation.

## Run

Run the server:

```
go run main.go -spec ../../app/openapi.yaml
```

Send a correct request:

```
$ curl -XPOST -i --url 'http://localhost:8080/api/v1/sum' \
    -H "Content-Type: application/json; charset=utf-8" \
    -d '{"number":2}'
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 25 Sep 2018 18:54:36 GMT
Content-Length: 10

{"sum":2}
```

Now try a request that does not meet spec parameters requirements:

```
$ curl -XPOST -i --url 'http://localhost:8080/api/v1/sum' \
    -H "Content-Type: application/json; charset=utf-8" \
    -d '{"number":"foo"}'
HTTP/1.1 400 Bad Request
Content-Type: application/json; charset=utf-8
Date: Tue, 25 Sep 2018 18:55:00 GMT
Content-Length: 106

{"errors":["request body does not match the schema: number in body must be of type integer: \"string\""]}
```
