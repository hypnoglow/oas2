// Package validate provides utilities that allow to validate request and
// response data against OpenAPI Specification parameter and schema definitions.
//
// Note that errors returned from validation functions are generally of type
// Error, so they can be asserted to corresponding interface(s) to retrieve
// error's field and value.
//  errs := validate.Query(ps, q)
//  for _, err := range errs {
//      if e, ok := err.(validate.Error) {
//          field, value := e.Field(), e.Value()
//          // ...
//      }
//  }
package validate

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/validate"

	"github.com/hypnoglow/oas2/convert"
)

// Query validates request query parameters by spec and returns errors
// if any.
func Query(ps []spec.Parameter, q url.Values) []error {
	errs := make(ValidationErrors, 0)

	// Iterate over spec parameters and validate each against the spec.
	for _, p := range ps {
		if p.In != "query" {
			// Validating only "query" parameters.
			continue
		}

		errs = append(errs, validateQueryParam(p, q)...)

		delete(q, p.Name) // to check not described parameters passed
	}

	// Check that no additional parameters passed.
	for name := range q {
		errs = append(errs, ValidationErrorf(name, q.Get(name), "parameter %s is unknown", name))
	}

	return errs.Errors()
}

// Body validates request body by spec and returns errors if any.
func Body(ps []spec.Parameter, data interface{}) []error {
	errs := make(ValidationErrors, 0)

	for _, p := range ps {
		if p.In != "body" {
			// Validating only "body" parameters.
			continue
		}

		errs = append(errs, validateBodyParam(p, data)...)
	}

	return errs.Errors()
}

// BySchema validates data by spec and returns errors if any.
func BySchema(sch *spec.Schema, data interface{}) []error {
	return validatebySchema(sch, data).Errors()
}

// ValidationError describes validation error.
type ValidationError interface {
	error

	// Field returns field name where error occurred.
	Field() string

	// Value returns original value passed by client on field where error
	// occurred.
	Value() interface{}
}

// ValidationErrorf returns a new formatted ValidationError.
func ValidationErrorf(field string, value interface{}, format string, args ...interface{}) ValidationError {
	return valErr{
		message: fmt.Sprintf(format, args...),
		field:   field,
		value:   value,
	}
}

// ValidationErrors is a set of validation errors.
type ValidationErrors []ValidationError

// Errors returns ValidationErrors in form of Go builtin errors.
func (es ValidationErrors) Errors() []error {
	if len(es) == 0 {
		return nil
	}

	errs := make([]error, len(es))
	for i, e := range es {
		errs[i] = e
	}
	return errs
}

func validateQueryParam(p spec.Parameter, q url.Values) (errs ValidationErrors) {
	_, ok := q[p.Name]
	if !ok {
		if p.Required {
			errs = append(errs, ValidationErrorf(p.Name, nil, "param %s is required", p.Name))
		}
		return errs
	}

	value, err := convert.Parameter(q[p.Name], &p)
	if err != nil {
		// TODO: q.Get(p.Name) relies on type that is not array/file.
		return append(errs, ValidationErrorf(p.Name, q.Get(p.Name), "param %s: %s", p.Name, err))
	}

	if result := validate.NewParamValidator(&p, formatRegistry).Validate(value); result != nil {
		for _, e := range result.Errors {
			errs = append(errs, ValidationErrorf(p.Name, value, e.Error()))
		}
	}

	return errs
}

func validateBodyParam(p spec.Parameter, data interface{}) (errs ValidationErrors) {
	return validatebySchema(p.Schema, data)
}

func validatebySchema(sch *spec.Schema, data interface{}) (errs ValidationErrors) {
	err := validate.AgainstSchema(sch, data, formatRegistry)
	ves, ok := err.(*errors.CompositeError)
	if ok && len(ves.Errors) > 0 {
		for _, e := range ves.Errors {
			ve := e.(*errors.Validation)
			errs = append(errs, ValidationErrorf(strings.TrimPrefix(ve.Name, "."), nil, strings.TrimPrefix(ve.Error(), ".")))
		}
	}

	return errs
}

// valErr implements ValidationError.
type valErr struct {
	message string
	field   string
	value   interface{}
}

func (v valErr) Error() string {
	return v.message
}

func (v valErr) Field() string {
	return v.field
}

func (v valErr) Value() interface{} {
	return v.value
}
