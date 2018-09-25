package gorilla

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hypnoglow/oas2"
)

// NewPathParamExtractor returns a new path param extractor that extracts
// path parameter from the request using gorilla/mux route context.
func NewPathParamExtractor() oas.PathParamExtractor {
	return &pathParamsExtractor{}
}

type pathParamsExtractor struct{}

// PathParam returns path parameter by key from the request context.
func (e pathParamsExtractor) PathParam(req *http.Request, key string) string {
	return mux.Vars(req)[key]
}
