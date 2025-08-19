package utils

import (
	"math/rand/v2"

	"github.com/lucasb-eyer/go-colorful"
)

type ColorPair struct {
	Foreground string
	Background string
}

var AllowedColors = []string{
	"#FF5733", // Bright Red-Orange
	"#33FF57", // Bright Green
	"#3357FF", // Bright Blue
	"#FF33A8", // Vibrant Pink
	"#FFD733", // Bright Yellow
}

func GenerateRandomColorPair() ColorPair {
	fg := AllowedColors[rand.IntN(len(AllowedColors))]

	color, err := colorful.Hex(fg)
	if err != nil {
		return ColorPair{"#ffffff", "#cccccc"} // fallback color
	}

	bg := color.BlendRgb(colorful.Color{R: 0, G: 0, B: 0}, 0.3).Hex()

	return ColorPair{fg, bg}
}
