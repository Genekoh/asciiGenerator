package image

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// GetImage opens a file and parses it, returning an image.Image
func GetImage(file *os.File) (image.Image, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		err = errors.New("unable to seek file")
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		err = errors.New("unable to decode image")
		return nil, err
	}

	return img, nil
}

func GetImageDimensions(img image.Image) (w, h uint) {
	size := img.Bounds().Size()
	width := uint(size.X)
	height := uint(size.Y)

	return width, height
}

// GetImageDimensionsFromFile opens a files and returns its width and height
func GetImageFileDimensions(file *os.File) (w, h uint, e error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		err = errors.New("unable to seek file")
		return 0, 0, err
	}

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		err = errors.New("unable to decode image config")
		return 0, 0, err
	}

	return uint(config.Width), uint(config.Height), nil
}

// ResizeImageAspect receives an image.Image and returns resized image with the newWidth and keeping its aspect ratio
func ResizeImageAspect(img image.Image, newWidth uint) (image.Image, uint) {
	w, h := GetImageDimensions(img)
	ratio := float64(h) / float64(w)
	newWidthFloat := float64(newWidth)
	newHeight := uint(ratio * newWidthFloat)

	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Bicubic)

	return resizedImage, newHeight
}

func ResizeImage(img image.Image, newWidth, newHeight uint) image.Image {
	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Bicubic)

	return resizedImage
}
