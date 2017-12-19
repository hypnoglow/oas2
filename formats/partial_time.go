package formats

import "time"

const (
	// RFC3339PartialTime is partial-time format as described in RFC3339
	// https://xml2rfc.tools.ietf.org/public/rfc/html/rfc3339.html#anchor14
	RFC3339PartialTime = "15:04:05"
)

// PartialTime represents a partial-time defined in RFC3339.
//
// swagger:strfmt partial-time
type PartialTime time.Time

// String implements fmt.Stringer
func (pt PartialTime) String() string {
	return time.Time(pt).Format(RFC3339PartialTime)
}

// MarshalText implements encoding.TextMarshaler
func (pt PartialTime) MarshalText() (text []byte, err error) {
	return []byte(pt.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (pt *PartialTime) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	partialTime, err := time.Parse(RFC3339PartialTime, string(text))
	if err != nil {
		return err
	}

	*pt = PartialTime(partialTime)
	return nil
}

// IsPartialTime returns true when the string is a valid partial time.
func IsPartialTime(value string) bool {
	_, err := time.Parse(RFC3339PartialTime, value)
	return err == nil
}
