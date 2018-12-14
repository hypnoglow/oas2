package oas

import (
	"sync"
)

// Adapter ... TODO
type Adapter interface {
	Resolver(meta interface{}) Resolver
	OperationRouter(meta interface{}) OperationRouter
	PathParamExtractor() PathParamExtractor

	// ...
}

var (
	adaptersMx sync.Mutex
	adapters   = make(map[string]Adapter)
)

// RegisterAdapter ... TODO
func RegisterAdapter(name string, adapter Adapter) {
	adaptersMx.Lock()
	defer adaptersMx.Unlock()

	if adapter == nil {
		panic("oas: RegisterAdapter adapter is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("oas: RegisterAdapter called twice for adapter " + name)
	}

	adapters[name] = adapter
}

func mustGetAdapter(name string) Adapter {
	adaptersMx.Lock()
	defer adaptersMx.Unlock()

	a, ok := adapters[name]
	if !ok {
		panic("oas: no adapter registered for name " + name)
	}
	return a
}
