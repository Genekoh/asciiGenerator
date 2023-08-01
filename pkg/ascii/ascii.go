package ascii

import (
	"image"
	"image/color"
	"reflect"
	"strings"
	"sync"

	"github.com/Genekoh/asciiGenerator/pkg/utils"
)

const (
	// DefaultCharSet = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
	DefaultCharSet          = "MN8@O$Zbe*+!:.,  "
	CorrectionRatio float64 = 10.0 / 22.0
)

type Frame struct {
	AsciiString string
	Delay       int
	Width       uint
	Height      uint
}

type StoredContent struct {
	Frames    []Frame
	LoopCount int
}

func NewFrame(asciiString string, delay int, width, height uint) Frame {
	return Frame{asciiString, delay, width, height}
}

func NewStoredContent(frames []Frame, loopCount int) *StoredContent {
	return &StoredContent{frames, loopCount}
}

func convertToAscii(img image.Image, charSet string, inverted bool) *[]string {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	xs := make([]string, h)

	var wg sync.WaitGroup
	for y := range xs {
		wg.Add(1)
		// convert each row using different goroutines
		go func(index int) {
			defer wg.Done()
			row := ""
			for x := 0; x < w; x++ {
				// get brightness of pixel
				gray_pixel := color.GrayModel.Convert(img.At(x, index))
				brightness := reflect.ValueOf(gray_pixel).FieldByName("Y").Uint()
				if inverted {
					brightness = 256 - brightness
				}

				// convert brightness of pixel to position in the character set
				l := uint64(len(charSet) - 1)
				char_index := uint(brightness * l / 255)
				row += string(charSet[char_index])
			}
			xs[index] = row
		}(y)
	}

	wg.Wait()

	return &xs
}

func GenerateAscii(img image.Image, charSet string, inverted bool, delay int) Frame {
	w, h := uint(img.Bounds().Dx()), uint(img.Bounds().Dy())
	h = uint(CorrectionRatio * float64(h))

	resizedImg := utils.ResizeImage(img, w, h)
	xs := convertToAscii(resizedImg, charSet, inverted)
	asciiString := strings.Join(*xs, "\n")

	return NewFrame(asciiString, delay, w, h)
}
