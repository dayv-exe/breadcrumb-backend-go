package utils

import (
	"breadcrumb-backend-go/constants"
	"testing"
	"time"
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

func TestBirthdateValid(t *testing.T) {
	// DD/MM/YYYY
	tests := []struct {
		name      string
		birthdate string
		wantValid bool
		wantError bool
	}{
		{
			name:      "*min age - 1 month* 12 years and 11 months old",
			birthdate: time.Now().AddDate(-constants.MIN_AGE, 1, 0).Format(constants.DATE_LAYOUT),
			wantValid: false,
			wantError: false,
		},
		{
			name:      "today years old",
			birthdate: time.Now().Format(constants.DATE_LAYOUT),
			wantValid: false,
			wantError: false,
		},
		{
			name:      "invalid format",
			birthdate: time.Now().Format("2006/01/02"),
			wantValid: false,
			wantError: true,
		},
		{
			name:      "*min age* 18 year old",
			birthdate: time.Now().AddDate(-constants.MIN_AGE, 0, 0).Format(constants.DATE_LAYOUT),
			wantValid: true,
			wantError: false,
		},
		{
			name:      "too damn old",
			birthdate: time.Now().AddDate(-100, 0, 0).Format(constants.DATE_LAYOUT),
			wantValid: false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BirthdateIsValid(tt.birthdate)
			if (err != nil) != tt.wantError {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.wantValid {
				t.Errorf("got %v, want %v", got, tt.wantValid)
			}
		})
	}
}
