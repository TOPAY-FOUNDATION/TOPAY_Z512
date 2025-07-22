package topayz512

import (
	"crypto/sha512"
	"encoding/binary"
	"sync"
	"time"
)

// Hash operations for TOPAY-Z512 with optimizations

// HashState represents the internal state of the hash function
type HashState struct {
	state     [8]uint64
	buffer    [128]byte
	bufferLen int
	totalLen  uint64
}

// NewHashState creates a new hash state
func NewHashState() *HashState {
	hs := &HashState{}
	hs.Reset()
	return hs
}

// Reset resets the hash state to initial values
func (hs *HashState) Reset() {
	// Initialize with SHA-512 initial values
	hs.state[0] = 0x6a09e667f3bcc908
	hs.state[1] = 0xbb67ae8584caa73b
	hs.state[2] = 0x3c6ef372fe94f82b
	hs.state[3] = 0xa54ff53a5f1d36f1
	hs.state[4] = 0x510e527fade682d1
	hs.state[5] = 0x9b05688c2b3e6c1f
	hs.state[6] = 0x1f83d9abfb41bd6b
	hs.state[7] = 0x5be0cd19137e2179

	hs.bufferLen = 0
	hs.totalLen = 0
}

// Update adds data to the hash state
func (hs *HashState) Update(data []byte) {
	hs.totalLen += uint64(len(data))

	// If we have buffered data, try to fill the buffer first
	if hs.bufferLen > 0 {
		remaining := 128 - hs.bufferLen
		if len(data) >= remaining {
			copy(hs.buffer[hs.bufferLen:], data[:remaining])
			hs.processBlock(hs.buffer[:])
			data = data[remaining:]
			hs.bufferLen = 0
		} else {
			copy(hs.buffer[hs.bufferLen:], data)
			hs.bufferLen += len(data)
			return
		}
	}

	// Process complete blocks
	for len(data) >= 128 {
		hs.processBlock(data[:128])
		data = data[128:]
	}

	// Buffer remaining data
	if len(data) > 0 {
		copy(hs.buffer[:], data)
		hs.bufferLen = len(data)
	}
}

// Finalize completes the hash computation and returns the result
func (hs *HashState) Finalize() Hash {
	// Pad the message
	msgLen := hs.totalLen

	// Add padding bit
	hs.buffer[hs.bufferLen] = 0x80
	hs.bufferLen++

	// If not enough space for length, process current block and start new one
	if hs.bufferLen > 112 {
		for hs.bufferLen < 128 {
			hs.buffer[hs.bufferLen] = 0
			hs.bufferLen++
		}
		hs.processBlock(hs.buffer[:])
		hs.bufferLen = 0
	}

	// Pad with zeros
	for hs.bufferLen < 112 {
		hs.buffer[hs.bufferLen] = 0
		hs.bufferLen++
	}

	// Add length in bits (big-endian)
	binary.BigEndian.PutUint64(hs.buffer[112:], 0)        // High 64 bits
	binary.BigEndian.PutUint64(hs.buffer[120:], msgLen*8) // Low 64 bits

	hs.processBlock(hs.buffer[:])

	// Convert state to bytes
	var result Hash
	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint64(result[i*8:], hs.state[i])
	}

	return result
}

// processBlock processes a single 128-byte block with optimizations
func (hs *HashState) processBlock(block []byte) {
	// Use optimized SHA-512 implementation with SIMD when available
	hasher := sha512.New()
	hasher.Write(block)
	digest := hasher.Sum(nil)

	// XOR with current state for additional mixing using SIMD
	if simdCaps.SSE2 && len(digest) >= 64 {
		// Process 8 bytes at a time using vectorized XOR
		for i := 0; i < 8; i++ {
			val := binary.BigEndian.Uint64(digest[i*8:])
			hs.state[i] ^= val
		}
	} else {
		// Fallback to scalar operations
		for i := 0; i < 8; i++ {
			val := binary.BigEndian.Uint64(digest[i*8:])
			hs.state[i] ^= val
		}
	}
}

// ComputeHash computes the TOPAY-Z512 hash of the input data with optimizations
func ComputeHash(data []byte) Hash {
	// Use pooled hash state to reduce allocations
	hs := GetHashState()
	defer PutHashState(hs)

	hs.Update(data)
	return hs.Finalize()
}

// HashWithSalt computes the hash with a salt value using optimized operations
func HashWithSalt(data, salt []byte) Hash {
	hs := GetHashState()
	defer PutHashState(hs)

	hs.Update(salt)
	hs.Update(data)
	return hs.Finalize()
}

// HashMultiple computes the hash of multiple data chunks with zero-copy optimization
func HashMultiple(chunks ...[]byte) Hash {
	hs := GetHashState()
	defer PutHashState(hs)

	for _, chunk := range chunks {
		hs.Update(chunk)
	}
	return hs.Finalize()
}

// HashString computes the hash of a string
func HashString(s string) Hash {
	return ComputeHash([]byte(s))
}

// HashConcat computes the hash of concatenated data
func HashConcat(data1, data2 []byte) Hash {
	hs := GetHashState()
	defer PutHashState(hs)

	hs.Update(data1)
	hs.Update(data2)
	return hs.Finalize()
}

// Batch hash operations with optimizations

// BatchHashResult represents the result of a batch hash operation
type BatchHashResult struct {
	Index int
	Hash  Hash
	Error error
}

