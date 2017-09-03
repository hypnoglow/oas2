package oas2

import (
	"testing"

	"io"
	"net/http"

	"github.com/go-chi/chi"
)

func TestChiRouter_Route(t *testing.T) {
	r := ChiAdapter(chi.NewRouter())

	hf := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "handler func")
	}

	r.Route("GET", "/path", http.HandlerFunc(hf))
}

func TestChiAdapter(t *testing.T) {
	ChiAdapter(chi.NewRouter())
}

func TestChiAdapterFactory(t *testing.T) {
	f := ChiAdapterFactory(chi.NewRouter())
	f()
}
