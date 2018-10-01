package oas_chi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/hypnoglow/oas2"
)

func TestOperationRouter_implementation(t *testing.T) {
	var _ oas.OperationRouter = &OperationRouter{}
}

func TestOperationRouter(t *testing.T) {
	doc, err := oas.LoadFile("testdata/petstore.yml")
	assert.NoError(t, err)

	r := chi.NewRouter()
	basis := oas.NewResolvingBasis(doc, NewResolver(doc))

	var notHandledOps []string

	err = NewOperationRouter(r).
		WithDocument(doc).
		WithOperationHandlers(map[string]http.Handler{
			"getPetById": getPetHandler{},
		}).
		WithMiddleware(basis.OperationContext()).
		WithMiddleware(basis.PathParamsContext(NewPathParamExtractor())).
		WithMissingOperationHandlerFunc(func(s string) {
			notHandledOps = append(notHandledOps, s)
		}).
		Route()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v2/pet/12", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"name": "Hooch", "age": 3, "debug": true}`, w.Body.String())
	assert.ElementsMatch(t, []string{"addPet", "loginUser"}, notHandledOps)
}

type getPetHandler struct{}

func (h getPetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id, ok := oas.GetPathParam(req, "petId").(int64)
	if !ok {
		fmt.Fprintf(os.Stderr, "DEBUG: %#v\n", id)
		os.Exit(1)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if id != 12 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"name":  "Hooch",
		"age":   3,
		"debug": true,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
