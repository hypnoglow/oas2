package oas

import (
	"context"
	"net/http"

	"github.com/hypnoglow/oas2/convert"
)

// PathParameterExtractor returns new Middleware that extracts parameters
// defined in OpenAPI 2.0 spec as path parameters from path.
func PathParameterExtractor(extractor func(r *http.Request, key string) string) Middleware {
	return pathParameterExtractor{extractor}.chain
}

type pathParameterExtractor struct {
	extractor func(r *http.Request, key string) string
}

func (m pathParameterExtractor) chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// It's better to panic than to silently skip validation.
		params := MustParams(req)

		for _, p := range params {
			if p.In != "path" {
				continue
			}

			value, err := convert.Primitive(m.extractor(req, p.Name), p.Type, p.Format)
			if err == nil {
				req = WithPathParam(req, p.Name, value)
			}
		}

		next.ServeHTTP(w, req)
	})
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
