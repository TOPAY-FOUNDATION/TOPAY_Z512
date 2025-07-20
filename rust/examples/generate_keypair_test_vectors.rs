//! TOPAY-Z512 Key Pair Test Vector Generator
//!
//! This example generates actual public key values for the test vectors in the test-vectors directory.

use std::fs::File;
use std::io::{Read, Write};
use std::path::Path;

use serde_json::{json, Value};
use topayz512::{PrivateKey, PublicKey};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("TOPAY-Z512 Key Pair Test Vector Generator\n");
    
    // Load the test vector file
    let test_vector_path = Path::new("../test-vectors/003_keypair.json");
    let mut file = File::open(&test_vector_path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;
    
    // Parse the JSON
    let mut json_data: Value = serde_json::from_str(&contents)?;
    
    // Process the key pair test vectors
    if let Some(test_vectors) = json_data["keypair_test_vectors"].as_array_mut() {
        for test_vector in test_vectors {
            if let (Some(private_key_hex), Some(test_id)) = (
                test_vector["private_key"].as_str(),
                test_vector["test_id"].as_str(),
            ) {
                println!("Processing test vector: {}", test_id);
                
                // Parse the private key
                let private_key = PrivateKey::from_hex(private_key_hex)?;
                
                // Derive the public key
                let public_key = PublicKey::from_private_key(&private_key);
                
                // Update the expected public key
                test_vector["expected_public_key"] = json!(public_key.to_hex());
            }
        }
    }
    
    // Update the notes
    if let Some(notes) = json_data["notes"].as_array_mut() {
        if let Some(note) = notes.get_mut(0) {
            *note = json!("The expected_public_key values are actual SHA3-512 hash outputs of the private keys.");
        }
    }
    
    // Write the updated test vectors back to the file
    let mut file = File::create(&test_vector_path)?;
    let formatted_json = serde_json::to_string_pretty(&json_data)?;
    file.write_all(formatted_json.as_bytes())?;
    
    println!("\nTest vectors updated successfully!");
    
    Ok(())
}