package oas2

import (
	"net/http"

	"github.com/go-chi/chi"
)

type chiRouter struct {
	chi.Router
}

func (r chiRouter) Route(method, pathPattern string, handler http.Handler) {
	r.Method(method, pathPattern, handler)
}

// ChiAdapter returns a BaseRouter made from chi.BaseRouter.
// More about router: github.com/go-chi/chi
func ChiAdapter(router chi.Router) BaseRouter {
	return chiRouter{
		router,
	}
}

// ChiAdapterFactory returns a BaseRouterFactory that produces BaseRouter made from chi.BaseRouter.
// More about router: github.com/go-chi/chi
func ChiAdapterFactory(router chi.Router) func() BaseRouter {
	return func() BaseRouter {
		return ChiAdapter(router)
	}
}

func defaultBaseRouterFactory() func() BaseRouter {
	return func() BaseRouter {
		return ChiAdapter(chi.NewRouter())
	}
}
