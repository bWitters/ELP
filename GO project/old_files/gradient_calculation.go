package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

// Helper function to load an image from a file
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Helper function to convert an image to grayscale
func toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y))
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

// Sobel filter implementation
func sobelFilters(img *image.Gray) (*image.Gray, *image.Gray) {
	bounds := img.Bounds()
	gradient := image.NewGray(bounds)
	theta := image.NewGray(bounds)

	kx := [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	ky := [3][3]int{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1},
	}

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			var gx, gy int
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					gray := img.GrayAt(x+j, y+i).Y
					gx += int(gray) * kx[i+1][j+1]
					gy += int(gray) * ky[i+1][j+1]
				}
			}

			magnitude := math.Sqrt(float64(gx*gx + gy*gy))
			angle := math.Atan2(float64(gy), float64(gx))

			gradient.SetGray(x, y, color.Gray{Y: uint8(magnitude)})
			theta.SetGray(x, y, color.Gray{Y: uint8(angle * 255 / (2 * math.Pi))})
		}
	}

	return gradient, theta
}

func main_grad() {
	img, err := loadImage("trump.jpg")
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	grayImg := toGrayscale(img)
	gradient, theta := sobelFilters(grayImg)

	// Save the gradient and theta images
	saveImage("gradient.jpg", gradient)
	saveImage("theta.jpg", theta)

	fmt.Println("Gradient and theta images saved.")
}

// Helper function to save an image to a file
func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}
