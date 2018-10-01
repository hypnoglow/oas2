package oas

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryValidator(t *testing.T) {
	testCases := map[string]struct {
		query          string
		expectedStatus int
		expectedBody   string
	}{
		"valid query": {
			query:          "username=johndoe&password=123",
			expectedStatus: http.StatusOK,
			expectedBody:   "username: johndoe, password: 123",
		},
		"missing required query parameter": {
			query:          "username=johndoe",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"message":"param password is required","field":"password"}]}`,
		},
	}

	doc := loadDocFile(t, "testdata/petstore_1.yml")
	params := doc.Analyzer.ParametersFor("loginUser")

	v := &queryValidator{
		next:              http.HandlerFunc(handleUserLogin),
		problemHandler:    problemHandlerResponseWriter(),
		continueOnProblem: false,
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v2/user/login?"+tc.query, nil)
			w := httptest.NewRecorder()
			v.ServeHTTP(w, req, params, true)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func handleUserLogin(w http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	password := req.URL.Query().Get("password")

	// Never do this! This is just for testing purposes.
	fmt.Fprintf(w, "username: %s, password: %s", username, password)
}
