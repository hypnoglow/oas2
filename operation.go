package oas2

import (
	"context"
	"net/http"

	"github.com/go-openapi/spec"
)

// OperationID is an operation identifier.
type OperationID string

// String implements fmt.Stringer interface.
func (oid OperationID) String() string {
	return string(oid)
}

// OperationHandlers maps OperationID to its handler.
type OperationHandlers map[OperationID]http.Handler

func newOperationMiddleware(op *spec.Operation) Middleware {
	return operationMiddleware{op}
}

type operationMiddleware struct {
	operation *spec.Operation
}

func (m operationMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = withOperation(req, m.operation)
		next.ServeHTTP(w, req)
	})
}

func withOperation(req *http.Request, op *spec.Operation) *http.Request {
	return req.WithContext(
		context.WithValue(req.Context(), contextKeyOperation{}, op),
	)
}

// GetOperation returns *spec.Operation from the request's context.
// In case of operation not found GetOperation returns nil.
func GetOperation(req *http.Request) *spec.Operation {
	op, ok := req.Context().Value(contextKeyOperation{}).(*spec.Operation)
	if ok {
		return op
	}

	return nil
}

// MustOperation returns *spec.Operation from the request's context.
// In case of operation not found MustOperation panics.
func MustOperation(req *http.Request) *spec.Operation {
	op := GetOperation(req)
	if op != nil {
		return op
	}

	panic("OAS operation not found in request context")
}

type contextKeyOperation struct{}
