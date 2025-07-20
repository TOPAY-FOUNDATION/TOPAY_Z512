//! TOPAY-Z512 Private to Public Key Conversion Example
//!
//! This example demonstrates how to convert a private key to a public key
//! using the TOPAY-Z512 library.

use rand::rngs::OsRng;
use topayz512::{PrivateKey, private_to_public};

fn main() {
    println!("TOPAY-Z512 Private to Public Key Conversion Example\n");
    
    // Generate a random private key
    let mut rng = OsRng;
    let private_key = PrivateKey::generate(&mut rng);
    println!("Generated private key: {}", private_key.to_hex());
    
    // Convert private key to public key using the convenience function
    let public_key = private_to_public(&private_key);
    println!("Derived public key: {}", public_key.to_hex());
    
    // Create a private key from a hex string
    let hex_private_key = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef";
    let private_key = PrivateKey::from_hex(hex_private_key).unwrap();
    println!("\nPredefined private key: {}", private_key.to_hex());
    
    // Convert private key to public key
    let public_key = private_to_public(&private_key);
    println!("Derived public key: {}", public_key.to_hex());
}