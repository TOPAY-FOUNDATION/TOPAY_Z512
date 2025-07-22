package topayz512

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"sync"
	"time"
)

// KEM (Key Encapsulation Mechanism) operations for TOPAY-Z512 with optimizations

// KEMResult represents the result of key encapsulation
type KEMResult struct {
	Ciphertext    []byte    `json:"ciphertext"`
	SharedSecret  []byte    `json:"shared_secret"`
	Timestamp     time.Time `json:"timestamp"`
	KeySize       uint32    `json:"key_size"`
	SecurityLevel uint32    `json:"security_level"`
}

// KEMDecryptResult represents the result of key decapsulation
type KEMDecryptResult struct {
	SharedSecret []byte    `json:"shared_secret"`
	IsValid      bool      `json:"is_valid"`
	Timestamp    time.Time `json:"timestamp"`
	KeySize      uint32    `json:"key_size"`
}

// Encapsulate performs key encapsulation with the public key using optimizations
func Encapsulate(publicKey *PublicKey) (*KEMResult, error) {
	if publicKey == nil {
		return nil, errors.New("public key cannot be nil")
	}

	// Validate public key
	if !IsValidPublicKey(*publicKey) {
		return nil, errors.New("invalid public key")
	}

	// Generate shared secret using pooled buffer
	sharedSecret := GetBuffer(SharedSecretSize)
	if _, err := rand.Read(sharedSecret); err != nil {
		PutBuffer(sharedSecret)
		return nil, err
	}

	// Create ciphertext using optimized encryption
	ciphertext, err := encryptSharedSecret(sharedSecret, publicKey)
	if err != nil {
		PutBuffer(sharedSecret)
		return nil, err
	}

	return &KEMResult{
		Ciphertext:    ciphertext,
		SharedSecret:  sharedSecret,
		Timestamp:     time.Now(),
		KeySize:       uint32(len(sharedSecret)),
		SecurityLevel: SecurityLevel,
	}, nil
}

// Decapsulate performs key decapsulation with the private key using optimizations
func Decapsulate(ciphertext []byte, privateKey *PrivateKey) (*KEMDecryptResult, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}

	if len(ciphertext) == 0 {
		return nil, errors.New("ciphertext cannot be empty")
	}

	// Validate private key
	if !IsValidPrivateKey(*privateKey) {
		return nil, errors.New("invalid private key")
	}

	// Decrypt shared secret using optimized decryption
	sharedSecret, err := decryptSharedSecret(ciphertext, privateKey)
	if err != nil {
		return &KEMDecryptResult{
			SharedSecret: nil,
			IsValid:      false,
			Timestamp:    time.Now(),
			KeySize:      0,
		}, err
	}

	return &KEMDecryptResult{
		SharedSecret: sharedSecret,
		IsValid:      true,
		Timestamp:    time.Now(),
		KeySize:      uint32(len(sharedSecret)),
	}, nil
}

// encryptSharedSecret encrypts the shared secret with the public key using optimizations
func encryptSharedSecret(sharedSecret []byte, publicKey *PublicKey) ([]byte, error) {
	// Use pooled buffer for ciphertext
	ciphertext := GetBuffer(CiphertextSize)

	// Combine shared secret with public key data using SIMD operations
	keyData := publicKey.Bytes()

	// Use optimized XOR operation with SIMD when available
	if simdCaps.SSE2 && len(sharedSecret) >= 16 && len(keyData) >= len(sharedSecret) {
		VectorizedXOR(ciphertext[:len(sharedSecret)], sharedSecret, keyData[:len(sharedSecret)])
	} else {
		// Fallback to scalar XOR
		for i := 0; i < len(sharedSecret) && i < len(keyData); i++ {
			ciphertext[i] = sharedSecret[i] ^ keyData[i]
		}
	}

	// Add additional encryption layers
	hash := ComputeHash(ciphertext[:len(sharedSecret)])
	copy(ciphertext[len(sharedSecret):], hash[:])

	// Pad to full ciphertext size with secure random data
	if len(ciphertext) > len(sharedSecret)+HashSize {
		padding := ciphertext[len(sharedSecret)+HashSize:]
		if _, err := rand.Read(padding); err != nil {
			PutBuffer(ciphertext)
			return nil, err
		}
	}

	return ciphertext, nil
}

