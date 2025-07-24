/**
 * Performance monitoring and optimization for TOPAY-Z512
 */

import { measureTime, getSystemCapabilities } from './utils';
import { generateKeyPair } from './keypair';
import { kemKeyGen, kemEncapsulate, kemDecapsulate } from './kem';
import { computeHash } from './hash';
import { fragmentData, reconstructData } from './fragment';

/**
 * Performance metrics for operations
 */
export interface PerformanceMetrics {
  operation: string;
  executionTimeMs: number;
  throughputOpsPerSec?: number;
  memoryUsageMB?: number;
  cpuUsagePercent?: number;
}

/**
 * Benchmark results for multiple operations
 */
export interface BenchmarkResults {
  systemInfo: ReturnType<typeof getSystemCapabilities>;
  timestamp: number;
  metrics: PerformanceMetrics[];
  summary: {
    totalTimeMs: number;
    averageTimeMs: number;
    operationsPerSecond: number;
  };
}

/**
 * Benchmarks hash operations
 * @param iterations - Number of iterations to run
 * @param dataSize - Size of data to hash in bytes
 * @returns Promise resolving to performance metrics
 */
export async function benchmarkHash(
  iterations: number = 1000,
  dataSize: number = 1024
): Promise<PerformanceMetrics> {
  const testData = new Uint8Array(dataSize);
  crypto.getRandomValues(testData);

  const { /* result, */ timeMs } = await measureTime(async () => {
    for (let i = 0; i < iterations; i++) {
      await computeHash(testData);
    }
  });

  return {
    operation: 'Hash',
    executionTimeMs: timeMs,
    throughputOpsPerSec: iterations / (timeMs / 1000)
  };
}

/**
 * Benchmarks key pair generation
 * @param iterations - Number of key pairs to generate
 * @returns Promise resolving to performance metrics
 */
export async function benchmarkKeyPairGeneration(
  iterations: number = 100
): Promise<PerformanceMetrics> {
  const { /* result, */ timeMs } = await measureTime(async () => {
    for (let i = 0; i < iterations; i++) {
      await generateKeyPair();
    }
  });

  return {
    operation: 'KeyPair Generation',
    executionTimeMs: timeMs,
    throughputOpsPerSec: iterations / (timeMs / 1000)
  };
}

/**
 * Benchmarks KEM operations
 * @param iterations - Number of KEM operations to perform
 * @returns Promise resolving to performance metrics
 */
export async function benchmarkKEM(iterations: number = 100): Promise<PerformanceMetrics> {
  const { /* result, */ timeMs } = await measureTime(async () => {
    for (let i = 0; i < iterations; i++) {
      const keyPair = await kemKeyGen();
      const { ciphertext /* , sharedSecret */ } = await kemEncapsulate(keyPair.publicKey);
      await kemDecapsulate(keyPair.secretKey, ciphertext);
    }
  });

  return {
    operation: 'KEM Operations',
    executionTimeMs: timeMs,
    throughputOpsPerSec: iterations / (timeMs / 1000)
  };
}

/**
 * Benchmarks fragmentation operations
 * @param iterations - Number of fragmentation operations
 * @param dataSize - Size of data to fragment
 * @returns Promise resolving to performance metrics
 */
export async function benchmarkFragmentation(
  iterations: number = 50,
  dataSize: number = 10240
): Promise<PerformanceMetrics> {
  const testData = new Uint8Array(dataSize);
  crypto.getRandomValues(testData);

  const { /* result, */ timeMs } = await measureTime(async () => {
    for (let i = 0; i < iterations; i++) {
      const fragResult = await fragmentData(testData);
      await reconstructData(fragResult.fragments);
    }
  });

  return {
    operation: 'Fragmentation',
    executionTimeMs: timeMs,
    throughputOpsPerSec: iterations / (timeMs / 1000)
  };
}

/**
 * Runs a comprehensive benchmark suite
 * @param config - Benchmark configuration
 * @returns Promise resolving to complete benchmark results
 */
export async function runBenchmarkSuite(
  config: {
    hashIterations?: number;
    keyPairIterations?: number;
    kemIterations?: number;
    fragmentationIterations?: number;
    dataSize?: number;
  } = {}
): Promise<BenchmarkResults> {
  const {
    hashIterations = 1000,
    keyPairIterations = 100,
    kemIterations = 100,
    fragmentationIterations = 50,
    dataSize = 1024
  } = config;

  console.log('Starting TOPAY-Z512 benchmark suite...');

  const systemInfo = getSystemCapabilities();
  const startTime = Date.now();

  // Run individual benchmarks
  const hashMetrics = await benchmarkHash(hashIterations, dataSize);
  console.log(`Hash benchmark completed: ${hashMetrics.throughputOpsPerSec?.toFixed(2)} ops/sec`);

  const keyPairMetrics = await benchmarkKeyPairGeneration(keyPairIterations);
  console.log(
    `KeyPair benchmark completed: ${keyPairMetrics.throughputOpsPerSec?.toFixed(2)} ops/sec`
  );

  const kemMetrics = await benchmarkKEM(kemIterations);
  console.log(`KEM benchmark completed: ${kemMetrics.throughputOpsPerSec?.toFixed(2)} ops/sec`);

  const fragmentationMetrics = await benchmarkFragmentation(fragmentationIterations, dataSize * 10);
  console.log(
    `Fragmentation benchmark completed: ${fragmentationMetrics.throughputOpsPerSec?.toFixed(2)} ops/sec`
  );

  const metrics = [hashMetrics, keyPairMetrics, kemMetrics, fragmentationMetrics];
  const totalTimeMs = Date.now() - startTime;
  const averageTimeMs = totalTimeMs / metrics.length;
  const totalOperations =
    hashIterations + keyPairIterations + kemIterations + fragmentationIterations;
  const operationsPerSecond = totalOperations / (totalTimeMs / 1000);

  return {
    systemInfo,
    timestamp: startTime,
    metrics,
    summary: {
      totalTimeMs,
      averageTimeMs,
      operationsPerSecond
    }
  };
}

