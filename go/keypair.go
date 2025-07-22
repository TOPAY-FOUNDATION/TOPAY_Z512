package topayz512

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"sync"
	"time"
)

// Key pair generation and management for TOPAY-Z512 with optimizations

// KeyPair represents a public/private key pair
type KeyPair struct {
	PublicKey  *PublicKey  `json:"public_key"`
	PrivateKey *PrivateKey `json:"private_key"`
	Timestamp  time.Time   `json:"timestamp"`
	KeyID      []byte      `json:"key_id"`
}

// GenerateKeyPair generates a new cryptographic key pair
func GenerateKeyPair() (PrivateKey, PublicKey, error) {
	// Generate random private key
	privateBytes, err := SecureRandom(PrivateKeySize)
	if err != nil {
		return PrivateKey{}, PublicKey{}, err
	}

	var privateKey PrivateKey
	copy(privateKey[:], privateBytes)

	// Derive public key from private key
	publicKey := DerivePublicKey(privateKey)

	return privateKey, publicKey, nil
}

// GenerateKeyPairAdvanced generates a new TOPAY-Z512 key pair with optimizations
func GenerateKeyPairAdvanced() (*KeyPair, error) {
	// Use pooled buffers for key generation
	privateKeyData := GetBuffer(PrivateKeySize)
	publicKeyData := GetBuffer(PublicKeySize)
	keyID := GetBuffer(16) // 128-bit key ID

	// Generate secure random private key
	if _, err := rand.Read(privateKeyData); err != nil {
		PutBuffer(privateKeyData)
		PutBuffer(publicKeyData)
		PutBuffer(keyID)
		return nil, err
	}

	// Derive public key from private key using optimized operations
	if err := derivePublicKeyAdvanced(publicKeyData, privateKeyData); err != nil {
		// Secure erase private key on error
		SecureZero(privateKeyData)
		PutBuffer(privateKeyData)
		PutBuffer(publicKeyData)
		PutBuffer(keyID)
		return nil, err
	}

	// Generate unique key ID
	if _, err := rand.Read(keyID); err != nil {
		SecureZero(privateKeyData)
		PutBuffer(privateKeyData)
		PutBuffer(publicKeyData)
		PutBuffer(keyID)
		return nil, err
	}

	timestamp := time.Now()

	// Create public key
	var pubKey PublicKey
	copy(pubKey[:], publicKeyData)

	// Create private key
	var privKey PrivateKey
	copy(privKey[:], privateKeyData)

	// Return buffers to pool
	PutBuffer(privateKeyData)
	PutBuffer(publicKeyData)
	PutBuffer(keyID)

	return &KeyPair{
		PublicKey:  &pubKey,
		PrivateKey: &privKey,
		Timestamp:  timestamp,
		KeyID:      make([]byte, 16), // Copy keyID data
	}, nil
}

// derivePublicKeyAdvanced derives the public key from private key using optimized operations
func derivePublicKeyAdvanced(publicKey, privateKey []byte) error {
	if len(publicKey) != PublicKeySize || len(privateKey) != PrivateKeySize {
		return errors.New("invalid key sizes")
	}

	// Use optimized hash-based key derivation
	hash := ComputeHash(privateKey)

	// Use SIMD operations for key derivation when available
	if simdCaps.SSE2 && len(publicKey) >= 16 {
		// Process in 16-byte chunks using vectorized operations
		for i := 0; i < len(publicKey); i += 16 {
			end := i + 16
			if end > len(publicKey) {
				end = len(publicKey)
			}

			// Use hash bytes cyclically for key derivation
			hashOffset := (i / 16) % (HashSize / 16)
			hashChunk := hash[hashOffset*16 : hashOffset*16+16]

			if end-i == 16 {
				VectorizedXOR(publicKey[i:end], privateKey[i:end], hashChunk)
			} else {
				// Fallback for partial chunks
				for j := i; j < end; j++ {
					publicKey[j] = privateKey[j] ^ hashChunk[j-i]
				}
			}
		}
	} else {
		// Fallback to scalar operations
		for i := 0; i < len(publicKey); i++ {
			hashByte := hash[i%HashSize]
			publicKey[i] = privateKey[i] ^ hashByte
		}
	}

	return nil
}

// DerivePublicKey derives a public key from a private key
func DerivePublicKey(privateKey PrivateKey) PublicKey {
	// Use SHA-256 based key derivation
	// In production, use proper elliptic curve or lattice-based derivation
	hasher := sha256.New()
	hasher.Write(privateKey[:])
	hasher.Write([]byte("TOPAY-Z512-PUBLIC-KEY-DERIVATION"))

	firstHash := hasher.Sum(nil)

	// Second round for additional security
	hasher.Reset()
	hasher.Write(firstHash)
	hasher.Write(privateKey[:])
	secondHash := hasher.Sum(nil)

	var publicKey PublicKey
	copy(publicKey[:], secondHash)

	return publicKey
}

