package oas2

import (
	"io/ioutil"
	"net/http"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

// NewRouter returns http.Handler that routes requests based on OAS 2.0 spec.
func NewRouter(
	sw *spec.Swagger,
	handlers OperationHandlers,
	options ...Option,
) (http.Handler, error) {
	// Default options.
	opts := Options{
		logger:     &logrus.Logger{Out: ioutil.Discard},
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
				opts.logger.Warnf("oas2 router: no handler registered for operation %s", op.ID)
				continue
			}

			// Apply custom middleware before the operationIDMiddleware so
			// they can use the OptionID.
			for _, mwf := range opts.mws {
				handler = mwf(handler)
			}

			opts.logger.Debugf("oas2 router: handle: %s %s", method, path)
			handler = operationIDMiddleware(handler, OperationID(op.ID))
			subrouter.Route(method, path, handler)
		}
	}

	// Mount the subrouter under the spec's basePath.
	router := opts.baseRouter
	router.Mount(sw.BasePath, subrouter)
	return router, nil
}

// Options is options for oas2 router.
type Options struct {
	logger     logrus.FieldLogger
	baseRouter BaseRouter
	mws        []MiddlewareFn
}

// Option is an option for oas2 router.
type Option func(*Options)

// LoggerOpt returns an option that sets a logger for oas2 router.
func LoggerOpt(logger logrus.FieldLogger) Option {
	return func(args *Options) {
		args.logger = logger
	}
}

// BaseRouterOpt returns an option that sets a BaseRouter for oas2 router.
// It allows to plug-in your favorite router to the oas2 router.
func BaseRouterOpt(br BaseRouter) Option {
	return func(args *Options) {
		args.baseRouter = br
	}
}

// MiddlewareOpt returns an option that sets a middleware for router operations.
func MiddlewareOpt(mw MiddlewareFn) Option {
	return func(args *Options) {
		args.mws = append(args.mws, mw)
	}
}

// BaseRouter is an underlying router used in oas2 router.
type BaseRouter interface {
	Route(method string, pathPattern string, handler http.Handler)
	Mount(path string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
