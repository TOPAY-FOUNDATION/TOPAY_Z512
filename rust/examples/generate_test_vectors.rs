//! Utility to generate test vectors for the topayz512 library.

use topayz512::{decapsulate, encapsulate, keygen, Ciphertext, PublicKey, SecretKey};
use rand::{RngCore, SeedableRng};
use rand_chacha::ChaCha20Rng;
use serde::{Deserialize, Serialize};
use std::fs::File;
use std::io::{Read, Write};
use std::path::Path;

#[cfg(feature = "fragmentation")]
use topayz512::{fragmented, FragmentedCiphertext, FragmentedPublicKey, FragmentedSecretKey};

#[derive(Serialize, Deserialize)]
struct Parameters {
    N: usize,
    Q: u32,
    sigma: f64,
    secret_length: usize,
}

#[derive(Serialize, Deserialize)]
struct TestVector {
    test_id: String,
    description: String,
    seed: String,
    public_key: String,
    secret_key: String,
    message: String,
    ciphertext: String,
    shared_secret: String,
}

#[derive(Serialize, Deserialize)]
struct FragmentedTestVector {
    test_id: String,
    description: String,
    num_fragments: usize,
    seed: String,
    fragmented_public_key: Vec<String>,
    fragmented_secret_key: Vec<String>,
    message: String,
    fragmented_ciphertext: Vec<String>,
    shared_secret: String,
}

#[derive(Serialize, Deserialize)]
struct TestVectors {
    description: String,
    version: String,
    parameters: Parameters,
    test_vectors: Vec<TestVector>,
    #[serde(default)]
    fragmented_test_vectors: Vec<FragmentedTestVector>,
}

/// Converts bytes to a hex string.
fn bytes_to_hex(bytes: &[u8]) -> String {
    bytes.iter().map(|b| format!("{:02x}", b)).collect()
}

/// Converts a hex string to bytes.
fn hex_to_bytes(hex: &str) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let mut bytes = Vec::new();
    let mut chars = hex.chars();
    
    while let (Some(a), Some(b)) = (chars.next(), chars.next()) {
        let byte = u8::from_str_radix(&format!("{}{}", a, b), 16)?;
        bytes.push(byte);
    }
    
    Ok(bytes)
}

