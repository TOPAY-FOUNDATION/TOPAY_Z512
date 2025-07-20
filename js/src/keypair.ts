/**
 * TOPAY-Z512 Key Pair implementation
 * 
 * This module provides a 512-bit cryptographic key pair implementation
 * that can be used as part of the TOPAY-Z512 cryptographic suite.
 */

import { sha3_512 } from 'js-sha3';
import { HASH_SIZE_BYTES } from './hash';

/** The size of TOPAY-Z512 keys in bytes (512 bits = 64 bytes) */
export const KEY_SIZE_BYTES = HASH_SIZE_BYTES;

/**
 * Represents a TOPAY-Z512 private key (512 bits)
 */
export class PrivateKey {
  private readonly bytes: Uint8Array;

  /**
   * Creates a new PrivateKey instance
   * @param bytes - The 64-byte private key value
   */
  constructor(bytes: Uint8Array) {
    if (bytes.length !== KEY_SIZE_BYTES) {
      throw new Error(`Invalid private key length: ${bytes.length}, expected ${KEY_SIZE_BYTES}`);
    }
    this.bytes = new Uint8Array(bytes);
  }

  /**
   * Generates a new private key using secure random data
   * @returns A new PrivateKey instance
   */
  static generate(): PrivateKey {
    const bytes = new Uint8Array(KEY_SIZE_BYTES);
    crypto.getRandomValues(bytes);
    return new PrivateKey(bytes);
  }

  /**
   * Returns the private key as a byte array
   * @returns The private key bytes
   */
  getBytes(): Uint8Array {
    return new Uint8Array(this.bytes);
  }

  /**
   * Converts the private key to a hexadecimal string
   * @returns The private key as a hex string
   */
  toHex(): string {
    return Array.from(this.bytes)
      .map(b => b.toString(16).padStart(2, '0'))
      .join('');
  }

  /**
   * Creates a private key from a hexadecimal string
   * @param hex - The hex string to convert
   * @returns A new PrivateKey instance
   */
  static fromHex(hex: string): PrivateKey {
    if (hex.length !== KEY_SIZE_BYTES * 2) {
      throw new Error(`Invalid hex length: ${hex.length}, expected ${KEY_SIZE_BYTES * 2}`);
    }

    if (!/^[0-9a-fA-F]+$/.test(hex)) {
      throw new Error('Invalid hex string: contains non-hex characters');
    }

    const bytes = new Uint8Array(KEY_SIZE_BYTES);
    for (let i = 0; i < KEY_SIZE_BYTES; i++) {
      bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
    }

    return new PrivateKey(bytes);
  }
}

/**
 * Represents a TOPAY-Z512 public key (512 bits)
 */
export class PublicKey {
  private readonly bytes: Uint8Array;

  /**
   * Creates a new PublicKey instance
   * @param bytes - The 64-byte public key value
   */
  constructor(bytes: Uint8Array) {
    if (bytes.length !== KEY_SIZE_BYTES) {
      throw new Error(`Invalid public key length: ${bytes.length}, expected ${KEY_SIZE_BYTES}`);
    }
    this.bytes = new Uint8Array(bytes);
  }

  /**
   * Derives a public key from a private key
   * @param privateKey - The private key to derive from
   * @returns A new PublicKey instance
   */
  static fromPrivateKey(privateKey: PrivateKey): PublicKey {
    const privateKeyBytes = privateKey.getBytes();
    const hashValue = sha3_512.arrayBuffer(privateKeyBytes);
    return new PublicKey(new Uint8Array(hashValue));
  }

  /**
   * Returns the public key as a byte array
   * @returns The public key bytes
   */
  getBytes(): Uint8Array {
    return new Uint8Array(this.bytes);
  }

  /**
   * Converts the public key to a hexadecimal string
   * @returns The public key as a hex string
   */
  toHex(): string {
    return Array.from(this.bytes)
      .map(b => b.toString(16).padStart(2, '0'))
      .join('');
  }

  /**
   * Creates a public key from a hexadecimal string
   * @param hex - The hex string to convert
   * @returns A new PublicKey instance
   */
  static fromHex(hex: string): PublicKey {
    if (hex.length !== KEY_SIZE_BYTES * 2) {
      throw new Error(`Invalid hex length: ${hex.length}, expected ${KEY_SIZE_BYTES * 2}`);
    }

    if (!/^[0-9a-fA-F]+$/.test(hex)) {
      throw new Error('Invalid hex string: contains non-hex characters');
    }

    const bytes = new Uint8Array(KEY_SIZE_BYTES);
    for (let i = 0; i < KEY_SIZE_BYTES; i++) {
      bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
    }

    return new PublicKey(bytes);
  }

  /**
   * Compares this public key with another public key for equality
   * @param other - The other public key to compare with
   * @returns True if the public keys are equal, false otherwise
   */
  equals(other: PublicKey): boolean {
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
}

/**
 * Represents a TOPAY-Z512 key pair (private key and public key)
 */
export class KeyPair {
  /**
   * Creates a new KeyPair instance
   * @param privateKey - The private key
   * @param publicKey - The public key
   */
  constructor(
    public readonly privateKey: PrivateKey,
    public readonly publicKey: PublicKey
  ) {}

  /**
   * Generates a new key pair
   * @returns A new KeyPair instance
   */
  static generate(): KeyPair {
    const privateKey = PrivateKey.generate();
    const publicKey = PublicKey.fromPrivateKey(privateKey);
    return new KeyPair(privateKey, publicKey);
  }
}

/**
 * Convenience function to generate a key pair
 * @returns A new KeyPair instance
 */
export function generateKeyPair(): KeyPair {
  return KeyPair.generate();
}

/**
 * Convenience function to derive a public key from a private key
 * @param privateKey - The private key to derive from
 * @returns A new PublicKey instance
 */
export function privateToPublic(privateKey: PrivateKey): PublicKey {
  return PublicKey.fromPrivateKey(privateKey);
}