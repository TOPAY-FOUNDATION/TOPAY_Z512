package topayz512

import (
	"bytes"
	"testing"
	"time"
)

// Test basic hash functionality
func TestHash(t *testing.T) {
	data := []byte("Hello, TOPAY-Z512!")
	hash := ComputeHash(data)

	// Verify hash is not empty
	var zero Hash
	if HashEqual(hash, zero) {
		t.Error("Hash should not be zero")
	}

	// Test deterministic behavior
	hash2 := ComputeHash(data)
	if !HashEqual(hash, hash2) {
		t.Error("Hash should be deterministic")
	}

	// Test different data produces different hash
	differentData := []byte("Different data")
	differentHash := ComputeHash(differentData)
	if HashEqual(hash, differentHash) {
		t.Error("Different data should produce different hash")
	}
}

func TestHashWithSalt(t *testing.T) {
	data := []byte("test data")
	salt := []byte("test salt")

	hash1 := HashWithSalt(data, salt)
	hash2 := HashWithSalt(data, salt)

	if !HashEqual(hash1, hash2) {
		t.Error("Hash with salt should be deterministic")
	}

	// Different salt should produce different hash
	differentSalt := []byte("different salt")
	hash3 := HashWithSalt(data, differentSalt)
	if HashEqual(hash1, hash3) {
		t.Error("Different salt should produce different hash")
	}
}

func TestBatchHash(t *testing.T) {
	inputs := [][]byte{
		[]byte("input1"),
		[]byte("input2"),
		[]byte("input3"),
	}

	results := BatchHash(inputs)

	if len(results) != len(inputs) {
		t.Errorf("Expected %d results, got %d", len(inputs), len(results))
	}

	// Verify each result matches individual hash
	for i, input := range inputs {
		expected := ComputeHash(input)
		if !HashEqual(results[i], expected) {
			t.Errorf("Batch hash result %d doesn't match individual hash", i)
		}
	}
}

// Test key pair functionality
func TestGenerateKeyPair(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Verify keys are valid
	if !IsValidPrivateKey(privateKey) {
		t.Error("Generated private key is invalid")
	}

	if !IsValidPublicKey(publicKey) {
		t.Error("Generated public key is invalid")
	}

	// Verify key pair consistency
	if !VerifyKeyPair(privateKey, publicKey) {
		t.Error("Generated key pair is inconsistent")
	}
}

func TestDerivePublicKey(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	derivedPublic := DerivePublicKey(privateKey)
	if !PublicKeyEqual(publicKey, derivedPublic) {
		t.Error("Derived public key doesn't match original")
	}
}

func TestGenerateKeyPairFromSeed(t *testing.T) {
	seed := []byte("this is a test seed that is long enough")

	privateKey1, publicKey1, err := GenerateKeyPairFromSeed(seed)
	if err != nil {
		t.Fatalf("Failed to generate key pair from seed: %v", err)
	}

	privateKey2, publicKey2, err := GenerateKeyPairFromSeed(seed)
	if err != nil {
		t.Fatalf("Failed to generate key pair from seed: %v", err)
	}

	// Should be deterministic
	if !PrivateKeyEqual(privateKey1, privateKey2) {
		t.Error("Key pair generation from seed should be deterministic")
	}

	if !PublicKeyEqual(publicKey1, publicKey2) {
		t.Error("Key pair generation from seed should be deterministic")
	}
}

func TestBatchGenerateKeyPairs(t *testing.T) {
	count := 10
	privateKeys, publicKeys, err := BatchGenerateKeyPairs(count)
	if err != nil {
		t.Fatalf("Failed to batch generate key pairs: %v", err)
	}

	if len(privateKeys) != count {
		t.Errorf("Expected %d private keys, got %d", count, len(privateKeys))
	}

	if len(publicKeys) != count {
		t.Errorf("Expected %d public keys, got %d", count, len(publicKeys))
	}

	// Verify all key pairs are valid and unique
	for i := 0; i < count; i++ {
		if !IsValidPrivateKey(privateKeys[i]) {
			t.Errorf("Private key %d is invalid", i)
		}

		if !IsValidPublicKey(publicKeys[i]) {
			t.Errorf("Public key %d is invalid", i)
		}

		if !VerifyKeyPair(privateKeys[i], publicKeys[i]) {
			t.Errorf("Key pair %d is inconsistent", i)
		}

		// Check uniqueness
		for j := i + 1; j < count; j++ {
			if PrivateKeyEqual(privateKeys[i], privateKeys[j]) {
				t.Errorf("Private keys %d and %d are identical", i, j)
			}
		}
	}
}

