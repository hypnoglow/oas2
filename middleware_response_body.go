package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-openapi/spec"

	"github.com/hypnoglow/oas2/validate"
)

// responseBodyValidator is a middleware that validates response body by OpenAPI
// operation definition.
type responseBodyValidator struct {
	next http.Handler

	// jsonSelectors represent content-type selectors. If any selector
	// matches content-type of the response, then response body will be validated.
	// Otherwise no validation is performed.
	jsonSelectors []*regexp.Regexp

	problemHandler ProblemHandler
}

func (mw *responseBodyValidator) ServeHTTP(w http.ResponseWriter, req *http.Request, responses *spec.Responses, ok bool) {
	if !ok {
		mw.next.ServeHTTP(w, req)
		return
	}

	respBuf := &bytes.Buffer{}
	rr := newWrapResponseWriter(w, 1)
	rr.Tee(respBuf)

	mw.next.ServeHTTP(rr, req)

	// First of all, check if response is defined for the status code.
	responseSpec, ok := responses.StatusCodeResponses[rr.Status()]
	if !ok {
		// If no response is explicitly defined for the status code, consider it
		// is ok.
		//
		// Quote from OpenAPI 2.0 spec:
		// > It is not expected from the documentation to necessarily cover all
		// > possible HTTP response codes, since they may not be known in advance.
		return
	}

	if responseSpec.Schema == nil {
		// This may be ok for example for HTTP 204 responses, but any response
		// with a body should explicitly define a schema.
		//
		// Quote from OpenAPI 2.0 spec:
		// > If this field does not exist, it means no content is returned as
		// > part of the response.
		if respBuf.Len() > 0 {
			e := fmt.Errorf("response has non-emtpy body, but the operation does not define response schema for code %d", rr.Status())
			mw.problemHandler.HandleProblem(NewProblem(w, req, e))
		}
		return
	}

	// Check the content type of the response. If it does not match any selector,
	// don't validate the response.
	if !mw.matchContentType(rr.Header()) {
		return
	}

	var body interface{}
	if err := json.NewDecoder(respBuf).Decode(&body); err != nil {
		e := fmt.Errorf("response body contains invalid json: %s", err)
		mw.problemHandler.HandleProblem(NewProblem(w, req, e))
		return
	}

	if errs := validate.BySchema(responseSpec.Schema, body); len(errs) > 0 {
		me := newMultiError("response body does not match the schema", errs...)
		mw.problemHandler.HandleProblem(NewProblem(w, req, me))
		return
	}
}

// matchContentType checks if content type of the request matches any selector.
func (mw *responseBodyValidator) matchContentType(hdr http.Header) bool {
	contentType := hdr.Get("Content-Type")
	for _, selector := range mw.jsonSelectors {
		if selector.MatchString(contentType) {
			return true
		}
	}

	return false
}
