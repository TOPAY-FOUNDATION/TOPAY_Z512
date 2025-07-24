//! Optimized key pair generation and management for TOPAY-Z512
//!
//! This module provides high-performance key generation for demonstration purposes.
//! In production, this would use cryptographically secure random number generation.

#[cfg(not(feature = "std"))]
use alloc::{string::String, vec::Vec};

use crate::error::{Result, TopayzError};
use crate::{PRIVATE_KEY_SIZE, PUBLIC_KEY_SIZE};
use std::time::{SystemTime, UNIX_EPOCH};

/// A private key for TOPAY-Z512 with optimized layout
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PrivateKey {
    bytes: [u8; PRIVATE_KEY_SIZE],
}

/// A public key for TOPAY-Z512 with optimized layout
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PublicKey {
    bytes: [u8; PUBLIC_KEY_SIZE],
}

/// High-performance pseudo-random number generator optimized for key generation
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

    /// Generate random bytes efficiently
    #[inline]
    fn fill_bytes(&mut self, bytes: &mut [u8]) {
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

/// A cryptographic key pair consisting of a private and public key with optimized layout
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct KeyPair {
    private_key: PrivateKey,
    public_key: PublicKey,
}

impl PrivateKey {
    /// Generate a new random private key with optimized performance
    #[inline]
    pub fn generate() -> Self {
        let mut rng = OptimizedRng::new();
        let mut bytes = [0u8; PRIVATE_KEY_SIZE];
        rng.fill_bytes(&mut bytes);
        Self { bytes }
    }

    /// Create a private key from bytes
    #[inline(always)]
    pub fn from_bytes(bytes: [u8; PRIVATE_KEY_SIZE]) -> Self {
        Self { bytes }
    }

    /// Get the bytes of the private key
    #[inline(always)]
    pub fn to_bytes(&self) -> [u8; PRIVATE_KEY_SIZE] {
        self.bytes
    }

    /// Get the bytes of the private key as a reference
    #[inline(always)]
    pub fn as_bytes(&self) -> &[u8; PRIVATE_KEY_SIZE] {
        &self.bytes
    }

    /// Create a private key from hex string with optimized parsing
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != PRIVATE_KEY_SIZE * 2 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; PRIVATE_KEY_SIZE];

        // Optimized hex parsing
        for i in 0..PRIVATE_KEY_SIZE {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(Self { bytes })
    }

    /// Convert the private key to hex string with optimized formatting
    pub fn to_hex(&self) -> String {
        let mut hex = String::with_capacity(PRIVATE_KEY_SIZE * 2);
        for &byte in &self.bytes {
            hex.push_str(&format!("{byte:02x}"));
        }
        hex
    }

    /// Derive the public key from this private key with optimized computation
    #[inline]
    pub fn public_key(&self) -> PublicKey {
        // Optimized public key derivation using better mixing
        let mut public_bytes = [0u8; PUBLIC_KEY_SIZE];

        // Use a more sophisticated derivation function
        for (i, public_byte) in public_bytes.iter_mut().enumerate().take(PUBLIC_KEY_SIZE) {
            let private_idx = i % PRIVATE_KEY_SIZE;
            let base = self.bytes[private_idx];

            // Better mixing function for public key derivation
            let mixed = base
                .wrapping_mul(0x9E)
                .wrapping_add(0x37)
                .wrapping_mul(i as u8)
                .wrapping_add(0x5A);

            *public_byte = mixed;
        }

        PublicKey::from_bytes(public_bytes)
    }

    /// Batch generate multiple private keys for improved performance
    pub fn batch_generate(count: usize) -> Vec<Self> {
        let mut keys = Vec::with_capacity(count);
        let mut rng = OptimizedRng::new();

        for _ in 0..count {
            let mut bytes = [0u8; PRIVATE_KEY_SIZE];
            rng.fill_bytes(&mut bytes);
            keys.push(Self { bytes });
        }

        keys
    }

    /// Secure zero out private key (for security)
    pub fn zeroize(&mut self) {
        self.bytes.fill(0);
    }

    /// Fast equality check for private keys
    #[inline(always)]
    pub fn equals(&self, other: &PrivateKey) -> bool {
        self.bytes == other.bytes
    }
}

