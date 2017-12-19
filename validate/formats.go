package validate

import (
	"github.com/go-openapi/strfmt"

	"github.com/hypnoglow/oas2/formats"
)

var (
	formatRegistry strfmt.Registry
)

func init() {
	formatRegistry = strfmt.Default

	RegisterFormat("partialtime", &formats.PartialTime{}, formats.IsPartialTime)
}

// RegisterFormat registers custom format and validator for it.
// See default OAS formatRegistry here:
// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types
func RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	formatRegistry.Add(name, format, validator)
}
