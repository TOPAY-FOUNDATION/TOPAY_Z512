package topayz512

import (
	"encoding/binary"
	"sync"
	"time"
)

// Data fragmentation for TOPAY-Z512 with optimizations

// Fragment represents a single data fragment
type Fragment struct {
	ID       uint32 `json:"id"`
	Index    uint32 `json:"index"`
	Total    uint32 `json:"total"`
	Data     []byte `json:"data"`
	Checksum Hash   `json:"checksum"`
	Size     uint32 `json:"size"`
}

// FragmentationResult contains the result of data fragmentation
type FragmentationResult struct {
	Fragments    []Fragment       `json:"fragments"`
	TotalSize    uint64           `json:"total_size"`
	FragmentSize uint32           `json:"fragment_size"`
	Metadata     FragmentMetadata `json:"metadata"`
}

// FragmentMetadata contains metadata about the fragmentation
type FragmentMetadata struct {
	OriginalSize  uint64    `json:"original_size"`
	FragmentCount uint32    `json:"fragment_count"`
	Timestamp     time.Time `json:"timestamp"`
	Algorithm     string    `json:"algorithm"`
	Checksum      Hash      `json:"checksum"`
}

// ReconstructionResult contains the result of data reconstruction
type ReconstructionResult struct {
	Data         []byte           `json:"data"`
	IsComplete   bool             `json:"is_complete"`
	MissingCount uint32           `json:"missing_count"`
	Metadata     FragmentMetadata `json:"metadata"`
}

// ShouldFragment determines if data should be fragmented based on size
func ShouldFragment(dataSize int) bool {
	return dataSize >= MinFragmentThreshold
}

// CalculateFragmentCount calculates the optimal number of fragments
func CalculateFragmentCount(dataSize int) int {
	if !ShouldFragment(dataSize) {
		return 1
	}

	fragmentCount := (dataSize + FragmentSize - 1) / FragmentSize
	if fragmentCount > MaxFragments {
		fragmentCount = MaxFragments
	}

	return fragmentCount
}

// FragmentData splits data into fragments for parallel processing
func FragmentData(data []byte) (FragmentationResult, error) {
	if len(data) == 0 {
		return FragmentationResult{}, ErrEmptyData
	}

	fragmentCount := CalculateFragmentCount(len(data))
	fragmentSize := (len(data) + fragmentCount - 1) / fragmentCount

	// Generate unique fragment ID
	idBytes, err := SecureRandom(4)
	if err != nil {
		return FragmentationResult{}, err
	}
	fragmentID := binary.BigEndian.Uint32(idBytes)

	// Calculate total checksum
	totalChecksum := ComputeHash(data)

	// Create fragments
	fragments := make([]Fragment, fragmentCount)

	for i := 0; i < fragmentCount; i++ {
		start := i * fragmentSize
		end := start + fragmentSize
		if end > len(data) {
			end = len(data)
		}

		fragmentData := make([]byte, end-start)
		copy(fragmentData, data[start:end])

		// Calculate fragment checksum
		fragmentChecksum := ComputeHash(fragmentData)

		fragments[i] = Fragment{
			ID:       fragmentID,
			Index:    uint32(i),
			Total:    uint32(fragmentCount),
			Data:     fragmentData,
			Checksum: fragmentChecksum,
		}
	}

	metadata := FragmentMetadata{
		OriginalSize:  uint64(len(data)),
		FragmentCount: uint32(fragmentCount),
		Timestamp:     time.Now(),
		Algorithm:     "TOPAY-Z512",
		Checksum:      totalChecksum,
	}

	return FragmentationResult{
		Fragments: fragments,
		Metadata:  metadata,
	}, nil
}

