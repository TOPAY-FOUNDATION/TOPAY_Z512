//! TOPAY-Z512 Hash Example
//!
//! This example demonstrates how to use the TOPAY-Z512 hash functionality.
//! It uses the current time as input to show dynamic hashing.

use std::time::{SystemTime, UNIX_EPOCH};
use topayz512::{Hash, hash, hash_combine};

fn main() {
    println!("TOPAY-Z512 Hash Example with Time-based Input\n");
    
    // Get current time as input
    let now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap();
    let time_str = format!("TOPAY-Z512 time: {} seconds", now.as_secs());
    let data = time_str.as_bytes();
    let hash_value = Hash::new(data);
    
    println!("Input: {}", String::from_utf8_lossy(data));
    println!("Hash: {}", hash_value);
    println!("Hash (hex): {}", hash_value.to_hex());
    println!("Hash size: {} bytes\n", hash_value.as_bytes().len());
    
    // Hash combination with time components
    let time_millis = now.as_millis().to_string();
    let data1 = format!("TOPAY-{}", time_millis).as_bytes().to_vec();
    let data2 = format!("Z512-{}", now.as_nanos() % 1000000).as_bytes().to_vec();
    
    let combined_hash = Hash::combine(&data1, &data2);
    
    println!("Input 1: {}", String::from_utf8_lossy(&data1));
    println!("Input 2: {}", String::from_utf8_lossy(&data2));
    println!("Combined Hash: {}", combined_hash);
    
    // Concatenated hash (different from combined hash)
    let mut concatenated = Vec::new();
    concatenated.extend_from_slice(&data1);
    concatenated.extend_from_slice(&data2);
    let concat_hash = Hash::new(&concatenated);
    
    println!("Concatenated Hash: {}", concat_hash);
    println!("Are they equal? {}\n", combined_hash == concat_hash);
    
    // Convenience functions
    let hash_bytes = hash(data);
    println!("Hash bytes (first 8): {:?}", &hash_bytes[0..8]);
    
    let combined_bytes = hash_combine(&data1, &data2);
    println!("Combined hash bytes (first 8): {:?}", &combined_bytes[0..8]);
    
    // Hex conversion
    let hex = hash_value.to_hex();
    println!("\nHex string: {}", hex);
    
    match Hash::from_hex(&hex) {
        Ok(hash_from_hex) => {
            println!("Converted back from hex: {}", hash_from_hex);
            println!("Equal to original: {}", hash_value == hash_from_hex);
        },
        Err(e) => println!("Error: {}", e),
    }
}