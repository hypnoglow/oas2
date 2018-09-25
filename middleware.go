package oas

import (
	"net/http"
	"regexp"
)

// Middleware describes a middleware that can be applied to a http.handler.
type Middleware func(next http.Handler) http.Handler

// MiddlewareOptions represent options for middleware.
type MiddlewareOptions struct {
	jsonSelectors     []*regexp.Regexp
	problemHandler    ProblemHandler
	continueOnProblem bool
}

// MiddlewareOption represent option for middleware.
type MiddlewareOption func(*MiddlewareOptions)

// WithJSONSelectors returns a middleware option that sets JSON Content-Type selectors.
func WithJSONSelectors(selectors ...*regexp.Regexp) MiddlewareOption {
	return func(opts *MiddlewareOptions) {
		opts.jsonSelectors = append(opts.jsonSelectors, selectors...)
	}
}

// WithProblemHandler returns a middleware option that sets problem handler.
func WithProblemHandler(h ProblemHandler) MiddlewareOption {
	return func(opts *MiddlewareOptions) {
		opts.problemHandler = h
	}
}

// WithProblemHandlerFunc returns a middleware option that sets problem handler.
func WithProblemHandlerFunc(f ProblemHandlerFunc) MiddlewareOption {
	return func(opts *MiddlewareOptions) {
		opts.problemHandler = f
	}
}

// WithContinueOnProblem returns a middleware option that defines if middleware
// should continue when error occurs.
func WithContinueOnProblem(continue_ bool) MiddlewareOption {
	return func(opts *MiddlewareOptions) {
		opts.continueOnProblem = continue_
	}
}

func parseMiddlewareOptions(opts ...MiddlewareOption) MiddlewareOptions {
	options := MiddlewareOptions{
		jsonSelectors:     nil,
		continueOnProblem: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.jsonSelectors == nil {
		defaultJSONSelectors()(&options)
	}

	return options
}

func defaultJSONSelectors() MiddlewareOption {
	return func(opts *MiddlewareOptions) {
		opts.jsonSelectors = []*regexp.Regexp{
			contentTypeSelectorRegexJSON,
			contentTypeSelectorRegexJSONAPI,
		}
	}
}
