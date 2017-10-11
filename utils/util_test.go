package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestResponseRecorder(t *testing.T) {
	w := httptest.NewRecorder()
	rr := NewResponseRecorder(w)

	h := rr.Header()
	if !reflect.DeepEqual(h, w.Header()) {
		t.Error("Expected header to be equal")
	}

	rr.Write([]byte("test response body"))
	if !bytes.Equal(rr.Payload(), w.Body.Bytes()) {
		t.Errorf("Expected body to be equal")
	}

	rr.WriteHeader(http.StatusOK)
	if w.Result().StatusCode != rr.Status() {
		t.Error("Expected status to be equal")
	}
}
