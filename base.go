package oas

import "net/http"

// BaseRouter is an underlying base router used in oas router.
// Any third-party router can be a base Router by using an adapter pattern.
type BaseRouter interface {
	http.Handler
	Compose(middlewares []Middleware, routes []Route)
}

// Route describes a routing item.
type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}
