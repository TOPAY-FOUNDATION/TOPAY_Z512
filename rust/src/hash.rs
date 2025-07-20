//! TOPAY-Z512 Hash implementation
//!
//! This module provides a 512-bit cryptographic hash function implementation
//! that can be used as part of the TOPAY-Z512 cryptographic suite.

use sha3::{Digest, Sha3_512};
use zeroize::Zeroize;

/// The size of TOPAY-Z512 hash output in bytes (512 bits = 64 bytes)
pub const HASH_SIZE_BYTES: usize = 64;

/// Represents a TOPAY-Z512 hash value (512 bits)
#[derive(Clone, PartialEq, Eq)]
pub struct Hash([u8; HASH_SIZE_BYTES]);

impl Hash {
    /// Creates a new hash from the given data
    pub fn new(data: &[u8]) -> Self {
        let mut hasher = Sha3_512::new();
        hasher.update(data);
        let result = hasher.finalize();
        
        let mut hash_bytes = [0u8; HASH_SIZE_BYTES];
        hash_bytes.copy_from_slice(&result);
        
        Self(hash_bytes)
    }
    
    /// Creates a new hash by combining two input values
    pub fn combine(data1: &[u8], data2: &[u8]) -> Self {
        let mut hasher = Sha3_512::new();
        hasher.update(data1);
        hasher.update(data2);
        let result = hasher.finalize();
        
        let mut hash_bytes = [0u8; HASH_SIZE_BYTES];
        hash_bytes.copy_from_slice(&result);
        
        Self(hash_bytes)
    }
    
    /// Returns the hash value as a byte slice
    pub fn as_bytes(&self) -> &[u8; HASH_SIZE_BYTES] {
        &self.0
    }
    
    /// Returns the hash value as a mutable byte slice
    pub fn as_bytes_mut(&mut self) -> &mut [u8; HASH_SIZE_BYTES] {
        &mut self.0
    }
    
    /// Converts the hash to a hexadecimal string
    pub fn to_hex(&self) -> String {
        self.0.iter()
            .map(|b| format!("{:02x}", b))
            .collect()
    }
    
    /// Creates a hash from a hexadecimal string
    pub fn from_hex(hex: &str) -> Result<Self, String> {
        if hex.len() != HASH_SIZE_BYTES * 2 {
            return Err(format!("Invalid hex length: {}, expected {}", 
                              hex.len(), HASH_SIZE_BYTES * 2));
        }
        
        let mut bytes = [0u8; HASH_SIZE_BYTES];
        
        for i in 0..HASH_SIZE_BYTES {
            let byte_str = &hex[i*2..i*2+2];
            bytes[i] = u8::from_str_radix(byte_str, 16)
                .map_err(|e| format!("Invalid hex character: {}", e))?;
        }
        
        Ok(Self(bytes))
    }
}

impl Drop for Hash {
    fn drop(&mut self) {
        self.0.zeroize();
    }
}

impl std::fmt::Debug for Hash {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Hash({})...", self.to_hex().get(0..16).unwrap_or("<error>"))
    }
}

impl std::fmt::Display for Hash {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.to_hex())
    }
}

/// A convenience function to hash data
pub fn hash(data: &[u8]) -> [u8; HASH_SIZE_BYTES] {
    Hash::new(data).0
}

/// A convenience function to hash two pieces of data together
pub fn hash_combine(data1: &[u8], data2: &[u8]) -> [u8; HASH_SIZE_BYTES] {
    Hash::combine(data1, data2).0
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_hash_size() {
        let data = b"TOPAY-Z512 test data";
        let hash = Hash::new(data);
        assert_eq!(hash.as_bytes().len(), HASH_SIZE_BYTES);
    }
    
    #[test]
    fn test_hash_deterministic() {
        let data = b"TOPAY-Z512 test data";
        let hash1 = Hash::new(data);
        let hash2 = Hash::new(data);
        assert_eq!(hash1, hash2);
    }
    
    #[test]
    fn test_hash_different_inputs() {
        let data1 = b"TOPAY-Z512 test data 1";
        let data2 = b"TOPAY-Z512 test data 2";
        let hash1 = Hash::new(data1);
        let hash2 = Hash::new(data2);
        assert_ne!(hash1, hash2);
    }
    
    #[test]
    fn test_hash_combine() {
        let data1 = b"TOPAY-Z512";
        let data2 = b"test data";
        
        // Combined hash
        let combined = Hash::combine(data1, data2);
        
        // Concatenated hash
        let mut concatenated = Vec::new();
        concatenated.extend_from_slice(data1);
        concatenated.extend_from_slice(data2);
        let concat_hash = Hash::new(&concatenated);
        
        // These should be different
        assert_ne!(combined, concat_hash);
    }
    
    #[test]
    fn test_hex_conversion() {
        let data = b"TOPAY-Z512 hex conversion test";
        let hash = Hash::new(data);
        let hex = hash.to_hex();
        let hash2 = Hash::from_hex(&hex).unwrap();
        assert_eq!(hash, hash2);
    }
    
    #[test]
    fn test_invalid_hex() {
        let result = Hash::from_hex("invalid");
        assert!(result.is_err());
        
        let result = Hash::from_hex(&"00".repeat(HASH_SIZE_BYTES + 1));
        assert!(result.is_err());
        
        let result = Hash::from_hex(&"zz".repeat(HASH_SIZE_BYTES));
        assert!(result.is_err());
    }
}