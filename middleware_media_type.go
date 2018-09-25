package oas

import (
	"fmt"
	"net/http"
)

// matchMediaType checks if media type matches any allowed media type.
func matchMediaType(mediaType string, allowed []string) bool {
	if len(allowed) == 0 {
		// If no media types are explicitly defined, allow all.
		return true
	}

	if mediaType == "" {
		// If no media type given, consider it is ok.
		// This is useful for HTTP 204 responses, as well as for
		// assuming "defaults". RFC 7231 does not strictly requires
		// Content-Type to be defined: https://tools.ietf.org/html/rfc7231#section-3.1.1.5
		return true
	}

	for _, a := range allowed {
		if a == mediaTypeWildcard {
			return true
		}
		if mediaType == a {
			return true
		}
	}

	return false
}

// matchMediaTypes checks if any of the media types matches any allowed
// media type.
func matchMediaTypes(mediaTypes []string, allowed []string) bool {
	if len(allowed) == 0 {
		// If no media types are explicitly defined, allow all.
		return true
	}

	if len(mediaTypes) == 0 {
		// If no media types are explicitly requested, allow all.
		return true
	}

	for _, mediaType := range mediaTypes {
		if mediaType == mediaTypeWildcard {
			return true
		}
		for _, a := range allowed {
			if a == mediaTypeWildcard {
				return true
			}
			if a == mediaType {
				return true
			}
		}
	}

	return false
}

// requestContentTypeValidator validates request against media types which
// are defined by the corresponding operation.
type requestContentTypeValidator struct {
	next http.Handler
}

func (mw *requestContentTypeValidator) ServeHTTP(w http.ResponseWriter, req *http.Request, consumes []string, produces []string, ok bool) {
	if !ok {
		mw.next.ServeHTTP(w, req)
		return
	}

	if req.ContentLength > 0 {
		if !matchMediaType(req.Header.Get("Content-Type"), consumes) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
	}

	if !matchMediaTypes(req.Header["Accept"], produces) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	mw.next.ServeHTTP(w, req)
}

// requestContentTypeValidator validates response against media types which
// are defined by the corresponding operation.
type responseContentTypeValidator struct {
	next http.Handler

	problemHandler ProblemHandler
}

func (mw *responseContentTypeValidator) ServeHTTP(w http.ResponseWriter, req *http.Request, produces []string, ok bool) {
	mw.next.ServeHTTP(w, req)
	if !ok {
		return
	}

	ct := w.Header().Get("Content-Type")

	if !matchMediaType(ct, req.Header["Accept"]) {
		err := fmt.Errorf("Content-Type header of the response does not match Accept header of the request")
		mw.problemHandler.HandleProblem(NewProblem(w, req, err))
	}

	if !matchMediaType(ct, produces) {
		err := fmt.Errorf("Content-Type header of the response does not match any of the media types the operation can produce")
		mw.problemHandler.HandleProblem(NewProblem(w, req, err))
	}
}
