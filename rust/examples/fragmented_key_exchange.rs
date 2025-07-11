//! Example demonstrating fragmented key exchange using topayz512.

use topayz512::fragmented::{decapsulate, encapsulate, keygen};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("topayz512 Fragmented Key Exchange Example");
    println!("==========================================");

    // Number of fragments to use
    let num_fragments = 4;
    println!("Using {} fragments for key exchange", num_fragments);

    // Alice generates a fragmented key pair
    println!("\nAlice: Generating fragmented key pair...");
    let (alice_pk, alice_sk) = keygen(num_fragments)?;
    println!("Alice: Fragmented key pair generated successfully.");

    // Bob encapsulates a shared secret using Alice's fragmented public key
    println!("\nBob: Encapsulating shared secret using Alice's fragmented public key...");
    let (ciphertext, bob_shared_secret) = encapsulate(&alice_pk)?;
    println!("Bob: Shared secret encapsulated successfully.");
    println!("Bob's shared secret (first 16 bytes): {:?}", &bob_shared_secret[0..16]);

    // Alice decapsulates the shared secret using her fragmented secret key and Bob's fragmented ciphertext
    println!("\nAlice: Decapsulating shared secret using Bob's fragmented ciphertext...");
    let alice_shared_secret = decapsulate(&alice_sk, &ciphertext)?;
    println!("Alice: Shared secret decapsulated successfully.");
    println!("Alice's shared secret (first 16 bytes): {:?}", &alice_shared_secret[0..16]);

    // Verify that both parties have the same shared secret
    if alice_shared_secret == bob_shared_secret {
        println!("\nSuccess! Both parties have the same shared secret.");
    } else {
        println!("\nError! The shared secrets do not match.");
    }

    // Compare performance with different numbers of fragments
    println!("\nPerformance comparison with different numbers of fragments:");
    println!("--------------------------------------------------------");
    
    use std::time::Instant;
    
    for &frags in &[2, 4, 8] {
        println!("\nTesting with {} fragments:", frags);
        
        // Measure key generation time
        let start = Instant::now();
        let (pk, sk) = keygen(frags)?;
        let keygen_time = start.elapsed();
        println!("  Key generation time: {:?}", keygen_time);
        
        // Measure encapsulation time
        let start = Instant::now();
        let (ct, _) = encapsulate(&pk)?;
        let encap_time = start.elapsed();
        println!("  Encapsulation time: {:?}", encap_time);
        
        // Measure decapsulation time
        let start = Instant::now();
        let _ = decapsulate(&sk, &ct)?;
        let decap_time = start.elapsed();
        println!("  Decapsulation time: {:?}", decap_time);
    }

    Ok(())
}