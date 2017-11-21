package oas

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/go-openapi/spec"

	"github.com/hypnoglow/oas2/convert"
)

const (
	tag = "oas"
)

// reflectValueOfDestination ensures that destination is a pointer to a struct
// and returns reflect.Value of destination struct.
func reflectValueOfDestination(dst interface{}) (reflect.Value, error) {
	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr {
		return dv, fmt.Errorf("dst is not a pointer to struct (cannot modify)")
	}

	dv = dv.Elem()
	if dv.Kind() != reflect.Struct {
		return dv, fmt.Errorf("dst is not a pointer to struct (cannot modify)")
	}

	return dv, nil
}

func setFieldValue(v interface{}, f reflect.StructField, dst reflect.Value) error {
	isPointer := false
	fieldType := f.Type

	valueType := reflect.TypeOf(v)
	valueValue := reflect.ValueOf(v)

	if f.Type.Kind() == reflect.Ptr && valueType.Kind() != reflect.Ptr {
		isPointer = true
		fieldType = f.Type.Elem()
	}

	// Get field value and check if it is settable.
	fieldVal := dst.FieldByName(f.Name)
	if !fieldVal.CanSet() {
		return fmt.Errorf("field %s of type %s is not settable", f.Name, dst.Type().Name())
	}

	// Check if tag in struct can accept value of type v.
	if !valueType.AssignableTo(fieldType) {
		return fmt.Errorf("value of type %s is not assignable to field %s of type %s", reflect.TypeOf(v).String(), f.Name, f.Type.String())
	}

	// Set the value. Pay attention to pointers.

	if isPointer {
		fieldVal.Set(reflect.New(fieldType))
		fieldVal.Elem().Set(valueValue)
		return nil
	}

	fieldVal.Set(valueValue)
	return nil
}

// fieldMap returns v fields mapped by their tags.
func fieldMap(rv reflect.Value) map[string]reflect.StructField {
	rt := rv.Type()

	m := make(map[string]reflect.StructField)
	n := rt.NumField()
	for i := 0; i < n; i++ {
		f := rt.Field(i)
		tag, ok := f.Tag.Lookup(tag)
		if !ok {
			continue
		}

		m[tag] = f
	}

	return m
}

// fromValues fetches and converts a value by parameter from url.Values.
func fromValues(values url.Values, p spec.Parameter) (interface{}, error) {
	vals, ok := values[p.Name]
	if !ok {
		// No such value in query/form.
		// This is not an error because validator should mark request
		// as invalid if this parameter is required.

		// Just return default value.
		// It may be nil, which means no default value.
		return p.Default, nil
	}

	// Convert value by type+format in parameter.
	v, err := convert.Parameter(vals, &p)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot use values %v as parameter %s with type %s and format %s",
			vals,
			p.Name,
			p.Type,
			p.Format,
		)
	}

	return v, nil
}
