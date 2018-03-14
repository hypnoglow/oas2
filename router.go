package oas

import (
	"net/http"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// Router routes requests based on OAS 2.0 spec operations.
type Router struct {
	debugLog   LogWriter
	baseRouter BaseRouter
	mws        []Middleware

	// serveSpec, if nonzero, makes router serve its spec.
	serveSpec SpecHandlerType
}

// ServeHTTP implements http.Handler.
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.baseRouter.ServeHTTP(w, req)
}

// NewRouter returns a new Router.
func NewRouter(
	doc *loads.Document,
	handlers OperationHandlers,
	options ...RouterOption,
) (Router, error) {
	// Apply argument options.
	router := Router{}
	for _, o := range options {
		o(&router)
	}

	// Default options
	if router.debugLog == nil {
		router.debugLog = func(format string, args ...interface{}) {}
	}
	if router.baseRouter == nil {
		router.baseRouter = defaultBaseRouter()
	}

	// Router handles all the spec operations.
	base := router.baseRouter

	// Serve the specification itself if enabled.
	if router.serveSpec != 0 {
		var specHandler http.Handler
		switch router.serveSpec {
		case SpecHandlerTypeDynamic:
			specHandler = DynamicSpecHandler(doc.OrigSpec())
		case SpecHandlerTypeStatic:
			specHandler = StaticSpecHandler(doc.OrigSpec())
		}
		base.Route(http.MethodGet, doc.Spec().BasePath, specHandler)
	}

	for method, pathOps := range analysis.New(doc.Spec()).Operations() {
		for path, operation := range pathOps {
			handler, ok := handlers[OperationID(operation.ID)]
			if !ok {
				router.debugLog("oas: no handler registered for operation %s", operation.ID)
				continue
			}

			// Apply custom middleware before the operationIDMiddleware so
			// they can use the OptionID.
			for _, mwf := range router.mws {
				handler = mwf(handler)
			}

			// Copy operation to keep original operation unmodified.
			op := &spec.Operation{}
			if err := copyOperation(op, operation); err != nil {
				return router, err
			}

			// Add all path parameters to operation parameters
			// so operation in request context will be self-sufficient.
			// This is required for middlewares that use OpenAPI operation.
			for _, pathParam := range doc.Spec().Paths.Paths[path].Parameters {
				op.AddParam(&pathParam)
			}

			// Apply middleware to inject operation into every request
			// to make middlewares able to use it.
			handler = newOperationMiddleware(op)(handler)

			router.debugLog("oas: handle %s %s", method, doc.Spec().BasePath+path)
			base.Route(method, doc.Spec().BasePath+path, handler)
		}
	}

	return router, nil
}

// BaseRouter is an underlying router used in oas router.
// Any third-party router can be a BaseRouter by using adapter pattern.
type BaseRouter interface {
	http.Handler
	Route(method string, pathPattern string, handler http.Handler)
}

// LogWriter logs router operations that will be handled and what will be not
// during router creation. Useful for debugging.
type LogWriter func(format string, args ...interface{})

// RouterOption is an option for oas router.
type RouterOption func(*Router)

// DebugLog returns an option that sets a debug log for oas router.
// Debug log may help to see what router operations will be handled and what
// will be not.
func DebugLog(lw LogWriter) RouterOption {
	return func(args *Router) {
		args.debugLog = lw
	}
}

// Base returns an option that sets a BaseRouter for oa2 router.
// It allows to plug-in your favorite router to the oas router.
func Base(br BaseRouter) RouterOption {
	return func(args *Router) {
		args.baseRouter = br
	}
}

// Use returns an option that sets a middleware for router operations.
func Use(mw Middleware) RouterOption {
	return func(args *Router) {
		args.mws = append(args.mws, mw)
	}
}

// ServeSpec returns an option that makes router serve its spec.
func ServeSpec(t SpecHandlerType) RouterOption {
	return func(r *Router) {
		r.serveSpec = t
	}
}
