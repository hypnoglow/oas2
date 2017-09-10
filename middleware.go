package oas2

import (
	"net/http"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
)

// MiddlewareFn describes middleware function.
type MiddlewareFn func(next http.Handler) http.Handler

// Middleware describes a middleware that can be applied to a http.handler.
type Middleware interface {
	Apply(next http.Handler) http.Handler
}

// NewQueryValidator returns new Middleware that validates request query
// parameters against OpenAPI 2.0 spec.
func NewQueryValidator(sp *spec.Swagger, errHandler func(w http.ResponseWriter, errs []error)) Middleware {
	return queryValidatorMiddleware{
		sp:              sp,
		an:              analysis.New(sp),
		errHandler:      errHandler,
		continueOnError: false, // TODO: make controllable
	}
}

type queryValidatorMiddleware struct {
	sp              *spec.Swagger
	an              *analysis.Spec
	errHandler      func(w http.ResponseWriter, errs []error)
	continueOnError bool
}

func (m queryValidatorMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		opID := GetOperationID(req)
		if opID == "" {
			next.ServeHTTP(w, req)
			return
		}

		_, _, op, ok := m.an.OperationForName(opID.String())
		if !ok {
			next.ServeHTTP(w, req)
			return
		}

		if errs := ValidateQuery(op.Parameters, req.URL.Query()); len(errs) > 0 {
			m.errHandler(w, errs)
			if !m.continueOnError {
				return
			}
		}

		next.ServeHTTP(w, req)
	})
}