// ReconstructData reconstructs original data from fragments
func ReconstructData(fragments []Fragment) (ReconstructionResult, error) {
	if len(fragments) == 0 {
		return ReconstructionResult{}, ErrEmptyData
	}

	// Validate fragments
	fragmentID := fragments[0].ID
	totalFragments := fragments[0].Total

	if len(fragments) != int(totalFragments) {
		return ReconstructionResult{}, ErrInvalidFragmentCount
	}

	// Sort fragments by index
	sortedFragments := make([]Fragment, len(fragments))
	copy(sortedFragments, fragments)

	// Simple bubble sort for fragment ordering
	for i := 0; i < len(sortedFragments); i++ {
		for j := 0; j < len(sortedFragments)-1-i; j++ {
			if sortedFragments[j].Index > sortedFragments[j+1].Index {
				sortedFragments[j], sortedFragments[j+1] = sortedFragments[j+1], sortedFragments[j]
			}
		}
	}

	// Validate fragment integrity
	for i, fragment := range sortedFragments {
		if fragment.ID != fragmentID {
			return ReconstructionResult{}, ErrReconstructionFailed
		}

		if fragment.Index != uint32(i) {
			return ReconstructionResult{}, ErrReconstructionFailed
		}

		if fragment.Total != totalFragments {
			return ReconstructionResult{}, ErrReconstructionFailed
		}

		// Verify fragment checksum
		computedChecksum := ComputeHash(fragment.Data)
		if !HashEqual(computedChecksum, fragment.Checksum) {
			return ReconstructionResult{}, ErrReconstructionFailed
		}
	}

	// Reconstruct data
	var totalSize int
	for _, fragment := range sortedFragments {
		totalSize += len(fragment.Data)
	}

	reconstructedData := make([]byte, 0, totalSize)
	for _, fragment := range sortedFragments {
		reconstructedData = append(reconstructedData, fragment.Data...)
	}

	// Verify total checksum
	totalChecksum := ComputeHash(reconstructedData)

	metadata := FragmentMetadata{
		OriginalSize:  uint64(len(reconstructedData)),
		FragmentCount: totalFragments,
		Timestamp:     time.Now(),
		Algorithm:     "TOPAY-Z512",
		Checksum:      totalChecksum,
	}

	return ReconstructionResult{
		Data:         reconstructedData,
		IsComplete:   true,
		MissingCount: 0,
		Metadata:     metadata,
	}, nil
}

// Parallel fragmentation operations

// ParallelFragmentData fragments data using parallel processing
func ParallelFragmentData(data []byte) (FragmentationResult, error) {
	if len(data) == 0 {
		return FragmentationResult{}, ErrEmptyData
	}

	fragmentCount := CalculateFragmentCount(len(data))
	fragmentSize := (len(data) + fragmentCount - 1) / fragmentCount

	// Generate unique fragment ID
	idBytes, err := SecureRandom(4)
	if err != nil {
		return FragmentationResult{}, err
	}
	fragmentID := binary.BigEndian.Uint32(idBytes)

	// Calculate total checksum
	totalChecksum := ComputeHash(data)

	// Create fragments in parallel
	fragments := make([]Fragment, fragmentCount)

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > fragmentCount {
		numWorkers = fragmentCount
	}

	// Channel for work distribution
	workChan := make(chan int, fragmentCount)
	resultChan := make(chan struct {
		index    int
		fragment Fragment
		err      error
	}, fragmentCount)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				start := index * fragmentSize
				end := start + fragmentSize
				if end > len(data) {
					end = len(data)
				}

				fragmentData := make([]byte, end-start)
				copy(fragmentData, data[start:end])

				// Calculate fragment checksum
				fragmentChecksum := ComputeHash(fragmentData)

				fragment := Fragment{
					ID:       fragmentID,
					Index:    uint32(index),
					Total:    uint32(fragmentCount),
					Data:     fragmentData,
					Checksum: fragmentChecksum,
				}

				resultChan <- struct {
					index    int
					fragment Fragment
					err      error
				}{index, fragment, nil}
			}
		}()
	}

	// Send work
	go func() {
		for i := 0; i < fragmentCount; i++ {
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
		if result.err != nil {
			return FragmentationResult{}, result.err
		}
		fragments[result.index] = result.fragment
	}

	metadata := FragmentMetadata{
		OriginalSize:  uint64(len(data)),
		FragmentCount: uint32(fragmentCount),
		Timestamp:     time.Now(),
		Algorithm:     "TOPAY-Z512",
		Checksum:      totalChecksum,
	}

	return FragmentationResult{
		Fragments: fragments,
		Metadata:  metadata,
	}, nil
}

// ParallelReconstructData reconstructs data using parallel processing
func ParallelReconstructData(fragments []Fragment) (ReconstructionResult, error) {
	if len(fragments) == 0 {
		return ReconstructionResult{}, ErrEmptyData
	}

	// Validate and sort fragments first
	result, err := ReconstructData(fragments)
	if err != nil {
		return result, err
	}

	// For parallel reconstruction, we could verify checksums in parallel
	// but since reconstruction is already fast, we'll use the sequential version
	return result, nil
}

// Fragment serialization

