package utils

import "testing"

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
