//! KeyPair example for TOPAY-Z512
//!
//! This example demonstrates key pair generation and management functionality.

use topayz512::{KeyPair, PrivateKey, PublicKey};

fn main() {
    println!("=== TOPAY-Z512 KeyPair Example ===\n");

    // Example 1: Generate a new key pair
    println!("1. Key Pair Generation:");
    let keypair = KeyPair::generate();
    println!(
        "   Private Key Size: {} bytes",
        keypair.private_key().to_bytes().len()
    );
    println!(
        "   Public Key Size: {} bytes",
        keypair.public_key().to_bytes().len()
    );
    println!("   Private Key: {}", keypair.private_key().to_hex());
    println!("   Public Key:  {}", keypair.public_key().to_hex());
    println!();

    // Example 2: Derive public key from private key
    println!("2. Public Key Derivation:");
    let private_key = PrivateKey::generate();
    let public_key1 = private_key.public_key();
    let public_key2 = private_key.public_key();

    println!("   Private Key: {}", private_key.to_hex());
    println!("   Derived Public Key 1: {}", public_key1.to_hex());
    println!("   Derived Public Key 2: {}", public_key2.to_hex());
    println!(
        "   Derivation is Deterministic: {}",
        public_key1 == public_key2
    );
    println!();

    // Example 3: Key pair from existing private key
    println!("3. KeyPair from Private Key:");
    let keypair_from_private = KeyPair::from_private_key(private_key.clone());
    println!("   Original Private: {}", private_key.to_hex());
    println!(
        "   KeyPair Private:  {}",
        keypair_from_private.private_key().to_hex()
    );
    println!(
        "   Keys Match: {}",
        private_key == *keypair_from_private.private_key()
    );
    println!();

    // Example 4: Hex serialization and deserialization
    println!("4. Hex Serialization:");
    let (private_hex, public_hex) = keypair.to_hex();
    println!("   Private Hex: {}", private_hex);
    println!("   Public Hex:  {}", public_hex);

    let private_restored = PrivateKey::from_hex(&private_hex).unwrap();
    let public_restored = PublicKey::from_hex(&public_hex).unwrap();

    println!(
        "   Private Key Restored: {}",
        private_restored == *keypair.private_key()
    );
    println!(
        "   Public Key Restored:  {}",
        public_restored == *keypair.public_key()
    );
    println!();

    // Example 5: Test with known values
    println!("5. Known Value Tests:");

    // Test with all zeros
    let zero_private = PrivateKey::from_bytes([0u8; 64]);
    let zero_public = zero_private.public_key();
    println!("   All-zero private key public: {}", zero_public.to_hex());

    // Test with all ones
    let ones_private = PrivateKey::from_bytes([0xFFu8; 64]);
    let ones_public = ones_private.public_key();
    println!("   All-ones private key public: {}", ones_public.to_hex());

    // Test deterministic behavior
    let zero_public2 = zero_private.public_key();
    println!(
        "   Zero key derivation consistent: {}",
        zero_public == zero_public2
    );
    println!();

    // Example 6: Multiple key pairs
    println!("6. Multiple Key Pairs:");
    let mut keypairs = Vec::new();

    for i in 0..5 {
        let kp = KeyPair::generate();
        println!(
            "   KeyPair {}: Private={}, Public={}",
            i + 1,
            &kp.private_key().to_hex()[..16],
            &kp.public_key().to_hex()[..16]
        );
        keypairs.push(kp);
    }

    // Verify all are different
    let mut all_different = true;
    for i in 0..keypairs.len() {
        for j in i + 1..keypairs.len() {
            if keypairs[i].private_key() == keypairs[j].private_key() {
                all_different = false;
                break;
            }
        }
    }
    println!("   All private keys different: {}", all_different);
    println!();

    // Example 7: Key recovery test
    println!("7. Key Recovery Test:");
    let original_keypair = KeyPair::generate();
    let private_hex = original_keypair.private_key().to_hex();

    // Simulate storing and retrieving the private key
    let recovered_private = PrivateKey::from_hex(&private_hex).unwrap();
    let recovered_public = recovered_private.public_key();

    println!(
        "   Original Private:  {}",
        original_keypair.private_key().to_hex()
    );
    println!("   Recovered Private: {}", recovered_private.to_hex());
    println!(
        "   Original Public:   {}",
        original_keypair.public_key().to_hex()
    );
    println!("   Recovered Public:  {}", recovered_public.to_hex());
    println!(
        "   Private Key Recovery: {}",
        *original_keypair.private_key() == recovered_private
    );
    println!(
        "   Public Key Recovery:  {}",
        *original_keypair.public_key() == recovered_public
    );
    println!();

    // Example 8: Performance test
    println!("8. Performance Test:");
    let start = std::time::Instant::now();

    for i in 0..100 {
        let _keypair = KeyPair::generate();
        if i == 0 {
            println!("   First key pair generated successfully");
        }
    }

    let duration = start.elapsed();
    println!("   100 key pair generations took: {:?}", duration);
    println!("   Average per generation: {:?}", duration / 100);
    println!();

    println!("=== KeyPair Example Complete ===");
    println!("All key operations completed successfully!");
}
