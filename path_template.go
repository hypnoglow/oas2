package oas

import (
	"context"
	"net/http"
)

// PathTemplateFunc defines function that returns path template from the request.
//
// For example, if a route "/fruits/{name}" was registered in the router,
// and the current request has path "/fruits/apple", then this function
// should return "/fruits/{name}".
type PathTemplateFunc func(req *http.Request) string

type contextKeyPathTemplate struct{}

func WithPathTemplate(req *http.Request, pathTemplate string) *http.Request {
	return req.WithContext(
		context.WithValue(
			req.Context(), contextKeyPathTemplate{}, pathTemplate,
		),
	)
}

func GetPathTemplate(req *http.Request) (string, bool) {
	pt, ok := req.Context().Value(contextKeyPathTemplate{}).(string)
	return pt, ok
}

func MustPathTemplate(req *http.Request) string {
	pt, ok := GetPathTemplate(req)
	if ok {
		return pt
	}

	panic("request has no path template in its context")
}
