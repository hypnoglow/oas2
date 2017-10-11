package oas2

import (
	"net/http"

	"github.com/hypnoglow/oas2/validate"
)

// NewQueryValidator returns new Middleware that validates request query
// parameters against OpenAPI 2.0 spec.
func NewQueryValidator(errHandler func(w http.ResponseWriter, errs []error)) Middleware {
	return queryValidatorMiddleware{
		errHandler:      errHandler,
		continueOnError: false, // TODO: make controllable
	}
}

type queryValidatorMiddleware struct {
	errHandler      func(w http.ResponseWriter, errs []error)
	continueOnError bool
}

func (m queryValidatorMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		op := GetOperation(req)
		if op == nil {
			next.ServeHTTP(w, req)
			return
		}

		if errs := validate.Query(op.Parameters, req.URL.Query()); len(errs) > 0 {
			m.errHandler(w, errs)
			if !m.continueOnError {
				return
			}
		}

		next.ServeHTTP(w, req)
	})
}
