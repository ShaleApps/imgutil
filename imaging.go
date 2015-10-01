// Package imgutil provides an easy way to take an image and convert any color in
// it into another color. This also allows the conversion of a color image into
// a gray scale image. Supports the image formats image supports such as png and
// jpeg. This exclusively works with RGBA's from image.
package imgutil

import (
	"fmt"
	"image"
	"image/color"
	"regexp"

	"github.com/nfnt/resize"
)

// ConvertImageColor takes an image and a color model and returns a copy
// of the image's color with provided replacement color.
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

// HexToRBGA returns a color.Color if it is provided with a hex code in the correct
// format. https://en.wikipedia.org/wiki/Web_colors
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

// ColorConverter compares a color to another. IF it matches, return
// the new color.
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

// ResizeImage takes the dimensions of an image and then the image and
// returns a new image resized to that dimension. If either width or height
// is set as 0, the aspect ratio will be preserved and the counter part
// will be used as a determining factor for the size.
// Uses github.com/nfnt/resize
func ResizeImage(maxWidth, maxHeight int, img image.Image) image.Image {
	return resize.Resize(uint(maxWidth), uint(maxHeight), img, resize.NearestNeighbor)
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
