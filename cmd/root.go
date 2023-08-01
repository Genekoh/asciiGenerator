package cmd

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"github.com/Genekoh/asciiGenerator/pkg/ascii"
	imagePkg "github.com/Genekoh/asciiGenerator/pkg/image"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
)

var (
	path, output string
	width        uint
	read         bool

	rootCmd = &cobra.Command{
		Use:   "asciiGenerator",
		Short: "Converts an image to ASCII",
		Long:  "A CLI that can create ASCII from images",
		Run: func(cmd *cobra.Command, args []string) {
			command(path, width, output, read)
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "Path to image to be converted to ASCII")
	rootCmd.Flags().UintVarP(&width, "width", "w", 0, "Width of the ascii file")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Path to output file")
	rootCmd.Flags().BoolVarP(&read, "read", "r", false, "Determines whether cli converts image to ascii or reads from existing ascii file")
}

func command(path string, width uint, output string, read bool) {
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		fmt.Println("Unable to find file")
		os.Exit(1)
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to open file")
		os.Exit(1)
	}

	if !read {
		switch extension := mime.Extension(); extension {
		case ".jpg", ".jpeg", ".png", ".webp":
			img, err := imagePkg.GetImage(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if width != 0 {
				img, _ = imagePkg.ResizeImageAspect(img, width)
			}
			buf := ascii.GenerateAscii(img)

			// output ascii to terminal if no output is defined
			if strings.TrimSpace(output) == "" {
				fmt.Println(buf.String())
			} else {
				err = os.WriteFile(output, buf.Bytes(), 0666)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

		case ".gif":
			gif, err := imagePkg.GetGif(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var asciiStrings []string
			var imgs []image.Image
			for i := 0; i < len(gif.Image); i++ {
				var img image.Image = gif.Image[i]
				if width != 0 {
					img, _ = imagePkg.ResizeImageAspect(gif.Image[i], width)
				}

				buf := ascii.GenerateAscii(img)
				asciiStrings = append(asciiStrings, buf.String())
				imgs = append(imgs, img)
			}

			// output ascii to terminal if no output is defined
			if strings.TrimSpace(output) == "" {
				for i := 0; i <= 20; i++ {
					for j, img := range imgs {
						if i != 0 || j != 0 {
							h := uint(float64(img.Bounds().Dy()) * ascii.CorrectionRatio)
							fmt.Printf("\x1b[%dF", h)
						}

						fmt.Print(asciiStrings[j])
						time.Sleep(time.Duration(gif.Delay[j]) * 10 * time.Millisecond)
					}
				}
			} else {
				if !strings.HasSuffix(output, ".json") {
					output += ".json"
				}

				content := ascii.NewStoredAscii(asciiStrings, gif.Delay, uint(imgs[0].Bounds().Dy()))
				json_content, err := json.Marshal(content)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				err = os.WriteFile(output, json_content, 0666)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

		default:
			fmt.Println("extension:\t", extension)
			fmt.Println("Not a supported file type")
			os.Exit(1)
		}
	} else {
		var content ascii.StoredAscii
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&content)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for i := 0; i < 20; i++ {
			for j, img := range content.Frames {
				if i != 0 || j != 0 {
					h := uint(float64(content.Height) * ascii.CorrectionRatio)
					fmt.Printf("\x1b[%dF", h)
				}

				fmt.Print(img)
				time.Sleep(time.Duration(content.Delay[j]) * 10 * time.Millisecond)
			}
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
