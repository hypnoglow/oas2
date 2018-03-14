package oas

import "github.com/go-openapi/spec"

func copyOperation(dst *spec.Operation, src *spec.Operation) error {
	b, err := src.MarshalJSON()
	if err != nil {
		return err
	}

	return dst.UnmarshalJSON(b)
}
