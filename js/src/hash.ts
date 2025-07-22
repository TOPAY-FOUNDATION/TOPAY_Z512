/**
 * TOPAY-Z512 Hash implementation
 * 
 * This module provides a 512-bit cryptographic hash function implementation
 * that can be used as part of the TOPAY-Z512 cryptographic suite.
 */

import { sha3_512 } from 'js-sha3';

/** The size of TOPAY-Z512 hash output in bytes (512 bits = 64 bytes) */
export const HASH_SIZE_BYTES = 64;

/**
 * Represents a TOPAY-Z512 hash value (512 bits)
 */
export class Hash {
  private readonly bytes: Uint8Array;

  /**
   * Creates a new Hash instance
   * @param bytes - The 64-byte hash value
   */
  constructor(bytes: Uint8Array) {
    if (bytes.length !== HASH_SIZE_BYTES) {
      throw new Error(`Invalid hash length: ${bytes.length}, expected ${HASH_SIZE_BYTES}`);
    }
    this.bytes = new Uint8Array(bytes);
  }

  /**
   * Creates a new hash from the given data
   * @param data - The data to hash
   * @returns A new Hash instance
   */
  static new(data: Uint8Array | string): Hash {
    const inputData = typeof data === 'string' ? new TextEncoder().encode(data) : data;
    const hashValue = sha3_512.arrayBuffer(inputData);
    return new Hash(new Uint8Array(hashValue));
  }

  /**
   * Creates a new hash by combining two input values
   * @param data1 - The first data to hash
   * @param data2 - The second data to hash
   * @returns A new Hash instance
   */
  static combine(data1: Uint8Array | string, data2: Uint8Array | string): Hash {
    const inputData1 = typeof data1 === 'string' ? new TextEncoder().encode(data1) : data1;
    const inputData2 = typeof data2 === 'string' ? new TextEncoder().encode(data2) : data2;
    
    // First hash the individual inputs
    const hash1 = Hash.new(inputData1);
    const hash2 = Hash.new(inputData2);
    
    // Then combine the hashes with a separator to ensure it's different from concatenation
    const combinedBytes = new Uint8Array(hash1.getBytes().length + hash2.getBytes().length + 1);
    combinedBytes.set(hash1.getBytes(), 0);
    combinedBytes[hash1.getBytes().length] = 0xFF; // Add a separator byte
    combinedBytes.set(hash2.getBytes(), hash1.getBytes().length + 1);
    
    // Hash the combined result
    return Hash.new(combinedBytes);
  }

  /**
   * Returns the hash value as a byte array
   * @returns The hash bytes
   */
  getBytes(): Uint8Array {
    return new Uint8Array(this.bytes);
  }

  /**
   * Converts the hash to a hexadecimal string
   * @returns The hash as a hex string
   */
  toHex(): string {
    return Array.from(this.bytes)
      .map(b => b.toString(16).padStart(2, '0'))
      .join('');
  }

  /**
   * Creates a hash from a hexadecimal string
   * @param hex - The hex string to convert
   * @returns A new Hash instance
   */
  static fromHex(hex: string): Hash {
    if (hex.length !== HASH_SIZE_BYTES * 2) {
      throw new Error(`Invalid hex length: ${hex.length}, expected ${HASH_SIZE_BYTES * 2}`);
    }

    if (!/^[0-9a-fA-F]+$/.test(hex)) {
      throw new Error('Invalid hex string: contains non-hex characters');
    }

    const bytes = new Uint8Array(HASH_SIZE_BYTES);
    for (let i = 0; i < HASH_SIZE_BYTES; i++) {
      bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
    }

    return new Hash(bytes);
  }

  /**
   * Compares this hash with another hash for equality
   * @param other - The other hash to compare with
   * @returns True if the hashes are equal, false otherwise
   */
  equals(other: Hash): boolean {
    if (this.bytes.length !== other.bytes.length) {
      return false;
    }

    for (let i = 0; i < this.bytes.length; i++) {
      if (this.bytes[i] !== other.bytes[i]) {
        return false;
      }
    }

    return true;
  }

  /**
   * Returns a string representation of the hash
   * @returns The hash as a hex string
   */
  toString(): string {
    return this.toHex();
  }

  /**
   * Creates a new hash using the current time as input
   * This is useful for generating random-like hashes when no specific input is available
   * @returns A new Hash instance
   */
  static newWithTime(): Hash {
    const timeData = Date.now().toString();
    return Hash.new(timeData);
  }
}

/**
 * A convenience function to hash data
 * @param data - The data to hash
 * @returns The hash bytes
 */
export function hash(data: Uint8Array | string): Uint8Array {
  return Hash.new(data).getBytes();
}

/**
 * A convenience function to hash two pieces of data together
 * @param data1 - The first data to hash
 * @param data2 - The second data to hash
 * @returns The hash bytes
 */
export function hashCombine(data1: Uint8Array | string, data2: Uint8Array | string): Uint8Array {
  return Hash.combine(data1, data2).getBytes();
}

/**
 * Creates a new hash using the current time as input
 * This is useful for generating random-like hashes when no specific input is available
 * @returns A new Hash instance
 */
export function newWithTime(): Hash {
  return Hash.newWithTime();
}

/**
 * A convenience function to hash using the current time
 * @returns The hash bytes
 */
export function hashWithTime(): Uint8Array {
  return newWithTime().getBytes();
}