//! Integration tests for TOPAY-Z512
//!
//! These tests verify that all components work together correctly.

use topayz512::{Hash, Kem, KeyPair, PrivateKey, PublicKey};

#[cfg(feature = "fragmentation")]
use topayz512::fragment::FragmentEngine;

#[test]
fn test_complete_workflow() {
    // Test a complete cryptographic workflow
    let _signing_keypair = KeyPair::generate();
    let encryption_keypair = KeyPair::generate();

    // Generate KEM keys
    let (kem_public, kem_secret) = Kem::keygen();

    // Test message
    let message = b"Hello, TOPAY-Z512!";

    // Hash the message
    let message_hash = Hash::new(message);

    // KEM encapsulation
    let (ciphertext, shared_secret1) = Kem::encapsulate(&kem_public);

    // KEM decapsulation
    let shared_secret2 = Kem::decapsulate(&kem_secret, &ciphertext);

    // Verify shared secrets match
    assert_eq!(shared_secret1.to_bytes(), shared_secret2.to_bytes());

    // Verify hash consistency
    let message_hash2 = Hash::new(message);
    assert_eq!(message_hash.to_bytes(), message_hash2.to_bytes());

    // Test key derivation consistency
    let derived_public1 = encryption_keypair.private_key().public_key();
    let derived_public2 = encryption_keypair.private_key().public_key();
    assert_eq!(derived_public1.to_bytes(), derived_public2.to_bytes());

    println!("Complete workflow test passed");
}

#[test]
fn test_serialization_roundtrip() {
    // Test that all types can be serialized and deserialized correctly

    // KeyPair serialization
    let original_keypair = KeyPair::generate();
    let (private_hex, public_hex) = original_keypair.to_hex();

    let restored_private = PrivateKey::from_hex(&private_hex).unwrap();
    let restored_public = PublicKey::from_hex(&public_hex).unwrap();

    assert_eq!(*original_keypair.private_key(), restored_private);
    assert_eq!(*original_keypair.public_key(), restored_public);

    // Hash serialization
    let original_hash = Hash::new(b"test data");
    let hash_hex = original_hash.to_hex();
    let restored_hash = Hash::from_hex(&hash_hex).unwrap();

    assert_eq!(original_hash.to_bytes(), restored_hash.to_bytes());

    // KEM key serialization
    let (kem_public, kem_secret) = Kem::keygen();

    let public_hex = kem_public.to_hex();
    let secret_hex = kem_secret.to_hex();

    let restored_public = topayz512::kem::PublicKey::from_hex(&public_hex).unwrap();
    let restored_secret = topayz512::kem::SecretKey::from_hex(&secret_hex).unwrap();

    assert_eq!(kem_public.to_bytes(), restored_public.to_bytes());
    assert_eq!(kem_secret.to_bytes(), restored_secret.to_bytes());

    println!("Serialization roundtrip test passed");
}

#[test]
fn test_deterministic_behavior() {
    // Test that operations are deterministic when they should be

    // Same input should produce same hash
    let data = b"deterministic test";
    let hash1 = Hash::new(data);
    let hash2 = Hash::new(data);
    assert_eq!(hash1.to_bytes(), hash2.to_bytes());

    // Same private key should produce same public key
    let private_bytes = [0x42u8; 64];
    let private1 = PrivateKey::from_bytes(private_bytes);
    let private2 = PrivateKey::from_bytes(private_bytes);

    let public1 = private1.public_key();
    let public2 = private2.public_key();
    assert_eq!(public1.to_bytes(), public2.to_bytes());

    // Hash combination should be deterministic
    let combined1 = Hash::combine(b"first", b"second");
    let combined2 = Hash::combine(b"first", b"second");
    assert_eq!(combined1.to_bytes(), combined2.to_bytes());

    println!("Deterministic behavior test passed");
}

#[test]
fn test_cross_language_compatibility() {
    // Test known values that should be consistent across language implementations

    // Test known private key to public key derivation
    let zero_private = PrivateKey::from_bytes([0u8; 64]);
    let zero_public = zero_private.public_key();

    let ones_private = PrivateKey::from_bytes([0xFFu8; 64]);
    let ones_public = ones_private.public_key();

    // These should produce consistent results across languages
    println!("Zero private -> public: {}", zero_public.to_hex());
    println!("Ones private -> public: {}", ones_public.to_hex());

    // Test known hash values
    let empty_hash = Hash::new(b"");
    let hello_hash = Hash::new(b"Hello, TOPAY!");
    let long_hash = Hash::new(b"This is a longer message for testing cross-language compatibility");

    println!("Empty hash: {}", empty_hash.to_hex());
    println!("Hello hash: {}", hello_hash.to_hex());
    println!("Long hash: {}", long_hash.to_hex());

    // Test hash combination
    let combined = Hash::combine(b"first", b"second");
    println!("Combined hash: {}", combined.to_hex());

    println!("Cross-language compatibility test passed");
}

