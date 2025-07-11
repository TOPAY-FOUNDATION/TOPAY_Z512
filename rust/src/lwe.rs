//! Learning With Errors (LWE) implementation.

use crate::error::Error;
use crate::params::{COEFF_BITS, N, Q, SEED_LENGTH};
use crate::utils::{create_seeded_rng, decode_matrix, mod_add, mod_mul, mod_sub, random_error_vector, random_matrix};
use byteorder::{ByteOrder, LittleEndian};
use rand::{CryptoRng, RngCore};

/// Generates an LWE key pair.
///
/// # Returns
///
/// A tuple containing the matrix A, vector b, and vector s.
pub fn keygen() -> Result<(Vec<Vec<u32>>, Vec<u32>, Vec<u32>), Error> {
    // Use the system's secure random number generator
    let mut rng = rand::thread_rng();
    
    // Generate a random seed
    let mut seed = vec![0u8; SEED_LENGTH];
    rng.fill_bytes(&mut seed);
    
    keygen_with_seed(&seed)
}

/// Generates an LWE key pair using a seed.
///
/// # Arguments
///
/// * `seed` - The seed for the random number generator
///
/// # Returns
///
/// A tuple containing the matrix A, vector b, and vector s.
pub fn keygen_with_seed(seed: &[u8]) -> Result<(Vec<Vec<u32>>, Vec<u32>, Vec<u32>), Error> {
    // Create a deterministic random number generator from the seed
    let mut rng = create_seeded_rng(seed)?;
    
    // Generate a random matrix A
    let a = random_matrix(&mut rng, N, N);
    
    // Generate a random secret vector s
    let s = random_error_vector(&mut rng, N);
    
    // Generate a random error vector e
    let e = random_error_vector(&mut rng, N);
    
    // Compute b = A·s + e
    let mut b = vec![0u32; N];
    for i in 0..N {
        for j in 0..N {
            b[i] = mod_add(b[i], mod_mul(a[i][j], s[j]));
        }
        b[i] = mod_add(b[i], e[i]);
    }
    
    Ok((a, b, s))
}

/// Encrypts a message using LWE.
///
/// # Arguments
///
/// * `public_key_bytes` - The encoded public key (matrix A and vector b)
/// * `message` - The message to encrypt
/// * `seed` - The seed for the random number generator
///
/// # Returns
///
/// The encrypted ciphertext vector.
pub fn encrypt(public_key_bytes: &[u8], message: &[u8], seed: &[u8]) -> Result<Vec<u32>, Error> {
    // Decode the public key
    let a_size = N * N * 2; // 2 bytes per coefficient
    let b_size = N * 2;
    
    if public_key_bytes.len() < a_size + b_size {
        return Err(Error::InvalidKeyFormat(
            "Public key bytes are too short".to_string(),
        ));
    }
    
    let a_bytes = &public_key_bytes[0..a_size];
    let b_bytes = &public_key_bytes[a_size..a_size + b_size];
    
    let a = decode_matrix(a_bytes, N, N)?;
    let b = decode_matrix(b_bytes, 1, N)?[0].clone();
    
    // Create a deterministic random number generator from the seed
    let mut rng = create_seeded_rng(seed)?;
    
    // Generate a random vector r
    let r = random_error_vector(&mut rng, N);
    
    // Compute v = r·A
    let mut v = vec![0u32; N];
    for i in 0..N {
        for j in 0..N {
            v[j] = mod_add(v[j], mod_mul(r[i], a[i][j]));
        }
    }
    
    // Compute c = r·b + encode(message)
    let mut c = 0u32;
    for i in 0..N {
        c = mod_add(c, mod_mul(r[i], b[i]));
    }
    
    // Encode the message
    let message_bits = message.len() * 8;
    let bits_per_coeff = Q.ilog2() as usize - 1; // Leave 1 bit for noise
    let coeffs_needed = (message_bits + bits_per_coeff - 1) / bits_per_coeff;
    
    if coeffs_needed > 1 {
        return Err(Error::Encapsulation(
            "Message is too large for single coefficient encryption".to_string(),
        ));
    }
    
    // Encode the message into a single coefficient
    let mut message_val = 0u32;
    for (i, &byte) in message.iter().enumerate() {
        message_val |= (byte as u32) << (i * 8);
    }
    
    // Scale the message to fit in the modulus
    let scaling_factor = (Q / 256) as u32;
    let encoded_message = (message_val * scaling_factor) % Q;
    
    // Add the encoded message to c
    c = mod_add(c, encoded_message);
    
    // Return the ciphertext (v, c)
    let mut ciphertext = v;
    ciphertext.push(c);
    
    Ok(ciphertext)
}