// Test KEM functionality
func TestKEMKeyGen(t *testing.T) {
	publicKey, secretKey, err := KEMKeyGen()
	if err != nil {
		t.Fatalf("Failed to generate KEM key pair: %v", err)
	}

	if !IsValidKEMPublicKey(publicKey) {
		t.Error("Generated KEM public key is invalid")
	}

	if !IsValidKEMSecretKey(secretKey) {
		t.Error("Generated KEM secret key is invalid")
	}

	if !VerifyKEMKeyPair(publicKey, secretKey) {
		t.Error("Generated KEM key pair is inconsistent")
	}
}

func TestKEMEncapsulateDecapsulate(t *testing.T) {
	publicKey, secretKey, err := KEMKeyGen()
	if err != nil {
		t.Fatalf("Failed to generate KEM key pair: %v", err)
	}

	ciphertext, sharedSecret1, err := KEMEncapsulate(publicKey)
	if err != nil {
		t.Fatalf("Failed to encapsulate: %v", err)
	}

	if !IsValidCiphertext(ciphertext) {
		t.Error("Generated ciphertext is invalid")
	}

	if !IsValidSharedSecret(sharedSecret1) {
		t.Error("Generated shared secret is invalid")
	}

	sharedSecret2, err := KEMDecapsulate(secretKey, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decapsulate: %v", err)
	}

	if !SharedSecretEqual(sharedSecret1, sharedSecret2) {
		t.Error("Shared secrets don't match")
	}
}

func TestBatchKEMOperations(t *testing.T) {
	count := 5

	// Test batch key generation
	publicKeys, secretKeys, err := BatchKEMKeyGen(count)
	if err != nil {
		t.Fatalf("Failed to batch generate KEM keys: %v", err)
	}

	if len(publicKeys) != count || len(secretKeys) != count {
		t.Errorf("Expected %d keys, got %d public and %d secret", count, len(publicKeys), len(secretKeys))
	}

	// Test batch encapsulation
	ciphertexts, sharedSecrets1, err := BatchKEMEncapsulate(publicKeys)
	if err != nil {
		t.Fatalf("Failed to batch encapsulate: %v", err)
	}

	// Test batch decapsulation
	sharedSecrets2, err := BatchKEMDecapsulate(secretKeys, ciphertexts)
	if err != nil {
		t.Fatalf("Failed to batch decapsulate: %v", err)
	}

	// Verify shared secrets match
	for i := 0; i < count; i++ {
		if !SharedSecretEqual(sharedSecrets1[i], sharedSecrets2[i]) {
			t.Errorf("Shared secrets %d don't match", i)
		}
	}
}

// Test fragmentation functionality
func TestFragmentData(t *testing.T) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}

	result, err := FragmentData(data)
	if err != nil {
		t.Fatalf("Failed to fragment data: %v", err)
	}

	if len(result.Fragments) == 0 {
		t.Error("No fragments created")
	}

	if result.Metadata.OriginalSize != uint64(len(data)) {
		t.Errorf("Expected original size %d, got %d", len(data), result.Metadata.OriginalSize)
	}

	// Verify fragment integrity
	for i, fragment := range result.Fragments {
		if fragment.Index != uint32(i) {
			t.Errorf("Fragment %d has wrong index %d", i, fragment.Index)
		}

		if fragment.Total != uint32(len(result.Fragments)) {
			t.Errorf("Fragment %d has wrong total %d", i, fragment.Total)
		}

		// Verify fragment checksum
		computedChecksum := ComputeHash(fragment.Data)
		if !HashEqual(computedChecksum, fragment.Checksum) {
			t.Errorf("Fragment %d has invalid checksum", i)
		}
	}
}