/// Generates test vectors for the KEM operations.
fn generate_test_vectors() -> Result<TestVectors, Box<dyn std::error::Error>> {
    let mut test_vectors = Vec::new();
    
    // Test vector 1: Basic KEM operation with fixed seed
    let seed_hex = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f";
    let seed = hex_to_bytes(seed_hex)?;
    let message_hex = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef";
    let message = hex_to_bytes(message_hex)?;
    
    // Create a deterministic RNG from the seed
    let mut seed_array = [0u8; 32];
    seed_array.copy_from_slice(&seed);
    let mut rng = ChaCha20Rng::from_seed(seed_array);
    
    // Generate a key pair
    let (pk, sk) = keygen()?;
    
    // Encapsulate a shared secret
    let (ct, shared_secret) = encapsulate(&pk)?;
    
    test_vectors.push(TestVector {
        test_id: "kem_basic_1".to_string(),
        description: "Basic KEM operation with fixed seed".to_string(),
        seed: seed_hex.to_string(),
        public_key: bytes_to_hex(pk.as_bytes()),
        secret_key: bytes_to_hex(sk.as_bytes()),
        message: message_hex.to_string(),
        ciphertext: bytes_to_hex(ct.as_bytes()),
        shared_secret: bytes_to_hex(&shared_secret),
    });
    
    // Test vector 2: Basic KEM operation with different fixed seed
    let seed_hex = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
    let seed = hex_to_bytes(seed_hex)?;
    let message_hex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";
    let message = hex_to_bytes(message_hex)?;
    
    // Create a deterministic RNG from the seed
    let mut seed_array = [0u8; 32];
    seed_array.copy_from_slice(&seed);
    let mut rng = ChaCha20Rng::from_seed(seed_array);
    
    // Generate a key pair
    let (pk, sk) = keygen()?;
    
    // Encapsulate a shared secret
    let (ct, shared_secret) = encapsulate(&pk)?;
    
    test_vectors.push(TestVector {
        test_id: "kem_basic_2".to_string(),
        description: "Basic KEM operation with different fixed seed".to_string(),
        seed: seed_hex.to_string(),
        public_key: bytes_to_hex(pk.as_bytes()),
        secret_key: bytes_to_hex(sk.as_bytes()),
        message: message_hex.to_string(),
        ciphertext: bytes_to_hex(ct.as_bytes()),
        shared_secret: bytes_to_hex(&shared_secret),
    });
    
    // Generate fragmented test vectors if the feature is enabled
    let fragmented_test_vectors = if cfg!(feature = "fragmentation") {
        #[cfg(feature = "fragmentation")]
        {
            let mut vectors = Vec::new();
            
            // Test vector 1: Fragmented KEM operation with 2 fragments
            let seed_hex = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f";
            let seed = hex_to_bytes(seed_hex)?;
            let message_hex = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef";
            let message = hex_to_bytes(message_hex)?;
            let num_fragments = 2;
            
            // Create a deterministic RNG from the seed
            let mut seed_array = [0u8; 32];
            seed_array.copy_from_slice(&seed);
            let mut rng = ChaCha20Rng::from_seed(seed_array);
            
            // Generate a fragmented key pair
            let (fpk, fsk) = fragmented::keygen(num_fragments)?;
            
            // Encapsulate a shared secret
            let (fct, shared_secret) = fragmented::encapsulate(&fpk)?;
            
            vectors.push(FragmentedTestVector {
                test_id: "fragmented_kem_1".to_string(),
                description: "Fragmented KEM operation with 2 fragments".to_string(),
                num_fragments,
                seed: seed_hex.to_string(),
                fragmented_public_key: fpk.fragments().iter().map(|pk| bytes_to_hex(pk.as_bytes())).collect(),
                fragmented_secret_key: fsk.fragments().iter().map(|sk| bytes_to_hex(sk.as_bytes())).collect(),
                message: message_hex.to_string(),
                fragmented_ciphertext: fct.fragments().iter().map(|ct| bytes_to_hex(ct.as_bytes())).collect(),
                shared_secret: bytes_to_hex(&shared_secret),
            });
            
            // Test vector 2: Fragmented KEM operation with 4 fragments
            let seed_hex = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
            let seed = hex_to_bytes(seed_hex)?;
            let message_hex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";
            let message = hex_to_bytes(message_hex)?;
            let num_fragments = 4;
            
            // Create a deterministic RNG from the seed
            let mut seed_array = [0u8; 32];
            seed_array.copy_from_slice(&seed);
            let mut rng = ChaCha20Rng::from_seed(seed_array);
            
            // Generate a fragmented key pair
            let (fpk, fsk) = fragmented::keygen(num_fragments)?;
            
            // Encapsulate a shared secret
            let (fct, shared_secret) = fragmented::encapsulate(&fpk)?;
            
            vectors.push(FragmentedTestVector {
                test_id: "fragmented_kem_2".to_string(),
                description: "Fragmented KEM operation with 4 fragments".to_string(),
                num_fragments,
                seed: seed_hex.to_string(),
                fragmented_public_key: fpk.fragments().iter().map(|pk| bytes_to_hex(pk.as_bytes())).collect(),
                fragmented_secret_key: fsk.fragments().iter().map(|sk| bytes_to_hex(sk.as_bytes())).collect(),
                message: message_hex.to_string(),
                fragmented_ciphertext: fct.fragments().iter().map(|ct| bytes_to_hex(ct.as_bytes())).collect(),
                shared_secret: bytes_to_hex(&shared_secret),
            });
            
            vectors
        }
        #[cfg(not(feature = "fragmentation"))]
        {
            Vec::new()
        }
    } else {
        Vec::new()
    };
    
    Ok(TestVectors {
        description: "Test vectors for topayz512 KEM operations".to_string(),
        version: "1.0.0".to_string(),
        parameters: Parameters {
            N: topayz512::params::N,
            Q: topayz512::params::Q,
            sigma: topayz512::params::SIGMA,
            secret_length: topayz512::params::SECRET_LENGTH,
        },
        test_vectors,
        fragmented_test_vectors,
    })
}

