// Package topayz512 provides quantum-safe cryptographic primitives for the TOPAY Foundation blockchain ecosystem.
//
// TOPAY-Z512 offers post-quantum security with ≥512-bit classical security (~256-bit quantum resistance)
// using lattice-based cryptography. It includes support for key generation, hashing, key encapsulation
// mechanisms (KEM), and fragmented processing for parallel computation.
//
// Key Features:
//   - Quantum-safe cryptographic operations
//   - High-performance implementations with batch operations
//   - Fragmented architecture for mobile and IoT device support
//   - Cross-platform compatibility
//   - Developer-friendly API design
//
// Basic Usage:
//
//	// Generate a key pair
//	privateKey, publicKey, err := topayz512.GenerateKeyPair()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Hash data
//	data := []byte("Hello, TOPAY-Z512!")
//	hash := topayz512.ComputeHash(data)
//
//	// KEM operations
//	kemPublic, kemSecret, err := topayz512.KEMKeyGen()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	ciphertext, sharedSecret, err := topayz512.KEMEncapsulate(kemPublic)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	decapsulatedSecret, err := topayz512.KEMDecapsulate(kemSecret, ciphertext)
//	if err != nil {
//	    log.Fatal(err)
//	}
package topayz512

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// Version information
const (
	// Version of the TOPAY-Z512 library
	Version = "0.1.0"

	// SecurityLevel represents the classical security level in bits
	SecurityLevel = 512

	// QuantumSecurityLevel represents the quantum security level in bits
	QuantumSecurityLevel = 256
)

// Cryptographic constants
const (
	// PrivateKeySize is the size of a private key in bytes
	PrivateKeySize = 64

	// PublicKeySize is the size of a public key in bytes
	PublicKeySize = 64

	// HashSize is the size of a hash in bytes
	HashSize = 64

	// KEMPublicKeySize is the size of a KEM public key in bytes
	KEMPublicKeySize = 64

	// KEMSecretKeySize is the size of a KEM secret key in bytes
	KEMSecretKeySize = 64

	// CiphertextSize is the size of a KEM ciphertext in bytes
	CiphertextSize = 64

	// SharedSecretSize is the size of a shared secret in bytes
	SharedSecretSize = 64
)

// Performance constants
const (
	// DefaultBatchSize for batch operations
	DefaultBatchSize = 32

	// CacheLineSize for memory alignment
	CacheLineSize = 64

	// SIMDWidth for vectorized operations
	SIMDWidth = 256

	// PrefetchDistance for memory prefetching
	PrefetchDistance = 64
)

// Fragmentation constants
const (
	// FragmentSize is the default size of each fragment in bytes
	FragmentSize = 256

	// MinFragmentThreshold is the minimum data size to consider fragmentation
	MinFragmentThreshold = 512

	// MaxFragments is the maximum number of fragments allowed
	MaxFragments = 1024
)

// Core cryptographic types

// PrivateKey represents a private key
type PrivateKey [PrivateKeySize]byte

// PublicKey represents a public key
type PublicKey [PublicKeySize]byte

// Hash represents a cryptographic hash
type Hash [HashSize]byte

// KEMPublicKey represents a public key for key encapsulation
type KEMPublicKey [KEMPublicKeySize]byte

// KEMSecretKey represents a secret key for key encapsulation
type KEMSecretKey [KEMSecretKeySize]byte

// Ciphertext represents an encapsulated key
type Ciphertext [CiphertextSize]byte

// SharedSecret represents a shared secret from KEM
type SharedSecret [SharedSecretSize]byte

// KEMKeyPair represents a complete KEM key pair
type KEMKeyPair struct {
	Public KEMPublicKey
	Secret KEMSecretKey
}

// Error definitions
var (
	// ErrInvalidKeySize indicates an invalid key size
	ErrInvalidKeySize = errors.New("invalid key size")

	// ErrInvalidHashSize indicates an invalid hash size
	ErrInvalidHashSize = errors.New("invalid hash size")

	// ErrInvalidCiphertextSize indicates an invalid ciphertext size
	ErrInvalidCiphertextSize = errors.New("invalid ciphertext size")

	// ErrInvalidHexEncoding indicates invalid hex encoding
	ErrInvalidHexEncoding = errors.New("invalid hex encoding")

	// ErrDecapsulationFailed indicates KEM decapsulation failure
	ErrDecapsulationFailed = errors.New("decapsulation failed")

	// ErrFragmentationFailed indicates fragmentation failure
	ErrFragmentationFailed = errors.New("fragmentation failed")

	// ErrReconstructionFailed indicates reconstruction failure
	ErrReconstructionFailed = errors.New("reconstruction failed")

	// ErrEmptyData indicates empty data was provided
	ErrEmptyData = errors.New("empty data provided")

	// ErrInvalidFragmentCount indicates invalid fragment count
	ErrInvalidFragmentCount = errors.New("invalid fragment count")
)

