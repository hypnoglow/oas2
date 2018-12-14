package oas_chi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/adapter/chi"
	_ "github.com/hypnoglow/oas2/adapter/chi/init"
)

func TestOperationRouter_implementation(t *testing.T) {
	var _ oas.OperationRouter = &oas_chi.OperationRouter{}
}

func TestOperationRouter(t *testing.T) {
	doc, err := oas.LoadFile("testdata/petstore.yml")
	assert.NoError(t, err)

	r := chi.NewRouter()
	basis := oas.NewResolvingBasis("chi", doc)

	var notHandledOps []string

	err = basis.OperationRouter(r).
		WithOperationHandlers(map[string]http.Handler{
			"getPetById": getPetHandler{},
		}).
		WithMiddleware(basis.PathParamsContext()).
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
