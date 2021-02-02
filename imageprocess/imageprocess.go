package imageprocess

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/Eleron8/filestackTestApi/models"
	"go.uber.org/zap"
)

type ImageProcess struct {
	ImageFile  image.Image
	Transforms []models.Transform
	logger     *zap.Logger
}

// type Action struct {
// 	Degrees float64
// 	Width   int
// 	Height  int
// }

func OpenImage(filename string, transforms []models.Transform, logger *zap.Logger) (ImageProcess, string, error) {
	handleErr := func(err error) (ImageProcess, string, error) {
		return ImageProcess{}, "", nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return handleErr(err)
	}
	defer file.Close()
	img, format, err := image.Decode(file)
	if err != nil {
		return handleErr(err)
	}
	if format != "jpeg" && format != "png" {
		return handleErr(errors.New("image format is not correct"))
	}
	imgProcess := ImageProcess{
		ImageFile:  img,
		Transforms: transforms,
		logger:     logger,
	}
	return imgProcess, format, nil
}

func (imgPr ImageProcess) ImageTransform(param models.Transform, format string, name string, radius int, formatPart string) (string, error) {
	switch param.Type {
	case models.Rotate:
		size := imgPr.ImageFile.Bounds().Size()
		newimage := imgPr.rotateByAngle(param.Params.Degrees, size, radius)
		filename := fmt.Sprintf("%s_%s_%f", name, models.Rotate, param.Params.Degrees)
		filename = filename + formatPart
		imgPr.logger.Info("in rotate case ", zap.String("filename", filename))
		if err := imgPr.createImageRGBA(newimage, filename, format); err != nil {
			return "", err
		}
		return filename, nil
	case models.Crop:
		rect := image.Rect(0, 0, param.Params.Width, param.Params.Height)
		newImg, err := imgPr.cropImage(imgPr.ImageFile, rect)
		if err != nil {
			return "", err
		}
		filename := fmt.Sprintf("%s_%s_%d_%d", name, models.Crop, param.Params.Width, param.Params.Height)
		filename = filename + formatPart
		imgPr.logger.Info("in crop case ", zap.String("filename", filename))
		if err := imgPr.createImage(newImg, filename, format); err != nil {
			return "", err
		}
		return filename, nil
	case models.RemoveExif:
		filename := fmt.Sprintf("%s_%s_%d", name, models.RemoveExif, rand.Intn(100))
		filename = filename + formatPart
		imgPr.logger.Info("in remove exif case ", zap.String("filename", filename))
		if err := imgPr.createImage(imgPr.ImageFile, filename, format); err != nil {
			return "", err
		}
		return filename, nil
	}
	return "", nil
}

func (imgPr ImageProcess) ImageTransforms(format string, name string) error {
	_, radius := imgPr.GetPixels()
	for _, v := range imgPr.Transforms {
		switch v.Type {
		case models.Rotate:
			size := imgPr.ImageFile.Bounds().Size()
			newimage := imgPr.rotateByAngle(v.Params.Degrees, size, radius)
			filename := fmt.Sprintf("%s_%s_%f", name, models.Rotate, v.Params.Degrees)
			if err := imgPr.createImageRGBA(newimage, filename, format); err != nil {
				return err
			}
		case models.Crop:
			rect := image.Rect(0, 0, v.Params.Width, v.Params.Height)
			newImg, err := imgPr.cropImage(imgPr.ImageFile, rect)
			if err != nil {
				return err
			}
			filename := fmt.Sprintf("%s_%s_%d_%d", name, models.Crop, v.Params.Width, v.Params.Height)
			if err := imgPr.createImage(newImg, filename, format); err != nil {
				return err
			}
		}

	}
	return nil
}

func (imgPr ImageProcess) rotateByAngle(angle float64, size image.Point, radius int) *image.RGBA {
	var pixels [][]color.Color
	for i := 0; i < size.X; i++ {
		var y []color.Color
		for j := 0; j < size.Y; j++ {
			y = append(y, imgPr.ImageFile.At(i, j))
		}
		pixels = append(pixels, y)
	}

	cos := math.Cos(angle)
	sin := math.Sin(angle)

	var newImage [][]color.Color

	newImage = make([][]color.Color, 2*radius+1)
	for i := 0; i < len(newImage); i++ {
		newImage[i] = make([]color.Color, 2*radius+1)
	}

	fmt.Println("x:", len(newImage), "y:", len(newImage[0]), "radius:", radius)

	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			xNew := int(float64(i)*cos-float64(j)*sin) + radius
			yNew := int(float64(i)*sin+float64(j)*cos) + radius

			newImage[xNew][yNew] = pixels[i][j]
		}
	}
	rect := image.Rect(0, 0, len(newImage), len(newImage[0]))
	nImg := image.NewRGBA(rect)

	for x := 0; x < len(newImage); x++ {
		for y := 0; y < len(newImage[0]); y++ {
			q := newImage[x]
			if q == nil {
				continue
			}
			p := newImage[x][y]
			if p == nil {
				continue
			}
			original, ok := color.RGBAModel.Convert(p).(color.RGBA)
			if ok {
				nImg.Set(x, y, original)
			}
		}
	}
	return nImg
}

func (i ImageProcess) cropImage(img image.Image, rect image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	crop, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}
	return crop.SubImage(rect), nil
}

func (imgPr ImageProcess) createImageRGBA(newImage *image.RGBA, name string, format string) error {
	fg, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fg.Close()
	if format == "jpeg" {
		err := jpeg.Encode(fg, newImage, nil)
		if err != nil {
			return err
		}
	}
	if format == "png" {
		err := png.Encode(fg, newImage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (imgPr ImageProcess) createImage(newImage image.Image, name string, format string) error {
	fg, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fg.Close()
	if format == "jpeg" {
		err := jpeg.Encode(fg, newImage, nil)
		if err != nil {
			return err
		}
	}
	if format == "png" {
		err := png.Encode(fg, newImage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (imgPr ImageProcess) GetPixels() ([][]color.Color, int) {
	size := imgPr.ImageFile.Bounds().Size()
	var pixels [][]color.Color
	for i := 0; i < size.X; i++ {
		var y []color.Color
		for j := 0; j < size.Y; j++ {

			y = append(y, imgPr.ImageFile.At(i, j))
		}
		pixels = append(pixels, y)
	}
	radius := int(math.Sqrt(float64(size.X*size.X) + float64(size.Y*size.Y)))
	return pixels, radius
}
