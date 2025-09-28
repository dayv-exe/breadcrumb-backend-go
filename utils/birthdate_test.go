package utils

import (
	"breadcrumb-backend-go/constants"
	"testing"
	"time"
)

func TestBirthdateValid(t *testing.T) {
	// DD/MM/YYYY
	tests := []struct {
		name      string
		birthdate string
		wantValid bool
		wantError bool
	}{
		{
			name:      "(min age - 1 month) which should be 12 years and 11 months old",
			birthdate: time.Now().AddDate(-constants.MIN_AGE, 1, 0).Format(constants.BIRTHDATE_ONLY_DATE_LAYOUT),
			wantValid: false,
			wantError: false,
		},
		{
			name:      "today years old",
			birthdate: time.Now().Format(constants.BIRTHDATE_ONLY_DATE_LAYOUT),
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
			name:      "(min age) 13 years old",
			birthdate: time.Now().AddDate(-constants.MIN_AGE, 0, 0).Format(constants.BIRTHDATE_ONLY_DATE_LAYOUT),
			wantValid: true,
			wantError: false,
		},
		{
			name:      "too damn old",
			birthdate: time.Now().AddDate(-constants.MAX_AGE-1, 0, 0).Format(constants.BIRTHDATE_ONLY_DATE_LAYOUT),
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
