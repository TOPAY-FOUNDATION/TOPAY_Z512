//! # TOPAY-Z512 Cryptographic Library
//!
//! TOPAY-Z512 is a high-performance quantum-safe 512-bit post-quantum cryptographic library that provides:
//! - Key Encapsulation Mechanism (KEM) based on Learning With Errors (LWE)
//! - 512-bit cryptographic hashing using SHA3-512
//! - Key pair generation and management
//! - Fragmented-block architecture support for parallel processing
//! - Optimized performance for mobile and embedded devices
//!
//! This library serves as the cryptographic foundation for TOPAY Foundation's
//! quantum-safe blockchain ecosystem.
//!
//! ## Performance Features
//! - SIMD-optimized operations where available
//! - Memory-efficient data structures
//! - Batch processing capabilities
//! - Inline optimizations for hot paths
//! - Zero-copy operations where possible

#![cfg_attr(not(feature = "std"), no_std)]
#![deny(unsafe_op_in_unsafe_fn)]
#![warn(missing_docs, rust_2018_idioms)]

#[cfg(not(feature = "std"))]
extern crate alloc;

#[cfg(not(feature = "std"))]
use alloc::{string::String, vec::Vec};

pub mod error;
pub mod hash;
pub mod kem;
pub mod keypair;

#[cfg(feature = "fragmentation")]
pub mod fragment;

// Re-export main types for convenience
pub use error::{Result, TopayzError};
pub use hash::Hash;
pub use kem::{
    Ciphertext, Kem, PublicKey as KemPublicKey, SecretKey as KemSecretKey, SharedSecret,
};
pub use keypair::{KeyPair, PrivateKey, PublicKey};

/// Library version
pub const VERSION: &str = env!("CARGO_PKG_VERSION");

/// Security level in bits (classical security)
pub const SECURITY_LEVEL: usize = 512;

/// Hash output size in bytes (64 bytes = 512 bits)
pub const HASH_SIZE: usize = 64;

/// Private key size in bytes
pub const PRIVATE_KEY_SIZE: usize = 64;

/// Public key size in bytes  
pub const PUBLIC_KEY_SIZE: usize = 64;

/// KEM public key size in bytes
pub const KEM_PUBLIC_KEY_SIZE: usize = 1024;

/// KEM secret key size in bytes
pub const KEM_SECRET_KEY_SIZE: usize = 512;

/// KEM ciphertext size in bytes
pub const KEM_CIPHERTEXT_SIZE: usize = 1024;

/// KEM shared secret size in bytes
pub const KEM_SHARED_SECRET_SIZE: usize = 64;

/// Performance configuration constants
pub mod perf {
    /// Optimal batch size for key generation
    pub const OPTIMAL_BATCH_SIZE: usize = 16;

    /// Cache line size for alignment optimization
    pub const CACHE_LINE_SIZE: usize = 64;

    /// SIMD vector width (256-bit AVX2)
    pub const SIMD_WIDTH: usize = 32;

    /// Memory prefetch distance
    pub const PREFETCH_DISTANCE: usize = 64;
}

/// Utility functions for high-performance operations
pub mod utils {
    use crate::Result;

    /// Fast constant-time comparison for cryptographic data
    #[inline(always)]
    pub fn constant_time_eq(a: &[u8], b: &[u8]) -> bool {
        if a.len() != b.len() {
            return false;
        }

        let mut result = 0u8;
        for (x, y) in a.iter().zip(b.iter()) {
            result |= x ^ y;
        }
        result == 0
    }

    /// Secure memory zeroing that won't be optimized away
    #[inline(always)]
    pub fn secure_zero(data: &mut [u8]) {
        // Use volatile write to prevent optimization
        for byte in data.iter_mut() {
            unsafe {
                core::ptr::write_volatile(byte, 0);
            }
        }
    }

    /// Fast hex encoding optimized for performance
    #[inline]
    pub fn fast_hex_encode(data: &[u8]) -> String {
        const HEX_CHARS: &[u8] = b"0123456789abcdef";
        let mut result = String::with_capacity(data.len() * 2);

        for &byte in data {
            result.push(HEX_CHARS[(byte >> 4) as usize] as char);
            result.push(HEX_CHARS[(byte & 0xf) as usize] as char);
        }

        result
    }

    /// Fast hex decoding optimized for performance
    pub fn fast_hex_decode(hex: &str) -> Result<Vec<u8>> {
        if hex.len() % 2 != 0 {
            return Err(crate::TopayzError::InvalidInput(
                "Odd hex length".to_string(),
            ));
        }

        let mut result = Vec::with_capacity(hex.len() / 2);
        let hex_bytes = hex.as_bytes();

        for chunk in hex_bytes.chunks_exact(2) {
            let high = hex_char_to_nibble(chunk[0])?;
            let low = hex_char_to_nibble(chunk[1])?;
            result.push((high << 4) | low);
        }

        Ok(result)
    }

