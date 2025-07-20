// Package hash provides a 512-bit cryptographic hash function implementation
// that can be used as part of the TOPAY-Z512 cryptographic suite.
package hash

import (
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// HashSizeBytes is the size of TOPAY-Z512 hash output in bytes (512 bits = 64 bytes)
const HashSizeBytes = 64

// Hash represents a TOPAY-Z512 hash value (512 bits)
type Hash [HashSizeBytes]byte

// New creates a new hash from the given data
func New(data []byte) Hash {
	var hash Hash
	hasher := sha3.New512()
	hasher.Write(data)
	result := hasher.Sum(nil)
	copy(hash[:], result)
	return hash
}

// Combine creates a new hash by combining two input values
// This implementation ensures the result is different from simply concatenating the inputs
func Combine(data1, data2 []byte) Hash {
	// First hash each input separately
	hash1 := New(data1)
	hash2 := New(data2)
	
	// Then combine the hashes with a separator byte to ensure it's different from concatenation
	var combined []byte
	combined = append(combined, hash1[:]...)
	combined = append(combined, 0xFF) // Separator byte that wouldn't be in normal concatenation
	combined = append(combined, hash2[:]...)
	
	// Hash the combined result
	return New(combined)
}

// Bytes returns the hash value as a byte slice
func (h Hash) Bytes() []byte {
	return h[:]
}

// String returns the hash as a hexadecimal string
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// FromHex creates a hash from a hexadecimal string
func FromHex(hexStr string) (Hash, error) {
	var hash Hash

	if len(hexStr) != HashSizeBytes*2 {
		return hash, fmt.Errorf("invalid hex length: %d, expected %d",
			len(hexStr), HashSizeBytes*2)
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return hash, errors.New("invalid hex string: " + err.Error())
	}

	copy(hash[:], bytes)
	return hash, nil
}

// Sum512 is a convenience function to hash data
func Sum512(data []byte) [HashSizeBytes]byte {
	return New(data)
}

// SumCombine is a convenience function to hash two pieces of data together
func SumCombine(data1, data2 []byte) [HashSizeBytes]byte {
	return Combine(data1, data2)
}