// Utility functions

// SecureRandom generates cryptographically secure random bytes
func SecureRandom(size int) ([]byte, error) {
	data := make([]byte, size)
	_, err := rand.Read(data)
	return data, err
}

// ConstantTimeEqual performs constant-time comparison of two byte slices
func ConstantTimeEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// SecureZero securely zeros a byte slice
func SecureZero(data []byte) {
	for i := range data {
		data[i] = 0
	}

	// Force memory barrier to prevent compiler optimization
	runtime.KeepAlive(data)
}

// FastHexEncode encodes bytes to hex string with optimized performance
func FastHexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// FastHexDecode decodes hex string to bytes with optimized performance
func FastHexDecode(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// System capability detection

// HasSIMDSupport detects if SIMD instructions are available
func HasSIMDSupport() bool {
	// Simplified detection - in production, use proper CPU feature detection
	return runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64"
}

// HasHardwareRNG detects if hardware random number generation is available
func HasHardwareRNG() bool {
	// Simplified detection - in production, use proper hardware detection
	return runtime.GOOS != "js"
}

// OptimalThreadCount returns the optimal number of threads for parallel processing
func OptimalThreadCount() int {
	numCPU := runtime.NumCPU()
	if numCPU <= 2 {
		return numCPU
	}
	// Use 75% of available CPUs for optimal performance
	return (numCPU * 3) / 4
}

// Performance monitoring

// MemoryProfiler provides memory usage profiling
type MemoryProfiler struct {
	startTime time.Time
	startMem  runtime.MemStats
}

// NewMemoryProfiler creates a new memory profiler
func NewMemoryProfiler() *MemoryProfiler {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &MemoryProfiler{
		startTime: time.Now(),
		startMem:  m,
	}
}

// Report returns a memory usage report
func (mp *MemoryProfiler) Report() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	duration := time.Since(mp.startTime)
	allocDiff := m.TotalAlloc - mp.startMem.TotalAlloc

	return fmt.Sprintf("Duration: %v, Memory allocated: %d bytes, GC cycles: %d",
		duration, allocDiff, m.NumGC-mp.startMem.NumGC)
}

// String methods for types

// String returns the hex representation of a PrivateKey
func (pk PrivateKey) String() string {
	return FastHexEncode(pk[:])
}

// String returns the hex representation of a PublicKey
func (pk PublicKey) String() string {
	return FastHexEncode(pk[:])
}

// String returns the hex representation of a Hash
func (h Hash) String() string {
	return FastHexEncode(h[:])
}

// String returns the hex representation of a KEMPublicKey
func (kpk KEMPublicKey) String() string {
	return FastHexEncode(kpk[:])
}

// String returns the hex representation of a KEMSecretKey
func (ksk KEMSecretKey) String() string {
	return FastHexEncode(ksk[:])
}

// String returns the hex representation of a Ciphertext
func (ct Ciphertext) String() string {
	return FastHexEncode(ct[:])
}

// String returns the hex representation of a SharedSecret
func (ss SharedSecret) String() string {
	return FastHexEncode(ss[:])
}

// Bytes methods for types

// Bytes returns the byte representation of a PrivateKey
func (pk PrivateKey) Bytes() []byte {
	return pk[:]
}

// Bytes returns the byte representation of a PublicKey
func (pk PublicKey) Bytes() []byte {
	return pk[:]
}

// Bytes returns the byte representation of a Hash
func (h Hash) Bytes() []byte {
	return h[:]
}

// Bytes returns the byte representation of a KEMPublicKey
func (kpk KEMPublicKey) Bytes() []byte {
	return kpk[:]
}

// Bytes returns the byte representation of a KEMSecretKey
func (ksk KEMSecretKey) Bytes() []byte {
	return ksk[:]
}

// Bytes returns the byte representation of a Ciphertext
func (ct Ciphertext) Bytes() []byte {
	return ct[:]
}

// Bytes returns the byte representation of a SharedSecret
func (ss SharedSecret) Bytes() []byte {
	return ss[:]
}

// FromBytes methods for types

