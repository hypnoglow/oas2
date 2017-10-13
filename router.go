package oas2

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
)

// Router routes requests based on OAS 2.0 spec operations.
type Router interface {
	http.Handler
}

// NewRouter returns http.Handler that routes requests based on OAS 2.0 spec.
func NewRouter(
	sw *spec.Swagger,
	handlers OperationHandlers,
	options ...RouterOption,
) (Router, error) {
	// Default options.
	opts := RouterOptions{
		logger:     nil,
		baseRouter: defaultBaseRouter(),
		mws:        make([]MiddlewareFn, 0),
	}

	// Apply argument options.
	for _, o := range options {
		o(&opts)
	}

	// Subrouter handles all the spec operations.
	subrouter := opts.baseRouter
	for method, pathOps := range analysis.New(sw).Operations() {
		for path, op := range pathOps {
			handler, ok := handlers[OperationID(op.ID)]
			if !ok {
				logf(opts.logger, "oas2 router: no handler registered for operation %s", op.ID)
				continue
			}

			// Apply custom middleware before the operationIDMiddleware so
			// they can use the OptionID.
			for _, mwf := range opts.mws {
				handler = mwf(handler)
			}

			logf(opts.logger, "oas2 router: handle: %s %s", method, path)
			handler = NewOperationMiddleware(op).Apply(handler)
			subrouter.Route(method, path, handler)
		}
	}

	// Mount the subrouter under the spec's basePath.
	router := opts.baseRouter
	router.Mount(sw.BasePath, subrouter)
	return router, nil
}

// RouterOptions is options for oas2 router.
type RouterOptions struct {
	logger     io.Writer
	baseRouter BaseRouter
	mws        []MiddlewareFn
}

// RouterOption is an option for oas2 router.
type RouterOption func(*RouterOptions)

// LoggerOpt returns an option that sets a logger for oas2 router.
func LoggerOpt(logger io.Writer) RouterOption {
	return func(args *RouterOptions) {
		args.logger = logger
	}
}

// BaseRouterOpt returns an option that sets a BaseRouter for oas2 router.
// It allows to plug-in your favorite router to the oas2 router.
func BaseRouterOpt(br BaseRouter) RouterOption {
	return func(args *RouterOptions) {
		args.baseRouter = br
	}
}

// MiddlewareOpt returns an option that sets a middleware for router operations.
func MiddlewareOpt(mw Middleware) RouterOption {
	return func(args *RouterOptions) {
		args.mws = append(args.mws, mw.Apply)
	}
}

// MiddlewareFnOpt returns an option that sets a middleware for router operations.
func MiddlewareFnOpt(mw MiddlewareFn) RouterOption {
	return func(args *RouterOptions) {
		args.mws = append(args.mws, mw)
	}
}

// BaseRouter is an underlying router used in oas2 router.
type BaseRouter interface {
	Route(method string, pathPattern string, handler http.Handler)
	Mount(path string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func logf(w io.Writer, format string, args ...interface{}) {
	if w == nil {
		return
	}
	fmt.Fprintf(w, format, args)
}
