//! TOPAY-Z512 Key Pair implementation
//!
//! This module provides a 512-bit cryptographic key pair implementation
//! that can be used as part of the TOPAY-Z512 cryptographic suite.

use rand::{CryptoRng, RngCore};
use sha3::{Digest, Sha3_512};
use zeroize::Zeroize;

use crate::hash::{Hash, HASH_SIZE_BYTES};

/// The size of TOPAY-Z512 private key in bytes (512 bits = 64 bytes)
pub const PRIVATE_KEY_SIZE_BYTES: usize = HASH_SIZE_BYTES;

/// The size of TOPAY-Z512 public key in bytes (512 bits = 64 bytes)
pub const PUBLIC_KEY_SIZE_BYTES: usize = HASH_SIZE_BYTES;

/// Represents a TOPAY-Z512 private key (512 bits)
#[derive(Clone, Zeroize)]
#[zeroize(drop)]
pub struct PrivateKey([u8; PRIVATE_KEY_SIZE_BYTES]);

/// Represents a TOPAY-Z512 public key (512 bits)
#[derive(Clone, PartialEq, Eq)]
pub struct PublicKey([u8; PUBLIC_KEY_SIZE_BYTES]);

/// Represents a TOPAY-Z512 key pair (private key and public key)
pub struct KeyPair {
    /// The private key
    pub private_key: PrivateKey,
    /// The public key
    pub public_key: PublicKey,
}

impl PrivateKey {
    /// Creates a new private key from random data
    pub fn generate<R: RngCore + CryptoRng>(rng: &mut R) -> Self {
        let mut key_bytes = [0u8; PRIVATE_KEY_SIZE_BYTES];
        rng.fill_bytes(&mut key_bytes);
        Self(key_bytes)
    }
    
    /// Creates a new private key from existing bytes
    pub fn from_bytes(bytes: [u8; PRIVATE_KEY_SIZE_BYTES]) -> Self {
        Self(bytes)
    }
    
    /// Returns the private key as a byte slice
    pub fn as_bytes(&self) -> &[u8; PRIVATE_KEY_SIZE_BYTES] {
        &self.0
    }
    
    /// Returns the private key as a mutable byte slice
    pub fn as_bytes_mut(&mut self) -> &mut [u8; PRIVATE_KEY_SIZE_BYTES] {
        &mut self.0
    }
    
    /// Converts the private key to a hexadecimal string
    pub fn to_hex(&self) -> String {
        let mut result = String::with_capacity(PRIVATE_KEY_SIZE_BYTES * 2);
        for byte in self.0.iter() {
            result.push_str(&format!("{:02x}", byte));
        }
        result
    }
    
    /// Creates a private key from a hexadecimal string
    pub fn from_hex(hex_str: &str) -> Result<Self, &'static str> {
        if hex_str.len() != PRIVATE_KEY_SIZE_BYTES * 2 {
            return Err("Invalid hex string length");
        }
        
        let mut bytes = [0u8; PRIVATE_KEY_SIZE_BYTES];
        for i in 0..PRIVATE_KEY_SIZE_BYTES {
            let byte_str = &hex_str[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(byte_str, 16)
                .map_err(|_| "Invalid hex string")?;
        }
        
        Ok(Self(bytes))
    }
}

impl PublicKey {
    /// Derives a public key from a private key
    pub fn from_private_key(private_key: &PrivateKey) -> Self {
        let mut hasher = Sha3_512::new();
        hasher.update(private_key.as_bytes());
        let result = hasher.finalize();
        
        let mut key_bytes = [0u8; PUBLIC_KEY_SIZE_BYTES];
        key_bytes.copy_from_slice(&result);
        
        Self(key_bytes)
    }
    
    /// Creates a new public key from existing bytes
    pub fn from_bytes(bytes: [u8; PUBLIC_KEY_SIZE_BYTES]) -> Self {
        Self(bytes)
    }
    
    /// Returns the public key as a byte slice
    pub fn as_bytes(&self) -> &[u8; PUBLIC_KEY_SIZE_BYTES] {
        &self.0
    }
    
