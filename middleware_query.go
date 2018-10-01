package oas

import (
	"net/http"

	"github.com/go-openapi/spec"

	"github.com/hypnoglow/oas2/validate"
)

// queryValidator is a middleware that validates request query by OpenAPI operation
// definition.
type queryValidator struct {
	next http.Handler

	problemHandler    ProblemHandler
	continueOnProblem bool
}

func (mw *queryValidator) ServeHTTP(w http.ResponseWriter, req *http.Request, params []spec.Parameter, ok bool) {
	if !ok {
		mw.next.ServeHTTP(w, req)
		return
	}

	if errs := validate.Query(params, req.URL.Query()); len(errs) > 0 {
		me := newMultiError("query params do not match the schema", errs...)
		mw.problemHandler.HandleProblem(NewProblem(w, req, me))
		if !mw.continueOnProblem {
			return
		}
	}

	mw.next.ServeHTTP(w, req)
}
