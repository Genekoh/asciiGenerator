package cmd

import (
	"bytes"
	"fmt"
	"github.com/Genekoh/asciiGenerator/pkg/ascii"
	"github.com/Genekoh/asciiGenerator/pkg/image"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	path, output string
	width        uint

	rootCmd = &cobra.Command{
		Use:   "asciiGenerator",
		Short: "Video To Ascii Convertor",
		Long:  "A CLI that can create ASCII from videos",
		Run: func(cmd *cobra.Command, args []string) {
			command(path, width, output)
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "Path to image to be converted to ASCII")
	rootCmd.Flags().UintVarP(&width, "width", "w", 0, "Width of the ascii file")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Path to output file")
}

func outputAscii(file *os.File, newWidth uint) (*bytes.Buffer, error) {
	img, err := image.GetImage(file)
	if err != nil {
		return new(bytes.Buffer), err
	}

	if newWidth != 0 {
		img, _ = image.ResizeImageAspect(img, newWidth)
	}

	asciiBuffer := ascii.GenerateAscii(img)
	fmt.Println(asciiBuffer.String())
	return asciiBuffer, nil
}

func command(path string, width uint, output string) {
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
		buf, err := outputAscii(file, width)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) == "" {
			break
		}

		err = os.WriteFile(output, buf.Bytes(), 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		break

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
