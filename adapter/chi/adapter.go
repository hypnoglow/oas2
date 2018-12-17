package oas_chi

import (
	"github.com/go-chi/chi"

	"github.com/hypnoglow/oas2"
)

// adapter implements oas.Adapter using chi router.
type adapter struct{}

// Resolver returns a resolver based on chi router context.
func (a adapter) Resolver(meta interface{}) oas.Resolver {
	doc, ok := meta.(*oas.Document)
	if !ok {
		panic("oas_chi: Resolver meta is not *oas.Document")
	}

	return NewResolver(doc)
}

// OperationRouter returns an operation router based on chi router.
func (a adapter) OperationRouter(meta interface{}) oas.OperationRouter {
	r, ok := meta.(chi.Router)
	if !ok {
		panic("oas_chi: OperationRouter meta is not chi.Router")
	}

	return NewOperationRouter(r)
}

// PathParamExtractor returns a new path param extractor based on chi router
// context.
func (a adapter) PathParamExtractor() oas.PathParamExtractor {
	return NewPathParamExtractor()
}

// NewAdapter returns a new adapter based on chi router.
func NewAdapter() oas.Adapter {
	return adapter{}
}
