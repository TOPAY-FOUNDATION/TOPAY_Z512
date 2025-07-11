//! Example demonstrating key exchange using topayz512.

use topayz512::{decapsulate, encapsulate, keygen};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("topayz512 Key Exchange Example");
    println!("=================================");

    // Alice generates a key pair
    println!("Alice: Generating key pair...");
    let (alice_pk, alice_sk) = keygen()?;
    println!("Alice: Key pair generated successfully.");

    // Bob encapsulates a shared secret using Alice's public key
    println!("\nBob: Encapsulating shared secret using Alice's public key...");
    let (ciphertext, bob_shared_secret) = encapsulate(&alice_pk)?;
    println!("Bob: Shared secret encapsulated successfully.");
    println!("Bob's shared secret (first 16 bytes): {:?}", &bob_shared_secret[0..16]);

    // Alice decapsulates the shared secret using her secret key and Bob's ciphertext
    println!("\nAlice: Decapsulating shared secret using Bob's ciphertext...");
    let alice_shared_secret = decapsulate(&alice_sk, &ciphertext)?;
    println!("Alice: Shared secret decapsulated successfully.");
    println!("Alice's shared secret (first 16 bytes): {:?}", &alice_shared_secret[0..16]);

    // Verify that both parties have the same shared secret
    if alice_shared_secret == bob_shared_secret {
        println!("\nSuccess! Both parties have the same shared secret.");
    } else {
        println!("\nError! The shared secrets do not match.");
    }

    Ok(())
}