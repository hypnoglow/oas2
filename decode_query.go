package oas

import (
	"net/url"

	"github.com/go-openapi/spec"
)

// DecodeQuery decodes query parameters by their spec to the dst.
//
// If you previously validated query parameters against the spec (using middleware
// or manually), then the returned error should be considered as a server error.
func DecodeQuery(ps []spec.Parameter, q url.Values, dst interface{}) error {
	dv, err := reflectValueOfDestination(dst)
	if err != nil {
		return err // just propagate
	}

	fields := fieldMap(dv)

	for _, p := range ps {
		// No such tag in struct - no need to populate.
		f, ok := fields[p.Name]
		if !ok {
			continue
		}

		v, err := fromValues(q, p)
		if err != nil {
			return err // just propagate
		}
		if v == nil {
			continue
		}

		if err := setFieldValue(v, f, dv); err != nil {
			return err
		}
	}

	return nil
}
