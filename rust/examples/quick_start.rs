// Quick Start Example for TOPAY-Z512
// This shows the most common use cases in a simple format

use topayz512::{Hash, Kem, KeyPair};

fn main() {
    println!("🚀 TOPAY-Z512 Quick Start Guide");
    println!("================================\n");

    // 1. Generate a key pair
    println!("1. Key Pair Generation:");
    let keypair = KeyPair::generate();
    println!("   ✅ Generated new key pair");
    println!(
        "   Private key (first 16 chars): {}",
        &keypair.private_key().to_hex()[..16]
    );
    println!(
        "   Public key (first 16 chars):  {}\n",
        &keypair.public_key().to_hex()[..16]
    );

    // 2. Hash some data
    println!("2. Hashing:");
    let message = b"Hello, TOPAY-Z512!";
    let hash = Hash::new(message);
    println!("   Message: {:?}", std::str::from_utf8(message).unwrap());
    println!("   Hash (first 32 chars): {}\n", &hash.to_hex()[..32]);

    // 3. Key Encapsulation Mechanism
    println!("3. Key Encapsulation:");
    let (kem_public, kem_secret) = Kem::keygen();
    let (ciphertext, shared_secret1) = Kem::encapsulate(&kem_public);
    let shared_secret2 = Kem::decapsulate(&kem_secret, &ciphertext);

    println!("   ✅ Generated KEM keys");
    println!("   ✅ Encapsulated shared secret");
    println!("   ✅ Decapsulated shared secret");
    println!(
        "   Secrets match: {}\n",
        shared_secret1.to_bytes() == shared_secret2.to_bytes()
    );

    // 4. Serialization
    println!("4. Serialization:");
    let (private_hex, public_hex) = keypair.to_hex();
    println!("   ✅ Serialized keys to hex");
    println!("   Private hex length: {} chars", private_hex.len());
    println!("   Public hex length:  {} chars\n", public_hex.len());

    // 5. Deserialization
    println!("5. Deserialization:");
    let restored_private = topayz512::keypair::PrivateKey::from_hex(&private_hex).unwrap();
    let restored_public = topayz512::keypair::PublicKey::from_hex(&public_hex).unwrap();
    println!("   ✅ Restored keys from hex");
    println!(
        "   Keys match: {}\n",
        *keypair.private_key() == restored_private && *keypair.public_key() == restored_public
    );

    println!("🎉 Quick start complete! You're ready to use TOPAY-Z512.");
    println!("\nNext steps:");
    println!(
        "- Run other examples: cargo run --example [keypair_example|kem_example|hash_example]"
    );
    println!("- Run tests: cargo test");
    println!("- Build for production: cargo build --release");
}