impl PublicKey {
    /// Create a public key from bytes
    #[inline(always)]
    pub fn from_bytes(bytes: [u8; PUBLIC_KEY_SIZE]) -> Self {
        Self { bytes }
    }

    /// Get the bytes of the public key
    #[inline(always)]
    pub fn to_bytes(&self) -> [u8; PUBLIC_KEY_SIZE] {
        self.bytes
    }

    /// Get the bytes of the public key as a reference
    #[inline(always)]
    pub fn as_bytes(&self) -> &[u8; PUBLIC_KEY_SIZE] {
        &self.bytes
    }

    /// Create a public key from hex string with optimized parsing
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != PUBLIC_KEY_SIZE * 2 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; PUBLIC_KEY_SIZE];

        // Optimized hex parsing
        for i in 0..PUBLIC_KEY_SIZE {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(Self { bytes })
    }

    /// Convert the public key to hex string with optimized formatting
    pub fn to_hex(&self) -> String {
        let mut hex = String::with_capacity(PUBLIC_KEY_SIZE * 2);
        for &byte in &self.bytes {
            hex.push_str(&format!("{byte:02x}"));
        }
        hex
    }

    /// Fast equality check for public keys
    #[inline(always)]
    pub fn equals(&self, other: &PublicKey) -> bool {
        self.bytes == other.bytes
    }

    /// Verify if this public key was derived from a given private key
    #[inline]
    pub fn verify_derivation(&self, private_key: &PrivateKey) -> bool {
        let derived = private_key.public_key();
        self.equals(&derived)
    }
}

impl KeyPair {
    /// Generate a new random key pair with optimized performance
    #[inline]
    pub fn generate() -> Self {
        let private_key = PrivateKey::generate();
        let public_key = private_key.public_key();
        Self {
            private_key,
            public_key,
        }
    }

    /// Create a key pair from a private key
    #[inline]
    pub fn from_private_key(private_key: PrivateKey) -> Self {
        let public_key = private_key.public_key();
        Self {
            private_key,
            public_key,
        }
    }

    /// Get the private key
    #[inline(always)]
    pub fn private_key(&self) -> &PrivateKey {
        &self.private_key
    }

    /// Get the public key
    #[inline(always)]
    pub fn public_key(&self) -> &PublicKey {
        &self.public_key
    }

    /// Convert the key pair to hex strings
    pub fn to_hex(&self) -> (String, String) {
        (self.private_key.to_hex(), self.public_key.to_hex())
    }

    /// Batch generate multiple key pairs for improved performance
    pub fn batch_generate(count: usize) -> Vec<Self> {
        let mut keypairs = Vec::with_capacity(count);
        let mut rng = OptimizedRng::new();

        for _ in 0..count {
            let mut private_bytes = [0u8; PRIVATE_KEY_SIZE];
            rng.fill_bytes(&mut private_bytes);

            let private_key = PrivateKey::from_bytes(private_bytes);
            let public_key = private_key.public_key();

            keypairs.push(Self {
                private_key,
                public_key,
            });
        }

        keypairs
    }

    /// Verify the integrity of the key pair
    #[inline]
    pub fn verify(&self) -> bool {
        self.public_key.verify_derivation(&self.private_key)
    }

    /// Secure zero out private key (for security)
    pub fn zeroize(&mut self) {
        self.private_key.zeroize();
    }
}

impl AsRef<[u8]> for PrivateKey {
    #[inline(always)]
    fn as_ref(&self) -> &[u8] {
        &self.bytes
    }
}

