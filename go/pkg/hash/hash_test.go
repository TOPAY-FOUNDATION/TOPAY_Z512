package hash

import (
	"bytes"
	"testing"
)

func TestHashSize(t *testing.T) {
	data := []byte("TOPAY-Z512 test data")
	// Create hash and directly check its length without storing the variable
	if size := len(New(data)); size != HashSizeBytes {
		t.Errorf("Expected hash size %d, got %d", HashSizeBytes, size)
	}
}

func TestHashDeterministic(t *testing.T) {
	data := []byte("TOPAY-Z512 test data")
	hash1 := New(data)
	hash2 := New(data)

	if hash1 != hash2 {
		t.Errorf("Expected identical hashes for identical inputs")
	}
}

func TestHashDifferentInputs(t *testing.T) {
	data1 := []byte("TOPAY-Z512 test data 1")
	data2 := []byte("TOPAY-Z512 test data 2")
	hash1 := New(data1)
	hash2 := New(data2)

	if hash1 == hash2 {
		t.Errorf("Expected different hashes for different inputs")
	}
}

func TestHashCombine(t *testing.T) {
	data1 := []byte("TOPAY-Z512")
	data2 := []byte("test data")

	// Combined hash
	combined := Combine(data1, data2)

	// Concatenated hash
	concatenated := append([]byte{}, data1...)
	concatenated = append(concatenated, data2...)
	concatHash := New(concatenated)

	// These should be different
	if combined == concatHash {
		t.Errorf("Expected combined hash to be different from concatenated hash")
	}
}

func TestHexConversion(t *testing.T) {
	data := []byte("TOPAY-Z512 hex conversion test")
	hash := New(data)
	hexStr := hash.String()

	hash2, err := FromHex(hexStr)
	if err != nil {
		t.Errorf("Error converting from hex: %v", err)
	}

	if hash != hash2 {
		t.Errorf("Expected identical hashes after hex conversion")
	}
}

func TestInvalidHex(t *testing.T) {
	// Test invalid length
	_, err := FromHex("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid hex length")
	}

	// Test invalid characters
	_, err = FromHex("zz" + repeat("00", HashSizeBytes-1))
	if err == nil {
		t.Errorf("Expected error for invalid hex characters")
	}
}

// Helper function to repeat a string
func repeat(s string, count int) string {
	var result string
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func TestBytes(t *testing.T) {
	data := []byte("TOPAY-Z512 bytes test")
	hash := New(data)
	hashBytes := hash.Bytes()

	if !bytes.Equal(hashBytes, hash[:]) {
		t.Errorf("Expected Bytes() to return the hash bytes")
	}
}