// decryptSharedSecret decrypts the shared secret with the private key using optimizations
func decryptSharedSecret(ciphertext []byte, privateKey *PrivateKey) ([]byte, error) {
	if len(ciphertext) < SharedSecretSize+HashSize {
		return nil, errors.New("ciphertext too short")
	}

	// Extract encrypted shared secret
	encryptedSecret := ciphertext[:SharedSecretSize]
	expectedHash := ciphertext[SharedSecretSize : SharedSecretSize+HashSize]

	// Use pooled buffer for decrypted secret
	sharedSecret := GetBuffer(SharedSecretSize)

	// Decrypt using private key with SIMD optimization
	keyData := privateKey.Bytes()

	if simdCaps.SSE2 && len(encryptedSecret) >= 16 && len(keyData) >= len(encryptedSecret) {
		VectorizedXOR(sharedSecret, encryptedSecret, keyData[:len(encryptedSecret)])
	} else {
		// Fallback to scalar XOR
		for i := 0; i < len(encryptedSecret) && i < len(keyData); i++ {
			sharedSecret[i] = encryptedSecret[i] ^ keyData[i]
		}
	}

	// Verify integrity using constant-time comparison
	computedHash := ComputeHash(encryptedSecret)
	if !ConstantTimeEqual(expectedHash, computedHash[:HashSize]) {
		PutBuffer(sharedSecret)
		return nil, errors.New("integrity check failed")
	}

	return sharedSecret, nil
}

// KEMKeyGen generates a new KEM key pair
func KEMKeyGen() (KEMPublicKey, KEMSecretKey, error) {
	// Generate random secret key
	secretBytes, err := SecureRandom(KEMSecretKeySize)
	if err != nil {
		return KEMPublicKey{}, KEMSecretKey{}, err
	}

	var secretKey KEMSecretKey
	copy(secretKey[:], secretBytes)

	// Derive public key from secret key
	publicKey := deriveKEMPublicKey(secretKey)

	return publicKey, secretKey, nil
}

// deriveKEMPublicKey derives a KEM public key from a secret key
func deriveKEMPublicKey(secretKey KEMSecretKey) KEMPublicKey {
	// Use SHA-256 based key derivation for KEM public key
	hasher := sha256.New()
	hasher.Write(secretKey[:])
	hasher.Write([]byte("TOPAY-Z512-KEM-PUBLIC-KEY"))

	firstHash := hasher.Sum(nil)

	// Second round for additional security
	hasher.Reset()
	hasher.Write(firstHash)
	hasher.Write(secretKey[:])
	hasher.Write([]byte("KEM-DERIVATION-ROUND-2"))
	secondHash := hasher.Sum(nil)

	var publicKey KEMPublicKey
	copy(publicKey[:], secondHash)

	return publicKey
}

// KEMEncapsulate encapsulates a shared secret using the public key
func KEMEncapsulate(publicKey KEMPublicKey) (Ciphertext, SharedSecret, error) {
	// Generate random ephemeral key
	ephemeralBytes, err := SecureRandom(32)
	if err != nil {
		return Ciphertext{}, SharedSecret{}, err
	}

	// Derive shared secret from ephemeral key and public key
	hasher := sha256.New()
	hasher.Write(ephemeralBytes)
	hasher.Write(publicKey[:])
	hasher.Write([]byte("TOPAY-Z512-KEM-SHARED-SECRET"))

	sharedSecretHash := hasher.Sum(nil)

	var sharedSecret SharedSecret
	copy(sharedSecret[:], sharedSecretHash)

	// Create ciphertext by encrypting ephemeral key with public key
	ciphertext := createCiphertext(ephemeralBytes, publicKey)

	return ciphertext, sharedSecret, nil
}

// KEMDecapsulate decapsulates the shared secret using the secret key
func KEMDecapsulate(secretKey KEMSecretKey, ciphertext Ciphertext) (SharedSecret, error) {
	// Derive public key from secret key for verification
	publicKey := deriveKEMPublicKey(secretKey)

	// Decrypt ephemeral key from ciphertext
	ephemeralBytes, err := decryptCiphertext(ciphertext, secretKey)
	if err != nil {
		return SharedSecret{}, ErrDecapsulationFailed
	}

	// Derive shared secret from ephemeral key and public key
	hasher := sha256.New()
	hasher.Write(ephemeralBytes)
	hasher.Write(publicKey[:])
	hasher.Write([]byte("TOPAY-Z512-KEM-SHARED-SECRET"))

	sharedSecretHash := hasher.Sum(nil)

	var sharedSecret SharedSecret
	copy(sharedSecret[:], sharedSecretHash)

	return sharedSecret, nil
}

