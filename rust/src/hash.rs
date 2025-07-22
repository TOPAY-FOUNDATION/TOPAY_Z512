//! 512-bit cryptographic hash implementation for TOPAY-Z512
//! 
//! This module provides an optimized hash implementation for demonstration purposes.
//! In production, this would use SHA3-512 or another quantum-resistant hash function.

#[cfg(not(feature = "std"))]
use alloc::{vec::Vec, string::String};

use crate::{error::{TopayzError, Result}, HASH_SIZE};
use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash as StdHash, Hasher};

/// A 512-bit cryptographic hash with optimized operations
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Hash {
    bytes: [u8; HASH_SIZE],
}

impl Hash {
    /// Create a new hash from input data with optimized hashing
    #[inline]
    pub fn new(data: &[u8]) -> Self {
        let mut hasher = DefaultHasher::new();
        data.hash(&mut hasher);
        let hash_value = hasher.finish();
        
        // Optimized 512-bit expansion using SIMD-friendly operations
        let mut bytes = [0u8; HASH_SIZE];
        let hash_bytes = hash_value.to_le_bytes();
        
        // Unrolled loop for better performance
        unsafe {
            let bytes_ptr = bytes.as_mut_ptr();
            let hash_ptr = hash_bytes.as_ptr();
            
            // Fill 8 chunks of 8 bytes each with optimized mixing
            for i in 0..8 {
                let offset = i * 8;
                let mix_factor = (i as u8).wrapping_mul(17);
                
                for j in 0..8 {
                    *bytes_ptr.add(offset + j) = 
                        (*hash_ptr.add(j)).wrapping_add(mix_factor).wrapping_add(j as u8);
                }
            }
        }
        
        // Optimized length mixing
        let len_hash = (data.len() as u64).wrapping_mul(0x9e3779b97f4a7c15);
        let len_bytes = len_hash.to_le_bytes();
        
        // XOR length into first and last 8 bytes
        for i in 0..8 {
            bytes[i] ^= len_bytes[i];
            bytes[HASH_SIZE - 8 + i] ^= len_bytes[i];
        }
        
        Hash { bytes }
    }

    /// Create a hash from raw bytes (zero-cost)
    #[inline(always)]
    pub const fn from_bytes(bytes: [u8; HASH_SIZE]) -> Self {
        Hash { bytes }
    }

    /// Create a hash from a hex string with optimized parsing
    pub fn from_hex(hex: &str) -> Result<Self> {
        if hex.len() != HASH_SIZE * 2 {
            return Err(TopayzError::InvalidInput("Invalid hex length".to_string()));
        }

        let mut bytes = [0u8; HASH_SIZE];
        let hex_bytes = hex.as_bytes();
        
        // Optimized hex parsing
        for i in 0..HASH_SIZE {
            let idx = i * 2;
            let high = Self::hex_char_to_byte(hex_bytes[idx])?;
            let low = Self::hex_char_to_byte(hex_bytes[idx + 1])?;
            bytes[i] = (high << 4) | low;
        }
        
        Ok(Hash { bytes })
    }

    /// Fast hex character to byte conversion
    #[inline]
    fn hex_char_to_byte(c: u8) -> Result<u8> {
        match c {
            b'0'..=b'9' => Ok(c - b'0'),
            b'a'..=b'f' => Ok(c - b'a' + 10),
            b'A'..=b'F' => Ok(c - b'A' + 10),
            _ => Err(TopayzError::InvalidInput("Invalid hex character".to_string())),
        }
    }

    /// Get the hash as a byte array reference (zero-cost)
    #[inline(always)]
    pub const fn as_bytes(&self) -> &[u8; HASH_SIZE] {
        &self.bytes
    }

    /// Get the hash as a byte array copy
    #[inline(always)]
    pub const fn to_bytes(&self) -> [u8; HASH_SIZE] {
        self.bytes
    }

    /// Convert the hash to a hex string with optimized formatting
    pub fn to_hex(&self) -> String {
        const HEX_CHARS: &[u8] = b"0123456789abcdef";
        let mut result = String::with_capacity(HASH_SIZE * 2);
        
        unsafe {
            let result_bytes = result.as_mut_vec();
            result_bytes.reserve_exact(HASH_SIZE * 2);
            
            for &byte in &self.bytes {
                result_bytes.push(HEX_CHARS[(byte >> 4) as usize]);
                result_bytes.push(HEX_CHARS[(byte & 0xf) as usize]);
            }
            
            result.as_mut_vec().set_len(HASH_SIZE * 2);
        }
        
        result
    }

    /// Combine two pieces of data into a single hash (optimized, no allocation)
    #[inline]
    pub fn combine(data1: &[u8], data2: &[u8]) -> Self {
        let mut hasher = DefaultHasher::new();
        
        // Hash lengths first for domain separation
        data1.len().hash(&mut hasher);
        data2.len().hash(&mut hasher);
        
        // Hash data
        data1.hash(&mut hasher);
        data2.hash(&mut hasher);
        
        let hash_value = hasher.finish();
        let mut bytes = [0u8; HASH_SIZE];
        let hash_bytes = hash_value.to_le_bytes();
        
        // Optimized expansion with better mixing
        for i in 0..8 {
            let offset = i * 8;
            let mix = (i as u64).wrapping_mul(0x9e3779b97f4a7c15);
            let mix_bytes = mix.to_le_bytes();
            
            for j in 0..8 {
                bytes[offset + j] = hash_bytes[j] ^ mix_bytes[j];
            }
        }
        
        Hash { bytes }
    }

    /// Concatenate multiple hashes into a single hash (optimized)
    pub fn concat(hashes: &[&Hash]) -> Self {
        if hashes.is_empty() {
            return Hash::new(&[]);
        }
        
        let mut hasher = DefaultHasher::new();
        hashes.len().hash(&mut hasher);
        
        // Hash all input hashes efficiently
        for hash in hashes {
            hash.bytes.hash(&mut hasher);
        }
        
        let hash_value = hasher.finish();
        let mut bytes = [0u8; HASH_SIZE];
        let hash_bytes = hash_value.to_le_bytes();
        
        // Optimized expansion
        for i in 0..8 {
            let offset = i * 8;
            let mix = (i as u64).wrapping_mul(0xc6a4a7935bd1e995);
            let mix_bytes = mix.to_le_bytes();
            
            for j in 0..8 {
                bytes[offset + j] = hash_bytes[j] ^ mix_bytes[j];
            }
        }
        
        Hash { bytes }
    }

    /// Hash binary data (optimized convenience method)
    #[inline(always)]
    pub fn hash_binary(data: &[u8]) -> Self {
        Self::new(data)
    }

    /// Hash a string (UTF-8 encoded, optimized)
    #[inline(always)]
    pub fn hash_string(s: &str) -> Self {
        Self::new(s.as_bytes())
    }

    /// Fast equality check for first N bytes (useful for quick comparisons)
    #[inline]
    pub fn starts_with_bytes(&self, prefix: &[u8]) -> bool {
        if prefix.len() > HASH_SIZE {
            return false;
        }
        self.bytes[..prefix.len()] == *prefix
    }

    /// XOR two hashes together (useful for combining operations)
    #[inline]
    pub fn xor(&self, other: &Hash) -> Hash {
        let mut result = [0u8; HASH_SIZE];
        for i in 0..HASH_SIZE {
            result[i] = self.bytes[i] ^ other.bytes[i];
        }
        Hash { bytes: result }
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