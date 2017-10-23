package validate

import "github.com/go-openapi/strfmt"

var (
	formats strfmt.Registry
)

func init() {
	formats = strfmt.Default
	formats.Add("partialTime", &PartialTime{}, IsPartialTime)
}