// createCiphertext creates a ciphertext by encrypting ephemeral key with public key
func createCiphertext(ephemeralKey []byte, publicKey KEMPublicKey) Ciphertext {
	// Simple XOR-based encryption for demonstration
	// In production, use proper lattice-based encryption
	hasher := sha256.New()
	hasher.Write(publicKey[:])
	hasher.Write([]byte("TOPAY-Z512-KEM-ENCRYPTION-KEY"))

	encryptionKey := hasher.Sum(nil)

	var ciphertext Ciphertext

	// XOR ephemeral key with derived encryption key
	for i := 0; i < len(ephemeralKey) && i < CiphertextSize; i++ {
		ciphertext[i] = ephemeralKey[i] ^ encryptionKey[i%len(encryptionKey)]
	}

	// Fill remaining bytes with hash of the encrypted portion
	if len(ephemeralKey) < CiphertextSize {
		hasher.Reset()
		hasher.Write(ciphertext[:len(ephemeralKey)])
		hasher.Write(publicKey[:])
		fillHash := hasher.Sum(nil)

		for i := len(ephemeralKey); i < CiphertextSize; i++ {
			ciphertext[i] = fillHash[i%len(fillHash)]
		}
	}

	return ciphertext
}

// decryptCiphertext decrypts the ciphertext to recover ephemeral key
func decryptCiphertext(ciphertext Ciphertext, secretKey KEMSecretKey) ([]byte, error) {
	// Derive public key from secret key
	publicKey := deriveKEMPublicKey(secretKey)

	// Derive encryption key
	hasher := sha256.New()
	hasher.Write(publicKey[:])
	hasher.Write([]byte("TOPAY-Z512-KEM-ENCRYPTION-KEY"))

	encryptionKey := hasher.Sum(nil)

	// Decrypt ephemeral key (assuming 32 bytes)
	ephemeralKey := make([]byte, 32)
	for i := 0; i < len(ephemeralKey); i++ {
		ephemeralKey[i] = ciphertext[i] ^ encryptionKey[i%len(encryptionKey)]
	}

	return ephemeralKey, nil
}

// Batch KEM operations

// BatchKEMResult represents the result of batch KEM operations
type BatchKEMResult struct {
	Index        int
	PublicKey    KEMPublicKey
	SecretKey    KEMSecretKey
	Ciphertext   Ciphertext
	SharedSecret SharedSecret
	Error        error
}

// BatchKEMKeyGen generates multiple KEM key pairs in parallel
func BatchKEMKeyGen(count int) ([]KEMPublicKey, []KEMSecretKey, error) {
	if count <= 0 {
		return nil, nil, ErrInvalidFragmentCount
	}

	publicKeys := make([]KEMPublicKey, count)
	secretKeys := make([]KEMSecretKey, count)

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > count {
		numWorkers = count
	}

	// Channel for work distribution
	workChan := make(chan int, count)
	resultChan := make(chan BatchKEMResult, count)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				publicKey, secretKey, err := KEMKeyGen()
				resultChan <- BatchKEMResult{
					Index:     index,
					PublicKey: publicKey,
					SecretKey: secretKey,
					Error:     err,
				}
			}
		}()
	}

	// Send work
	go func() {
		for i := 0; i < count; i++ {
			workChan <- i
		}
		close(workChan)
	}()

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		if result.Error != nil {
			return nil, nil, result.Error
		}
		publicKeys[result.Index] = result.PublicKey
		secretKeys[result.Index] = result.SecretKey
	}

	return publicKeys, secretKeys, nil
}

