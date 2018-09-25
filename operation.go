package oas

import (
	"context"
	"net/http"

	"github.com/go-openapi/spec"
)

type operationInfo struct {
	operation *spec.Operation

	// params include all applicable operation params, even those defined
	// on the path operation belongs to.
	params []spec.Parameter

	// consumes is either operation-defined "consumes" property or spec-wide
	// "consumes" property.
	consumes []string

	// produces is either operation-defined "produces" property or spec-wide
	// "produces" property.
	produces []string
}

// operationContext is a middleware that adds operation info to the request
// context and calls next.
type operationContext struct {
	next http.Handler
}

func (mw *operationContext) ServeHTTP(w http.ResponseWriter, req *http.Request, oi operationInfo, ok bool) {
	if ok {
		req = withOperationInfo(req, oi)
	}

	mw.next.ServeHTTP(w, req)
}

type contextKeyOperationInfo struct{}

// withOperationInfo returns request with context value defining *spec.Operation.
func withOperationInfo(req *http.Request, info operationInfo) *http.Request {
	return req.WithContext(
		context.WithValue(req.Context(), contextKeyOperationInfo{}, info),
	)
}

// getOperationInfo returns *spec.Operation from the request's context.
// In case of operation not found GetOperation returns nil.
func getOperationInfo(req *http.Request) (operationInfo, bool) {
	op, ok := req.Context().Value(contextKeyOperationInfo{}).(operationInfo)
	return op, ok
}

// mustOperationInfo returns *spec.Operation from the request's context.
// In case of operation not found MustOperation panics.
//
// nolint
func mustOperationInfo(req *http.Request) operationInfo {
	op, ok := getOperationInfo(req)
	if ok {
		return op
	}

	panic("request has no OpenAPI operation spec in its context")
}
