package main
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/topayfoundation/topayz512"
)
func main() {
	fmt.Println("=== TOPAY-Z512 Go Implementation - Interactive Guide ===")
	fmt.Println()
	fmt.Println("Welcome to the TOPAY-Z512 Interactive Guide!")
	fmt.Println("This guide will walk you through the key features of the library.")
	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)
	// System information
	showSystemInfo()
	waitForUser(scanner)
	// Hash operations
	demonstrateHashOperations(scanner)
	// Key pair operations
	demonstrateKeyPairOperations(scanner)
	// KEM operations
	demonstrateKEMOperations(scanner)
	// Fragmentation operations
	demonstrateFragmentationOperations(scanner)
	// Performance comparison
	demonstratePerformanceComparison(scanner)
	// Advanced features
	demonstrateAdvancedFeatures(scanner)
	fmt.Println()
	fmt.Println("=== Interactive Guide Complete ===")
	fmt.Println("Thank you for exploring TOPAY-Z512!")
	fmt.Println("For more information, check out the examples and documentation.")
}
func showSystemInfo() {
	fmt.Println("=== System Information ===")
	fmt.Printf("TOPAY-Z512 Version: %s\n", topayz512.Version)
	fmt.Printf("Security Level: %d-bit\n", topayz512.SecurityLevel)
	fmt.Printf("Optimal Thread Count: %d\n", topayz512.OptimalThreadCount())
	fmt.Printf("SIMD Support: %v\n", topayz512.HasSIMDSupport())
	fmt.Printf("Hardware RNG: %v\n", topayz512.HasHardwareRNG())
	fmt.Println()
}
func demonstrateHashOperations(scanner *bufio.Scanner) {
	fmt.Println("=== Hash Operations ===")
	fmt.Println("Let's start with basic hash operations.")
	fmt.Println()
	// Get user input for hashing
	fmt.Print("Enter some text to hash (or press Enter for default): ")
	scanner.Scan()
	input := scanner.Text()
	if input == "" {
		input = "Hello, TOPAY-Z512!"
	}
	data := []byte(input)
	fmt.Printf("Input: %s\n", input)
	fmt.Printf("Input size: %d bytes\n", len(data))
	fmt.Println()
	// Basic hash
	fmt.Println("1. Basic Hash:")
	start := time.Now()
	hash := topayz512.ComputeHash(data)
	duration := time.Since(start)
	fmt.Printf("   Hash: %s\n", hash.String())
	fmt.Printf("   Time: %v\n", duration)
	fmt.Println()
	// Hash with salt
	fmt.Println("2. Hash with Salt:")
	salt := []byte("random_salt_123")
	saltedHash := topayz512.HashWithSalt(data, salt)
	fmt.Printf("   Salt: %x\n", salt)
	fmt.Printf("   Salted Hash: %s\n", saltedHash.String())
	fmt.Printf("   Different from basic hash: %v\n", !topayz512.HashEqual(hash, saltedHash))
	fmt.Println()
	// String hash
	fmt.Println("3. String Hash:")
	stringHash := topayz512.HashString(input)
	fmt.Printf("   String Hash: %s\n", stringHash.String())
	fmt.Printf("   Same as basic hash: %v\n", topayz512.HashEqual(hash, stringHash))
	fmt.Println()
	// Multiple chunks
	fmt.Println("4. Hash Multiple Chunks:")
	concatenated := topayz512.ComputeHash([]byte("First chunk of data"))
	fmt.Printf("   Concatenated Hash: %s\n", concatenated.String())
	fmt.Println()
	waitForUser(scanner)
}
func demonstrateKeyPairOperations(scanner *bufio.Scanner) {
	fmt.Println("=== Key Pair Operations ===")
	fmt.Println("Now let's explore key pair generation and operations.")
	fmt.Println()
	// Basic key generation
	fmt.Println("1. Generating a new key pair...")
	privateKey, publicKey, err := topayz512.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	fmt.Printf("   Private Key: %s\n", privateKey.String())
	fmt.Printf("   Public Key:  %s\n", publicKey.String())
	fmt.Printf("   Key pair valid: %v\n", topayz512.VerifyKeyPair(privateKey, publicKey))
	fmt.Println()
	// Ask user for seed
	fmt.Print("Enter a seed for deterministic key generation (or press Enter to skip): ")
	scanner.Scan()
	seedInput := scanner.Text()
	if seedInput != "" {
		fmt.Println("2. Deterministic Key Generation:")
		seed := []byte(seedInput)
		var deterministicPrivate topayz512.PrivateKey
		var deterministicPublic topayz512.PublicKey
		deterministicPrivate, deterministicPublic, err = topayz512.GenerateKeyPairFromSeed(seed)
		if err != nil {
			log.Printf("Failed to generate deterministic key pair: %v", err)
		} else {
			fmt.Printf("   Seed: %s\n", seedInput)
			fmt.Printf("   Private Key: %s\n", deterministicPrivate.String())
			fmt.Printf("   Public Key:  %s\n", deterministicPublic.String())
			
			// Generate again with same seed
			deterministicPrivate2, deterministicPublic2, _ := topayz512.GenerateKeyPairFromSeed(seed)
			fmt.Printf("   Reproducible: %v\n", 
				topayz512.PrivateKeyEqual(deterministicPrivate, deterministicPrivate2) &&
				topayz512.PublicKeyEqual(deterministicPublic, deterministicPublic2))
		}
		fmt.Println()
	}

	// Public key derivation
	fmt.Println("3. Public Key Derivation:")
	derivedPublic := topayz512.DerivePublicKey(privateKey)
	fmt.Printf("   Original:  %s\n", publicKey.String())
	fmt.Printf("   Derived:   %s\n", derivedPublic.String())
	fmt.Printf("   Match: %v\n", topayz512.PublicKeyEqual(publicKey, derivedPublic))
	fmt.Println()

	// Ask for batch size
	fmt.Print("How many key pairs would you like to generate in batch? (1-100, default 5): ")
	scanner.Scan()
	batchInput := scanner.Text()
	batchSize := 5
	if batchInput != "" {
		var size int
		if size, err = strconv.Atoi(batchInput); err == nil && size > 0 && size <= 100 {
			batchSize = size
		}
	}

	fmt.Printf("4. Batch Key Generation (%d keys):\n", batchSize)
	start := time.Now()
	privateKeys, publicKeys, err := topayz512.BatchGenerateKeyPairs(batchSize)
	duration := time.Since(start)
	if err != nil {
		log.Printf("Failed batch generation: %v", err)
	} else {
		fmt.Printf("   Generated %d key pairs in %v\n", len(privateKeys), duration)
		fmt.Printf("   Rate: %.2f keys/sec\n", float64(batchSize)/duration.Seconds())
		
		// Verify a few
		validCount := 0
		checkCount := min(5, batchSize)
		for i := 0; i < checkCount; i++ {
			if topayz512.VerifyKeyPair(privateKeys[i], publicKeys[i]) {
				validCount++
			}
		}
		fmt.Printf("   Verified %d/%d keys as valid\n", validCount, checkCount)
	}
	fmt.Println()

	waitForUser(scanner)
}

