// +build e2e

package query_validator

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/e2e/testdata"
)

// TestQueryValidatorMiddleware tests that router created with query validator
// middleware will validate query.
func TestQueryValidatorMiddleware(t *testing.T) {
	doc := testdata.GreeterSpec(t)

	handlers := oas.OperationHandlers{
		"greet": testdata.GreetHandler{},
	}

	router := oas.NewRouter(
		oas.RouterMiddleware(oas.QueryValidator(handleValidationError)),
	)
	err := router.AddSpec(doc, handlers)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	srv := httptest.NewServer(router)
	defer srv.Close()

	// We expect to get 400 error on empty name.
	resp, err := srv.Client().Get(srv.URL + "/api/greeting?name=")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Unexpected response status: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !strings.Contains(string(b), "name in query is required") {
		t.Fatalf("Unexpected response body")
	}
}

func handleValidationError(w http.ResponseWriter, req *http.Request, err error) (resume bool) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, err.Error())
	return false
}