// PrivateKeyFromBytes creates a PrivateKey from bytes
func PrivateKeyFromBytes(data []byte) (PrivateKey, error) {
	if len(data) != PrivateKeySize {
		return PrivateKey{}, ErrInvalidKeySize
	}

	var pk PrivateKey
	copy(pk[:], data)
	return pk, nil
}

// PublicKeyFromBytes creates a PublicKey from bytes
func PublicKeyFromBytes(data []byte) (PublicKey, error) {
	if len(data) != PublicKeySize {
		return PublicKey{}, ErrInvalidKeySize
	}

	var pk PublicKey
	copy(pk[:], data)
	return pk, nil
}

// HashFromBytes creates a Hash from bytes
func HashFromBytes(data []byte) (Hash, error) {
	if len(data) != HashSize {
		return Hash{}, ErrInvalidHashSize
	}

	var h Hash
	copy(h[:], data)
	return h, nil
}

// KEMPublicKeyFromBytes creates a KEMPublicKey from bytes
func KEMPublicKeyFromBytes(data []byte) (KEMPublicKey, error) {
	if len(data) != KEMPublicKeySize {
		return KEMPublicKey{}, ErrInvalidKeySize
	}

	var kpk KEMPublicKey
	copy(kpk[:], data)
	return kpk, nil
}

// KEMSecretKeyFromBytes creates a KEMSecretKey from bytes
func KEMSecretKeyFromBytes(data []byte) (KEMSecretKey, error) {
	if len(data) != KEMSecretKeySize {
		return KEMSecretKey{}, ErrInvalidKeySize
	}

	var ksk KEMSecretKey
	copy(ksk[:], data)
	return ksk, nil
}

// CiphertextFromBytes creates a Ciphertext from bytes
func CiphertextFromBytes(data []byte) (Ciphertext, error) {
	if len(data) != CiphertextSize {
		return Ciphertext{}, ErrInvalidCiphertextSize
	}

	var ct Ciphertext
	copy(ct[:], data)
	return ct, nil
}

// SharedSecretFromBytes creates a SharedSecret from bytes
func SharedSecretFromBytes(data []byte) (SharedSecret, error) {
	if len(data) != SharedSecretSize {
		return SharedSecret{}, ErrInvalidKeySize
	}

	var ss SharedSecret
	copy(ss[:], data)
	return ss, nil
}

// FromHex methods for types

// PrivateKeyFromHex creates a PrivateKey from hex string
func PrivateKeyFromHex(hexStr string) (PrivateKey, error) {
	data, err := FastHexDecode(hexStr)
	if err != nil {
		return PrivateKey{}, ErrInvalidHexEncoding
	}
	return PrivateKeyFromBytes(data)
}

// PublicKeyFromHex creates a PublicKey from hex string
func PublicKeyFromHex(hexStr string) (PublicKey, error) {
	data, err := FastHexDecode(hexStr)
	if err != nil {
		return PublicKey{}, ErrInvalidHexEncoding
	}
	return PublicKeyFromBytes(data)
}

// HashFromHex creates a Hash from hex string
func HashFromHex(hexStr string) (Hash, error) {
	data, err := FastHexDecode(hexStr)
	if err != nil {
		return Hash{}, ErrInvalidHexEncoding
	}
	return HashFromBytes(data)
}

// KEMPublicKeyFromHex creates a KEMPublicKey from hex string
func KEMPublicKeyFromHex(hexStr string) (KEMPublicKey, error) {
	data, err := FastHexDecode(hexStr)
	if err != nil {
		return KEMPublicKey{}, ErrInvalidHexEncoding
	}
	return KEMPublicKeyFromBytes(data)
}

// KEMSecretKeyFromHex creates a KEMSecretKey from hex string
func KEMSecretKeyFromHex(hexStr string) (KEMSecretKey, error) {
	data, err := FastHexDecode(hexStr)
	if err != nil {
		return KEMSecretKey{}, ErrInvalidHexEncoding
	}
	return KEMSecretKeyFromBytes(data)
}

// CiphertextFromHex creates a Ciphertext from hex string
func CiphertextFromHex(hexStr string) (Ciphertext, error) {
	data, err := FastHexDecode(hexStr)
	if err != nil {
		return Ciphertext{}, ErrInvalidHexEncoding
	}
	return CiphertextFromBytes(data)
}

// SharedSecretFromHex creates a SharedSecret from hex string
func SharedSecretFromHex(hexStr string) (SharedSecret, error) {
	data, err := FastHexDecode(hexStr)
	if err != nil {
		return SharedSecret{}, ErrInvalidHexEncoding
	}
	return SharedSecretFromBytes(data)
}
