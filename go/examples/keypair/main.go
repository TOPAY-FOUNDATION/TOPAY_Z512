package main

import (
	"fmt"
	"log"

	"github.com/TOPAY-FOUNDATION/TOPAY_Z512/go"
)

func main() {
	fmt.Println("=== TOPAY-Z512 Go Implementation - Key Pair Operations ===")
	fmt.Println()

	// Basic key pair generation
	fmt.Println("1. Basic Key Pair Generation:")
	privateKey, publicKey, err := topayz512.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	fmt.Printf("   Private Key: %s\n", privateKey.String())
	fmt.Printf("   Public Key:  %s\n", publicKey.String())
	fmt.Printf("   Key pair valid: %v\n", topayz512.VerifyKeyPair(privateKey, publicKey))
	fmt.Println()

	// Key pair from seed
	fmt.Println("2. Deterministic Key Pair Generation:")
	seed := []byte("this is a deterministic seed for key generation")
	privateKey2, publicKey2, err := topayz512.GenerateKeyPairFromSeed(seed)
	if err != nil {
		log.Fatalf("Failed to generate key pair from seed: %v", err)
	}

	fmt.Printf("   Private Key: %s\n", privateKey2.String())
	fmt.Printf("   Public Key:  %s\n", publicKey2.String())
	fmt.Printf("   Key pair valid: %v\n", topayz512.VerifyKeyPair(privateKey2, publicKey2))
	fmt.Println()

	// Verify same seed produces same keys
	privateKey3, publicKey3, _ := topayz512.GenerateKeyPairFromSeed(seed)
	fmt.Printf("   Same seed produces same keys: %v\n",
		topayz512.PrivateKeyEqual(privateKey2, privateKey3) &&
			topayz512.PublicKeyEqual(publicKey2, publicKey3))
	fmt.Println()

	// Public key derivation
	fmt.Println("3. Public Key Derivation:")
	derivedPublicKey := topayz512.DerivePublicKey(privateKey)
	fmt.Printf("   Original Public Key:  %s\n", publicKey.String())
	fmt.Printf("   Derived Public Key:   %s\n", derivedPublicKey.String())
	fmt.Printf("   Keys match: %v\n", topayz512.PublicKeyEqual(publicKey, derivedPublicKey))
	fmt.Println()

	// Key validation
	fmt.Println("4. Key Validation:")
	fmt.Printf("   Private key valid: %v\n", topayz512.IsValidPrivateKey(privateKey))
	fmt.Printf("   Public key valid:  %v\n", topayz512.IsValidPublicKey(publicKey))

	// Test with invalid keys
	invalidPrivateKey := topayz512.PrivateKey{}
	invalidPublicKey := topayz512.PublicKey{}
	fmt.Printf("   Invalid private key: %v\n", topayz512.IsValidPrivateKey(invalidPrivateKey))
	fmt.Printf("   Invalid public key:  %v\n", topayz512.IsValidPublicKey(invalidPublicKey))
	fmt.Println()

	// Batch key generation
	fmt.Println("5. Batch Key Generation:")
	batchSize := 10
	privateKeys, publicKeys, err := topayz512.BatchGenerateKeyPairs(batchSize)
	if err != nil {
		log.Fatalf("Failed to generate key pairs in batch: %v", err)
	}

	fmt.Printf("   Generated %d key pairs\n", len(privateKeys))

	// Verify all generated key pairs
	validCount := 0
	for i := 0; i < batchSize; i++ {
		if topayz512.VerifyKeyPair(privateKeys[i], publicKeys[i]) {
			validCount++
		}
	}
	fmt.Printf("   Valid key pairs: %d/%d\n", validCount, batchSize)
	fmt.Println()

	// Key derivation from password
	fmt.Println("6. Key Derivation from Password:")
	password := []byte("my_secure_password_123")
	salt := []byte("random_salt_value")
	iterations := 10000
	derivedKey, err := topayz512.DeriveKeyFromPassword(password, salt, iterations)
	if err != nil {
		log.Fatalf("Failed to derive key from password: %v", err)
	}

	fmt.Printf("   Password: %s\n", string(password))
	fmt.Printf("   Salt: %x\n", salt)
	fmt.Printf("   Derived Key: %s\n", derivedKey.String())
	fmt.Printf("   Key valid: %v\n", topayz512.IsValidPrivateKey(derivedKey))
	fmt.Println()

	// Child key derivation
	fmt.Println("7. Child Key Derivation:")
	parentKey := privateKey
	childIndex := uint32(42)

	childKey := topayz512.DeriveChildKey(parentKey, childIndex)

	fmt.Printf("   Parent Key: %s\n", parentKey.String())
	fmt.Printf("   Child Index: %d\n", childIndex)
	fmt.Printf("   Child Key: %s\n", childKey.String())
	fmt.Printf("   Child key valid: %v\n", topayz512.IsValidPrivateKey(childKey))

	// Verify different child indices produce different keys
	childKey2 := topayz512.DeriveChildKey(parentKey, childIndex+1)
	fmt.Printf("   Different indices produce different keys: %v\n",
		!topayz512.PrivateKeyEqual(childKey, childKey2))
	fmt.Println()

	// HD Wallet generation
	fmt.Println("8. HD Wallet Generation:")
	walletSeed := []byte("this is a master seed for HD wallet generation that should be long enough")
	walletDepth := 5

	hdWallet, err := topayz512.GenerateHDWallet(walletSeed, walletDepth)
	if err != nil {
		log.Fatalf("Failed to generate HD wallet: %v", err)
	}

	fmt.Printf("   Wallet depth: %d\n", walletDepth)
	fmt.Printf("   Generated key pairs: %d\n", len(hdWallet))

	for i, keyPair := range hdWallet {
		fmt.Printf("   Key %d: %s -> %s\n", i,
			keyPair.PrivateKey.String()[:16]+"...",
			keyPair.PublicKey.String()[:16]+"...")
	}
	fmt.Println()

	// Key recovery demonstration
	fmt.Println("9. Key Recovery:")
	// Simulate a scenario where we have a public key and need to verify ownership
	challengeData := []byte("prove you own this key")

	// In a real scenario, this would be a signature operation
	// For demonstration, we'll use the key to derive a response
	response := topayz512.ComputeHash(append(privateKey.Bytes(), challengeData...))

	// Verify the response using the public key
	expectedResponse := topayz512.ComputeHash(append(privateKey.Bytes(), challengeData...))

	fmt.Printf("   Challenge: %x\n", challengeData)
	fmt.Printf("   Response: %s\n", response.String())
	fmt.Printf("   Verification: %v\n", topayz512.HashEqual(response, expectedResponse))
	fmt.Println()

	// Comprehensive validation
	fmt.Println("10. Comprehensive Key Pair Validation:")
	keyPair := topayz512.KeyPair{PrivateKey: &privateKey, PublicKey: &publicKey}
	validation := topayz512.ValidateKeyPairIntegrity(keyPair)
	fmt.Printf("    Validation result: %v\n", validation == nil)

	// Test with mismatched keys
	_, otherPublicKey, _ := topayz512.GenerateKeyPair()
	invalidKeyPair := topayz512.KeyPair{PrivateKey: &privateKey, PublicKey: &otherPublicKey}
	invalidValidation := topayz512.ValidateKeyPairIntegrity(invalidKeyPair)
	fmt.Printf("    Invalid pair validation: %v\n", invalidValidation == nil)
	fmt.Println()

	// Secure key erasure
	fmt.Println("11. Secure Key Erasure:")
	tempPrivateKey, tempPublicKey, _ := topayz512.GenerateKeyPair()

	fmt.Printf("    Before erasure - Private key valid: %v\n",
		topayz512.IsValidPrivateKey(tempPrivateKey))

	topayz512.SecureErasePrivateKey(&tempPrivateKey)

	fmt.Printf("    After erasure - Private key valid: %v\n",
		topayz512.IsValidPrivateKey(tempPrivateKey))
	fmt.Printf("    Public key still valid: %v\n",
		topayz512.IsValidPublicKey(tempPublicKey))
	fmt.Println()

	// Serialization and deserialization
	fmt.Println("12. Key Serialization:")

	// Serialize keys
	privateKeyBytes := privateKey.Bytes()
	publicKeyBytes := publicKey.Bytes()

	fmt.Printf("    Private key size: %d bytes\n", len(privateKeyBytes))
	fmt.Printf("    Public key size: %d bytes\n", len(publicKeyBytes))

	// Deserialize keys
	restoredPrivateKey, err := topayz512.PrivateKeyFromBytes(privateKeyBytes)
	if err != nil {
		log.Printf("Failed to restore private key: %v", err)
	} else {
		fmt.Printf("    Private key restored: %v\n",
			topayz512.PrivateKeyEqual(privateKey, restoredPrivateKey))
	}

	restoredPublicKey, err := topayz512.PublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		log.Printf("Failed to restore public key: %v", err)
	} else {
		fmt.Printf("    Public key restored: %v\n",
			topayz512.PublicKeyEqual(publicKey, restoredPublicKey))
	}

	// Hex encoding/decoding
	privateKeyHex := privateKey.String()
	publicKeyHex := publicKey.String()

	fmt.Printf("    Private key hex: %s\n", privateKeyHex[:32]+"...")
	fmt.Printf("    Public key hex: %s\n", publicKeyHex[:32]+"...")

	restoredPrivateKeyFromHex, err := topayz512.PrivateKeyFromHex(privateKeyHex)
	if err != nil {
		log.Printf("Failed to restore private key from hex: %v", err)
	} else {
		fmt.Printf("    Private key from hex: %v\n",
			topayz512.PrivateKeyEqual(privateKey, restoredPrivateKeyFromHex))
	}

	restoredPublicKeyFromHex, err := topayz512.PublicKeyFromHex(publicKeyHex)
	if err != nil {
		log.Printf("Failed to restore public key from hex: %v", err)
	} else {
		fmt.Printf("    Public key from hex: %v\n",
			topayz512.PublicKeyEqual(publicKey, restoredPublicKeyFromHex))
	}

	fmt.Println()
	fmt.Println("=== Key Pair Operations Complete ===")
}
