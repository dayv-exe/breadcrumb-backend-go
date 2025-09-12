package utils

import (
	"breadcrumb-backend-go/constants"
	"regexp"
	"strings"
)

// not in use yet
var bannedNicknameSubstrings = []string{"user", ""}

func NicknameValid(nickname string) bool {
	if len(nickname) < constants.MIN_USERNAME_CHARS || len(nickname) > constants.MAX_USERNAME_CHARS {
		return false
	}

	// rules: must start and end with letter or number, only one of dot or underscore, must be 3 to 15 characters long
	match, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9._]*[a-zA-Z0-9]$`, nickname)
	if !match {
		return false
	}

	// to allow striping the username and placing it in the search index table for easier searching

	// block patterns like a_1, b.2, a_b, c.d
	if len(strings.Replace(nickname, ".", "", 1)) < 3 || len(strings.Replace(nickname, "_", "", 1)) < 3 {
		return false
	}

	// block if username contains no letters
	if ok, _ := regexp.MatchString(`[a-zA-Z]`, nickname); !ok {
		return false
	}

	// can only contain on dot or underscore
	if strings.Contains(nickname, "..") || strings.Contains(nickname, "__") {
		return false
	}

	if strings.Contains(nickname, "_.") || strings.Contains(nickname, "._") {
		return false
	}

	return true
}
