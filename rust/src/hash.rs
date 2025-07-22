//! 512-bit cryptographic hash implementation for TOPAY-Z512
//! 
//! This module provides a simplified hash implementation for demonstration purposes.
//! In production, this would use SHA3-512 or another quantum-resistant hash function.

#[cfg(not(feature = "std"))]
use alloc::{vec::Vec, string::String};

use crate::{error::{TopayzError, Result}, HASH_SIZE};
use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash as StdHash, Hasher};

/// A 512-bit cryptographic hash
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Hash {
    bytes: [u8; HASH_SIZE],
}

impl Hash {
    /// Create a new hash from input data
    pub fn new(data: &[u8]) -> Self {
        let mut hasher = DefaultHasher::new();
        data.hash(&mut hasher);
        let hash_value = hasher.finish();
        
        // Expand the 64-bit hash to 512 bits using a simple method
        let mut bytes = [0u8; HASH_SIZE];
        let hash_bytes = hash_value.to_le_bytes();
        
        // Fill the 512-bit array by repeating and XORing the 64-bit hash
        for i in 0..8 {
            let offset = i * 8;
            for j in 0..8 {
                bytes[offset + j] = hash_bytes[j] ^ ((i as u8) * 17 + j as u8);
            }
        }
        
        // Add some additional mixing based on data length
        let len_bytes = (data.len() as u64).to_le_bytes();
        for i in 0..8 {
            bytes[i] ^= len_bytes[i];
            bytes[56 + i] ^= len_bytes[i];
        }
        
        Hash { bytes }
    }

    /// Create a hash from raw bytes
    pub fn from_bytes(bytes: [u8; HASH_SIZE]) -> Self {
        Hash { bytes }
    }

    /// Create a hash from a hex string
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != HASH_SIZE * 2 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; HASH_SIZE];
        for i in 0..HASH_SIZE {
            let hex_byte = &hex[i * 2..i * 2 + 2];
            bytes[i] = u8::from_str_radix(hex_byte, 16)
                .map_err(|_| TopayzError::InvalidInput("Invalid hex character".to_string()))?;
        }
        Ok(Hash { bytes })
    }

    /// Get the hash as a byte array
    pub fn as_bytes(&self) -> &[u8; HASH_SIZE] {
        &self.bytes
    }

    /// Get the hash as a byte slice
    pub fn to_bytes(&self) -> [u8; HASH_SIZE] {
        self.bytes
    }

    /// Convert the hash to a hex string
    pub fn to_hex(&self) -> String {
        self.bytes.iter().map(|b| format!("{:02x}", b)).collect()
    }

    /// Combine two pieces of data into a single hash
    pub fn combine(data1: &[u8], data2: &[u8]) -> Self {
        // Simple combination without Vec - just hash the concatenated length and data
        let mut hasher = DefaultHasher::new();
        data1.len().hash(&mut hasher);
        data1.hash(&mut hasher);
        data2.len().hash(&mut hasher);
        data2.hash(&mut hasher);
        
        let hash_value = hasher.finish();
        let mut bytes = [0u8; HASH_SIZE];
        let hash_bytes = hash_value.to_le_bytes();
        
        // Fill the array with the combined hash
        for i in 0..8 {
            let offset = i * 8;
            for j in 0..8 {
                bytes[offset + j] = hash_bytes[j] ^ ((i as u8) * 23 + j as u8);
            }
        }
        
        Hash { bytes }
    }

    /// Concatenate multiple hashes into a single hash
    pub fn concat(hashes: &[&Hash]) -> Self {
        let mut hasher = DefaultHasher::new();
        hashes.len().hash(&mut hasher);
        for hash in hashes {
            hash.bytes.hash(&mut hasher);
        }
        
        let hash_value = hasher.finish();
        let mut bytes = [0u8; HASH_SIZE];
        let hash_bytes = hash_value.to_le_bytes();
        
        // Fill the array with the concatenated hash
        for i in 0..8 {
            let offset = i * 8;
            for j in 0..8 {
                bytes[offset + j] = hash_bytes[j] ^ ((i as u8) * 31 + j as u8);
            }
        }
        
        Hash { bytes }
    }

    /// Hash binary data (convenience method)
    pub fn hash_binary(data: &[u8]) -> Self {
        Self::new(data)
    }

    /// Hash a string (UTF-8 encoded)
    pub fn hash_string(s: &str) -> Self {
        Self::new(s.as_bytes())
    }
}

impl AsRef<[u8]> for Hash {
    fn as_ref(&self) -> &[u8] {
        &self.bytes
    }
}

impl From<[u8; HASH_SIZE]> for Hash {
    fn from(bytes: [u8; HASH_SIZE]) -> Self {
        Hash { bytes }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_hash_creation() {
        let data = b"Hello, TOPAY!";
        let hash = Hash::new(data);
        assert_eq!(hash.as_bytes().len(), HASH_SIZE);
    }

    #[test]
    fn test_hash_hex_conversion() {
        let data = b"test data";
        let hash = Hash::new(data);
        let hex_str = hash.to_hex();
        let hash2 = Hash::from_hex(&hex_str).unwrap();
        assert_eq!(hash, hash2);
    }

    #[test]
    fn test_hash_combine() {
        let data1 = b"first";
        let data2 = b"second";
        let combined = Hash::combine(data1, data2);
        
        // Should be different from individual hashes
        let hash1 = Hash::new(data1);
        let hash2 = Hash::new(data2);
        assert_ne!(combined, hash1);
        assert_ne!(combined, hash2);
    }

    #[test]
    fn test_hash_concat() {
        let hash1 = Hash::new(b"first");
        let hash2 = Hash::new(b"second");
        let hash3 = Hash::new(b"third");
        
        let concatenated = Hash::concat(&[&hash1, &hash2, &hash3]);
        assert_eq!(concatenated.as_bytes().len(), HASH_SIZE);
    }

    #[test]
    fn test_invalid_hex() {
        let result = Hash::from_hex("invalid_hex");
        assert!(result.is_err());
        
        let result = Hash::from_hex("deadbeef"); // Too short
        assert!(result.is_err());
    }
}