// GenerateKeyPairFromSeed generates a deterministic key pair from a seed
func GenerateKeyPairFromSeed(seed []byte) (PrivateKey, PublicKey, error) {
	if len(seed) < 32 {
		return PrivateKey{}, PublicKey{}, ErrInvalidKeySize
	}

	// Derive private key from seed
	hasher := sha256.New()
	hasher.Write(seed)
	hasher.Write([]byte("TOPAY-Z512-PRIVATE-KEY-SEED"))

	privateHash := hasher.Sum(nil)

	// Extend to full private key size
	var privateKey PrivateKey
	for i := 0; i < PrivateKeySize; i++ {
		privateKey[i] = privateHash[i%len(privateHash)]
	}

	// Derive public key
	publicKey := DerivePublicKey(privateKey)

	return privateKey, publicKey, nil
}

// VerifyKeyPair verifies that a private and public key form a valid pair
func VerifyKeyPair(privateKey PrivateKey, publicKey PublicKey) bool {
	derivedPublic := DerivePublicKey(privateKey)
	return ConstantTimeEqual(publicKey[:], derivedPublic[:])
}

// IsValidPrivateKey checks if a private key is valid
func IsValidPrivateKey(privateKey PrivateKey) bool {
	// Check if private key is not all zeros
	var zero PrivateKey
	if ConstantTimeEqual(privateKey[:], zero[:]) {
		return false
	}

	// Check if private key is not all ones
	var ones PrivateKey
	for i := range ones {
		ones[i] = 0xFF
	}
	return !ConstantTimeEqual(privateKey[:], ones[:])
}

// IsValidPublicKey checks if a public key is valid
func IsValidPublicKey(publicKey PublicKey) bool {
	// Check if public key is not all zeros
	var zero PublicKey
	if ConstantTimeEqual(publicKey[:], zero[:]) {
		return false
	}

	// Check if public key is not all ones
	var ones PublicKey
	for i := range ones {
		ones[i] = 0xFF
	}
	return !ConstantTimeEqual(publicKey[:], ones[:])
}

// Batch key pair generation

// BatchKeyPairResult represents the result of batch key pair generation
type BatchKeyPairResult struct {
	Index      int
	PrivateKey PrivateKey
	PublicKey  PublicKey
	Error      error
}

// BatchGenerateKeyPairs generates multiple key pairs in parallel
func BatchGenerateKeyPairs(count int) ([]PrivateKey, []PublicKey, error) {
	if count <= 0 {
		return nil, nil, ErrInvalidFragmentCount
	}

	privateKeys := make([]PrivateKey, count)
	publicKeys := make([]PublicKey, count)

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > count {
		numWorkers = count
	}

	// Channel for work distribution
	workChan := make(chan int, count)
	resultChan := make(chan BatchKeyPairResult, count)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				privateKey, publicKey, err := GenerateKeyPair()
				resultChan <- BatchKeyPairResult{
					Index:      index,
					PrivateKey: privateKey,
					PublicKey:  publicKey,
					Error:      err,
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
		privateKeys[result.Index] = result.PrivateKey
		publicKeys[result.Index] = result.PublicKey
	}

	return privateKeys, publicKeys, nil
}

// BatchGenerateKeyPairsFromSeeds generates key pairs from multiple seeds in parallel
func BatchGenerateKeyPairsFromSeeds(seeds [][]byte) ([]PrivateKey, []PublicKey, error) {
	if len(seeds) == 0 {
		return nil, nil, ErrEmptyData
	}

	privateKeys := make([]PrivateKey, len(seeds))
	publicKeys := make([]PublicKey, len(seeds))

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > len(seeds) {
		numWorkers = len(seeds)
	}

	// Channel for work distribution
	workChan := make(chan int, len(seeds))
	resultChan := make(chan BatchKeyPairResult, len(seeds))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				privateKey, publicKey, err := GenerateKeyPairFromSeed(seeds[index])
				resultChan <- BatchKeyPairResult{
					Index:      index,
					PrivateKey: privateKey,
					PublicKey:  publicKey,
					Error:      err,
				}
			}
		}()
	}

	// Send work
	go func() {
		for i := range seeds {
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
		privateKeys[result.Index] = result.PrivateKey
		publicKeys[result.Index] = result.PublicKey
	}

	return privateKeys, publicKeys, nil
}

// Key pair utilities

// KeyPairEqual compares two key pairs for equality
func KeyPairEqual(kp1, kp2 KeyPair) bool {
	return ConstantTimeEqual((*kp1.PrivateKey)[:], (*kp2.PrivateKey)[:]) &&
		ConstantTimeEqual((*kp1.PublicKey)[:], (*kp2.PublicKey)[:])
}

// PrivateKeyEqual compares two private keys for equality
func PrivateKeyEqual(pk1, pk2 PrivateKey) bool {
	return ConstantTimeEqual(pk1[:], pk2[:])
}

// PublicKeyEqual compares two public keys for equality
func PublicKeyEqual(pk1, pk2 PublicKey) bool {
	return ConstantTimeEqual(pk1[:], pk2[:])
}

// SecureErasePrivateKey securely erases a private key from memory
func SecureErasePrivateKey(privateKey *PrivateKey) {
	SecureZero(privateKey[:])
}

