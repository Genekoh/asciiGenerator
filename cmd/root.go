package cmd

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"time"

	"github.com/Genekoh/asciiGenerator/pkg/ascii"
	"github.com/Genekoh/asciiGenerator/pkg/utils"
	"github.com/spf13/cobra"
	_ "golang.org/x/image/webp"
)

var (
	path, output, charSet string
	width                 uint
	inverted, read        bool

	rootCmd = &cobra.Command{
		Use:   "asciiGenerator",
		Short: "Converts an image to ASCII",
		Long:  "A CLI that can create ASCII from images",
		Run: func(cmd *cobra.Command, args []string) {
			command(path, output, charSet, width, inverted, read)
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "Path to image to be converted to ASCII")
	rootCmd.Flags().UintVarP(&width, "width", "w", 0, "Width of the ascii file")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Path to output file")
	rootCmd.Flags().StringVarP(&charSet, "char", "c", ascii.DefaultCharSet, "The character set for ")
	rootCmd.Flags().BoolVarP(&inverted, "inverted", "i", false, "Determines whether ascii is inverted or not")
	rootCmd.Flags().BoolVarP(&read, "read", "r", false, "Determines whether cli converts image to ascii or reads from existing ascii file")
}

func command(path, output, charSet string, width uint, inverted, read bool) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to open file")
		os.Exit(1)
	}

	if !read {
		img, ext, err := image.Decode(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch ext {
		case "jpeg", "png", "webp":
			if width != 0 {
				img = utils.ResizeImageAspect(img, width)
			}

			frame := ascii.GenerateAscii(img, charSet, inverted, 0)

			// output asciiString to terminal if no output is defined
			if strings.TrimSpace(output) == "" {
				fmt.Println(frame.AsciiString)
			} else {
				// else encode content to json and store to output file
				content := ascii.NewStoredContent([]ascii.Frame{frame}, -1)
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

		case "gif":
			_, err := file.Seek(0, 0)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			g, err := gif.DecodeAll(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			frames := make([]ascii.Frame, len(g.Image))
			// convert each gif image to a Frame object
			for i, x := range g.Image {
				// optimize to use go routines!!!
				var img image.Image
				if width != 0 {
					img = utils.ResizeImageAspect(x, width)
				} else {
					img = x
				}

				frames[i] = ascii.GenerateAscii(img, charSet, inverted, g.Delay[i])
			}

			// output ascii gif to terminal if no output is defined
			if strings.TrimSpace(output) == "" {
				n := g.LoopCount + 1
				if g.LoopCount == 0 {
					n = 900 // should be infinite??
				} else if g.LoopCount == -1 {
					n = 1
				}

				for i := 0; i < n; i++ {
					for j, f := range frames {
						if i != 0 || j != 0 {
							fmt.Printf("\x1b[%dF", f.Height-1)
						}

						fmt.Print(f.AsciiString)
						time.Sleep(time.Duration(g.Delay[j]) * 10 * time.Millisecond)
					}
					i++
				}
			} else {
				// else encode content to json and store to output file
				content := ascii.NewStoredContent(frames, g.LoopCount)
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
			fmt.Printf("%s is not a supported file type", ext)
			// os.Exit(1)
		}
	} else {
		var content ascii.StoredContent
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&content)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		n := content.LoopCount + 1
		if content.LoopCount == 0 {
			n = 900 // should be infinite??
		} else if content.LoopCount == -1 {
			n = 1
		}

		for i := 0; i < n; i++ {
			for j, f := range content.Frames {
				if i != 0 || j != 0 {
					fmt.Printf("\x1b[%dF", f.Height-1)
				}

				fmt.Print(f.AsciiString)
				time.Sleep(time.Duration(f.Delay) * 10 * time.Millisecond)
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
