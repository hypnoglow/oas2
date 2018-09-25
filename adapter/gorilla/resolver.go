package gorilla

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/hypnoglow/oas2"
)

// NewResolver returns a resolver that resolves OpenAPI operation ID using
// gorilla/mux request route context. It should be used in conjunction with
// gorilla/mux router, and only with it.
func NewResolver(doc *oas.Document) oas.Resolver {
	return &resolver{
		doc: doc,
	}
}

// resolver implements Resolver using gorilla's mux CurrentRoute method
// that extracts path template from the request.
type resolver struct {
	doc *oas.Document
}

// Resolve resolves operation id from the request using gorilla/mux route
// context.
func (r *resolver) Resolve(req *http.Request) (string, bool) {
	cr := mux.CurrentRoute(req)
	if cr == nil {
		// WARNING: this can happen because of improper package vendoring.
		return "", false
	}

	pt, err := cr.GetPathTemplate()
	if err != nil {
		return "", false
	}

	p := strings.TrimPrefix(pt, r.doc.BasePath())
	op, ok := r.doc.Analyzer.OperationFor(req.Method, p)
	if !ok {
		return "", false
	}

	return op.ID, true
}
