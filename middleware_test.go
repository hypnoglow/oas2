package oas2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	ve := ValidationError{
		error: errors.New("got validation errors"),
		errs:  []error{fmt.Errorf("field is invalid")},
	}

	s := ve.Error()
	expectedString := "got validation errors: - field is invalid"
	if s != expectedString {
		t.Errorf("Expected %q but got %q", expectedString, s)
	}
}

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

func makeErrorHandler() RequestErrorHandler {
	return func(w http.ResponseWriter, req *http.Request, err error) (resume bool) {

		switch err.(type) {
		case ValidationError:
			e := err.(ValidationError)
			p := convertErrs(e.Errors())
			b, err := json.Marshal(p)
			if err != nil {
				panic(err)
			}

			if _, err := w.Write(b); err != nil {
				panic(err)
			}
			return false // do not continue

		case JsonError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"errors":[{"message":"Body contains invalid json"}]}`))
			return false // do not continue

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false // do not continue
		}
	}
}

func responseErrorHandler(buffer *bytes.Buffer) ResponseErrorHandler {
	l := log.New(buffer, "", 0)

	return func(w http.ResponseWriter, req *http.Request, err error) {
		switch err.(type) {
		case ValidationError:
			ve := err.(ValidationError)
			for _, e := range convertErrs(ve.Errors()).Errors {
				l.Printf(
					"response data does not match the schema: field=%s value=%v message=%s",
					e.Field,
					e.Value,
					e.Message,
				)
			}
		default:
			l.Print(err)
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
