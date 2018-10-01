package oas

import (
	"log"
	"net/http"
)

// NewProblem returns a new problem occurred while processing the request.
func NewProblem(w http.ResponseWriter, req *http.Request, err error) Problem {
	return Problem{
		w:   w,
		req: req,
		err: err,
	}
}

// Problem describes a problem occurred while processing the request (or the response).
// In most cases, the problem represents a validation error.
type Problem struct {
	w   http.ResponseWriter
	req *http.Request
	err error
}

// Cause returns the underlying error that represents the problem.
func (p Problem) Cause() error {
	return p.err
}

// ResponseWriter retruns the ResponseWriter relative to the request.
func (p Problem) ResponseWriter() http.ResponseWriter {
	return p.w
}

// Request returns the request on which the problem occured.
func (p Problem) Request() *http.Request {
	return p.req
}

// ProblemHandlerFunc is a function that handles problems occurred in a middleware
// while processing a request or a response.
//
// This function implements ProblemHandler.
type ProblemHandlerFunc func(Problem)

// HandleProblem handles the problem.
func (f ProblemHandlerFunc) HandleProblem(problem Problem) {
	f(problem)
}

// ProblemHandler can handle problems occurred in a middleware while processing
// a request or a response.
//
// Validation error depends on the middleware type, e.g. for query validator
// middleware the error will describe query validation failure. Usually, the
// handler should not wrap the error with a message like "query validation failure",
// because the message will be already present in such error.
type ProblemHandler interface {
	HandleProblem(problem Problem)
}

// newProblemHandlerErrorResponder is a very simple ProblemHandler that
// writes problem error message to the response.
func newProblemHandlerErrorResponder() ProblemHandlerFunc {
	return func(p Problem) {
		p.ResponseWriter().Header().Set("Content-Type", "text/plain; charset=utf-8")
		p.ResponseWriter().WriteHeader(http.StatusBadRequest)
		p.ResponseWriter().Write([]byte(p.err.Error())) // nolint
	}
}

// newProblemHandlerWarnLogger is a very simple ProblemHandler that writes
// problem error to the standard logger with a warning prefix.
func newProblemHandlerWarnLogger(kind string) ProblemHandlerFunc {
	return func(p Problem) {
		log.Printf("[WARN] oas %s problem on \"%s %s\": %v", kind, p.Request().Method, p.Request().URL.String(), p.Cause())
	}
}
