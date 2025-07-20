// TOPAY-Z512 Key Pair Example
//
// This example demonstrates how to generate and use key pairs in TOPAY-Z512.
// It shows key generation, public key derivation, and hex conversion.
package examples

import (
	"fmt"
	"log"

	"github.com/topay-foundation/topayz512/pkg/keypair"
)

// RunKeypair demonstrates the key pair functionality
func RunKeypair() {
	fmt.Println("TOPAY-Z512 Key Pair Example")

	// Generate a new key pair
	kp, err := keypair.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	fmt.Println("Generated Key Pair:")
	fmt.Printf("Private Key: %s\n", kp.PrivateKey.String())
	fmt.Printf("Public Key:  %s\n\n", kp.PublicKey.String())

	// Demonstrate deriving a public key from a private key
	derivedPublicKey := keypair.DerivePublicKey(kp.PrivateKey)
	fmt.Printf("Derived Public Key: %s\n", derivedPublicKey.String())
	fmt.Printf("Keys match: %v\n\n", derivedPublicKey.String() == kp.PublicKey.String())

	// Demonstrate hex conversion
	privateHex := kp.PrivateKey.String()
	recoveredPrivateKey, err := keypair.PrivateKeyFromHex(privateHex)
	if err != nil {
		log.Fatalf("Failed to recover private key from hex: %v", err)
	}

	publicHex := kp.PublicKey.String()
	recoveredPublicKey, err := keypair.PublicKeyFromHex(publicHex)
	if err != nil {
		log.Fatalf("Failed to recover public key from hex: %v", err)
	}

	fmt.Println("Hex Conversion Test:")
	fmt.Printf("Private Key Recovered: %v\n", privateHex == recoveredPrivateKey.String())
	fmt.Printf("Public Key Recovered:  %v\n", publicHex == recoveredPublicKey.String())

	// Create a new key pair from existing keys
	newKeyPair := keypair.NewKeyPair(recoveredPrivateKey, recoveredPublicKey)
	fmt.Println("\nRecreated Key Pair:")
	fmt.Printf("Private Key: %s\n", newKeyPair.PrivateKey.String())
	fmt.Printf("Public Key:  %s\n", newKeyPair.PublicKey.String())
}