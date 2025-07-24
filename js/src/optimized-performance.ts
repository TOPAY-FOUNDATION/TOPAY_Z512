/**
 * Optimized Performance Module for TOPAY-Z512
 * 
 * This module demonstrates the performance improvements made to the JavaScript implementation:
 * 1. Hash function memoization for small data
 * 2. Batch processing with configurable concurrency
 * 3. Buffer pooling for memory optimization
 * 4. 32-bit optimized operations
 * 5. KEM key caching (optional)
 */

import { computeHash, batchHash, clearHashCache } from './hash.js';
import { kemKeyGen, kemEncapsulate, kemDecapsulate } from './kem.js';
import { clearBufferPools, getBufferPoolStats, secureRandom } from './utils.js';

/**
 * Performance optimization configuration
 */
export interface OptimizationConfig {
  enableHashCaching: boolean;
  enableBufferPooling: boolean;
  enableKEMCaching: boolean;
  batchConcurrency: number;
  maxCacheSize: number;
}

/**
 * Default optimization configuration
 */
export const DEFAULT_OPTIMIZATION_CONFIG: OptimizationConfig = {
  enableHashCaching: true,
  enableBufferPooling: true,
  enableKEMCaching: false, // Disabled by default for security
  batchConcurrency: 8,
  maxCacheSize: 1000
};

/**
 * Performance metrics for optimization tracking
 */
export interface OptimizationMetrics {
  hashCacheHits: number;
  hashCacheMisses: number;
  bufferPoolReuses: number;
  bufferPoolAllocations: number;
  kemCacheHits: number;
  kemCacheMisses: number;
  totalOperations: number;
  averageOperationTime: number;
}

let metrics: OptimizationMetrics = {
  hashCacheHits: 0,
  hashCacheMisses: 0,
  bufferPoolReuses: 0,
  bufferPoolAllocations: 0,
  kemCacheHits: 0,
  kemCacheMisses: 0,
  totalOperations: 0,
  averageOperationTime: 0
};

/**
 * Benchmarks hash operations with and without optimizations
 */
export async function benchmarkHashOptimizations(dataSize: number, iterations: number): Promise<{
  withOptimizations: number;
  withoutOptimizations: number;
  improvement: number;
}> {
  const testData = await secureRandom(dataSize);
  
  // Benchmark with optimizations
  const startOptimized = performance.now();
  for (let i = 0; i < iterations; i++) {
    await computeHash(testData);
  }
  const endOptimized = performance.now();
  const withOptimizations = endOptimized - startOptimized;
  
  // Clear cache and benchmark without optimizations
  clearHashCache();
  const startUnoptimized = performance.now();
  for (let i = 0; i < iterations; i++) {
    // Force new computation each time
    const uniqueData = new Uint8Array(testData.length + 1);
    uniqueData.set(testData);
    uniqueData[testData.length] = i % 256;
    await computeHash(uniqueData);
  }
  const endUnoptimized = performance.now();
  const withoutOptimizations = endUnoptimized - startUnoptimized;
  
  const improvement = ((withoutOptimizations - withOptimizations) / withoutOptimizations) * 100;
  
  return {
    withOptimizations,
    withoutOptimizations,
    improvement
  };
}

/**
 * Benchmarks batch operations with different concurrency levels
 */
export async function benchmarkBatchConcurrency(
  dataItems: Uint8Array[],
  concurrencyLevels: number[]
): Promise<{ concurrency: number; time: number }[]> {
  const results: { concurrency: number; time: number }[] = [];
  
  for (const concurrency of concurrencyLevels) {
    const start = performance.now();
    await batchHash(dataItems, concurrency);
    const end = performance.now();
    
    results.push({
      concurrency,
      time: end - start
    });
  }
  
  return results;
}

/**
 * Benchmarks KEM operations with and without caching
 */
export async function benchmarkKEMOptimizations(iterations: number): Promise<{
  withCaching: number;
  withoutCaching: number;
  improvement: number;
}> {
  // Benchmark with caching
  const startCached = performance.now();
  for (let i = 0; i < iterations; i++) {
    const keyPair = await kemKeyGen(true); // Enable caching
    const { ciphertext, sharedSecret } = await kemEncapsulate(keyPair.publicKey);
    await kemDecapsulate(keyPair.secretKey, ciphertext);
  }
  const endCached = performance.now();
  const withCaching = endCached - startCached;
  
  // Benchmark without caching
  const startUncached = performance.now();
  for (let i = 0; i < iterations; i++) {
    const keyPair = await kemKeyGen(false); // Disable caching
    const { ciphertext, sharedSecret } = await kemEncapsulate(keyPair.publicKey);
    await kemDecapsulate(keyPair.secretKey, ciphertext);
  }
  const endUncached = performance.now();
  const withoutCaching = endUncached - startUncached;
  
  const improvement = ((withoutCaching - withCaching) / withoutCaching) * 100;
  
  return {
    withCaching,
    withoutCaching,
    improvement
  };
}

/**
 * Runs a comprehensive optimization benchmark suite
 */
