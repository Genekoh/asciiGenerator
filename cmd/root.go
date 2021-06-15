package cmd

import (
	"fmt"
	"github.com/Genekoh/asciiGenerator/pkg/ascii"
	"github.com/Genekoh/asciiGenerator/pkg/image"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
	"os"
)

var (
	path  string
	width uint

	rootCmd = &cobra.Command{
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

func output(file *os.File, newWidth uint) error {
	img, err := image.GetImage(file)
	if err != nil {
		return err
	}

	if newWidth != 0 {
		img, _ = image.ResizeImageAspect(img, newWidth)
	}

	asciiImage := ascii.GenerateAscii(img)
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
