package imgutil_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/ShaleApps/imgutil"
)

func TestConversion(t *testing.T) {
	// colors to test
	baseImgs := []image.Image{
		makeImage(255, 0, 0), // red
		makeImage(0, 255, 0), // green
		makeImage(0, 0, 255), // blue
		// test to see if greyscale works
		makeImage(0, 0, 255), // blue
	}

	// colors that should be after conversion
	expectedImgs := []image.Image{
		makeImage(0, 0, 255), // to blue
		makeImage(255, 0, 0), // to red
		makeImage(0, 255, 0), // to green
		// test to see if greyscale works
		makeGray(0, 0, 255), // gray
	}

	red, err := imgutil.HexToRBGA("FF0000")
	if err != nil {
		t.Fatal(err)
	}
	blue, err := imgutil.HexToRBGA("0000Ff")
	if err != nil {
		t.Fatal(err)
	}
	green, err := imgutil.HexToRBGA("00FF00")
	if err != nil {
		t.Fatal(err)
	}

	goodConvertedImgs := []image.Image{
		imgutil.ConvertImageColor(baseImgs[0], imgutil.ColorConverter(red, blue)),
		imgutil.ConvertImageColor(baseImgs[1], imgutil.ColorConverter(green, red)),
		imgutil.ConvertImageColor(baseImgs[2], imgutil.ColorConverter(blue, green)),
		imgutil.ConvertImageColor(baseImgs[2], image.NewGray(image.Rect(0, 0, 1920, 1200)).ColorModel()),
	}

	badConvertedImgs := []image.Image{
		imgutil.ConvertImageColor(baseImgs[0], imgutil.ColorConverter(blue, green)),
		imgutil.ConvertImageColor(baseImgs[1], imgutil.ColorConverter(red, blue)),
		imgutil.ConvertImageColor(baseImgs[2], imgutil.ColorConverter(green, red)),
		imgutil.ConvertImageColor(baseImgs[2], imgutil.ColorConverter(green, red)),
	}

	for i, img := range expectedImgs {
		if !compare(img, goodConvertedImgs[i]) {
			t.Fatalf("was supposed to pass: image %d did not match", i)
		}
		if compare(img, badConvertedImgs[i]) {
			t.Fatalf("was supposed to fail: image %d did match", i)
		}
	}
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
		_, err := imgutil.HexToRBGA(color)
		if err != nil {
			t.Fatalf("color %s failed when it was supposed to pass", color)
		}
	}

	for _, badColor := range badColors {
		_, err := imgutil.HexToRBGA(badColor)
		if err == nil {
			t.Fatalf("color %s passed when it was supposed to fail", badColor)
		}
	}
}

func TestResize(t *testing.T) {
	// We create an image but we don't care about the color, only it's size.
	// The ratio of this base image is 1.6
	baseImage := makeImage(255, 255, 255)

	assortedRects := []image.Rectangle{
		image.Rect(0, 0, 2560, 1600),
		image.Rect(0, 0, 1680, 1050),
		image.Rect(0, 0, 123, 908),
	}

	expectedImages := []image.Image{
		image.NewRGBA(assortedRects[0]),
		image.NewRGBA(assortedRects[1]),
		image.NewRGBA(assortedRects[2]),
	}

	badImages := []image.Image{
		image.NewRGBA(assortedRects[1]),
		image.NewRGBA(assortedRects[2]),
		image.NewRGBA(assortedRects[0]),
	}

	resultImages := []image.Image{
		imgutil.ResizeImage(2560, 0, baseImage),
		imgutil.ResizeImage(0, 1050, baseImage),
		imgutil.ResizeImage(123, 908, baseImage),
	}

	for i, resultImg := range resultImages {
		if !sameSize(resultImg, expectedImages[i]) {
			t.Fatalf("resultImages[%d] does not match size for expectedImages[%d]", i, i)
		}
	}

	for i, resultImg := range resultImages {
		if sameSize(resultImg, badImages[i]) {
			t.Fatalf("resultImages[%d] does match size for badImages[%d]", i, i)
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

func makeImage(r, g, b uint8) image.Image {
	rect := image.Rect(0, 0, 1920, 1200)
	img := image.NewRGBA(rect)

	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

func sameSize(imgA, imgB image.Image) bool {
	aBounds := imgA.Bounds()
	bBounds := imgB.Bounds()
	aHeight, aWidth := aBounds.Max.X, aBounds.Max.Y
	bHeight, bWidth := bBounds.Max.X, bBounds.Max.Y

	if aHeight != bHeight || aWidth != bWidth {
		return false
	}
	return true
}

func makeGray(r, g, b uint8) image.Image {
	rect := image.Rect(0, 0, 1920, 1200)
	img := image.NewGray(rect)

	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}
