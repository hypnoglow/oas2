// +build e2e

package multiple_specs

import (
	"net/http/httptest"
	"testing"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/e2e/testdata"
)

// TestMultipleSpecs tests that router properly works with multiple specs
// assigned to it.
func TestMultipleSpecs(t *testing.T) {
	greeterDoc := testdata.GreeterSpec(t)
	greeterHandlers := oas.OperationHandlers{
		"greet": testdata.GreetHandler{},
	}

	router := oas.NewRouter(
		oas.RouterMiddleware(oas.QueryValidator(nil)),
		oas.ServeSpec(oas.SpecHandlerTypeDynamic),
	)
	err := router.AddSpec(greeterDoc, greeterHandlers)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	adderDoc := testdata.AdderSpec(t)
	adderHandlers := oas.OperationHandlers{
		"add": testdata.AddHandler{},
	}

	if err := router.AddSpec(adderDoc, adderHandlers); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	srv := httptest.NewServer(router)
	defer srv.Close()

	testdata.TestGreeter(t, srv)
	testdata.TestAdder(t, srv)
}
