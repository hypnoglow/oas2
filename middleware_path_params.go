package oas

import (
	"context"
	"net/http"

	"github.com/go-openapi/spec"

	"github.com/hypnoglow/oas2/convert"
)

// GetPathParam returns a path parameter by name from a request.
// For example, a handler defined on a path "/pet/{id}" gets a request with
// path "/pet/12" - in this case GetPathParam(req, "id") returns 12.
func GetPathParam(req *http.Request, name string) interface{} {
	return req.Context().Value(contextKeyPathParam(name))
}

// WithPathParam returns request with context value defining path parameter name
// set to value.
func WithPathParam(req *http.Request, name string, value interface{}) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), contextKeyPathParam(name), value))
}

type contextKeyPathParam string

// pathParamsExtractor is a middleware that extracts parameters
// defined in OpenAPI 2.0 spec as path parameters from path and adds
// them to the request context.
type pathParamsExtractor struct {
	next http.Handler

	extractor func(req *http.Request, key string) string
}

func (mw *pathParamsExtractor) ServeHTTP(w http.ResponseWriter, req *http.Request, params []spec.Parameter, ok bool) {
	if !ok {
		mw.next.ServeHTTP(w, req)
		return
	}

	for _, p := range params {
		if p.In != "path" {
			continue
		}

		value, err := convert.Primitive(mw.extractor(req, p.Name), p.Type, p.Format)
		if err == nil {
			req = WithPathParam(req, p.Name, value)
		}
	}

	mw.next.ServeHTTP(w, req)
}
