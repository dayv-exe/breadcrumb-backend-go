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

func TestNameChangeAllowed(t *testing.T) {
	tests := []struct {
		name              string
		lastChangeDate    string
		expectAllowChange bool
	}{
		{
			// yyyy/mm/dd
			name:              "last changed a week ago",
			lastChangeDate:    time.Now().AddDate(0, 0, -7).Format(constants.FULL_DATE_TIME_LAYOUT),
			expectAllowChange: false,
		},
		{
			// yyyy/mm/dd
			name:              "last changed 2 years ago",
			lastChangeDate:    time.Now().AddDate(-2, 0, 0).Format(constants.FULL_DATE_TIME_LAYOUT),
			expectAllowChange: true,
		},
		{
			// yyyy/mm/dd
			name:              "last changed 2 days ago",
			lastChangeDate:    time.Now().AddDate(0, 0, -2).Format(constants.FULL_DATE_TIME_LAYOUT),
			expectAllowChange: false,
		},
		{
			// yyyy/mm/dd
			name:              "last changed 31 days ago",
			lastChangeDate:    time.Now().AddDate(0, 0, -31).Format(constants.FULL_DATE_TIME_LAYOUT),
			expectAllowChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NameChangeAllowed(tt.lastChangeDate)
			if err != nil {
				t.Errorf("An unexpected error occurred %v", err)
			}
			if result != tt.expectAllowChange {
				t.Errorf("Expected %v but got %v instead!", tt.expectAllowChange, result)
			}
		})
	}
}
