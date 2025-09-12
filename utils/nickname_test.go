package utils

import (
	"testing"
)

func TestIsNicknameValid(t *testing.T) {
	tests := []struct {
		nickname string
		valid    bool
	}{
		{"john_doe", true},
		{"_john", false},
		{"john..doe", false},
		{"j", false},
		{"averylongnamethatshouldfail", false},
		{"john.doe_", false},
		{"john__doe", false},
		{"j.doe", true},
		{"j_doe", true},
		{"JohnDoe99", true},
		{"john_.doe", false},
		{"john.doe_", false},
		{".johndoe", false},
		{"4ksf_sqmd1", true},
		{"john.doe001234ed", false},
		{"14792384913", false},
		{"a.1", false},
		{"ab.1", true},
		{"p_12345", true},
	}

	for _, tt := range tests {
		t.Run(tt.nickname, func(t *testing.T) {
			got := NicknameValid(tt.nickname)
			if got != tt.valid {
				t.Errorf("isNicknameValid(%q) = %v; want %v", tt.nickname, got, tt.valid)
			}
		})
	}
}
