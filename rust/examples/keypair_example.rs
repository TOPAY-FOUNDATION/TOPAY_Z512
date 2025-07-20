//! TOPAY-Z512 Key Pair Example
//!
//! This example demonstrates how to generate and use key pairs in TOPAY-Z512.
//! It shows key generation, public key derivation, and hex conversion.

use rand::rngs::OsRng;
use topayz512::{KeyPair, PrivateKey, PublicKey, generate_keypair};

fn main() {
    println!("TOPAY-Z512 Key Pair Example\n");
    
    // Generate a new key pair using the OS random number generator
    let mut rng = OsRng;
    let keypair = generate_keypair(&mut rng);
    
    println!("Generated Key Pair:");
    println!("Private Key: {}", keypair.private_key.to_hex());
    println!("Public Key:  {}\n", keypair.public_key.to_hex());
    
    // Demonstrate deriving a public key from a private key
    let derived_public_key = PublicKey::from_private_key(&keypair.private_key);
    println!("Derived Public Key: {}", derived_public_key.to_hex());
    println!("Keys match: {}\n", derived_public_key.to_hex() == keypair.public_key.to_hex());
    
    // Demonstrate hex conversion
    let private_hex = keypair.private_key.to_hex();
    let recovered_private_key = PrivateKey::from_hex(&private_hex).unwrap();
    
    let public_hex = keypair.public_key.to_hex();
    let recovered_public_key = PublicKey::from_hex(&public_hex).unwrap();
    
    println!("Hex Conversion Test:");
    println!("Private Key Recovered: {}", 
             private_hex == PrivateKey::from_hex(&private_hex).unwrap().to_hex());
    println!("Public Key Recovered:  {}", 
             public_hex == PublicKey::from_hex(&public_hex).unwrap().to_hex());
    
    // Create a new key pair from existing keys
    let new_keypair = KeyPair::from_keys(recovered_private_key, recovered_public_key);
    println!("\nRecreated Key Pair:");
    println!("Private Key: {}", new_keypair.private_key.to_hex());
    println!("Public Key:  {}", new_keypair.public_key.to_hex());
}