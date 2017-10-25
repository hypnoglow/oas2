package oas2

import (
	"net/http"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
)

// Router routes requests based on OAS 2.0 spec operations.
type Router struct {
	debugLog   LogWriter
	baseRouter BaseRouter
	mws        []MiddlewareFunc
}

// ServeHTTP implements http.Handler.
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.baseRouter.ServeHTTP(w, req)
}

// NewRouter returns a new Router.
func NewRouter(
	sw *spec.Swagger,
	handlers OperationHandlers,
	options ...RouterOption,
) (Router, error) {
	// Default options.
	router := Router{
		debugLog:   nil,
		baseRouter: defaultBaseRouter(),
		mws:        make([]MiddlewareFunc, 0),
	}

	// Apply argument options.
	for _, o := range options {
		o(&router)
	}

	// Router handles all the spec operations.
	base := router.baseRouter
	for method, pathOps := range analysis.New(sw).Operations() {
		for path, op := range pathOps {
			handler, ok := handlers[OperationID(op.ID)]
			if !ok {
				logf(router.debugLog, "oas2: no handler registered for operation %s", op.ID)
				continue
			}

			// Apply custom middleware before the operationIDMiddleware so
			// they can use the OptionID.
			for _, mwf := range router.mws {
				handler = mwf(handler)
			}

			// Add all path parameters to operation parameters.
			for _, pathParam := range sw.Paths.Paths[path].Parameters {
				op.AddParam(&pathParam)
			}

			logf(router.debugLog, "oas2: handle %s %s", method, sw.BasePath+path)
			handler = newOperationMiddleware(op).Apply(handler)
			base.Route(method, sw.BasePath+path, handler)
		}
	}

	return router, nil
}

// BaseRouter is an underlying router used in oas2 router.
// Any third-party router can be a BaseRouter by using adapter pattern.
type BaseRouter interface {
	http.Handler
	Route(method string, pathPattern string, handler http.Handler)
}

// LogWriter logs router operations that will be handled and what will be not
// during router creation. Useful for debugging.
type LogWriter func(format string, args ...interface{})

// RouterOption is an option for oas2 router.
type RouterOption func(*Router)

// DebugLog returns an option that sets a debug log for oas2 router.
// Debug log may help to see what router operations will be handled and what
// will be not.
func DebugLog(lw LogWriter) RouterOption {
	return func(args *Router) {
		args.debugLog = lw
	}
}

// Base returns an option that sets a BaseRouter for oas2 router.
// It allows to plug-in your favorite router to the oas2 router.
func Base(br BaseRouter) RouterOption {
	return func(args *Router) {
		args.baseRouter = br
	}
}

// Use returns an option that sets a middleware for router operations.
func Use(mw Middleware) RouterOption {
	return func(args *Router) {
		args.mws = append(args.mws, mw.Apply)
	}
}

// UseFunc returns an option that sets a middleware for router operations.
func UseFunc(mw MiddlewareFunc) RouterOption {
	return func(args *Router) {
		args.mws = append(args.mws, mw)
	}
}

func logf(w LogWriter, format string, args ...interface{}) {
	if w == nil {
		return
	}
	w(format, args...)
}
