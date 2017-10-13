package oas2

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

// DecodeQuery decodes query parameters by their spec to the dst.
func DecodeQuery(ps []spec.Parameter, q url.Values, dst interface{}) error {
	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr {
		return fmt.Errorf("dst is not a pointer to struct (cannot modify)")
	}

	dv = dv.Elem()
	if dv.Kind() != reflect.Struct {
		return fmt.Errorf("dst is not a pointer to struct (cannot modify)")
	}

	fields := fieldMap(dv)

	for _, p := range ps {
		// No such tag in struct - no need to populate.
		f, ok := fields[p.Name]
		if !ok {
			continue
		}

		vals, ok := q[p.Name]
		if !ok {
			// No such value in query.
			if p.Default != nil {
				// Populate with default value.
				if err := set(p.Default, f, dv); err != nil {
					return err
				}
			}
			continue
		}

		// Convert value by type+format in parameter.
		v, err := convert.Parameter(vals, &p)
		if err != nil {
			return fmt.Errorf(
				"cannot use values %v as parameter %s with type %s and format %s",
				vals,
				p.Name,
				p.Type,
				p.Format,
			)
		}

		if err := set(v, f, dv); err != nil {
			return err
		}
	}

	return nil
}

func set(v interface{}, f reflect.StructField, dst reflect.Value) error {
	// Check if tag in struct can accept value of type v.
	if !f.Type.AssignableTo(reflect.TypeOf(v)) {
		return fmt.Errorf("field %s is not assignable to %s", f.Name, reflect.TypeOf(v).Name())
	}

	fieldVal := dst.FieldByName(f.Name)
	if !fieldVal.CanSet() {
		return fmt.Errorf("field %s of type %s is not settable", f.Name, dst.Type().Name())
	}

	fieldVal.Set(reflect.ValueOf(v))
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
