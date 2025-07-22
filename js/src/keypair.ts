/**
 * Key pair generation and management for TOPAY-Z512
 */

import { 
  PRIVATE_KEY_SIZE, 
  PUBLIC_KEY_SIZE, 
  PrivateKey, 
  PublicKey, 
  KeyPair,
  InvalidKeySizeError 
} from './index';
import { secureRandom, validateSize, copyBytes, constantTimeEqual, secureZero } from './utils';
import { computeHash, deriveKey } from './hash';

/**
 * Generates a new cryptographic key pair
 * @returns Promise resolving to a new key pair
 */
export async function generateKeyPair(): Promise<KeyPair> {
  // Generate private key from secure random bytes
  const privateKey = await secureRandom(PRIVATE_KEY_SIZE);
  
  // Derive public key from private key using hash-based derivation
  const publicKey = await derivePublicKey(privateKey);
  
  return {
    privateKey,
    publicKey
  };
}

/**
 * Derives a public key from a private key
 * @param privateKey - Private key to derive from
 * @returns Promise resolving to the derived public key
 * @throws InvalidKeySizeError if private key size is invalid
 */
export async function derivePublicKey(privateKey: PrivateKey): Promise<PublicKey> {
  validateSize(privateKey, PRIVATE_KEY_SIZE, 'private key');
  
  // Use hash-based public key derivation for quantum safety
  const publicKey = await computeHash(privateKey);
  validateSize(publicKey, PUBLIC_KEY_SIZE, 'public key');
  
  return publicKey;
}

/**
 * Generates a key pair from a seed
 * @param seed - Seed bytes for deterministic generation
 * @returns Promise resolving to a key pair
 */
export async function generateKeyPairFromSeed(seed: Uint8Array): Promise<KeyPair> {
  if (seed.length < 32) {
    throw new Error('Seed must be at least 32 bytes');
  }
  
  // Derive private key from seed
  const privateKey = await deriveKey(seed, new Uint8Array(16), 10000, PRIVATE_KEY_SIZE);
  
  // Derive public key
  const publicKey = await derivePublicKey(privateKey);
  
  return {
    privateKey,
    publicKey
  };
}

/**
 * Generates multiple key pairs in batch
 * @param count - Number of key pairs to generate
 * @returns Promise resolving to array of key pairs
 */
export async function batchGenerateKeyPairs(count: number): Promise<KeyPair[]> {
  if (count <= 0) {
    throw new Error('Count must be positive');
  }
  
  const promises: Promise<KeyPair>[] = [];
  for (let i = 0; i < count; i++) {
    promises.push(generateKeyPair());
  }
  
  return Promise.all(promises);
}

/**
 * Validates a key pair by checking if public key matches private key
 * @param keyPair - Key pair to validate
 * @returns Promise resolving to true if valid
 */
export async function validateKeyPair(keyPair: KeyPair): Promise<boolean> {
  try {
    validateSize(keyPair.privateKey, PRIVATE_KEY_SIZE, 'private key');
    validateSize(keyPair.publicKey, PUBLIC_KEY_SIZE, 'public key');
    
    const derivedPublic = await derivePublicKey(keyPair.privateKey);
    return constantTimeEqual(keyPair.publicKey, derivedPublic);
  } catch {
    return false;
  }
}

/**
 * Derives a child key pair from a parent private key and index
 * @param parentPrivateKey - Parent private key
 * @param index - Child key index
 * @returns Promise resolving to child key pair
 */
export async function deriveChildKeyPair(parentPrivateKey: PrivateKey, index: number): Promise<KeyPair> {
  validateSize(parentPrivateKey, PRIVATE_KEY_SIZE, 'parent private key');
  
  // Create derivation data
  const indexBytes = new Uint8Array(4);
  new DataView(indexBytes.buffer).setUint32(0, index, false);
  
  // Combine parent key and index
  const derivationData = new Uint8Array(parentPrivateKey.length + indexBytes.length);
  derivationData.set(parentPrivateKey);
  derivationData.set(indexBytes, parentPrivateKey.length);
  
  // Derive child private key
  const childPrivateKey = await computeHash(derivationData);
  validateSize(childPrivateKey, PRIVATE_KEY_SIZE, 'child private key');
  
  // Derive child public key
  const childPublicKey = await derivePublicKey(childPrivateKey);
  
  return {
    privateKey: childPrivateKey,
    publicKey: childPublicKey
  };
}

/**
 * Generates an HD (Hierarchical Deterministic) wallet
 * @param masterSeed - Master seed for the wallet
 * @param count - Number of key pairs to generate
 * @returns Promise resolving to array of key pairs
 */
export async function generateHDWallet(masterSeed: Uint8Array, count: number): Promise<KeyPair[]> {
  if (count <= 0) {
    throw new Error('Count must be positive');
  }
  
  // Generate master key pair
  const masterKeyPair = await generateKeyPairFromSeed(masterSeed);
  
  // Generate child key pairs
  const keyPairs: KeyPair[] = [masterKeyPair];
  
  for (let i = 1; i < count; i++) {
    const childKeyPair = await deriveChildKeyPair(masterKeyPair.privateKey, i);
    keyPairs.push(childKeyPair);
  }
  
  return keyPairs;
}

/**
 * Derives a key pair from a password
 * @param password - Password string
 * @param salt - Salt for key derivation
 * @param iterations - PBKDF2 iterations (default: 100000)
 * @returns Promise resolving to key pair
 */
export async function deriveKeyPairFromPassword(
  password: string,
  salt: Uint8Array,
  iterations: number = 100000
): Promise<KeyPair> {
  const passwordBytes = new TextEncoder().encode(password);
  const privateKey = await deriveKey(passwordBytes, salt, iterations, PRIVATE_KEY_SIZE);
  const publicKey = await derivePublicKey(privateKey);
  
  return {
    privateKey,
    publicKey
  };
}

/**
 * Securely erases a key pair from memory
 * @param keyPair - Key pair to erase
 */
export function secureEraseKeyPair(keyPair: KeyPair): void {
  secureZero(keyPair.privateKey);
  secureZero(keyPair.publicKey);
}

/**
 * Creates a backup of a key pair (copies the data)
 * @param keyPair - Key pair to backup
 * @returns New key pair with copied data
 */
export function backupKeyPair(keyPair: KeyPair): KeyPair {
  return {
    privateKey: copyBytes(keyPair.privateKey),
    publicKey: copyBytes(keyPair.publicKey)
  };
}

/**
 * Serializes a key pair to a portable format
 * @param keyPair - Key pair to serialize
 * @returns Serialized key pair as JSON string
 */
export function serializeKeyPair(keyPair: KeyPair): string {
  return JSON.stringify({
    privateKey: Array.from(keyPair.privateKey),
    publicKey: Array.from(keyPair.publicKey)
  });
}

/**
 * Deserializes a key pair from a portable format
 * @param serialized - Serialized key pair JSON string
 * @returns Deserialized key pair
 * @throws Error if deserialization fails
 */
export function deserializeKeyPair(serialized: string): KeyPair {
  try {
    const data = JSON.parse(serialized);
    
    if (!data.privateKey || !data.publicKey) {
      throw new Error('Missing key data');
    }
    
    const privateKey = new Uint8Array(data.privateKey);
    const publicKey = new Uint8Array(data.publicKey);
    
    validateSize(privateKey, PRIVATE_KEY_SIZE, 'private key');
    validateSize(publicKey, PUBLIC_KEY_SIZE, 'public key');
    
    return { privateKey, publicKey };
  } catch (error) {
    throw new Error(`Failed to deserialize key pair: ${error}`);
  }
}