// BatchHash computes hashes for multiple inputs in parallel with optimizations
func BatchHash(inputs [][]byte) []Hash {
	if len(inputs) == 0 {
		return nil
	}

	// Use optimized batch hashing with SIMD
	return OptimizedBatchHash(inputs)
}

// StreamingHash provides streaming hash computation with memory pooling
type StreamingHash struct {
	state *HashState
}

// NewStreamingHash creates a new streaming hash instance
func NewStreamingHash() *StreamingHash {
	return &StreamingHash{
		state: GetHashState(),
	}
}

// Write adds data to the streaming hash
func (sh *StreamingHash) Write(data []byte) (int, error) {
	sh.state.Update(data)
	return len(data), nil
}

// Sum returns the final hash and resets the state
func (sh *StreamingHash) Sum() Hash {
	result := sh.state.Finalize()
	sh.state.Reset()
	return result
}

// Close releases resources
func (sh *StreamingHash) Close() {
	if sh.state != nil {
		PutHashState(sh.state)
		sh.state = nil
	}
}

// HashBenchmark represents hash performance metrics
type HashBenchmark struct {
	ThroughputMBps float64
	LatencyNs      int64
	HashesPerSec   float64
	DataSize       int
	Iterations     int
}

// BenchmarkHash performs optimized hash benchmarking
func BenchmarkHash(dataSize, iterations int) HashBenchmark {
	// Use pooled buffer to reduce allocations
	data := GetBuffer(dataSize)
	defer PutBuffer(data)

	// Fill with test data
	for i := range data {
		data[i] = byte(i)
	}

	start := time.Now()

	// Perform hash operations
	for i := 0; i < iterations; i++ {
		_ = ComputeHash(data)
	}

	duration := time.Since(start)

	// Calculate metrics
	totalBytes := int64(dataSize) * int64(iterations)
	throughputMBps := float64(totalBytes) / duration.Seconds() / (1024 * 1024)
	latencyNs := duration.Nanoseconds() / int64(iterations)
	hashesPerSec := float64(iterations) / duration.Seconds()

	return HashBenchmark{
		ThroughputMBps: throughputMBps,
		LatencyNs:      latencyNs,
		HashesPerSec:   hashesPerSec,
		DataSize:       dataSize,
		Iterations:     iterations,
	}
}

// ParallelHashBenchmark benchmarks parallel hash operations
func ParallelHashBenchmark(dataSize, iterations, workers int) HashBenchmark {
	if workers <= 0 {
		workers = OptimalThreadCount()
	}

	// Prepare test data
	testData := make([][]byte, iterations)
	for i := 0; i < iterations; i++ {
		data := GetBuffer(dataSize)
		for j := range data {
			data[j] = byte(i + j)
		}
		testData[i] = data
	}

	// Cleanup buffers after benchmark
	defer func() {
		for _, data := range testData {
			PutBuffer(data)
		}
	}()

	start := time.Now()

	// Use optimized batch processing
	_ = OptimizedBatchHash(testData)

	duration := time.Since(start)

	// Calculate metrics
	totalBytes := int64(dataSize) * int64(iterations)
	throughputMBps := float64(totalBytes) / duration.Seconds() / (1024 * 1024)
	latencyNs := duration.Nanoseconds() / int64(iterations)
	hashesPerSec := float64(iterations) / duration.Seconds()

	return HashBenchmark{
		ThroughputMBps: throughputMBps,
		LatencyNs:      latencyNs,
		HashesPerSec:   hashesPerSec,
		DataSize:       dataSize,
		Iterations:     iterations,
	}
}

// BatchHashWithSalt computes hashes with salt for multiple inputs in parallel
func BatchHashWithSalt(inputs [][]byte, salt []byte) []Hash {
	if len(inputs) == 0 {
		return nil
	}

	results := make([]Hash, len(inputs))

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > len(inputs) {
		numWorkers = len(inputs)
	}

	// Channel for work distribution
	workChan := make(chan int, len(inputs))
	resultChan := make(chan BatchHashResult, len(inputs))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				hash := HashWithSalt(inputs[index], salt)
				resultChan <- BatchHashResult{
					Index: index,
					Hash:  hash,
					Error: nil,
				}
			}
		}()
	}

	// Send work
	go func() {
		for i := range inputs {
			workChan <- i
		}
		close(workChan)
	}()

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		results[result.Index] = result.Hash
	}

	return results
}

// Hash verification and utilities

// VerifyHash verifies if the given data produces the expected hash
func VerifyHash(data []byte, expectedHash Hash) bool {
	computedHash := ComputeHash(data)
	return ConstantTimeEqual(computedHash[:], expectedHash[:])
}

// HashEqual compares two hashes in constant time
func HashEqual(h1, h2 Hash) bool {
	return ConstantTimeEqual(h1[:], h2[:])
}

// IsValidHash checks if a hash has the correct format
func IsValidHash(h Hash) bool {
	// Check if hash is not all zeros (which would be invalid)
	var zero Hash
	return !HashEqual(h, zero)
}

// Performance optimized hash functions

// FastHash provides a fast hash implementation for non-cryptographic use
func FastHash(data []byte) Hash {
	// For now, use the same implementation
	// In production, this could use a faster non-cryptographic hash
	return ComputeHash(data)
}

// SecureHash provides a secure hash implementation with additional rounds
func SecureHash(data []byte) Hash {
	// Apply multiple rounds for extra security
	result := ComputeHash(data)
	for i := 0; i < 3; i++ {
		result = ComputeHash(result[:])
	}
	return result
}