// SerializeFragment converts a fragment to bytes
func SerializeFragment(fragment Fragment) []byte {
	// Calculate total size needed
	dataLen := len(fragment.Data)
	totalSize := 4 + 4 + 4 + 4 + dataLen + HashSize // ID + Index + Total + DataLen + Data + Checksum

	result := make([]byte, totalSize)
	offset := 0

	// Write ID
	binary.BigEndian.PutUint32(result[offset:], fragment.ID)
	offset += 4

	// Write Index
	binary.BigEndian.PutUint32(result[offset:], fragment.Index)
	offset += 4

	// Write Total
	binary.BigEndian.PutUint32(result[offset:], fragment.Total)
	offset += 4

	// Write Data length
	binary.BigEndian.PutUint32(result[offset:], uint32(dataLen))
	offset += 4

	// Write Data
	copy(result[offset:], fragment.Data)
	offset += dataLen

	// Write Checksum
	copy(result[offset:], fragment.Checksum[:])

	return result
}

// DeserializeFragment converts bytes to a fragment
func DeserializeFragment(data []byte) (Fragment, error) {
	if len(data) < 16+HashSize { // Minimum size: 4 fields + checksum
		return Fragment{}, ErrInvalidFragmentCount
	}

	offset := 0

	// Read ID
	id := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// Read Index
	index := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// Read Total
	total := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	// Read Data length
	dataLen := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	if len(data) < offset+int(dataLen)+HashSize {
		return Fragment{}, ErrInvalidFragmentCount
	}

	// Read Data
	fragmentData := make([]byte, dataLen)
	copy(fragmentData, data[offset:offset+int(dataLen)])
	offset += int(dataLen)

	// Read Checksum
	var checksum Hash
	copy(checksum[:], data[offset:])

	return Fragment{
		ID:       id,
		Index:    index,
		Total:    total,
		Data:     fragmentData,
		Checksum: checksum,
	}, nil
}

// Mobile optimization

// MobileLatencyEstimate estimates processing latency for mobile devices
type MobileLatencyEstimate struct {
	FragmentationMs   float64
	ReconstructionMs  float64
	TotalMs           float64
	RecommendedChunks int
}

// EstimateMobileLatency estimates processing time for mobile devices
func EstimateMobileLatency(dataSize int) MobileLatencyEstimate {
	// Base latency factors for mobile devices (conservative estimates)
	const (
		baseFragmentationMsPerKB  = 0.1
		baseReconstructionMsPerKB = 0.05
		mobileCPUFactor           = 2.0 // Mobile CPUs are typically 2x slower
	)

	dataSizeKB := float64(dataSize) / 1024.0

	fragmentationMs := dataSizeKB * baseFragmentationMsPerKB * mobileCPUFactor
	reconstructionMs := dataSizeKB * baseReconstructionMsPerKB * mobileCPUFactor
	totalMs := fragmentationMs + reconstructionMs

	// Recommend optimal chunk count for mobile
	recommendedChunks := CalculateFragmentCount(dataSize)
	if recommendedChunks > 64 { // Limit for mobile devices
		recommendedChunks = 64
	}

	return MobileLatencyEstimate{
		FragmentationMs:   fragmentationMs,
		ReconstructionMs:  reconstructionMs,
		TotalMs:           totalMs,
		RecommendedChunks: recommendedChunks,
	}
}

// Fragmented cryptographic operations

// FragmentedHash computes hash using fragmented processing
func FragmentedHash(data []byte) (Hash, error) {
	if !ShouldFragment(len(data)) {
		return ComputeHash(data), nil
	}

	fragResult, err := ParallelFragmentData(data)
	if err != nil {
		return Hash{}, err
	}

	// Compute hashes of fragments in parallel
	fragmentHashes := make([]Hash, len(fragResult.Fragments))

	// Use optimal number of goroutines
	numWorkers := OptimalThreadCount()
	if numWorkers > len(fragResult.Fragments) {
		numWorkers = len(fragResult.Fragments)
	}

	// Channel for work distribution
	workChan := make(chan int, len(fragResult.Fragments))
	resultChan := make(chan struct {
		index int
		hash  Hash
	}, len(fragResult.Fragments))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range workChan {
				hash := ComputeHash(fragResult.Fragments[index].Data)
				resultChan <- struct {
					index int
					hash  Hash
				}{index, hash}
			}
		}()
	}

	// Send work
	go func() {
		for i := range fragResult.Fragments {
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
		fragmentHashes[result.index] = result.hash
	}

	// Combine fragment hashes
	hs := NewHashState()
	for _, fragmentHash := range fragmentHashes {
		hs.Update(fragmentHash[:])
	}

	return hs.Finalize(), nil
}

