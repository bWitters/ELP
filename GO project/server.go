package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
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
	fmt.Println("image_opened")

	// Decode the image
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil
	}
	fmt.Println("Image format:", "png")
	return img
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

func splitImage(width, height, rows, cols, padding int) []image.Rectangle {
	parts := []image.Rectangle{}
	cellWidth := width / cols
	cellHeight := height / rows

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			x0 := col*cellWidth - padding
			y0 := row*cellHeight - padding
			x1 := (col+1)*cellWidth + padding
			y1 := (row+1)*cellHeight + padding

			// Empêcher les débordements hors de l'image
			if x0 < 0 {
				x0 = 0
			}
			if y0 < 0 {
				y0 = 0
			}
			if x1 > width {
				x1 = width
			}
			if y1 > height {
				y1 = height
			}

			parts = append(parts, image.Rect(x0, y0, x1, y1))
		}
	}

	return parts
}

// Convolution sur une image.Image avec un noyau
func ConvolveImage(img image.Image, kernel [][]float64) *image.Gray {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	kernelSize := len(kernel)
	pad := kernelSize / 2

	// Créer une nouvelle image en niveaux de gris pour stocker le résultat
	grayImg := image.NewGray(bounds)
	draw.Draw(grayImg, bounds, img, bounds.Min, draw.Src)

	result := image.NewGray(bounds)

	// Appliquer la convolution
	for y := pad; y < height-pad; y++ {
		for x := pad; x < width-pad; x++ {
			sum := 0.0
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					// Calculer les coordonnées de l'image source
					srcX := x + kx - pad
					srcY := y + ky - pad

					// Récupérer l'intensité en niveau de gris
					grayValue := float64(grayImg.GrayAt(srcX, srcY).Y)

					// Appliquer le filtre
					sum += grayValue * kernel[ky][kx]
				}
			}
			// Clamper les valeurs entre 0 et 255
			clamped := uint8(max(0, min(255, int(sum))))
			result.SetGray(x, y, color.Gray{Y: clamped})
		}
	}
	return result
}

// Fonctions utilitaires pour éviter les débordements
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Fonction pour créer un noyau de flou gaussien
func gaussianKernel(sigma float64) [][]float64 {
	// Déterminer la taille du noyau (taille impaire pour centrer sur le pixel)
	size := int(math.Ceil(6 * sigma)) // Taille ≈ 6σ
	if size%2 == 0 {
		size++ // S'assurer que la taille est impaire
	}

	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
	}

	// Calculer les valeurs du noyau
	sum := 0.0
	half := size / 2
	for i := -half; i <= half; i++ {
		for j := -half; j <= half; j++ {
			exponent := -(float64(i*i+j*j) / (2 * sigma * sigma))
			value := math.Exp(exponent) / (2 * math.Pi * sigma * sigma)
			kernel[i+half][j+half] = value
			sum += value
		}
	}

	// Normalisation (pour que la somme du noyau = 1)
	for i := range kernel {
		for j := range kernel[i] {
			kernel[i][j] /= sum
		}
	}

	return kernel
}

func BlurKernel(size int) [][]float64 {
	if size%2 == 0 || size < 3 {
		size = 3 // Forcer un kernel impair >= 3
	}

	kernel := make([][]float64, size)
	coeff := 1.0 / float64(size*size)

	for i := range kernel {
		kernel[i] = make([]float64, size)
		for j := range kernel[i] {
			kernel[i][j] = coeff
		}
	}
	return kernel
}

func SobelXKernel() [][]float64 {
	return [][]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
}

func SobelYKernel() [][]float64 {
	return [][]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}
}

func SharpenKernel() [][]float64 {
	return [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
}

func applyFilters(chosenFilter string) image.Image {
	pic := image_opener()
	imageGray := grey_scale(pic)

	// Vieux
	// kernel := gaussianKernel(5)
	// kernelSize := len(kernel)
	// -Vieux

	// Appliquer le filtre choisi
	var kernel [][]float64

	// Déterminer le noyau de convolution basé sur le filtre choisi
	switch chosenFilter {
	case "blur":
		kernel = BlurKernel(5) // Flou
	case "sobelX":
		kernel = SobelXKernel() // Noyau de Sobel X
	case "sobelY":
		kernel = SobelYKernel() // Noyau de Sobel Y
	case "sharpen":
		kernel = SharpenKernel() // Noyau de sharpening
	case "gaussianFilter":
		kernel = gaussianKernel(5) // Noyau de flou gaussien
	default:
		fmt.Println("Filtre non reconnu.")
		return nil
	}

	kernelSize := len(kernel)

	width := imageGray.Bounds().Dx()
	height := imageGray.Bounds().Dy()
	rows, cols := 1, 1 // Bizarre quand on divise l'image ça ralenti, cause potentiel : cout de la division, cout du padding, cout de la création des workers
	padding := kernelSize / 2

	start := time.Now()
	parts := splitImage(width, height, rows, cols, padding)

	var wg sync.WaitGroup
	result := make(chan image.Image, len(parts))

	// Appliquer le filtre en parallèle
	for _, rect := range parts {
		wg.Add(1)
		go func(r image.Rectangle) {
			defer wg.Done()
			subImg := ConvolveImage(imageGray, kernel) // Appliquer filtre
			result <- subImg
		}(rect)
	}

	wg.Wait()
	close(result)

	// Reconstruction de l'image
	finalImg := image.NewRGBA(imageGray.Bounds())
	for subImg := range result {
		subBounds := subImg.Bounds()
		for y := subBounds.Min.Y + padding; y < subBounds.Max.Y-padding; y++ {
			for x := subBounds.Min.X + padding; x < subBounds.Max.X-padding; x++ {
				c := subImg.At(x, y)
				finalImg.Set(x, y, c)
			}
		}
	}
	duration := time.Since(start)
	fmt.Println(duration)

	return finalImg
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connecté :", conn.RemoteAddr())

	// // Lire le type de filtre choisi par le client
	// var filterType [50]byte
	// _, err := conn.Read(filterType[:])
	// if err != nil {
	// 	fmt.Println("Erreur lecture du type de filtre:", err)
	// 	return
	// }
	// chosenFilter := string(filterType[:])

	// Ouvrir un fichier temporaire pour stocker l'image reçue
	tempFile, err := os.Create("received.png")
	if err != nil {
		fmt.Println("Erreur création du fichier:", err)
		return
	}
	defer tempFile.Close()

	// Lire les données en utilisant un buffer
	buffer := make([]byte, 4096) // Taille du buffer
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Fin de la transmission
			}
			fmt.Println("Erreur lecture des données:", err)
			return
		}
		_, writeErr := tempFile.Write(buffer[:n])
		if writeErr != nil {
			fmt.Println("Erreur écriture du fichier:", writeErr)
			return
		}
	}

	fmt.Println("✅ Fichier reçu, application du filtre...")

	fmt.Println("Fichier reçu, application du filtre...")

	// Appliquer le filtre
	image_created := applyFilters("gaussianFilter")

	fmt.Println("Filtre appliqué, sauvegarde...")

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

	// Réinitialiser la position du fichier avant l'envoi
	outFile.Seek(0, 0)

	// Envoyer l'image modifiée au client
	processedFile, err := os.Open("processed.png")
	if err != nil {
		fmt.Println("Erreur ouverture du fichier modifié:", err)
		return
	}
	defer processedFile.Close()

	_, err = io.Copy(conn, processedFile) // Envoi via io.Copy()
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
