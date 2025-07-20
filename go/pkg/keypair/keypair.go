// Package keypair provides a 512-bit cryptographic key pair implementation
// that can be used as part of the TOPAY-Z512 cryptographic suite.
package keypair

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// KeySizeBytes is the size of TOPAY-Z512 keys in bytes (512 bits = 64 bytes)
const KeySizeBytes = 64

// PrivateKey represents a TOPAY-Z512 private key (512 bits)
type PrivateKey [KeySizeBytes]byte

// PublicKey represents a TOPAY-Z512 public key (512 bits)
type PublicKey [KeySizeBytes]byte

// KeyPair represents a TOPAY-Z512 key pair (private key and public key)
type KeyPair struct {
	PrivateKey PrivateKey
	PublicKey  PublicKey
}

// GenerateKeyPair generates a new key pair using secure random data
func GenerateKeyPair() (KeyPair, error) {
	var privateKey PrivateKey
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to generate random private key: %w", err)
	}

	publicKey := DerivePublicKey(privateKey)

	return KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// DerivePublicKey derives a public key from a private key
func DerivePublicKey(privateKey PrivateKey) PublicKey {
	var publicKey PublicKey
	hasher := sha3.New512()
	hasher.Write(privateKey[:])
	result := hasher.Sum(nil)
	copy(publicKey[:], result)
	return publicKey
}

// Bytes returns the private key as a byte slice
func (pk PrivateKey) Bytes() []byte {
	return pk[:]
}

// String returns the private key as a hexadecimal string
func (pk PrivateKey) String() string {
	return hex.EncodeToString(pk[:])
}

// PrivateKeyFromHex creates a private key from a hexadecimal string
func PrivateKeyFromHex(hexStr string) (PrivateKey, error) {
	var privateKey PrivateKey

	if len(hexStr) != KeySizeBytes*2 {
		return privateKey, errors.New("invalid hex string length")
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return privateKey, fmt.Errorf("invalid hex string: %w", err)
	}

	copy(privateKey[:], bytes)
	return privateKey, nil
}

// Bytes returns the public key as a byte slice
func (pk PublicKey) Bytes() []byte {
	return pk[:]
}

// String returns the public key as a hexadecimal string
func (pk PublicKey) String() string {
	return hex.EncodeToString(pk[:])
}

// PublicKeyFromHex creates a public key from a hexadecimal string
func PublicKeyFromHex(hexStr string) (PublicKey, error) {
	var publicKey PublicKey

	if len(hexStr) != KeySizeBytes*2 {
		return publicKey, errors.New("invalid hex string length")
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return publicKey, fmt.Errorf("invalid hex string: %w", err)
	}

	copy(publicKey[:], bytes)
	return publicKey, nil
}

// NewKeyPair creates a key pair from existing private and public keys
func NewKeyPair(privateKey PrivateKey, publicKey PublicKey) KeyPair {
	return KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

// PrivateToPublic is a convenience function to derive a public key from a private key
func PrivateToPublic(privateKey PrivateKey) PublicKey {
	return DerivePublicKey(privateKey)
}