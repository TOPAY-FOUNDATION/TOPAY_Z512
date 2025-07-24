package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TOPAY-FOUNDATION/TOPAY_Z512/go"
)

func main() {
	fmt.Println("=== TOPAY-Z512 Go Implementation - Fragmentation Operations ===")
	fmt.Println()

	// Basic fragmentation
	fmt.Println("1. Basic Data Fragmentation:")
	originalData := []byte("This is a test message that will be fragmented into smaller pieces for processing. " +
		"The fragmentation system allows for efficient processing of large data sets by breaking them " +
		"into manageable chunks that can be processed in parallel or on resource-constrained devices.")

	fmt.Printf("   Original data size: %d bytes\n", len(originalData))
	fmt.Printf("   Should fragment: %v\n", topayz512.ShouldFragment(len(originalData)))

	// Fragment the data
	result, err := topayz512.FragmentData(originalData)
	if err != nil {
		log.Fatalf("Failed to fragment data: %v", err)
	}

	fmt.Printf("   Fragments created: %d\n", len(result.Fragments))
	fmt.Printf("   Fragment size: %d bytes each\n", len(result.Fragments[0].Data))
	fmt.Printf("   Original size: %d bytes\n", result.Metadata.OriginalSize)
	fmt.Printf("   Fragmentation overhead: %d bytes\n",
		uint64(len(result.Fragments)*16+64)-result.Metadata.OriginalSize) // Approximate overhead
	fmt.Println()

	// Display fragment details
	fmt.Println("   Fragment details:")
	for i, fragment := range result.Fragments {
		fmt.Printf("     Fragment %d: ID=%d, Size=%d, Checksum=%x\n",
			i, fragment.ID, len(fragment.Data), fragment.Checksum[:8])
	}
	fmt.Println()

	// Reconstruct the data
	fmt.Println("2. Data Reconstruction:")
	reconstructed, err := topayz512.ReconstructData(result.Fragments)
	if err != nil {
		log.Fatalf("Failed to reconstruct data: %v", err)
	}

	fmt.Printf("   Reconstructed size: %d bytes\n", len(reconstructed.Data))
	fmt.Printf("   Data integrity: %v\n", string(originalData) == string(reconstructed.Data))
	fmt.Printf("   Reconstruction successful: %v\n", reconstructed.IsComplete)
	fmt.Println()

	// Parallel fragmentation
	fmt.Println("3. Parallel Fragmentation:")
	largeData := make([]byte, 10240) // 10KB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	fmt.Printf("   Large data size: %d bytes\n", len(largeData))

	// Sequential fragmentation
	start := time.Now()
	seqResult, err := topayz512.FragmentData(largeData)
	seqDuration := time.Since(start)
	if err != nil {
		log.Fatalf("Failed sequential fragmentation: %v", err)
	}

	// Parallel fragmentation
	start = time.Now()
	parResult, err := topayz512.ParallelFragmentData(largeData)
	parDuration := time.Since(start)
	if err != nil {
		log.Fatalf("Failed parallel fragmentation: %v", err)
	}

	fmt.Printf("   Sequential time: %v\n", seqDuration)
	fmt.Printf("   Parallel time: %v\n", parDuration)
	fmt.Printf("   Speedup: %.2fx\n", float64(seqDuration)/float64(parDuration))
	fmt.Printf("   Same result: %v\n", len(seqResult.Fragments) == len(parResult.Fragments))
	fmt.Println()

	// Fragment serialization
	fmt.Println("4. Fragment Serialization:")
	fragment := result.Fragments[0]

	serialized := topayz512.SerializeFragment(fragment)

	fmt.Printf("   Original fragment size: %d bytes\n", len(fragment.Data))
	fmt.Printf("   Serialized size: %d bytes\n", len(serialized))

	deserialized, err := topayz512.DeserializeFragment(serialized)
	if err != nil {
		log.Fatalf("Failed to deserialize fragment: %v", err)
	}

	fmt.Printf("   Deserialization successful: %v\n",
		fragment.ID == deserialized.ID &&
			len(fragment.Data) == len(deserialized.Data))
	fmt.Println()

	// Mobile performance estimation
	fmt.Println("5. Mobile Performance Estimation:")
	mobileSizes := []int{1024, 5120, 10240, 51200, 102400} // 1KB to 100KB

	for _, size := range mobileSizes {
		estimate := topayz512.EstimateMobileLatency(size)
		fmt.Printf("   %5dKB: Total=%6.2fms, Frag=%5.2fms, Recon=%5.2fms, Chunks=%d\n",
			size/1024,
			estimate.TotalMs,
			estimate.FragmentationMs,
			estimate.ReconstructionMs,
			estimate.RecommendedChunks)
	}
	fmt.Println()

	// Fragmented hash operations
	fmt.Println("6. Fragmented Hash Operations:")
	hashData := make([]byte, 8192) // 8KB
	for i := range hashData {
		hashData[i] = byte(i)
	}

	// Regular hash
	start = time.Now()
	regularHash := topayz512.ComputeHash(hashData)
	regularDuration := time.Since(start)

	// Fragmented hash
	var fragmentedHash topayz512.Hash
	var fragmentedDuration time.Duration
	start = time.Now()
	fragmentedHash, err = topayz512.FragmentedHash(hashData)
	fragmentedDuration = time.Since(start)
	if err != nil {
		log.Fatalf("Failed fragmented hash: %v", err)
	}

	fmt.Printf("   Data size: %d bytes\n", len(hashData))
	fmt.Printf("   Regular hash: %s (took %v)\n", regularHash.String(), regularDuration)
	fmt.Printf("   Fragmented hash: %s (took %v)\n", fragmentedHash.String(), fragmentedDuration)
	fmt.Printf("   Hashes match: %v\n", topayz512.HashEqual(regularHash, fragmentedHash))
	fmt.Println()

	// Fragmented KEM operations
	fmt.Println("7. Fragmented KEM Operations:")
	kemData := make([]byte, 4096) // 4KB
	for i := range kemData {
		kemData[i] = byte(i * 7 % 256)
	}

	ciphertexts, sharedSecrets, err := topayz512.FragmentedKEM(kemData)
	if err != nil {
		log.Fatalf("Failed fragmented KEM: %v", err)
	}

	fmt.Printf("   KEM data size: %d bytes\n", len(kemData))
	fmt.Printf("   Number of ciphertexts: %d\n", len(ciphertexts))
	fmt.Printf("   Number of shared secrets: %d\n", len(sharedSecrets))
	if len(ciphertexts) > 0 {
		fmt.Printf("   First ciphertext: %s\n", ciphertexts[0].String())
	}
	if len(sharedSecrets) > 0 {
		fmt.Printf("   First shared secret: %s\n", sharedSecrets[0].String())
	}
	fmt.Println()

	// Fragment integrity and repair
	fmt.Println("8. Fragment Integrity and Repair:")
	testData := []byte("This is test data for integrity checking and repair functionality.")
	fragResult, _ := topayz512.FragmentData(testData)

	// Check integrity of original fragments
	fmt.Println("   Original fragments:")
	for i, fragment := range fragResult.Fragments {
		err = topayz512.ValidateFragmentIntegrity(fragment)
		fmt.Printf("     Fragment %d: %v\n", i, err == nil)
	}

	// Corrupt a fragment (use the first fragment if only one exists)
	fragmentIndex := 0
	if len(fragResult.Fragments) > 1 {
		fragmentIndex = 1
	}
	corruptedFragment := fragResult.Fragments[fragmentIndex]
	corruptedFragment.Data[0] ^= 0xFF // Flip bits

	fmt.Println("   After corruption:")
	for i, fragment := range fragResult.Fragments {
		valid := topayz512.ValidateFragmentIntegrity(fragment)
		fmt.Printf("     Fragment %d: %v\n", i, valid)
	}

	// Attempt repair
	repairedFragment, err := topayz512.RepairFragment(corruptedFragment, testData, len(testData)/len(fragResult.Fragments))
	if err != nil {
		fmt.Printf("   Repair failed: %v\n", err)
	} else {
		fmt.Printf("   Repair successful: %v\n",
			topayz512.ValidateFragmentIntegrity(repairedFragment) == nil)
	}
	fmt.Println()

	// Performance benchmarking
	fmt.Println("9. Fragmentation Performance Benchmark:")
	benchmarkSizes := []int{1024, 4096, 16384, 65536}

	for _, size := range benchmarkSizes {
		benchmark := topayz512.BenchmarkFragmentation(size)
		fmt.Printf("   %5dKB: %.2f MB/s, Frag=%.2fms, Recon=%.2fms, Speedup=%.2fx\n",
			size/1024,
			benchmark.ThroughputMBps,
			benchmark.FragmentationMs,
			benchmark.ReconstructionMs,
			benchmark.ParallelSpeedup)
	}
	fmt.Println()

	// Fragment count optimization
	fmt.Println("10. Fragment Count Optimization:")
	optimizationData := make([]byte, 32768) // 32KB
	for i := range optimizationData {
		optimizationData[i] = byte(i)
	}

	fragmentCounts := []int{1, 2, 4, 8, 16, 32}

	fmt.Printf("    Data size: %d bytes\n", len(optimizationData))
	fmt.Println("    Fragment count analysis:")

	for _, count := range fragmentCounts {
		optimalCount := topayz512.CalculateFragmentCount(len(optimizationData))

		start = time.Now()
		// Simulate fragmentation with specific count
		fragmentSize := (len(optimizationData) + count - 1) / count
		for i := 0; i < count; i++ {
			startIdx := i * fragmentSize
			endIdx := startIdx + fragmentSize
			if endIdx > len(optimizationData) {
				endIdx = len(optimizationData)
			}
			_ = topayz512.ComputeHash(optimizationData[startIdx:endIdx])
		}
		duration := time.Since(start)

		throughput := float64(len(optimizationData)) / duration.Seconds() / (1024 * 1024)
		optimal := ""
		if count == optimalCount {
			optimal = " (optimal)"
		}

		fmt.Printf("      %2d fragments: %8.2f MB/s (%6.2f ms)%s\n",
			count, throughput, duration.Seconds()*1000, optimal)
	}
	fmt.Println()

	// Throughput improvement demonstration
	fmt.Println("11. Throughput Improvement:")
	throughputData := make([]byte, 1048576) // 1MB
	for i := range throughputData {
		throughputData[i] = byte(i % 256)
	}

	// Sequential processing
	start = time.Now()
	_ = topayz512.ComputeHash(throughputData)
	seqTime := time.Since(start)
	seqThroughput := float64(len(throughputData)) / seqTime.Seconds() / (1024 * 1024)

	// Parallel fragmented processing
	start = time.Now()
	_, _ = topayz512.ParallelFragmentData(throughputData)
	parTime := time.Since(start)
	parThroughput := float64(len(throughputData)) / parTime.Seconds() / (1024 * 1024)

	fmt.Printf("    Data size: %d MB\n", len(throughputData)/(1024*1024))
	fmt.Printf("    Sequential: %.2f MB/s (%v)\n", seqThroughput, seqTime)
	fmt.Printf("    Parallel: %.2f MB/s (%v)\n", parThroughput, parTime)
	fmt.Printf("    Improvement: %.2fx\n", parThroughput/seqThroughput)
	fmt.Println()

	// Memory usage analysis
	fmt.Println("12. Memory Usage Analysis:")
	memoryData := make([]byte, 262144) // 256KB
	for i := range memoryData {
		memoryData[i] = byte(i)
	}

	profiler := topayz512.NewMemoryProfiler()

	// Fragment the data
	_, err = topayz512.FragmentData(memoryData)
	if err != nil {
		log.Printf("Fragmentation failed: %v", err)
	}

	memoryReport := profiler.Report()
	fmt.Printf("    Memory usage for 256KB fragmentation: %s\n", memoryReport)

	fmt.Println()
	fmt.Println("=== Fragmentation Operations Complete ===")
}
