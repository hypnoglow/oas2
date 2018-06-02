package oas

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNewRouter(t *testing.T) {
	doc := loadDocFile(t, "testdata/petstore_1.yml")

	handlers := OperationHandlers{
		"addPet": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "Hello, addPet!")
		}),
	}

	// router with default base router
	r, err := NewRouter(doc, handlers)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	server := httptest.NewServer(r)
	client := server.Client()

	body := bytes.NewBufferString(`{"name": "Rex"}`)
	resp, err := client.Post(server.URL+"/v2/pet", "application/json", body)
	if err != nil {
		t.Fatalf("HTTP POST request failed: %v", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if !bytes.Equal([]byte("Hello, addPet!"), respBody) {
		t.Fatalf("Expected response body to be %s but got %s", "Hello, addPet!", string(respBody))
	}

	// tear down

	server.Close()
}

func TestDebugLog(t *testing.T) {
	buf := &bytes.Buffer{}
	lg := log.New(buf, "", log.LstdFlags)

	opt := DebugLog(lg.Printf)

	router := &Router{}
	opt(router)

	router.debugLog("hello %s", "debugLog")

	if !strings.Contains(buf.String(), "hello debugLog") {
		t.Fatalf("Expected buf to contain hello debugLog")
	}
}

func TestBaseRouterOpt(t *testing.T) {
	baseRouter := defaultBaseRouter()
	opt := Base(baseRouter)

	router := &Router{}

	opt(router)

	if !reflect.DeepEqual(router.baseRouter, baseRouter) {
		t.Fatalf("Expected base router to be %v but got %v", baseRouter, router.baseRouter)
	}
}

func TestUse(t *testing.T) {
	r := &Router{}
	mw := QueryValidator(nil)
	opt := Use(mw)
	opt(r)

	if len(r.mws) == 0 {
		t.Errorf("Expected to apply middleware")
	}
}

func TestServeSpec(t *testing.T) {
	doc, err := LoadFile("testdata/petstore_1.yml")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	r, err := NewRouter(
		doc,
		OperationHandlers{},
		ServeSpec(SpecHandlerTypeStatic),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2", nil)
	r.ServeHTTP(w, req)

	var paths []string
	for name := range doc.Spec().Paths.Paths {
		paths = append(paths, name)
	}
	sort.Strings(paths)

	expectedPaths := []string{"/pet", "/pet/{petId}", "/user/login"}
	if !reflect.DeepEqual(expectedPaths, paths) {
		t.Errorf("Expected output spec paths to be %v but got %v", expectedPaths, paths)
	}
}
