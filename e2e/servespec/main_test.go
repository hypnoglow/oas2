// +build e2e

package servespec

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/e2e/testdata"
)

// TestMiddlewareIsAppliedToServedSpec tests that middleware passed to the router
// with Wrap option is applied to the router's served spec.
func TestMiddlewareIsAppliedToServedSpec(t *testing.T) {
	doc := testdata.GreeterSpec(t)

	r := mux.NewRouter()
	r.Path("/api").
		Methods(http.MethodGet).
		Handler(oas.NewDynamicSpecHandler(doc))

	srv := httptest.NewServer(r)
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
}
