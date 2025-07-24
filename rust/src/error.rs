//! Error types for TOPAY-Z512

#[cfg(not(feature = "std"))]
use alloc::{format, string::String};

use core::fmt;

/// Result type for TOPAY-Z512 operations
pub type Result<T> = core::result::Result<T, TopayzError>;

/// Error types for TOPAY-Z512 operations
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum TopayzError {
    /// Invalid input data
    InvalidInput(String),
    /// Invalid key format or size
    InvalidKey(String),
    /// Invalid hash format or size
    InvalidHash(String),
    /// Cryptographic operation failed
    CryptoError(String),
    /// Serialization/deserialization error
    SerializationError(String),
    /// Random number generation failed
    RandomError(String),
    /// Fragmentation error
    FragmentationError(String),
}

impl fmt::Display for TopayzError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            TopayzError::InvalidInput(msg) => write!(f, "Invalid input: {msg}"),
            TopayzError::InvalidKey(msg) => write!(f, "Invalid key: {msg}"),
            TopayzError::InvalidHash(msg) => write!(f, "Invalid hash: {msg}"),
            TopayzError::CryptoError(msg) => write!(f, "Cryptographic error: {msg}"),
            TopayzError::SerializationError(msg) => write!(f, "Serialization error: {msg}"),
            TopayzError::RandomError(msg) => write!(f, "Random error: {msg}"),
            TopayzError::FragmentationError(msg) => write!(f, "Fragmentation error: {msg}"),
        }
    }
}

#[cfg(feature = "std")]
impl std::error::Error for TopayzError {}

impl From<hex::FromHexError> for TopayzError {
    fn from(err: hex::FromHexError) -> Self {
        TopayzError::SerializationError(format!("Hex decode error: {err}"))
    }
}
