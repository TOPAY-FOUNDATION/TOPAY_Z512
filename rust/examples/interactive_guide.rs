// Interactive TOPAY-Z512 Usage Example
// This shows you exactly how to use each component!

use topayz512::{Hash, Kem, KeyPair};

fn main() {
    println!("🎯 TOPAY-Z512 Interactive Usage Guide");
    println!("=====================================\n");

    // 🔑 1. KEY PAIR USAGE
    println!("🔑 1. KEY PAIR OPERATIONS");
    println!("-------------------------");

    // Generate a new key pair
    let my_keypair = KeyPair::generate();
    println!("✅ Generated new key pair!");

    // Get the keys
    let private_key = my_keypair.private_key();
    let public_key = my_keypair.public_key();

    println!("📋 Private Key: {}", private_key.to_hex());
    println!("📋 Public Key:  {}", public_key.to_hex());

    // Save keys to hex strings (for storage/transmission)
    let (private_hex, public_hex) = my_keypair.to_hex();
    println!("💾 Keys saved as hex strings\n");

    // 🔐 2. KEY ENCAPSULATION MECHANISM (KEM)
    println!("🔐 2. KEM OPERATIONS (Secure Key Exchange)");
    println!("-------------------------------------------");

    // Generate KEM keys
    let (kem_public, kem_secret) = Kem::keygen();
    println!("✅ Generated KEM key pair");

    // Alice encapsulates a shared secret using Bob's public key
    let (ciphertext, alice_secret) = Kem::encapsulate(&kem_public);
    println!("📦 Alice encapsulated shared secret");
    println!("   Ciphertext: {}", ciphertext.to_hex());

    // Bob decapsulates the shared secret using his private key
    let bob_secret = Kem::decapsulate(&kem_secret, &ciphertext);
    println!("🔓 Bob decapsulated shared secret");

    // Verify they have the same secret
    let secrets_match = alice_secret.to_bytes() == bob_secret.to_bytes();
    println!("🤝 Shared secrets match: {secrets_match}");
    println!("   Alice's secret: {}", alice_secret.to_hex());
    println!("   Bob's secret:   {}\n", bob_secret.to_hex());

    // 🔨 3. HASHING
    println!("🔨 3. HASH OPERATIONS");
    println!("---------------------");

    // Hash some data
    let message = b"Hello from TOPAY-Z512!";
    let hash1 = Hash::new(message);
    println!("📝 Message: {:?}", std::str::from_utf8(message).unwrap());
    println!("🔍 Hash: {}", hash1.to_hex());

    // Combine two pieces of data
    let data1 = b"First part";
    let data2 = b"Second part";
    let combined_hash = Hash::combine(data1, data2);
    println!(
        "🔗 Combined hash of two parts: {}\n",
        combined_hash.to_hex()
    );

    // 💾 4. SERIALIZATION & STORAGE
    println!("💾 4. SERIALIZATION (Save/Load)");
    println!("--------------------------------");

    // Save everything to hex strings
    let kem_public_hex = kem_public.to_hex();
    let kem_secret_hex = kem_secret.to_hex();
    let hash_hex = hash1.to_hex();

    println!("✅ All data serialized to hex strings");
    println!("📁 You can now save these to files or databases:");
    println!("   - Private key: {} chars", private_hex.len());
    println!("   - Public key:  {} chars", public_hex.len());
    println!("   - KEM public:  {} chars", kem_public_hex.len());
    println!("   - KEM secret:  {} chars", kem_secret_hex.len());
    println!("   - Hash:        {} chars\n", hash_hex.len());

    // 🔄 5. RESTORATION
    println!("🔄 5. RESTORATION (Load from hex)");
    println!("----------------------------------");

    // Restore from hex strings
    let restored_private = topayz512::keypair::PrivateKey::from_hex(&private_hex).unwrap();
    let restored_public = topayz512::keypair::PublicKey::from_hex(&public_hex).unwrap();
    let restored_hash = Hash::from_hex(&hash_hex).unwrap();

    println!("✅ Successfully restored all data from hex!");
    println!("🔍 Verification:");
    println!(
        "   Private key matches: {}",
        *private_key == restored_private
    );
    println!("   Public key matches:  {}", *public_key == restored_public);
    println!(
        "   Hash matches:        {}\n",
        hash1.to_bytes() == restored_hash.to_bytes()
    );

    // 🚀 6. READY TO USE!
    println!("🚀 YOU'RE READY TO USE TOPAY-Z512!");
    println!("===================================");
    println!("✨ You now know how to:");
    println!("   🔑 Generate and manage key pairs");
    println!("   🔐 Perform secure key exchange with KEM");
    println!("   🔨 Hash data securely");
    println!("   💾 Save and load all cryptographic data");
    println!("   🔄 Serialize for storage or transmission");
    println!("\n🎯 Next steps:");
    println!("   - Integrate into your application");
    println!("   - Run 'cargo test' to see all features");
    println!("   - Check other examples for advanced usage");
}
