//! TOPAY-Z512 KEM Example
//!
//! This example demonstrates the usage of the TOPAY-Z512 Key Encapsulation Mechanism (KEM).

use rand::rngs::OsRng;
use topayz512::{keygen, encapsulate, decapsulate};

fn main() {
    println!("TOPAY-Z512 KEM Example");
    println!("=======================");
    
    // Generate a key pair
    let mut rng = OsRng;
    let keypair = keygen(&mut rng);
    
    println!("Generated key pair:");
    println!("  Private key: {}...", &keypair.private_key.to_hex()[..16]);
    println!("  Public key: {}...", &keypair.public_key.to_hex()[..16]);
    
    // Encapsulate a shared secret
    let (ciphertext, shared_secret1) = encapsulate(&keypair.public_key, &mut rng);
    
    println!("\nEncapsulated shared secret:");
    println!("  Ciphertext: {}...", &ciphertext.to_hex()[..32]);
    println!("  Shared secret: {}...", &shared_secret1.to_hex()[..32]);
    
    // Decapsulate the shared secret
    let shared_secret2 = decapsulate(&keypair.private_key, &ciphertext);
    
    println!("\nDecapsulated shared secret:");
    println!("  Shared secret: {}...", &shared_secret2.to_hex()[..32]);
    
    // Verify that the shared secrets match
    if shared_secret1.to_hex() == shared_secret2.to_hex() {
        println!("\nSuccess! The shared secrets match.");
    } else {
        println!("\nError! The shared secrets do not match.");
    }
    
    // Demonstrate serialization
    println!("\nSerialization Example:");
    
    // Convert ciphertext to hex and back
    let ciphertext_hex = ciphertext.to_hex();
    let ciphertext2 = topayz512::Ciphertext::from_hex(&ciphertext_hex).unwrap();
    
    println!("  Original ciphertext: {}...", &ciphertext.to_hex()[..32]);
    println!("  Deserialized ciphertext: {}...", &ciphertext2.to_hex()[..32]);
    
    // Convert shared secret to hex and back
    let shared_secret_hex = shared_secret1.to_hex();
    let shared_secret3 = topayz512::SharedSecret::from_hex(&shared_secret_hex).unwrap();
    
    println!("  Original shared secret: {}...", &shared_secret1.to_hex()[..32]);
    println!("  Deserialized shared secret: {}...", &shared_secret3.to_hex()[..32]);
}