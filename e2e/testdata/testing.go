package testdata

import (
	"fmt"
	"net/http"
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"greeting":"Hello, %s!"}`, query.Name)
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
