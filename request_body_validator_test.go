package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestBodyValidator(t *testing.T) {
	handlers := OperationHandlers{
		"addPet": http.HandlerFunc(handleAddPet),
	}
	errHandler := makeErrorHandler()

	bv := BodyValidator(errHandler, ContentTypeRegexSelector(contentTypeSelectorRegexJSON))

	router := NewRouter(RouterMiddleware(bv))
	err := router.AddSpec(loadDocFile(t, "testdata/petstore_1.yml"), handlers)
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

	t.Run("should skip body validation for not application/json content type", func(t *testing.T) {
		server := httptest.NewServer(router)
		client := server.Client()
		defer server.Close()

		resp, err := client.Post(server.URL+"/v2/pet", "text/plain", bytes.NewBufferString(`some`))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected response status code to be %d but got %d", http.StatusBadRequest, resp.StatusCode)
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

	t.Run("request an url which handler does not provide operation context", func(t *testing.T) {
		resourceHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, "hit no operation resource")
		})
		var panicmsg string
		handler := PanicRecover(BodyValidator(errHandler)(resourceHandler), &panicmsg)
		noopRouter := chi.NewRouter()
		noopRouter.Handle("/resource", handler)

		helperPost(t, noopRouter, "/resource", bytes.NewBufferString(`{"name":"johndoe`))
		expectedPanic := "request has no OpenAPI parameters in its context"
		if panicmsg != expectedPanic {
			t.Fatalf("Expected panic %q but got %q", expectedPanic, panicmsg)
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
