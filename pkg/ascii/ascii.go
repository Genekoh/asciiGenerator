package ascii

import (
	"bytes"
	imagePkg "github.com/Genekoh/asciiGenerator/pkg/image"
	"image"
	"image/color"
	"reflect"
)

var (
	Characters = "MN8@O$Zbe*+!:.,  "
	Table      = []byte(Characters)
)

// GenerateAscii takes in an image
func GenerateAscii(img image.Image) string {
	uintW, uintH := imagePkg.GetImageDimensions(img)
	w := uintW
	h := uintH
	ratio := 10.0 / 22.0 // Characters in terminal aren't equally tall as it is wide
	h = uint(ratio * float64(h))

	i := imagePkg.ResizeImage(img, w, h)

	buffer := WriteAsciiBytes(i, int(w), int(h))
	return buffer.String()
}

// WriteAsciiBytes returns a *bytes.Buffer of asciiCharacters that represents the brightness of each pixel of the image.Image given
func WriteAsciiBytes(img image.Image, w, h int) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := color.GrayModel.Convert(img.At(x, y))
			brightness := reflect.ValueOf(p).FieldByName("Y").Uint()
			length := uint64(len(Characters) - 1)
			index := int(brightness * length / 255)
			buffer.WriteByte(Table[index])
		}
		buffer.WriteByte('\n')
	}

	return buffer
}
