package oas_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hypnoglow/oas2"
)

func TestDecodeQueryParams(t *testing.T) {
	doc, err := oas.LoadFile(getSpecPath(t))
	assert.NoError(t, err)

	oas.RegisterAdapter("fake", fakeAdapter{})

	basis := oas.NewResolvingBasis("fake", doc)
	h := basis.OperationContext()(getPetsHandler{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/pets", nil)
	h.ServeHTTP(w, req)

	var result struct {
		Limit int64 `json:"limit"`
	}

	assert.Equal(t, http.StatusOK, w.Code)

	if err = json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assert.EqualValues(t, 10, result.Limit)
}

type fakeAdapter struct{}

func (fakeAdapter) Resolver(meta interface{}) oas.Resolver {
	return fakeResolver{}
}

func (fakeAdapter) OperationRouter(meta interface{}) oas.OperationRouter {
	panic("implement me")
}

func (fakeAdapter) PathParamExtractor() oas.PathParamExtractor {
	panic("implement me")
}

type fakeResolver struct{}

func (fakeResolver) Resolve(req *http.Request) (string, bool) {
	return "getPets", true
}

type getPetsHandler struct{}

func (getPetsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var query struct {
		Limit int64 `json:"limit" oas:"limit"`
	}

	if err := oas.DecodeQuery(req, &query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getSpecPath(t *testing.T) string {
	t.Helper()

	spec := `
swagger: "2.0"
info:
  title: Test API
  version: 0.1.0
basePath: "/api"
paths:
  /pets:
    get:
      summary: Find pets
      operationId: getPets
      parameters:
      - name: limit
        in: query
        required: false
        type: integer
        format: int64
        default: 10
      responses:
        200:
          description: Pets found
          schema:
            type: object
            properties: 
              limit:
                title: Requested limit
                type: integer
                format: int64
        500:
          description: "Internal Server Error"
`
	p := "/tmp/spec.yaml"

	// write to file
	if err := ioutil.WriteFile(p, []byte(spec), 0755); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	return p
}
