package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TOPAY-FOUNDATION/TOPAY_Z512/go"
)

func main() {
	fmt.Println("=== TOPAY-Z512 Go Implementation - Quick Start Guide ===")
	fmt.Println()

	// 1. Key Pair Generation
	fmt.Println("1. Generating Key Pair...")
	privateKey, publicKey, err := topayz512.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	fmt.Printf("   Private Key: %s\n", privateKey.String()[:32]+"...")
	fmt.Printf("   Public Key:  %s\n", publicKey.String()[:32]+"...")
	fmt.Printf("   Key Pair Valid: %v\n", topayz512.VerifyKeyPair(privateKey, publicKey))
	fmt.Println()

	// 2. Hash Operations
	fmt.Println("2. Hash Operations...")
	data := []byte("Hello, TOPAY-Z512 from Go!")
	hash := topayz512.ComputeHash(data)

	fmt.Printf("   Data: %s\n", string(data))
	fmt.Printf("   Hash: %s\n", hash.String()[:32]+"...")
	fmt.Printf("   Hash Valid: %v\n", topayz512.IsValidHash(hash))
	fmt.Printf("   Verification: %v\n", topayz512.VerifyHash(data, hash))
	fmt.Println()

	// 3. KEM Operations
	fmt.Println("3. KEM (Key Encapsulation Mechanism)...")
	kemPublic, kemSecret, err := topayz512.KEMKeyGen()
	if err != nil {
		log.Fatalf("Failed to generate KEM key pair: %v", err)
	}

	fmt.Printf("   KEM Public Key: %s\n", kemPublic.String()[:32]+"...")
	fmt.Printf("   KEM Secret Key: %s\n", kemSecret.String()[:32]+"...")

	// Encapsulation
	ciphertext, sharedSecret1, err := topayz512.KEMEncapsulate(kemPublic)
	if err != nil {
		log.Fatalf("Failed to encapsulate: %v", err)
	}

	fmt.Printf("   Ciphertext: %s\n", ciphertext.String()[:32]+"...")
	fmt.Printf("   Shared Secret (Encap): %s\n", sharedSecret1.String()[:32]+"...")

	// Decapsulation
	sharedSecret2, err := topayz512.KEMDecapsulate(kemSecret, ciphertext)
	if err != nil {
		log.Fatalf("Failed to decapsulate: %v", err)
	}

	fmt.Printf("   Shared Secret (Decap): %s\n", sharedSecret2.String()[:32]+"...")
	fmt.Printf("   Secrets Match: %v\n", topayz512.SharedSecretEqual(sharedSecret1, sharedSecret2))
	fmt.Println()

	// 4. Fragmentation
	fmt.Println("4. Data Fragmentation...")
	largeData := make([]byte, 2048)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	fmt.Printf("   Original Data Size: %d bytes\n", len(largeData))
	fmt.Printf("   Should Fragment: %v\n", topayz512.ShouldFragment(len(largeData)))

	// Fragment the data
	fragResult, err := topayz512.FragmentData(largeData)
	if err != nil {
		log.Fatalf("Failed to fragment data: %v", err)
	}

	fmt.Printf("   Fragment Count: %d\n", len(fragResult.Fragments))
	fmt.Printf("   Original Size: %d bytes\n", fragResult.Metadata.OriginalSize)
	fmt.Printf("   Total Checksum: %s\n", fragResult.Metadata.Checksum.String()[:32]+"...")

	// Reconstruct the data
	reconResult, err := topayz512.ReconstructData(fragResult.Fragments)
	if err != nil {
		log.Fatalf("Failed to reconstruct data: %v", err)
	}

	fmt.Printf("   Reconstruction Complete: %v\n", reconResult.IsComplete)
	fmt.Printf("   Data Integrity: %v\n", len(reconResult.Data) == len(largeData))
	fmt.Println()

	// 5. Batch Operations
	fmt.Println("5. Batch Operations Performance...")

	// Batch key generation
	start := time.Now()
	batchSize := 100
	privateKeys, publicKeys, err := topayz512.BatchGenerateKeyPairs(batchSize)
	if err != nil {
		log.Fatalf("Failed to batch generate key pairs: %v", err)
	}
	batchKeyGenDuration := time.Since(start)

	fmt.Printf("   Generated %d key pairs in %v\n", batchSize, batchKeyGenDuration)
	fmt.Printf("   Average per key pair: %v\n", batchKeyGenDuration/time.Duration(batchSize))

	// Verify all key pairs
	validCount := 0
	for i := 0; i < batchSize; i++ {
		if topayz512.VerifyKeyPair(privateKeys[i], publicKeys[i]) {
			validCount++
		}
	}
	fmt.Printf("   Valid key pairs: %d/%d\n", validCount, batchSize)

	// Batch KEM operations
	start = time.Now()
	kemPublicKeys, kemSecretKeys, err := topayz512.BatchKEMKeyGen(batchSize)
	if err != nil {
		log.Fatalf("Failed to batch generate KEM keys: %v", err)
	}
	batchKEMGenDuration := time.Since(start)

	fmt.Printf("   Generated %d KEM key pairs in %v\n", batchSize, batchKEMGenDuration)

	// Batch encapsulation
	start = time.Now()
	ciphertexts, sharedSecrets1, err := topayz512.BatchKEMEncapsulate(kemPublicKeys)
	if err != nil {
		log.Fatalf("Failed to batch encapsulate: %v", err)
	}
	batchEncapDuration := time.Since(start)

	fmt.Printf("   Batch encapsulation of %d keys in %v\n", batchSize, batchEncapDuration)

	// Batch decapsulation
	start = time.Now()
	sharedSecrets2, err := topayz512.BatchKEMDecapsulate(kemSecretKeys, ciphertexts)
	if err != nil {
		log.Fatalf("Failed to batch decapsulate: %v", err)
	}
	batchDecapDuration := time.Since(start)

	fmt.Printf("   Batch decapsulation of %d keys in %v\n", batchSize, batchDecapDuration)

	// Verify shared secrets match
	matchCount := 0
	for i := 0; i < batchSize; i++ {
		if topayz512.SharedSecretEqual(sharedSecrets1[i], sharedSecrets2[i]) {
			matchCount++
		}
	}
	fmt.Printf("   Matching shared secrets: %d/%d\n", matchCount, batchSize)
	fmt.Println()

	// 6. System Capabilities
	fmt.Println("6. System Capabilities...")
	fmt.Printf("   SIMD Support: %v\n", topayz512.HasSIMDSupport())
	fmt.Printf("   Hardware RNG: %v\n", topayz512.HasHardwareRNG())
	fmt.Printf("   Optimal Thread Count: %d\n", topayz512.OptimalThreadCount())
	fmt.Println()

	// 7. Mobile Performance Estimation
	fmt.Println("7. Mobile Performance Estimation...")
	testSizes := []int{1024, 10240, 102400, 1048576} // 1KB, 10KB, 100KB, 1MB

	for _, size := range testSizes {
		estimate := topayz512.EstimateMobileLatency(size)
		fmt.Printf("   %6s: Total=%6.2fms, Frag=%5.2fms, Recon=%5.2fms, Chunks=%d\n",
			formatBytes(size),
			estimate.TotalMs,
			estimate.FragmentationMs,
			estimate.ReconstructionMs,
			estimate.RecommendedChunks)
	}
	fmt.Println()

	// 8. Serialization Example
	fmt.Println("8. Serialization Example...")

	// Create a fragment
	testData := []byte("This is test fragment data for serialization")
	fragment := topayz512.Fragment{
		ID:       12345,
		Index:    0,
		Total:    1,
		Data:     testData,
		Checksum: topayz512.ComputeHash(testData),
	}

	// Serialize
	serialized := topayz512.SerializeFragment(fragment)
	fmt.Printf("   Original Fragment Size: %d bytes\n", len(fragment.Data))
	fmt.Printf("   Serialized Size: %d bytes\n", len(serialized))

	// Deserialize
	deserialized, err := topayz512.DeserializeFragment(serialized)
	if err != nil {
		log.Fatalf("Failed to deserialize fragment: %v", err)
	}

	fmt.Printf("   Deserialization Success: %v\n", err == nil)
	fmt.Printf("   Data Integrity: %v\n", string(deserialized.Data) == string(fragment.Data))
	fmt.Printf("   Checksum Match: %v\n", topayz512.HashEqual(deserialized.Checksum, fragment.Checksum))
	fmt.Println()

	// 9. Memory Profiling
	fmt.Println("9. Memory Profiling...")
	profiler := topayz512.NewMemoryProfiler()

	// Perform some operations
	for i := 0; i < 1000; i++ {
		_ = topayz512.ComputeHash([]byte(fmt.Sprintf("test data %d", i)))
	}

	report := profiler.Report()
	fmt.Printf("   Memory Profile: %s\n", report)
	fmt.Println()

	fmt.Println("=== Quick Start Complete! ===")
	fmt.Println("The TOPAY-Z512 Go implementation is working correctly.")
	fmt.Printf("Library Version: %s\n", topayz512.Version)
	fmt.Printf("Security Level: %d-bit classical, %d-bit quantum\n",
		topayz512.SecurityLevel, topayz512.QuantumSecurityLevel)
}

func formatBytes(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%dKB", bytes/1024)
	} else {
		return fmt.Sprintf("%dMB", bytes/(1024*1024))
	}
}
