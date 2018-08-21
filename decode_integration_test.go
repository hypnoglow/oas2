package oas_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hypnoglow/oas2"
)

func TestDecodeQueryParams(t *testing.T) {
	doc, err := oas.LoadFile(getSpecPath(t))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	handlers := oas.OperationHandlers{
		"getPets": getPetsHandler{},
	}

	r := oas.NewRouter()
	if err = r.AddSpec(doc, handlers); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	srv := httptest.NewServer(r)

	u := srv.URL + "/api/pets"
	resp, err := http.Get(u)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Limit int64 `json:"limit"`
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status code to be 200 but got %v", resp.StatusCode)
		if b, err := ioutil.ReadAll(resp.Body); err == nil {
			t.Errorf("Response body: %v", string(b))
		}
		t.Fatalf("Request test failed")
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Limit != 10 {
		t.Fatalf("Expected limit to be 10 but got %v", result.Limit)
	}
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
