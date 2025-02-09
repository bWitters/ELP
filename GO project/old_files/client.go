package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Vérifier que les arguments sont fournis
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run client.go <adresse:port> <chemin_image>")
		return
	}

	serverAddr := os.Args[1] // Ex: "localhost:8080"
	imagePath := os.Args[2]  // Ex: "image.jpg"

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

	// Envoyer l'image au serveur
	_, err = file.Seek(0, 0) // Réinitialiser la lecture du fichier
	if err != nil {
		fmt.Println("Erreur réinitialisation du fichier:", err)
		return
	}

	_, err = file.WriteTo(conn) // Envoyer l'image
	if err != nil {
		fmt.Println("Erreur envoi de l'image:", err)
		return
	}

	// Réception de l'image modifiée
	outFile, err := os.Create("image_modifiee.png")
	if err != nil {
		fmt.Println("Erreur création du fichier de sortie:", err)
		return
	}
	defer outFile.Close()

	_, err = outFile.ReadFrom(conn)
	if err != nil {
		fmt.Println("Erreur réception de l'image modifiée:", err)
		return
	}

	fmt.Println("Image modifiée reçue et enregistrée sous 'image_modifiee.jpg'")
}
