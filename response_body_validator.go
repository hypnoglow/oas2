package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	chimw "github.com/go-chi/chi/middleware"

	"github.com/hypnoglow/oas2/validate"
)

// ResponseBodyValidator returns new Middleware that validates response body
// against schema defined in OpenAPI 2.0 spec.
func ResponseBodyValidator(errHandler ResponseErrorHandler, opts ...MiddlewareOption) Middleware {
	options := MiddlewareOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	selectors := []*regexp.Regexp{contentTypeSelectorRegexJSON, contentTypeSelectorRegexJSONAPI}
	if options.contentTypeRegexSelectors != nil {
		selectors = options.contentTypeRegexSelectors
	}

	return responseBodyValidator{errHandler, selectors}.chain
}

type responseBodyValidator struct {
	errHandler ResponseErrorHandler
	selectors  []*regexp.Regexp
}

func (m responseBodyValidator) chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// It's better to panic than to silently skip validation.
		op := MustOperation(req)

		responseBodyBuffer := &bytes.Buffer{}
		rr := chimw.NewWrapResponseWriter(w, 1)
		rr.Tee(responseBodyBuffer)

		next.ServeHTTP(rr, req)

		// Check content type of the response. If it does not match any selector,
		// don't validate the request.
		contentTypeMatch := false
		contentType := w.Header().Get("Content-Type")
		for _, selector := range m.selectors {
			if selector.MatchString(contentType) {
				contentTypeMatch = true
				break
			}
		}

		if !contentTypeMatch {
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
