package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	// Create a 32-byte slice to store the secret
	key := make([]byte, 32)

	// Fill the slice with cryptographically secure random bytes
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalf("Failed to generate random key: %v", err)
	}

	// Encode the key as a base64 string for easy copying and use
	secret := base64.StdEncoding.EncodeToString(key)

	fmt.Println("Generated 32-byte secret key:")
	fmt.Println(secret)
}
