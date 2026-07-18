package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func manipulateImage(src image.Image) *image.Gray {
	bounds := src.Bounds()
	row, col := bounds.Max.Y, bounds.Max.X

	greyImage := image.NewGray(bounds)
	for i := bounds.Min.Y; i < row; i++ {
		for j := bounds.Min.X; j < col; j++ {
			r, g, b, _ := src.At(j, i).RGBA() // 16 bit

			// convert to 8 bit
			R := r >> 8
			G := g >> 8
			B := b >> 8

			grayScale := (0.299 * float64(R)) + (0.587 * float64(G)) + (0.114 * float64(B))
			greyImage.Set(j, i, color.Gray{Y: uint8(grayScale)})
		}
	}

	return greyImage

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <image-file>")
		return
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Cannot decode the image:", err)
		return
	}

	greyImage := manipulateImage(img)

	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)
	outputFile := base + "-output.png"

	newFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Unable to create output file:", err)
		return
	}
	defer newFile.Close()

	if err := png.Encode(newFile, greyImage); err != nil {
		fmt.Println("Unable to encode PNG:", err)
		return
	}

	fmt.Println("Saved:", outputFile)
}
