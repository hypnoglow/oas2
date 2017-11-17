package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hypnoglow/oas2/validate"
)

// NewBodyValidator returns new Middleware that validates request body
// against parameters defined in OpenAPI 2.0 spec.
func NewBodyValidator(errHandler RequestErrorHandler) Middleware {
	return bodyValidatorMiddleware{
		errHandler: errHandler,
	}
}

type bodyValidatorMiddleware struct {
	errHandler RequestErrorHandler
}

func (m bodyValidatorMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// TODO
		if req.Header.Get("Content-Type") != "application/json" {
			// Do not validate multipart/form.
			// There will be built-in validation in oas2 package,
			// but currently it's cumbersome to implement.
			next.ServeHTTP(w, req)
			return
		}

		if req.Body == http.NoBody {
			next.ServeHTTP(w, req)
			return
		}

		op := GetOperation(req)
		if op == nil {
			next.ServeHTTP(w, req)
			return
		}

		// Read req.Body using io.TeeReader, so it can be read again
		// in the actual request handler.

		var b bytes.Buffer
		tr := io.TeeReader(req.Body, &b)
		defer req.Body.Close()

		var body interface{}
		if err := json.NewDecoder(tr).Decode(&body); err != nil {
			err = JsonError{error: fmt.Errorf("json decode: %s", err)}
			if !m.errHandler(w, req, err) {
				return
			}
		}

		// Validate body
		if errs := validate.Body(op.Parameters, body); len(errs) > 0 {
			err := ValidationError{error: fmt.Errorf("validation error"), errs: errs}
			if !m.errHandler(w, req, err) {
				return
			}
		}

		// Replace the body so it can be read again.
		req.Body = ioutil.NopCloser(&b)

		next.ServeHTTP(w, req)
	})
}
