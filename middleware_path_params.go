package oas

import (
	"context"
	"net/http"

	"github.com/go-openapi/spec"

	"github.com/hypnoglow/oas2/convert"
)

// PathParamExtractorFunc is a function that extracts path parameters by key
// from the request.
type PathParamExtractorFunc func(req *http.Request, key string) string

// PathParam implements PathParamExtractor.
func (f PathParamExtractorFunc) PathParam(req *http.Request, key string) string {
	return f(req, key)
}

// PathParamExtractor can extract path parameters by key from the request.
type PathParamExtractor interface {
	PathParam(req *http.Request, key string) string
}

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

// pathParamExtractor is a middleware that extracts parameters
// defined in OpenAPI 2.0 spec as path parameters from path and adds
// them to the request context.
type pathParamExtractor struct {
	next http.Handler

	extractor PathParamExtractor
}

func (mw *pathParamExtractor) ServeHTTP(w http.ResponseWriter, req *http.Request, params []spec.Parameter, ok bool) {
	if !ok {
		mw.next.ServeHTTP(w, req)
		return
	}

	for _, p := range params {
		if p.In != "path" {
			continue
		}

		value, err := convert.Primitive(mw.extractor.PathParam(req, p.Name), p.Type, p.Format)
		if err == nil {
			req = WithPathParam(req, p.Name, value)
		}
	}

	mw.next.ServeHTTP(w, req)
}
