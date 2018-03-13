package oas

import "regexp"

var (
	contentTypeSelectorRegexJSON    *regexp.Regexp
	contentTypeSelectorRegexJSONAPI *regexp.Regexp
)

func init() {
	contentTypeSelectorRegexJSON = regexp.MustCompile(`^application\/json`)
	contentTypeSelectorRegexJSONAPI = regexp.MustCompile(`^application\/vnd\.api\+json$`)
}

// MiddlewareOptions represent options for middleware.
type MiddlewareOptions struct {
	contentTypeRegexSelectors []*regexp.Regexp
}

// MiddlewareOption represent option for middleware.
type MiddlewareOption func(*MiddlewareOptions)

// ContentTypeRegexSelector select requests/responses based on Content-Type
// header. If any selector matches Content-Type of the request/response, then it
// will be validated. Otherwise, validator skips validation of the request/response.
//
// This options can be applied to the following middlewares:
// - BodyValidator
// - ResponseBodyValidator
func ContentTypeRegexSelector(selector *regexp.Regexp) MiddlewareOption {
	return func(opts *MiddlewareOptions) {
		opts.contentTypeRegexSelectors = append(opts.contentTypeRegexSelectors, selector)
	}
}
