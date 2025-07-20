/**
 * TOPAY-Z512 Private to Public Key Conversion Example
 *
 * This example demonstrates how to convert a private key to a public key
 * using the TOPAY-Z512 library.
 */

import { PrivateKey, privateToPublic } from '../src/keypair';

// Main function
async function main() {
  console.log('TOPAY-Z512 Private to Public Key Conversion Example\n');

  // Generate a random private key
  const privateKey = await PrivateKey.generate();
  console.log(`Generated private key: ${privateKey.toHex()}`);

  // Convert private key to public key using the convenience function
  const publicKey = privateToPublic(privateKey);
  console.log(`Derived public key: ${publicKey.toHex()}`);

  // Create a private key from a hex string
  const hexPrivateKey = '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef';
  const predefinedPrivateKey = PrivateKey.fromHex(hexPrivateKey);
  console.log(`\nPredefined private key: ${predefinedPrivateKey.toHex()}`);

  // Convert private key to public key
  const derivedPublicKey = privateToPublic(predefinedPrivateKey);
  console.log(`Derived public key: ${derivedPublicKey.toHex()}`);
}

// Run the main function
main().catch(error => {
  console.error('Error:', error);
});