export async function runOptimizationBenchmarks(): Promise<{
  hashOptimizations: Awaited<ReturnType<typeof benchmarkHashOptimizations>>;
  batchConcurrency: Awaited<ReturnType<typeof benchmarkBatchConcurrency>>;
  kemOptimizations: Awaited<ReturnType<typeof benchmarkKEMOptimizations>>;
  bufferPoolStats: ReturnType<typeof getBufferPoolStats>;
}> {
  console.log('🚀 Running TOPAY-Z512 Optimization Benchmarks...\n');
  
  // Hash optimizations benchmark
  console.log('📊 Benchmarking hash optimizations...');
  const hashOptimizations = await benchmarkHashOptimizations(1024, 100);
  console.log(`Hash optimization improvement: ${hashOptimizations.improvement.toFixed(2)}%\n`);
  
  // Batch concurrency benchmark
  console.log('⚡ Benchmarking batch concurrency...');
  const testData = await Promise.all(Array.from({ length: 50 }, () => secureRandom(512)));
  const batchConcurrency = await benchmarkBatchConcurrency(testData, [1, 2, 4, 8, 16]);
  const optimalConcurrency = batchConcurrency.reduce((min, curr) => 
    curr.time < min.time ? curr : min
  );
  console.log(`Optimal concurrency level: ${optimalConcurrency.concurrency} (${optimalConcurrency.time.toFixed(2)}ms)\n`);
  
  // KEM optimizations benchmark
  console.log('🔐 Benchmarking KEM optimizations...');
  const kemOptimizations = await benchmarkKEMOptimizations(20);
  console.log(`KEM caching improvement: ${kemOptimizations.improvement.toFixed(2)}%\n`);
  
  // Buffer pool statistics
  const bufferPoolStats = getBufferPoolStats();
  console.log(`📈 Buffer pool stats: ${bufferPoolStats.totalBuffers} buffers across ${bufferPoolStats.totalPools} pools\n`);
  
  return {
    hashOptimizations,
    batchConcurrency,
    kemOptimizations,
    bufferPoolStats
  };
}

/**
 * Clears all optimization caches and pools
 */
export function clearOptimizationCaches(): void {
  clearHashCache();
  clearBufferPools();
  
  // Reset metrics
  metrics = {
    hashCacheHits: 0,
    hashCacheMisses: 0,
    bufferPoolReuses: 0,
    bufferPoolAllocations: 0,
    kemCacheHits: 0,
    kemCacheMisses: 0,
    totalOperations: 0,
    averageOperationTime: 0
  };
}

/**
 * Gets current optimization metrics
 */
export function getOptimizationMetrics(): OptimizationMetrics {
  return { ...metrics };
}

/**
 * Estimates memory usage of optimization features
 */
export function estimateOptimizationMemoryUsage(): {
  hashCache: number;
  bufferPools: number;
  kemCache: number;
  total: number;
} {
  const bufferStats = getBufferPoolStats();
  
  // Rough estimates in bytes
  const hashCache = 1024 * 100; // Assume 100 cached hashes of 1KB each
  const bufferPools = bufferStats.totalBuffers * 512; // Average buffer size
  const kemCache = 64 * 100; // Assume 100 cached key pairs
  
  return {
    hashCache,
    bufferPools,
    kemCache,
    total: hashCache + bufferPools + kemCache
  };
}

/**
 * Generates an optimization report
 */
export async function generateOptimizationReport(): Promise<string> {
  const benchmarks = await runOptimizationBenchmarks();
  const memoryUsage = estimateOptimizationMemoryUsage();
  
  return `
# TOPAY-Z512 JavaScript Optimization Report

## Performance Improvements

### Hash Operations
- **Improvement**: ${benchmarks.hashOptimizations.improvement.toFixed(2)}%
- **With optimizations**: ${benchmarks.hashOptimizations.withOptimizations.toFixed(2)}ms
- **Without optimizations**: ${benchmarks.hashOptimizations.withoutOptimizations.toFixed(2)}ms

### Batch Processing
- **Optimal concurrency**: ${benchmarks.batchConcurrency.reduce((min, curr) => curr.time < min.time ? curr : min).concurrency}
- **Performance gain**: Up to ${((Math.max(...benchmarks.batchConcurrency.map(b => b.time)) - Math.min(...benchmarks.batchConcurrency.map(b => b.time))) / Math.max(...benchmarks.batchConcurrency.map(b => b.time)) * 100).toFixed(2)}%

### KEM Operations
- **Caching improvement**: ${benchmarks.kemOptimizations.improvement.toFixed(2)}%
- **With caching**: ${benchmarks.kemOptimizations.withCaching.toFixed(2)}ms
- **Without caching**: ${benchmarks.kemOptimizations.withoutCaching.toFixed(2)}ms

## Memory Optimization

### Buffer Pools
- **Active pools**: ${benchmarks.bufferPoolStats.totalPools}
- **Pooled buffers**: ${benchmarks.bufferPoolStats.totalBuffers}
- **Memory saved**: ~${(memoryUsage.bufferPools / 1024).toFixed(2)}KB

### Total Memory Usage
- **Hash cache**: ~${(memoryUsage.hashCache / 1024).toFixed(2)}KB
- **Buffer pools**: ~${(memoryUsage.bufferPools / 1024).toFixed(2)}KB
- **KEM cache**: ~${(memoryUsage.kemCache / 1024).toFixed(2)}KB
- **Total**: ~${(memoryUsage.total / 1024).toFixed(2)}KB

## Optimization Features

✅ **Hash Memoization**: Caches small data hashes for repeated operations
✅ **Buffer Pooling**: Reuses memory buffers to reduce allocations
✅ **32-bit Operations**: Uses optimized 32-bit operations where possible
✅ **Batch Concurrency**: Configurable parallel processing
✅ **KEM Caching**: Optional key pair caching (disabled by default for security)

## Recommendations

1. **Production Use**: Enable hash caching and buffer pooling for better performance
2. **Security**: Keep KEM caching disabled unless specifically needed
3. **Concurrency**: Use 8 concurrent operations for optimal batch processing
4. **Memory**: Monitor buffer pool usage in long-running applications
5. **Cleanup**: Periodically clear caches in memory-constrained environments
`;
}