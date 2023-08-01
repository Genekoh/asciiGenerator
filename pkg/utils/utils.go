package utils

import (
	"image"

	"github.com/nfnt/resize"
)

func ResizeImage(img image.Image, newWidth, newHeight uint) image.Image {
	return resize.Resize(newWidth, newHeight, img, resize.Bicubic)
}

func ResizeImageAspect(img image.Image, newWidth uint) image.Image {
	w, h := uint(img.Bounds().Dx()), uint(img.Bounds().Dy())
	ratio := float64(h) / float64(w)
	newHeight := uint(ratio * float64(newWidth))

	return ResizeImage(img, newWidth, newHeight)
}
