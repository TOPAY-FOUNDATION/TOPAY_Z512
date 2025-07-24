/**
 * Utility functions for TOPAY-Z512
 */

import { createHash, randomBytes } from 'crypto';

/**
 * Generates cryptographically secure random bytes
 * @param size - Number of bytes to generate
 * @returns Promise resolving to random bytes
 */
export async function secureRandom(size: number): Promise<Uint8Array> {
  return new Promise((resolve, reject) => {
    try {
      const buffer = randomBytes(size);
      resolve(new Uint8Array(buffer));
    } catch (error) {
      reject(error);
    }
  });
}

/**
 * Optimized constant-time comparison with early termination prevention
 * @param a - First byte array
 * @param b - Second byte array
 * @returns True if arrays are equal, false otherwise
 */
export function constantTimeEqual(a: Uint8Array, b: Uint8Array): boolean {
  if (a.length !== b.length) {
    return false;
  }

  let result = 0;
  
  // Use 32-bit operations for better performance
  const len32 = Math.floor(a.length / 4) * 4;
  const view32A = new Uint32Array(a.buffer, a.byteOffset, len32 / 4);
  const view32B = new Uint32Array(b.buffer, b.byteOffset, len32 / 4);
  
  for (let i = 0; i < view32A.length; i++) {
    result |= view32A[i]! ^ view32B[i]!;
  }
  
  // Handle remaining bytes
  for (let i = len32; i < a.length; i++) {
    result |= a[i]! ^ b[i]!;
  }

  return result === 0;
}

/**
 * Securely zeros a byte array
 * @param data - Byte array to zero
 */
export function secureZero(data: Uint8Array): void {
  for (let i = 0; i < data.length; i++) {
    data[i] = 0;
  }
}

/**
 * Converts a byte array to hexadecimal string
 * @param data - Byte array to convert
 * @returns Hexadecimal string representation
 */
export function toHex(data: Uint8Array): string {
  return Array.from(data)
    .map(byte => byte.toString(16).padStart(2, '0'))
    .join('');
}

/**
 * Converts a hexadecimal string to byte array
 * @param hex - Hexadecimal string to convert
 * @returns Byte array representation
 * @throws Error if hex string is invalid
 */
export function fromHex(hex: string): Uint8Array {
  if (hex.length % 2 !== 0) {
    throw new Error('Invalid hex string length');
  }

  const result = new Uint8Array(hex.length / 2);
  for (let i = 0; i < hex.length; i += 2) {
    const byte = parseInt(hex.substr(i, 2), 16);
    if (isNaN(byte)) {
      throw new Error('Invalid hex character');
    }
    result[i / 2] = byte;
  }

  return result;
}

/**
 * Validates that a byte array has the expected size
 * @param data - Byte array to validate
 * @param expectedSize - Expected size in bytes
 * @param name - Name of the data for error messages
 * @throws Error if size is invalid
 */
export function validateSize(data: Uint8Array, expectedSize: number, name: string): void {
  if (data.length !== expectedSize) {
    throw new Error(`Invalid ${name} size: expected ${expectedSize}, got ${data.length}`);
  }
}

/**
 * Creates a deep copy of a byte array
 * @param data - Byte array to copy
 * @returns New byte array with copied data
 */
export function copyBytes(data: Uint8Array): Uint8Array {
  return new Uint8Array(data);
}

/**
 * Optimized concatenation of multiple byte arrays with pooled buffers
 * @param arrays - Arrays to concatenate
 * @returns Concatenated byte array
 */
export function concatBytes(...arrays: Uint8Array[]): Uint8Array {
  const totalLength = arrays.reduce((sum, arr) => sum + arr.length, 0);
  const result = getPooledBuffer(totalLength);
  
  let offset = 0;
  for (const array of arrays) {
    result.set(array, offset);
    offset += array.length;
  }
  
  return result;
}

/**
 * XORs two byte arrays of equal length
 * @param a - First byte array
 * @param b - Second byte array
 * @returns XOR result
 * @throws Error if arrays have different lengths
 */
// Buffer pool for memory optimization
const bufferPool = new Map<number, Uint8Array[]>();
const MAX_POOL_SIZE = 50;

/**
 * Gets a buffer from the pool or creates a new one
 * @param size - Size of buffer needed
 * @returns Uint8Array buffer
 */
function getPooledBuffer(size: number): Uint8Array {
  const pool = bufferPool.get(size);
  if (pool && pool.length > 0) {
    return pool.pop()!;
  }
  return new Uint8Array(size);
}

/**
 * Returns a buffer to the pool for reuse
 * @param buffer - Buffer to return to pool
 */
function returnToPool(buffer: Uint8Array): void {
  const size = buffer.length;
  let pool = bufferPool.get(size);
  
  if (!pool) {
    pool = [];
    bufferPool.set(size, pool);
  }
  
  if (pool.length < MAX_POOL_SIZE) {
    // Clear the buffer before returning to pool
    buffer.fill(0);
    pool.push(buffer);
  }
}

export function xorBytes(a: Uint8Array, b: Uint8Array): Uint8Array {
  if (a.length !== b.length) {
    throw new Error('Arrays must have the same length');
  }
  
  const result = getPooledBuffer(a.length);
  
  // Use simple byte-by-byte XOR to avoid alignment issues
  for (let i = 0; i < a.length; i++) {
    result[i] = a[i]! ^ b[i]!;
  }
  
  return result;
}

/**
 * Generates a timestamp in milliseconds
 * @returns Current timestamp
 */
export function timestamp(): number {
  return Date.now();
}

/**
 * Sleeps for the specified number of milliseconds
 * @param ms - Milliseconds to sleep
 * @returns Promise that resolves after the delay
 */
export function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

/**
 * Measures execution time of an async function
 * @param fn - Function to measure
 * @returns Object containing result and execution time
 */
export async function measureTime<T>(fn: () => Promise<T>): Promise<{ result: T; timeMs: number }> {
  const start = performance.now();
  const result = await fn();
  const timeMs = performance.now() - start;
  return { result, timeMs };
}

/**
 * Checks if the current environment supports WebCrypto
 * @returns True if WebCrypto is available
 */
export function hasWebCrypto(): boolean {
  return typeof crypto !== 'undefined' && typeof crypto.subtle !== 'undefined';
}

/**
 * Gets system capabilities and performance information
 * @returns System information object
 */
export function getSystemCapabilities(): {
  platform: string;
  nodeVersion?: string;
  hasWebCrypto: boolean;
  memoryUsage?: NodeJS.MemoryUsage;
  cpuCount?: number;
} {
  const capabilities: any = {
    platform: typeof process !== 'undefined' ? process.platform : 'browser',
    hasWebCrypto: hasWebCrypto(),
  };

  if (typeof process !== 'undefined') {
    capabilities.nodeVersion = process.version;
    capabilities.memoryUsage = process.memoryUsage();
    capabilities.cpuCount = require('os').cpus().length;
  }

  return capabilities;
}

/**
 * Clears buffer pools to free memory
 */
export function clearBufferPools(): void {
  bufferPool.clear();
}

/**
 * Gets buffer pool statistics
 * @returns Pool statistics
 */
export function getBufferPoolStats(): { totalPools: number; totalBuffers: number; sizes: number[] } {
  let totalBuffers = 0;
  const sizes: number[] = [];
  
  for (const [size, pool] of Array.from(bufferPool.entries())) {
    totalBuffers += pool.length;
    sizes.push(size);
  }
  
  return {
    totalPools: bufferPool.size,
    totalBuffers,
    sizes: sizes.sort((a, b) => a - b)
  };
}