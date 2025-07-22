//! Key pair generation and management for TOPAY-Z512
//! 
//! This module provides simplified key generation for demonstration purposes.
//! In production, this would use cryptographically secure random number generation.

#[cfg(not(feature = "std"))]
use alloc::{vec::Vec, string::String};

use std::time::{SystemTime, UNIX_EPOCH};
use crate::error::{TopayzError, Result};
use crate::{PRIVATE_KEY_SIZE, PUBLIC_KEY_SIZE};

/// A private key for TOPAY-Z512
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PrivateKey {
    bytes: [u8; PRIVATE_KEY_SIZE],
}

/// A public key for TOPAY-Z512
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PublicKey {
    bytes: [u8; PUBLIC_KEY_SIZE],
}

/// A simple pseudo-random number generator for demonstration purposes
struct SimpleRng {
    state: u64,
}

impl SimpleRng {
    fn new() -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_nanos() as u64;
        Self { state: now }
    }

    fn next_u8(&mut self) -> u8 {
        // Linear congruential generator
        self.state = self.state.wrapping_mul(1103515245).wrapping_add(12345);
        (self.state >> 24) as u8
    }
}

/// A cryptographic key pair consisting of a private and public key
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct KeyPair {
    private_key: PrivateKey,
    public_key: PublicKey,
}

impl PrivateKey {
    /// Generate a new random private key
    pub fn generate() -> Self {
        let mut rng = SimpleRng::new();
        let mut bytes = [0u8; PRIVATE_KEY_SIZE];
        for byte in &mut bytes {
            *byte = rng.next_u8();
        }
        Self { bytes }
    }

    /// Create a private key from bytes
    pub fn from_bytes(bytes: [u8; PRIVATE_KEY_SIZE]) -> Self {
        Self { bytes }
    }

    /// Get the bytes of the private key
    pub fn to_bytes(&self) -> [u8; PRIVATE_KEY_SIZE] {
        self.bytes
    }

    /// Create a private key from hex string
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != PRIVATE_KEY_SIZE * 2 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }
        
        let mut bytes = [0u8; PRIVATE_KEY_SIZE];
        for i in 0..PRIVATE_KEY_SIZE {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(Self { bytes })
    }

    /// Convert the private key to hex string
    pub fn to_hex(&self) -> String {
        self.bytes.iter().map(|b| format!("{:02x}", b)).collect()
    }

    /// Derive the public key from this private key
    pub fn public_key(&self) -> PublicKey {
        // Simplified public key derivation - in practice this would use proper cryptographic operations
        let mut public_bytes = [0u8; PUBLIC_KEY_SIZE];
        for i in 0..PUBLIC_KEY_SIZE {
            public_bytes[i] = self.bytes[i % PRIVATE_KEY_SIZE].wrapping_mul(3).wrapping_add(7);
        }
        PublicKey::from_bytes(public_bytes)
    }
}

impl PublicKey {
    /// Create a public key from bytes
    pub fn from_bytes(bytes: [u8; PUBLIC_KEY_SIZE]) -> Self {
        Self { bytes }
    }

    /// Get the bytes of the public key
    pub fn to_bytes(&self) -> [u8; PUBLIC_KEY_SIZE] {
        self.bytes
    }

    /// Create a public key from hex string
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != PUBLIC_KEY_SIZE * 2 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }
        
        let mut bytes = [0u8; PUBLIC_KEY_SIZE];
        for i in 0..PUBLIC_KEY_SIZE {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(Self { bytes })
    }

    /// Convert the public key to hex string
    pub fn to_hex(&self) -> String {
        self.bytes.iter().map(|b| format!("{:02x}", b)).collect()
    }
}

impl KeyPair {
    /// Generate a new random key pair
    pub fn generate() -> Self {
        let private_key = PrivateKey::generate();
        let public_key = private_key.public_key();
        Self { private_key, public_key }
    }

    /// Create a key pair from a private key
    pub fn from_private_key(private_key: PrivateKey) -> Self {
        let public_key = private_key.public_key();
        Self {
            private_key,
            public_key,
        }
    }

    /// Get the private key
    pub fn private_key(&self) -> &PrivateKey {
        &self.private_key
    }

    /// Get the public key
    pub fn public_key(&self) -> &PublicKey {
        &self.public_key
    }

    /// Convert the key pair to hex strings
    pub fn to_hex(&self) -> (String, String) {
        (self.private_key.to_hex(), self.public_key.to_hex())
    }
}

impl AsRef<[u8]> for PrivateKey {
    fn as_ref(&self) -> &[u8] {
        &self.bytes
    }
}

impl AsRef<[u8]> for PublicKey {
    fn as_ref(&self) -> &[u8] {
        &self.bytes
    }
}

impl From<[u8; PRIVATE_KEY_SIZE]> for PrivateKey {
    fn from(bytes: [u8; PRIVATE_KEY_SIZE]) -> Self {
        PrivateKey { bytes }
    }
}

impl From<[u8; PUBLIC_KEY_SIZE]> for PublicKey {
    fn from(bytes: [u8; PUBLIC_KEY_SIZE]) -> Self {
        PublicKey { bytes }
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
}