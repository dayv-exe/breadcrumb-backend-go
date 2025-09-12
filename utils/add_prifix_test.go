package utils

import "testing"

func TestAddPrefix(t *testing.T) {
	tests := []struct {
		word     string
		prefix   string
		expected string
	}{
		{
			word:     "test",
			prefix:   "user#",
			expected: "user#test",
		},
		{
			word:     "user#test",
			prefix:   "user#",
			expected: "user#test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := AddPrefix(tt.prefix, tt.word)
			if result != tt.expected {
				t.Fatalf("Expected %v but got %v instead", tt.expected, result)
			}
		})
	}
}
