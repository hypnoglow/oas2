package oas2

import (
	"net/http"
	"io/ioutil"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

type (
	// OperationID is an operation identifier.
	OperationID string

	// OperationHandlers maps OperationID to its handler.
	OperationHandlers map[OperationID]http.Handler
)

// NewRouter returns http.Handler that routes requests based on OAS 2.0 spec.
func NewRouter(
	sw *spec.Swagger,
	handlers OperationHandlers,
	options ...Option,
) (http.Handler, error) {
	// Default options.
	opts := Options{
		logger:            &logrus.Logger{Out: ioutil.Discard},
		baseRouterFactory: defaultBaseRouterFactory(),
	}

	// Apply argument options.
	for _, o := range options {
		o(&opts)
	}

	// Subrouter handles all the spec operations.
	subrouter := opts.baseRouterFactory()
	for method, pathOps := range analysis.New(sw).Operations() {
		for path, op := range pathOps {
			handler, ok := handlers[OperationID(op.ID)]
			if !ok {
				opts.logger.Warnf("oas3 router: no handler registered for operation %s", op.ID)
				continue
			}

			opts.logger.Debugf("oas3 router: handle: %s %s", method, path)
			subrouter.Route(method, path, handler)
		}
	}

	// Mount the subrouter under the spec's basePath.
	router := opts.baseRouterFactory()
	router.Mount(sw.BasePath, subrouter)
	return router, nil
}

// Options is options for oas2 router.
type Options struct {
	logger            logrus.FieldLogger
	baseRouterFactory func() BaseRouter
}

// Option is an option for oas2 router.
type Option func(*Options)

// Logger returns an option that sets a logger for oas2 router.
func Logger(logger logrus.FieldLogger) Option {
	return func(args *Options) {
		args.logger = logger
	}
}

// BaseRouterFactory returns an option that sets a BaseRouter factory for oas2
// router. It allows to plug-in your favorite router to the oas2 router.
func BaseRouterFactory(factory func() BaseRouter) Option {
	return func(args *Options) {
		args.baseRouterFactory = factory
	}
}

// BaseRouter is an underlying router used in oas2 router.
type BaseRouter interface {
	Route(method string, pathPattern string, handler http.Handler)
	Mount(path string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
