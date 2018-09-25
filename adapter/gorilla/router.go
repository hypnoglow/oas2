package gorilla

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hypnoglow/oas2"
)

// NewOperationRouter returns a new operation router based on gorilla/mux
// router.
func NewOperationRouter(r *mux.Router) oas.OperationRouter {
	return &OperationRouter{
		router: r,
	}
}

// OperationRouter is an operation router based on gorilla/mux router.
type OperationRouter struct {
	router   *mux.Router
	doc      *oas.Document
	mws      []oas.Middleware
	handlers map[string]http.Handler

	// onMissingOperationHandler is invoked with operation name
	// when operation handler is missing.
	onMissingOperationHandler func(op string)
}

func (r *OperationRouter) WithDocument(doc *oas.Document) oas.OperationRouter {
	r.doc = doc
	return r
}

func (r *OperationRouter) WithMiddleware(mws ...oas.Middleware) oas.OperationRouter {
	r.mws = append(r.mws, mws...)
	return r
}

func (r *OperationRouter) WithOperationHandlers(hh map[string]http.Handler) oas.OperationRouter {
	r.handlers = hh
	return r
}

func (r *OperationRouter) WithMissingOperationHandlerFunc(fn func(string)) oas.OperationRouter {
	r.onMissingOperationHandler = fn
	return r
}

func (r *OperationRouter) Route() error {
	if r.doc == nil {
		return fmt.Errorf("no doc is given")
	}
	if r.handlers == nil {
		return fmt.Errorf("no operation handlers given")
	}

	router := r.router.
		PathPrefix(r.doc.BasePath()).
		Subrouter()

	mws := make([]mux.MiddlewareFunc, len(r.mws))
	for i, mw := range r.mws {
		mws[i] = mux.MiddlewareFunc(mw)
	}
	router.Use(mws...)

	for method, pathOps := range r.doc.Analyzer.Operations() {
		for path, operation := range pathOps {
			h, ok := r.handlers[operation.ID]
			if !ok {
				if r.onMissingOperationHandler != nil {
					r.onMissingOperationHandler(operation.ID)
				}
				continue
			}

			router.Path(path).Methods(method).Handler(h)
		}
	}

	return nil
}
