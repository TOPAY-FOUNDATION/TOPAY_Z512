//! Hash example for TOPAY-Z512
//!
//! This example demonstrates the hash functionality of the TOPAY-Z512 library.

use topayz512::{Hash, Result};

fn main() -> Result<()> {
    println!("=== TOPAY-Z512 Hash Example ===\n");

    // Example 1: Basic hash operations
    println!("1. Basic Hash Operations:");
    let input1 = b"Hello, TOPAY Foundation!";
    let hash1 = Hash::new(input1);
    println!("   Input: {}", std::str::from_utf8(input1).unwrap());
    println!("   Hash:  {}", hash1.to_hex());
    println!();

    // Example 2: Hash different inputs
    println!("2. Different Inputs:");
    let inputs = [
        b"" as &[u8],
        b"a",
        b"TOPAY-Z512",
        b"Quantum-safe blockchain for the future",
    ];

    for (i, input) in inputs.iter().enumerate() {
        let hash = Hash::new(input);
        println!(
            "   Input {}: {:?}",
            i + 1,
            std::str::from_utf8(input).unwrap_or("<binary>")
        );
        println!("   Hash:    {}", hash.to_hex());
    }
    println!();

    // Example 3: Hash combination
    println!("3. Hash Combination:");
    let data1 = b"first part";
    let data2 = b"second part";
    let combined_hash = Hash::combine(data1, data2);
    println!("   Data 1: {}", std::str::from_utf8(data1).unwrap());
    println!("   Data 2: {}", std::str::from_utf8(data2).unwrap());
    println!("   Combined Hash: {}", combined_hash.to_hex());
    println!();

    // Example 4: Hash concatenation
    println!("4. Hash Concatenation:");
    let hash_a = Hash::new(b"apple");
    let hash_b = Hash::new(b"banana");
    let hash_c = Hash::new(b"cherry");

    println!("   Hash A: {}", hash_a.to_hex());
    println!("   Hash B: {}", hash_b.to_hex());
    println!("   Hash C: {}", hash_c.to_hex());

    let concatenated = Hash::concat(&[&hash_a, &hash_b, &hash_c]);
    println!("   Concatenated: {}", concatenated.to_hex());
    println!();

    // Example 5: Hex conversion
    println!("5. Hex Conversion:");
    let original_hash = Hash::new(b"test data for hex conversion");
    let hex_string = original_hash.to_hex();
    let restored_hash = Hash::from_hex(&hex_string)?;

    println!("   Original:  {}", original_hash.to_hex());
    println!("   Hex:       {hex_string}");
    println!("   Restored:  {}", restored_hash.to_hex());
    println!("   Match:     {}", original_hash == restored_hash);
    println!();

    // Example 6: Binary data hashing
    println!("6. Binary Data Hashing:");
    let binary_data = vec![0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC];
    let binary_hash = Hash::hash_binary(&binary_data);
    println!("   Binary Data: {binary_data:02x?}");
    println!("   Hash:        {}", binary_hash.to_hex());
    println!();

    // Example 7: String hashing
    println!("7. String Hashing:");
    let message = "TOPAY Foundation: Building the future of quantum-safe finance";
    let string_hash = Hash::hash_string(message);
    println!("   Message: {message}");
    println!("   Hash:    {}", string_hash.to_hex());
    println!();

    println!("=== Hash Example Complete ===");
    Ok(())
}
