// Package kem provides a post-quantum Key Encapsulation Mechanism (KEM)
// based on the Learning With Errors (LWE) problem.
package kem

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/topay-foundation/topayz512/pkg/hash"
	"github.com/topay-foundation/topayz512/pkg/keypair"
	"golang.org/x/crypto/sha3"
)

// Parameters for the LWE problem
const (
	// N is the lattice dimension for the LWE problem
	N = 1024

	// Q is the modulus for the LWE problem
	Q = 65537

	// Sigma is the standard deviation for the error distribution
	Sigma = 3.2

	// SecretLength is the length of the shared secret in bytes
	SecretLength = hash.HashSizeBytes

	// CiphertextSizeBytes is the length of the ciphertext in bytes
	CiphertextSizeBytes = N*2 + SecretLength
)

// Ciphertext represents a TOPAY-Z512 ciphertext
type Ciphertext struct {
	// B is the first component of the ciphertext (vector b)
	B []uint16
	// V is the second component of the ciphertext (vector v)
	V []byte
}

// SharedSecret represents a TOPAY-Z512 shared secret
type SharedSecret struct {
	bytes []byte
}

// NewCiphertext creates a new ciphertext from components
func NewCiphertext(b []uint16, v []byte) (*Ciphertext, error) {
	if len(b) != N {
		return nil, errors.New("invalid vector b length")
	}
	if len(v) != SecretLength {
		return nil, errors.New("invalid vector v length")
	}

	return &Ciphertext{
		B: b,
		V: v,
	}, nil
}

// Bytes returns the ciphertext as bytes
func (c *Ciphertext) Bytes() []byte {
	bytes := make([]byte, CiphertextSizeBytes)

	// Convert vector b to bytes
	for i, val := range c.B {
		binary.BigEndian.PutUint16(bytes[i*2:], val)
	}

	// Add vector v
	copy(bytes[N*2:], c.V)

	return bytes
}

// FromBytes creates a ciphertext from bytes
func CiphertextFromBytes(bytes []byte) (*Ciphertext, error) {
	if len(bytes) != CiphertextSizeBytes {
		return nil, errors.New("invalid ciphertext length")
	}

	b := make([]uint16, N)

	// Extract vector b
	for i := 0; i < N; i++ {
		b[i] = binary.BigEndian.Uint16(bytes[i*2 : i*2+2])
	}

	// Extract vector v
	v := make([]byte, SecretLength)
	copy(v, bytes[N*2:N*2+SecretLength])

	return &Ciphertext{
		B: b,
		V: v,
	}, nil
}

// String converts the ciphertext to a hexadecimal string
func (c *Ciphertext) String() string {
	return hex.EncodeToString(c.Bytes())
}

// CiphertextFromHex creates a ciphertext from a hexadecimal string
func CiphertextFromHex(hexStr string) (*Ciphertext, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %v", err)
	}

	return CiphertextFromBytes(bytes)
}

// NewSharedSecret creates a new shared secret from bytes
func NewSharedSecret(bytes []byte) (*SharedSecret, error) {
	if len(bytes) != SecretLength {
		return nil, errors.New("invalid shared secret length")
	}

	secretBytes := make([]byte, SecretLength)
	copy(secretBytes, bytes)

	return &SharedSecret{
		bytes: secretBytes,
	}, nil
}

// Bytes returns the shared secret as bytes
func (s *SharedSecret) Bytes() []byte {
	bytes := make([]byte, SecretLength)
	copy(bytes, s.bytes)
	return bytes
}

// String converts the shared secret to a hexadecimal string
func (s *SharedSecret) String() string {
	return hex.EncodeToString(s.bytes)
}

// SharedSecretFromHex creates a shared secret from a hexadecimal string
func SharedSecretFromHex(hexStr string) (*SharedSecret, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %v", err)
	}

	return NewSharedSecret(bytes)
}

// Keygen generates a key pair for the KEM
func Keygen() (*keypair.KeyPair, error) {
	// For now, we'll use the existing keypair generation
	// In a real implementation, this would generate LWE-specific keys
	keyPair, err := keypair.GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	return &keyPair, nil
}

// Encapsulate encapsulates a shared secret using a public key
func Encapsulate(publicKey *keypair.PublicKey) (*Ciphertext, *SharedSecret, error) {
	// This is a placeholder implementation
	// In a real implementation, this would use LWE encapsulation

	// Generate a random message
	message := make([]byte, SecretLength)
	_, err := rand.Read(message)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random message: %v", err)
	}

	// Hash the message with the public key to create the shared secret
	h := sha3.New512()
	h.Write(message)
	h.Write(publicKey.Bytes())
	result := h.Sum(nil)

	sharedSecretBytes := make([]byte, SecretLength)
	copy(sharedSecretBytes, result)

	sharedSecret, err := NewSharedSecret(sharedSecretBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create shared secret: %v", err)
	}

	// Create a dummy ciphertext (this would be the actual LWE ciphertext in a real implementation)
	b := make([]uint16, N) // Placeholder
	v := make([]byte, SecretLength)
	copy(v, message) // Placeholder

	ciphertext, err := NewCiphertext(b, v)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ciphertext: %v", err)
	}

	return ciphertext, sharedSecret, nil
}

// Decapsulate decapsulates a shared secret using a private key and ciphertext
func Decapsulate(privateKey *keypair.PrivateKey, ciphertext *Ciphertext) (*SharedSecret, error) {
	// This is a placeholder implementation
	// In a real implementation, this would use LWE decapsulation

	// Derive the public key from the private key
	publicKey := keypair.PrivateToPublic(*privateKey)

	// Hash the ciphertext.V (which contains the message) with the public key
	// to recreate the shared secret - must match the encapsulation method
	h := sha3.New512()
	h.Write(ciphertext.V)
	h.Write(publicKey.Bytes())
	result := h.Sum(nil)

	sharedSecretBytes := make([]byte, SecretLength)
	copy(sharedSecretBytes, result)

	return NewSharedSecret(sharedSecretBytes)
}
