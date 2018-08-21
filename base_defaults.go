package oas

import "github.com/go-chi/chi"

// DefaultBaseRouter is the base router used by default when no specific
// base router passed to NewRouter().
func DefaultBaseRouter() BaseRouter {
	return ChiRouter(chi.NewRouter())
}

// DefaultPathTemplateFunc is the path template func that should be used with
// the default router.
var DefaultPathTemplateFunc = ChiPathTemplate

// DefaultPathParamFunc is the path parameter extractor func that should be used
// with the default router.
var DefaultPathParamFunc = ChiPathParam

// DefaultExtractorFunc is the path parameter extractor func that should be used
// with the default router. This function is kept for compatibility reasons.
// Deprecated.
var DefaultExtractorFunc = ChiPathParam
