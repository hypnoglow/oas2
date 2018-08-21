package oas

import (
	"context"
	"net/http"

	"github.com/go-openapi/spec"
)

type contextKeyParams struct{}

func WithParams(req *http.Request, params []spec.Parameter) *http.Request {
	return req.WithContext(
		context.WithValue(req.Context(), contextKeyParams{}, params),
	)
}

func GetParams(req *http.Request) ([]spec.Parameter, bool) {
	params, ok := req.Context().Value(contextKeyParams{}).([]spec.Parameter)
	return params, ok
}

func MustParams(req *http.Request) []spec.Parameter {
	params, ok := GetParams(req)
	if ok {
		return params
	}

	panic("request has no OpenAPI parameters in its context")
}