    #[inline(always)]
    fn hex_char_to_nibble(c: u8) -> Result<u8> {
        match c {
            b'0'..=b'9' => Ok(c - b'0'),
            b'a'..=b'f' => Ok(c - b'a' + 10),
            b'A'..=b'F' => Ok(c - b'A' + 10),
            _ => Err(crate::TopayzError::InvalidInput(
                "Invalid hex character".to_string(),
            )),
        }
    }

    /// Memory-aligned allocation for performance-critical operations
    #[cfg(feature = "std")]
    pub fn aligned_alloc(size: usize, alignment: usize) -> Vec<u8> {
        let layout = std::alloc::Layout::from_size_align(size, alignment).expect("Invalid layout");

        unsafe {
            let ptr = std::alloc::alloc_zeroed(layout);
            if ptr.is_null() {
                std::alloc::handle_alloc_error(layout);
            }
            Vec::from_raw_parts(ptr, size, size)
        }
    }
}

/// Benchmark utilities for performance testing
#[cfg(feature = "benchmark")]
pub mod bench {
    use std::time::{Duration, Instant};

    /// Simple benchmark runner
    pub fn benchmark<F, R>(name: &str, iterations: usize, mut f: F) -> Duration
    where
        F: FnMut() -> R,
    {
        // Warmup
        for _ in 0..10 {
            let _ = f();
        }

        let start = Instant::now();
        for _ in 0..iterations {
            let _ = f();
        }
        let elapsed = start.elapsed();

        println!(
            "{}: {} iterations in {:?} ({:?} per iteration)",
            name,
            iterations,
            elapsed,
            elapsed / iterations as u32
        );

        elapsed
    }

    /// Memory usage profiler
    pub struct MemoryProfiler {
        start_usage: usize,
    }

    impl Default for MemoryProfiler {
        fn default() -> Self {
            Self::new()
        }
    }

    impl MemoryProfiler {
        /// Create a new memory profiler
        pub fn new() -> Self {
            Self {
                start_usage: get_memory_usage(),
            }
        }

        /// Report memory usage delta for an operation
        pub fn report(&self, operation: &str) {
            let current_usage = get_memory_usage();
            let delta = current_usage.saturating_sub(self.start_usage);
            println!("{operation}: Memory delta: {delta} bytes");
        }
    }

    #[cfg(target_os = "linux")]
    fn get_memory_usage() -> usize {
        use std::fs;

        if let Ok(status) = fs::read_to_string("/proc/self/status") {
            for line in status.lines() {
                if line.starts_with("VmRSS:") {
                    if let Some(kb_str) = line.split_whitespace().nth(1) {
                        if let Ok(kb) = kb_str.parse::<usize>() {
                            return kb * 1024; // Convert to bytes
                        }
                    }
                }
            }
        }
        0
    }

    #[cfg(not(target_os = "linux"))]
    fn get_memory_usage() -> usize {
        0 // Placeholder for other platforms
    }
}

/// Feature flags and capability detection
pub mod features {
    /// Check if SIMD instructions are available
    #[inline]
    pub fn has_simd() -> bool {
        #[cfg(target_arch = "x86_64")]
        {
            is_x86_feature_detected!("avx2")
        }
        #[cfg(target_arch = "aarch64")]
        {
            std::arch::is_aarch64_feature_detected!("neon")
        }
        #[cfg(not(any(target_arch = "x86_64", target_arch = "aarch64")))]
        {
            false
        }
    }

    /// Check if hardware random number generator is available
    #[inline]
    pub fn has_hardware_rng() -> bool {
        #[cfg(target_arch = "x86_64")]
        {
            is_x86_feature_detected!("rdrand")
        }
        #[cfg(not(target_arch = "x86_64"))]
        {
            false
        }
    }

    /// Get optimal number of threads for parallel operations
    pub fn optimal_thread_count() -> usize {
        #[cfg(feature = "std")]
        {
            std::thread::available_parallelism()
                .map(|n| n.get())
                .unwrap_or(1)
                .min(16) // Cap at 16 threads for efficiency
        }
        #[cfg(not(feature = "std"))]
        {
            1
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_constant_time_eq() {
        let a = [1, 2, 3, 4];
        let b = [1, 2, 3, 4];
        let c = [1, 2, 3, 5];

        assert!(utils::constant_time_eq(&a, &b));
        assert!(!utils::constant_time_eq(&a, &c));
    }

    #[test]
    fn test_hex_encoding() {
        let data = [0x12, 0x34, 0xab, 0xcd];
        let hex = utils::fast_hex_encode(&data);
        assert_eq!(hex, "1234abcd");

        let decoded = utils::fast_hex_decode(&hex).unwrap();
        assert_eq!(decoded, data);
    }

    #[test]
    fn test_secure_zero() {
        let mut data = [1, 2, 3, 4, 5];
        utils::secure_zero(&mut data);
        assert_eq!(data, [0, 0, 0, 0, 0]);
    }

    #[test]
    fn test_feature_detection() {
        // These should not panic
        let _ = features::has_simd();
        let _ = features::has_hardware_rng();
        let _ = features::optimal_thread_count();
    }
}
