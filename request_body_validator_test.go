package oas2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func TestBodyValidator(t *testing.T) {
	handlers := OperationHandlers{
		"addPet": http.HandlerFunc(handleAddPet),
	}
	errHandler := makeErrorHandler()
	router, err := NewRouter(loadDoc().Spec(), handlers, Use(NewBodyValidator(errHandler)))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("positive", func(t *testing.T) {
		resp := helperPost(t, router, "/v2/pet", bytes.NewBufferString(`{"name":"johndoe", "age":7}`))
		expectedPayload := "pet name: johndoe"
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})

	t.Run(`required field "name" is missing`, func(t *testing.T) {
		resp := helperPost(t, router, "/v2/pet", bytes.NewBufferString(`{"age":7}`))
		expectedPayload := `{"errors":[{"message":"name in body is required","field":"name"}]}`
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})

	t.Run(`value for field "age" is incorrect`, func(t *testing.T) {
		resp := helperPost(t, router, "/v2/pet", bytes.NewBufferString(`{"name":"johndoe","age":"abc"}`))
		expectedPayload := `{"errors":[{"message":"age in body must be of type integer: \"string\"","field":"age"}]}`
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})

	t.Run("no body", func(t *testing.T) {
		resp := helperPost(t, router, "/v2/pet", nil)
		expectedPayload := ""
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		resp := helperPost(t, router, "/v2/pet", bytes.NewBufferString(`{"name":"johndoe`))
		expectedPayload := `{"errors":[{"message":"Body contains invalid json"}]}`
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})

	resourceHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})
	handler := NewBodyValidator(errHandler).Apply(resourceHandler)
	noopRouter := chi.NewRouter()
	noopRouter.Handle("/resource", handler)

	t.Run("request an url which handler does not provide operation context", func(t *testing.T) {
		resp := helperPost(t, noopRouter, "/resource", bytes.NewBufferString(`{"name":"johndoe`))
		expectedPayload := "hit no operation resource"
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})
}

func handleAddPet(w http.ResponseWriter, req *http.Request) {
	type pet struct {
		Name      string   `json:"name"`
		PhotoURLs []string `json:"photoUrls"`
	}

	var p pet
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "pet name: %s", p.Name)
}
