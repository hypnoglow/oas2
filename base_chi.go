package oas

import (
	"net/http"

	"github.com/go-chi/chi"
)

// ChiRouter returns a BaseRouter made from chi router (github.com/go-chi/chi).
//
// This can be considered as an example adapter implementation. You can implement
// your own adapter for a desired router if you need it.
func ChiRouter(router chi.Router) BaseRouter {
	return &chiRouter{
		router,
	}
}

// Router implements oas.BaseRouter using chi router.
type chiRouter struct {
	chi.Router
}

// Compose is an implementation of oas.BaseRouter.
func (r *chiRouter) Compose(middlewares []Middleware, routes []Route) {
	mws := make([]func(handler http.Handler) http.Handler, len(middlewares))
	for i, m := range middlewares {
		mws[i] = m
	}

	sub := r.Router.With(mws...)
	for _, route := range routes {
		sub.Method(route.Method, route.Path, route.Handler)
	}
}

// ChiPathTemplate is a PathTemplateFunc implemented by chi router.
func ChiPathTemplate(req *http.Request) string {
	return chi.RouteContext(req.Context()).RoutePattern()
}

// ChiPathParam is a PathParamFunc implemented by chi router.
var ChiPathParam = chi.URLParam