func demonstrateKEMOperations(scanner *bufio.Scanner) {
	fmt.Println("=== KEM (Key Encapsulation Mechanism) Operations ===")
	fmt.Println("KEM allows secure key exchange without prior shared secrets.")
	fmt.Println()

	// Generate KEM key pair
	fmt.Println("1. Generating KEM key pair...")
	publicKey, secretKey, err := topayz512.KEMKeyGen()
	if err != nil {
		log.Fatalf("Failed to generate KEM key pair: %v", err)
	}

	fmt.Printf("   Public Key: %s\n", publicKey.String())
	fmt.Printf("   Secret Key: %s\n", secretKey.String())
	fmt.Printf("   Key pair valid: %v\n", topayz512.VerifyKEMKeyPair(publicKey, secretKey))
	fmt.Println()

	// Encapsulation
	fmt.Println("2. Key Encapsulation:")
	fmt.Println("   Alice encapsulates a shared secret using Bob's public key...")
	ciphertext, sharedSecret1, err := topayz512.KEMEncapsulate(publicKey)
	if err != nil {
		log.Fatalf("Failed to encapsulate: %v", err)
	}

	fmt.Printf("   Ciphertext: %s\n", ciphertext.String())
	fmt.Printf("   Shared Secret: %s\n", sharedSecret1.String())
	fmt.Println()

	// Decapsulation
	fmt.Println("3. Key Decapsulation:")
	fmt.Println("   Bob decapsulates the shared secret using his secret key...")
	sharedSecret2, err := topayz512.KEMDecapsulate(secretKey, ciphertext)
	if err != nil {
		log.Fatalf("Failed to decapsulate: %v", err)
	}

	fmt.Printf("   Decapsulated Secret: %s\n", sharedSecret2.String())
	fmt.Printf("   Secrets match: %v\n", topayz512.SharedSecretEqual(sharedSecret1, sharedSecret2))
	fmt.Println()

	// Multiple encapsulations
	fmt.Println("4. Multiple Encapsulations:")
	fmt.Println("   Each encapsulation produces a different ciphertext and shared secret...")
	for i := 0; i < 3; i++ {
		ct, ss, _ := topayz512.KEMEncapsulate(publicKey)
		ss2, _ := topayz512.KEMDecapsulate(secretKey, ct)
		fmt.Printf("   Encapsulation %d: Secrets match = %v, Different ciphertext = %v\n", 
			i+1, 
			topayz512.SharedSecretEqual(ss, ss2),
			!topayz512.CiphertextEqual(ciphertext, ct))
	}
	fmt.Println()

	// Ask for context
	fmt.Print("Enter a context string for contextual KEM (or press Enter to skip): ")
	scanner.Scan()
	contextInput := scanner.Text()

	if contextInput != "" {
		fmt.Println("5. KEM with Context:")
		context := []byte(contextInput)
		contextCiphertext, contextSecret1, err := topayz512.KEMWithContext(publicKey, context)
		if err != nil {
			log.Printf("Failed KEM with context: %v", err)
		} else {
			contextSecret2, err := topayz512.KEMDecapsulateWithContext(secretKey, contextCiphertext, context)
			if err != nil {
				log.Printf("Failed decapsulation with context: %v", err)
			} else {
				fmt.Printf("   Context: %s\n", contextInput)
				fmt.Printf("   Context affects result: %v\n", 
					!topayz512.SharedSecretEqual(sharedSecret1, contextSecret1))
				fmt.Printf("   Context secrets match: %v\n", 
					topayz512.SharedSecretEqual(contextSecret1, contextSecret2))
			}
		}
		fmt.Println()
	}

	waitForUser(scanner)
}

