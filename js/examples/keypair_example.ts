/**
 * TOPAY-Z512 Key Pair Example
 * 
 * This example demonstrates how to generate and use key pairs in TOPAY-Z512.
 * It shows key generation, public key derivation, and hex conversion.
 */

import { KeyPair, PrivateKey, PublicKey, generateKeyPair } from '../src/keypair';

console.log('TOPAY-Z512 Key Pair Example\n');

// Generate a new key pair
const keypair = generateKeyPair();

console.log('Generated Key Pair:');
console.log(`Private Key: ${keypair.privateKey.toHex()}`);
console.log(`Public Key:  ${keypair.publicKey.toHex()}\n`);

// Demonstrate deriving a public key from a private key
const derivedPublicKey = PublicKey.fromPrivateKey(keypair.privateKey);
console.log(`Derived Public Key: ${derivedPublicKey.toHex()}`);
console.log(`Keys match: ${derivedPublicKey.toHex() === keypair.publicKey.toHex()}\n`);

// Demonstrate hex conversion
const privateHex = keypair.privateKey.toHex();
const recoveredPrivateKey = PrivateKey.fromHex(privateHex);

const publicHex = keypair.publicKey.toHex();
const recoveredPublicKey = PublicKey.fromHex(publicHex);

console.log('Hex Conversion Test:');
console.log(`Private Key Recovered: ${privateHex === recoveredPrivateKey.toHex()}`);
console.log(`Public Key Recovered:  ${publicHex === recoveredPublicKey.toHex()}`);

// Create a new key pair from existing keys
const newKeypair = new KeyPair(recoveredPrivateKey, recoveredPublicKey);
console.log('\nRecreated Key Pair:');
console.log(`Private Key: ${newKeypair.privateKey.toHex()}`);
console.log(`Public Key:  ${newKeypair.publicKey.toHex()}`);