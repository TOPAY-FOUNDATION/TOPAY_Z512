/**
 * TOPAY-Z512 Optimization Demo
 *
 * This example demonstrates the performance improvements made to the JavaScript implementation.
 */

import {
  runOptimizationBenchmarks,
  generateOptimizationReport,
  clearOptimizationCaches
} from '../optimized-performance.js';
import { computeHash, batchHash } from '../hash.js';
import { kemKeyGen, kemEncapsulate, kemDecapsulate } from '../kem.js';
import { secureRandom, getBufferPoolStats } from '../utils.js';

async function demonstrateOptimizations() {
  console.log('🎯 TOPAY-Z512 JavaScript Optimization Demo\n');
  console.log('This demo showcases the performance improvements made to the implementation.\n');

  // Clear any existing caches
  clearOptimizationCaches();

  // 1. Hash Memoization Demo
  console.log('1️⃣ Hash Memoization Demo');
  console.log('='.repeat(50));

  const smallData = await secureRandom(512); // Small data that will be cached

  console.log('First hash computation (cache miss):');
  const start1 = performance.now();
  const hash1 = await computeHash(smallData);
  const end1 = performance.now();
  console.log(`Time: ${(end1 - start1).toFixed(3)}ms`);

  console.log('Second hash computation (cache hit):');
  const start2 = performance.now();
  const hash2 = await computeHash(smallData);
  const end2 = performance.now();
  console.log(`Time: ${(end2 - start2).toFixed(3)}ms`);
  console.log(`Speedup: ${((end1 - start1) / (end2 - start2)).toFixed(2)}x\n`);

  // 2. Batch Processing Demo
  console.log('2️⃣ Batch Processing Demo');
  console.log('='.repeat(50));

  const batchData = await Promise.all(Array.from({ length: 20 }, () => secureRandom(256)));

  console.log('Sequential processing:');
  const startSeq = performance.now();
  for (const data of batchData) {
    await computeHash(data);
  }
  const endSeq = performance.now();
  console.log(`Time: ${(endSeq - startSeq).toFixed(3)}ms`);

  console.log('Batch processing (concurrency=8):');
  const startBatch = performance.now();
  await batchHash(batchData, 8);
  const endBatch = performance.now();
  console.log(`Time: ${(endBatch - startBatch).toFixed(3)}ms`);
  console.log(`Speedup: ${((endSeq - startSeq) / (endBatch - startBatch)).toFixed(2)}x\n`);

  // 3. Buffer Pool Demo
  console.log('3️⃣ Buffer Pool Demo');
  console.log('='.repeat(50));

  console.log('Initial buffer pool stats:');
  let poolStats = getBufferPoolStats();
  console.log(`Pools: ${poolStats.totalPools}, Buffers: ${poolStats.totalBuffers}`);

  // Perform operations that use buffers
  const operations = await Promise.all(Array.from({ length: 10 }, () => secureRandom(1024)));
  for (const data of operations) {
    await computeHash(data);
  }

  console.log('After operations:');
  poolStats = getBufferPoolStats();
  console.log(`Pools: ${poolStats.totalPools}, Buffers: ${poolStats.totalBuffers}`);
  console.log(`Buffer sizes: [${poolStats.sizes.join(', ')}]\n`);

  // 4. KEM Caching Demo (optional)
  console.log('4️⃣ KEM Caching Demo');
  console.log('='.repeat(50));

  console.log('KEM operations without caching:');
  const startKemNoCache = performance.now();
  for (let i = 0; i < 5; i++) {
    const keyPair = await kemKeyGen(false);
    const { ciphertext } = await kemEncapsulate(keyPair.publicKey);
    await kemDecapsulate(keyPair.secretKey, ciphertext);
  }
  const endKemNoCache = performance.now();
  console.log(`Time: ${(endKemNoCache - startKemNoCache).toFixed(3)}ms`);

  console.log('KEM operations with caching:');
  const startKemCache = performance.now();
  for (let i = 0; i < 5; i++) {
    const keyPair = await kemKeyGen(true);
    const { ciphertext } = await kemEncapsulate(keyPair.publicKey);
    await kemDecapsulate(keyPair.secretKey, ciphertext);
  }
  const endKemCache = performance.now();
  console.log(`Time: ${(endKemCache - startKemCache).toFixed(3)}ms`);
  console.log(
    `Speedup: ${((endKemNoCache - startKemNoCache) / (endKemCache - startKemCache)).toFixed(2)}x\n`
  );

  // 5. Comprehensive Benchmark
  console.log('5️⃣ Comprehensive Benchmark Suite');
  console.log('='.repeat(50));

  const benchmarkResults = await runOptimizationBenchmarks();

  // 6. Generate Report
  console.log('6️⃣ Optimization Report');
  console.log('='.repeat(50));

  const report = await generateOptimizationReport();
  console.log(report);

  console.log('\n✅ Optimization demo completed!');
  console.log('\nKey takeaways:');
  console.log('• Hash memoization provides significant speedup for repeated small data');
  console.log('• Batch processing with optimal concurrency improves throughput');
  console.log('• Buffer pooling reduces memory allocations and GC pressure');
  console.log('• 32-bit optimized operations improve performance on modern hardware');
  console.log('• KEM caching can speed up repeated operations (use with caution)');
}

// Run the demo
demonstrateOptimizations().catch(console.error);
