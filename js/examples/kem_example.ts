/**
 * TOPAY-Z512 KEM Example
 * 
 * This example demonstrates the usage of the TOPAY-Z512 Key Encapsulation Mechanism (KEM).
 */

import { keygen, encapsulate, decapsulate, Ciphertext, SharedSecret } from '../src/kem';

// Main example function
function kemExample() {
  console.log('TOPAY-Z512 KEM Example');
  console.log('=======================');
  
  // Generate a key pair
  const keypair = keygen();
  
  console.log('Generated key pair:');
  console.log(`  Private key: ${keypair.privateKey.toHex().substring(0, 16)}...`);
  console.log(`  Public key: ${keypair.publicKey.toHex().substring(0, 16)}...`);
  
  // Encapsulate a shared secret
  const { ciphertext, sharedSecret: sharedSecret1 } = encapsulate(keypair.publicKey);
  
  console.log('\nEncapsulated shared secret:');
  console.log(`  Ciphertext: ${ciphertext.toHex().substring(0, 32)}...`);
  console.log(`  Shared secret: ${sharedSecret1.toHex().substring(0, 32)}...`);
  
  // Decapsulate the shared secret
  const sharedSecret2 = decapsulate(keypair.privateKey, ciphertext);
  
  console.log('\nDecapsulated shared secret:');
  console.log(`  Shared secret: ${sharedSecret2.toHex().substring(0, 32)}...`);
  
  // Verify that the shared secrets match
  if (sharedSecret1.toHex() === sharedSecret2.toHex()) {
    console.log('\nSuccess! The shared secrets match.');
  } else {
    console.log('\nError! The shared secrets do not match.');
  }
  
  // Demonstrate serialization
  console.log('\nSerialization Example:');
  
  // Convert ciphertext to hex and back
  const ciphertextHex = ciphertext.toHex();
  const ciphertext2 = Ciphertext.fromHex(ciphertextHex);
  
  console.log(`  Original ciphertext: ${ciphertext.toHex().substring(0, 32)}...`);
  console.log(`  Deserialized ciphertext: ${ciphertext2.toHex().substring(0, 32)}...`);
  
  // Convert shared secret to hex and back
  const sharedSecretHex = sharedSecret1.toHex();
  const sharedSecret3 = SharedSecret.fromHex(sharedSecretHex);
  
  console.log(`  Original shared secret: ${sharedSecret1.toHex().substring(0, 32)}...`);
  console.log(`  Deserialized shared secret: ${sharedSecret3.toHex().substring(0, 32)}...`);
}

// Run the example
kemExample();