func TestReconstructData(t *testing.T) {
	originalData := make([]byte, 1024)
	for i := range originalData {
		originalData[i] = byte(i)
	}

	// Fragment the data
	fragResult, err := FragmentData(originalData)
	if err != nil {
		t.Fatalf("Failed to fragment data: %v", err)
	}

	// Reconstruct the data
	reconResult, err := ReconstructData(fragResult.Fragments)
	if err != nil {
		t.Fatalf("Failed to reconstruct data: %v", err)
	}

	if !reconResult.IsComplete {
		t.Error("Reconstruction not complete")
	}

	if !bytes.Equal(originalData, reconResult.Data) {
		t.Error("Reconstructed data doesn't match original")
	}
}

func TestParallelFragmentation(t *testing.T) {
	data := make([]byte, 2048)
	for i := range data {
		data[i] = byte(i)
	}

	// Test parallel fragmentation
	result1, err := FragmentData(data)
	if err != nil {
		t.Fatalf("Failed to fragment data: %v", err)
	}

	result2, err := ParallelFragmentData(data)
	if err != nil {
		t.Fatalf("Failed to parallel fragment data: %v", err)
	}

	// Results should be equivalent (though fragment IDs may differ)
	if len(result1.Fragments) != len(result2.Fragments) {
		t.Error("Parallel fragmentation produced different number of fragments")
	}

	if result1.Metadata.OriginalSize != result2.Metadata.OriginalSize {
		t.Error("Parallel fragmentation produced different metadata")
	}
}

func TestFragmentSerialization(t *testing.T) {
	data := []byte("test fragment data")
	fragment := Fragment{
		ID:       12345,
		Index:    0,
		Total:    1,
		Data:     data,
		Checksum: ComputeHash(data),
	}

	// Serialize
	serialized := SerializeFragment(fragment)

	// Deserialize
	deserialized, err := DeserializeFragment(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize fragment: %v", err)
	}

	// Verify
	if deserialized.ID != fragment.ID {
		t.Error("Deserialized fragment ID doesn't match")
	}

	if deserialized.Index != fragment.Index {
		t.Error("Deserialized fragment index doesn't match")
	}

	if deserialized.Total != fragment.Total {
		t.Error("Deserialized fragment total doesn't match")
	}

	if !bytes.Equal(deserialized.Data, fragment.Data) {
		t.Error("Deserialized fragment data doesn't match")
	}

	if !HashEqual(deserialized.Checksum, fragment.Checksum) {
		t.Error("Deserialized fragment checksum doesn't match")
	}
}

func TestMobileLatencyEstimate(t *testing.T) {
	dataSize := 1024 * 1024 // 1MB
	estimate := EstimateMobileLatency(dataSize)

	if estimate.TotalMs <= 0 {
		t.Error("Mobile latency estimate should be positive")
	}

	if estimate.RecommendedChunks <= 0 {
		t.Error("Recommended chunks should be positive")
	}

	if estimate.FragmentationMs <= 0 {
		t.Error("Fragmentation latency should be positive")
	}

	if estimate.ReconstructionMs <= 0 {
		t.Error("Reconstruction latency should be positive")
	}
}

// Test utility functions
func TestConstantTimeEqual(t *testing.T) {
	a := []byte{1, 2, 3, 4}
	b := []byte{1, 2, 3, 4}
	c := []byte{1, 2, 3, 5}
	d := []byte{1, 2, 3}

	if !ConstantTimeEqual(a, b) {
		t.Error("Equal slices should return true")
	}

	if ConstantTimeEqual(a, c) {
		t.Error("Different slices should return false")
	}

	if ConstantTimeEqual(a, d) {
		t.Error("Different length slices should return false")
	}
}

func TestSecureZero(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	SecureZero(data)

	for i, b := range data {
		if b != 0 {
			t.Errorf("Byte %d not zeroed: %d", i, b)
		}
	}
}

func TestHexEncoding(t *testing.T) {
	data := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}

	encoded := FastHexEncode(data)
	decoded, err := FastHexDecode(encoded)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Error("Hex encoding/decoding roundtrip failed")
	}
}

func TestSystemCapabilities(t *testing.T) {
	// These tests just verify the functions don't panic
	_ = HasSIMDSupport()
	_ = HasHardwareRNG()

	threadCount := OptimalThreadCount()
	if threadCount <= 0 {
		t.Error("Optimal thread count should be positive")
	}
}

