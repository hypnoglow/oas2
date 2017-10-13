package oas2

import (
	"encoding/json"
	"net/http"

	"github.com/hypnoglow/oas2/utils"
	"github.com/hypnoglow/oas2/validate"
)

// NewResponseBodyValidator returns new Middleware that validates response body
// against schema defined in OpenAPI 2.0 spec.
func NewResponseBodyValidator(errHandler func(w http.ResponseWriter, errs []error)) Middleware {
	return responseBodyValidator{errHandler}
}

type responseBodyValidator struct {
	errHandler func(w http.ResponseWriter, errs []error)
}

func (m responseBodyValidator) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		op := GetOperation(req)
		if op == nil {
			next.ServeHTTP(w, req)
			return
		}

		rr := utils.NewResponseRecorder(w)

		next.ServeHTTP(rr, req)

		responseSpec, ok := op.Responses.StatusCodeResponses[rr.Status()]
		if !ok {
			// TODO: should notify package user that there is no response spec.
			return
		}

		if responseSpec.Schema == nil {
			// It may be ok for responses like 204.
			return
		}

		var body interface{}
		if err := json.Unmarshal(rr.Payload(), &body); err != nil {
			// TODO: should notify package user about the error.
			return
		}

		if errs := validate.BySchema(responseSpec.Schema, body); len(errs) > 0 {
			m.errHandler(w, errs)
		}
	})
}
