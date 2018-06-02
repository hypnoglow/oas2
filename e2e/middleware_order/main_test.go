// +build e2e

package middleware_order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hypnoglow/oas2"
)

func TestMiddlewareExecutionOrder(t *testing.T) {
	doc, err := oas.LoadFile("testdata/greeter.yaml")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	handlers := oas.OperationHandlers{
		"greet": GreetHandler{},
	}

	buffer := &bytes.Buffer{}

	router, err := oas.NewRouter(
		doc,
		handlers,
		// We are testing that RequestIDLogger will have access to the request id
		// in the request created by RequestID middleware.
		oas.Use(RequestID),
		oas.Use(RequestIDLogger(log.New(buffer, "", 0))),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	srv := httptest.NewServer(router)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/api/greeting?name=Andrew")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected response status: %s", resp.Status)
	}

	var reply struct {
		Greeting string `json:"greeting"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedGreeting := "Hello, Andrew!"
	if reply.Greeting != expectedGreeting {
		t.Fatalf("Expected greeting to be %q but got %q", expectedGreeting, reply.Greeting)
	}

	expectedLogEntry := "request with id 1234567890\n"
	if buffer.String() != expectedLogEntry {
		t.Fatalf("Expected log entry to be %q but got %q", expectedLogEntry, buffer.String())
	}
}

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

type ctxKeyRequestID struct{}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(
			req.Context(),
			ctxKeyRequestID{},
			"1234567890",
		)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func GetRequestID(req *http.Request) string {
	id, ok := req.Context().Value(ctxKeyRequestID{}).(string)
	if !ok {
		return ""
	}
	return id
}

func RequestIDLogger(log *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if id := GetRequestID(req); id != "" {
				log.Printf("request with id %s", id)
			} else {
				log.Printf("request with no id")
			}
			next.ServeHTTP(w, req)
		})
	}
}