// Test type conversions
func TestTypeConversions(t *testing.T) {
	// Test PrivateKey conversions
	privateKey, _, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Test bytes conversion
	privateBytes := privateKey.Bytes()
	if len(privateBytes) != PrivateKeySize {
		t.Errorf("Expected %d bytes, got %d", PrivateKeySize, len(privateBytes))
	}

	// Test from bytes
	privateKey2, err := PrivateKeyFromBytes(privateBytes)
	if err != nil {
		t.Fatalf("Failed to create private key from bytes: %v", err)
	}

	if !PrivateKeyEqual(privateKey, privateKey2) {
		t.Error("Private key conversion failed")
	}

	// Test hex conversion
	hexStr := privateKey.String()
	privateKey3, err := PrivateKeyFromHex(hexStr)
	if err != nil {
		t.Fatalf("Failed to create private key from hex: %v", err)
	}

	if !PrivateKeyEqual(privateKey, privateKey3) {
		t.Error("Private key hex conversion failed")
	}
}

// Benchmark tests
func BenchmarkHashTest(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComputeHash(data)
	}
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = GenerateKeyPair()
	}
}

func BenchmarkKEMEncapsulate(b *testing.B) {
	publicKey, _, _ := KEMKeyGen()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = KEMEncapsulate(publicKey)
	}
}

func BenchmarkFragmentData(b *testing.B) {
	data := make([]byte, 4096)
	for i := range data {
		data[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FragmentData(data)
	}
}

// Integration tests
func TestCompleteWorkflow(t *testing.T) {
	// Generate key pairs
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Generate KEM key pair
	kemPublic, kemSecret, err := KEMKeyGen()
	if err != nil {
		t.Fatalf("Failed to generate KEM key pair: %v", err)
	}

	// Test data
	data := []byte("This is a test message for the complete workflow")

	// Hash the data
	hash := ComputeHash(data)

	// KEM operations
	ciphertext, sharedSecret1, err := KEMEncapsulate(kemPublic)
	if err != nil {
		t.Fatalf("Failed to encapsulate: %v", err)
	}

	sharedSecret2, err := KEMDecapsulate(kemSecret, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decapsulate: %v", err)
	}

	if !SharedSecretEqual(sharedSecret1, sharedSecret2) {
		t.Error("Shared secrets don't match")
	}

	// Fragment the data
	fragResult, err := FragmentData(data)
	if err != nil {
		t.Fatalf("Failed to fragment data: %v", err)
	}

	// Reconstruct the data
	reconResult, err := ReconstructData(fragResult.Fragments)
	if err != nil {
		t.Fatalf("Failed to reconstruct data: %v", err)
	}

	if !bytes.Equal(data, reconResult.Data) {
		t.Error("Reconstructed data doesn't match original")
	}

	// Verify hash
	if !VerifyHash(data, hash) {
		t.Error("Hash verification failed")
	}

	// Verify key pair
	if !VerifyKeyPair(privateKey, publicKey) {
		t.Error("Key pair verification failed")
	}

	// Verify KEM key pair
	if !VerifyKEMKeyPair(kemPublic, kemSecret) {
		t.Error("KEM key pair verification failed")
	}
}

func TestErrorHandling(t *testing.T) {
	// Test invalid key sizes
	_, err := PrivateKeyFromBytes([]byte{1, 2, 3}) // Too short
	if err == nil {
		t.Error("Should fail with invalid key size")
	}

	// Test invalid hex
	_, err = PrivateKeyFromHex("invalid hex")
	if err == nil {
		t.Error("Should fail with invalid hex")
	}

	// Test empty data fragmentation
	_, err = FragmentData([]byte{})
	if err == nil {
		t.Error("Should fail with empty data")
	}

	// Test invalid fragment reconstruction
	fragments := []Fragment{
		{ID: 1, Index: 0, Total: 2, Data: []byte("test"), Checksum: ComputeHash([]byte("test"))},
		// Missing fragment 1
	}
	_, err = ReconstructData(fragments)
	if err == nil {
		t.Error("Should fail with incomplete fragments")
	}
}

func TestMemoryProfiler(t *testing.T) {
	profiler := NewMemoryProfiler()

	// Do some work
	time.Sleep(10 * time.Millisecond)
	data := make([]byte, 1024)
	_ = ComputeHash(data)

	report := profiler.Report()
	if len(report) == 0 {
		t.Error("Memory profiler should return a report")
	}
}
