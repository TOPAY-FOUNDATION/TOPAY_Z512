// KEM (Key Encapsulation Mechanism) Explained with Examples
// This demonstrates what KEM is and why it's important

use topayz512::Kem;

fn main() {
    println!("🔐 Understanding KEM (Key Encapsulation Mechanism)");
    println!("==================================================\n");

    println!("📚 WHAT IS KEM?");
    println!("---------------");
    println!("KEM is a cryptographic method that allows two parties to securely");
    println!("establish a shared secret key over an insecure communication channel.\n");

    println!("🎯 WHY IS KEM IMPORTANT?");
    println!("------------------------");
    println!("1. 🛡️  Quantum-Safe: Protects against future quantum computers");
    println!("2. 🔒 Secure: Even if someone intercepts the communication, they can't get the secret");
    println!("3. 🚀 Efficient: Fast key exchange for secure communication");
    println!("4. 🌐 Internet-Ready: Perfect for secure web communications\n");

    println!("🔄 HOW KEM WORKS (Step by Step):");
    println!("---------------------------------");

    // Step 1: Generate keys
    println!("Step 1: Bob generates a key pair");
    let (bob_public_key, bob_secret_key) = Kem::keygen();
    println!("   ✅ Bob has: Public Key (shareable) + Secret Key (private)");
    println!("   📤 Bob shares his public key with Alice\n");

    // Step 2: Alice encapsulates
    println!("Step 2: Alice wants to send a secure message to Bob");
    println!("   🔐 Alice uses Bob's public key to 'encapsulate' a shared secret");
    let (ciphertext, alice_shared_secret) = Kem::encapsulate(&bob_public_key);
    println!("   ✅ Alice gets: Ciphertext + Shared Secret");
    println!("   📤 Alice sends the ciphertext to Bob (safe to intercept!)\n");

    // Step 3: Bob decapsulates
    println!("Step 3: Bob receives the ciphertext from Alice");
    println!("   🔓 Bob uses his secret key to 'decapsulate' the shared secret");
    let bob_shared_secret = Kem::decapsulate(&bob_secret_key, &ciphertext);
    println!("   ✅ Bob gets the same shared secret as Alice!\n");

    // Step 4: Verification
    println!("Step 4: Verification");
    let secrets_match = alice_shared_secret.to_bytes() == bob_shared_secret.to_bytes();
    println!("   🤝 Alice and Bob have the same secret: {}", secrets_match);
    println!("   🔒 They can now use this secret for secure communication!\n");

    println!("📊 TECHNICAL DETAILS:");
    println!("---------------------");
    println!("   Public Key Size:  {} bytes", bob_public_key.to_bytes().len());
    println!("   Secret Key Size:  {} bytes", bob_secret_key.to_bytes().len());
    println!("   Ciphertext Size:  {} bytes", ciphertext.to_bytes().len());
    println!("   Shared Secret:    {} bytes\n", alice_shared_secret.to_bytes().len());

    println!("🌟 REAL-WORLD APPLICATIONS:");
    println!("---------------------------");
    println!("   🌐 HTTPS websites (secure web browsing)");
    println!("   💬 Secure messaging apps");
    println!("   🏦 Banking and financial transactions");
    println!("   📱 Mobile app security");
    println!("   🔐 VPN connections");
    println!("   ⚡ Cryptocurrency transactions\n");

    println!("🚀 WHY TOPAY-Z512 KEM IS SPECIAL:");
    println!("----------------------------------");
    println!("   🛡️  Post-Quantum Ready: Safe against quantum computer attacks");
    println!("   ⚡ Fast Performance: ~30µs per operation");
    println!("   🔒 512-bit Security: Maximum security level");
    println!("   🌐 Cross-Platform: Works on any system");
    println!("   💾 Serializable: Easy to store and transmit\n");

    println!("🎯 SIMPLE ANALOGY:");
    println!("------------------");
    println!("   Think of KEM like a magic lockbox:");
    println!("   1. 📦 Bob gives Alice a special lockbox (public key)");
    println!("   2. 🔐 Alice puts a secret inside and locks it (encapsulation)");
    println!("   3. 📤 Alice sends the locked box to Bob (ciphertext)");
    println!("   4. 🔓 Bob unlocks it with his key (decapsulation)");
    println!("   5. 🤝 Both now have the same secret!\n");

    println!("✨ You now understand KEM! It's the foundation of modern secure communication.");
}