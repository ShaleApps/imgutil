package img_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/ShaleApps/image-util"
)

func TestConversion(t *testing.T) {
	// colors to test
	baseImgs := []image.Image{
		makeImage(255, 0, 0), // red
		makeImage(0, 255, 0), // green
		makeImage(0, 0, 255), // blue
	}

	// colors that should be after conversion
	expectedImgs := []image.Image{
		makeImage(0, 0, 255), // to blue
		makeImage(255, 0, 0), // to red
		makeImage(0, 255, 0), // to green
	}

	red, err := img.HexToRBGA("FF0000")
	if err != nil {
		t.Fatal(err)
	}
	blue, err := img.HexToRBGA("0000FF")
	if err != nil {
		t.Fatal(err)
	}
	green, err := img.HexToRBGA("00FF00")
	if err != nil {
		t.Fatal(err)
	}

	convertedImgs := []image.Image{
		img.ConvertImageColor(baseImgs[0], img.ColorConverter(red, blue)),
		img.ConvertImageColor(baseImgs[1], img.ColorConverter(green, red)),
		img.ConvertImageColor(baseImgs[2], img.ColorConverter(blue, green)),
	}

	badConvertedImgs := []image.Image{
		img.ConvertImageColor(baseImgs[0], img.ColorConverter(blue, green)),
		img.ConvertImageColor(baseImgs[1], img.ColorConverter(red, blue)),
		img.ConvertImageColor(baseImgs[2], img.ColorConverter(green, red)),
	}

	for i, img := range expectedImgs {
		if !compare(img, convertedImgs[i]) {
			t.Fatalf("was supposed to pass: image %d did not match", i)
		}
		if compare(img, badConvertedImgs[i]) {
			t.Fatalf("was supposed to fail: image %d did match", i)
		}
	}
}

func compare(a, b image.Image) bool {
	aBounds := a.Bounds()
	bBounds := b.Bounds()
	width, height := aBounds.Max.X, bBounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			aR, aG, aB, aA := a.At(x, y).RGBA()
			bR, bG, bB, bA := b.At(x, y).RGBA()
			if aR != bR || aG != bG || aB != bB || aA != bA {
				return false
			}
		}
	}
	return true
}

func TestHexConversion(t *testing.T) {
	goodColors := []string{
		"FF0000",
		"#0000FF",
		"00FF00",
	}

	badColors := []string{
		"000FF00",
		"0000FF#",
		"F0000",
		"1234f5F",
	}

	for _, color := range goodColors {
		_, err := img.HexToRBGA(color)
		if err != nil {
			t.Fatalf("color %s failed when it was supposed to pass", color)
		}
	}

	for _, badColor := range badColors {
		_, err := img.HexToRBGA(badColor)
		if err == nil {
			t.Fatalf("color %s passed when it was supposed to fail", badColor)
		}
	}
}

func makeImage(r, g, b uint8) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1920, 1200))

	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}
