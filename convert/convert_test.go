package convert

import (
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

func TestParameter(t *testing.T) {
	t.Run("ok for primitive type", func(t *testing.T) {
		values := []string{"John"}
		param := spec.QueryParam("name").Typed("string", "")

		v, err := Parameter(values, param)
		assertConversionResult(t, "John", v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for string array", func(t *testing.T) {
		values := []string{"Nicolas", "Jonathan"}
		param := spec.QueryParam("names").Typed("array", "")
		param.Items = spec.NewItems().Typed("string", "")
		param.SimpleSchema.CollectionFormat = "multi"

		v, err := Parameter(values, param)
		assertConversionResult(t, []string{"Nicolas", "Jonathan"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for space-separated array in string", func(t *testing.T) {
		values := []string{"Peter Lois"}
		param := spec.QueryParam("names").Typed("array", "")
		param.Items = spec.NewItems().Typed("string", "")
		param.SimpleSchema.CollectionFormat = "ssv"

		v, err := Parameter(values, param)
		assertConversionResult(t, []string{"Peter", "Lois"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for tab-separated array in string", func(t *testing.T) {
		values := []string{"Brian\tStewie"}
		param := spec.QueryParam("names").Typed("array", "")
		param.Items = spec.NewItems().Typed("string", "")
		param.SimpleSchema.CollectionFormat = "tsv"

		v, err := Parameter(values, param)
		assertConversionResult(t, []string{"Brian", "Stewie"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for pipe-separated array in string", func(t *testing.T) {
		values := []string{"Meg|Chris"}
		param := spec.QueryParam("names").Typed("array", "")
		param.Items = spec.NewItems().Typed("string", "")
		param.SimpleSchema.CollectionFormat = "pipes"

		v, err := Parameter(values, param)
		assertConversionResult(t, []string{"Meg", "Chris"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for comma-separated array in string", func(t *testing.T) {
		values := []string{"Stan,Francine"}
		param := spec.QueryParam("names").Typed("array", "")
		param.Items = spec.NewItems().Typed("string", "")
		param.SimpleSchema.CollectionFormat = "csv"

		v, err := Parameter(values, param)
		assertConversionResult(t, []string{"Stan", "Francine"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for comma-separated array in string as a default behavior", func(t *testing.T) {
		values := []string{"Steve,Hayley"}
		param := spec.QueryParam("names").Typed("array", "")
		param.Items = spec.NewItems().Typed("string", "")

		v, err := Parameter(values, param)
		assertConversionResult(t, []string{"Steve", "Hayley"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("fail for array that has no items type", func(t *testing.T) {
		values := []string{"does not matter"}
		param := spec.QueryParam("names").Typed("array", "")

		v, err := Parameter(values, param)
		assertConversionResult(t, nil, v)
		assertConversionError(t, true, err)
	})

	t.Run("fail for file (not implemented yet)", func(t *testing.T) {
		values := []string{"does not matter"}
		param := spec.FormDataParam("photo").Typed("file", "")

		v, err := Parameter(values, param)
		assertConversionResult(t, nil, v)
		assertConversionError(t, true, err)
	})

	t.Run("fail for multiple values on primitive type", func(t *testing.T) {
		values := []string{"John", "Edvard"}
		param := spec.QueryParam("name").Typed("string", "")

		v, err := Parameter(values, param)
		assertConversionResult(t, nil, v)
		assertConversionError(t, true, err)
	})
}

func assertConversionResult(t *testing.T, expectedValue interface{}, v interface{}) {
	t.Helper()

	if !reflect.DeepEqual(expectedValue, v) {
		t.Errorf(
			"Expected value to be %v (%T) but got %v (%T)",
			expectedValue,
			expectedValue,
			v,
			v,
		)
	}
}

func assertConversionError(t *testing.T, expectError bool, err error) {
	t.Helper()

	if err != nil && !expectError {
		t.Errorf("Unexpected error: %v", err)
	}
	if err == nil && expectError {
		t.Error("Expected error, but got nil")
	}
}

func TestPrimitive(t *testing.T) {
	cases := []struct {
		value         string
		typ           string
		format        string
		expectedValue interface{}
		expectError   bool
	}{
		{
			value:         "Igor",
			typ:           "string",
			format:        "",
			expectedValue: "Igor",
		},
		{
			value:         "123",
			typ:           "integer",
			format:        "int32",
			expectedValue: int32(123),
		},
		{
			value:         "456",
			typ:           "integer",
			format:        "int64",
			expectedValue: int64(456),
		},
		{
			value:         "44.44",
			typ:           "number",
			format:        "float",
			expectedValue: float32(44.44),
		},
		{
			value:         "55.55",
			typ:           "number",
			format:        "double",
			expectedValue: float64(55.55),
		},
		{
			value:         "true",
			typ:           "boolean",
			expectedValue: true,
		},
		{
			value:         "1",
			typ:           "boolean",
			expectedValue: true,
		},
		{
			value:         "yes",
			typ:           "boolean",
			expectedValue: true,
		},
		{
			value:         "false",
			typ:           "boolean",
			expectedValue: false,
		},
		{
			// unknown string format
			value:       "some",
			typ:         "string",
			format:      "xml",
			expectError: true,
		},
		{
			// unknown number format
			value:       "$15.50",
			typ:         "number",
			format:      "currency",
			expectError: true,
		},
		{
			// unknown integer format
			value:       "i15",
			typ:         "integer",
			format:      "imaginary",
			expectError: true,
		},
		{
			// wrong value for number float
			value:       "44.44abc",
			typ:         "number",
			format:      "float",
			expectError: true,
		},
		{
			// wrong value for number double
			value:       "55.55abc",
			typ:         "number",
			format:      "double",
			expectError: true,
		},
		{
			// wrong value for integer int32
			value:       "123abc",
			typ:         "integer",
			format:      "int32",
			expectError: true,
		},
		{
			// wrong value for integer int64
			value:       "456abc",
			typ:         "integer",
			format:      "int64",
			expectError: true,
		},
		{
			// unknown type
			value:       "a",
			typ:         "char",
			expectError: true,
		},
	}

	for _, c := range cases {
		v, err := Primitive(c.value, c.typ, c.format)
		if err != nil && !c.expectError {
			t.Errorf("Unexpected error: %v", err)
		}
		if err == nil && c.expectError {
			t.Error("Expected error, but got nil")
		}

		if !reflect.DeepEqual(c.expectedValue, v) {
			t.Errorf(
				"Expected value to be %v (%T) but got %v (%T)",
				c.expectedValue,
				c.expectedValue,
				v,
				v,
			)
		}
	}
}
