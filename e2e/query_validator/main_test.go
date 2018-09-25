// +build e2e

package query_validator

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/adapter/gorilla"
	"github.com/hypnoglow/oas2/e2e/testdata"
)

// TestQueryValidatorMiddleware tests that router created with query validator
// middleware will validate query.
func TestQueryValidatorMiddleware(t *testing.T) {
	doc := testdata.GreeterSpec(t)
	basis := oas.NewResolvingBasis(doc, gorilla.NewResolver(doc))

	r := mux.NewRouter()
	err := gorilla.NewOperationRouter(r).
		WithDocument(doc).
		WithOperationHandlers(map[string]http.Handler{
			"greet": testdata.GreetHandler{},
		}).
		WithMiddleware(
			basis.OperationContext(),
			basis.QueryValidator(
				oas.WithProblemHandler(oas.ProblemHandlerFunc(handleValidationError)),
			),
		).
		Route()
	assert.NoError(t, err)

	srv := httptest.NewServer(r)
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

func handleValidationError(p oas.Problem) {
	p.ResponseWriter().WriteHeader(http.StatusBadRequest)
	io.WriteString(p.ResponseWriter(), p.Cause().Error())
}
