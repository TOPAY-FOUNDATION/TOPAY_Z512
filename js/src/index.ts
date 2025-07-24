/**
 * TOPAY-Z512 Cryptographic Library
 *
 * Quantum-safe cryptographic primitives for the TOPAY Foundation blockchain ecosystem.
 *
 * TOPAY-Z512 offers post-quantum security with ≥512-bit classical security (~256-bit quantum resistance)
 * using lattice-based cryptography. It includes support for key generation, hashing, key encapsulation
 * mechanisms (KEM), and fragmented processing for parallel computation.
 *
 * Key Features:
 * - Quantum-safe cryptographic operations
 * - High-performance implementations with batch operations
 * - Fragmented architecture for mobile and IoT device support
 * - Cross-platform compatibility
 * - Developer-friendly API design
 *
 * @example
 * ```typescript
 * import { generateKeyPair, computeHash, kemKeyGen, kemEncapsulate, kemDecapsulate } from '@topayfoundation/topayz512';
 *
 * // Generate a key pair
 * const keyPair = await generateKeyPair();
 *
 * // Hash data
 * const data = new Uint8Array([1, 2, 3, 4, 5]);
 * const hash = await computeHash(data);
 *
 * // KEM operations
 * const kemKeyPair = await kemKeyGen();
 * const { ciphertext, sharedSecret } = await kemEncapsulate(kemKeyPair.publicKey);
 * const decapsulatedSecret = await kemDecapsulate(kemKeyPair.secretKey, ciphertext);
 * ```
 *
 * @packageDocumentation
 */

// Version information
export const VERSION = '0.1.0';
export const SECURITY_LEVEL = 512;
export const QUANTUM_SECURITY_LEVEL = 256;

// Cryptographic constants
export const PRIVATE_KEY_SIZE = 64;
export const PUBLIC_KEY_SIZE = 64;
export const HASH_SIZE = 64;
export const KEM_PUBLIC_KEY_SIZE = 64;
export const KEM_SECRET_KEY_SIZE = 64;
export const KEM_CIPHERTEXT_SIZE = 64;
export const KEM_SHARED_SECRET_SIZE = 32;
export const CIPHERTEXT_SIZE = 64;
export const SHARED_SECRET_SIZE = 64;

// Performance constants
export const DEFAULT_BATCH_SIZE = 32;
export const CACHE_LINE_SIZE = 64;
export const SIMD_WIDTH = 256;
export const PREFETCH_DISTANCE = 64;

// Fragmentation constants
export const DEFAULT_FRAGMENT_SIZE = 1024;
export const MAX_FRAGMENT_SIZE = 65536;
export const MIN_FRAGMENT_SIZE = 256;
export const FRAGMENT_SIZE = 256;
export const MIN_FRAGMENT_THRESHOLD = 512;
export const MAX_FRAGMENTS = 1024;

// Core cryptographic types
export type PrivateKey = Uint8Array;
export type PublicKey = Uint8Array;
export type Hash = Uint8Array;
export type KEMPublicKey = Uint8Array;
export type KEMSecretKey = Uint8Array;
export type Ciphertext = Uint8Array;
export type SharedSecret = Uint8Array;

/**
 * Represents a complete key pair
 */
export interface KeyPair {
  privateKey: PrivateKey;
  publicKey: PublicKey;
}

/**
 * Represents a complete KEM key pair
 */
export interface KEMKeyPair {
  publicKey: KEMPublicKey;
  secretKey: KEMSecretKey;
}

/**
 * Represents the result of KEM encapsulation
 */
export interface EncapsulationResult {
  ciphertext: Ciphertext;
  sharedSecret: SharedSecret;
}

/**
 * Represents metadata for a fragment
 */
export interface FragmentMetadata {
  originalSize: number;
  fragmentCount: number;
  timestamp: number;
  algorithm: string;
  checksum: Hash;
}

/**
 * Represents a single fragment
 */
export interface Fragment {
  index: number;
  data: Uint8Array;
  metadata: FragmentMetadata;
}

/**
 * Represents the result of fragmentation
 */
export interface FragmentationResult {
  fragments: Fragment[];
  metadata: FragmentMetadata;
}

/**
 * Represents the result of reconstruction
 */
export interface ReconstructionResult {
  data: Uint8Array;
  isComplete: boolean;
  missingCount: number;
  metadata: FragmentMetadata;
}

/**
 * Custom error types for TOPAY-Z512
 */
export class TopayZ512Error extends Error {
  constructor(
    message: string,
    public readonly code: string
  ) {
    super(message);
    this.name = 'TopayZ512Error';
  }
}

export class InvalidKeySizeError extends TopayZ512Error {
  constructor() {
    super('Invalid key size', 'INVALID_KEY_SIZE');
  }
}

export class InvalidHashSizeError extends TopayZ512Error {
  constructor() {
    super('Invalid hash size', 'INVALID_HASH_SIZE');
  }
}

export class InvalidCiphertextSizeError extends TopayZ512Error {
  constructor() {
    super('Invalid ciphertext size', 'INVALID_CIPHERTEXT_SIZE');
  }
}

export class DecapsulationFailedError extends TopayZ512Error {
  constructor() {
    super('Decapsulation failed', 'DECAPSULATION_FAILED');
  }
}

export class FragmentationFailedError extends TopayZ512Error {
  constructor() {
    super('Fragmentation failed', 'FRAGMENTATION_FAILED');
  }
}

export class ReconstructionFailedError extends TopayZ512Error {
  constructor() {
    super('Reconstruction failed', 'RECONSTRUCTION_FAILED');
  }
}

export class EmptyDataError extends TopayZ512Error {
  constructor() {
    super('Empty data provided', 'EMPTY_DATA');
  }
}

export class InvalidFragmentCountError extends TopayZ512Error {
  constructor() {
    super('Invalid fragment count', 'INVALID_FRAGMENT_COUNT');
  }
}

// Re-export all functionality
export * from './hash';
export * from './keypair';
export * from './kem';
export * from './fragment';
export * from './utils';
export * from './performance';
