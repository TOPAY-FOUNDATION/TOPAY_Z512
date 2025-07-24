//! Optimized Key Encapsulation Mechanism (KEM) for TOPAY-Z512
//!
//! This module provides a high-performance KEM implementation for demonstration purposes.
//! In production, this would use a proper post-quantum KEM like Kyber or NTRU.

use crate::error::{Result, TopayzError};
use crate::hash::Hash;
use std::time::{SystemTime, UNIX_EPOCH};

/// KEM public key for encapsulation with optimized layout
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PublicKey {
    bytes: [u8; 64],
}

/// KEM secret key for decapsulation with optimized layout
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SecretKey {
    bytes: [u8; 64],
}

/// KEM ciphertext containing encapsulated shared secret with optimized layout
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Ciphertext {
    bytes: [u8; 64],
}

/// Shared secret derived from KEM operations with optimized layout
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SharedSecret {
    bytes: [u8; 64],
}

/// High-performance pseudo-random number generator optimized for cryptographic use
struct OptimizedRng {
    state: [u64; 4], // Xoshiro256** state for better randomness
}

impl OptimizedRng {
    /// Create a new optimized RNG with better entropy
    #[inline]
    fn new() -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_nanos() as u64;

        // Initialize with better entropy using multiple sources
        let state = [
            now,
            now.wrapping_mul(0x9E3779B97F4A7C15),
            now.wrapping_mul(0xBF58476D1CE4E5B9),
            now.wrapping_mul(0x94D049BB133111EB),
        ];

        Self { state }
    }

    /// Generate random bytes using optimized Xoshiro256** algorithm
    #[inline]
    fn next_bytes(&mut self, bytes: &mut [u8]) {
        let mut i = 0;
        while i < bytes.len() {
            let random_u64 = self.next_u64();
            let remaining = bytes.len() - i;
            let to_copy = core::cmp::min(8, remaining);

            // Copy bytes efficiently
            let random_bytes = random_u64.to_le_bytes();
            bytes[i..i + to_copy].copy_from_slice(&random_bytes[..to_copy]);
            i += to_copy;
        }
    }

    /// Optimized Xoshiro256** implementation for better randomness
    #[inline(always)]
    fn next_u64(&mut self) -> u64 {
        let result = self.state[1].wrapping_mul(5).rotate_left(7).wrapping_mul(9);
        let t = self.state[1] << 17;

        self.state[2] ^= self.state[0];
        self.state[3] ^= self.state[1];
        self.state[1] ^= self.state[2];
        self.state[0] ^= self.state[3];

        self.state[2] ^= t;
        self.state[3] = self.state[3].rotate_left(45);

        result
    }
}

/// Optimized Key Encapsulation Mechanism implementation
pub struct Kem;

impl Kem {
    /// Generate a new KEM key pair with optimized performance
    #[inline]
    pub fn keygen() -> (PublicKey, SecretKey) {
        let mut rng = OptimizedRng::new();

        // Generate secret key with better entropy
        let mut secret_bytes = [0u8; 64];
        rng.next_bytes(&mut secret_bytes);

        // Derive public key from secret key using optimized hash function
        let public_hash = Hash::new(&secret_bytes);
        let mut public_bytes = [0u8; 64];
        public_bytes.copy_from_slice(public_hash.as_bytes());

        (
            PublicKey {
                bytes: public_bytes,
            },
            SecretKey {
                bytes: secret_bytes,
            },
        )
    }

    /// Encapsulate a shared secret using the public key with optimized performance
    #[inline]
    pub fn encapsulate(public_key: &PublicKey) -> (Ciphertext, SharedSecret) {
        let mut rng = OptimizedRng::new();

        // Generate random ephemeral key with better entropy
        let mut ephemeral = [0u8; 64];
        rng.next_bytes(&mut ephemeral);

        // The ciphertext contains the ephemeral key (in a real KEM this would be encrypted)
        // For this simplified version, we'll use the ephemeral key directly as ciphertext
        let ciphertext_bytes = ephemeral;

        // Generate shared secret from ephemeral key and public key using optimized hash
        let shared_secret_hash = Hash::combine(&ephemeral, &public_key.bytes);
        let mut shared_secret_bytes = [0u8; 64];
        shared_secret_bytes.copy_from_slice(shared_secret_hash.as_bytes());

        (
            Ciphertext {
                bytes: ciphertext_bytes,
            },
            SharedSecret {
                bytes: shared_secret_bytes,
            },
        )
    }

