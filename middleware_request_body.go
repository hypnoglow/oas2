package oas

import (
	"bytes"
	//"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/go-openapi/spec"
	"github.com/json-iterator/go"

	"github.com/hypnoglow/oas2/validate"
)

// requestBodyValidator is a middleware that validates request body by OpenAPI
// operation definition.
type requestBodyValidator struct {
	next http.Handler

	// jsonSelectors represent content-type selectors. If any selector
	// matches content-type of the request, then request body will be validated.
	// Otherwise no validation is performed.
	jsonSelectors []*regexp.Regexp

	problemHandler    ProblemHandler
	continueOnProblem bool
}

func (mw *requestBodyValidator) ServeHTTP(w http.ResponseWriter, req *http.Request, params []spec.Parameter, ok bool) {
	if !ok {
		mw.next.ServeHTTP(w, req)
		return
	}

	if req.Body == http.NoBody {
		for _, param := range params {
			if param.In == "body" && param.Required {
				// No request body found, but operation actually requires body.
				e := fmt.Errorf("request body is empty, but the operation requires non-empty body")
				mw.problemHandler.HandleProblem(NewProblem(w, req, e))
				if !mw.continueOnProblem {
					return
				}
			}
		}
		mw.next.ServeHTTP(w, req)
		return
	}

	if !mw.matchContentType(req) {
		mw.next.ServeHTTP(w, req)
		return
	}

	// Read req.Body using io.TeeReader, so it can be read again
	// in the actual request handler.
	body, err := bodyPayload(req)
	if err != nil {
		e := fmt.Errorf("request body contains invalid json: %s", err)
		mw.problemHandler.HandleProblem(NewProblem(w, req, e))
		if !mw.continueOnProblem {
			return
		}
	}

	if errs := validate.Body(params, body); len(errs) > 0 {
		me := newMultiError("request body does not match the schema", errs...)
		mw.problemHandler.HandleProblem(NewProblem(w, req, me))
		if !mw.continueOnProblem {
			return
		}
	}

	mw.next.ServeHTTP(w, req)
}

// matchContentType checks if content type of the request matches any selector.
func (mw *requestBodyValidator) matchContentType(req *http.Request) bool {
	contentType := req.Header.Get("Content-Type")
	for _, selector := range mw.jsonSelectors {
		if selector.MatchString(contentType) {
			return true
		}
	}

	return false
}

// bodyPayload reads req.Body and returns it. Request body can be
// read again later.
func bodyPayload(req *http.Request) (interface{}, error) {
	buf := &bytes.Buffer{}
	tr := io.TeeReader(req.Body, buf)
	defer req.Body.Close()

	// OK
	b, err := ioutil.ReadAll(tr)
	if err != nil {
		return nil, err
	}

	var payload interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, err
	}

	req.Body = ioutil.NopCloser(buf)
	return payload, nil
}