    /// Returns the public key as a mutable byte slice
    pub fn as_bytes_mut(&mut self) -> &mut [u8; PUBLIC_KEY_SIZE_BYTES] {
        &mut self.0
    }
    
    /// Converts the public key to a hexadecimal string
    pub fn to_hex(&self) -> String {
        let mut result = String::with_capacity(PUBLIC_KEY_SIZE_BYTES * 2);
        for byte in self.0.iter() {
            result.push_str(&format!("{:02x}", byte));
        }
        result
    }
    
    /// Creates a public key from a hexadecimal string
    pub fn from_hex(hex_str: &str) -> Result<Self, &'static str> {
        if hex_str.len() != PUBLIC_KEY_SIZE_BYTES * 2 {
            return Err("Invalid hex string length");
        }
        
        let mut bytes = [0u8; PUBLIC_KEY_SIZE_BYTES];
        for i in 0..PUBLIC_KEY_SIZE_BYTES {
            let byte_str = &hex_str[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(byte_str, 16)
                .map_err(|_| "Invalid hex string")?;
        }
        
        Ok(Self(bytes))
    }
}

impl KeyPair {
    /// Generates a new key pair using the provided random number generator
    pub fn generate<R: RngCore + CryptoRng>(rng: &mut R) -> Self {
        let private_key = PrivateKey::generate(rng);
        let public_key = PublicKey::from_private_key(&private_key);
        
        Self {
            private_key,
            public_key,
        }
    }
    
    /// Creates a key pair from existing private and public keys
    pub fn from_keys(private_key: PrivateKey, public_key: PublicKey) -> Self {
        Self {
            private_key,
            public_key,
        }
    }
}

/// Convenience function to generate a key pair
pub fn generate_keypair<R: RngCore + CryptoRng>(rng: &mut R) -> KeyPair {
    KeyPair::generate(rng)
}

/// Convenience function to derive a public key from a private key
pub fn private_to_public(private_key: &PrivateKey) -> PublicKey {
    PublicKey::from_private_key(private_key)
}

#[cfg(test)]
mod tests {
    use super::*;
    use rand::rngs::OsRng;
    
    #[test]
    fn test_keypair_generation() {
        let mut rng = OsRng;
        let keypair = KeyPair::generate(&mut rng);
        
        // Verify key sizes
        assert_eq!(keypair.private_key.as_bytes().len(), PRIVATE_KEY_SIZE_BYTES);
        assert_eq!(keypair.public_key.as_bytes().len(), PUBLIC_KEY_SIZE_BYTES);
    }
    
    #[test]
    fn test_public_key_derivation() {
        let mut rng = OsRng;
        let private_key = PrivateKey::generate(&mut rng);
        let public_key = PublicKey::from_private_key(&private_key);
        
        // Verify that the same private key produces the same public key
        let public_key2 = PublicKey::from_private_key(&private_key);
        assert_eq!(public_key.as_bytes(), public_key2.as_bytes());
    }
    
    #[test]
    fn test_hex_conversion() {
        let mut rng = OsRng;
        let keypair = KeyPair::generate(&mut rng);
        
        // Test private key hex conversion
        let private_hex = keypair.private_key.to_hex();
        let private_key2 = PrivateKey::from_hex(&private_hex).unwrap();
        assert_eq!(keypair.private_key.as_bytes(), private_key2.as_bytes());
        
        // Test public key hex conversion
        let public_hex = keypair.public_key.to_hex();
        let public_key2 = PublicKey::from_hex(&public_hex).unwrap();
        assert_eq!(keypair.public_key.as_bytes(), public_key2.as_bytes());
    }
    
    #[test]
    fn test_invalid_hex() {
        // Test invalid hex length
        let result = PrivateKey::from_hex("invalid");
        assert!(result.is_err());
        
        // Test invalid hex characters
        let result = PrivateKey::from_hex(&"z".repeat(PRIVATE_KEY_SIZE_BYTES * 2));
        assert!(result.is_err());
    }
}