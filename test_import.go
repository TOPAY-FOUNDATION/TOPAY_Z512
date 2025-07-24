package main

import (
	"fmt"
	"github.com/TOPAY-FOUNDATION/TOPAY_Z512/go"
)

func main() {
	// Test importing the module
	fmt.Println("Testing TOPAY Z512 Go module import...")
	
	// Test key pair generation
	privateKey, publicKey, err := topayz512.GenerateKeyPair()
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		return
	}
	
	fmt.Printf("Successfully generated key pair!\n")
	fmt.Printf("Private key: %s\n", privateKey.String()[:32]+"...")
	fmt.Printf("Public key: %s\n", publicKey.String()[:32]+"...")
	
	// Test hash computation
	data := []byte("Hello, TOPAY-Z512!")
	hash := topayz512.ComputeHash(data)
	fmt.Printf("Hash of test data: %s\n", hash.String()[:32]+"...")
	
	fmt.Println("Go module import and functionality test successful!")
}