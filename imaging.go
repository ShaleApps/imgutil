package imgutil

import (
	"fmt"
	"image"
	"image/color"
	"regexp"
)

// ConvertImageColor takes an image, a color.Color to remove with another color.Color
// You can use HexToRBGA to get a color.Color via a hex code.
func ConvertImageColor(img image.Image, m color.Model) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	img2 := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newColor := m.Convert(img.At(x, y))
			img2.Set(x, y, newColor)
		}
	}
	return img2
}

// HexToRBGA returns a color.Color if it is provided with a hex code in the format
// of #?[0-9a-fA-F]{6}
func HexToRBGA(hex1 string) (color.Color, error) {
	validHex := regexp.MustCompile(`#?\b[0-9a-fA-F]{6}\b$`)
	if !validHex.MatchString(hex1) {
		return color.RGBA{}, fmt.Errorf("invalid hex: %s", hex1)
	}
	if hex1[0] == '#' {
		hex1 = hex1[1:]
	}
	return color.RGBA{
		R: hexDigitToInt(byte(hex1[1])) + hexDigitToInt(byte(hex1[0]))*16,
		G: hexDigitToInt(byte(hex1[3])) + hexDigitToInt(byte(hex1[2]))*16,
		B: hexDigitToInt(byte(hex1[5])) + hexDigitToInt(byte(hex1[4]))*16,
		A: 255,
	}, nil
}

// ColorConverter does stuff
func ColorConverter(c1 color.Color, c2 color.Color) color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		cR, cG, cB, cA := c.RGBA()
		c1R, c1G, c1B, c1A := c1.RGBA()
		if cR == c1R && cG == c1G && cB == c1B && cA == c1A {
			return c2
		}
		return c
	})
}

func hexDigitToInt(hex byte) byte {
	switch {
	case '0' <= hex && hex <= '9':
		return hex - '0'
	case 'a' <= hex && hex <= 'f':
		return hex - 'a' + 10
	case 'A' <= hex && hex <= 'F':
		return hex - 'A' + 10
	}
	return 0
}
