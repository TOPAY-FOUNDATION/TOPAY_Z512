// TOPAY-Z512 Private to Public Key Conversion Example
//
// This example demonstrates how to convert a private key to a public key
// using the TOPAY-Z512 library.

package examples

import (
	"fmt"

	"github.com/topay-foundation/topayz512/pkg/keypair"
)

// RunPrivateToPublic demonstrates the private to public key conversion functionality
func RunPrivateToPublic() {
	fmt.Println("TOPAY-Z512 Private to Public Key Conversion Example")

	// Generate a new key pair to get a private key
	kp, err := keypair.GenerateKeyPair()
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		return
	}
	privateKey := kp.PrivateKey
	fmt.Printf("Generated private key: %s\n", privateKey.String())

	// Convert private key to public key using the convenience function
	publicKey := keypair.PrivateToPublic(privateKey)
	fmt.Printf("Derived public key: %s\n", publicKey.String())

	// Create a private key from a hex string
	hexPrivateKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	privateKey, err = keypair.PrivateKeyFromHex(hexPrivateKey)
	if err != nil {
		fmt.Printf("Error creating private key from hex: %v\n", err)
		return
	}
	fmt.Printf("\nPredefined private key: %s\n", privateKey.String())

	// Convert private key to public key
	publicKey = keypair.PrivateToPublic(privateKey)
	fmt.Printf("Derived public key: %s\n", publicKey.String())
}