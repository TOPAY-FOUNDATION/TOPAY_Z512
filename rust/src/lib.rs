//! # TOPAY-Z512 Cryptographic Library
//!
//! TOPAY-Z512 is a quantum-safe 512-bit post-quantum cryptographic library that provides:
//! - Key Encapsulation Mechanism (KEM) based on Learning With Errors (LWE)
//! - 512-bit cryptographic hashing using SHA3-512
//! - Key pair generation and management
//! - Fragmented-block architecture support for parallel processing
//!
//! This library serves as the cryptographic foundation for TOPAY Foundation's
//! quantum-safe blockchain ecosystem.

#![cfg_attr(not(feature = "std"), no_std)]

#[cfg(not(feature = "std"))]
extern crate alloc;

#[cfg(not(feature = "std"))]
use alloc::{vec::Vec, string::String};

pub mod hash;
pub mod kem;
pub mod keypair;
pub mod error;

#[cfg(feature = "fragmentation")]
pub mod fragment;

// Re-export main types for convenience
pub use hash::Hash;
pub use kem::{Kem, PublicKey as KemPublicKey, SecretKey as KemSecretKey, Ciphertext, SharedSecret};
pub use keypair::{KeyPair, PrivateKey, PublicKey};
pub use error::{TopayzError, Result};

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