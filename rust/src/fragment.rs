//! Fragmentation support for the topayz512 library.
//!
//! This module provides functionality to split large cryptographic operations
//! into smaller workloads to support high throughput on desktops and
//! acceptable latency on smartphones and embedded devices.

use crate::error::Error;
use crate::kem::{Ciphertext, PublicKey, SecretKey};
use crate::params::{MAX_FRAGMENTS, MIN_FRAGMENTS, N};
use crate::utils::hash;
use rand::{CryptoRng, RngCore};
use zeroize::Zeroize;

/// A fragmented public key for the KEM.
#[derive(Clone, Debug)]
pub struct FragmentedPublicKey {
    /// The fragments of the public key.
    fragments: Vec<PublicKey>,
}

/// A fragmented secret key for the KEM.
#[derive(Clone, Debug, Zeroize)]
#[zeroize(drop)]
pub struct FragmentedSecretKey {
    /// The fragments of the secret key.
    fragments: Vec<SecretKey>,
}

/// A fragmented ciphertext for the KEM.
#[derive(Clone, Debug)]
pub struct FragmentedCiphertext {
    /// The fragments of the ciphertext.
    fragments: Vec<Ciphertext>,
}

impl FragmentedPublicKey {
    /// Creates a new fragmented public key from fragments.
    ///
    /// # Arguments
    ///
    /// * `fragments` - The public key fragments
    ///
    /// # Returns
    ///
    /// A new fragmented public key.
    pub fn new(fragments: Vec<PublicKey>) -> Result<Self, Error> {
        if fragments.len() < MIN_FRAGMENTS || fragments.len() > MAX_FRAGMENTS {
            return Err(Error::Fragmentation(
                format!(
                    "Number of fragments must be between {} and {}",
                    MIN_FRAGMENTS, MAX_FRAGMENTS
                )
            ));
        }
        Ok(Self { fragments })
    }

    /// Returns the number of fragments.
    pub fn num_fragments(&self) -> usize {
        self.fragments.len()
    }

    /// Returns a reference to the fragments.
    pub fn fragments(&self) -> &[PublicKey] {
        &self.fragments
    }
}

impl FragmentedSecretKey {
    /// Creates a new fragmented secret key from fragments.
    ///
    /// # Arguments
    ///
    /// * `fragments` - The secret key fragments
    ///
    /// # Returns
    ///
    /// A new fragmented secret key.
    pub fn new(fragments: Vec<SecretKey>) -> Result<Self, Error> {
        if fragments.len() < MIN_FRAGMENTS || fragments.len() > MAX_FRAGMENTS {
            return Err(Error::Fragmentation(
                format!(
                    "Number of fragments must be between {} and {}",
                    MIN_FRAGMENTS, MAX_FRAGMENTS
                )
            ));
        }
        Ok(Self { fragments })
    }

    /// Returns the number of fragments.
    pub fn num_fragments(&self) -> usize {
        self.fragments.len()
    }

    /// Returns a reference to the fragments.
    pub fn fragments(&self) -> &[SecretKey] {
        &self.fragments
    }
}

impl FragmentedCiphertext {
    /// Creates a new fragmented ciphertext from fragments.
    ///
    /// # Arguments
    ///
    /// * `fragments` - The ciphertext fragments
    ///
    /// # Returns
    ///
    /// A new fragmented ciphertext.
    pub fn new(fragments: Vec<Ciphertext>) -> Result<Self, Error> {
        if fragments.len() < MIN_FRAGMENTS || fragments.len() > MAX_FRAGMENTS {
            return Err(Error::Fragmentation(
                format!(
                    "Number of fragments must be between {} and {}",
                    MIN_FRAGMENTS, MAX_FRAGMENTS
                )
            ));
        }
        Ok(Self { fragments })
    }

    /// Returns the number of fragments.
    pub fn num_fragments(&self) -> usize {
        self.fragments.len()
    }

    /// Returns a reference to the fragments.
    pub fn fragments(&self) -> &[Ciphertext] {
        &self.fragments
    }
}

/// Generates a fragmented key pair for the KEM.
///
/// # Arguments
///
/// * `num_fragments` - The number of fragments to split the keys into
///
/// # Returns
///
/// A tuple containing the fragmented public key and fragmented secret key.
pub fn keygen(num_fragments: usize) -> Result<(FragmentedPublicKey, FragmentedSecretKey), Error> {
    // Use the system's secure random number generator
    let mut rng = rand::thread_rng();
    keygen_with_rng(&mut rng, num_fragments)
}

