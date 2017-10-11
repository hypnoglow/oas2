package utils

import (
	"bytes"
	"io"
	"net/http"
)

// ResponseRecorder is a http.ResponseWriter that provides a way to fetch
// written status and payload.
type ResponseRecorder interface {
	http.ResponseWriter
	Status() int
	Payload() []byte
}

type responseRecorder struct {
	origin        http.ResponseWriter
	status        int
	statusWritten bool
	payload       *bytes.Buffer
}

// NewResponseRecorder returns a new ResponseRecorder.
func NewResponseRecorder(origin http.ResponseWriter) ResponseRecorder {
	return &responseRecorder{
		origin:        origin,
		status:        http.StatusOK,
		statusWritten: false,
		payload:       new(bytes.Buffer),
	}
}

func (r *responseRecorder) Header() http.Header {
	return r.origin.Header()
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return io.MultiWriter(r.origin, r.payload).Write(b)
}

func (r *responseRecorder) WriteHeader(status int) {
	if !r.statusWritten {
		r.status = status
	}
	r.origin.WriteHeader(status)
}

func (r *responseRecorder) Status() int {
	return r.status
}

func (r *responseRecorder) Payload() []byte {
	return r.payload.Bytes()
}
