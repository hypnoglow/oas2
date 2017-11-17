package oas

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func TestChiRouter_Route(t *testing.T) {
	r := ChiAdapter(chi.NewRouter())

	hf := func(w http.ResponseWriter, req *http.Request) {
		if _, err := io.WriteString(w, "handler func"); err != nil {
			log.Fatal(err)
		}
	}

	r.Route("GET", "/path", http.HandlerFunc(hf))
}

func TestChiAdapter(t *testing.T) {
	ChiAdapter(chi.NewRouter())
}
