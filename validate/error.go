package validate

import "fmt"

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