// SecureEraseKeyPair securely erases a key pair from memory
func SecureEraseKeyPair(keyPair *KeyPair) {
	SecureZero((*keyPair.PrivateKey)[:])
	SecureZero((*keyPair.PublicKey)[:])
}

// Key derivation functions

// DeriveKeyFromPassword derives a private key from a password using PBKDF2
func DeriveKeyFromPassword(password, salt []byte, iterations int) (PrivateKey, error) {
	if len(password) == 0 {
		return PrivateKey{}, ErrEmptyData
	}

	if len(salt) < 16 {
		return PrivateKey{}, ErrInvalidKeySize
	}

	// Simple PBKDF2-like derivation
	derived := make([]byte, PrivateKeySize)

	current := append(password, salt...)
	for i := 0; i < iterations; i++ {
		hash := ComputeHash(current)
		current = hash[:]
	}

	copy(derived, current[:PrivateKeySize])

	var privateKey PrivateKey
	copy(privateKey[:], derived)

	// Ensure the derived key is valid
	if !IsValidPrivateKey(privateKey) {
		// If invalid, hash again with a different salt
		modifiedSalt := append(salt, 0x01)
		return DeriveKeyFromPassword(password, modifiedSalt, iterations)
	}

	return privateKey, nil
}

// DeriveChildKey derives a child key from a parent private key and index
func DeriveChildKey(parentKey PrivateKey, index uint32) PrivateKey {
	// Simple child key derivation
	hasher := sha256.New()
	hasher.Write(parentKey[:])
	hasher.Write([]byte("TOPAY-Z512-CHILD-KEY"))

	// Add index as bytes
	indexBytes := make([]byte, 4)
	indexBytes[0] = byte(index >> 24)
	indexBytes[1] = byte(index >> 16)
	indexBytes[2] = byte(index >> 8)
	indexBytes[3] = byte(index)
	hasher.Write(indexBytes)

	childHash := hasher.Sum(nil)

	var childKey PrivateKey
	copy(childKey[:], childHash)

	return childKey
}

// Key pair performance benchmarking

// KeyPairBenchmark represents key pair generation performance metrics
type KeyPairBenchmark struct {
	KeyPairsPerSec    float64
	AvgLatencyMs      float64
	BatchSpeedupRatio float64
}

// BenchmarkKeyPairGeneration measures key pair generation performance
func BenchmarkKeyPairGeneration(iterations int) KeyPairBenchmark {
	// Single key pair generation benchmark
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, _, _ = GenerateKeyPair()
	}
	singleDuration := time.Since(start)

	// Batch key pair generation benchmark
	start = time.Now()
	_, _, _ = BatchGenerateKeyPairs(iterations)
	batchDuration := time.Since(start)

	keyPairsPerSec := float64(iterations) / singleDuration.Seconds()
	avgLatencyMs := singleDuration.Seconds() * 1000 / float64(iterations)
	batchSpeedupRatio := singleDuration.Seconds() / batchDuration.Seconds()

	return KeyPairBenchmark{
		KeyPairsPerSec:    keyPairsPerSec,
		AvgLatencyMs:      avgLatencyMs,
		BatchSpeedupRatio: batchSpeedupRatio,
	}
}

// Advanced key pair operations

// RecoverPublicKey attempts to recover a public key from a private key
func RecoverPublicKey(privateKey PrivateKey) (PublicKey, error) {
	if !IsValidPrivateKey(privateKey) {
		return PublicKey{}, ErrInvalidKeySize
	}

	publicKey := DerivePublicKey(privateKey)
	return publicKey, nil
}

// ValidateKeyPairIntegrity performs comprehensive validation of a key pair
func ValidateKeyPairIntegrity(keyPair KeyPair) error {
	// Check private key validity
	if !IsValidPrivateKey(*keyPair.PrivateKey) {
		return ErrInvalidKeySize
	}

	// Check public key validity
	if !IsValidPublicKey(*keyPair.PublicKey) {
		return ErrInvalidKeySize
	}

	// Verify key pair consistency
	if !VerifyKeyPair(*keyPair.PrivateKey, *keyPair.PublicKey) {
		return ErrInvalidKeySize
	}

	return nil
}

// GenerateHDWallet generates a hierarchical deterministic wallet
func GenerateHDWallet(seed []byte, depth int) ([]KeyPair, error) {
	if depth <= 0 || depth > 256 {
		return nil, ErrInvalidFragmentCount
	}

	masterPrivate, masterPublic, err := GenerateKeyPairFromSeed(seed)
	if err != nil {
		return nil, err
	}

	keyPairs := make([]KeyPair, depth)
	keyPairs[0] = KeyPair{PrivateKey: &masterPrivate, PublicKey: &masterPublic}

	currentPrivate := masterPrivate
	for i := 1; i < depth; i++ {
		childPrivate := DeriveChildKey(currentPrivate, uint32(i))
		childPublic := DerivePublicKey(childPrivate)

		keyPairs[i] = KeyPair{PrivateKey: &childPrivate, PublicKey: &childPublic}
		currentPrivate = childPrivate
	}

	return keyPairs, nil
}
