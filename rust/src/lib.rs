//! # topayz512
//!
//! A 512-bit post-quantum cryptography library based on Learning With Errors (LWE).
//! This library provides a Key Encapsulation Mechanism (KEM) with ≥512-bit classical security
//! and ~256-bit quantum resistance.

#![cfg_attr(not(feature = "std"), no_std)]

mod error;
mod kem;
mod lwe;
mod params;
mod utils;

#[cfg(feature = "fragmentation")]
mod fragment;

pub use error::Error;
pub use kem::{Ciphertext, PublicKey, SecretKey};

/// Key generation function that produces a key pair (public key, secret key)
/// for the KEM.
///
/// # Returns
///
/// A tuple containing the public key and secret key.
pub fn keygen() -> Result<(PublicKey, SecretKey), Error> {
    kem::keygen()
}

/// Encapsulates a shared secret using the recipient's public key.
///
/// # Arguments
///
/// * `public_key` - The recipient's public key
///
/// # Returns
///
/// A tuple containing the ciphertext and the shared secret.
pub fn encapsulate(public_key: &PublicKey) -> Result<(Ciphertext, Vec<u8>), Error> {
    kem::encapsulate(public_key)
}

/// Decapsulates a shared secret using the recipient's secret key and the sender's ciphertext.
///
/// # Arguments
///
/// * `secret_key` - The recipient's secret key
/// * `ciphertext` - The ciphertext from the sender
///
/// # Returns
///
/// The shared secret.
pub fn decapsulate(secret_key: &SecretKey, ciphertext: &Ciphertext) -> Result<Vec<u8>, Error> {
    kem::decapsulate(secret_key, ciphertext)
}

#[cfg(feature = "fragmentation")]
pub use fragment::{FragmentedCiphertext, FragmentedPublicKey, FragmentedSecretKey};

#[cfg(feature = "fragmentation")]
pub mod fragmented {
    use super::*;

    /// Key generation function that produces a fragmented key pair for the KEM.
    ///
    /// # Arguments
    ///
    /// * `num_fragments` - The number of fragments to split the keys into
    ///
    /// # Returns
    ///
    /// A tuple containing the fragmented public key and fragmented secret key.
    pub fn keygen(num_fragments: usize) -> Result<(FragmentedPublicKey, FragmentedSecretKey), Error> {
        fragment::keygen(num_fragments)
    }

    /// Encapsulates a shared secret using a fragmented public key.
    ///
    /// # Arguments
    ///
    /// * `public_key` - The fragmented public key
    ///
    /// # Returns
    ///
    /// A tuple containing the fragmented ciphertext and the shared secret.
    pub fn encapsulate(
        public_key: &FragmentedPublicKey,
    ) -> Result<(FragmentedCiphertext, Vec<u8>), Error> {
        fragment::encapsulate(public_key)
    }

    /// Decapsulates a shared secret using a fragmented secret key and ciphertext.
    ///
    /// # Arguments
    ///
    /// * `secret_key` - The fragmented secret key
    /// * `ciphertext` - The fragmented ciphertext
    ///
    /// # Returns
    ///
    /// The shared secret.
    pub fn decapsulate(
        secret_key: &FragmentedSecretKey,
        ciphertext: &FragmentedCiphertext,
    ) -> Result<Vec<u8>, Error> {
        fragment::decapsulate(secret_key, ciphertext)
    }
}

#[cfg(test)]
mod tests {
    use crate::{decapsulate, encapsulate, keygen};
    
    #[test]
    fn test_kem_roundtrip() {
        // Generate a key pair
        let (pk, sk) = keygen().unwrap();
        
        // Encapsulate a shared secret
        let (ct, secret1) = encapsulate(&pk).unwrap();
        
        // Decapsulate the shared secret
        let secret2 = decapsulate(&sk, &ct).unwrap();
        
        // Verify that the shared secrets match
        assert_eq!(secret1, secret2);
    }
    
    #[cfg(feature = "fragmentation")]
    #[test]
    fn test_fragmented_kem_roundtrip() {
        use crate::fragmented;
        
        // Number of fragments
        let num_fragments = 4;
        
        // Generate a fragmented key pair
        let (pk, sk) = fragmented::keygen(num_fragments).unwrap();
        
        // Encapsulate a shared secret
        let (ct, secret1) = fragmented::encapsulate(&pk).unwrap();
        
        // Decapsulate the shared secret
        let secret2 = fragmented::decapsulate(&sk, &ct).unwrap();
        
        // Verify that the shared secrets match
        assert_eq!(secret1, secret2);
    }
}