/// Verifies the test vectors by performing the KEM operations and comparing the results.
fn verify_test_vectors(vectors: &TestVectors) -> Result<(), Box<dyn std::error::Error>> {
    println!("Verifying test vectors...");
    
    for vector in &vectors.test_vectors {
        println!("Verifying test vector: {}", vector.test_id);
        
        // Decode the keys and ciphertext
        let pk_bytes = hex_to_bytes(&vector.public_key)?;
        let sk_bytes = hex_to_bytes(&vector.secret_key)?;
        let ct_bytes = hex_to_bytes(&vector.ciphertext)?;
        
        let pk = PublicKey::new(pk_bytes)?;
        let sk = SecretKey::new(sk_bytes)?;
        let ct = Ciphertext::new(ct_bytes)?;
        
        // Decapsulate the shared secret
        let shared_secret = decapsulate(&sk, &ct)?;
        let shared_secret_hex = bytes_to_hex(&shared_secret);
        
        // Verify that the shared secret matches
        if shared_secret_hex != vector.shared_secret {
            return Err(format!(
                "Shared secret mismatch for test vector {}: expected {}, got {}",
                vector.test_id, vector.shared_secret, shared_secret_hex
            ).into());
        }
    }
    
    // Verify fragmented test vectors if the feature is enabled
    #[cfg(feature = "fragmentation")]
    for vector in &vectors.fragmented_test_vectors {
        println!("Verifying fragmented test vector: {}", vector.test_id);
        
        // Decode the fragmented keys and ciphertext
        let mut pk_fragments = Vec::new();
        let mut sk_fragments = Vec::new();
        let mut ct_fragments = Vec::new();
        
        for pk_hex in &vector.fragmented_public_key {
            let pk_bytes = hex_to_bytes(pk_hex)?;
            pk_fragments.push(PublicKey::new(pk_bytes)?);
        }
        
        for sk_hex in &vector.fragmented_secret_key {
            let sk_bytes = hex_to_bytes(sk_hex)?;
            sk_fragments.push(SecretKey::new(sk_bytes)?);
        }
        
        for ct_hex in &vector.fragmented_ciphertext {
            let ct_bytes = hex_to_bytes(ct_hex)?;
            ct_fragments.push(Ciphertext::new(ct_bytes)?);
        }
        
        let fpk = FragmentedPublicKey::new(pk_fragments)?;
        let fsk = FragmentedSecretKey::new(sk_fragments)?;
        let fct = FragmentedCiphertext::new(ct_fragments)?;
        
        // Decapsulate the shared secret
        let shared_secret = fragmented::decapsulate(&fsk, &fct)?;
        let shared_secret_hex = bytes_to_hex(&shared_secret);
        
        // Verify that the shared secret matches
        if shared_secret_hex != vector.shared_secret {
            return Err(format!(
                "Shared secret mismatch for fragmented test vector {}: expected {}, got {}",
                vector.test_id, vector.shared_secret, shared_secret_hex
            ).into());
        }
    }
    
    println!("All test vectors verified successfully!");
    Ok(())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Generate test vectors
    let test_vectors = generate_test_vectors()?;
    
    // Verify the test vectors
    verify_test_vectors(&test_vectors)?;
    
    // Write the test vectors to a file
    let output_path = Path::new("../test-vectors/001_basic.json");
    let mut file = File::create(output_path)?;
    let json = serde_json::to_string_pretty(&test_vectors)?;
    file.write_all(json.as_bytes())?;
    
    println!("Test vectors written to {}", output_path.display());
    
    Ok(())
}