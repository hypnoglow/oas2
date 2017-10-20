package oas2

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func TestQueryValidator(t *testing.T) {
	handlers := OperationHandlers{
		"loginUser": http.HandlerFunc(handleUserLogin),
	}
	errHandler := makeErrorHandler()
	router, err := NewRouter(loadDoc().Spec(), handlers, Use(NewQueryValidator(errHandler)))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("positive", func(t *testing.T) {
		resp, _ := helperGet(t, router, "/v2/user/login?username=johndoe&password=123")
		expectedpayload := "username: johndoe, password: 123"
		if !bytes.Equal([]byte(expectedpayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedpayload, string(resp))
		}
	})

	t.Run("validation error", func(t *testing.T) {
		resp, _ := helperGet(t, router, "/v2/user/login?username=johndoe")
		expectedPayload := `{"errors":[{"message":"param password is required","field":"password"}]}`
		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}
	})

	resourceHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})
	handler := NewQueryValidator(errHandler).Apply(resourceHandler)
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

func handleUserLogin(w http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	password := req.URL.Query().Get("password")

	// Never do this! This is just for testing purposes.
	fmt.Fprintf(w, "username: %s, password: %s", username, password)
}