// BatchKEMEncapsulate performs multiple encapsulations in parallel
func BatchKEMEncapsulate(publicKeys []KEMPublicKey) ([]Ciphertext, []SharedSecret, error) {
	if len(publicKeys) == 0 {
		return nil, nil, ErrEmptyData
	}

	ciphertexts := make([]Ciphertext, len(publicKeys))
	sharedSecrets := make([]SharedSecret, len(publicKeys))

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > len(publicKeys) {
		numWorkers = len(publicKeys)
	}

	// Channel for work distribution
	workChan := make(chan int, len(publicKeys))
	resultChan := make(chan BatchKEMResult, len(publicKeys))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				ciphertext, sharedSecret, err := KEMEncapsulate(publicKeys[index])
				resultChan <- BatchKEMResult{
					Index:        index,
					Ciphertext:   ciphertext,
					SharedSecret: sharedSecret,
					Error:        err,
				}
			}
		}()
	}

	// Send work
	go func() {
		for i := range publicKeys {
			workChan <- i
		}
		close(workChan)
	}()

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		if result.Error != nil {
			return nil, nil, result.Error
		}
		ciphertexts[result.Index] = result.Ciphertext
		sharedSecrets[result.Index] = result.SharedSecret
	}

	return ciphertexts, sharedSecrets, nil
}

// BatchKEMDecapsulate performs multiple decapsulations in parallel
func BatchKEMDecapsulate(secretKeys []KEMSecretKey, ciphertexts []Ciphertext) ([]SharedSecret, error) {
	if len(secretKeys) != len(ciphertexts) {
		return nil, ErrInvalidFragmentCount
	}

	if len(secretKeys) == 0 {
		return nil, ErrEmptyData
	}

	sharedSecrets := make([]SharedSecret, len(secretKeys))

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > len(secretKeys) {
		numWorkers = len(secretKeys)
	}

	// Channel for work distribution
	workChan := make(chan int, len(secretKeys))
	resultChan := make(chan BatchKEMResult, len(secretKeys))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				sharedSecret, err := KEMDecapsulate(secretKeys[index], ciphertexts[index])
				resultChan <- BatchKEMResult{
					Index:        index,
					SharedSecret: sharedSecret,
					Error:        err,
				}
			}
		}()
	}

	// Send work
	go func() {
		for i := range secretKeys {
			workChan <- i
		}
		close(workChan)
	}()

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		if result.Error != nil {
			return nil, result.Error
		}
		sharedSecrets[result.Index] = result.SharedSecret
	}

	return sharedSecrets, nil
}

// KEM validation and utilities

// VerifyKEMKeyPair verifies that a KEM public and secret key form a valid pair
func VerifyKEMKeyPair(publicKey KEMPublicKey, secretKey KEMSecretKey) bool {
	derivedPublic := deriveKEMPublicKey(secretKey)
	return ConstantTimeEqual(publicKey[:], derivedPublic[:])
}

// IsValidKEMPublicKey checks if a KEM public key is valid
func IsValidKEMPublicKey(publicKey KEMPublicKey) bool {
	var zero KEMPublicKey
	return !ConstantTimeEqual(publicKey[:], zero[:])
}

// IsValidKEMSecretKey checks if a KEM secret key is valid
func IsValidKEMSecretKey(secretKey KEMSecretKey) bool {
	var zero KEMSecretKey
	return !ConstantTimeEqual(secretKey[:], zero[:])
}

// IsValidCiphertext checks if a ciphertext is valid
func IsValidCiphertext(ciphertext Ciphertext) bool {
	var zero Ciphertext
	return !ConstantTimeEqual(ciphertext[:], zero[:])
}

// IsValidSharedSecret checks if a shared secret is valid
func IsValidSharedSecret(sharedSecret SharedSecret) bool {
	var zero SharedSecret
	return !ConstantTimeEqual(sharedSecret[:], zero[:])
}

// KEMEqual compares KEM components for equality

// KEMPublicKeyEqual compares two KEM public keys for equality
func KEMPublicKeyEqual(pk1, pk2 KEMPublicKey) bool {
	return ConstantTimeEqual(pk1[:], pk2[:])
}

// KEMSecretKeyEqual compares two KEM secret keys for equality
func KEMSecretKeyEqual(sk1, sk2 KEMSecretKey) bool {
	return ConstantTimeEqual(sk1[:], sk2[:])
}

// CiphertextEqual compares two ciphertexts for equality
func CiphertextEqual(ct1, ct2 Ciphertext) bool {
	return ConstantTimeEqual(ct1[:], ct2[:])
}

// SharedSecretEqual compares two shared secrets for equality
func SharedSecretEqual(ss1, ss2 SharedSecret) bool {
	return ConstantTimeEqual(ss1[:], ss2[:])
}

// Secure erasure functions

// SecureEraseKEMSecretKey securely erases a KEM secret key from memory
func SecureEraseKEMSecretKey(secretKey *KEMSecretKey) {
	SecureZero(secretKey[:])
}

