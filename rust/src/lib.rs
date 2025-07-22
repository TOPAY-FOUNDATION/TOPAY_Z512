//! TOPAY-Z512 Cryptographic Library
//!
//! A 512-bit post-quantum cryptography library based on LWE.
//! This library provides implementations for key encapsulation mechanisms (KEM)
//! and cryptographic hashing functions.

#![forbid(unsafe_code)]
#![warn(missing_docs, rust_2018_idioms)]

pub mod hash;
pub mod keypair;
pub mod kem;

/// Re-export commonly used types and functions
pub use hash::{Hash, hash, hash_combine, HASH_SIZE_BYTES};
pub use keypair::{KeyPair, PrivateKey, PublicKey, generate_keypair, private_to_public,
                  PRIVATE_KEY_SIZE_BYTES, PUBLIC_KEY_SIZE_BYTES};
pub use kem::{Ciphertext, SharedSecret, keygen, encapsulate, decapsulate,
              N, Q, SIGMA, SECRET_LENGTH, CIPHERTEXT_SIZE_BYTES};

/// Library version information
pub const VERSION: &str = env!("CARGO_PKG_VERSION");

/// Library name
pub const NAME: &str = env!("CARGO_PKG_NAME");