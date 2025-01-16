package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const size int = 5
const image_size int = 10

func image_opener() image.Image {
	// Open the image file
	file, err := os.Open("Valve_original_(1).PNG") // Replace with your image file path
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	fmt.Println(file)
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

func gaussian_matrix(sigma float64, gaussian_m [size][size]float64) [size][size]float64 {
	for i := -size / 2; i <= size/2; i++ {
		for j := -size / 2; j <= size/2; j++ {
			var pos_a float64 = float64(i)
			var pos_b float64 = float64(j)
			gaussian_m[i+size/2][j+size/2] = math.Exp(-(pos_a*pos_a+pos_b*pos_b)/(2*sigma)) / (2 * math.Pi * (sigma * sigma))
		}
	}
	return gaussian_m
}

/* func gaussian_filter(op [size][size]float64, data [image_size][image_size] float64) {
	// https://janth.home.xs4all.nl/Publications/html/conv2d.html for 2D convolution algorithm
	oplx = size
	oply = size
	ny = image_size
	nx = image_size

	hoplx = (oplx + 1)/2
	hoply = (oply+1)/2

	for iy = 0; iy <ny ; iy ++ {
		starty = math.Max(iy-hoply+1, 0)
		endy   = math.Min(iy+hoply, ny)

		dumr = dumi = 0.0
		k = math.Max(hoply-1-iy, 0)
		for i = starty; i < endy; i++ {
			l = math.Max(hoplx-1-ix, 0)
			for j = startx; j < endx; j++ {
				dumr += data[i*nx+j].r*opx[k*oplx+l].r
				dumr += data[i*nx+j].i*opx[k*oplx+l].i
				dumi += data[i*nx+j].i*opx[k*oplx+l].r
				dumi -= data[i*nx+j].r*opx[k*oplx+l].i
				l++
			}
			k++
		}
		convr[iy*nx+ix].r = dumr
		convr[iy*nx+ix].i = dumi
	}
} */

func gaussian_filter(op [size][size]float64, data image.Image) image.Image {
	var bounds = data.Bounds()
	var image_with_gauss = image.NewRGBA(bounds)
	for i := 0; i < image_size; i++ {
		for j := 0; j < image_size; j++ {
			var sum float64 = 0
			for k := 0; k < size; k++ {
				for l := 0; l < size; l++ {
					switch {
					case i-size/2+k < 0, j-size/2+l < 0, i-size/2+k-1 > image_size, j-size/2+l-1 > image_size:
						continue
					default:
						g, _, _, _ := data.At(i+k, j+l).RGBA()
						fmt.Println(g)
						sum += op[k][l] * float64(g)
						fmt.Println(sum)
					}
				}
			}
			fmt.Println(uint8(sum))
			image_with_gauss.Set(i, j, color.Gray{uint8(sum)})
			fmt.Println(image_with_gauss)
		}
	}
	return image_with_gauss
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
	var gaussian_mat [size][size]float64
	fmt.Println("coucou")
	fmt.Println(gaussian_mat)
	fmt.Println(size / 2)
	gaussian_mat = gaussian_matrix(1, gaussian_mat)

	var pic = image_opener()

	var image_gray_scale image.Image = grey_scale(pic)

	fmt.Println(image_gray_scale)

	fmt.Println(image_gray_scale.At(10, 10).RGBA())

	image_gauss := gaussian_filter(gaussian_mat, image_gray_scale)

	//fmt.Println(image_gauss)

	outFile, err := os.Create("output.png")
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
	fmt.Println("Image saved as output.png")
	// res_1 := gaussian_filter(gaussian_mat, image)

	// fmt.Println(res_1)
}
