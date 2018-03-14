package oas

import (
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

// LoadSpecOption is an option for LoadSpec.
type LoadSpecOption struct {
	optionType optionType
	value      interface{}
}

// Validation is an option for LoadSpec that tells loader
// whether validation should be performed or not.
func Validation(val bool) *LoadSpecOption {
	return &LoadSpecOption{
		optionType: validationOption,
		value:      val,
	}
}

// LoadSpec opens an OpenAPI Specification v2.0 document from file.
func LoadSpec(fpath string, opts ...*LoadSpecOption) (document *loads.Document, err error) {
	validate := true

	for _, opt := range opts {
		if val, ok := opt.value.(bool); ok && opt.optionType == validationOption {
			validate = val
		}
	}

	// Load regularly.

	document, err = loads.Spec(fpath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load spec")
	}

	if validate {
		if err := ValidateSpec(fpath); err != nil {
			return nil, errors.Wrap(err, "Spec is invalid")
		}
	}

	return document, nil
}

// ValidateSpec opens an OpenAPI Specification v2.0 document from file, expands all
// references within it, then validates the spec.
func ValidateSpec(fpath string) error {
	document, err := loads.Spec(fpath)
	if err != nil {
		return errors.Wrap(err, "failed to load spec")
	}

	document, err = document.Expanded(&spec.ExpandOptions{RelativeBase: fpath})
	if err != nil {
		return errors.Wrap(err, "failed to expand spec")
	}

	if err = validate.Spec(document, strfmt.Default); err != nil {
		return errors.Wrap(err, "spec is invalid")
	}

	return nil
}

type optionType string

const (
	validationOption optionType = "skip validation"
)
