package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestResponseBodyValidator(t *testing.T) {
	handlers := OperationHandlers{
		"getPetById": http.HandlerFunc(handleGetPetByIDFaked),
	}
	logBuffer := &bytes.Buffer{}
	errHandler := responseErrorHandler(logBuffer)

	router, err := NewRouter(loadDoc().Spec(), handlers, Use(NewResponseBodyValidator(errHandler)))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("positive", func(t *testing.T) {
		logBuffer.Reset()
		resp, statusCode := helperGet(t, router, "/v2/pet/12")
		respBody := strings.TrimSpace(string(resp))

		expectedStatusCode := http.StatusOK
		expectedPayload := `{"id":123,"name":"Kitty"}`
		expectedLogBuffer := "response data does not match the schema: field=age value=<nil> message=age in body is required"

		if expectedStatusCode != statusCode {
			t.Errorf("Expected status code to be %v but got %v", expectedStatusCode, statusCode)
		}

		if expectedPayload != respBody {
			t.Errorf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, respBody)
		}

		actualLogBuff := strings.TrimSpace(logBuffer.String())
		if expectedLogBuffer != actualLogBuff {
			t.Errorf("Expected log buffer to be\n%v\nbut got\n%v\n", expectedLogBuffer, actualLogBuff)
		}
	})

	t.Run("no spec for 500", func(t *testing.T) {
		logBuffer.Reset()
		resp, statusCode := helperGet(t, router, "/v2/pet/500")
		respBody := strings.TrimSpace(string(resp))

		expectedStatusCode := http.StatusInternalServerError
		expectedPayload := `Internal Server Error`
		expectedLogBuffer := ""

		if expectedStatusCode != statusCode {
			t.Errorf("Expected status code to be %v but got %v", expectedStatusCode, statusCode)
		}

		if expectedPayload != respBody {
			t.Errorf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, respBody)
		}

		actualLogBuff := strings.TrimSpace(logBuffer.String())
		if expectedLogBuffer != actualLogBuff {
			t.Errorf("Expected log buffer to be\n%v\nbut got\n%v\n", expectedLogBuffer, actualLogBuff)
		}
	})

	t.Run("no schema for 404", func(t *testing.T) {
		logBuffer.Reset()
		resp, statusCode := helperGet(t, router, "/v2/pet/404")
		respBody := strings.TrimSpace(string(resp))

		expectedStatusCode := http.StatusNotFound
		expectedPayload := "404 page not found"
		expectedLogBuffer := ""

		if expectedStatusCode != statusCode {
			t.Errorf("Expected status code to be %v but got %v", expectedStatusCode, statusCode)
		}

		if expectedPayload != respBody {
			t.Errorf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, respBody)
		}

		actualLogBuff := strings.TrimSpace(logBuffer.String())
		if expectedLogBuffer != actualLogBuff {
			t.Errorf("Expected log buffer to be\n%v\nbut got\n%v\n", expectedLogBuffer, actualLogBuff)
		}
	})

	t.Run("bad json body", func(t *testing.T) {
		logBuffer.Reset()
		resp, statusCode := helperGet(t, router, "/v2/pet/badjson")

		expectedStatusCode := http.StatusOK
		expectedPayload := `{"name":`
		expectedLogBuffer := "json decode: unexpected end of JSON input"

		if expectedStatusCode != statusCode {
			t.Errorf("Expected status code to be %v but got %v", expectedStatusCode, statusCode)
		}

		if !bytes.Equal([]byte(expectedPayload), resp) {
			t.Errorf("Expected response body to be\n%s\nbut got\n%s", expectedPayload, string(resp))
		}

		actualLogBuff := strings.TrimSpace(logBuffer.String())
		if expectedLogBuffer != actualLogBuff {
			t.Errorf("Expected log buffer to be\n%v\nbut got\n%v\n", expectedLogBuffer, actualLogBuff)
		}
	})

	resourceHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})
	handler := NewResponseBodyValidator(errHandler).Apply(resourceHandler)
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

func handleGetPetByIDFaked(w http.ResponseWriter, req *http.Request) {
	// fake not found
	if req.URL.Path == "/v2/pet/404" {
		http.NotFound(w, req)
		return
	}

	// fake for server error {
	if req.URL.Path == "/v2/pet/500" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fake for bad json
	if req.URL.Path == "/v2/pet/badjson" {
		w.Write([]byte(`{"name":`))
		return
	}

	// normal

	type pet struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	p := pet{123, "Kitty"}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
}
