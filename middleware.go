package oas2

import (
	"net/http"
)

// MiddlewareFn describes middleware function.
type MiddlewareFn func(next http.Handler) http.Handler

// Middleware describes a middleware that can be applied to a http.handler.
type Middleware interface {
	Apply(next http.Handler) http.Handler
}
