package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/TOPAY-FOUNDATION/TOPAY_Z512/go"
)

func main() {
	fmt.Println("=== TOPAY-Z512 Go Implementation - Performance Benchmark ===")
	fmt.Println()

	// System Information
	fmt.Println("System Information:")
	fmt.Printf("  Go Version: %s\n", runtime.Version())
	fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("  Optimal Threads: %d\n", topayz512.OptimalThreadCount())
	fmt.Printf("  SIMD Support: %v\n", topayz512.HasSIMDSupport())
	fmt.Printf("  Hardware RNG: %v\n", topayz512.HasHardwareRNG())
	fmt.Println()

	// Hash Performance
	fmt.Println("=== Hash Performance ===")
	runHashBenchmarks()
	fmt.Println()

	// Key Pair Performance
	fmt.Println("=== Key Pair Performance ===")
	runKeyPairBenchmarks()
	fmt.Println()

	// KEM Performance
	fmt.Println("=== KEM Performance ===")
	runKEMBenchmarks()
	fmt.Println()

	// Fragmentation Performance
	fmt.Println("=== Fragmentation Performance ===")
	runFragmentationBenchmarks()
	fmt.Println()

	// Memory Usage Analysis
	fmt.Println("=== Memory Usage Analysis ===")
	runMemoryAnalysis()
	fmt.Println()

	// Scalability Test
	fmt.Println("=== Scalability Test ===")
	runScalabilityTest()
	fmt.Println()

	fmt.Println("=== Benchmark Complete ===")
}

func runHashBenchmarks() {
	dataSizes := []int{64, 256, 1024, 4096, 16384, 65536}
	iterations := 10000

	fmt.Println("Hash Throughput (different data sizes):")
	fmt.Printf("%-10s %-15s %-15s %-15s %-15s\n", "Size", "Throughput", "Latency", "Hashes/sec", "MB/s")
	fmt.Println(strings.Repeat("-", 75))

	for _, size := range dataSizes {
		benchmark := topayz512.BenchmarkHash(size, iterations)

		fmt.Printf("%-10s %-15.2f %-15d %-15.0f %-15.2f\n",
			formatBytes(size),
			benchmark.ThroughputMBps,
			benchmark.LatencyNs,
			benchmark.HashesPerSec,
			benchmark.ThroughputMBps)
	}

	// Batch hash performance
	fmt.Println("\nBatch Hash Performance:")
	batchSizes := []int{10, 50, 100, 500, 1000}
	dataSize := 1024

	for _, batchSize := range batchSizes {
		inputs := make([][]byte, batchSize)
		for i := range inputs {
			inputs[i] = make([]byte, dataSize)
			for j := range inputs[i] {
				inputs[i][j] = byte(i + j)
			}
		}

		start := time.Now()
		results := topayz512.BatchHash(inputs)
		duration := time.Since(start)

		throughput := float64(batchSize*dataSize) / duration.Seconds() / (1024 * 1024)
		hashesPerSec := float64(batchSize) / duration.Seconds()

		fmt.Printf("  Batch Size %4d: %8.2f MB/s, %8.0f hashes/sec, %d results\n",
			batchSize, throughput, hashesPerSec, len(results))
	}
}

func runKeyPairBenchmarks() {
	iterations := 1000

	benchmark := topayz512.BenchmarkKeyPairGeneration(iterations)

	fmt.Printf("Single Key Pair Generation:\n")
	fmt.Printf("  Key Pairs/sec: %.2f\n", benchmark.KeyPairsPerSec)
	fmt.Printf("  Avg Latency: %.2f ms\n", benchmark.AvgLatencyMs)
	fmt.Printf("  Batch Speedup: %.2fx\n", benchmark.BatchSpeedupRatio)

	// Test different batch sizes
	fmt.Println("\nBatch Key Pair Generation:")
	batchSizes := []int{10, 50, 100, 500, 1000}

	for _, batchSize := range batchSizes {
		start := time.Now()
		privateKeys, publicKeys, err := topayz512.BatchGenerateKeyPairs(batchSize)
		duration := time.Since(start)

		if err != nil {
			log.Printf("Failed batch generation for size %d: %v", batchSize, err)
			continue
		}

		keyPairsPerSec := float64(batchSize) / duration.Seconds()
		avgLatencyMs := duration.Seconds() * 1000 / float64(batchSize)

		fmt.Printf("  Batch Size %4d: %8.2f keys/sec, %6.2f ms avg, %d generated\n",
			batchSize, keyPairsPerSec, avgLatencyMs, len(privateKeys))

		// Verify a few random key pairs
		validCount := 0
		checkCount := min(10, batchSize)
		for i := 0; i < checkCount; i++ {
			if topayz512.VerifyKeyPair(privateKeys[i], publicKeys[i]) {
				validCount++
			}
		}
		fmt.Printf("    Verification: %d/%d valid\n", validCount, checkCount)
	}

	// HD Wallet performance
	fmt.Println("\nHD Wallet Generation:")
	seed := []byte("this is a test seed for HD wallet generation that is long enough")
	depths := []int{10, 50, 100, 256}

	for _, depth := range depths {
		start := time.Now()
		keyPairs, err := topayz512.GenerateHDWallet(seed, depth)
		duration := time.Since(start)

		if err != nil {
			log.Printf("Failed HD wallet generation for depth %d: %v", depth, err)
			continue
		}

		keysPerSec := float64(depth) / duration.Seconds()
		fmt.Printf("  Depth %3d: %8.2f keys/sec, %6.2f ms total, %d generated\n",
			depth, keysPerSec, duration.Seconds()*1000, len(keyPairs))
	}
}

