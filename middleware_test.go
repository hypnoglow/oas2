package oas2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestQueryValidatorMiddleware_Apply(t *testing.T) {
	cases := []struct {
		url             string
		expectedPayload string
	}{
		// ok
		{
			url:             "/v2/user/login?username=johndoe&password=123",
			expectedPayload: "username: johndoe, password: 123",
		},
		// required parameter is not passed
		{
			url:             "/v2/user/login?username=johndoe",
			expectedPayload: `{"errors":[{"message":"param password is required","field":"password","value":null}]}`,
		},
		// request an url which handler does not provide operation id context
		{
			url:             "/no_operation_id_resource",
			expectedPayload: "hit no operation id resource",
		},
		// request an url which handler has an operation id but no spec for that operation
		{
			url:             "/no_operation_for_operation_id_resource",
			expectedPayload: "hit no operation for operation id resource",
		},
	}

	// set up

	doc := loadDoc()

	handlers := OperationHandlers{"loginUser": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		username := req.URL.Query().Get("username")
		password := req.URL.Query().Get("password")
		fmt.Fprintf(w, "username: %s, password: %s", username, password)
	})}

	noOpIDHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation id resource")
	})
	noOpByOpIDHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation for operation id resource")
	})

	qv := NewQueryValidator(doc.Spec(), errHandler)
	opts := []Option{MiddlewareOpt(qv.Apply)}

	operationsRouter, err := NewRouter(doc.Spec(), handlers, opts...)
	if err != nil {
		t.Fatal(err)
	}

	finalRouter := chi.NewRouter()
	finalRouter.Mount("/", operationsRouter)
	finalRouter.Handle("/no_operation_id_resource", qv.Apply(noOpIDHandler))
	finalRouter.Handle("/no_operation_for_operation_id_resource", operationIDMiddleware(qv.Apply(noOpByOpIDHandler), "someOperationID"))

	server := httptest.NewServer(finalRouter)
	client := server.Client()

	// test

	for _, c := range cases {
		resp, err := client.Get(server.URL + c.url)
		if err != nil {
			t.Fatal(err)
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal([]byte(c.expectedPayload), respBody) {
			t.Fatalf("Expected response body to be\n%s\nbut got\n%s", c.expectedPayload, string(respBody))
		}
	}

	// tear down

	server.Close()
}

func errHandler(w http.ResponseWriter, errs []error) {
	// This is an example of composing an error for response from validation
	// errors.

	type (
		errorItem struct {
			Message string      `json:"message"`
			Field   string      `json:"field"`
			Value   interface{} `json:"value"`
		}
		payload struct {
			Errors []errorItem `json:"errors"`
		}
	)

	type fielder interface {
		Field() string
	}

	type valuer interface {
		Value() interface{}
	}

	p := payload{Errors: make([]errorItem, 0)}
	for _, e := range errs {
		item := errorItem{Message: e.Error()}
		if fe, ok := e.(fielder); ok {
			item.Field = fe.Field()
		}
		if ve, ok := e.(valuer); ok {
			item.Value = ve.Value()
		}
		p.Errors = append(p.Errors, item)
	}

	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}
