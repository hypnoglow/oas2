package oas

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseBodyValidator(t *testing.T) {
	testCases := map[string]struct {
		url               string
		expectedStatus    int
		expectedBody      string
		expectedLogBuffer string
	}{
		"logs validation error": {
			url:               "/v2/pet/12",
			expectedStatus:    http.StatusOK,
			expectedBody:      `{"id":123,"name":"Kitty"}`,
			expectedLogBuffer: "problem handler: response body does not match the schema: field=age value=<nil> message=age in body is required",
		},
		"no logs when no response spec defined": {
			url:               "/v2/pet/500",
			expectedStatus:    http.StatusInternalServerError,
			expectedBody:      `{"error":"foo"}`,
			expectedLogBuffer: "",
		},
		"logs validation warning when no schema defined": {
			url:               "/v2/pet/404",
			expectedStatus:    http.StatusNotFound,
			expectedBody:      `{"error":"not found"}`,
			expectedLogBuffer: "problem handler: response has non-emtpy body, but the operation does not define response schema for code 404",
		},
		"logs validation error when response body is bad json": {
			url:               "/v2/pet/badjson",
			expectedStatus:    http.StatusOK,
			expectedBody:      `{"name":`,
			expectedLogBuffer: "problem handler: response body contains invalid json: unexpected EOF",
		},
	}

	doc := loadDocFile(t, "testdata/petstore_1.yml")
	_, _, op, ok := doc.Analyzer.OperationForName("getPetById")
	assert.True(t, ok)

	logBuffer := &bytes.Buffer{}

	v := &responseBodyValidator{
		next:           http.HandlerFunc(handleGetPetByIDFaked),
		jsonSelectors:  []*regexp.Regexp{contentTypeSelectorRegexJSON},
		problemHandler: problemHandlerBufferLogger(logBuffer),
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			logBuffer.Reset()
			defer logBuffer.Reset()

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			v.ServeHTTP(w, req, op.Responses, true)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expectedLogBuffer, strings.TrimSpace(logBuffer.String()))
		})
	}
}

func handleGetPetByIDFaked(w http.ResponseWriter, req *http.Request) {
	// set Content-Type for all responses to ensure validator does not filter
	// them out prior to other checks like response spec schema presence.
	w.Header().Set("Content-Type", "application/json")

	// fake not found
	if req.URL.Path == "/v2/pet/404" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))  // nolint: errcheck
		return
	}

	// fake for server error {
	if req.URL.Path == "/v2/pet/500" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"foo"}`))  // nolint: errcheck
		return
	}

	// fake for bad json
	if req.URL.Path == "/v2/pet/badjson" {
		w.Write([]byte(`{"name":`))  // nolint: errcheck
		return
	}

	// normal

	type pet struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	p := pet{123, "Kitty"}

	err := json.NewEncoder(w).Encode(p)
	assertNoError(err)
}
