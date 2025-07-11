//! Error types for the topayz512 library.

use thiserror::Error;

/// Error type for the topayz512 library.
#[derive(Error, Debug)]
pub enum Error {
    /// Error during key generation.
    #[error("Key generation failed: {0}")]
    KeyGeneration(String),

    /// Error during encapsulation.
    #[error("Encapsulation failed: {0}")]
    Encapsulation(String),

    /// Error during decapsulation.
    #[error("Decapsulation failed: {0}")]
    Decapsulation(String),

    /// Error during fragmentation.
    #[cfg(feature = "fragmentation")]
    #[error("Fragmentation failed: {0}")]
    Fragmentation(String),

    /// Error during random number generation.
    #[error("Random number generation failed: {0}")]
    RandomGeneration(String),

    /// Invalid parameter error.
    #[error("Invalid parameter: {0}")]
    InvalidParameter(String),

    /// Invalid key format error.
    #[error("Invalid key format: {0}")]
    InvalidKeyFormat(String),

    /// Invalid ciphertext format error.
    #[error("Invalid ciphertext format: {0}")]
    InvalidCiphertextFormat(String),
}