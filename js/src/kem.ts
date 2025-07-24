/**
 * Key Encapsulation Mechanism (KEM) for TOPAY-Z512
 */

import {
  KEM_PUBLIC_KEY_SIZE,
  KEM_SECRET_KEY_SIZE,
  CIPHERTEXT_SIZE,
  // SHARED_SECRET_SIZE,
  KEM_SHARED_SECRET_SIZE,
  KEMKeyPair,
  KEMPublicKey,
  KEMSecretKey,
  Ciphertext,
  SharedSecret,
  EncapsulationResult,
  DecapsulationFailedError
} from './index';
import {
  secureRandom,
  validateSize,
  constantTimeEqual,
  secureZero,
  copyBytes,
  xorBytes
} from './utils';
import { computeHash /* , computeHmac */ } from './hash';

// Key pair cache for improved performance
const keyPairCache = new Map<
  string,
  { publicKey: Uint8Array; secretKey: Uint8Array; timestamp: number }
>();
const KEY_CACHE_TTL = 300000; // 5 minutes
const KEY_CACHE_MAX_SIZE = 100;

/**
 * Clears expired entries from the key pair cache
 */
function cleanupKeyCache(): void {
  const now = Date.now();
  for (const [key, value] of Array.from(keyPairCache.entries())) {
    if (now - value.timestamp > KEY_CACHE_TTL) {
      keyPairCache.delete(key);
    }
  }
}

/**
 * Optimized KEM key generation with optional caching
 * @param useCache - Whether to use caching (default: false for security)
 * @returns Promise resolving to a new KEM key pair
 */
export async function kemKeyGen(useCache: boolean = false): Promise<KEMKeyPair> {
  if (useCache) {
    // Clean up expired entries periodically
    if (keyPairCache.size > 0 && Math.random() < 0.1) {
      cleanupKeyCache();
    }

    // For demonstration purposes, we'll cache based on a simple key
    // In production, this should be more sophisticated
    const cacheKey = 'default';
    const cached = keyPairCache.get(cacheKey);

    if (cached && Date.now() - cached.timestamp < KEY_CACHE_TTL) {
      return {
        publicKey: new Uint8Array(cached.publicKey),
        secretKey: new Uint8Array(cached.secretKey)
      };
    }
  }

  // Generate secret key from secure random bytes
  const secretKey = await secureRandom(KEM_SECRET_KEY_SIZE);

  // Derive public key from secret key
  const publicKey = await deriveKEMPublicKey(secretKey);

  if (useCache) {
    // Manage cache size
    if (keyPairCache.size >= KEY_CACHE_MAX_SIZE) {
      const firstKey = keyPairCache.keys().next().value;
      if (firstKey !== undefined) {
        keyPairCache.delete(firstKey);
      }
    }

    keyPairCache.set('default', {
      publicKey: new Uint8Array(publicKey),
      secretKey: new Uint8Array(secretKey),
      timestamp: Date.now()
    });
  }

  return {
    publicKey,
    secretKey
  };
}

/**
 * Derives a KEM public key from a secret key
 * @param secretKey - Secret key to derive from
 * @returns Promise resolving to the derived public key
 * @throws InvalidKeySizeError if secret key size is invalid
 */
export async function deriveKEMPublicKey(secretKey: KEMSecretKey): Promise<KEMPublicKey> {
  validateSize(secretKey, KEM_SECRET_KEY_SIZE, 'KEM secret key');

  // Use hash-based public key derivation
  const publicKey = await computeHash(secretKey);
  validateSize(publicKey, KEM_PUBLIC_KEY_SIZE, 'KEM public key');

  return publicKey;
}

/**
 * Encapsulates a shared secret using a public key
 * @param publicKey - KEM public key for encapsulation
 * @returns Promise resolving to encapsulation result (ciphertext and shared secret)
 * @throws InvalidKeySizeError if public key size is invalid
 */
