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
	router, err := NewRouter(loadDoc().Spec(), handlers, Use(NewPathParameterExtractor(chi.URLParam)))
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

	resourceHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})
	handler := NewPathParameterExtractor(chi.URLParam).Apply(resourceHandler)
	noopRouter := chi.NewRouter()
	noopRouter.Handle("/resource", handler)

	t.Run("request an url which handler does not provide operation context", func(t *testing.T) {
		resp, _ := helperGet(t, noopRouter, "/resource")
		expectedPayload := "hit no operation resource"
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
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
