package oas

import (
	"context"
	"net/http"
)

// OperationID is an operation identifier.
type OperationID string

// String implements fmt.Stringer interface.
func (oid OperationID) String() string {
	return string(oid)
}

// OperationHandlers maps OperationID to its handler.
type OperationHandlers map[OperationID]http.Handler

func newOperationMiddleware(op *Operation) Middleware {
	return operationMiddleware{op}.chain
}

type operationMiddleware struct {
	operation *Operation
}

func (m operationMiddleware) chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = WithOperation(req, m.operation)
		next.ServeHTTP(w, req)
	})
}

// WithOperation returns request with context value defining *spec.Operation.
func WithOperation(req *http.Request, op *Operation) *http.Request {
	return req.WithContext(
		context.WithValue(req.Context(), contextKeyOperation{}, op),
	)
}

// GetOperation returns *spec.Operation from the request's context.
// In case of operation not found GetOperation returns nil.
func GetOperation(req *http.Request) *Operation {
	op, ok := req.Context().Value(contextKeyOperation{}).(*Operation)
	if ok {
		return op
	}

	return nil
}

// MustOperation returns *spec.Operation from the request's context.
// In case of operation not found MustOperation panics.
func MustOperation(req *http.Request) *Operation {
	op := GetOperation(req)
	if op != nil {
		return op
	}

	panic("request has no OpenAPI operation spec in its context")
}

type contextKeyOperation struct{}
