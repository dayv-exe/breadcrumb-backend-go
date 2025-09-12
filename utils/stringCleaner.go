package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// remove diacritics (è → e, ç → c, etc.)
func removeDiacritics(s string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, _ := transform.String(t, s)
	return strings.TrimSpace(result)
}

// remove emojis and other symbols
func removeEmojis(s string) string {
	// Allow only letters, numbers, and whitespace
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	s = re.ReplaceAllString(s, "")

	reSpaces := regexp.MustCompile(`\s+`)
	s = reSpaces.ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}

func NormalizeString(s string) string {
	s = removeDiacritics(s)
	s = removeEmojis(s)

	return s
}

func SplitOnDelimiter(s string, delimiters ...string) []string {
	tokens := []string{s}
	for _, d := range delimiters {
		if strings.Contains(s, d) {
			tokens = append(tokens, strings.Split(s, d)...)
			break
		}
	}

	return tokens
}