export async function kemEncapsulate(publicKey: KEMPublicKey): Promise<EncapsulationResult> {
  validateSize(publicKey, KEM_PUBLIC_KEY_SIZE, 'KEM public key');

  // Generate ephemeral key material
  const ephemeralKey = await secureRandom(32);

  // Derive shared secret using key derivation
  const keyMaterial = new Uint8Array(publicKey.length + ephemeralKey.length);
  keyMaterial.set(publicKey);
  keyMaterial.set(ephemeralKey, publicKey.length);

  const fullHash = await computeHash(keyMaterial);
  // Truncate to the expected KEM shared secret size
  const sharedSecret = fullHash.slice(0, KEM_SHARED_SECRET_SIZE);
  validateSize(sharedSecret, KEM_SHARED_SECRET_SIZE, 'shared secret');

  // Create ciphertext by encrypting ephemeral key with public key
  const ciphertext = await encryptEphemeralKey(ephemeralKey, publicKey);
  validateSize(ciphertext, CIPHERTEXT_SIZE, 'ciphertext');

  return {
    ciphertext,
    sharedSecret
  };
}

/**
 * Decapsulates a shared secret using a secret key and ciphertext
 * @param secretKey - KEM secret key for decapsulation
 * @param ciphertext - Ciphertext to decapsulate
 * @returns Promise resolving to the shared secret
 * @throws InvalidKeySizeError if key or ciphertext size is invalid
 * @throws DecapsulationFailedError if decapsulation fails
 */
export async function kemDecapsulate(
  secretKey: KEMSecretKey,
  ciphertext: Ciphertext
): Promise<SharedSecret> {
  validateSize(secretKey, KEM_SECRET_KEY_SIZE, 'KEM secret key');
  validateSize(ciphertext, CIPHERTEXT_SIZE, 'ciphertext');

  try {
    // Derive public key from secret key
    const publicKey = await deriveKEMPublicKey(secretKey);

    // Decrypt ephemeral key from ciphertext
    const ephemeralKey = await decryptEphemeralKey(ciphertext, secretKey);

    // Derive shared secret using the same process as encapsulation
    const keyMaterial = new Uint8Array(publicKey.length + ephemeralKey.length);
    keyMaterial.set(publicKey);
    keyMaterial.set(ephemeralKey, publicKey.length);

    const fullHash = await computeHash(keyMaterial);
    // Truncate to the expected KEM shared secret size
    const sharedSecret = fullHash.slice(0, KEM_SHARED_SECRET_SIZE);
    validateSize(sharedSecret, KEM_SHARED_SECRET_SIZE, 'shared secret');

    return sharedSecret;
  } catch (error) {
    throw new DecapsulationFailedError();
  }
}

/**
 * Performs batch KEM key generation
 * @param count - Number of key pairs to generate
 * @returns Promise resolving to array of KEM key pairs
 */
export async function batchKEMKeyGen(count: number): Promise<KEMKeyPair[]> {
  if (count <= 0) {
    throw new Error('Count must be positive');
  }

  const promises: Promise<KEMKeyPair>[] = [];
  for (let i = 0; i < count; i++) {
    promises.push(kemKeyGen());
  }

  return Promise.all(promises);
}

/**
 * Performs batch KEM encapsulation
 * @param publicKeys - Array of public keys for encapsulation
 * @returns Promise resolving to array of encapsulation results
 */
export async function batchKEMEncapsulate(
  publicKeys: KEMPublicKey[]
): Promise<EncapsulationResult[]> {
  const promises = publicKeys.map(publicKey => kemEncapsulate(publicKey));
  return Promise.all(promises);
}

/**
 * Performs batch KEM decapsulation
 * @param operations - Array of {secretKey, ciphertext} pairs
 * @returns Promise resolving to array of shared secrets
 */
export async function batchKEMDecapsulate(
  operations: Array<{ secretKey: KEMSecretKey; ciphertext: Ciphertext }>
): Promise<SharedSecret[]> {
  const promises = operations.map(({ secretKey, ciphertext }) =>
    kemDecapsulate(secretKey, ciphertext)
  );
  return Promise.all(promises);
}

/**
 * Validates a KEM key pair
 * @param keyPair - KEM key pair to validate
 * @returns Promise resolving to true if valid
 */
export async function validateKEMKeyPair(keyPair: KEMKeyPair): Promise<boolean> {
  try {
    validateSize(keyPair.secretKey, KEM_SECRET_KEY_SIZE, 'KEM secret key');
    validateSize(keyPair.publicKey, KEM_PUBLIC_KEY_SIZE, 'KEM public key');

    const derivedPublic = await deriveKEMPublicKey(keyPair.secretKey);
    return constantTimeEqual(keyPair.publicKey, derivedPublic);
  } catch {
    return false;
  }
}

