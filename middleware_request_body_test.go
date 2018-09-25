package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBodyValidator(t *testing.T) {
	testCases := map[string]struct {
		contentType    string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		"valid json body": {
			contentType:    "application/json",
			body:           `{"name":"johndoe","age":7}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "pet name: johndoe",
		},
		"required field \"name\" is missing": {
			contentType:    "application/json",
			body:           `{"age":7}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"message":"name in body is required","field":"name"}]}`,
		},
		"value for field \"age\" is incorrect": {
			contentType:    "application/json",
			body:           `{"name":"johndoe","age":"abc"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"message":"age in body must be of type integer: \"string\"","field":"age"}]}`,
		},
		"no body": {
			contentType:    "application/json",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"message":"request body is empty, but the operation requires non-empty body"}]}`,
		},
		"invalid json body": {
			contentType:    "application/json",
			body:           `{"name":"johndoe`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"message":"request body contains invalid json: unexpected EOF"}]}`,
		},
		"skip body validation for not application/json content type": {
			contentType:    "text/plain",
			body:           "some",
			expectedStatus: http.StatusUnsupportedMediaType, // returned from the actual handler
			expectedBody:   "",
		},
	}

	doc := loadDocFile(t, "testdata/petstore_1.yml")
	params := doc.Analyzer.ParametersFor("addPet")

	v := &requestBodyValidator{
		next:              http.HandlerFunc(handleAddPet),
		jsonSelectors:     []*regexp.Regexp{contentTypeSelectorRegexJSON},
		problemHandler:    problemHandlerResponseWriter(),
		continueOnProblem: false,
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var body io.Reader
			if tc.body != "" {
				body = bytes.NewBufferString(tc.body)
			}
			req := httptest.NewRequest(http.MethodPost, "/v2/pet", body)
			req.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()
			v.ServeHTTP(w, req, params, true)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func handleAddPet(w http.ResponseWriter, req *http.Request) {
	type pet struct {
		Name      string   `json:"name"`
		PhotoURLs []string `json:"photoUrls"`
	}

	var p pet
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	fmt.Fprintf(w, "pet name: %s", p.Name)
}
