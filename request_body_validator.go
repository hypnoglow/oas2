package oas2

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
func NewBodyValidator(errHandler func(w http.ResponseWriter, errs []error)) Middleware {
	return bodyValidatorMiddleware{
		errHandler: errHandler,
	}
}

type bodyValidatorMiddleware struct {
	errHandler func(w http.ResponseWriter, errs []error)
}

func (m bodyValidatorMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
			m.errHandler(w, []error{fmt.Errorf("Body contains invalid json")})
			return
		}

		// Validate body
		if errs := validate.Body(op.Parameters, body); len(errs) > 0 {
			m.errHandler(w, errs)
			return
		}

		// Replace the body so it can be read again.
		req.Body = ioutil.NopCloser(&b)

		next.ServeHTTP(w, req)
	})
}
