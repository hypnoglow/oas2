package convert

import (
	"fmt"
	"testing"
)

func ExampleArray() {
	values := []string{"123", "456"}
	v, _ := Array(values, "integer", "int64")

	fmt.Printf("%#v", v)

	// Output:
	// []int64{123, 456}
}

func TestArray(t *testing.T) {
	t.Run("ok for string array", func(t *testing.T) {
		values := []string{"Nicolas", "Jonathan"}
		v, err := Array(values, "string", "")

		assertConversionResult(t, []string{"Nicolas", "Jonathan"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for string array with any format", func(t *testing.T) {
		values := []string{"Nicolas", "Jonathan"}
		v, err := Array(values, "string", "phone")

		assertConversionResult(t, []string{"Nicolas", "Jonathan"}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for int64 array", func(t *testing.T) {
		values := []string{"123", "456"}

		v, err := Array(values, "integer", "int64")
		assertConversionResult(t, []int64{123, 456}, v)
		assertConversionError(t, false, err)
	})

	t.Run("ok for int32 array", func(t *testing.T) {
		values := []string{"123", "456"}

		v, err := Array(values, "integer", "int32")
		assertConversionResult(t, []int32{123, 456}, v)
		assertConversionError(t, false, err)
	})

	t.Run("fail on mixed types in int64 array", func(t *testing.T) {
		values := []string{"123", "Max"}

		v, err := Array(values, "integer", "int64")
		assertConversionResult(t, nil, v)
		assertConversionError(t, true, err)
	})

	t.Run("fail on mixed types in int32 array", func(t *testing.T) {
		values := []string{"123", "Max"}

		v, err := Array(values, "integer", "int32")
		assertConversionResult(t, nil, v)
		assertConversionError(t, true, err)
	})

	t.Run("ok for int array with any other format", func(t *testing.T) {
		values := []string{"123", "456"}

		v, err := Array(values, "integer", "years")
		assertConversionResult(t, []int64{123, 456}, v)
		assertConversionError(t, false, err)
	})

	t.Run("fail on unsupported type", func(t *testing.T) {
		values := []string{"true", "false"}

		v, err := Array(values, "boolean", "")
		assertConversionResult(t, nil, v)
		assertConversionError(t, true, err)
	})
}
