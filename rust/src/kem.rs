//! Key Encapsulation Mechanism (KEM) implementation.

use crate::error::Error;
use crate::lwe;
use crate::params::{CIPHERTEXT_BYTES, N, PUBLIC_KEY_BYTES, SECRET_KEY_BYTES, SECRET_LENGTH, SEED_LENGTH};
use crate::utils::{create_seeded_rng, encode_matrix, hash};
use rand::{CryptoRng, RngCore};
use zeroize::Zeroize;

/// A public key for the KEM.
#[derive(Clone, Debug)]
pub struct PublicKey {
    /// The encoded public key data.
    data: Vec<u8>,
}

/// A secret key for the KEM.
#[derive(Clone, Debug, Zeroize)]
#[zeroize(drop)]
pub struct SecretKey {
    /// The encoded secret key data.
    data: Vec<u8>,
}

/// A ciphertext for the KEM.
#[derive(Clone, Debug)]
pub struct Ciphertext {
    /// The encoded ciphertext data.
    data: Vec<u8>,
}

impl PublicKey {
    /// Creates a new public key from raw data.
    ///
    /// # Arguments
    ///
    /// * `data` - The raw public key data
    ///
    /// # Returns
    ///
    /// A new public key.
    pub fn new(data: Vec<u8>) -> Result<Self, Error> {
        if data.len() != PUBLIC_KEY_BYTES {
            return Err(Error::InvalidKeyFormat(
                format!("Public key must be {} bytes", PUBLIC_KEY_BYTES)
            ));
        }
        Ok(Self { data })
    }

    /// Returns the raw public key data.
    pub fn as_bytes(&self) -> &[u8] {
        &self.data
    }
}

impl SecretKey {
    /// Creates a new secret key from raw data.
    ///
    /// # Arguments
    ///
    /// * `data` - The raw secret key data
    ///
    /// # Returns
    ///
    /// A new secret key.
    pub fn new(data: Vec<u8>) -> Result<Self, Error> {
        if data.len() != SECRET_KEY_BYTES {
            return Err(Error::InvalidKeyFormat(
                format!("Secret key must be {} bytes", SECRET_KEY_BYTES)
            ));
        }
        Ok(Self { data })
    }

    /// Returns the raw secret key data.
    pub fn as_bytes(&self) -> &[u8] {
        &self.data
    }
}

impl Ciphertext {
    /// Creates a new ciphertext from raw data.
    ///
    /// # Arguments
    ///
    /// * `data` - The raw ciphertext data
    ///
    /// # Returns
    ///
    /// A new ciphertext.
    pub fn new(data: Vec<u8>) -> Result<Self, Error> {
        if data.len() != CIPHERTEXT_BYTES {
            return Err(Error::InvalidCiphertextFormat(
                format!("Ciphertext must be {} bytes", CIPHERTEXT_BYTES)
            ));
        }
        Ok(Self { data })
    }

    /// Returns the raw ciphertext data.
    pub fn as_bytes(&self) -> &[u8] {
        &self.data
    }
}

/// Generates a key pair for the KEM.
///
/// # Returns
///
/// A tuple containing the public key and secret key.
pub fn keygen() -> Result<(PublicKey, SecretKey), Error> {
    // Use the system's secure random number generator
    let mut rng = rand::thread_rng();
    keygen_with_rng(&mut rng)
}

/// Generates a key pair for the KEM using the provided random number generator.
///
/// # Arguments
///
/// * `rng` - Random number generator
///
/// # Returns
///
/// A tuple containing the public key and secret key.
pub fn keygen_with_rng<R: RngCore + CryptoRng>(
    rng: &mut R,
) -> Result<(PublicKey, SecretKey), Error> {
    // Generate a random seed for the LWE key generation
    let mut seed = vec![0u8; SEED_LENGTH];
    rng.fill_bytes(&mut seed);

    // Generate LWE key pair
    let (a, b, s) = lwe::keygen_with_seed(&seed)?;

    // Encode the public key: matrix A and vector b
    let a_bytes = encode_matrix(&a);
    let b_bytes = encode_matrix(&[b]);
    let mut pk_data = Vec::with_capacity(PUBLIC_KEY_BYTES);
    pk_data.extend_from_slice(&a_bytes);
    pk_data.extend_from_slice(&b_bytes);
    pk_data.extend_from_slice(&seed);

    // Encode the secret key: vector s and seed
    let s_bytes = encode_matrix(&[s]);
    let mut sk_data = Vec::with_capacity(SECRET_KEY_BYTES);
    sk_data.extend_from_slice(&s_bytes);
    sk_data.extend_from_slice(&seed);

    Ok((PublicKey::new(pk_data)?, SecretKey::new(sk_data)?))
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
    // Use the system's secure random number generator
    let mut rng = rand::thread_rng();
    encapsulate_with_rng(public_key, &mut rng)
}

/// Encapsulates a shared secret using the recipient's public key and the provided random number generator.
///
/// # Arguments
///
/// * `public_key` - The recipient's public key
/// * `rng` - Random number generator
///
/// # Returns
///
/// A tuple containing the ciphertext and the shared secret.
pub fn encapsulate_with_rng<R: RngCore + CryptoRng>(
    public_key: &PublicKey,
    rng: &mut R,
) -> Result<(Ciphertext, Vec<u8>), Error> {
    // Extract the public key components
    let pk_bytes = public_key.as_bytes();
    let a_b_bytes = &pk_bytes[0..pk_bytes.len() - SEED_LENGTH];
    let seed = &pk_bytes[pk_bytes.len() - SEED_LENGTH..];

    // Generate a random message
    let mut message = vec![0u8; SECRET_LENGTH];
    rng.fill_bytes(&mut message);

    // Encrypt the message using LWE
    let ciphertext_vector = lwe::encrypt(a_b_bytes, &message, seed)?;

    // Encode the ciphertext
    let c_bytes = encode_matrix(&[ciphertext_vector]);
    let mut ct_data = Vec::with_capacity(CIPHERTEXT_BYTES);
    ct_data.extend_from_slice(&c_bytes);
    ct_data.extend_from_slice(&message);

    // Hash the message to derive the shared secret
    let shared_secret = hash(&message).to_vec();

    Ok((Ciphertext::new(ct_data)?, shared_secret))
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
    // Extract the secret key components
    let sk_bytes = secret_key.as_bytes();
    let s_bytes = &sk_bytes[0..sk_bytes.len() - SEED_LENGTH];

    // Extract the ciphertext components
    let ct_bytes = ciphertext.as_bytes();
    let c_bytes = &ct_bytes[0..ct_bytes.len() - SECRET_LENGTH];
    let message = &ct_bytes[ct_bytes.len() - SECRET_LENGTH..];

    // Decrypt the message using LWE
    let decrypted = lwe::decrypt(c_bytes, s_bytes)?;

    // Verify the decryption
    if decrypted != message {
        return Err(Error::Decapsulation("Decapsulation failed: invalid ciphertext".to_string()));
    }

    // Hash the message to derive the shared secret
    let shared_secret = hash(message).to_vec();

    Ok(shared_secret)
}

#[cfg(test)]
mod tests {
    use super::*;

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
}