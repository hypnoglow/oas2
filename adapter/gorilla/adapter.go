package oas_gorilla

import (
	"github.com/gorilla/mux"

	"github.com/hypnoglow/oas2"
)

// adapter implements oas.Adapter using gorilla mux router.
type adapter struct{}

// Resolver returns a resolver based on gorilla mux router context.
func (a adapter) Resolver(meta interface{}) oas.Resolver {
	doc, ok := meta.(*oas.Document)
	if !ok {
		panic("oas_chi: Resolver meta is not *oas.Document")
	}

	return NewResolver(doc)
}

// OperationRouter returns an operation router based on gorilla mux router.
func (a adapter) OperationRouter(meta interface{}) oas.OperationRouter {
	r, ok := meta.(*mux.Router)
	if !ok {
		panic("oas_chi: OperationRouter meta is not *mux.Router")
	}

	return NewOperationRouter(r)
}

// PathParamExtractor returns a new path param extractor based on gorilla mux
// router context.
func (a adapter) PathParamExtractor() oas.PathParamExtractor {
	return NewPathParamExtractor()
}

// NewAdapter returns a new adapter based on gorilla mux router.
func NewAdapter() oas.Adapter {
	return adapter{}
}
