// +build e2e

package middleware_order

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/adapter/gorilla"
	"github.com/hypnoglow/oas2/e2e/testdata"
)

// TestMiddlewareExecutionOrder tests that middleware passed to the router
// is executed in correct order.
func TestMiddlewareExecutionOrder(t *testing.T) {
	doc := testdata.GreeterSpec(t)
	basis := oas.NewResolvingBasis(doc, oas_gorilla.NewResolver(doc))

	t.Run("middleware passed inline with RouterMiddleware()", func(t *testing.T) {
		buffer := &bytes.Buffer{}

		r := mux.NewRouter()
		err := oas_gorilla.NewOperationRouter(r).
			WithDocument(doc).
			WithOperationHandlers(map[string]http.Handler{
				"greet": testdata.GreetHandler{},
			}).
			WithMiddleware(
				basis.OperationContext(),
				// We are testing that RequestIDLogger will have access to the request id
				// in the request created by RequestID middleware.
				RequestID,
				RequestIDLogger(log.New(buffer, "", 0)),
			).Route()
		assert.NoError(t, err)

		testRouterMiddleware(t, r, buffer)
	})

	t.Run("middleware passed as two options with RouterMiddleware()", func(t *testing.T) {
		buffer := &bytes.Buffer{}

		r := mux.NewRouter()
		err := oas_gorilla.NewOperationRouter(r).
			WithDocument(doc).
			WithOperationHandlers(map[string]http.Handler{
				"greet": testdata.GreetHandler{},
			}).
			WithMiddleware(basis.OperationContext()).
			// We are testing that RequestIDLogger will have access to the request id
			// in the request created by RequestID middleware.
			WithMiddleware(RequestID).
			WithMiddleware(RequestIDLogger(log.New(buffer, "", 0))).
			Route()
		assert.NoError(t, err)

		testRouterMiddleware(t, r, buffer)
	})
}

func testRouterMiddleware(t *testing.T, router *mux.Router, buf *bytes.Buffer) {
	t.Helper()

	srv := httptest.NewServer(router)
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL + "/api/greeting?name=Andrew")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var reply struct {
		Greeting string `json:"greeting"`
	}

	if err := json.Unmarshal(b, &reply); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedGreeting := "Hello, Andrew!"
	if reply.Greeting != expectedGreeting {
		t.Fatalf("Expected greeting to be %q but got %q", expectedGreeting, reply.Greeting)
	}

	expectedLogEntry := "request with id 1234567890\n"
	if buf.String() != expectedLogEntry {
		t.Fatalf("Expected log entry to be %q but got %q", expectedLogEntry, buf.String())
	}
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