// SecureEraseSharedSecret securely erases a shared secret from memory
func SecureEraseSharedSecret(sharedSecret *SharedSecret) {
	SecureZero(sharedSecret[:])
}

// SecureEraseKEMKeyPair securely erases a KEM key pair from memory
func SecureEraseKEMKeyPair(keyPair *KEMKeyPair) {
	SecureZero(keyPair.Public[:])
	SecureZero(keyPair.Secret[:])
}

// KEM performance benchmarking

// KEMBenchmark represents KEM performance metrics
type KEMBenchmark struct {
	KeyGenPerSec      float64
	EncapsulatePerSec float64
	DecapsulatePerSec float64
	BatchSpeedupRatio float64
	AvgLatencyMs      float64
}

// BenchmarkKEM measures KEM performance
func BenchmarkKEM(iterations int) KEMBenchmark {
	// Key generation benchmark
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, _, _ = KEMKeyGen()
	}
	keyGenDuration := time.Since(start)

	// Generate test key pair
	publicKey, secretKey, _ := KEMKeyGen()

	// Encapsulation benchmark
	start = time.Now()
	var ciphertext Ciphertext
	for i := 0; i < iterations; i++ {
		ciphertext, _, _ = KEMEncapsulate(publicKey)
	}
	encapDuration := time.Since(start)

	// Decapsulation benchmark
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, _ = KEMDecapsulate(secretKey, ciphertext)
	}
	decapDuration := time.Since(start)

	// Batch benchmark
	start = time.Now()
	_, _, _ = BatchKEMKeyGen(iterations)
	batchDuration := time.Since(start)

	keyGenPerSec := float64(iterations) / keyGenDuration.Seconds()
	encapsulatePerSec := float64(iterations) / encapDuration.Seconds()
	decapsulatePerSec := float64(iterations) / decapDuration.Seconds()
	batchSpeedupRatio := keyGenDuration.Seconds() / batchDuration.Seconds()
	avgLatencyMs := keyGenDuration.Seconds() * 1000 / float64(iterations)

	return KEMBenchmark{
		KeyGenPerSec:      keyGenPerSec,
		EncapsulatePerSec: encapsulatePerSec,
		DecapsulatePerSec: decapsulatePerSec,
		BatchSpeedupRatio: batchSpeedupRatio,
		AvgLatencyMs:      avgLatencyMs,
	}
}

// Advanced KEM operations

// KEMWithContext performs KEM operations with additional context data
func KEMWithContext(publicKey KEMPublicKey, context []byte) (Ciphertext, SharedSecret, error) {
	// Generate random ephemeral key
	ephemeralBytes, err := SecureRandom(32)
	if err != nil {
		return Ciphertext{}, SharedSecret{}, err
	}

	// Derive shared secret with context
	hasher := sha256.New()
	hasher.Write(ephemeralBytes)
	hasher.Write(publicKey[:])
	hasher.Write(context)
	hasher.Write([]byte("TOPAY-Z512-KEM-CONTEXT-SECRET"))

	sharedSecretHash := hasher.Sum(nil)

	var sharedSecret SharedSecret
	copy(sharedSecret[:], sharedSecretHash)

	// Create ciphertext
	ciphertext := createCiphertext(ephemeralBytes, publicKey)

	return ciphertext, sharedSecret, nil
}

// KEMDecapsulateWithContext decapsulates with additional context data
func KEMDecapsulateWithContext(secretKey KEMSecretKey, ciphertext Ciphertext, context []byte) (SharedSecret, error) {
	// Derive public key from secret key
	publicKey := deriveKEMPublicKey(secretKey)

	// Decrypt ephemeral key from ciphertext
	ephemeralBytes, err := decryptCiphertext(ciphertext, secretKey)
	if err != nil {
		return SharedSecret{}, ErrDecapsulationFailed
	}

	// Derive shared secret with context
	hasher := sha256.New()
	hasher.Write(ephemeralBytes)
	hasher.Write(publicKey[:])
	hasher.Write(context)
	hasher.Write([]byte("TOPAY-Z512-KEM-CONTEXT-SECRET"))

	sharedSecretHash := hasher.Sum(nil)

	var sharedSecret SharedSecret
	copy(sharedSecret[:], sharedSecretHash)

	return sharedSecret, nil
}
