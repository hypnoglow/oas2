package oas2

import (
	"net/http"
)

// MiddlewareFunc describes middleware function.
type MiddlewareFunc func(next http.Handler) http.Handler

// Middleware describes a middleware that can be applied to a http.handler.
type Middleware interface {
	Apply(next http.Handler) http.Handler
}
