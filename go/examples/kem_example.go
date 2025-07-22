// TOPAY-Z512 KEM Example
//
// This example demonstrates the usage of the TOPAY-Z512 Key Encapsulation Mechanism (KEM).
package examples

import (
	"fmt"

	"github.com/topay-foundation/topayz512/pkg/kem"
)

func KemExample() {
	fmt.Println("TOPAY-Z512 KEM Example")
	fmt.Println("=======================")

	// Generate a key pair
	keyPair, err := kem.Keygen()
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		return
	}

	fmt.Println("Generated key pair:")
	fmt.Printf("  Private key: %s...\n", keyPair.PrivateKey.String()[:16])
	fmt.Printf("  Public key: %s...\n", keyPair.PublicKey.String()[:16])

	// Encapsulate a shared secret
	ciphertext, sharedSecret1, err := kem.Encapsulate(&keyPair.PublicKey)
	if err != nil {
		fmt.Printf("Error encapsulating: %v\n", err)
		return
	}

	fmt.Println("\nEncapsulated shared secret:")
	fmt.Printf("  Ciphertext: %s...\n", ciphertext.String()[:32])
	fmt.Printf("  Shared secret: %s...\n", sharedSecret1.String()[:32])

	// Decapsulate the shared secret
	sharedSecret2, err := kem.Decapsulate(&keyPair.PrivateKey, ciphertext)
	if err != nil {
		fmt.Printf("Error decapsulating: %v\n", err)
		return
	}

	fmt.Println("\nDecapsulated shared secret:")
	fmt.Printf("  Shared secret: %s...\n", sharedSecret2.String()[:32])

	// Verify that the shared secrets match
	if sharedSecret1.String() == sharedSecret2.String() {
		fmt.Println("\nSuccess! The shared secrets match.")
	} else {
		fmt.Println("\nError! The shared secrets do not match.")
	}

	// Demonstrate serialization
	fmt.Println("\nSerialization Example:")

	// Convert ciphertext to hex and back
	ciphertextHex := ciphertext.String()
	ciphertext2, err := kem.CiphertextFromHex(ciphertextHex)
	if err != nil {
		fmt.Printf("Error deserializing ciphertext: %v\n", err)
		return
	}

	fmt.Printf("  Original ciphertext: %s...\n", ciphertext.String()[:32])
	fmt.Printf("  Deserialized ciphertext: %s...\n", ciphertext2.String()[:32])

	// Convert shared secret to hex and back
	sharedSecretHex := sharedSecret1.String()
	sharedSecret3, err := kem.SharedSecretFromHex(sharedSecretHex)
	if err != nil {
		fmt.Printf("Error deserializing shared secret: %v\n", err)
		return
	}

	fmt.Printf("  Original shared secret: %s...\n", sharedSecret1.String()[:32])
	fmt.Printf("  Deserialized shared secret: %s...\n", sharedSecret3.String()[:32])
}
