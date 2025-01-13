package main

import (
	"fmt"
	"math"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

const size int = 5
const image_size int = 10

func gaussian_matrix(sigma float64, gaussian_m [size][size]float64) ([size][size]float64) {
	for i:=-size/2; i<=size/2; i++{
		for j:=-size/2; j<=size/2; j++{
			var pos_a float64 = float64(i) 
			var pos_b float64 = float64(j)
			gaussian_m[i+size/2][j+size/2] = math.Exp(-(pos_a*pos_a+pos_b*pos_b)/(2*sigma))/(2*math.Pi*(sigma*sigma))
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

func gaussian_filter(op [size][size]float64, data [image_size][image_size] float64) ([10][10]float64) {
	var image_with_gauss [image_size][image_size]float64
	for i := 0; i<image_size; i++ {
		for j := 0; j< image_size; j++{
			var sum float64 = 0
			for k := 0; k<size; k++{
				for l := 0; l<size; l++{
					switch{
					case i-size/2+k < 0, j-size/2+l < 0, i-size/2+k-1 > image_size, j-size/2+l-1 > image_size:
						continue
					default :
						sum += op[i-size/2+k][j-size/2+l]*data[i-size/2+k][j-size/2+l]
					}
				}
			}
			image_with_gauss[i][j] = sum
		}
	}
	return image_with_gauss
}

func main() {
	var gaussian_mat [size][size]float64
	fmt.Println("coucou")
	fmt.Println(gaussian_mat)
	fmt.Println(size/2)
	gaussian_mat = gaussian_matrix(1,gaussian_mat)
	fmt.Println(gaussian_mat)

	var image = [10][10]float64{{0.34, 0.85, 0.91, 0.44, 0.56, 0.32, 0.72, 0.61, 0.45, 0.88},
    {0.12, 0.67, 0.79, 0.24, 0.92, 0.85, 0.31, 0.77, 0.40, 0.62},
    {0.54, 0.33, 0.11, 0.29, 0.89, 0.45, 0.66, 0.38, 0.49, 0.93},
    {0.25, 0.78, 0.37, 0.14, 0.82, 0.60, 0.47, 0.53, 0.94, 0.50},
    {0.68, 0.99, 0.36, 0.23, 0.12, 0.55, 0.74, 0.18, 0.67, 0.81},
    {0.87, 0.46, 0.63, 0.21, 0.35, 0.95, 0.29, 0.40, 0.88, 0.76},
    {0.41, 0.83, 0.90, 0.39, 0.61, 0.44, 0.73, 0.27, 0.52, 0.19},
    {0.30, 0.98, 0.58, 0.20, 0.85, 0.71, 0.64, 0.49, 0.53, 0.97},
    {0.11, 0.69, 0.75, 0.56, 0.34, 0.93, 0.62, 0.48, 0.72, 0.57},
    {0.77, 0.43, 0.80, 0.15, 0.28, 0.66, 0.38, 0.84, 0.70, 0.31}}

	res_1 := gaussian_filter(gaussian_mat, image)

	fmt.Println(res_1)
}