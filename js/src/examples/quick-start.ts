/**
 * Quick Start Example for TOPAY-Z512
 * 
 * This example demonstrates the basic usage of TOPAY-Z512 cryptographic library
 * including key generation, hashing, KEM operations, and data fragmentation.
 */

import {
  generateKeyPair,
  computeHash,
  kemKeyGen,
  kemEncapsulate,
  kemDecapsulate,
  fragmentData,
  reconstructData,
  toHex
} from '../index';

async function quickStartExample(): Promise<void> {
  console.log('🚀 TOPAY-Z512 Quick Start Example');
  console.log('='.repeat(50));

  try {
    // 1. Generate a key pair
    console.log('\n1. Generating key pair...');
    const keyPair = await generateKeyPair();
    console.log(`   Private key: ${toHex(keyPair.privateKey).substring(0, 32)}...`);
    console.log(`   Public key:  ${toHex(keyPair.publicKey).substring(0, 32)}...`);

    // 2. Hash some data
    console.log('\n2. Computing hash...');
    const testData = new TextEncoder().encode('Hello, TOPAY-Z512!');
    const hash = await computeHash(testData);
    console.log(`   Data: "Hello, TOPAY-Z512!"`);
    console.log(`   Hash: ${toHex(hash).substring(0, 32)}...`);

    // 3. KEM operations
    console.log('\n3. KEM operations...');
    const kemKeyPair = await kemKeyGen();
    console.log('   Generated KEM key pair');

    const { ciphertext, sharedSecret } = await kemEncapsulate(kemKeyPair.publicKey);
    console.log('   Encapsulated shared secret');
    console.log(`   Ciphertext: ${toHex(ciphertext).substring(0, 32)}...`);
    console.log(`   Shared secret: ${toHex(sharedSecret).substring(0, 32)}...`);

    const decapsulatedSecret = await kemDecapsulate(kemKeyPair.secretKey, ciphertext);
    console.log('   Decapsulated shared secret');
    console.log(`   Secrets match: ${toHex(sharedSecret) === toHex(decapsulatedSecret)}`);

    // 4. Data fragmentation
    console.log('\n4. Data fragmentation...');
    const largeData = new Uint8Array(1024);
    crypto.getRandomValues(largeData);
    
    const fragResult = await fragmentData(largeData);
    console.log(`   Original size: ${fragResult.metadata.originalSize} bytes`);
    console.log(`   Fragment count: ${fragResult.metadata.fragmentCount}`);
    console.log(`   Checksum: ${toHex(fragResult.metadata.checksum).substring(0, 32)}...`);

    const reconResult = await reconstructData(fragResult.fragments);
    console.log(`   Reconstruction complete: ${reconResult.isComplete}`);
    console.log(`   Data integrity: ${largeData.length === reconResult.data.length}`);

    console.log('\n✅ Quick start example completed successfully!');
    console.log('\nNext steps:');
    console.log('- Explore other examples in the examples/ directory');
    console.log('- Read the documentation for advanced features');
    console.log('- Check out the benchmark example for performance testing');

  } catch (error) {
    console.error('❌ Error in quick start example:', error);
    process.exit(1);
  }
}

// Run the example
if (require.main === module) {
  quickStartExample().catch(console.error);
}

export { quickStartExample };