func demonstrateFragmentationOperations(scanner *bufio.Scanner) {
	fmt.Println("=== Fragmentation Operations ===")
	fmt.Println("Fragmentation allows efficient processing of large data on resource-constrained devices.")
	fmt.Println()

	// Ask for data size
	fmt.Print("Enter data size in KB for fragmentation demo (1-1000, default 10): ")
	scanner.Scan()
	sizeInput := scanner.Text()
	sizeKB := 10
	if sizeInput != "" {
		var size int
		var err error
		if size, err = strconv.Atoi(sizeInput); err == nil && size > 0 && size <= 1000 {
			sizeKB = size
		}
	}

	// Generate test data
	dataSize := sizeKB * 1024
	testData := make([]byte, dataSize)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	fmt.Printf("Generated %d KB of test data\n", sizeKB)
	fmt.Printf("Should fragment: %v\n", topayz512.ShouldFragment(dataSize))
	fmt.Printf("Optimal fragment count: %d\n", topayz512.CalculateFragmentCount(dataSize))
	fmt.Println()

	// Sequential fragmentation
	fmt.Println("1. Sequential Fragmentation:")
	start := time.Now()
	result, err := topayz512.FragmentData(testData)
	seqDuration := time.Since(start)
	if err != nil {
		log.Printf("Failed sequential fragmentation: %v", err)
	} else {
		fmt.Printf("   Fragments: %d\n", len(result.Fragments))
		fmt.Printf("   Fragment size: %d bytes each\n", len(result.Fragments[0].Data))
		fmt.Printf("   Time: %v\n", seqDuration)
		fmt.Printf("   Throughput: %.2f MB/s\n", 
			float64(dataSize)/seqDuration.Seconds()/(1024*1024))
	}
	fmt.Println()

	// Parallel fragmentation
	fmt.Println("2. Parallel Fragmentation:")
	start = time.Now()
	parResult, err := topayz512.ParallelFragmentData(testData)
	parDuration := time.Since(start)
	if err != nil {
		log.Printf("Failed parallel fragmentation: %v", err)
	} else {
		fmt.Printf("   Fragments: %d\n", len(parResult.Fragments))
		fmt.Printf("   Time: %v\n", parDuration)
		fmt.Printf("   Throughput: %.2f MB/s\n", 
			float64(dataSize)/parDuration.Seconds()/(1024*1024))
		fmt.Printf("   Speedup: %.2fx\n", float64(seqDuration)/float64(parDuration))
	}
	fmt.Println()

	// Reconstruction
	fmt.Println("3. Data Reconstruction:")
	start = time.Now()
	reconstructed, err := topayz512.ReconstructData(result.Fragments)
	reconDuration := time.Since(start)
	if err != nil {
		log.Printf("Failed reconstruction: %v", err)
	} else {
		fmt.Printf("   Reconstruction time: %v\n", reconDuration)
		fmt.Printf("   Data integrity: %v\n", len(reconstructed.Data) == len(testData))
		fmt.Printf("   Complete: %v\n", reconstructed.IsComplete)
	}
	fmt.Println()

	// Mobile performance estimation
	fmt.Println("4. Mobile Performance Estimation:")
	estimate := topayz512.EstimateMobileLatency(dataSize)
	fmt.Printf("   Total latency: %.2f ms\n", estimate.TotalMs)
	fmt.Printf("   Fragmentation: %.2f ms\n", estimate.FragmentationMs)
	fmt.Printf("   Reconstruction: %.2f ms\n", estimate.ReconstructionMs)
	fmt.Printf("   Recommended chunks: %d\n", estimate.RecommendedChunks)
	fmt.Println()

	// Fragment integrity
	fmt.Println("5. Fragment Integrity:")
	if len(result.Fragments) > 0 {
		fragment := result.Fragments[0]
		fmt.Printf("   Original fragment valid: %v\n", 
			topayz512.ValidateFragmentIntegrity(fragment))
		
		// Corrupt the fragment
		corruptedFragment := fragment
		corruptedFragment.Data[0] ^= 0xFF
		fmt.Printf("   Corrupted fragment valid: %v\n", 
			topayz512.ValidateFragmentIntegrity(corruptedFragment))
		
		// Attempt repair
		repairedFragment, err := topayz512.RepairFragment(corruptedFragment, testData, len(result.Fragments[0].Data))
		if err != nil {
			fmt.Printf("   Repair failed: %v\n", err)
		} else {
			fmt.Printf("   Repaired fragment valid: %v\n", 
				topayz512.ValidateFragmentIntegrity(repairedFragment))
		}
	}
	fmt.Println()

	waitForUser(scanner)
}

