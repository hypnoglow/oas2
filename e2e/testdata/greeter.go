package testdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"testing"

	"github.com/hypnoglow/oas2"
)

// GreetHandler is a simple handler that greets using a name.
type GreetHandler struct{}

// ServeHTTP implements http.Handler.
func (GreetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var query struct {
		Name string `oas:"name"`
	}
	if err := oas.DecodeQuery(req, &query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, `{"greeting":"Hello, %s!"}`, query.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TestGreeter(t *testing.T, srv *httptest.Server) {
	t.Helper()

	resp, err := srv.Client().Get(srv.URL + "/api/greeting?name=Foo")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected response, status=%s, body=%q", resp.Status, string(b))
	}

	expected := `{"greeting":"Hello, Foo!"}`
	if string(b) != expected {
		t.Fatalf("Expected %q but got %q", expected, string(b))
	}
}

// GreeterSpec returns an OpenAPI spec for greeter server.
func GreeterSpec(t *testing.T) *oas.Document {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	doc, err := oas.LoadFile(path.Join(path.Dir(filename), "greeter.yaml"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	return doc
}