/// Generates a fragmented key pair for the KEM using the provided random number generator.
///
/// # Arguments
///
/// * `rng` - Random number generator
/// * `num_fragments` - The number of fragments to split the keys into
///
/// # Returns
///
/// A tuple containing the fragmented public key and fragmented secret key.
pub fn keygen_with_rng<R: RngCore + CryptoRng>(
    rng: &mut R,
    num_fragments: usize,
) -> Result<(FragmentedPublicKey, FragmentedSecretKey), Error> {
    if num_fragments < MIN_FRAGMENTS || num_fragments > MAX_FRAGMENTS {
        return Err(Error::Fragmentation(
            format!(
                "Number of fragments must be between {} and {}",
                MIN_FRAGMENTS, MAX_FRAGMENTS
            )
        ));
    }

    // Calculate the size of each fragment
    let fragment_size = (N + num_fragments - 1) / num_fragments;

    // Generate key pairs for each fragment
    let mut pk_fragments = Vec::with_capacity(num_fragments);
    let mut sk_fragments = Vec::with_capacity(num_fragments);

    for _ in 0..num_fragments {
        let (pk, sk) = crate::kem::keygen_with_rng(rng)?;
        pk_fragments.push(pk);
        sk_fragments.push(sk);
    }

    Ok((FragmentedPublicKey::new(pk_fragments)?, FragmentedSecretKey::new(sk_fragments)?))
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
    // Use the system's secure random number generator
    let mut rng = rand::thread_rng();
    encapsulate_with_rng(public_key, &mut rng)
}

/// Encapsulates a shared secret using a fragmented public key and the provided random number generator.
///
/// # Arguments
///
/// * `public_key` - The fragmented public key
/// * `rng` - Random number generator
///
/// # Returns
///
/// A tuple containing the fragmented ciphertext and the shared secret.
pub fn encapsulate_with_rng<R: RngCore + CryptoRng>(
    public_key: &FragmentedPublicKey,
    rng: &mut R,
) -> Result<(FragmentedCiphertext, Vec<u8>), Error> {
    let num_fragments = public_key.num_fragments();
    let mut ct_fragments = Vec::with_capacity(num_fragments);
    let mut shared_secrets = Vec::with_capacity(num_fragments);

    // Encapsulate for each fragment
    for pk in public_key.fragments() {
        let (ct, ss) = crate::kem::encapsulate_with_rng(pk, rng)?;
        ct_fragments.push(ct);
        shared_secrets.push(ss);
    }

    // Combine the shared secrets
    let mut combined_secret = Vec::new();
    for ss in &shared_secrets {
        combined_secret.extend_from_slice(ss);
    }

    // Hash the combined secret to get the final shared secret
    let final_secret = hash(&combined_secret).to_vec();

    Ok((FragmentedCiphertext::new(ct_fragments)?, final_secret))
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
    let num_fragments = secret_key.num_fragments();

    if num_fragments != ciphertext.num_fragments() {
        return Err(Error::Fragmentation(
            "Number of fragments in secret key and ciphertext must match".to_string(),
        ));
    }

    let mut shared_secrets = Vec::with_capacity(num_fragments);

    // Decapsulate for each fragment
    for i in 0..num_fragments {
        let ss = crate::kem::decapsulate(&secret_key.fragments()[i], &ciphertext.fragments()[i])?;
        shared_secrets.push(ss);
    }

    // Combine the shared secrets
    let mut combined_secret = Vec::new();
    for ss in &shared_secrets {
        combined_secret.extend_from_slice(ss);
    }

    // Hash the combined secret to get the final shared secret
    let final_secret = hash(&combined_secret).to_vec();

    Ok(final_secret)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_fragmented_kem_roundtrip() {
        // Number of fragments
        let num_fragments = 4;

        // Generate a fragmented key pair
        let (pk, sk) = keygen(num_fragments).unwrap();

        // Encapsulate a shared secret
        let (ct, secret1) = encapsulate(&pk).unwrap();

        // Decapsulate the shared secret
        let secret2 = decapsulate(&sk, &ct).unwrap();

        // Verify that the shared secrets match
        assert_eq!(secret1, secret2);
    }
}