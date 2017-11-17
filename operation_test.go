package oas

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

func ExampleID_String() {
	opID := OperationID("addPet")

	fmt.Fprint(os.Stdout, opID.String())

	// Output:
	// addPet
}

func TestMustOperation(t *testing.T) {
	t.Run("ok on request with operation", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/path", nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		op := &spec.Operation{OperationProps: spec.OperationProps{Description: "some"}}
		req = withOperation(req, op)

		actualOperation := MustOperation(req)
		if !reflect.DeepEqual(actualOperation, op) {
			t.Error("Expected operations to be equal")
		}
	})

	t.Run("panic on missing operation", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatalf("Expected panic")
			}
		}()

		req, err := http.NewRequest(http.MethodGet, "/path", nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		_ = MustOperation(req)
	})

}
