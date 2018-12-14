package oas_gorilla_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/adapter/gorilla"
	_ "github.com/hypnoglow/oas2/adapter/gorilla/init"
)

func TestOperationRouter_implementation(t *testing.T) {
	var _ oas.OperationRouter = &oas_gorilla.OperationRouter{}
}

func TestOperationRouter(t *testing.T) {
	doc, err := oas.LoadFile("testdata/petstore.yml")
	assert.NoError(t, err)

	r := mux.NewRouter()
	basis := oas.NewResolvingBasis("gorilla", doc)

	var notHandledOps []string

	err = basis.OperationRouter(r).
		WithOperationHandlers(map[string]http.Handler{
			"addPet": addPetHandler2{},
		}).
		WithMissingOperationHandlerFunc(func(s string) {
			notHandledOps = append(notHandledOps, s)
		}).
		Route()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v2/pet", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"foo":"bar"}`, w.Body.String())
	assert.ElementsMatch(t, []string{"getPetById", "loginUser"}, notHandledOps)
}

type addPetHandler2 struct{}

func (h addPetHandler2) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`{"foo":"bar"}`))
}