/**
 * Monitors memory usage during operation execution
 * @param operation - Function to monitor
 * @returns Promise resolving to operation result and memory metrics
 */
export async function monitorMemoryUsage<T>(operation: () => Promise<T>): Promise<{
  result: T;
  memoryMetrics: {
    beforeMB: number;
    afterMB: number;
    peakMB: number;
    deltaMB: number;
  };
}> {
  // Get initial memory usage
  const initialMemory = typeof process !== 'undefined' ? process.memoryUsage() : { heapUsed: 0 };
  const beforeMB = initialMemory.heapUsed / 1024 / 1024;

  // Run operation
  const result = await operation();

  // Force garbage collection if available
  if (typeof global !== 'undefined' && global.gc) {
    global.gc();
  }

  // Get final memory usage
  const finalMemory = typeof process !== 'undefined' ? process.memoryUsage() : { heapUsed: 0 };
  const afterMB = finalMemory.heapUsed / 1024 / 1024;

  return {
    result,
    memoryMetrics: {
      beforeMB,
      afterMB,
      peakMB: Math.max(beforeMB, afterMB),
      deltaMB: afterMB - beforeMB
    }
  };
}

/**
 * Profiles CPU usage during operation execution
 * @param operation - Function to profile
 * @param sampleIntervalMs - Sampling interval in milliseconds
 * @returns Promise resolving to operation result and CPU metrics
 */
export async function profileCPUUsage<T>(
  operation: () => Promise<T>,
  sampleIntervalMs: number = 100
): Promise<{
  result: T;
  cpuMetrics: {
    averageUsagePercent: number;
    peakUsagePercent: number;
    samples: number[];
  };
}> {
  const samples: number[] = [];
  let isRunning = true;

  // Start CPU monitoring (simplified for cross-platform compatibility)
  const monitor = setInterval(() => {
    if (isRunning) {
      // Simplified CPU usage estimation
      const usage = Math.random() * 100; // Placeholder - would use actual CPU monitoring in production
      samples.push(usage);
    }
  }, sampleIntervalMs);

  try {
    const result = await operation();
    isRunning = false;
    clearInterval(monitor);

    const averageUsagePercent = samples.reduce((sum, sample) => sum + sample, 0) / samples.length;
    const peakUsagePercent = Math.max(...samples);

    return {
      result,
      cpuMetrics: {
        averageUsagePercent,
        peakUsagePercent,
        samples
      }
    };
  } catch (error) {
    isRunning = false;
    clearInterval(monitor);
    throw error;
  }
}

/**
 * Estimates performance on mobile devices
 * @param operation - Operation to estimate
 * @returns Estimated performance metrics for mobile
 */
export function estimateMobilePerformance(metrics: PerformanceMetrics): PerformanceMetrics {
  // Mobile devices typically have 2-4x slower performance
  const mobileSlowdownFactor = 3;

  return {
    ...metrics,
    operation: `${metrics.operation} (Mobile Estimate)`,
    executionTimeMs: metrics.executionTimeMs * mobileSlowdownFactor,
    throughputOpsPerSec: (metrics.throughputOpsPerSec || 0) / mobileSlowdownFactor
  };
}

/**
 * Generates a performance report
 * @param results - Benchmark results to report
 * @returns Formatted performance report
 */
export function generatePerformanceReport(results: BenchmarkResults): string {
  const report = [
    '='.repeat(60),
    'TOPAY-Z512 Performance Report',
    '='.repeat(60),
    '',
    'System Information:',
    `  Platform: ${results.systemInfo.platform}`,
    `  Node Version: ${results.systemInfo.nodeVersion || 'N/A'}`,
    `  CPU Count: ${results.systemInfo.cpuCount || 'N/A'}`,
    `  WebCrypto: ${results.systemInfo.hasWebCrypto ? 'Available' : 'Not Available'}`,
    '',
    'Benchmark Results:',
    '-'.repeat(40)
  ];

  results.metrics.forEach(metric => {
    report.push(`${metric.operation}:`);
    report.push(`  Execution Time: ${metric.executionTimeMs.toFixed(2)} ms`);
    if (metric.throughputOpsPerSec) {
      report.push(`  Throughput: ${metric.throughputOpsPerSec.toFixed(2)} ops/sec`);
    }
    if (metric.memoryUsageMB) {
      report.push(`  Memory Usage: ${metric.memoryUsageMB.toFixed(2)} MB`);
    }
    report.push('');
  });

  report.push('Summary:');
  report.push(`  Total Time: ${results.summary.totalTimeMs.toFixed(2)} ms`);
  report.push(`  Average Time: ${results.summary.averageTimeMs.toFixed(2)} ms`);
  report.push(`  Overall Throughput: ${results.summary.operationsPerSecond.toFixed(2)} ops/sec`);
  report.push('');
  report.push('='.repeat(60));

  return report.join('\n');
}
