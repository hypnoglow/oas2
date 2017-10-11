package oas2

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func helperGet(t *testing.T, router http.Handler, url string) ([]byte, int) {
	server := httptest.NewServer(router)
	client := server.Client()
	defer server.Close()

	resp, err := client.Get(server.URL + url)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	return respBody, resp.StatusCode
}

func helperPost(t *testing.T, router http.Handler, url string, body io.Reader) []byte {
	server := httptest.NewServer(router)
	client := server.Client()
	defer server.Close()

	resp, err := client.Post(server.URL+url, "application/json", body)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	return respBody
}

func validationErrorsHandler(w http.ResponseWriter, errs []error) {
	p := convertErrs(errs)

	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		panic(err)
	}
}

func responseErrorHandler(buffer *bytes.Buffer) func(w http.ResponseWriter, errs []error) {
	l := log.New(buffer, "", 0)
	return func(w http.ResponseWriter, errs []error) {
		for _, e := range convertErrs(errs).Errors {
			l.Printf(
				"response data does not match the schema: field=%s value=%v message=%s",
				e.Field,
				e.Value,
				e.Message,
			)
		}
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

func convertErrs(errs []error) payload {
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
