package validate

import "testing"

func TestValidationError(t *testing.T) {
	ve := ValidationErrorf("name", nil, "name cannot be empty")

	if ve.Error() != "name cannot be empty" {
		t.Errorf("Unexpected error message")
	}

	if ve.Field() != "name" {
		t.Errorf("Unexpected error field")
	}

	if ve.Value() != nil {
		t.Errorf("Unexpected error value")
	}
}
