package oas_chi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/hypnoglow/oas2"
)

// NewResolver returns a resolver that resolves OpenAPI operation ID using
// chi request route context. It should be used in conjunction with
// chi router, and only with it.
func NewResolver(doc *oas.Document) oas.Resolver {
	return &resolver{
		doc: doc,
	}
}

// resolver implements Resolver using chi's mux RouteContext.
type resolver struct {
	doc *oas.Document
}

// Resolve resolves operation id from the request using chi route
// context.
func (r *resolver) Resolve(req *http.Request) (string, bool) {
	ctx := chi.RouteContext(req.Context())
	pt := ctx.RoutePattern()
	if pt == "" {
		return "", false
	}

	p := strings.TrimPrefix(pt, r.doc.BasePath())
	op, ok := r.doc.Analyzer.OperationFor(req.Method, p)
	if !ok {
		return "", false
	}

	return op.ID, true
}
