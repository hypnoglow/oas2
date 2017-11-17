package formats

import (
	"bytes"
	"testing"
	"time"
)

func TestPartialTime_String(t *testing.T) {
	date, err := time.Parse(time.RFC3339, "2017-01-01T14:25:00Z")
	if err != nil {
		t.Fatalf("Unexpeted error: %v", err)
	}

	pt := PartialTime(date)
	expectedValue := "14:25:00"
	if pt.String() != expectedValue {
		t.Fatalf("Expected %q but got %q", expectedValue, pt.String())
	}
}

func TestPartialTime_MarshalText(t *testing.T) {
	date, err := time.Parse(time.RFC3339, "2017-01-01T14:25:00Z")
	if err != nil {
		t.Fatalf("Unexpeted error: %v", err)
	}

	pt := PartialTime(date)
	expectedValue := []byte("14:25:00")
	txt, err := pt.MarshalText()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(txt, expectedValue) {
		t.Fatalf("Expected %v but got %v", expectedValue, txt)
	}
}

func TestPartialTime_UnmarshalText(t *testing.T) {
	t.Run("should unmarshal text for correct value", func(t *testing.T) {
		pt := PartialTime{}

		err := pt.UnmarshalText([]byte("14:25:00"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedValue := "14:25:00"
		if pt.String() != expectedValue {
			t.Fatalf("Expected %v but got %v", expectedValue, pt.String())
		}
	})

	t.Run("should be ok for zero length text", func(t *testing.T) {
		pt := PartialTime{}

		err := pt.UnmarshalText([]byte{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedValue := "00:00:00"
		if pt.String() != expectedValue {
			t.Fatalf("Expected %v but got %v", expectedValue, pt.String())
		}
	})

	t.Run("should fail on incorrect value", func(t *testing.T) {
		pt := PartialTime{}

		err := pt.UnmarshalText([]byte("123456789"))
		if err == nil {
			t.Errorf("Expected error but got <nil>")
		}
	})
}

func TestIsPartialTime(t *testing.T) {
	t.Run("should be partial time", func(t *testing.T) {
		isPartial := IsPartialTime("14:15:30")
		if isPartial == false {
			t.Errorf("Expected true but got %v", isPartial)
		}
	})

	t.Run("should not be partial time", func(t *testing.T) {
		isPartial := IsPartialTime("1234567890")
		if isPartial == true {
			t.Errorf("Expected false but got %v", isPartial)
		}
	})
}
