package middleware_servespec

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/cors"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/e2e/testdata"
)

// TestMiddlewareIsAppliedToServedSpec tests that middleware passed to the router
// with Wrap option is applied to the router's served spec.
func TestMiddlewareIsAppliedToServedSpec(t *testing.T) {
	doc := testdata.GreeterSpec(t)

	handlers := oas.OperationHandlers{
		"greet": testdata.GreetHandler{},
	}

	router, err := oas.NewRouter(
		doc,
		handlers,
		oas.Wrap(cors.AllowAll().Handler),
		oas.ServeSpec(oas.SpecHandlerTypeDynamic),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	srv := httptest.NewServer(router)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/api")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected response status: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// a very simple check that spec is actually served
	if !strings.Contains(string(b), "swagger") {
		t.Fatalf("Expected reply to contain \"swagger\"")
	}

	if resp.Header.Get("Vary") != "Origin" {
		t.Fatalf("Expected response to contain Vary header with value Origin, but got %q", resp.Header.Get("Vary"))
	}
}
