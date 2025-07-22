//! KEM example for TOPAY-Z512
//! 
//! This example demonstrates the Key Encapsulation Mechanism (KEM) functionality.

use topayz512::Kem;

fn main() {
    println!("=== TOPAY-Z512 KEM Example ===\n");

    // Example 1: Key generation
    println!("1. KEM Key Generation:");
    let (public_key, secret_key) = Kem::keygen();
    println!("   Public Key Size: {} bytes", public_key.to_bytes().len());
    println!("   Secret Key Size: {} bytes", secret_key.to_bytes().len());
    println!("   Public Key (first 32 bytes): {:02x?}", &public_key.to_bytes()[..32]);
    println!("   Secret Key (first 32 bytes): {:02x?}", &secret_key.to_bytes()[..32]);
    println!();

    // Example 2: Encapsulation
    println!("2. KEM Encapsulation:");
    let (ciphertext, shared_secret1) = Kem::encapsulate(&public_key);
    println!("   Ciphertext Size: {} bytes", ciphertext.to_bytes().len());
    println!("   Shared Secret Size: {} bytes", shared_secret1.to_bytes().len());
    println!("   Ciphertext (first 32 bytes): {:02x?}", &ciphertext.to_bytes()[..32]);
    println!("   Shared Secret: {}", shared_secret1.to_hex());
    println!();

    // Example 3: Decapsulation
    println!("3. KEM Decapsulation:");
    let shared_secret2 = Kem::decapsulate(&secret_key, &ciphertext);
    println!("   Decapsulated Secret: {}", shared_secret2.to_hex());
    println!("   Secrets Match: {}", shared_secret1.to_hex() == shared_secret2.to_hex());
    println!();

    // Example 4: Multiple encapsulations (should produce different ciphertexts)
    println!("4. Multiple Encapsulations:");
    let (ciphertext_a, secret_a) = Kem::encapsulate(&public_key);
    let (ciphertext_b, secret_b) = Kem::encapsulate(&public_key);
    
    println!("   Ciphertext A (first 16 bytes): {:02x?}", &ciphertext_a.to_bytes()[..16]);
    println!("   Ciphertext B (first 16 bytes): {:02x?}", &ciphertext_b.to_bytes()[..16]);
    println!("   Ciphertexts Different: {}", ciphertext_a.to_bytes() != ciphertext_b.to_bytes());
    println!("   Secret A: {}...", &secret_a.to_hex()[..32]);
    println!("   Secret B: {}...", &secret_b.to_hex()[..32]);
    println!();

    // Example 5: Serialization and deserialization
    println!("5. Serialization Example:");
    let public_hex = public_key.to_hex();
    let secret_hex = secret_key.to_hex();
    let ciphertext_hex = ciphertext.to_hex();
    let shared_secret_hex = shared_secret1.to_hex();
    
    println!("   Public Key Hex Length: {} chars", public_hex.len());
    println!("   Secret Key Hex Length: {} chars", secret_hex.len());
    println!("   Ciphertext Hex Length: {} chars", ciphertext_hex.len());
    println!("   Shared Secret Hex Length: {} chars", shared_secret_hex.len());
    println!();

    // Example 6: Deserialization
    println!("6. Deserialization Example:");
    let public_key_restored = topayz512::kem::PublicKey::from_hex(&public_hex).unwrap();
    let secret_key_restored = topayz512::kem::SecretKey::from_hex(&secret_hex).unwrap();
    let ciphertext_restored = topayz512::kem::Ciphertext::from_hex(&ciphertext_hex).unwrap();
    let shared_secret_restored = topayz512::kem::SharedSecret::from_hex(&shared_secret_hex).unwrap();
    
    println!("   Public Key Match: {}", public_key.to_bytes() == public_key_restored.to_bytes());
    println!("   Secret Key Match: {}", secret_key.to_bytes() == secret_key_restored.to_bytes());
    println!("   Ciphertext Match: {}", ciphertext.to_bytes() == ciphertext_restored.to_bytes());
    println!("   Shared Secret Match: {}", shared_secret1.to_bytes() == shared_secret_restored.to_bytes());
    println!();

    // Example 7: Performance simulation
    println!("7. Performance Simulation:");
    let start = std::time::Instant::now();
    
    for i in 0..10 {
        let (pk, sk) = Kem::keygen();
        let (ct, _ss1) = Kem::encapsulate(&pk);
        let _ss2 = Kem::decapsulate(&sk, &ct);
        
        if i == 0 {
            println!("   First operation completed successfully");
        }
    }
    
    let duration = start.elapsed();
    println!("   10 complete KEM operations took: {:?}", duration);
    println!("   Average per operation: {:?}", duration / 10);
    println!();

    println!("=== KEM Example Complete ===");
    println!("Note: This is a reference implementation for testing and development.");
    println!("For production use, implement proper LWE-based KEM with security analysis.");
}