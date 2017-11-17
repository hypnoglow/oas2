package validate

import (
	"net/url"
	"testing"

	"github.com/go-openapi/spec"
)

func TestPartialTimeValidation(t *testing.T) {
	param := spec.QueryParam("starts_at").Typed("string", "partial-time")

	t.Run("positive", func(t *testing.T) {
		values := url.Values{"starts_at": []string{"10:30:05"}}
		errs := Query([]spec.Parameter{*param}, values)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors %v", errs)
		}
	})

	t.Run("negative", func(t *testing.T) {
		values := url.Values{"starts_at": []string{"10abc"}}
		errs := Query([]spec.Parameter{*param}, values)
		if len(errs) == 0 {
			t.Errorf("Expected errors but got no errors")
		}
	})
}
