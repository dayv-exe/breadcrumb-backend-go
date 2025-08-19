package utils

import "testing"

func TestRandomColorGen(t *testing.T) {
	randFgColor := GenerateRandomColorPair().Foreground
	tests := []struct {
		name       string
		foreground string
		wantValid  bool
	}{
		{
			name:       "test random valid color",
			foreground: randFgColor,
			wantValid:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValidColor := contains(AllowedColors, tt.foreground)
			if tt.wantValid && !isValidColor {
				t.Errorf("want valid color, got invalid color")
			} else if !tt.wantValid && isValidColor {
				t.Errorf("want invalid color, got valid color")
			}
		})
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