/**
 * Tests KEM operations with a key pair
 * @param keyPair - KEM key pair to test
 * @returns Promise resolving to true if operations succeed
 */
export async function testKEMOperations(keyPair: KEMKeyPair): Promise<boolean> {
  try {
    // Test encapsulation
    const { ciphertext, sharedSecret } = await kemEncapsulate(keyPair.publicKey);

    // Test decapsulation
    const decapsulatedSecret = await kemDecapsulate(keyPair.secretKey, ciphertext);

    // Verify shared secrets match
    return constantTimeEqual(sharedSecret, decapsulatedSecret);
  } catch {
    return false;
  }
}

/**
 * Securely erases a KEM key pair from memory
 * @param keyPair - KEM key pair to erase
 */
export function secureEraseKEMKeyPair(keyPair: KEMKeyPair): void {
  secureZero(keyPair.secretKey);
  secureZero(keyPair.publicKey);
}

/**
 * Creates a backup of a KEM key pair
 * @param keyPair - KEM key pair to backup
 * @returns New KEM key pair with copied data
 */
export function backupKEMKeyPair(keyPair: KEMKeyPair): KEMKeyPair {
  return {
    publicKey: copyBytes(keyPair.publicKey),
    secretKey: copyBytes(keyPair.secretKey)
  };
}

/**
 * Serializes a KEM key pair to JSON
 * @param keyPair - KEM key pair to serialize
 * @returns Serialized key pair as JSON string
 */
export function serializeKEMKeyPair(keyPair: KEMKeyPair): string {
  return JSON.stringify({
    publicKey: Array.from(keyPair.publicKey),
    secretKey: Array.from(keyPair.secretKey)
  });
}

/**
 * Deserializes a KEM key pair from JSON
 * @param serialized - Serialized key pair JSON string
 * @returns Deserialized KEM key pair
 * @throws Error if deserialization fails
 */
export function deserializeKEMKeyPair(serialized: string): KEMKeyPair {
  try {
    const data = JSON.parse(serialized);

    if (!data.publicKey || !data.secretKey) {
      throw new Error('Missing key data');
    }

    const publicKey = new Uint8Array(data.publicKey);
    const secretKey = new Uint8Array(data.secretKey);

    validateSize(publicKey, KEM_PUBLIC_KEY_SIZE, 'KEM public key');
    validateSize(secretKey, KEM_SECRET_KEY_SIZE, 'KEM secret key');

    return { publicKey, secretKey };
  } catch (error) {
    throw new Error(`Failed to deserialize KEM key pair: ${error}`);
  }
}

// Helper functions for encryption/decryption of ephemeral keys

/**
 * Encrypts ephemeral key material with public key
 * @param ephemeralKey - Ephemeral key to encrypt
 * @param publicKey - Public key for encryption
 * @returns Promise resolving to ciphertext
 */
async function encryptEphemeralKey(
  ephemeralKey: Uint8Array,
  publicKey: KEMPublicKey
): Promise<Ciphertext> {
  // Simple XOR-based encryption for demonstration
  // In production, use proper lattice-based encryption
  const keyStream = await computeHash(publicKey);
  const encrypted = xorBytes(ephemeralKey, keyStream.slice(0, ephemeralKey.length));

  // Pad to ciphertext size
  const ciphertext = new Uint8Array(CIPHERTEXT_SIZE);
  ciphertext.set(encrypted);

  return ciphertext;
}

/**
 * Decrypts ephemeral key material with secret key
 * @param ciphertext - Ciphertext to decrypt
 * @param secretKey - Secret key for decryption
 * @returns Promise resolving to ephemeral key
 */
async function decryptEphemeralKey(
  ciphertext: Ciphertext,
  secretKey: KEMSecretKey
): Promise<Uint8Array> {
  // Derive public key and decrypt
  const publicKey = await deriveKEMPublicKey(secretKey);
  const keyStream = await computeHash(publicKey);

  const encryptedPart = ciphertext.slice(0, 32); // Extract encrypted ephemeral key
  const ephemeralKey = xorBytes(encryptedPart, keyStream.slice(0, 32));

  return ephemeralKey;
}
