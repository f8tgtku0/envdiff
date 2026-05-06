package exporter

import "fmt"

// ParseFormat converts a user-supplied string into a Format constant.
// It returns an error when the string does not match a known format.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatShell:
		return FormatShell, nil
	case FormatDotenv:
		return FormatDotenv, nil
	case FormatExport:
		return FormatExport, nil
	default:
		return "", fmt.Errorf("exporter: unknown format %q (valid: dotenv, shell, export)", s)
	}
}

// ValidFormats returns all supported format names as a slice of strings.
func ValidFormats() []string {
	return []string{
		string(FormatDotenv),
		string(FormatShell),
		string(FormatExport),
	}
}