/// Decrypts a ciphertext using LWE.
///
/// # Arguments
///
/// * `ciphertext_bytes` - The encoded ciphertext
/// * `secret_key_bytes` - The encoded secret key (vector s)
///
/// # Returns
///
/// The decrypted message.
pub fn decrypt(ciphertext_bytes: &[u8], secret_key_bytes: &[u8]) -> Result<Vec<u8>, Error> {
    // Decode the ciphertext
    let ct_size = (N + 1) * 2; // 2 bytes per coefficient
    
    if ciphertext_bytes.len() < ct_size {
        return Err(Error::InvalidCiphertextFormat(
            "Ciphertext bytes are too short".to_string(),
        ));
    }
    
    let ciphertext = decode_matrix(ciphertext_bytes, 1, N + 1)?[0].clone();
    let v = &ciphertext[0..N];
    let c = ciphertext[N];
    
    // Decode the secret key
    let sk_size = N * 2;
    
    if secret_key_bytes.len() < sk_size {
        return Err(Error::InvalidKeyFormat(
            "Secret key bytes are too short".to_string(),
        ));
    }
    
    let s = decode_matrix(secret_key_bytes, 1, N)?[0].clone();
    
    // Compute m = c - v·s
    let mut vs = 0u32;
    for i in 0..N {
        vs = mod_add(vs, mod_mul(v[i], s[i]));
    }
    
    let m = mod_sub(c, vs);
    
    // Decode the message
    let scaling_factor = (Q / 256) as u32;
    let message_val = (m + scaling_factor / 2) / scaling_factor; // Round to nearest
    
    // Convert to bytes
    let mut message = Vec::new();
    let mut remaining = message_val;
    
    while remaining > 0 || message.is_empty() {
        message.push((remaining & 0xFF) as u8);
        remaining >>= 8;
    }
    
    Ok(message)
}

#[cfg(test)]
mod tests {
    use super::*;
    use rand::rngs::OsRng;
    
    #[test]
    fn test_lwe_roundtrip() {
        // Generate a key pair
        let (a, b, s) = keygen().unwrap();
        
        // Create a test message
        let message = b"test message";
        
        // Generate a random seed
        let mut seed = vec![0u8; SEED_LENGTH];
        OsRng.fill_bytes(&mut seed);
        
        // Encode the public key
        let mut a_bytes = Vec::new();
        for row in &a {
            for &val in row {
                let mut buf = [0u8; 2];
                LittleEndian::write_u16(&mut buf, val as u16);
                a_bytes.extend_from_slice(&buf);
            }
        }
        
        let mut b_bytes = Vec::new();
        for &val in &b {
            let mut buf = [0u8; 2];
            LittleEndian::write_u16(&mut buf, val as u16);
            b_bytes.extend_from_slice(&buf);
        }
        
        let mut public_key_bytes = Vec::new();
        public_key_bytes.extend_from_slice(&a_bytes);
        public_key_bytes.extend_from_slice(&b_bytes);
        
        // Encode the secret key
        let mut s_bytes = Vec::new();
        for &val in &s {
            let mut buf = [0u8; 2];
            LittleEndian::write_u16(&mut buf, val as u16);
            s_bytes.extend_from_slice(&buf);
        }
        
        // Encrypt the message
        let ciphertext = encrypt(&public_key_bytes, message, &seed).unwrap();
        
        // Encode the ciphertext
        let mut ciphertext_bytes = Vec::new();
        for &val in &ciphertext {
            let mut buf = [0u8; 2];
            LittleEndian::write_u16(&mut buf, val as u16);
            ciphertext_bytes.extend_from_slice(&buf);
        }
        
        // Decrypt the message
        let decrypted = decrypt(&ciphertext_bytes, &s_bytes).unwrap();
        
        // Verify that the decrypted message matches the original
        assert_eq!(decrypted, message);
    }
}