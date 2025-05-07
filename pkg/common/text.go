package common

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveAccents removes accents from a string
func RemoveAccents(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isNonSpacingMark), norm.NFC)
	result, _, _ := transform.String(t, input)
	return result
}

// IsNonSpacingMark checks if a rune is a non-spacing mark
func isNonSpacingMark(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// ExtractMiddleName extracts the middle name from a location string
func ExtractMiddleName(localidade string) string {
	parts := strings.Split(localidade, " - ")
	if len(parts) > 1 {
		return parts[1]
	}
	return localidade
}
