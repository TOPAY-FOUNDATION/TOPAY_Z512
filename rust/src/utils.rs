//! Utility functions for the topayz512 library.

use crate::error::Error;
use crate::params::{COEFF_BITS, N, Q, SIGMA};
use byteorder::{ByteOrder, LittleEndian};
use rand::{CryptoRng, Rng, RngCore, SeedableRng};
use rand_chacha::ChaCha20Rng;

/// Samples a value from a discrete Gaussian distribution with standard deviation sigma.
///
/// # Arguments
///
/// * `rng` - Random number generator
/// * `sigma` - Standard deviation
///
/// # Returns
///
/// A sample from the discrete Gaussian distribution.
pub fn sample_gaussian<R: RngCore + CryptoRng>(rng: &mut R, sigma: f64) -> i32 {
    // Box-Muller transform to generate Gaussian samples
    let u1 = rng.gen::<f64>();
    let u2 = rng.gen::<f64>();
    
    let radius = (-2.0 * u1.ln()).sqrt();
    let theta = 2.0 * std::f64::consts::PI * u2;
    
    let gaussian = radius * theta.cos() * sigma;
    
    // Round to nearest integer
    gaussian.round() as i32
}

/// Samples a uniform random value in the range [0, Q-1].
///
/// # Arguments
///
/// * `rng` - Random number generator
///
/// # Returns
///
/// A uniform random value in the range [0, Q-1].
pub fn sample_uniform<R: RngCore + CryptoRng>(rng: &mut R) -> u32 {
    rng.gen_range(0..Q)
}

/// Creates a deterministic random number generator from a seed.
///
/// # Arguments
///
/// * `seed` - The seed bytes
///
/// # Returns
///
/// A deterministic random number generator.
pub fn create_seeded_rng(seed: &[u8]) -> Result<ChaCha20Rng, Error> {
    if seed.len() < 32 {
        return Err(Error::RandomGeneration(
            "Seed must be at least 32 bytes long".to_string(),
        ));
    }
    
    let mut seed_array = [0u8; 32];
    seed_array.copy_from_slice(&seed[0..32]);
    
    Ok(ChaCha20Rng::from_seed(seed_array))
}

/// Encodes a matrix of u32 values into a byte array.
///
/// # Arguments
///
/// * `matrix` - The matrix to encode
///
/// # Returns
///
/// The encoded byte array.
pub fn encode_matrix(matrix: &[Vec<u32>]) -> Vec<u8> {
    let rows = matrix.len();
    let cols = if rows > 0 { matrix[0].len() } else { 0 };
    let mut bytes = Vec::with_capacity(rows * cols * 2); // 2 bytes per u32 (16 bits)
    
    for row in matrix {
        for &val in row {
            let mut buf = [0u8; 2];
            LittleEndian::write_u16(&mut buf, val as u16);
            bytes.extend_from_slice(&buf);
        }
    }
    
    bytes
}

/// Decodes a byte array into a matrix of u32 values.
///
/// # Arguments
///
/// * `bytes` - The encoded byte array
/// * `rows` - The number of rows in the matrix
/// * `cols` - The number of columns in the matrix
///
/// # Returns
///
/// The decoded matrix.
pub fn decode_matrix(bytes: &[u8], rows: usize, cols: usize) -> Result<Vec<Vec<u32>>, Error> {
    if bytes.len() < rows * cols * 2 {
        return Err(Error::InvalidParameter(
            "Byte array is too small for the specified matrix dimensions".to_string(),
        ));
    }
    
    let mut matrix = vec![vec![0u32; cols]; rows];
    
    for i in 0..rows {
        for j in 0..cols {
            let idx = 2 * (i * cols + j);
            matrix[i][j] = LittleEndian::read_u16(&bytes[idx..idx + 2]) as u32;
        }
    }
    
    Ok(matrix)
}

/// Performs modular addition: (a + b) mod Q.
///
/// # Arguments
///
/// * `a` - First operand
/// * `b` - Second operand
///
/// # Returns
///
/// The result of (a + b) mod Q.
pub fn mod_add(a: u32, b: u32) -> u32 {
    (a + b) % Q
}

/// Performs modular subtraction: (a - b) mod Q.
///
/// # Arguments
///
/// * `a` - First operand
/// * `b` - Second operand
///
/// # Returns
///
/// The result of (a - b) mod Q.
pub fn mod_sub(a: u32, b: u32) -> u32 {
    (a + Q - (b % Q)) % Q
}

/// Performs modular multiplication: (a * b) mod Q.
///
/// # Arguments
///
/// * `a` - First operand
/// * `b` - Second operand
///
/// # Returns
///
/// The result of (a * b) mod Q.
pub fn mod_mul(a: u32, b: u32) -> u32 {
    ((a as u64 * b as u64) % (Q as u64)) as u32
}

/// Computes the hash of a message using SHA3-512.
///
/// # Arguments
///
/// * `message` - The message to hash
///
/// # Returns
///
/// The hash of the message.
pub fn hash(message: &[u8]) -> [u8; 64] {
    use sha3::{Digest, Sha3_512};
    
    let mut hasher = Sha3_512::new();
    hasher.update(message);
    let result = hasher.finalize();
    
    let mut output = [0u8; 64];
    output.copy_from_slice(&result);
    output
}

/// Generates a random matrix of dimensions rows x cols with entries in [0, Q-1].
///
/// # Arguments
///
/// * `rng` - Random number generator
/// * `rows` - Number of rows
/// * `cols` - Number of columns
///
/// # Returns
///
/// A random matrix.
pub fn random_matrix<R: RngCore + CryptoRng>(
    rng: &mut R,
    rows: usize,
    cols: usize,
) -> Vec<Vec<u32>> {
    let mut matrix = vec![vec![0u32; cols]; rows];
    
    for i in 0..rows {
        for j in 0..cols {
            matrix[i][j] = sample_uniform(rng);
        }
    }
    
    matrix
}

/// Generates a random error vector of length n with entries sampled from a discrete Gaussian.
///
/// # Arguments
///
/// * `rng` - Random number generator
/// * `n` - Length of the vector
///
/// # Returns
///
/// A random error vector.
pub fn random_error_vector<R: RngCore + CryptoRng>(rng: &mut R, n: usize) -> Vec<u32> {
    let mut vector = vec![0u32; n];
    
    for i in 0..n {
        let error = sample_gaussian(rng, SIGMA);
        // Convert to unsigned and apply modulo Q
        vector[i] = ((error % (Q as i32) + (Q as i32)) % (Q as i32)) as u32;
    }
    
    vector
}

/// Adds the missing dependency to Cargo.toml
/// ```
/// # Cargo.toml
/// [dependencies]
/// sha3 = "0.10"
/// ```