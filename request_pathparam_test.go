package oas

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func TestPathParameterExtractor(t *testing.T) {
	handlers := OperationHandlers{
		"getPetById": http.HandlerFunc(handleGetPetByID),
	}

	router := NewRouter(RouterMiddleware(PathParameterExtractor(DefaultExtractorFunc)))
	err := router.AddSpec(loadDocFile(t, "testdata/petstore_1.yml"), handlers)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("positive", func(t *testing.T) {
		resp, _ := helperGet(t, router, "/v2/pet/12")
		expectedPayload := "pet by id: 12"
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})

	t.Run("request an url which handler does not provide operation context", func(t *testing.T) {
		resourceHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, "hit no operation resource")
		})
		var panicmsg string
		handler := PanicRecover(PathParameterExtractor(DefaultExtractorFunc)(resourceHandler), &panicmsg)
		noopRouter := chi.NewRouter()
		noopRouter.Handle("/resource", handler)

		helperGet(t, noopRouter, "/resource")
		expectedPanic := "request has no OpenAPI parameters in its context"
		if panicmsg != expectedPanic {
			t.Fatalf("Expected panic %q but got %q", expectedPanic, panicmsg)
		}
	})
}

func handleGetPetByID(w http.ResponseWriter, req *http.Request) {
	id, ok := GetPathParam(req, "petId").(int64)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "pet by id: %d", id)
}
