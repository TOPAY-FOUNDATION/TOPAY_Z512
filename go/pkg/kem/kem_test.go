package kem

import (
	"testing"
)

func TestKEM(t *testing.T) {
	// Generate a key pair
	keyPair, err := Keygen()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Encapsulate a shared secret
	ciphertext, sharedSecret1, err := Encapsulate(&keyPair.PublicKey)
	if err != nil {
		t.Fatalf("Failed to encapsulate: %v", err)
	}

	// Decapsulate the shared secret
	sharedSecret2, err := Decapsulate(&keyPair.PrivateKey, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decapsulate: %v", err)
	}

	// Verify that the shared secrets match
	if string(sharedSecret1.Bytes()) != string(sharedSecret2.Bytes()) {
		t.Errorf("Shared secrets do not match")
	}
}

func TestCiphertextSerialization(t *testing.T) {
	// Generate a key pair
	keyPair, err := Keygen()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Encapsulate a shared secret
	ciphertext, _, err := Encapsulate(&keyPair.PublicKey)
	if err != nil {
		t.Fatalf("Failed to encapsulate: %v", err)
	}

	// Test byte serialization
	bytes := ciphertext.Bytes()
	ciphertext2, err := CiphertextFromBytes(bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize ciphertext: %v", err)
	}

	if string(ciphertext.Bytes()) != string(ciphertext2.Bytes()) {
		t.Errorf("Ciphertext serialization failed")
	}

	// Test hex serialization
	hexStr := ciphertext.String()
	ciphertext3, err := CiphertextFromHex(hexStr)
	if err != nil {
		t.Fatalf("Failed to deserialize ciphertext from hex: %v", err)
	}

	if string(ciphertext.Bytes()) != string(ciphertext3.Bytes()) {
		t.Errorf("Ciphertext hex serialization failed")
	}
}

func TestSharedSecretSerialization(t *testing.T) {
	// Generate a key pair
	keyPair, err := Keygen()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Encapsulate a shared secret
	_, sharedSecret, err := Encapsulate(&keyPair.PublicKey)
	if err != nil {
		t.Fatalf("Failed to encapsulate: %v", err)
	}

	// Test hex serialization
	hexStr := sharedSecret.String()
	sharedSecret2, err := SharedSecretFromHex(hexStr)
	if err != nil {
		t.Fatalf("Failed to deserialize shared secret from hex: %v", err)
	}

	if string(sharedSecret.Bytes()) != string(sharedSecret2.Bytes()) {
		t.Errorf("Shared secret hex serialization failed")
	}
}

func TestInvalidCiphertext(t *testing.T) {
	// Test invalid ciphertext length
	_, err := CiphertextFromBytes([]byte{1, 2, 3})
	if err == nil {
		t.Errorf("Expected error for invalid ciphertext length")
	}

	// Test invalid hex string
	_, err = CiphertextFromHex("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid hex string")
	}
}

func TestInvalidSharedSecret(t *testing.T) {
	// Test invalid shared secret length
	_, err := NewSharedSecret([]byte{1, 2, 3})
	if err == nil {
		t.Errorf("Expected error for invalid shared secret length")
	}

	// Test invalid hex string
	_, err = SharedSecretFromHex("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid hex string")
	}
}