func runKEMBenchmarks() {
	iterations := 1000

	benchmark := topayz512.BenchmarkKEM(iterations)

	fmt.Printf("KEM Operations (single):\n")
	fmt.Printf("  Key Generation: %.2f ops/sec, %.2f ms avg\n",
		benchmark.KeyGenPerSec, benchmark.AvgLatencyMs)
	fmt.Printf("  Encapsulation: %.2f ops/sec\n", benchmark.EncapsulatePerSec)
	fmt.Printf("  Decapsulation: %.2f ops/sec\n", benchmark.DecapsulatePerSec)
	fmt.Printf("  Batch Speedup: %.2fx\n", benchmark.BatchSpeedupRatio)

	// Batch KEM operations
	fmt.Println("\nBatch KEM Operations:")
	batchSizes := []int{10, 50, 100, 500, 1000}

	for _, batchSize := range batchSizes {
		// Key generation
		start := time.Now()
		publicKeys, secretKeys, err := topayz512.BatchKEMKeyGen(batchSize)
		keyGenDuration := time.Since(start)

		if err != nil {
			log.Printf("Failed batch KEM key generation for size %d: %v", batchSize, err)
			continue
		}

		// Encapsulation
		start = time.Now()
		ciphertexts, sharedSecrets1, err := topayz512.BatchKEMEncapsulate(publicKeys)
		encapDuration := time.Since(start)

		if err != nil {
			log.Printf("Failed batch encapsulation for size %d: %v", batchSize, err)
			continue
		}

		// Decapsulation
		start = time.Now()
		sharedSecrets2, err := topayz512.BatchKEMDecapsulate(secretKeys, ciphertexts)
		decapDuration := time.Since(start)

		if err != nil {
			log.Printf("Failed batch decapsulation for size %d: %v", batchSize, err)
			continue
		}

		keyGenPerSec := float64(batchSize) / keyGenDuration.Seconds()
		encapPerSec := float64(batchSize) / encapDuration.Seconds()
		decapPerSec := float64(batchSize) / decapDuration.Seconds()

		fmt.Printf("  Batch Size %4d:\n", batchSize)
		fmt.Printf("    KeyGen: %8.2f ops/sec (%6.2f ms)\n", keyGenPerSec, keyGenDuration.Seconds()*1000)
		fmt.Printf("    Encap:  %8.2f ops/sec (%6.2f ms)\n", encapPerSec, encapDuration.Seconds()*1000)
		fmt.Printf("    Decap:  %8.2f ops/sec (%6.2f ms)\n", decapPerSec, decapDuration.Seconds()*1000)

		// Verify shared secrets match
		matchCount := 0
		for i := 0; i < batchSize; i++ {
			if topayz512.SharedSecretEqual(sharedSecrets1[i], sharedSecrets2[i]) {
				matchCount++
			}
		}
		fmt.Printf("    Verification: %d/%d secrets match\n", matchCount, batchSize)
	}
}

