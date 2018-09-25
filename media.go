package oas

import "regexp"

var (
	contentTypeSelectorRegexJSON    *regexp.Regexp
	contentTypeSelectorRegexJSONAPI *regexp.Regexp
)

const (
	mediaTypeWildcard = "*/*"
)

func init() {
	contentTypeSelectorRegexJSON = regexp.MustCompile(`^application\/json`)
	contentTypeSelectorRegexJSONAPI = regexp.MustCompile(`^application\/vnd\.api\+json$`)
}
