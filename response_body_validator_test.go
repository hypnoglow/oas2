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

	rbv := ResponseBodyValidator(errHandler, ContentTypeRegexSelector(contentTypeSelectorRegexJSON))

	router, err := NewRouter(
		loadDoc(petstore).Spec(),
		handlers,
		Use(rbv),
	)
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
		expectedPayload := `{"error":"foo"}`
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
		expectedPayload := `{"error":"not found"}`
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
		expectedLogBuffer := "json decode: unexpected EOF"

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
	handler := ResponseBodyValidator(errHandler)(resourceHandler)
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
	// set Content-Type for all responses to ensure validator does not filter
	// them out prior to other checks like response spec schema presence.
	w.Header().Set("Content-Type", "application/json")

	// fake not found
	if req.URL.Path == "/v2/pet/404" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
		return
	}

	// fake for server error {
	if req.URL.Path == "/v2/pet/500" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"foo"}`))
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
