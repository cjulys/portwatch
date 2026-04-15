package reporter

import "fmt"

// String returns a human-readable name for the Format.
func (f Format) String() string {
	switch f {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	default:
		return fmt.Sprintf("unknown(%s)", string(f))
	}
}

// ParseFormat converts a string to a Format, returning an error for
// unrecognised values.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("reporter: unknown format %q (want \"text\" or \"json\")", s)
	}
}