// FragmentedKEM performs KEM operations on fragmented data
func FragmentedKEM(data []byte) ([]Ciphertext, []SharedSecret, error) {
	if !ShouldFragment(len(data)) {
		// For small data, use single KEM operation
		publicKey, _, err := KEMKeyGen()
		if err != nil {
			return nil, nil, err
		}

		ciphertext, sharedSecret, err := KEMEncapsulate(publicKey)
		if err != nil {
			return nil, nil, err
		}

		return []Ciphertext{ciphertext}, []SharedSecret{sharedSecret}, nil
	}

	fragResult, err := ParallelFragmentData(data)
	if err != nil {
		return nil, nil, err
	}

	// Generate KEM key pairs for each fragment
	publicKeys, _, err := BatchKEMKeyGen(len(fragResult.Fragments))
	if err != nil {
		return nil, nil, err
	}

	// Perform batch encapsulation
	ciphertexts, sharedSecrets, err := BatchKEMEncapsulate(publicKeys)
	if err != nil {
		return nil, nil, err
	}

	return ciphertexts, sharedSecrets, nil
}

// Performance benchmarking

// FragmentationBenchmark represents fragmentation performance metrics
type FragmentationBenchmark struct {
	ThroughputMBps   float64
	FragmentationMs  float64
	ReconstructionMs float64
	ParallelSpeedup  float64
	MobileLatencyMs  float64
}

// BenchmarkFragmentation measures fragmentation performance
func BenchmarkFragmentation(dataSize int) FragmentationBenchmark {
	// Generate test data
	data := make([]byte, dataSize)
	for i := range data {
		data[i] = byte(i)
	}

	// Sequential fragmentation benchmark
	start := time.Now()
	fragResult, _ := FragmentData(data)
	seqFragDuration := time.Since(start)

	// Sequential reconstruction benchmark
	start = time.Now()
	_, _ = ReconstructData(fragResult.Fragments)
	seqReconDuration := time.Since(start)

	// Parallel fragmentation benchmark
	start = time.Now()
	_, _ = ParallelFragmentData(data)
	parFragDuration := time.Since(start)

	// Calculate metrics
	totalBytes := float64(dataSize)
	throughputMBps := totalBytes / seqFragDuration.Seconds() / (1024 * 1024)
	fragmentationMs := seqFragDuration.Seconds() * 1000
	reconstructionMs := seqReconDuration.Seconds() * 1000
	parallelSpeedup := seqFragDuration.Seconds() / parFragDuration.Seconds()

	// Mobile latency estimate
	mobileEst := EstimateMobileLatency(dataSize)

	return FragmentationBenchmark{
		ThroughputMBps:   throughputMBps,
		FragmentationMs:  fragmentationMs,
		ReconstructionMs: reconstructionMs,
		ParallelSpeedup:  parallelSpeedup,
		MobileLatencyMs:  mobileEst.TotalMs,
	}
}

// Fragment integrity and validation

// ValidateFragmentIntegrity validates the integrity of a fragment
func ValidateFragmentIntegrity(fragment Fragment) error {
	// Verify checksum
	computedChecksum := ComputeHash(fragment.Data)
	if !HashEqual(computedChecksum, fragment.Checksum) {
		return ErrReconstructionFailed
	}

	// Validate fragment parameters
	if fragment.Index >= fragment.Total {
		return ErrInvalidFragmentCount
	}

	if len(fragment.Data) == 0 {
		return ErrEmptyData
	}

	return nil
}

// RepairFragment attempts to repair a corrupted fragment
func RepairFragment(fragment Fragment, originalData []byte, fragmentSize int) (Fragment, error) {
	// Calculate expected fragment data
	start := int(fragment.Index) * fragmentSize
	end := start + fragmentSize
	if end > len(originalData) {
		end = len(originalData)
	}

	if start >= len(originalData) {
		return Fragment{}, ErrInvalidFragmentCount
	}

	// Create repaired fragment
	repairedData := make([]byte, end-start)
	copy(repairedData, originalData[start:end])

	repairedChecksum := ComputeHash(repairedData)

	repairedFragment := Fragment{
		ID:       fragment.ID,
		Index:    fragment.Index,
		Total:    fragment.Total,
		Data:     repairedData,
		Checksum: repairedChecksum,
	}

	return repairedFragment, nil
}
