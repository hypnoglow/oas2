package oas

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathParamsExtractor(t *testing.T) {
	testCases := map[string]struct {
		url            string
		extractor      func(req *http.Request, key string) string
		expectedStatus int
		expectedBody   string
	}{
		"extracts parameters": {
			url: "/v2/pet/12",
			extractor: func(req *http.Request, key string) string {
				return "12"
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "pet by id: 12",
		},
	}

	doc := loadDocFile(t, "testdata/petstore_1.yml")
	params := doc.Analyzer.ParametersFor("getPetById")

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := &pathParamsExtractor{
				next:      http.HandlerFunc(handleGetPetByID),
				extractor: tc.extractor,
			}

			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req, params, true)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func handleGetPetByID(w http.ResponseWriter, req *http.Request) {
	id, ok := GetPathParam(req, "petId").(int64)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "pet by id: %d", id)
}
