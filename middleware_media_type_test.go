package oas

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestContentTypeValidator(t *testing.T) {
	testCases := map[string]struct {
		consumes            []string
		expectHandlerCalled bool
		expectedStatus      int
	}{
		"consumes application/json": {
			consumes: []string{
				"application/json",
			},
			expectHandlerCalled: true,
			expectedStatus:      http.StatusOK,
		},
		"consumes application/xml": {
			consumes: []string{
				"application/xml",
			},
			expectHandlerCalled: false,
			expectedStatus:      http.StatusUnsupportedMediaType,
		},
	}

	// TODO: test accept
	var produces []string

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := &fakeHandler{}
			v := &requestContentTypeValidator{next: h}

			w := httptest.NewRecorder()
			v.ServeHTTP(w, newRequest(nil), tc.consumes, produces, true)

			assert.Equal(t, tc.expectHandlerCalled, h.called)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestResponseContentTypeValidator(t *testing.T) {
	testCases := map[string]struct {
		accept         []string
		produces       []string
		expectedErrors int
	}{
		"accept and produces application/json": {
			accept:         []string{"application/json"},
			produces:       []string{"application/json"},
			expectedErrors: 0,
		},
		"accept application/xml": {
			accept:         []string{"application/xml"},
			expectedErrors: 1,
		},
		"produces application/xml": {
			produces:       []string{"application/xml"},
			expectedErrors: 1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var errs []error
			errHandler := func(problem Problem) {
				errs = append(errs, problem.Cause())
			}

			h := &fakeHandler{}
			v := &responseContentTypeValidator{
				next:           h,
				problemHandler: ProblemHandlerFunc(errHandler),
			}

			w := httptest.NewRecorder()
			req := newRequest(tc.accept)
			v.ServeHTTP(w, req, tc.produces, true)

			assert.Len(t, errs, tc.expectedErrors)
		})
	}
}

func newRequest(accept []string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/foo", bytes.NewBufferString(`
{
	"foo": "bar"
}
`))
	req.Header.Set("Content-Type", "application/json")
	if len(accept) > 0 {
		req.Header["Accept"] = accept
	}
	return req
}

type fakeHandler struct {
	called bool
}

func (h *fakeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.called = true
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"foo":"bar"}`))
}
