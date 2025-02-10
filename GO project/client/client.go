package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// Vérifier que les arguments sont fournis
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run client.go <adresse:port> <chemin_image> <type_de_filtre>")
		return
	}

	serverAddr := os.Args[1] // Ex: "localhost:8080"
	imagePath := os.Args[2]  // Ex: "image.png"
	// filterType := os.Args[3] // Ex: "blur", "sobelX", "sobelY", "sharpen", "gaussianFilter"

	// Vérifier que le fichier existe
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Erreur ouverture de l'image:", err)
		return
	}
	defer file.Close()

	// Se connecter au serveur
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Erreur connexion au serveur:", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅ Connecté au serveur:", serverAddr)

	// Nouveau
	// _, err = conn.Write([]byte(filterType)) // Envoie du filtre choisi
	// if err != nil {
	// 	fmt.Println("Erreur envoi du type de filtre:", err)
	// 	return
	// }

	// // Envoyer l'image
	// _, err = file.Seek(0, 0) // Réinitialiser la lecture du fichier
	// if err != nil {
	// 	fmt.Println("Erreur réinitialisation du fichier:", err)
	// 	return
	// }
	// -Nouveau

	// Envoyer l'image au serveur
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Erreur envoi de l'image:", err)
		return
	}

	// 🔥 IMPORTANT : Fermer la connexion côté écriture pour signaler la fin
	conn.(*net.TCPConn).CloseWrite()

	fmt.Println("📤 Image envoyée, attente du retour...")

	// Créer un fichier pour stocker l'image modifiée
	outFile, err := os.Create("image_modifiee.png")
	if err != nil {
		fmt.Println("❌ Erreur création du fichier de sortie:", err)
		return
	}
	defer outFile.Close()

	// Lire les données en utilisant un buffer pour éviter le blocage
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Fin de la transmission
			}
			fmt.Println("❌ Erreur réception de l'image modifiée:", err)
			return
		}
		_, writeErr := outFile.Write(buffer[:n])
		if writeErr != nil {
			fmt.Println("❌ Erreur écriture du fichier:", writeErr)
			return
		}
	}

	fmt.Println("✅ Image modifiée reçue et enregistrée sous 'image_modifiee.png'")
}
