package oas

import (
	"sync"
)

// Adapter interface defines a contract between oas and a router implementation.
type Adapter interface {
	Resolver(meta interface{}) Resolver
	OperationRouter(meta interface{}) OperationRouter
	PathParamExtractor() PathParamExtractor
}

var (
	adaptersMx sync.Mutex
	adapters   = make(map[string]Adapter)
)

// RegisterAdapter makes an adapter available by the provided name. If this
// function is called twice with the same name or if adapter is nil, it panics.
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

// mustGetAdapter returns previously registered adapter by the provided name.
// If no adapter is registered by the name, it panics.
func mustGetAdapter(name string) Adapter {
	adaptersMx.Lock()
	defer adaptersMx.Unlock()

	a, ok := adapters[name]
	if !ok {
		panic("oas: no adapter registered for name " + name)
	}
	return a
}