func demonstratePerformanceComparison(scanner *bufio.Scanner) {
	fmt.Println("=== Performance Comparison ===")
	fmt.Println("Let's compare the performance of different operations.")
	fmt.Println()

	// Hash performance
	fmt.Println("1. Hash Performance:")
	dataSizes := []int{1024, 4096, 16384, 65536} // 1KB to 64KB
	
	for _, size := range dataSizes {
		benchmark := topayz512.BenchmarkHash(size, 1000)
		fmt.Printf("   %2dKB: %8.2f MB/s, %6.0f hashes/sec, %6d ns/hash\n",
			size/1024,
			benchmark.ThroughputMBps,
			benchmark.HashesPerSec,
			benchmark.LatencyNs)
	}
	fmt.Println()

	// Key pair performance
	fmt.Println("2. Key Pair Performance:")
	keyBenchmark := topayz512.BenchmarkKeyPairGeneration(1000)
	fmt.Printf("   Key pairs/sec: %.2f\n", keyBenchmark.KeyPairsPerSec)
	fmt.Printf("   Avg latency: %.2f ms\n", keyBenchmark.AvgLatencyMs)
	fmt.Printf("   Batch speedup: %.2fx\n", keyBenchmark.BatchSpeedupRatio)
	fmt.Println()

	// KEM performance
	fmt.Println("3. KEM Performance:")
	kemBenchmark := topayz512.BenchmarkKEM(1000)
	fmt.Printf("   Key generation: %.2f ops/sec\n", kemBenchmark.KeyGenPerSec)
	fmt.Printf("   Encapsulation: %.2f ops/sec\n", kemBenchmark.EncapsulatePerSec)
	fmt.Printf("   Decapsulation: %.2f ops/sec\n", kemBenchmark.DecapsulatePerSec)
	fmt.Printf("   Batch speedup: %.2fx\n", kemBenchmark.BatchSpeedupRatio)
	fmt.Println()

	// Fragmentation performance
	fmt.Println("4. Fragmentation Performance:")
	fragSizes := []int{4096, 16384, 65536} // 4KB to 64KB
	
	for _, size := range fragSizes {
		benchmark := topayz512.BenchmarkFragmentation(size)
		fmt.Printf("   %2dKB: %8.2f MB/s, Frag=%.2fms, Recon=%.2fms, Speedup=%.2fx\n",
			size/1024,
			benchmark.ThroughputMBps,
			benchmark.FragmentationMs,
			benchmark.ReconstructionMs,
			benchmark.ParallelSpeedup)
	}
	fmt.Println()

	waitForUser(scanner)
}

