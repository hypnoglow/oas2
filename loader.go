package oas2

import (
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

// LoadSpec opens an OpenAPI Specification v2.0 document, expands all references within it,
// then validates the spec and returns spec document.
func LoadSpec(path string) (*loads.Document, error) {
	document, err := loads.Spec(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load spec")
	}

	if err := spec.ExpandSpec(document.Spec(), &spec.ExpandOptions{RelativeBase: path}); err != nil {
		return nil, errors.Wrap(err, "Failed to expand spec")
	}

	if err := validate.Spec(document, strfmt.Default); err != nil {
		return nil, errors.Wrap(err, "Spec is invalid")
	}

	return document, nil
}
