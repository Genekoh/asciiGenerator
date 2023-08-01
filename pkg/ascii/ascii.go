package ascii

import (
	"bytes"
	"image"
	"image/color"
	"reflect"

	imagePkg "github.com/Genekoh/asciiGenerator/pkg/image"
)

const (
	// CharacterSet = "MN8@O$Zbe*+!:.,  "
	CharacterSet            = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	CorrectionRatio float64 = 10.0 / 22.0
)

var Table = []byte(CharacterSet)

type StoredAscii struct {
	Frames []string
	Delay  []int
	Height uint
}

func NewStoredAscii(frames []string, delay []int, height uint) *StoredAscii {
	return &StoredAscii{frames, delay, height}
}

// GenerateAscii takes in an image, resizes it to maintain image aspect ratio and returns a *bytes.Buffer of asciiCharacters
func GenerateAscii(img image.Image) *bytes.Buffer {
	w, h := imagePkg.GetImageDimensions(img)
	h = uint(CorrectionRatio * float64(h)) // CharacterSet in terminal aren't equally tall as it is wide

	resizedImg := imagePkg.ResizeImage(img, w, h)

	buffer := writeAsciiBytes(resizedImg)
	return buffer
}

// writeAsciiBytes returns a *bytes.Buffer of asciiCharacters that represents the brightness of each pixel of the image.Image given
func writeAsciiBytes(img image.Image) *bytes.Buffer {
	buffer := new(bytes.Buffer)

	w, h := imagePkg.GetImageDimensionsInt(img)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := color.GrayModel.Convert(img.At(x, y))
			brightness := reflect.ValueOf(p).FieldByName("Y").Uint()
			length := uint64(len(CharacterSet) - 1)
			index := int(brightness * length / 255)
			buffer.WriteByte(Table[index])
		}
		buffer.WriteByte('\n')
	}

	return buffer
}
