/**
 * Optimized hash functions for TOPAY-Z512
 */

import { createHash, createHmac } from 'crypto';
import { HASH_SIZE } from './index';
import { Hash, EmptyDataError } from './index';
import { validateSize, concatBytes, constantTimeEqual } from './utils';

// Additional constants not exported from index
const SALT_SIZE = 32;
const HMAC_KEY_SIZE = 64;

// Hash cache for memoization of repeated operations
const hashCache = new Map<string, Uint8Array>();
const CACHE_MAX_SIZE = 1000;

/**
 * Clears the hash cache to prevent memory leaks
 */
export function clearHashCache(): void {
  hashCache.clear();
}

/**
 * Optimized hash computation with memoization
 * @param data - Data to hash
 * @param useCache - Whether to use caching (default: true for data < 1KB)
 * @returns Promise resolving to hash
 * @throws EmptyDataError if data is empty
 */
export async function computeHash(data: Uint8Array, useCache: boolean = data.length < 1024): Promise<Hash> {
  if (data.length === 0) {
    throw new EmptyDataError();
  }

  // Use cache for small, repeated computations
  if (useCache) {
    const key = Array.from(data).join(',');
    if (hashCache.has(key)) {
      // Create a proper copy to prevent corruption
      const cached = hashCache.get(key)!;
      return new Uint8Array(cached);
    }
    
    // Use SHA-512 as the base hash function for quantum resistance
    const hash = createHash('sha512');
    hash.update(data);
    const result = new Uint8Array(hash.digest());
    
    // Manage cache size
    if (hashCache.size >= CACHE_MAX_SIZE) {
      const firstKey = hashCache.keys().next().value;
      if (firstKey !== undefined) {
        hashCache.delete(firstKey);
      }
    }
    
    // Store a copy in cache to prevent corruption
    hashCache.set(key, new Uint8Array(result));
    validateSize(result, HASH_SIZE, 'hash');
    return result;
  }
  
  // Use SHA-512 as the base hash function for quantum resistance
  const hash = createHash('sha512');
  hash.update(data);
  const result = new Uint8Array(hash.digest());
  
  validateSize(result, HASH_SIZE, 'hash');
  return result;
}

/**
 * Computes a hash with a salt for additional security
 * @param data - Data to hash
 * @param salt - Salt to use
 * @returns Promise resolving to the salted hash
 * @throws EmptyDataError if data is empty
 */
export async function computeHashWithSalt(data: Uint8Array, salt: Uint8Array): Promise<Hash> {
  if (data.length === 0) {
    throw new EmptyDataError();
  }

  const combined = concatBytes(salt, data);
  return computeHash(combined);
}

/**
 * Computes HMAC of the input data with a key
 * @param key - HMAC key
 * @param data - Data to authenticate
 * @returns Promise resolving to the HMAC
 * @throws EmptyDataError if data is empty
 */
export async function computeHmac(key: Uint8Array, data: Uint8Array): Promise<Hash> {
  if (data.length === 0) {
    throw new EmptyDataError();
  }

  const hmac = createHmac('sha512', key);
  hmac.update(data);
  const result = new Uint8Array(hmac.digest());
  
  validateSize(result, HASH_SIZE, 'hmac');
  return result;
}

/**
 * Optimized batch hashing with controlled concurrency
 * @param dataItems - Array of data to hash
 * @param concurrency - Maximum concurrent operations (default: 8)
 * @returns Promise resolving to array of hashes
 */
export async function batchHash(dataItems: Uint8Array[], concurrency: number = 8): Promise<Hash[]> {
  if (dataItems.length === 0) return [];
  
  // Process in chunks to avoid overwhelming the system
  const results: Hash[] = [];
  
  for (let i = 0; i < dataItems.length; i += concurrency) {
    const chunk = dataItems.slice(i, i + concurrency);
    const chunkPromises = chunk.map(data => computeHash(data, data.length < 1024));
    const chunkResults = await Promise.all(chunkPromises);
    results.push(...chunkResults);
  }
  
  return results;
}

/**
 * Computes a Merkle tree root hash from an array of data
 * @param dataItems - Array of data items
 * @returns Promise resolving to the Merkle root hash
 * @throws EmptyDataError if no data items provided
 */
export async function computeMerkleRoot(dataItems: Uint8Array[]): Promise<Hash> {
  if (dataItems.length === 0) {
    throw new EmptyDataError();
  }

  // Compute leaf hashes
  let hashes = await batchHash(dataItems);

  // Build Merkle tree bottom-up
  while (hashes.length > 1) {
    const nextLevel: Hash[] = [];
    
    for (let i = 0; i < hashes.length; i += 2) {
      if (i + 1 < hashes.length) {
        // Pair exists, hash both
        const combined = concatBytes(hashes[i]!, hashes[i + 1]!);
        nextLevel.push(await computeHash(combined));
      } else {
        // Odd number, duplicate the last hash
        const combined = concatBytes(hashes[i]!, hashes[i]!);
        nextLevel.push(await computeHash(combined));
      }
    }
    
    hashes = nextLevel;
  }

  return hashes[0]!;
}

/**
 * Derives a key from a password using PBKDF2
 * @param password - Password to derive from
 * @param salt - Salt for key derivation
 * @param iterations - Number of iterations (default: 100000)
 * @param keyLength - Desired key length in bytes (default: HASH_SIZE)
 * @returns Promise resolving to the derived key
 */
export async function deriveKey(
  password: Uint8Array,
  salt: Uint8Array,
  iterations: number = 100000,
  keyLength: number = HASH_SIZE
): Promise<Uint8Array> {
  // Use Node.js crypto for PBKDF2
  const crypto = require('crypto');
  
  return new Promise((resolve, reject) => {
    crypto.pbkdf2(password, salt, iterations, keyLength, 'sha512', (err: Error | null, derivedKey: Buffer) => {
      if (err) {
        reject(err);
      } else {
        resolve(new Uint8Array(derivedKey));
      }
    });
  });
}

/**
 * Computes a hash chain of specified length
 * @param initialData - Initial data to start the chain
 * @param length - Length of the hash chain
 * @returns Promise resolving to array of hashes in the chain
 */
export async function computeHashChain(initialData: Uint8Array, length: number): Promise<Hash[]> {
  if (length <= 0) {
    throw new Error('Hash chain length must be positive');
  }

  const chain: Hash[] = [];
  let current = await computeHash(initialData);
  chain.push(current);

  for (let i = 1; i <= length; i++) {
    current = await computeHash(current);
    chain.push(current);
  }

  return chain;
}

/**
 * Verifies a hash chain
 * @param chain - Hash chain to verify
 * @param initialData - Initial data that started the chain
 * @returns Promise resolving to true if chain is valid
 */
export async function verifyHashChain(chain: Hash[], initialData: Uint8Array): Promise<boolean> {
  if (chain.length === 0) {
    return false;
  }

  // Verify first hash
  const expectedFirst = await computeHash(initialData);
  if (!constantTimeEqual(chain[0]!, expectedFirst)) {
    return false;
  }

  // Verify subsequent hashes
  for (let i = 1; i < chain.length; i++) {
    const expected = await computeHash(chain[i - 1]!);
    if (!constantTimeEqual(chain[i]!, expected)) {
      return false;
    }
  }

  return true;
}