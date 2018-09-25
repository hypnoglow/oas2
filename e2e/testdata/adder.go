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

// AddHandler is a simple handler that sums two numbers.
type AddHandler struct{}

// ServeHTTP implements http.Handler.
func (AddHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var query struct {
		A int64 `oas:"a"`
		B int64 `oas:"b"`
	}
	if err := oas.DecodeQuery(req, &query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, `{"sum":%d}`, query.A+query.B); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TestAdder tests adder server.
func TestAdder(t *testing.T, srv *httptest.Server) {
	t.Helper()

	resp, err := srv.Client().Get(srv.URL + "/api/adder/sum?a=1&b=2")
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

	expected := `{"sum":3}`
	if string(b) != expected {
		t.Fatalf("Expected %q but got %q", expected, string(b))
	}
}

// AdderSpec returns an OpenAPI spec for adder server.
func AdderSpec(t *testing.T) *oas.Document {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	doc, err := oas.LoadFile(path.Join(path.Dir(filename), "adder.yaml"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	return doc
}
