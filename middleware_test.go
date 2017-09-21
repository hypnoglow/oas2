package oas2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

	qv := NewQueryValidator(writeErrorsToResponseWriter)
	opts := []RouterOption{MiddlewareOpt(qv.Apply)}

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

	bodyValidator := NewBodyValidator(writeErrorsToResponseWriter)
	opts := []RouterOption{MiddlewareOpt(bodyValidator.Apply)}

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

func TestPathParameterExtractor_Apply(t *testing.T) {
	cases := []struct {
		url                string
		expectedStatusCode int
		expectedPayload    string
	}{
		// ok
		{
			url:                "/pet/12",
			expectedStatusCode: http.StatusOK,
			expectedPayload:    "pet by id: 12",
		},
		// request an url which handler does not provide operation context
		{
			url:                "/no_operation_resource",
			expectedStatusCode: http.StatusOK,
			expectedPayload:    "hit no operation resource",
		},
	}

	// set up

	doc := loadDoc()

	handlers := OperationHandlers{"getPetById": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id, ok := GetPathParam(req, "petId").(int64)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "pet by id: %d", id)
	})}
	noOpHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})

	pathParamExtractor := NewPathParameterExtractor(chi.URLParam)
	opts := []RouterOption{MiddlewareOpt(pathParamExtractor.Apply)}

	operationsRouter, err := NewRouter(doc.Spec(), handlers, opts...)
	if err != nil {
		t.Fatal(err)
	}

	finalRouter := chi.NewRouter()
	finalRouter.Mount("/", operationsRouter)
	finalRouter.Handle("/no_operation_resource", pathParamExtractor.Apply(noOpHandler))

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
			t.Errorf("Expected response body to be\n%s\nbut got\n%s", c.expectedPayload, string(respBody))
		}
	}

	// tear down

	server.Close()
}

func TestResponseBodyValidator_Apply(t *testing.T) {
	cases := []struct {
		url                string
		logBuffer          *bytes.Buffer
		expectedStatusCode int
		expectedPayload    string
		expectedLogBuffer  string
	}{
		// with validation errors
		{
			url:                "/pet/12",
			logBuffer:          &bytes.Buffer{},
			expectedStatusCode: http.StatusOK,
			expectedPayload:    `{"id":123,"name":"Kitty"}` + "\n",
			expectedLogBuffer:  "response data does not match the schema: field=age value=<nil> message=age in body is required",
		},
		// no spec for 500
		{
			url:                "/pet/500",
			logBuffer:          &bytes.Buffer{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedPayload:    "Internal Server Error\n",
		},
		// no schema for 404
		{
			url:                "/pet/13",
			logBuffer:          &bytes.Buffer{},
			expectedStatusCode: http.StatusNotFound,
			expectedPayload:    "404 page not found\n",
		},
		{
			url:                "/pet/badjson",
			logBuffer:          &bytes.Buffer{},
			expectedStatusCode: http.StatusOK,
			expectedPayload:    `{"name":`,
		},
		// request an url which handler does not provide operation context
		{
			url:                "/no_operation_resource",
			logBuffer:          &bytes.Buffer{},
			expectedStatusCode: http.StatusOK,
			expectedPayload:    "hit no operation resource",
			expectedLogBuffer:  "",
		},
	}

	// set up

	doc := loadDoc()

	handlers := OperationHandlers{"getPetById": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// fake not found
		if req.URL.Path == "/pet/13" {
			http.NotFound(w, req)
			return
		}

		// fake for server error {
		if req.URL.Path == "/pet/500" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// fake for bad json
		if req.URL.Path == "/pet/badjson" {
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
			t.Fatalf("Unexpected error: %s", err)
		}
	})}
	noOpHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hit no operation resource")
	})

	// test

	for _, c := range cases {
		respBodyValidator := NewResponseBodyValidator(errorLogger(c.logBuffer))
		opts := []RouterOption{MiddlewareOpt(respBodyValidator.Apply)}

		operationsRouter, err := NewRouter(doc.Spec(), handlers, opts...)
		if err != nil {
			t.Fatal(err)
		}

		finalRouter := chi.NewRouter()
		finalRouter.Mount("/", operationsRouter)
		finalRouter.Handle("/no_operation_resource", respBodyValidator.Apply(noOpHandler))

		server := httptest.NewServer(finalRouter)
		client := server.Client()

		resp, err := client.Get(server.URL + c.url)
		if err != nil {
			t.Fatal(err)
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if c.expectedStatusCode != resp.StatusCode {
			t.Errorf("Expected status code to be %v but got %v", c.expectedStatusCode, resp.StatusCode)
		}

		if !bytes.Equal([]byte(c.expectedPayload), respBody) {
			t.Errorf("Expected response body to be\n%s\nbut got\n%s", c.expectedPayload, string(respBody))
		}

		expectedLogBuff := strings.TrimSpace(c.expectedLogBuffer)
		actualLogBuff := strings.TrimSpace(c.logBuffer.String())
		if expectedLogBuff != actualLogBuff {
			t.Errorf("Expected log buffer to be\n%v\nbut got\n%v\n", expectedLogBuff, actualLogBuff)
		}

		// tear down
		server.Close()
	}
}

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

func helperConvertErrs(errs []error) payload {
	// This is an example of composing an error for response from validation
	// errors.

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

	return p
}

func writeErrorsToResponseWriter(w http.ResponseWriter, errs []error) {
	// This is an example of composing an error for response from validation
	// errors.

	p := helperConvertErrs(errs)

	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}

func errorLogger(buffer *bytes.Buffer) func(w http.ResponseWriter, errs []error) {
	l := log.New(buffer, "", 0)
	return func(w http.ResponseWriter, errs []error) {
		for _, e := range helperConvertErrs(errs).Errors {
			l.Printf(
				"response data does not match the schema: field=%s value=%v message=%s",
				e.Field,
				e.Value,
				e.Message,
			)
		}
	}
}
