package oas_chi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/hypnoglow/oas2"
)

// NewOperationRouter returns a new operation router based on chi router.
func NewOperationRouter(r chi.Router) oas.OperationRouter {
	return &OperationRouter{
		router: r,
	}
}

type OperationRouter struct {
	router chi.Router

	doc      *oas.Document
	mws      []oas.Middleware
	handlers map[string]http.Handler

	// onMissingOperationHandler is invoked with operation name
	// when operation handler is missing.
	onMissingOperationHandler func(op string)
}

// WithDocument sets the OpenAPI specification to build routes on.
// It returns the router for convenient chaining.
func (r *OperationRouter) WithDocument(doc *oas.Document) oas.OperationRouter {
	r.doc = doc
	return r
}

// WithMiddleware sets the middleware to build routing with.
// It returns the router for convenient chaining.
func (r *OperationRouter) WithMiddleware(mws ...oas.Middleware) oas.OperationRouter {
	r.mws = append(r.mws, mws...)
	return r
}

// WithOperationHandlers sets operation handlers to build routing with.
// It returns the router for convenient chaining.
func (r *OperationRouter) WithOperationHandlers(hh map[string]http.Handler) oas.OperationRouter {
	r.handlers = hh
	return r
}

// WithMissingOperationHandlerFunc sets the function that will be called
// for each operation that is present in the spec but missing from operation
// handlers. This is completely optional. You can use this method for example
// to simply log a warning or to throw a panic and stop route building.
// This method returns the router for convenient chaining.
func (r *OperationRouter) WithMissingOperationHandlerFunc(fn func(string)) oas.OperationRouter {
	r.onMissingOperationHandler = fn
	return r
}

// Route builds routing based on the previously provided specification,
// operation handlers, and other options.
func (r *OperationRouter) Route() error {
	if r.doc == nil {
		return fmt.Errorf("no doc is given")
	}
	if r.handlers == nil {
		return fmt.Errorf("no operation handlers given")
	}

	var router chi.Router = chi.NewRouter()

	mws := make([]func(http.Handler) http.Handler, len(r.mws))
	for i, mw := range r.mws {
		mws[i] = mw
	}

	router = router.With(mws...)

	for method, pathOps := range r.doc.Analyzer.Operations() {
		for path, operation := range pathOps {
			h, ok := r.handlers[operation.ID]
			if !ok {
				if r.onMissingOperationHandler != nil {
					r.onMissingOperationHandler(operation.ID)
				}
				continue
			}

			router.Method(method, path, h)
		}
	}

	if len(router.Routes()) == 0 {
		return nil
	}

	r.router.Mount(r.doc.BasePath(), router)

	return nil
}
