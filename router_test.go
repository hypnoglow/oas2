package oas2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewRouter(t *testing.T) {
	doc := loadDoc()

	handlers := OperationHandlers{
		"addPet": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "Hello, addPet!")
		}),
	}

	// router with default base router
	r, err := NewRouter(doc.Spec(), handlers)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(r)
	client := server.Client()

	body := bytes.NewBufferString(`{"name": "Rex"}`)
	resp, err := client.Post(server.URL+"/v2/pet", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte("Hello, addPet!"), respBody) {
		t.Fatalf("Expected response body to be %s but got %s", "Hello, addPet!", string(respBody))
	}

	// tear down

	server.Close()
}

func TestLoggerOpt(t *testing.T) {
	lg := &logrus.Logger{Out: ioutil.Discard}
	opt := LoggerOpt(lg)

	opts := &RouterOptions{}

	opt(opts)

	if !reflect.DeepEqual(opts.logger, lg) {
		t.Fatalf("Expected logger to be %v but got %v", lg, opts.logger)
	}
}

func TestBaseRouterOpt(t *testing.T) {
	baseRouter := defaultBaseRouter()
	opt := BaseRouterOpt(baseRouter)

	opts := &RouterOptions{}

	opt(opts)

	if !reflect.DeepEqual(opts.baseRouter, baseRouter) {
		t.Fatalf("Expected base router to be %v but got %v", baseRouter, opts.baseRouter)
	}
}