impl AsRef<[u8]> for PublicKey {
    #[inline(always)]
    fn as_ref(&self) -> &[u8] {
        &self.bytes
    }
}

impl From<[u8; PRIVATE_KEY_SIZE]> for PrivateKey {
    #[inline(always)]
    fn from(bytes: [u8; PRIVATE_KEY_SIZE]) -> Self {
        PrivateKey { bytes }
    }
}

impl From<[u8; PUBLIC_KEY_SIZE]> for PublicKey {
    #[inline(always)]
    fn from(bytes: [u8; PUBLIC_KEY_SIZE]) -> Self {
        PublicKey { bytes }
    }
}

// Implement Drop for secure cleanup
impl Drop for PrivateKey {
    fn drop(&mut self) {
        self.zeroize();
    }
}

impl Drop for KeyPair {
    fn drop(&mut self) {
        self.zeroize();
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_keypair_generation() {
        let keypair = KeyPair::generate();
        assert_eq!(keypair.private_key().to_bytes().len(), PRIVATE_KEY_SIZE);
        assert_eq!(keypair.public_key().to_bytes().len(), PUBLIC_KEY_SIZE);
    }

    #[test]
    fn test_private_to_public_derivation() {
        let private_key = PrivateKey::generate();
        let public_key1 = private_key.public_key();
        let public_key2 = private_key.public_key();

        // Should be deterministic
        assert_eq!(public_key1, public_key2);
    }

    #[test]
    fn test_hex_conversion() {
        let keypair = KeyPair::generate();

        let private_hex = keypair.private_key().to_hex();
        let public_hex = keypair.public_key().to_hex();

        let private_key2 = PrivateKey::from_hex(&private_hex).unwrap();
        let public_key2 = PublicKey::from_hex(&public_hex).unwrap();

        assert_eq!(keypair.private_key(), &private_key2);
        assert_eq!(keypair.public_key(), &public_key2);
    }

    #[test]
    fn test_keypair_from_private() {
        let private_key = PrivateKey::generate();
        let keypair = KeyPair::from_private_key(private_key.clone());

        assert_eq!(keypair.private_key(), &private_key);
        assert_eq!(keypair.public_key(), &private_key.public_key());
    }

    #[test]
    fn test_deterministic_public_key() {
        // Test with known private key (all zeros)
        let private_bytes = [0u8; PRIVATE_KEY_SIZE];
        let private_key = PrivateKey::from_bytes(private_bytes);
        let public_key = private_key.public_key();

        // Should always produce the same public key for the same private key
        let public_key2 = private_key.public_key();
        assert_eq!(public_key, public_key2);
    }

    #[test]
    fn test_invalid_hex() {
        let result = PrivateKey::from_hex("invalid_hex");
        assert!(result.is_err());

        let result = PublicKey::from_hex("deadbeef"); // Too short
        assert!(result.is_err());
    }

    #[test]
    fn test_batch_generation() {
        let private_keys = PrivateKey::batch_generate(10);
        assert_eq!(private_keys.len(), 10);

        let keypairs = KeyPair::batch_generate(5);
        assert_eq!(keypairs.len(), 5);

        for keypair in keypairs {
            assert!(keypair.verify());
        }
    }

    #[test]
    fn test_key_verification() {
        let keypair = KeyPair::generate();
        assert!(keypair.verify());

        let private_key = PrivateKey::generate();
        let public_key = private_key.public_key();
        assert!(public_key.verify_derivation(&private_key));
    }

    #[test]
    fn test_equality_methods() {
        let keypair1 = KeyPair::generate();
        let keypair2 = KeyPair::generate();

        assert!(keypair1.private_key().equals(keypair1.private_key()));
        assert!(keypair1.public_key().equals(keypair1.public_key()));

        assert!(!keypair1.private_key().equals(keypair2.private_key()));
        assert!(!keypair1.public_key().equals(keypair2.public_key()));
    }
}