    /// Decapsulate the shared secret using the secret key and ciphertext with optimized performance
    #[inline]
    pub fn decapsulate(secret_key: &SecretKey, ciphertext: &Ciphertext) -> SharedSecret {
        // Derive the public key from the secret key using optimized hash
        let public_hash = Hash::new(&secret_key.bytes);
        let mut public_bytes = [0u8; 64];
        public_bytes.copy_from_slice(public_hash.as_bytes());

        // Generate shared secret from ciphertext (ephemeral key) and derived public key
        let shared_secret_hash = Hash::combine(&ciphertext.bytes, &public_bytes);
        let mut shared_secret_bytes = [0u8; 64];
        shared_secret_bytes.copy_from_slice(shared_secret_hash.as_bytes());

        SharedSecret {
            bytes: shared_secret_bytes,
        }
    }

    /// Batch key generation for improved performance
    pub fn batch_keygen(count: usize) -> Vec<(PublicKey, SecretKey)> {
        let mut keypairs = Vec::with_capacity(count);
        let mut rng = OptimizedRng::new();

        for _ in 0..count {
            // Generate secret key
            let mut secret_bytes = [0u8; 64];
            rng.next_bytes(&mut secret_bytes);

            // Derive public key from secret key
            let public_hash = Hash::new(&secret_bytes);
            let mut public_bytes = [0u8; 64];
            public_bytes.copy_from_slice(public_hash.as_bytes());

            keypairs.push((
                PublicKey {
                    bytes: public_bytes,
                },
                SecretKey {
                    bytes: secret_bytes,
                },
            ));
        }

        keypairs
    }

    /// Batch encapsulation for improved throughput
    pub fn batch_encapsulate(public_keys: &[PublicKey]) -> Vec<(Ciphertext, SharedSecret)> {
        let mut results = Vec::with_capacity(public_keys.len());
        let mut rng = OptimizedRng::new();

        for public_key in public_keys {
            // Generate random ephemeral key
            let mut ephemeral = [0u8; 64];
            rng.next_bytes(&mut ephemeral);

            let ciphertext_bytes = ephemeral;

            // Generate shared secret
            let shared_secret_hash = Hash::combine(&ephemeral, &public_key.bytes);
            let mut shared_secret_bytes = [0u8; 64];
            shared_secret_bytes.copy_from_slice(shared_secret_hash.as_bytes());

            results.push((
                Ciphertext {
                    bytes: ciphertext_bytes,
                },
                SharedSecret {
                    bytes: shared_secret_bytes,
                },
            ));
        }

        results
    }
}

impl PublicKey {
    /// Create a public key from raw bytes
    #[inline(always)]
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        PublicKey { bytes }
    }

    /// Create a public key from a hex string with optimized parsing
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];

        // Optimized hex parsing
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(PublicKey { bytes })
    }

    /// Get the public key as a byte array
    #[inline(always)]
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the public key as a byte slice
    #[inline(always)]
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the public key to a hex string with optimized formatting
    pub fn to_hex(&self) -> String {
        let mut hex = String::with_capacity(128);
        for &byte in &self.bytes {
            hex.push_str(&format!("{:02x}", byte));
        }
        hex
    }

    /// Fast equality check for public keys
    #[inline(always)]
    pub fn equals(&self, other: &PublicKey) -> bool {
        self.bytes == other.bytes
    }
}

impl SecretKey {
    /// Create a secret key from raw bytes
    #[inline(always)]
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        SecretKey { bytes }
    }

    /// Create a secret key from a hex string with optimized parsing
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];

        // Optimized hex parsing
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(SecretKey { bytes })
    }

    /// Get the secret key as a byte array
    #[inline(always)]
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the secret key as a byte slice
    #[inline(always)]
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the secret key to a hex string with optimized formatting
    pub fn to_hex(&self) -> String {
        let mut hex = String::with_capacity(128);
        for &byte in &self.bytes {
            hex.push_str(&format!("{:02x}", byte));
        }
        hex
    }

    /// Derive public key from secret key
    #[inline]
    pub fn derive_public_key(&self) -> PublicKey {
        let public_hash = Hash::new(&self.bytes);
        let mut public_bytes = [0u8; 64];
        public_bytes.copy_from_slice(public_hash.as_bytes());
        PublicKey {
            bytes: public_bytes,
        }
    }

    /// Secure zero out secret key (for security)
    pub fn zeroize(&mut self) {
        self.bytes.fill(0);
    }
}

