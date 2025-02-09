package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/anthonynsimon/bild/convolution"
)

// Fonction pour appliquer un filtre
const size int = 10

func image_opener() image.Image {
	// Open the image file
	file, err := os.Open("received.png")
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

func gaussian_filter(op []float64, data image.Image, rect image.Rectangle, wg *sync.WaitGroup, result chan<- image.Image) { // With convolve func
	defer wg.Done()

	subImage := data.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	kernel := convolution.Kernel{op, size, size}
	k := kernel.Normalized() // ça c'est pas bien, ça casse tout, plutot logique
	o := &convolution.Options{Bias: 0, Wrap: false}
	test := convolution.Convolve(subImage, k, o)
	// Marche pas parce que ça floute sans utiliser les images voisines...
	result <- test
}

func gradient_para(data image.Image, rect image.Rectangle, wg *sync.WaitGroup, result chan<- image.Image) {

	defer wg.Done()

	subImage := data.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	vect := []float64{1, 0, -1}
	vertical := convolution.Kernel{vect, 1, 3}
	horizontal := convolution.Kernel{vect, 3, 1}
	vert_norm := vertical.Normalized()  // ça aussi ça doit tout casser
	hor_norm := horizontal.Normalized() // ...
	o := &convolution.Options{Bias: 0, Wrap: false}
	gradient_vert := convolution.Convolve(subImage, vert_norm, o)
	gradient_hor := convolution.Convolve(subImage, hor_norm, o)
	var grad_dir [][]float64
	var grad_mag [][]float64
	var grad_dir_temp []float64
	var grad_mag_temp []float64
	for i := subImage.Bounds().Min.X; i < subImage.Bounds().Max.X; i += 1 {
		grad_dir_temp = []float64{}
		grad_mag_temp = []float64{}
		for j := subImage.Bounds().Min.Y; j < subImage.Bounds().Max.Y; j += 1 {
			r_v, _, _, _ := gradient_vert.RGBA64At(i, j).RGBA()
			r_h, _, _, _ := gradient_hor.RGBA64At(i, j).RGBA()
			grad_dir_temp = append(grad_dir_temp, math.Atan(float64(r_v)/float64(r_h)))
			grad_mag_temp = append(grad_mag_temp, math.Sqrt(math.Pow(float64(r_v), 2)+math.Pow(float64(r_h), 2)))
		}
		grad_dir = append(grad_dir, grad_dir_temp)
		grad_mag = append(grad_mag, grad_mag_temp)
	}
	bounds := subImage.Bounds()
	newImg := image.NewRGBA(bounds)

	// Process each pixel from the source image and write to the new image
	// Marche pas à cause des bordures ?
	for y := bounds.Min.Y; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X; x < bounds.Max.X-1; x++ {
			a := grad_mag[x][y]
			_, _, _, alpha := subImage.At(x, y).RGBA()
			newImg.Set(x, y, color.RGBA{uint8(a), uint8(a), uint8(a), uint8(alpha)})
		}
	}
	result <- newImg
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
	for y := bounds.Min.Y; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X; x < bounds.Max.X-1; x++ {
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

func applyFilters() image.Image {
	var gaussian_mat []float64
	gaussian_mat = gaussian_matrix(1.4, gaussian_mat)

	var pic = image_opener()

	var image_gray_scale image.Image = grey_scale(pic)
	// Diviser l'image en 4 parties
	width := image_gray_scale.Bounds().Dx()
	height := image_gray_scale.Bounds().Dy()

	// Définir les sous-images
	parts := []image.Rectangle{
		{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: width/2 + size/2, Y: height/2 + size/2}},
		{Min: image.Point{X: width/2 - size/2, Y: 0}, Max: image.Point{X: width, Y: height/2 + size/2}},
		{Min: image.Point{X: 0, Y: height/2 - size/2}, Max: image.Point{X: width/2 + size/2, Y: height}},
		{Min: image.Point{X: width/2 - size/2, Y: height/2 - size/2}, Max: image.Point{X: width, Y: height}},
	}

	var wg sync.WaitGroup
	result := make(chan image.Image, len(parts))

	for _, rect := range parts {
		wg.Add(1)
		go gaussian_filter(gaussian_mat, image_gray_scale, rect, &wg, result)
	}

	wg.Wait()
	close(result)

	finalImg := image.NewRGBA(image_gray_scale.Bounds())
	for subImg := range result {
		// Récupérer les dimensions de la sous-image
		subBounds := subImg.Bounds()

		// Placer la sous-image au bon endroit dans finalImg
		for y := subBounds.Min.Y; y < subBounds.Max.Y; y++ {
			for x := subBounds.Min.X; x < subBounds.Max.X; x++ {
				// Copier le pixel de la sous-image vers l'image finale
				c := subImg.At(x, y)
				finalImg.Set(x, y, c)
			}
		}
	}

	outFile, err := os.Create("gaussian_filter.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, finalImg) // Use jpeg.Encode for JPEG images
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}
	fmt.Println("Image saved as gaussian_filter.png")

	// var wg_2 sync.WaitGroup
	// result_2 := make(chan image.Image, len(parts))

	// for _, rect_2 := range parts {
	// 	wg_2.Add(1)
	// 	go gradient(finalImg, rect_2, &wg_2, result_2)
	// }

	// wg_2.Wait()
	// close(result)

	// finalImg_2 := image.NewRGBA(finalImg.Bounds())
	// for subImg := range result_2 {
	// 	// Récupérer les dimensions de la sous-image
	// 	subBounds := subImg.Bounds()

	// 	// Placer la sous-image au bon endroit dans finalImg_2
	// 	for y := subBounds.Min.Y; y < subBounds.Max.Y; y++ {
	// 		for x := subBounds.Min.X; x < subBounds.Max.X; x++ {
	// 			// Copier le pixel de la sous-image vers l'image finale
	// 			c := subImg.At(x, y)
	// 			finalImg_2.Set(x, y, c)
	// 		}
	// 	}
	// }

	finalImg_2 := gradient(finalImg)

	return finalImg_2
	// outFile_grad, err := os.Create("gradient_filter.png")
	// if err != nil {
	// 	fmt.Println("Error creating file:", err)
	// 	return
	// }
	// defer outFile_grad.Close()

	// err = png.Encode(outFile_grad, finalImg_2) // Use jpeg.Encode for JPEG images
	// if err != nil {
	// 	fmt.Println("Error encoding image:", err)
	// 	return
	// }
	// fmt.Println("Image saved as gradient_filter.png")
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connecté :", conn.RemoteAddr())

	// Ouvrir un fichier temporaire pour stocker l'image reçue
	tempFile, err := os.Create("received.png")
	if err != nil {
		fmt.Println("Erreur création du fichier:", err)
		return
	}
	defer tempFile.Close()

	// Lire l'image depuis la connexion et sauvegarder dans le fichier
	_, err = tempFile.ReadFrom(conn)
	if err != nil {
		fmt.Println("Erreur lecture de l'image:", err)
		return
	}

	// Appliquer le filtre
	image_created := applyFilters()

	// Sauvegarder l'image modifiée
	outFile, err := os.Create("processed.png")
	if err != nil {
		fmt.Println("Erreur création du fichier modifié:", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, image_created)
	if err != nil {
		fmt.Println("Erreur encodage de l'image:", err)
		return
	}

	// Envoyer l'image modifiée au client
	outFile.Seek(0, 0)
	_, err = outFile.WriteTo(conn)
	if err != nil {
		fmt.Println("Erreur envoi de l'image:", err)
	}

	fmt.Println("Image traitée envoyée au client")
}

func main() {
	// Lire le port en argument
	port := "8080" // Valeur par défaut
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// Vérifier que le port est un nombre valide
	if _, err := strconv.Atoi(port); err != nil {
		fmt.Println("Erreur: Le port doit être un nombre valide")
		return
	}

	// Démarrer le serveur TCP
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Erreur démarrage serveur:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Serveur en écoute sur le port", port, "...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur connexion client:", err)
			continue
		}
		go handleClient(conn) // Gérer chaque client dans une goroutine
	}
}
