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
			expectedPayload: `{"errors":[{"message":"param password is required","field":"password"}]}`,
		},
		// request an url which handler does not provide operation context
		{
			url:             "/no_operation_resource",
			expectedPayload: "hit no operation resource",
		},
	}

	// set up

	doc := loadDoc()

	handlers := OperationHandlers{"loginUser": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		username := req.URL.Query().Get("username")
		password := req.URL.Query().Get("password")
		fmt.Fprintf(w, "username: %s, password: %s", username, password)
	})}
	noOpHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})

	qv := NewQueryValidator(doc.Spec(), errHandler)
	opts := []Option{MiddlewareOpt(qv.Apply)}

	operationsRouter, err := NewRouter(doc.Spec(), handlers, opts...)
	if err != nil {
		t.Fatal(err)
	}

	finalRouter := chi.NewRouter()
	finalRouter.Mount("/", operationsRouter)
	finalRouter.Handle("/no_operation_resource", qv.Apply(noOpHandler))

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

func TestBodyValidatorMiddleware_Apply(t *testing.T) {
	cases := []struct {
		url                string
		body               *bytes.Buffer
		expectedStatusCode int
		expectedPayload    string
	}{
		// ok
		{
			url:                "/pet",
			body:               bytes.NewBufferString(`{"name":"johndoe", "age":7}`),
			expectedStatusCode: http.StatusOK,
			expectedPayload:    "pet name: johndoe",
		},
		// required field "name" is missing
		{
			url:                "/pet",
			body:               bytes.NewBufferString(`{"age":7}`),
			expectedStatusCode: http.StatusBadRequest,
			expectedPayload:    `{"errors":[{"message":"name in body is required","field":"name"}]}`,
		},
		// value for field "age" is incorrect
		{
			url:                "/pet",
			body:               bytes.NewBufferString(`{"name":"johndoe","age":"abc"}`),
			expectedStatusCode: http.StatusBadRequest,
			expectedPayload:    `{"errors":[{"message":"age in body must be of type integer: \"string\"","field":"age"}]}`,
		},
		// no body
		{
			url:                "/pet",
			body:               &bytes.Buffer{},
			expectedStatusCode: http.StatusBadRequest,
		},
		// invalid json body
		{
			url:                "/pet",
			body:               bytes.NewBufferString(`{"name":"johndoe`),
			expectedStatusCode: http.StatusBadRequest,
			expectedPayload:    `{"errors":[{"message":"Body contains invalid json"}]}`,
		},
		// request an url which handler does not provide operation context
		{
			url:                "/no_operation_resource",
			body:               bytes.NewBufferString(`{"name":"johndoe`),
			expectedStatusCode: http.StatusOK,
			expectedPayload:    "hit no operation resource",
		},
	}

	// set up

	doc := loadDoc()

	handlers := OperationHandlers{"addPet": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
	})}
	noOpHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})

	bodyValidator := NewBodyValidator(doc.Spec(), errHandler)
	opts := []Option{MiddlewareOpt(bodyValidator.Apply)}

	operationsRouter, err := NewRouter(doc.Spec(), handlers, opts...)
	if err != nil {
		t.Fatal(err)
	}

	finalRouter := chi.NewRouter()
	finalRouter.Mount("/", operationsRouter)
	finalRouter.Handle("/no_operation_resource", bodyValidator.Apply(noOpHandler))

	server := httptest.NewServer(finalRouter)
	client := server.Client()

	// test

	for _, c := range cases {
		resp, err := client.Post(server.URL+c.url, "application/json", c.body)
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
			Field   string      `json:"field,omitempty"`
			Value   interface{} `json:"value,omitempty"`
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
