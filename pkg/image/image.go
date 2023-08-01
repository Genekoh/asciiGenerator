package imagePkg

import (
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp"
)

// GetImage opens a file and parses it, returning an image.Image
func GetImage(file *os.File) (image.Image, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func GetGif(file *os.File) (*gif.GIF, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	gifPtr, err := gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}

	return gifPtr, nil
}

func GetImageDimensions(img image.Image) (w, h uint) {
	size := img.Bounds().Size()
	w = uint(size.X)
	h = uint(size.Y)

	return w, h
}

func GetImageDimensionsInt(img image.Image) (w, h int) {
	uw, uh := GetImageDimensions(img)
	w, h = int(uw), int(uh)

	return w, h
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