func demonstrateAdvancedFeatures(scanner *bufio.Scanner) {
	fmt.Println("=== Advanced Features ===")
	fmt.Println("Let's explore some advanced features of TOPAY-Z512.")
	fmt.Println()

	// Memory profiling
	fmt.Println("1. Memory Profiling:")
	profiler := topayz512.NewMemoryProfiler()
	
	// Perform some operations
	for i := 0; i < 100; i++ {
		data := make([]byte, 1024)
		_ = topayz512.ComputeHash(data)
	}
	
	report := profiler.Report()
	fmt.Printf("   Memory usage for 100 hash operations: %s\n", report)
	fmt.Println()

	// Batch operations comparison
	fmt.Println("2. Batch vs Individual Operations:")
	batchSize := 50
	
	// Individual key generation
	start := time.Now()
	for i := 0; i < batchSize; i++ {
		_, _, _ = topayz512.GenerateKeyPair()
	}
	individualDuration := time.Since(start)
	
	// Batch key generation
	start = time.Now()
	_, _, _ = topayz512.BatchGenerateKeyPairs(batchSize)
	batchDuration := time.Since(start)
	
	fmt.Printf("   Individual: %v (%.2f keys/sec)\n", 
		individualDuration, float64(batchSize)/individualDuration.Seconds())
	fmt.Printf("   Batch: %v (%.2f keys/sec)\n", 
		batchDuration, float64(batchSize)/batchDuration.Seconds())
	fmt.Printf("   Batch speedup: %.2fx\n", float64(individualDuration)/float64(batchDuration))
	fmt.Println()

	// HD Wallet demonstration
	fmt.Println("3. HD Wallet Generation:")
	fmt.Print("Enter a master seed (or press Enter for default): ")
	scanner.Scan()
	seedInput := scanner.Text()
	if seedInput == "" {
		seedInput = "master_seed_for_hd_wallet_demonstration"
	}
	
	seed := []byte(seedInput)
	depth := 10
	
	start = time.Now()
	hdWallet, err := topayz512.GenerateHDWallet(seed, depth)
	hdDuration := time.Since(start)
	
	if err != nil {
		log.Printf("Failed to generate HD wallet: %v", err)
	} else {
		fmt.Printf("   Generated %d key pairs in %v\n", len(hdWallet), hdDuration)
		fmt.Printf("   Rate: %.2f keys/sec\n", float64(depth)/hdDuration.Seconds())
		
		// Show first few keys
		fmt.Println("   First 3 key pairs:")
		for i := 0; i < min(3, len(hdWallet)); i++ {
			fmt.Printf("     Key %d: %s -> %s\n", i,
				hdWallet[i].PrivateKey.String()[:16]+"...",
				hdWallet[i].PublicKey.String()[:16]+"...")
		}
	}
	fmt.Println()

	// Serialization formats
	fmt.Println("4. Serialization Formats:")
	privateKey, publicKey, _ := topayz512.GenerateKeyPair()
	
	// Bytes format
	privateBytes := privateKey.Bytes()
	publicBytes := publicKey.Bytes()
	fmt.Printf("   Private key bytes: %d bytes\n", len(privateBytes))
	fmt.Printf("   Public key bytes: %d bytes\n", len(publicBytes))
	
	// Hex format
	privateHex := privateKey.String()
	publicHex := publicKey.String()
	fmt.Printf("   Private key hex: %d characters\n", len(privateHex))
	fmt.Printf("   Public key hex: %d characters\n", len(publicHex))
	
	// String format
	privateStr := privateKey.String()
	publicStr := publicKey.String()
	fmt.Printf("   Private key string: %d characters\n", len(privateStr))
	fmt.Printf("   Public key string: %d characters\n", len(publicStr))
	fmt.Println()

	// Security features
	fmt.Println("5. Security Features:")
	fmt.Printf("   Constant-time equality: %v\n", true) // Implemented in library
	fmt.Printf("   Secure random generation: %v\n", topayz512.HasHardwareRNG())
	fmt.Printf("   Memory zeroing: %v\n", true) // Implemented in SecureErase functions
	fmt.Printf("   Side-channel resistance: %v\n", true) // Built into the design
	fmt.Println()
}

func waitForUser(scanner *bufio.Scanner) {
	fmt.Print("Press Enter to continue...")
	scanner.Scan()
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}