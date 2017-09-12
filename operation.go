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

// GetOperation returns *spec.Operation from the request's context.
func GetOperation(req *http.Request) *spec.Operation {
	op, ok := req.Context().Value(contextKeyOperation{}).(*spec.Operation)
	if ok {
		return op
	}

	return nil
}

type contextKeyOperation struct{}

func operationIDMiddleware(next http.Handler, op *spec.Operation) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(
			context.WithValue(req.Context(), contextKeyOperation{}, op),
		)
		next.ServeHTTP(w, req)
	})
}