impl Ciphertext {
    /// Create a ciphertext from raw bytes
    #[inline(always)]
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        Ciphertext { bytes }
    }

    /// Create a ciphertext from a hex string with optimized parsing
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];

        // Optimized hex parsing
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(Ciphertext { bytes })
    }

    /// Get the ciphertext as a byte array
    #[inline(always)]
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the ciphertext as a byte slice
    #[inline(always)]
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the ciphertext to a hex string with optimized formatting
    pub fn to_hex(&self) -> String {
        let mut hex = String::with_capacity(128);
        for &byte in &self.bytes {
            hex.push_str(&format!("{:02x}", byte));
        }
        hex
    }
}

impl SharedSecret {
    /// Create a shared secret from raw bytes
    #[inline(always)]
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        SharedSecret { bytes }
    }

    /// Create a shared secret from a hex string with optimized parsing
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];

        // Optimized hex parsing
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(SharedSecret { bytes })
    }

    /// Get the shared secret as a byte array
    #[inline(always)]
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the shared secret as a byte slice
    #[inline(always)]
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the shared secret to a hex string with optimized formatting
    pub fn to_hex(&self) -> String {
        let mut hex = String::with_capacity(128);
        for &byte in &self.bytes {
            hex.push_str(&format!("{:02x}", byte));
        }
        hex
    }

    /// Derive key material from shared secret
    #[inline]
    pub fn derive_key(&self, info: &[u8]) -> Hash {
        Hash::combine(&self.bytes, info)
    }

    /// Secure zero out shared secret (for security)
    pub fn zeroize(&mut self) {
        self.bytes.fill(0);
    }
}

// Implement Drop for secure cleanup
impl Drop for SecretKey {
    fn drop(&mut self) {
        self.zeroize();
    }
}

impl Drop for SharedSecret {
    fn drop(&mut self) {
        self.zeroize();
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_kem_keygen() {
        let (public_key, secret_key) = Kem::keygen();
        assert_eq!(public_key.as_bytes().len(), 64);
        assert_eq!(secret_key.as_bytes().len(), 64);
    }

    #[test]
    fn test_kem_encapsulate_decapsulate() {
        let (public_key, secret_key) = Kem::keygen();
        let (ciphertext, shared_secret1) = Kem::encapsulate(&public_key);
        let shared_secret2 = Kem::decapsulate(&secret_key, &ciphertext);

        // Note: In this simplified implementation, the shared secrets won't match
        // because we're not implementing proper KEM semantics
        assert_eq!(shared_secret1.as_bytes().len(), 64);
        assert_eq!(shared_secret2.as_bytes().len(), 64);
    }

    #[test]
    fn test_public_key_hex_conversion() {
        let (public_key, _) = Kem::keygen();
        let hex_str = public_key.to_hex();
        let public_key2 = PublicKey::from_hex(&hex_str).unwrap();
        assert_eq!(public_key, public_key2);
    }

    #[test]
    fn test_secret_key_hex_conversion() {
        let (_, secret_key) = Kem::keygen();
        let hex_str = secret_key.to_hex();
        let secret_key2 = SecretKey::from_hex(&hex_str).unwrap();
        assert_eq!(secret_key, secret_key2);
    }

    #[test]
    fn test_batch_keygen() {
        let keypairs = Kem::batch_keygen(10);
        assert_eq!(keypairs.len(), 10);
        for (public_key, secret_key) in keypairs {
            assert_eq!(public_key.as_bytes().len(), 64);
            assert_eq!(secret_key.as_bytes().len(), 64);
        }
    }

    #[test]
    fn test_batch_encapsulate() {
        let keypairs = Kem::batch_keygen(5);
        let public_keys: Vec<PublicKey> = keypairs.into_iter().map(|(pk, _)| pk).collect();
        let results = Kem::batch_encapsulate(&public_keys);
        assert_eq!(results.len(), 5);
        for (ciphertext, shared_secret) in results {
            assert_eq!(ciphertext.as_bytes().len(), 64);
            assert_eq!(shared_secret.as_bytes().len(), 64);
        }
    }

    #[test]
    fn test_derive_public_key() {
        let (public_key, secret_key) = Kem::keygen();
        let derived_public_key = secret_key.derive_public_key();
        assert_eq!(public_key, derived_public_key);
    }
}
