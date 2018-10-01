package oas_gorilla

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/hypnoglow/oas2"
)

func TestResolver(t *testing.T) {
	doc, err := oas.LoadFile("testdata/petstore.yml")
	assert.NoError(t, err)

	resolver := NewResolver(doc)

	h := addPetHandler{
		resolver: resolver,
		t:        t,
	}

	r := mux.NewRouter()
	r.Path("/v2/pet").
		Methods(http.MethodPost).
		Handler(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v2/pet", nil)
	r.ServeHTTP(w, req)
}

type addPetHandler struct {
	resolver oas.Resolver
	t        *testing.T
}

func (h addPetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	op, ok := h.resolver.Resolve(req)
	assert.True(h.t, ok)
	assert.Equal(h.t, "addPet", op)
}
