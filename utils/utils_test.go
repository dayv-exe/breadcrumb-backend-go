package utils

import "testing"

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
		{"8david_arubs9", true},
		{"john.doe001234ed", false},
		{"14792384913", true},
	}

	for _, tt := range tests {
		t.Run(tt.nickname, func(t *testing.T) {
			got := IsNicknameValid(tt.nickname)
			if got != tt.valid {
				t.Errorf("isNicknameValid(%q) = %v; want %v", tt.nickname, got, tt.valid)
			}
		})
	}
}

func TestEmailValid(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@gmail.com", true},
		{"testgmail.com", false},
		{".testgmail.com", false},
		{"testgmail", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			got := IsEmailValid(tt.email)
			if got != tt.valid {
				t.Errorf("isEmailValid(%q) = %v; want %v", tt.email, got, tt.valid)
			}
		})
	}
}
