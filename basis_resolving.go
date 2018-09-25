package oas

import (
	"fmt"
	"net/http"
)

// Resolver resolves operation id from the request.
type Resolver interface {
	Resolve(req *http.Request) (string, bool)
}

// NewResolvingBasis returns a new resolving basis.
func NewResolvingBasis(doc *Document, resolver Resolver) *ResolvingBasis {
	b := &ResolvingBasis{
		doc:      doc,
		resolver: resolver,
		strict:   false,
	}

	b.initCache()
	return b
}

// A ResolvingBasis provides fundamental oas middleware, which resolve operation
// context from the request using the Resolver.
type ResolvingBasis struct {
	doc      *Document
	resolver Resolver
	cache    map[string]operationInfo

	// common options for derived middlewares

	strict bool
}

func (b *ResolvingBasis) initCache() {
	b.cache = make(map[string]operationInfo)
	// _ is method
	for _, pathOps := range b.doc.Analyzer.Operations() {
		// _ is path
		for _, operation := range pathOps {
			key := operation.ID
			value := operationInfo{
				operation: operation,
				params:    b.doc.Analyzer.ParametersFor(operation.ID),
				consumes:  b.doc.Analyzer.ConsumesFor(operation),
				produces:  b.doc.Analyzer.ProducesFor(operation),
			}
			b.cache[key] = value
		}
	}
}

// OperationContext returns a middleware that adds OpenAPI operation context to
// the request.
func (b *ResolvingBasis) OperationContext() Middleware {
	return func(next http.Handler) http.Handler {
		return &resolvingOperationContext{
			oc: &operationContext{
				next: next,
			},
			resolver: b.resolver,
			cache:    b.cache,
			strict:   b.strict,
		}
	}
}

// resolvingOperationContext is a middleware that resolves operation context
// from the request and adds operation info to the request context.
type resolvingOperationContext struct {
	oc       *operationContext
	resolver Resolver
	cache    map[string]operationInfo
	strict   bool
}

func (mw *resolvingOperationContext) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id, ok := mw.resolver.Resolve(req)
	if !ok {
		if mw.strict {
			panic("operation context middleware: cannot resolve operation id from the request")
		}
		mw.oc.ServeHTTP(w, req, operationInfo{}, false)
		return
	}

	oi, ok := mw.cache[id]
	if !ok {
		if mw.strict {
			panic(fmt.Sprintf("operation context middleware: cannot find operation info by the operation id %q", id))
		}
		mw.oc.ServeHTTP(w, req, operationInfo{}, false)
		return
	}

	mw.oc.ServeHTTP(w, req, oi, true)
}

// QueryValidator returns a middleware that validates request query parameters.
func (b *ResolvingBasis) QueryValidator(opts ...MiddlewareOption) Middleware {
	options := parseMiddlewareOptions(opts...)
	if options.problemHandler == nil {
		options.problemHandler = newProblemHandlerErrorResponder()
	}

	return func(next http.Handler) http.Handler {
		return &resolvingQueryValidator{
			qv: &queryValidator{
				next:              next,
				problemHandler:    options.problemHandler,
				continueOnProblem: options.continueOnProblem,
			},
			strict: b.strict,
		}
	}
}

// resolvingQueryValidator is a middleware that resolves operation context
// from the request and validates request query.
type resolvingQueryValidator struct {
	qv *queryValidator

	// strict enforces validation. If false, then validation is not
	// applied to requests without operation context.
	strict bool
}

func (mw *resolvingQueryValidator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oi, ok := getOperationInfo(req)
	if !ok {
		if mw.strict {
			panic("query validator middleware: cannot find operation info in the request context")
		}
		mw.qv.ServeHTTP(w, req, nil, false)
		return
	}

	mw.qv.ServeHTTP(w, req, oi.params, true)
}

// RequestContentTypeValidator returns a middleware that validates
// Content-Type header of the request.
//
// In case of validation error, this middleware will respond with
// either 406 or 415.
func (b *ResolvingBasis) RequestContentTypeValidator(opts ...MiddlewareOption) Middleware {
	return func(next http.Handler) http.Handler {
		return &resolvingRequestContentTypeValidator{
			rctv: &requestContentTypeValidator{
				next: next,
			},
			strict: b.strict,
		}
	}
}

type resolvingRequestContentTypeValidator struct {
	rctv *requestContentTypeValidator

	// strict enforces validation. If false, then validation is not
	// applied to requests without operation context.
	strict bool
}

func (mw *resolvingRequestContentTypeValidator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oi, ok := getOperationInfo(req)
	if !ok {
		if mw.strict {
			panic("request content type validator middleware: cannot find operation info in the request context")
		}
		mw.rctv.ServeHTTP(w, req, nil, nil, false)
		return
	}

	mw.rctv.ServeHTTP(w, req, oi.consumes, oi.produces, true)
}

