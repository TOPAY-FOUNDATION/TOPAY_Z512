//! Cryptographic parameters for the topayz512 library.

/// The dimension of the LWE problem (number of variables).
/// This is a key security parameter that determines the size of matrices and vectors.
pub const N: usize = 1024;

/// The modulus for the LWE problem.
/// All arithmetic operations are performed modulo Q.
pub const Q: u32 = 65537; // 2^16 + 1, a prime number

/// The standard deviation for the error distribution.
/// This is a security parameter that affects the hardness of the LWE problem.
pub const SIGMA: f64 = 3.2;

/// The length of the shared secret in bytes.
pub const SECRET_LENGTH: usize = 64; // 512 bits

/// The length of the seed used for pseudorandom generation in bytes.
pub const SEED_LENGTH: usize = 32;

/// The number of bits used for encoding each coefficient.
pub const COEFF_BITS: usize = 16; // log2(Q)

/// The number of bytes needed to store each coefficient.
pub const COEFF_BYTES: usize = (COEFF_BITS + 7) / 8;

/// The size of the public key in bytes.
pub const PUBLIC_KEY_BYTES: usize = N * N * COEFF_BYTES + SEED_LENGTH;

/// The size of the secret key in bytes.
pub const SECRET_KEY_BYTES: usize = N * COEFF_BYTES + SEED_LENGTH;

/// The size of the ciphertext in bytes.
pub const CIPHERTEXT_BYTES: usize = N * COEFF_BYTES + SECRET_LENGTH;

/// The maximum number of fragments allowed for fragmentation.
#[cfg(feature = "fragmentation")]
pub const MAX_FRAGMENTS: usize = 16;

/// The minimum number of fragments required for fragmentation.
#[cfg(feature = "fragmentation")]
pub const MIN_FRAGMENTS: usize = 2;