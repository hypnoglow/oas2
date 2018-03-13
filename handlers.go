package oas

import (
	"fmt"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/go-openapi/spec"
)

// SpecHandlerType represents spec handler type.
type SpecHandlerType int

const (
	// SpecHandlerTypeDynamic represents dynamic spec handler.
	SpecHandlerTypeDynamic SpecHandlerType = iota + 1

	// SpecHandlerTypeStatic represents static spec handler.
	SpecHandlerTypeStatic
)

// DynamicSpecHandler returns HTTP handler for OpenAPI spec that
// changes its host and schemes dynamically based on incoming request.
func DynamicSpecHandler(s *spec.Swagger) http.Handler {
	return &dynamicSpecHandler{s: s}
}

type dynamicSpecHandler struct {
	s *spec.Swagger
}

func (h *dynamicSpecHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = req.Host
	}

	scheme := req.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = req.Header.Get("X-Scheme")
		if scheme == "" {
			scheme = "http"
		}
	}

	specShallowCopy := &spec.Swagger{
		VendorExtensible: h.s.VendorExtensible,
		SwaggerProps:     h.s.SwaggerProps,
	}
	specShallowCopy.Host = host
	specShallowCopy.Schemes = []string{scheme}

	writeSpec(w, specShallowCopy)
}

// StaticSpecHandler returns HTTP handler for static OpenAPI spec.
func StaticSpecHandler(s *spec.Swagger) http.Handler {
	return &staticSpecHandler{s: s}
}

type staticSpecHandler struct {
	s *spec.Swagger
}

func (h *staticSpecHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	writeSpec(w, h.s)
}

func writeSpec(w http.ResponseWriter, s *spec.Swagger) {
	b, err := yaml.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	w.Write(b) // nolint
}
