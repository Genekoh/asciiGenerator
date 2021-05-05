package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"reflect"

	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
	"github.com/spf13/cobra"
)

var (
	path  string
	width uint

	asciiCharacters = "MN8@O$Zbe*+!:.,  "
	asciiTable      = []byte(asciiCharacters)
	rootCmd         = &cobra.Command{
		Use:   "asciiGenerator",
		Short: "Video To Ascii Convertor",
		Long:  "A CLI that can create ASCII from videos",
		Run: func(cmd *cobra.Command, args []string) {
			command(path, width)
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "Path to your file to convert to ASCII")
	rootCmd.Flags().UintVarP(&width, "width", "w", 0, "Width of the ascii file")
}

func getImage(file *os.File) (image.Image, error) {
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

func getImageDimensions(file *os.File) (w uint, h uint, e error) {
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

func resizeImage(image image.Image, width, height, newWidth uint) (image.Image, uint) {
	ratio := float64(height) / float64(width)
	newWidthFloat := float64(newWidth)
	newHeight := uint(ratio * newWidthFloat)

	resizedImage := resize.Resize(newWidth, newHeight, image, resize.Bicubic)

	return resizedImage, newHeight
}

func generateAscii(img image.Image, w, h int) string {
	buffer := new(bytes.Buffer)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := color.GrayModel.Convert(img.At(x, y))
			brightness := reflect.ValueOf(p).FieldByName("Y").Uint()
			length := uint64(len(asciiCharacters) - 1)
			index := int(brightness * length / 255)
			buffer.WriteByte(asciiTable[index])
		}
		buffer.WriteByte('\n')
	}

	return buffer.String()
}

func output(file *os.File, newWidth uint) error {
	var finalWidth, finalHeight uint
	img, err := getImage(file)
	if err != nil {
		return err
	}
	w, h, err := getImageDimensions(file)
	if err != nil {
		return err
	}

	if newWidth != 0 {
		i, newHeight := resizeImage(img, w, h, newWidth)
		img = i
		finalWidth = newWidth
		finalHeight = newHeight
	} else {
		finalWidth = w
		finalHeight = h
	}
	asciiImage := generateAscii(img, int(finalWidth), int(finalHeight))
	fmt.Println(asciiImage)
	return nil
}

func command(path string, width uint) {
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		fmt.Println("Unable to find file")
		os.Exit(1)
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to open file")
		os.Exit(1)
	}

	switch extension := mime.Extension(); extension {
	case ".jpg", ".jpeg", ".png":
		err := output(file, width)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		fmt.Println("Not a supported file type")
		os.Exit(1)
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
