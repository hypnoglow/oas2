package oas_chi

import (
	"github.com/go-chi/chi"

	"github.com/hypnoglow/oas2"
)

type adapter struct {
	// TODO: cache existing objects in here?
}

func (a adapter) Resolver(meta interface{}) oas.Resolver {
	doc, ok := meta.(*oas.Document)
	if !ok {
		panic("oas_chi: Resolver meta is not *oas.Document")
	}

	return NewResolver(doc)
}

func (a adapter) OperationRouter(meta interface{}) oas.OperationRouter {
	r, ok := meta.(chi.Router)
	if !ok {
		panic("oas_chi: OperationRouter meta is not chi.Router")
	}

	return NewOperationRouter(r)
}

func (a adapter) PathParamExtractor() oas.PathParamExtractor {
	return NewPathParamExtractor()
}

// NewAdapter ... TODO
func NewAdapter() oas.Adapter {
	return adapter{}
}
