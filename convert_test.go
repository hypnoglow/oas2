package oas2

import (
	"reflect"
	"testing"
)

func TestConvertParameter(t *testing.T) {
	cases := []struct {
		values        []string
		typ           string
		format        string
		expectedValue interface{}
		expectError   bool
	}{
		{
			values:        []string{"John"},
			typ:           "string",
			expectedValue: "John",
			expectError:   false,
		},
		{
			values:        []string{"does not matter"},
			typ:           "array",
			expectedValue: nil,
			expectError:   true,
		},
		{
			values:        []string{"does not matter"},
			typ:           "file",
			expectedValue: nil,
			expectError:   true,
		},
		{
			values:        []string{"John", "Edvard"},
			typ:           "string",
			expectedValue: nil,
			expectError:   true,
		},
	}

	for _, c := range cases {
		v, err := ConvertParameter(c.values, c.typ, c.format)

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

func TestConvertPrimitive(t *testing.T) {
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
		v, err := ConvertPrimitive(c.value, c.typ, c.format)
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
