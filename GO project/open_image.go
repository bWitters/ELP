package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func main(){
	// Open the image file
	file, err := os.Open("image1.png") // Replace with your image file path
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}
	fmt.Println("Image format:", format)

	// Work with the image (example: get dimensions)
	bounds := img.Bounds()
	fmt.Printf("Width: %d, Height: %d\n", bounds.Dx(), bounds.Dy())

	// Optional: Re-encode and save the image (example for PNG)
	outFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, img) // Use jpeg.Encode for JPEG images
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
	fmt.Println("Image saved as output.png")
}