package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/hypnoglow/oas2/validate"
)

// BodyValidator returns new Middleware that validates request body
// against parameters defined in OpenAPI 2.0 spec.
func BodyValidator(errHandler RequestErrorHandler, opts ...MiddlewareOption) Middleware {
	options := MiddlewareOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	selectors := []*regexp.Regexp{contentTypeSelectorRegexJSON, contentTypeSelectorRegexJSONAPI}
	if options.contentTypeRegexSelectors != nil {
		selectors = options.contentTypeRegexSelectors
	}

	return bodyValidatorMiddleware{errHandler: errHandler, selectors: selectors}.chain
}

type bodyValidatorMiddleware struct {
	errHandler RequestErrorHandler
	selectors  []*regexp.Regexp
}

func (m bodyValidatorMiddleware) chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Check content type of the request. If it does not match any selector,
		// don't validate the request.
		contentTypeMatch := false
		contentType := req.Header.Get("Content-Type")
		for _, selector := range m.selectors {
			if selector.MatchString(contentType) {
				contentTypeMatch = true
				break
			}
		}

		if !contentTypeMatch {
			next.ServeHTTP(w, req)
			return
		}

		if req.Body == http.NoBody {
			next.ServeHTTP(w, req)
			return
		}

		// It's better to panic than to silently skip validation.
		op := MustOperation(req)

		// Read req.Body using io.TeeReader, so it can be read again
		// in the actual request handler.

		var b bytes.Buffer
		tr := io.TeeReader(req.Body, &b)
		defer req.Body.Close()

		var body interface{}
		if err := json.NewDecoder(tr).Decode(&body); err != nil {
			err = JSONError{error: fmt.Errorf("json decode: %s", err)}
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
