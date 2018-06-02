package testdata

import (
	"fmt"
	"net/http"
	"path"
	"runtime"
	"testing"

	"github.com/go-openapi/loads"

	"github.com/hypnoglow/oas2"
)

type GreetHandler struct{}

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

func GreeterSpec(t *testing.T) *loads.Document {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	doc, err := oas.LoadFile(path.Join(path.Dir(filename), "greeter.yaml"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	return doc
}
