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
	contentTypeSelectorRegexJSON = regexp.MustCompile(`(?i)^application\/json`)
	contentTypeSelectorRegexJSONAPI = regexp.MustCompile(`(?i)^application\/vnd\.api\+json$`)
}
