package oas2

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

// GetOperationID returns OperationID from the request's context.
func GetOperationID(req *http.Request) OperationID {
	id, ok := req.Context().Value(contextKeyOperationID{}).(OperationID)
	if ok {
		return id
	}

	return OperationID("")
}

type contextKeyOperationID struct{}

func operationIDMiddleware(next http.Handler, id OperationID) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(
			context.WithValue(req.Context(), contextKeyOperationID{}, id),
		)
		next.ServeHTTP(w, req)
	})
}
