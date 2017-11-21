package oas

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/spec"
)

// DecodeForm decodes form parameters by their spec to the dst.
//
// If you previously validated form parameters against the spec (using middleware
// or manually), then the returned error should be considered as a server error.
func DecodeForm(ps []spec.Parameter, req *http.Request, dst interface{}) error {
	isMultipart := true // assume multipart/form by default.
	if err := req.ParseMultipartForm(1024 * 1024 * 1024); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}

		// Form is actually not a multipart form.
		if err := req.ParseForm(); err != nil {
			return err
		}
		isMultipart = false
	}

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

		if p.Type == "file" {
			if !isMultipart {
				return fmt.Errorf("files can be only acquired through multipart/form-data")
			}

			files, ok := req.MultipartForm.File[p.Name]
			if !ok {
				// No such value in form.
				// This is not an error because validator should mark request
				// as invalid if this parameter is required.
				continue
			}

			if l := len(files); l != 1 {
				return fmt.Errorf("got %d files for param %s, want 1", l, p.Name)
			}

			file := files[0]
			if err := setFieldValue(file, f, dv); err != nil {
				return err
			}

			continue
		}

		var v interface{}
		if isMultipart {
			v, err = fromValues(req.MultipartForm.Value, p)
		} else {
			v, err = fromValues(req.Form, p)
		}

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
