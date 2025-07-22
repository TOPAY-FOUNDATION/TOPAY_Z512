package main

import (
	"fmt"
	"log"
	"time"

	"github.com/topayfoundation/topayz512"
)

func main() {
	fmt.Println("=== TOPAY-Z512 Go Implementation - KEM Operations ===")
	fmt.Println()

	// Basic KEM operations
	fmt.Println("1. Basic KEM Key Generation:")
	publicKey, secretKey, err := topayz512.KEMKeyGen()
	if err != nil {
		log.Fatalf("Failed to generate KEM key pair: %v", err)
	}

	fmt.Printf("   Public Key: %s\n", publicKey.String())
	fmt.Printf("   Secret Key: %s\n", secretKey.String())
	fmt.Printf("   Key pair valid: %v\n", topayz512.VerifyKEMKeyPair(publicKey, secretKey))
	fmt.Println()

	// KEM encapsulation
	fmt.Println("2. KEM Encapsulation:")
	ciphertext, sharedSecret1, err := topayz512.KEMEncapsulate(publicKey)
	if err != nil {
		log.Fatalf("Failed to encapsulate: %v", err)
	}

	fmt.Printf("   Ciphertext: %s\n", ciphertext.String())
	fmt.Printf("   Shared Secret: %s\n", sharedSecret1.String())
	fmt.Printf("   Ciphertext valid: %v\n", topayz512.IsValidCiphertext(ciphertext))
	fmt.Printf("   Shared secret valid: %v\n", topayz512.IsValidSharedSecret(sharedSecret1))
	fmt.Println()

	// KEM decapsulation
	fmt.Println("3. KEM Decapsulation:")
	sharedSecret2, err := topayz512.KEMDecapsulate(secretKey, ciphertext)
	if err != nil {
		log.Fatalf("Failed to decapsulate: %v", err)
	}

	fmt.Printf("   Decapsulated Secret: %s\n", sharedSecret2.String())
	fmt.Printf("   Secrets match: %v\n", topayz512.SharedSecretEqual(sharedSecret1, sharedSecret2))
	fmt.Println()

	// Multiple encapsulations with same key
	fmt.Println("4. Multiple Encapsulations:")
	var ss2 topayz512.SharedSecret
	var ct topayz512.Ciphertext
	var ss topayz512.SharedSecret
	for i := 0; i < 3; i++ {
		ct, ss, err = topayz512.KEMEncapsulate(publicKey)
		if err != nil {
			log.Printf("Encapsulation %d failed: %v", i+1, err)
			continue
		}

		// Decapsulate to verify
		ss2, err = topayz512.KEMDecapsulate(secretKey, ct)
		if err != nil {
			log.Printf("Decapsulation %d failed: %v", i+1, err)
			continue
		}

		fmt.Printf("   Encapsulation %d: Secrets match = %v\n", i+1,
			topayz512.SharedSecretEqual(ss, ss2))
	}
	fmt.Println()

	// Batch KEM operations
	fmt.Println("5. Batch KEM Operations:")
	batchSize := 10

	// Batch key generation
	start := time.Now()
	publicKeys, secretKeys, err := topayz512.BatchKEMKeyGen(batchSize)
	keyGenDuration := time.Since(start)
	if err != nil {
		log.Fatalf("Failed batch key generation: %v", err)
	}

	fmt.Printf("   Generated %d key pairs in %v\n", len(publicKeys), keyGenDuration)
	fmt.Printf("   Key generation rate: %.2f keys/sec\n",
		float64(batchSize)/keyGenDuration.Seconds())

	// Verify all generated key pairs
	validKeyPairs := 0
	for i := 0; i < batchSize; i++ {
		if topayz512.VerifyKEMKeyPair(publicKeys[i], secretKeys[i]) {
			validKeyPairs++
		}
	}
	fmt.Printf("   Valid key pairs: %d/%d\n", validKeyPairs, batchSize)
	fmt.Println()

	// Batch encapsulation
	fmt.Println("6. Batch Encapsulation:")
	start = time.Now()
	ciphertexts, sharedSecrets1, err := topayz512.BatchKEMEncapsulate(publicKeys)
	encapDuration := time.Since(start)
	if err != nil {
		log.Fatalf("Failed batch encapsulation: %v", err)
	}

	fmt.Printf("   Encapsulated %d ciphertexts in %v\n", len(ciphertexts), encapDuration)
	fmt.Printf("   Encapsulation rate: %.2f ops/sec\n",
		float64(batchSize)/encapDuration.Seconds())
	fmt.Println()

	// Batch decapsulation
	fmt.Println("7. Batch Decapsulation:")
	start = time.Now()
	sharedSecrets2, err := topayz512.BatchKEMDecapsulate(secretKeys, ciphertexts)
	decapDuration := time.Since(start)
	if err != nil {
		log.Fatalf("Failed batch decapsulation: %v", err)
	}

	fmt.Printf("   Decapsulated %d secrets in %v\n", len(sharedSecrets2), decapDuration)
	fmt.Printf("   Decapsulation rate: %.2f ops/sec\n",
		float64(batchSize)/decapDuration.Seconds())

	// Verify all shared secrets match
	matchingSecrets := 0
	for i := 0; i < batchSize; i++ {
		if topayz512.SharedSecretEqual(sharedSecrets1[i], sharedSecrets2[i]) {
			matchingSecrets++
		}
	}
	fmt.Printf("   Matching secrets: %d/%d\n", matchingSecrets, batchSize)
	fmt.Println()

	// KEM with context
	fmt.Println("8. KEM with Context:")
	context := []byte("secure_communication_session_2024")

	contextCiphertext, contextSharedSecret1, err := topayz512.KEMWithContext(publicKey, context)
	if err != nil {
		log.Fatalf("Failed KEM with context: %v", err)
	}

	contextSharedSecret2, err := topayz512.KEMDecapsulateWithContext(secretKey, contextCiphertext, context)
	if err != nil {
		log.Fatalf("Failed KEM decapsulation with context: %v", err)
	}

	fmt.Printf("   Context: %s\n", string(context))
	fmt.Printf("   Context ciphertext: %s\n", contextCiphertext.String())
	fmt.Printf("   Context secrets match: %v\n",
		topayz512.SharedSecretEqual(contextSharedSecret1, contextSharedSecret2))

	// Verify context affects the result
	differentContext := []byte("different_context")
	_, differentSecret, _ := topayz512.KEMWithContext(publicKey, differentContext)
	fmt.Printf("   Different context produces different secret: %v\n",
		!topayz512.SharedSecretEqual(contextSharedSecret1, differentSecret))
	fmt.Println()

	// Key validation
	fmt.Println("9. KEM Key Validation:")
	fmt.Printf("   Public key valid: %v\n", topayz512.IsValidKEMPublicKey(publicKey))
	fmt.Printf("   Secret key valid: %v\n", topayz512.IsValidKEMSecretKey(secretKey))

	// Test with invalid keys
	invalidPublicKey := topayz512.KEMPublicKey{}
	invalidSecretKey := topayz512.KEMSecretKey{}
	fmt.Printf("   Invalid public key: %v\n", topayz512.IsValidKEMPublicKey(invalidPublicKey))
	fmt.Printf("   Invalid secret key: %v\n", topayz512.IsValidKEMSecretKey(invalidSecretKey))
	fmt.Println()

	// Ciphertext and shared secret validation
	fmt.Println("10. Ciphertext and Shared Secret Validation:")
	fmt.Printf("    Ciphertext valid: %v\n", topayz512.IsValidCiphertext(ciphertext))
	fmt.Printf("    Shared secret valid: %v\n", topayz512.IsValidSharedSecret(sharedSecret1))

	// Test with invalid values
	invalidCiphertext := topayz512.Ciphertext{}
	invalidSharedSecret := topayz512.SharedSecret{}
	fmt.Printf("    Invalid ciphertext: %v\n", topayz512.IsValidCiphertext(invalidCiphertext))
	fmt.Printf("    Invalid shared secret: %v\n", topayz512.IsValidSharedSecret(invalidSharedSecret))
	fmt.Println()

	// Equality comparisons
	fmt.Println("11. Equality Comparisons:")

	// Generate another key pair for comparison
	publicKey2, secretKey2, _ := topayz512.KEMKeyGen()
	ciphertext2, sharedSecret3, _ := topayz512.KEMEncapsulate(publicKey2)

	fmt.Printf("    Same public keys: %v\n", topayz512.KEMPublicKeyEqual(publicKey, publicKey))
	fmt.Printf("    Different public keys: %v\n", topayz512.KEMPublicKeyEqual(publicKey, publicKey2))
	fmt.Printf("    Same secret keys: %v\n", topayz512.KEMSecretKeyEqual(secretKey, secretKey))
	fmt.Printf("    Different secret keys: %v\n", topayz512.KEMSecretKeyEqual(secretKey, secretKey2))
	fmt.Printf("    Same ciphertexts: %v\n", topayz512.CiphertextEqual(ciphertext, ciphertext))
	fmt.Printf("    Different ciphertexts: %v\n", topayz512.CiphertextEqual(ciphertext, ciphertext2))
	fmt.Printf("    Same shared secrets: %v\n", topayz512.SharedSecretEqual(sharedSecret1, sharedSecret1))
	fmt.Printf("    Different shared secrets: %v\n", topayz512.SharedSecretEqual(sharedSecret1, sharedSecret3))
	fmt.Println()

	// Serialization and deserialization
	fmt.Println("12. KEM Serialization:")

	// Serialize KEM components
	publicKeyBytes := publicKey.Bytes()
	secretKeyBytes := secretKey.Bytes()
	ciphertextBytes := ciphertext.Bytes()
	sharedSecretBytes := sharedSecret1.Bytes()

	fmt.Printf("    Public key size: %d bytes\n", len(publicKeyBytes))
	fmt.Printf("    Secret key size: %d bytes\n", len(secretKeyBytes))
	fmt.Printf("    Ciphertext size: %d bytes\n", len(ciphertextBytes))
	fmt.Printf("    Shared secret size: %d bytes\n", len(sharedSecretBytes))

	// Deserialize KEM components
	restoredPublicKey, err := topayz512.KEMPublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		log.Printf("Failed to restore public key: %v", err)
	} else {
		fmt.Printf("    Public key restored: %v\n",
			topayz512.KEMPublicKeyEqual(publicKey, restoredPublicKey))
	}

	restoredSecretKey, err := topayz512.KEMSecretKeyFromBytes(secretKeyBytes)
	if err != nil {
		log.Printf("Failed to restore secret key: %v", err)
	} else {
		fmt.Printf("    Secret key restored: %v\n",
			topayz512.KEMSecretKeyEqual(secretKey, restoredSecretKey))
	}

	restoredCiphertext, err := topayz512.CiphertextFromBytes(ciphertextBytes)
	if err != nil {
		log.Printf("Failed to restore ciphertext: %v", err)
	} else {
		fmt.Printf("    Ciphertext restored: %v\n",
			topayz512.CiphertextEqual(ciphertext, restoredCiphertext))
	}

	restoredSharedSecret, err := topayz512.SharedSecretFromBytes(sharedSecretBytes)
	if err != nil {
		log.Printf("Failed to restore shared secret: %v", err)
	} else {
		fmt.Printf("    Shared secret restored: %v\n",
			topayz512.SharedSecretEqual(sharedSecret1, restoredSharedSecret))
	}
	fmt.Println()

	// Hex encoding/decoding
	fmt.Println("13. Hex Encoding/Decoding:")
	publicKeyHex := publicKey.String()
	secretKeyHex := secretKey.String()
	ciphertextHex := ciphertext.String()
	sharedSecretHex := sharedSecret1.String()

	fmt.Printf("    Public key hex: %s...\n", publicKeyHex[:32])
	fmt.Printf("    Secret key hex: %s...\n", secretKeyHex[:32])
	fmt.Printf("    Ciphertext hex: %s...\n", ciphertextHex[:32])
	fmt.Printf("    Shared secret hex: %s...\n", sharedSecretHex[:32])

	// Restore from hex
	restoredPublicKeyFromHex, err := topayz512.KEMPublicKeyFromHex(publicKeyHex)
	if err != nil {
		log.Printf("Failed to restore public key from hex: %v", err)
	} else {
		fmt.Printf("    Public key from hex: %v\n",
			topayz512.KEMPublicKeyEqual(publicKey, restoredPublicKeyFromHex))
	}

	restoredSecretKeyFromHex, err := topayz512.KEMSecretKeyFromHex(secretKeyHex)
	if err != nil {
		log.Printf("Failed to restore secret key from hex: %v", err)
	} else {
		fmt.Printf("    Secret key from hex: %v\n",
			topayz512.KEMSecretKeyEqual(secretKey, restoredSecretKeyFromHex))
	}
	fmt.Println()

	// Secure erasure
	fmt.Println("14. Secure Erasure:")
	tempPublicKey, tempSecretKey, _ := topayz512.KEMKeyGen()
	tempCiphertext, tempSharedSecret, _ := topayz512.KEMEncapsulate(tempPublicKey)

	fmt.Printf("    Before erasure - Secret key valid: %v\n",
		topayz512.IsValidKEMSecretKey(tempSecretKey))
	fmt.Printf("    Before erasure - Shared secret valid: %v\n",
		topayz512.IsValidSharedSecret(tempSharedSecret))

	topayz512.SecureEraseKEMSecretKey(&tempSecretKey)
	topayz512.SecureEraseSharedSecret(&tempSharedSecret)

	fmt.Printf("    After erasure - Secret key valid: %v\n",
		topayz512.IsValidKEMSecretKey(tempSecretKey))
	fmt.Printf("    After erasure - Shared secret valid: %v\n",
		topayz512.IsValidSharedSecret(tempSharedSecret))
	fmt.Printf("    Public key still valid: %v\n",
		topayz512.IsValidKEMPublicKey(tempPublicKey))
	fmt.Printf("    Ciphertext still valid: %v\n",
		topayz512.IsValidCiphertext(tempCiphertext))
	fmt.Println()

	// Performance benchmark
	fmt.Println("15. KEM Performance Benchmark:")
	iterations := 1000
	benchmark := topayz512.BenchmarkKEM(iterations)

	fmt.Printf("    Key generation: %.2f ops/sec (%.2f ms avg)\n",
		benchmark.KeyGenPerSec, benchmark.AvgLatencyMs)
	fmt.Printf("    Encapsulation: %.2f ops/sec\n", benchmark.EncapsulatePerSec)
	fmt.Printf("    Decapsulation: %.2f ops/sec\n", benchmark.DecapsulatePerSec)
	fmt.Printf("    Batch speedup: %.2fx\n", benchmark.BatchSpeedupRatio)
	fmt.Println()

	// Error handling demonstration
	fmt.Println("16. Error Handling:")

	// Try to decapsulate with wrong secret key
	_, wrongSecretKey, _ := topayz512.KEMKeyGen()
	wrongSharedSecret, _ := topayz512.KEMDecapsulate(wrongSecretKey, ciphertext)

	fmt.Printf("    Decapsulation with wrong key produces different secret: %v\n",
		!topayz512.SharedSecretEqual(sharedSecret1, wrongSharedSecret))

	// Try operations with invalid data
	invalidCt := topayz512.Ciphertext{}
	_, err = topayz512.KEMDecapsulate(secretKey, invalidCt)
	fmt.Printf("    Decapsulation with invalid ciphertext fails: %v\n", err != nil)

	// Try batch operations with mismatched sizes
	shortSecretKeys := secretKeys[:5]
	_, err = topayz512.BatchKEMDecapsulate(shortSecretKeys, ciphertexts)
	fmt.Printf("    Batch decapsulation with mismatched sizes fails: %v\n", err != nil)

	fmt.Println()
	fmt.Println("=== KEM Operations Complete ===")
}
