package main

import (
	"fmt"
	"math"
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

func gaussian_filter(op [size][size]float64, data [image_size][image_size] float64) {
	for i int = 0; i<image_size; i++ {
		for j int = 0; j< image_size; j++{
			sum = 0
			for k int = 0; k<size; k++{
				for l int = 0; l<size; l++{
					sum += op[i-size/2+k][j-size/2+l]*data[i-size/2+k][j-size/2+l]
				}
			}
				
			switch {
			case 

			}
		}
	}
}

func main() {
	var gaussian_mat [size][size]float64
	fmt.Println("coucou")
	fmt.Println(gaussian_mat)
	fmt.Println(size/2)
	gaussian_mat = gaussian_matrix(1,gaussian_mat)
	fmt.Println(gaussian_mat)
}