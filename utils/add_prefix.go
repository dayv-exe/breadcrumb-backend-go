package utils

import "strings"

func AddPrefix(prefix string, s string) string {
	parts := strings.Split(s, prefix)
	return prefix + parts[len(parts)-1]
}

func AddSuffix(s string, suffix string) string {
	parts := strings.Split(s, suffix)
	return parts[len(parts)-1] + suffix
}
