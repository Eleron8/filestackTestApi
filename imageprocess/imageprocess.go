package imageprocess

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

func OpenImage(jpg string) (image.Image, error) {
	file, err := os.Open(jpg)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	if format != "jpeg" {
		return nil, errors.New("image format is not jpeg")
	}
	return img, nil
}

func RotateByAngle(angle float64, pixels *[][]color.Color, radius int) {
	ppixels := *pixels
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	var newImage [][]color.Color

	newImage = make([][]color.Color, 2*radius+1)
	for i := 0; i < len(newImage); i++ {
		newImage[i] = make([]color.Color, 2*radius+1)
	}

	fmt.Println("x:", len(newImage), "y:", len(newImage[0]), "radius:", radius)

	for i := 0; i < len(ppixels); i++ {
		for j := 0; j < len(ppixels[i]); j++ {
			xNew := int(float64(i)*cos-float64(j)*sin) + radius
			yNew := int(float64(i)*sin+float64(j)*cos) + radius

			newImage[xNew][yNew] = ppixels[i][j]
		}
	}
	*pixels = newImage
}

func CropImage(img image.Image, rect image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	crop, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}
	return crop.SubImage(rect), nil
}

func createImage(newImage *image.RGBA, name string) error {
	fg, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fg.Close()
	err = jpeg.Encode(fg, newImage, nil)
	if err != nil {
		return err
	}
	return nil
}
