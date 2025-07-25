package utils

import (
	"regexp"
	"strings"
)

func IsNicknameValid(nickname string) bool {
	if len(nickname) < 3 || len(nickname) > 15 {
		return false
	}

	// rules: must start and end with letter or number, only one of dot or underscore, must be 3 to 15 characters long
	match, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9._]*[a-zA-Z0-9]$`, nickname)
	if !match {
		return false
	}

	if strings.Contains(nickname, "..") || strings.Contains(nickname, "__") {
		return false
	}

	if strings.Contains(nickname, "_.") || strings.Contains(nickname, "._") {
		return false
	}

	return true
}
