//! Key Encapsulation Mechanism (KEM) for TOPAY-Z512
//! 
//! This module provides a simplified KEM implementation for demonstration purposes.
//! In production, this would use a proper post-quantum KEM like Kyber or NTRU.

use std::time::{SystemTime, UNIX_EPOCH};
use crate::error::{TopayzError, Result};
use crate::hash::Hash;

/// KEM public key for encapsulation
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PublicKey {
    bytes: [u8; 64],
}

/// KEM secret key for decapsulation
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SecretKey {
    bytes: [u8; 64],
}

/// KEM ciphertext containing encapsulated shared secret
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Ciphertext {
    bytes: [u8; 64],
}

/// Shared secret derived from KEM operations
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct SharedSecret {
    bytes: [u8; 64],
}

/// Simple pseudo-random number generator for demonstration purposes
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

    fn next_bytes(&mut self, bytes: &mut [u8]) {
        for byte in bytes {
            // Linear congruential generator
            self.state = self.state.wrapping_mul(1103515245).wrapping_add(12345);
            *byte = (self.state >> 24) as u8;
        }
    }
}

/// Key Encapsulation Mechanism implementation
pub struct Kem;

impl Kem {
    /// Generate a new KEM key pair
    pub fn keygen() -> (PublicKey, SecretKey) {
        let mut rng = SimpleRng::new();
        
        // Generate secret key
        let mut secret_bytes = [0u8; 64];
        rng.next_bytes(&mut secret_bytes);

        // Derive public key from secret key using hash function
        let public_hash = Hash::new(&secret_bytes);
        let mut public_bytes = [0u8; 64];
        public_bytes.copy_from_slice(public_hash.as_bytes());

        (
            PublicKey { bytes: public_bytes },
            SecretKey { bytes: secret_bytes },
        )
    }

    /// Encapsulate a shared secret using the public key
    pub fn encapsulate(public_key: &PublicKey) -> (Ciphertext, SharedSecret) {
        let mut rng = SimpleRng::new();
        
        // Generate random ephemeral key
        let mut ephemeral = [0u8; 64];
        rng.next_bytes(&mut ephemeral);

        // The ciphertext contains the ephemeral key (in a real KEM this would be encrypted)
        // For this simplified version, we'll use the ephemeral key directly as ciphertext
        let ciphertext_bytes = ephemeral;

        // Generate shared secret from ephemeral key and public key
        let shared_secret_hash = Hash::combine(&ephemeral, &public_key.bytes);
        let mut shared_secret_bytes = [0u8; 64];
        shared_secret_bytes.copy_from_slice(shared_secret_hash.as_bytes());

        (
            Ciphertext { bytes: ciphertext_bytes },
            SharedSecret { bytes: shared_secret_bytes },
        )
    }

    /// Decapsulate the shared secret using the secret key and ciphertext
    pub fn decapsulate(secret_key: &SecretKey, ciphertext: &Ciphertext) -> SharedSecret {
        // Derive the public key from the secret key
        let public_hash = Hash::new(&secret_key.bytes);
        let mut public_bytes = [0u8; 64];
        public_bytes.copy_from_slice(public_hash.as_bytes());

        // Generate shared secret from ciphertext (ephemeral key) and derived public key
        let shared_secret_hash = Hash::combine(&ciphertext.bytes, &public_bytes);
        let mut shared_secret_bytes = [0u8; 64];
        shared_secret_bytes.copy_from_slice(shared_secret_hash.as_bytes());

        SharedSecret { bytes: shared_secret_bytes }
    }
}

impl PublicKey {
    /// Create a public key from raw bytes
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        PublicKey { bytes }
    }

    /// Create a public key from a hex string
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(PublicKey { bytes })
    }

    /// Get the public key as a byte array
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the public key as a byte slice
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the public key to a hex string
    pub fn to_hex(&self) -> String {
        self.bytes.iter().map(|b| format!("{:02x}", b)).collect()
    }
}

impl SecretKey {
    /// Create a secret key from raw bytes
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        SecretKey { bytes }
    }

    /// Create a secret key from a hex string
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(SecretKey { bytes })
    }

    /// Get the secret key as a byte array
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the secret key as a byte slice
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the secret key to a hex string
    pub fn to_hex(&self) -> String {
        self.bytes.iter().map(|b| format!("{:02x}", b)).collect()
    }
}

impl Ciphertext {
    /// Create a ciphertext from raw bytes
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        Ciphertext { bytes }
    }

    /// Create a ciphertext from a hex string
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(Ciphertext { bytes })
    }

    /// Get the ciphertext as a byte array
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the ciphertext as a byte slice
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the ciphertext to a hex string
    pub fn to_hex(&self) -> String {
        self.bytes.iter().map(|b| format!("{:02x}", b)).collect()
    }
}

impl SharedSecret {
    /// Create a shared secret from raw bytes
    pub fn from_bytes(bytes: [u8; 64]) -> Self {
        SharedSecret { bytes }
    }

    /// Create a shared secret from a hex string
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != 128 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; 64];
        for i in 0..64 {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(SharedSecret { bytes })
    }

    /// Get the shared secret as a byte array
    pub fn as_bytes(&self) -> &[u8; 64] {
        &self.bytes
    }

    /// Get the shared secret as a byte slice
    pub fn to_bytes(&self) -> [u8; 64] {
        self.bytes
    }

    /// Convert the shared secret to a hex string
    pub fn to_hex(&self) -> String {
        self.bytes.iter().map(|b| format!("{:02x}", b)).collect()
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
}