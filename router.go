package oas

import (
	"net/http"
)

// NewRouter returns a new Router that can route requests by an OpenAPI spec.
func NewRouter(options ...RouterOption) *Router {
	router := &Router{
		debugLog:   func(format string, args ...interface{}) {},
		baseRouter: DefaultBaseRouter(),
		ptf:        DefaultPathTemplateFunc,
		mws:        nil,
		serveSpec:  0,
	}

	for _, o := range options {
		o(router)
	}

	return router
}

// RouterOption is an option for oas router.
type RouterOption func(*Router)

// LogWriter logs router operations that will be handled and what will be not
// during router creation. Useful for debugging.
type LogWriter func(format string, args ...interface{})

// DebugLog returns an option that sets a debug log for oas router.
// Debug log may help to see what router operations will be handled and what
// will be not.
func DebugLog(lw LogWriter) RouterOption {
	return func(args *Router) {
		args.debugLog = lw
	}
}

// Base returns an option that sets a BaseRouter and a PathTemplateFunc for oas2
// router. It allows to plug-in your favorite router to the oas router.
func Base(br BaseRouter, ptf PathTemplateFunc) RouterOption {
	return func(args *Router) {
		args.baseRouter = br
		args.ptf = ptf
	}
}

// RouterMiddleware returns an option that sets middlewares for the router.
//
// These middlewares are applied to the router itself. All specs added to this
// router will use these middlewares along with spec-scoped ones.
//
// These middlewares will can access path template, operation and params in the
// request context.
//
// Multiple middlewares will be executed exactly in the same order
// they were passed to the router. For example:
//  router := oas.NewRouter(
//      oas.RouterMiddleware(CORS),
//      oas.RouterMiddleware(RequestID, RequestLogger),
//  )
// Here the CORS one will be executed first, then RequestID and then RequestLogger.
// Thus, RequestLogger will be able to use request id that RequestID middleware
// stored in a request context.
func RouterMiddleware(mw ...Middleware) RouterOption {
	return func(r *Router) {
		r.mws = append(r.mws, mw...)
	}
}

// ServeSpec returns an option that makes router serve its spec.
func ServeSpec(t SpecHandlerType) RouterOption {
	return func(r *Router) {
		r.serveSpec = t
	}
}

// Router routes requests based on OAS 2.0 spec operations.
type Router struct {
	debugLog LogWriter

	baseRouter BaseRouter
	ptf        PathTemplateFunc

	mws []Middleware

	// serveSpec, if nonzero, makes router serve its spec.
	serveSpec SpecHandlerType
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.baseRouter.ServeHTTP(w, req)
}

// AddSpec adds routes from the spec to the router.
func (r *Router) AddSpec(doc *Document, handlers OperationHandlers, opts ...SpecOption) error {
	options := routeSpecOptions{}
	for _, o := range opts {
		o(&options)
	}

	return r.applySpecRouting(doc, handlers, options.mws)
}

type routeSpecOptions struct {
	mws []Middleware
}

// SpecOption is an option for spec routing.
type SpecOption func(*routeSpecOptions)

// SpecMiddleware returns an option that sets middlewares for the spec.
//
// These middlewares are applied only to the spec and its routes. Other specs
// will use their own middlewares, not these ones.
//
// These middlewares will can access path template, operation and params in the
// request context.
//
// Multiple middlewares will be executed exactly in the same order
// they were passed. For example:
//  _ = router.AddSpec(
//      doc, operationHandlers,
//      oas.SpecMiddleware(CORS),
//      oas.SpecMiddleware(RequestID, RequestLogger),
//  )
// Here the CORS one will be executed first, then RequestID and then RequestLogger.
// Thus, RequestLogger will be able to use request id that RequestID middleware
// stored in a request context.
func SpecMiddleware(mw ...Middleware) SpecOption {
	return func(o *routeSpecOptions) {
		o.mws = append(o.mws, mw...)
	}
}

func (r *Router) applySpecRouting(doc *Document, handlers OperationHandlers, mws []Middleware) error {
	var routes []Route

	// Serve the specification itself if enabled.
	if r.serveSpec != 0 {
		var specHandler http.Handler
		switch r.serveSpec {
		case SpecHandlerTypeDynamic:
			specHandler = DynamicSpecHandler(doc.OrigSpec())
		case SpecHandlerTypeStatic:
			specHandler = StaticSpecHandler(doc.OrigSpec())
		}
		path := doc.Spec().BasePath
		if path == "" {
			path = "/"
		}

		routes = append(routes, Route{
			Method:  http.MethodGet,
			Path:    path,
			Handler: specHandler,
		})
	}

	for method, pathOps := range doc.Analyzer.Operations() {
		for path, operation := range pathOps {
			handler, ok := handlers[OperationID(operation.ID)]
			if !ok {
				r.debugLog("oas: no handler registered for operation %s", operation.ID)
				continue
			}

			// Wrap the handler with all middleware provided by user so that
			// middleware handlers will be executed exactly in the same order
			// they there passed to the router.
			//for i := range r.mws {
			//	handler = r.mws[len(r.mws)-1-i](handler)
			//}

			r.debugLog("oas: handle %s %s", method, doc.Spec().BasePath+path)
			routes = append(routes, Route{
				Method:  method,
				Path:    doc.Spec().BasePath + path,
				Handler: handler,
			})
		}
	}

	middlewareStack := append([]Middleware{newMultiMiddleware(r.ptf, doc)}, r.mws...)
	middlewareStack = append(middlewareStack, mws...)
	r.baseRouter.Compose(middlewareStack, routes)

	return nil
}