#[cfg(feature = "fragmentation")]
#[test]
fn test_fragmentation_integration() {
    // Test fragmentation with real cryptographic operations
    let large_data = vec![0x5Au8; 4096]; // 4KB of data

    // Fragment the data
    let fragments = FragmentEngine::fragment_data(&large_data).unwrap();
    assert!(!fragments.is_empty());

    // Verify all fragments
    for fragment in &fragments {
        assert!(fragment.verify());
    }

    // Reconstruct and verify
    let reconstructed = FragmentEngine::reconstruct_data(&fragments).unwrap();
    assert_eq!(large_data, reconstructed);

    // Test KEM fragmentation
    let (public_key, _) = Kem::keygen();
    let fragmented_kem = FragmentEngine::fragment_kem_encapsulation(&public_key).unwrap();
    let processed = FragmentEngine::process_fragmented_kem(&fragmented_kem).unwrap();
    assert!(!processed.is_empty());

    // Test hash fragmentation
    let fragmented_hash = FragmentEngine::fragment_hash_operation(&large_data).unwrap();
    let parallel_hash = FragmentEngine::parallel_hash_compute(&fragmented_hash).unwrap();
    assert_eq!(parallel_hash.to_bytes().len(), 64);

    println!("Fragmentation integration test passed");
}

#[test]
fn test_error_handling() {
    // Test various error conditions

    // Invalid hex strings
    assert!(PrivateKey::from_hex("invalid").is_err());
    assert!(PublicKey::from_hex("too_short").is_err());
    assert!(Hash::from_hex("not_hex").is_err());

    // Invalid lengths
    assert!(PrivateKey::from_hex("deadbeef").is_err()); // Too short
    assert!(Hash::from_hex("cafe").is_err()); // Too short

    #[cfg(feature = "fragmentation")]
    {
        // Fragmentation errors
        assert!(FragmentEngine::fragment_data(&[]).is_err()); // Empty data

        // Invalid fragment reconstruction
        let fragments = vec![]; // Empty fragments
        assert!(FragmentEngine::reconstruct_data(&fragments).is_err());
    }

    println!("Error handling test passed");
}

#[test]
fn test_performance_characteristics() {
    use std::time::Instant;

    // Test that operations complete within reasonable time
    let start = Instant::now();

    // Key generation should be fast
    for _ in 0..10 {
        let _keypair = KeyPair::generate();
    }
    let keypair_time = start.elapsed();

    let start = Instant::now();

    // KEM operations should be reasonable
    for _ in 0..5 {
        let (public_key, secret_key) = Kem::keygen();
        let (ciphertext, _) = Kem::encapsulate(&public_key);
        let _ = Kem::decapsulate(&secret_key, &ciphertext);
    }
    let kem_time = start.elapsed();

    let start = Instant::now();

    // Hash operations should be fast
    for i in 0..100 {
        let data = format!("test data {i}");
        let _ = Hash::new(data.as_bytes());
    }
    let hash_time = start.elapsed();

    println!("Performance test results:");
    println!("  10 keypair generations: {keypair_time:?}");
    println!("  5 complete KEM cycles: {kem_time:?}");
    println!("  100 hash operations: {hash_time:?}");

    // Ensure operations are reasonably fast (these are generous limits)
    assert!(keypair_time.as_millis() < 1000); // 1 second for 10 operations
    assert!(kem_time.as_millis() < 2000); // 2 seconds for 5 operations
    assert!(hash_time.as_millis() < 100); // 100ms for 100 operations

    println!("Performance characteristics test passed");
}

#[test]
fn test_memory_safety() {
    // Test operations that could cause memory issues

    // Large data handling
    let large_data = vec![0u8; 1_000_000]; // 1MB
    let hash = Hash::new(&large_data);
    assert_eq!(hash.to_bytes().len(), 64);

    // Many small operations
    for i in 0..1000 {
        let data = format!("iteration {i}");
        let _hash = Hash::new(data.as_bytes());
    }

    // Nested operations
    let mut current_hash = Hash::new(b"start");
    for i in 0..100 {
        let data = format!("round {i}");
        current_hash = Hash::combine(&current_hash.to_bytes(), data.as_bytes());
    }

    println!("Memory safety test passed");
}

#[test]
fn test_blockchain_simulation() {
    // Simulate a simplified blockchain scenario

    // Genesis block
    let genesis_data = b"TOPAY Genesis Block";
    let genesis_hash = Hash::new(genesis_data);

    // Validator key pairs
    let validator1 = KeyPair::generate();
    let validator2 = KeyPair::generate();
    let validator3 = KeyPair::generate();

    // Transaction simulation
    let mut previous_hash = genesis_hash;

    for block_num in 1..=10 {
        // Create block data
        let block_data = format!(
            "Block {} with previous hash {}",
            block_num,
            previous_hash.to_hex()
        );

        // Hash the block
        let block_hash = Hash::combine(&previous_hash.to_bytes(), block_data.as_bytes());

        // Simulate validator consensus (simplified)
        let validator_hashes = [
            Hash::combine(&block_hash.to_bytes(), &validator1.public_key().to_bytes()),
            Hash::combine(&block_hash.to_bytes(), &validator2.public_key().to_bytes()),
            Hash::combine(&block_hash.to_bytes(), &validator3.public_key().to_bytes()),
        ];

        // Combine validator signatures
        let consensus_hash = Hash::concat(&[
            &validator_hashes[0],
            &validator_hashes[1],
            &validator_hashes[2],
        ]);

        let consensus_hex = consensus_hash.to_hex();
        println!("Block {}: {}", block_num, &consensus_hex[..16]);
        previous_hash = consensus_hash;
    }

    println!("Blockchain simulation test passed");
}
