package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	chimw "github.com/go-chi/chi/middleware"

	"github.com/hypnoglow/oas2/validate"
)

// ResponseBodyValidator returns new Middleware that validates response body
// against schema defined in OpenAPI 2.0 spec.
func ResponseBodyValidator(errHandler ResponseErrorHandler) Middleware {
	return responseBodyValidator{errHandler}.chain
}

type responseBodyValidator struct {
	errHandler ResponseErrorHandler
}

func (m responseBodyValidator) chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		op := GetOperation(req)
		if op == nil {
			next.ServeHTTP(w, req)
			return
		}

		responseBodyBuffer := &bytes.Buffer{}
		rr := chimw.NewWrapResponseWriter(w, 1)
		rr.Tee(responseBodyBuffer)

		next.ServeHTTP(rr, req)

		// Only json body can be validated currently.
		if w.Header().Get("Content-Type") != "application/json" {
			return
		}

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
		if err := json.NewDecoder(responseBodyBuffer).Decode(&body); err != nil {
			err = JSONError{error: fmt.Errorf("json decode: %s", err)}
			m.errHandler(w, req, err)
			return
		}

		if errs := validate.BySchema(responseSpec.Schema, body); len(errs) > 0 {
			err := ValidationError{error: fmt.Errorf("validation error"), errs: errs}
			m.errHandler(w, req, err)
		}
	})
}