// RequestBodyValidator returns a middleware that validates request body.
func (b *ResolvingBasis) RequestBodyValidator(opts ...MiddlewareOption) Middleware {
	options := parseMiddlewareOptions(opts...)
	if options.problemHandler == nil {
		options.problemHandler = newProblemHandlerErrorResponder()
	}

	return func(next http.Handler) http.Handler {
		return &resolvingRequestBodyValidator{
			rbv: &requestBodyValidator{
				next:              next,
				jsonSelectors:     options.jsonSelectors,
				problemHandler:    options.problemHandler,
				continueOnProblem: options.continueOnProblem,
			},
			strict: b.strict,
		}
	}
}

type resolvingRequestBodyValidator struct {
	rbv *requestBodyValidator

	// strict enforces validation. If false, then validation is not
	// applied to requests without operation context.
	strict bool
}

func (mw *resolvingRequestBodyValidator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oi, ok := getOperationInfo(req)
	if !ok {
		if mw.strict {
			panic("request body validator middleware: cannot find operation info in the request context")
		}
		mw.rbv.ServeHTTP(w, req, nil, false)
		return
	}

	mw.rbv.ServeHTTP(w, req, oi.params, true)
}

// ResponseContentTypeValidator returns a middleware that validates
// Content-Type header of the response.
func (b *ResolvingBasis) ResponseContentTypeValidator(opts ...MiddlewareOption) Middleware {
	options := parseMiddlewareOptions(opts...)
	if options.problemHandler == nil {
		options.problemHandler = newProblemHandlerWarnLogger("response")
	}

	return func(next http.Handler) http.Handler {
		return &resolvingResponseContentTypeValidator{
			rctv: &responseContentTypeValidator{
				next:           next,
				problemHandler: options.problemHandler,
			},
			strict: b.strict,
		}
	}
}

type resolvingResponseContentTypeValidator struct {
	rctv *responseContentTypeValidator

	// strict enforces validation. If false, then validation is not
	// applied to requests without operation context.
	strict bool
}

func (mw *resolvingResponseContentTypeValidator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oi, ok := getOperationInfo(req)
	if !ok {
		if mw.strict {
			panic("response content type validator middleware: cannot find operation info in the request context")
		}
		mw.rctv.ServeHTTP(w, req, nil, false)
		return
	}

	mw.rctv.ServeHTTP(w, req, oi.produces, true)
}

// ResponseBodyValidator returns a middleware that validates response body.
func (b *ResolvingBasis) ResponseBodyValidator(opts ...MiddlewareOption) Middleware {
	options := parseMiddlewareOptions(opts...)
	if options.problemHandler == nil {
		options.problemHandler = newProblemHandlerWarnLogger("response")
	}

	return func(next http.Handler) http.Handler {
		return &resolvingResponseBodyValidator{
			rbv: &responseBodyValidator{
				next:           next,
				jsonSelectors:  options.jsonSelectors,
				problemHandler: options.problemHandler,
			},
			strict: b.strict,
		}
	}
}

type resolvingResponseBodyValidator struct {
	rbv *responseBodyValidator

	// strict enforces validation. If false, then validation is not
	// applied to requests without operation context.
	strict bool
}

func (mw *resolvingResponseBodyValidator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oi, ok := getOperationInfo(req)
	if !ok {
		if mw.strict {
			panic("response body validator middleware: cannot find operation info in the request context")
		}
		mw.rbv.ServeHTTP(w, req, nil, false)
		return
	}

	mw.rbv.ServeHTTP(w, req, oi.operation.Responses, true)
}

// ContextualMiddleware represents a middleware that works based on request
// operation context.
type ContextualMiddleware interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request, op *Operation, ok bool)
}

// ContextualMiddleware returns a middleware that can work based on request
// operation context which will be resolved by the basis.
func (b *ResolvingBasis) ContextualMiddleware(m ContextualMiddleware) Middleware {
	return func(next http.Handler) http.Handler {
		return &resolvingContextualMiddleware{
			next:   m,
			strict: b.strict,
		}
	}
}

// resolvingContextualMiddleware is a contextual middleware that resolves
// operation context from the request.
type resolvingContextualMiddleware struct {
	next ContextualMiddleware

	// strict enforces validation. If false, then validation is not
	// applied to requests without operation context.
	strict bool
}

func (mw *resolvingContextualMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oi, ok := getOperationInfo(req)
	if !ok {
		if mw.strict {
			panic("contextual middleware: cannot find operation info in the request context")
		}
		mw.next.ServeHTTP(w, req, nil, false)
		return
	}

	mw.next.ServeHTTP(w, req, wrapOperation(oi.operation), true)
}
