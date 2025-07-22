package topayz512

import (
	"sync"
)

// Memory pool management for high-performance operations

// BytePool manages reusable byte slices to reduce GC pressure
type BytePool struct {
	pools map[int]*sync.Pool
	mutex sync.RWMutex
}

// Global byte pools for common sizes
var (
	globalBytePool = NewBytePool()

	// Pre-defined pools for common sizes
	pool64   = &sync.Pool{New: func() interface{} { return make([]byte, 64) }}
	pool256  = &sync.Pool{New: func() interface{} { return make([]byte, 256) }}
	pool1024 = &sync.Pool{New: func() interface{} { return make([]byte, 1024) }}
	pool4096 = &sync.Pool{New: func() interface{} { return make([]byte, 4096) }}
)

// NewBytePool creates a new byte pool manager
func NewBytePool() *BytePool {
	return &BytePool{
		pools: make(map[int]*sync.Pool),
	}
}

// Get retrieves a byte slice from the pool
func (bp *BytePool) Get(size int) []byte {
	// Use pre-defined pools for common sizes
	switch {
	case size <= 64:
		buf := pool64.Get().([]byte)
		return buf[:size]
	case size <= 256:
		buf := pool256.Get().([]byte)
		return buf[:size]
	case size <= 1024:
		buf := pool1024.Get().([]byte)
		return buf[:size]
	case size <= 4096:
		buf := pool4096.Get().([]byte)
		return buf[:size]
	}

	// For larger or uncommon sizes, use dynamic pools
	bp.mutex.RLock()
	pool, exists := bp.pools[size]
	bp.mutex.RUnlock()

	if !exists {
		bp.mutex.Lock()
		// Double-check after acquiring write lock
		if pool, exists = bp.pools[size]; !exists {
			pool = &sync.Pool{
				New: func() interface{} {
					return make([]byte, size)
				},
			}
			bp.pools[size] = pool
		}
		bp.mutex.Unlock()
	}

	return pool.Get().([]byte)
}

// Put returns a byte slice to the pool
func (bp *BytePool) Put(buf []byte) {
	if buf == nil {
		return
	}

	size := cap(buf)

	// Clear the buffer for security
	for i := range buf {
		buf[i] = 0
	}

	// Use pre-defined pools for common sizes
	switch {
	case size == 64:
		pool64.Put(buf[:64])
		return
	case size == 256:
		pool256.Put(buf[:256])
		return
	case size == 1024:
		pool1024.Put(buf[:1024])
		return
	case size == 4096:
		pool4096.Put(buf[:4096])
		return
	}

	// For larger or uncommon sizes, use dynamic pools
	bp.mutex.RLock()
	pool, exists := bp.pools[size]
	bp.mutex.RUnlock()

	if exists {
		pool.Put(buf)
	}
	// If pool doesn't exist, just let GC handle it
}

// GetBuffer is a convenience function using the global pool
func GetBuffer(size int) []byte {
	return globalBytePool.Get(size)
}

// PutBuffer is a convenience function using the global pool
func PutBuffer(buf []byte) {
	globalBytePool.Put(buf)
}

// HashStatePool manages reusable hash states
type HashStatePool struct {
	pool sync.Pool
}

// NewHashStatePool creates a new hash state pool
func NewHashStatePool() *HashStatePool {
	return &HashStatePool{
		pool: sync.Pool{
			New: func() interface{} {
				return NewHashState()
			},
		},
	}
}

// Get retrieves a hash state from the pool
func (hsp *HashStatePool) Get() *HashState {
	hs := hsp.pool.Get().(*HashState)
	hs.Reset()
	return hs
}

// Put returns a hash state to the pool
func (hsp *HashStatePool) Put(hs *HashState) {
	if hs != nil {
		hs.Reset() // Clear state for security
		hsp.pool.Put(hs)
	}
}

// Global hash state pool
var globalHashStatePool = NewHashStatePool()

// GetHashState retrieves a hash state from the global pool
func GetHashState() *HashState {
	return globalHashStatePool.Get()
}

// PutHashState returns a hash state to the global pool
func PutHashState(hs *HashState) {
	globalHashStatePool.Put(hs)
}

// WorkerPool manages a pool of worker goroutines
type WorkerPool struct {
	workers   int
	workChan  chan func()
	closeChan chan struct{}
	wg        sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = OptimalThreadCount()
	}

	wp := &WorkerPool{
		workers:   workers,
		workChan:  make(chan func(), workers*2), // Buffered channel
		closeChan: make(chan struct{}),
	}

	// Start workers
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

// worker is the main worker goroutine
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case work := <-wp.workChan:
			work()
		case <-wp.closeChan:
			return
		}
	}
}

// Submit submits work to the pool
func (wp *WorkerPool) Submit(work func()) {
	select {
	case wp.workChan <- work:
	case <-wp.closeChan:
		// Pool is closed, execute work directly
		work()
	}
}

// Close closes the worker pool
func (wp *WorkerPool) Close() {
	close(wp.closeChan)
	wp.wg.Wait()
}

// Global worker pool
var globalWorkerPool *WorkerPool

// InitializeGlobalPools initializes global pools
func InitializeGlobalPools() {
	if globalWorkerPool == nil {
		globalWorkerPool = NewWorkerPool(OptimalThreadCount())
	}
}

// SubmitWork submits work to the global worker pool
func SubmitWork(work func()) {
	if globalWorkerPool == nil {
		InitializeGlobalPools()
	}
	globalWorkerPool.Submit(work)
}

// CleanupGlobalPools cleans up global pools
func CleanupGlobalPools() {
	if globalWorkerPool != nil {
		globalWorkerPool.Close()
		globalWorkerPool = nil
	}
}