func runFragmentationBenchmarks() {
	dataSizes := []int{1024, 4096, 16384, 65536, 262144, 1048576} // 1KB to 1MB

	fmt.Println("Fragmentation Performance:")
	fmt.Printf("%-10s %-12s %-12s %-12s %-12s %-12s\n",
		"Size", "Throughput", "Frag(ms)", "Recon(ms)", "Speedup", "Mobile(ms)")
	fmt.Println(strings.Repeat("-", 75))

	for _, size := range dataSizes {
		benchmark := topayz512.BenchmarkFragmentation(size)

		fmt.Printf("%-10s %-12.2f %-12.2f %-12.2f %-12.2fx %-12.2f\n",
			formatBytes(size),
			benchmark.ThroughputMBps,
			benchmark.FragmentationMs,
			benchmark.ReconstructionMs,
			benchmark.ParallelSpeedup,
			benchmark.MobileLatencyMs)
	}

	// Fragment count optimization
	fmt.Println("\nFragment Count Analysis:")
	testSize := 65536 // 64KB
	testData := make([]byte, testSize)
	for i := range testData {
		testData[i] = byte(i)
	}

	fragmentCounts := []int{1, 2, 4, 8, 16, 32, 64}

	for _, count := range fragmentCounts {
		// Manually fragment with specific count
		fragmentSize := (testSize + count - 1) / count

		start := time.Now()

		// Simulate fragmentation with specific count
		for i := 0; i < count; i++ {
			startIdx := i * fragmentSize
			endIdx := startIdx + fragmentSize
			if endIdx > testSize {
				endIdx = testSize
			}
			_ = topayz512.ComputeHash(testData[startIdx:endIdx])
		}

		duration := time.Since(start)
		throughput := float64(testSize) / duration.Seconds() / (1024 * 1024)

		fmt.Printf("  %2d fragments: %8.2f MB/s (%6.2f ms)\n",
			count, throughput, duration.Seconds()*1000)
	}

	// Mobile latency estimates
	fmt.Println("\nMobile Performance Estimates:")
	mobileSizes := []int{1024, 10240, 102400, 1048576} // 1KB, 10KB, 100KB, 1MB

	for _, size := range mobileSizes {
		estimate := topayz512.EstimateMobileLatency(size)
		fmt.Printf("  %6s: Total=%6.2fms, Frag=%5.2fms, Recon=%5.2fms, Chunks=%d\n",
			formatBytes(size),
			estimate.TotalMs,
			estimate.FragmentationMs,
			estimate.ReconstructionMs,
			estimate.RecommendedChunks)
	}
}

func runMemoryAnalysis() {
	profiler := topayz512.NewMemoryProfiler()

	// Test memory usage for different operations
	operations := []struct {
		name string
		fn   func()
	}{
		{"1000 Hashes", func() {
			data := make([]byte, 1024)
			for i := 0; i < 1000; i++ {
				_ = topayz512.ComputeHash(data)
			}
		}},
		{"100 Key Pairs", func() {
			for i := 0; i < 100; i++ {
				_, _, _ = topayz512.GenerateKeyPair()
			}
		}},
		{"100 KEM Operations", func() {
			for i := 0; i < 100; i++ {
				pub, sec, _ := topayz512.KEMKeyGen()
				ct, _, _ := topayz512.KEMEncapsulate(pub)
				_, _ = topayz512.KEMDecapsulate(sec, ct)
			}
		}},
		{"Large Fragmentation", func() {
			data := make([]byte, 1048576) // 1MB
			_, _ = topayz512.FragmentData(data)
		}},
	}

	for _, op := range operations {
		_ = topayz512.NewMemoryProfiler()
		op.fn()
		report := profiler.Report()
		fmt.Printf("  %-20s: %s\n", op.name, report)
	}

	// Overall memory usage
	fmt.Printf("\nOverall Memory Usage: %s\n", profiler.Report())
}

func runScalabilityTest() {
	// Test how performance scales with different numbers of goroutines
	testData := make([]byte, 4096)
	for i := range testData {
		testData[i] = byte(i)
	}

	threadCounts := []int{1, 2, 4, 8, 16, 32}
	iterations := 1000

	fmt.Println("Thread Scalability (Hash Operations):")
	fmt.Printf("%-10s %-15s %-15s %-15s\n", "Threads", "Duration", "Ops/sec", "Efficiency")
	fmt.Println(strings.Repeat("-", 60))

	baselineTime := time.Duration(0)

	for i, threadCount := range threadCounts {
		// Simulate work distribution
		start := time.Now()

		workPerThread := iterations / threadCount
		remaining := iterations % threadCount

		done := make(chan bool, threadCount)

		for t := 0; t < threadCount; t++ {
			work := workPerThread
			if t < remaining {
				work++
			}

			go func(workCount int) {
				for w := 0; w < workCount; w++ {
					_ = topayz512.ComputeHash(testData)
				}
				done <- true
			}(work)
		}

		for t := 0; t < threadCount; t++ {
			<-done
		}

		duration := time.Since(start)
		opsPerSec := float64(iterations) / duration.Seconds()

		if i == 0 {
			baselineTime = duration
		}

		efficiency := float64(baselineTime) / float64(duration) / float64(threadCount) * 100

		fmt.Printf("%-10d %-15v %-15.0f %-15.1f%%\n",
			threadCount, duration, opsPerSec, efficiency)
	}

	// Optimal thread count validation
	fmt.Printf("\nRecommended optimal threads: %d\n", topayz512.OptimalThreadCount())
	fmt.Printf("System CPU count: %d\n", runtime.NumCPU())
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
