package oas2

import (
	"net/http"

	"github.com/go-chi/chi"
)

// ChiAdapter returns a BaseRouter made from chi router (github.com/go-chi/chi).
// This is just an example adapter implementation, you should implement
// your own adapter for a desired router if you need it.
func ChiAdapter(router chi.Router) BaseRouter {
	return chiRouter{
		router,
	}
}

type chiRouter struct {
	chi.Router
}

func (r chiRouter) Route(method, pathPattern string, handler http.Handler) {
	r.Method(method, pathPattern, handler)
}

func defaultBaseRouter() BaseRouter {
	return ChiAdapter(chi.NewRouter())
}
