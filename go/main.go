// TOPAY-Z512 Examples Runner
//
// This file serves as the entry point to run all the examples in the TOPAY-Z512 Go implementation.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/topay-foundation/topayz512/pkg/hash"
	"github.com/topay-foundation/topayz512/pkg/keypair"
)

func main() {
	var choice int
	var err error

	// Check if an example number was provided as a command-line argument
	if len(os.Args) > 1 {
		choice, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Invalid example number: %s\n", os.Args[1])
			fmt.Println("Usage: go run main.go [example_number]")
			os.Exit(1)
		}
	} else {
		// No command-line argument, show the menu
		fmt.Println("TOPAY-Z512 Examples")
		fmt.Println("=================")
		fmt.Println("Please select an example to run:")
		fmt.Println("1. Hash Example")
		fmt.Println("2. Key Pair Example")
		fmt.Println("3. Private to Public Key Conversion Example")
		fmt.Println("0. Exit")
		fmt.Print("\nEnter your choice: ")

		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("\nInvalid input. Please enter a number.")
			return
		}

		fmt.Println() // Add a newline for better formatting
	}

	switch choice {
	case 0:
		fmt.Println("Exiting...")
		os.Exit(0)
	case 1:
		runHashExample()
	case 2:
		runKeypairExample()
	case 3:
		runPrivateToPublicExample()
	default:
		fmt.Printf("Invalid choice: %d\n", choice)
	}
}

func runHashExample() {
	fmt.Println("TOPAY-Z512 Hash Example with Time-based Input")

	// Get current time as input
	now := time.Now()
	timeStr := fmt.Sprintf("TOPAY-Z512 time: %d seconds", now.Unix())
	data := []byte(timeStr)
	hashValue := hash.New(data)

	fmt.Printf("Input: %s\n", string(data))
	fmt.Printf("Hash: %s\n", hashValue)
	fmt.Printf("Hash size: %d bytes\n\n", len(hashValue))

	// Hash combination with time components
	timeMillis := strconv.FormatInt(now.UnixNano()/int64(time.Millisecond), 10)
	data1 := []byte(fmt.Sprintf("TOPAY-%s", timeMillis))
	data2 := []byte(fmt.Sprintf("Z512-%d", now.UnixNano()%1000000))

	combinedHash := hash.Combine(data1, data2)

	fmt.Printf("Input 1: %s\n", string(data1))
	fmt.Printf("Input 2: %s\n", string(data2))
	fmt.Printf("Combined Hash: %s\n", combinedHash)

	// Concatenated hash (different from combined hash)
	concatenated := append([]byte{}, data1...)
	concatenated = append(concatenated, data2...)
	concatHash := hash.New(concatenated)

	fmt.Printf("Concatenated Hash: %s\n", concatHash)
	fmt.Printf("Are they equal? %t\n\n", combinedHash == concatHash)

	// Convenience functions
	hashBytes := hash.Sum512(data)
	fmt.Printf("Hash bytes (first 8): %v\n", hashBytes[:8])

	combinedBytes := hash.SumCombine(data1, data2)
	fmt.Printf("Combined hash bytes (first 8): %v\n", combinedBytes[:8])

	// Hex conversion
	hex := hashValue.String()
	fmt.Printf("\nHex string: %s\n", hex)

	hashFromHex, err := hash.FromHex(hex)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Converted back from hex: %s\n", hashFromHex)
		fmt.Printf("Equal to original: %t\n", hashValue == hashFromHex)
	}
}

func runKeypairExample() {
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

func runPrivateToPublicExample() {
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