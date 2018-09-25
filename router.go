package oas

import (
	"net/http"
)

// OperationRouter describes an OpenAPI operation router, which can build
// routing based on the specification, given operation handlers, and other
// options.
//
// See "adapter/*" packages for implementations.
type OperationRouter interface {
	// WithDocument sets the OpenAPI specification to build routes on.
	// It returns the router for convenient chaining.
	WithDocument(doc *Document) OperationRouter

	// WithMiddleware sets the middleware to build routing with.
	// It returns the router for convenient chaining.
	WithMiddleware(mws ...Middleware) OperationRouter

	// WithOperationHandlers sets operation handlers to build routing with.
	// It returns the router for convenient chaining.
	WithOperationHandlers(map[string]http.Handler) OperationRouter

	// WithMissingOperationHandlerFunc sets the function that will be called
	// for each operation that is present in the spec but missing from operation
	// handlers. This is completely optional. You can use this method for example
	// to simply log a warning or to throw a panic and stop route building.
	// This method returns the router for convenient chaining.
	WithMissingOperationHandlerFunc(fn func(string)) OperationRouter

	// Route builds routing based on the previously provided specification,
	// operation handlers, and other options.
	Route() error
}
