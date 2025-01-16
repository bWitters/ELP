package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/anthonynsimon/bild/convolution"
)

const size int = 5

func image_opener() image.Image {
	// Open the image file
	file, err := os.Open("Large_Scaled_Forest_Lizard.png") // Replace with your image file path
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Decode the image
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil
	}
	fmt.Println("Image format:", "png")
	return img
}

func gaussian_matrix(sigma float64, gaussian_m []float64) []float64 {
	k := 0
	for i := -size / 2; i <= size/2; i++ {
		for j := -size / 2; j <= size/2; j++ {
			var pos_a float64 = float64(i)
			var pos_b float64 = float64(j)
			gaussian_m = append(gaussian_m, math.Exp(-(pos_a*pos_a+pos_b*pos_b)/(2*sigma))/(2*math.Pi*(sigma*sigma)))
			k += 1
		}
	}
	return gaussian_m
}

func gaussian_filter(op []float64, data image.Image) image.Image { // With convolve func
	kernel := convolution.Kernel{op, size, size}
	k := kernel.Normalized()
	o := &convolution.Options{Bias: 0, Wrap: false}
	test := convolution.Convolve(data, k, o)
	return test
}

func gradient(data image.Image) image.Image {
	vect := []float64{1, 0, -1}
	vertical := convolution.Kernel{vect, 1, 3}
	horizontal := convolution.Kernel{vect, 3, 1}
	vert_norm := vertical.Normalized()
	hor_norm := horizontal.Normalized()
	o := &convolution.Options{Bias: 0, Wrap: false}
	gradient_vert := convolution.Convolve(data, vert_norm, o)
	gradient_hor := convolution.Convolve(data, hor_norm, o)
	var grad_dir [][]float64
	var grad_mag [][]float64
	var grad_dir_temp []float64
	var grad_mag_temp []float64
	for i := data.Bounds().Min.X; i < data.Bounds().Max.X; i += 1 {
		grad_dir_temp = []float64{}
		grad_mag_temp = []float64{}
		for j := data.Bounds().Min.Y; j < data.Bounds().Max.Y; j += 1 {
			r_v, _, _, _ := gradient_vert.RGBA64At(i, j).RGBA()
			r_h, _, _, _ := gradient_hor.RGBA64At(i, j).RGBA()
			grad_dir_temp = append(grad_dir_temp, math.Atan(float64(r_v)/float64(r_h)))
			grad_mag_temp = append(grad_mag_temp, math.Sqrt(math.Pow(float64(r_v), 2)+math.Pow(float64(r_h), 2)))
		}
		grad_dir = append(grad_dir, grad_dir_temp)
		grad_mag = append(grad_mag, grad_mag_temp)
	}
	bounds := data.Bounds()
	newImg := image.NewRGBA(bounds)

	// Process each pixel from the source image and write to the new image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			a := grad_mag[x][y]
			_, _, _, alpha := data.At(x, y).RGBA()
			newImg.Set(x, y, color.RGBA{uint8(a), uint8(a), uint8(a), uint8(alpha)})
		}
	}
	return newImg
}

func grey_scale(pic image.Image) image.Image {
	bounds := pic.Bounds()
	newImg := image.NewRGBA(bounds)

	// Process each pixel from the source image and write to the new image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, y, color.Gray16Model.Convert(pic.At(x, y)))
		}
	}
	return newImg
}

func main() {
	var gaussian_mat []float64
	gaussian_mat = gaussian_matrix(1.4, gaussian_mat)

	var pic = image_opener()

	var image_gray_scale image.Image = grey_scale(pic)

	image_gauss := gaussian_filter(gaussian_mat, image_gray_scale)

	outFile, err := os.Create("gaussian_filter.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, image_gauss) // Use jpeg.Encode for JPEG images
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
	fmt.Println("Image saved as gaussian_filter.png")

	grad_image := gradient(image_gauss)

	outFile_grad, err := os.Create("gradient_filter.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile_grad.Close()

	err = png.Encode(outFile_grad, grad_image) // Use jpeg.Encode for JPEG images
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
	fmt.Println("Image saved as gradient_filter.png")
}
