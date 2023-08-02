# asciiGenerator

A CLI tool for converting images to ascii art built using GO.

The currently supported image formats are:

- jpeg
- png
- webp
- gif

The asciiArt can either be displayed directly into the terminal or stored into a json file.

## Installation

The compiled file is already in the `./bin` directory.

## Building

If you want the build the tool yourself, clone this repository and run

```
go get
go build -o <where you want to output the binary>
```

or just

```
go get
go build
```

straight into your terminal.

## Usage

The general command will look like this

```
asciiGenerator -p <path to source file>   <other flags>
```

assuming that the executable is named `asciiGenerator.exe`.

### Flags:

**-p** \*\*Required\*\*

Takes in a string. Defines the path to the source input file for the CLI. In most cases, this will be the path to the image you want to convert to an image

Example:

```
asciiGenerator -p "testImages\amongus.jpg"
```

**-w**

Takes in a number. Defines the width of the ascii art output. The image will be scaled while preserving its aspect ratio. If not defined or set to 0, the CLI will keep the image in its original size.

Example:

```
asciiGenerator -p "testImages\amongus.jpg" -w 100
```

**-o**

Takes in a string. Defines the path to the output file to be stored as json. If not defined, the asciiArt will be outputted directly into the console.

Example:

```
asciiGenerator -p "testImages\amongus.jpg" -o ".\amongus.json"
```

**-c**

Takes in a string. Defines the character set used to assign each pixel of a certain brightness to. The default character set is `"MN8@O$Zbe*+!:.,  "`

Example:

```
asciiGenerator -p "testImages\amongus.jpg" -c "Wwli:,. "
```

**-i**

A boolean flag and its defualt value is false. If set to true, the CLI will reverse its character assignment to pixel.

Example:

```
asciiGenerator -p "testImages\amongus.jpg" -i
```

**-r**

Boolean flag. If set to true, the CLI will no longer convert images to ASCII art, but instead will read from a preexisting json file to display the ASCII to the console; the path flag will also represent the path the json file instead.

```
asciiGenerator -o ".\amongus.json" -r
```
