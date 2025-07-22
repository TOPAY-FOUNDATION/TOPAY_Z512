/**
 * Key Pair Management Example for TOPAY-Z512
 * 
 * This example demonstrates advanced key pair operations including:
 * - Basic key generation
 * - Deterministic key generation from seeds
 * - HD wallet generation
 * - Key derivation from passwords
 * - Key validation and serialization
 */

import {
  generateKeyPair,
  generateKeyPairFromSeed,
  batchGenerateKeyPairs,
  validateKeyPair,
  deriveChildKeyPair,
  generateHDWallet,
  deriveKeyPairFromPassword,
  serializeKeyPair,
  deserializeKeyPair,
  secureEraseKeyPair,
  backupKeyPair,
  toHex,
  fromHex,
  secureRandom
} from '../index';

async function keyPairExample(): Promise<void> {
  console.log('🔑 TOPAY-Z512 Key Pair Management Example');
  console.log('='.repeat(50));

  try {
    // 1. Basic key pair generation
    console.log('\n1. Basic Key Pair Generation');
    console.log('-'.repeat(30));
    
    const keyPair = await generateKeyPair();
    console.log(`   Private key: ${toHex(keyPair.privateKey).substring(0, 32)}...`);
    console.log(`   Public key:  ${toHex(keyPair.publicKey).substring(0, 32)}...`);
    
    const isValid = await validateKeyPair(keyPair);
    console.log(`   Key pair valid: ${isValid}`);

    // 2. Deterministic key generation
    console.log('\n2. Deterministic Key Generation');
    console.log('-'.repeat(30));
    
    const seed = await secureRandom(32);
    console.log(`   Seed: ${toHex(seed).substring(0, 32)}...`);
    
    const deterministicKeyPair1 = await generateKeyPairFromSeed(seed);
    const deterministicKeyPair2 = await generateKeyPairFromSeed(seed);
    
    const keysMatch = toHex(deterministicKeyPair1.privateKey) === toHex(deterministicKeyPair2.privateKey);
    console.log(`   Deterministic generation: ${keysMatch ? 'SUCCESS' : 'FAILED'}`);

    // 3. Batch key generation
    console.log('\n3. Batch Key Generation');
    console.log('-'.repeat(30));
    
    const startTime = Date.now();
    const keyPairs = await batchGenerateKeyPairs(10);
    const endTime = Date.now();
    
    console.log(`   Generated ${keyPairs.length} key pairs in ${endTime - startTime}ms`);
    console.log(`   Average time per key pair: ${(endTime - startTime) / keyPairs.length}ms`);

    // 4. Child key derivation
    console.log('\n4. Child Key Derivation');
    console.log('-'.repeat(30));
    
    const parentKeyPair = await generateKeyPair();
    console.log(`   Parent public key: ${toHex(parentKeyPair.publicKey).substring(0, 32)}...`);
    
    const childKeyPair1 = await deriveChildKeyPair(parentKeyPair.privateKey, 0);
    const childKeyPair2 = await deriveChildKeyPair(parentKeyPair.privateKey, 1);
    
    console.log(`   Child 0 public key: ${toHex(childKeyPair1.publicKey).substring(0, 32)}...`);
    console.log(`   Child 1 public key: ${toHex(childKeyPair2.publicKey).substring(0, 32)}...`);
    
    const child1Valid = await validateKeyPair(childKeyPair1);
    const child2Valid = await validateKeyPair(childKeyPair2);
    console.log(`   Child key pairs valid: ${child1Valid && child2Valid}`);

    // 5. HD Wallet generation
    console.log('\n5. HD Wallet Generation');
    console.log('-'.repeat(30));
    
    const masterSeed = await secureRandom(64);
    const hdWallet = await generateHDWallet(masterSeed, 5);
    
    console.log(`   Generated HD wallet with ${hdWallet.length} key pairs`);
    for (let i = 0; i < hdWallet.length; i++) {
      console.log(`   Key ${i}: ${toHex(hdWallet[i]!.publicKey).substring(0, 32)}...`);
    }

    // 6. Password-based key derivation
    console.log('\n6. Password-based Key Derivation');
    console.log('-'.repeat(30));
    
    const password = 'MySecurePassword123!';
    const salt = await secureRandom(32);
    
    const passwordKeyPair1 = await deriveKeyPairFromPassword(password, salt);
    const passwordKeyPair2 = await deriveKeyPairFromPassword(password, salt);
    
    const passwordKeysMatch = toHex(passwordKeyPair1.privateKey) === toHex(passwordKeyPair2.privateKey);
    console.log(`   Password-based derivation: ${passwordKeysMatch ? 'SUCCESS' : 'FAILED'}`);
    console.log(`   Derived public key: ${toHex(passwordKeyPair1.publicKey).substring(0, 32)}...`);

    // 7. Key serialization and deserialization
    console.log('\n7. Key Serialization');
    console.log('-'.repeat(30));
    
    const testKeyPair = await generateKeyPair();
    const serialized = serializeKeyPair(testKeyPair);
    console.log(`   Serialized length: ${serialized.length} characters`);
    
    const deserialized = deserializeKeyPair(serialized);
    const serializationValid = toHex(testKeyPair.privateKey) === toHex(deserialized.privateKey);
    console.log(`   Serialization round-trip: ${serializationValid ? 'SUCCESS' : 'FAILED'}`);

    // 8. Key backup and secure erasure
    console.log('\n8. Key Backup and Secure Erasure');
    console.log('-'.repeat(30));
    
    const originalKeyPair = await generateKeyPair();
    const originalPrivateKey = toHex(originalKeyPair.privateKey);
    
    const backup = backupKeyPair(originalKeyPair);
    secureEraseKeyPair(originalKeyPair);
    
    // Check if original is erased (all zeros)
    const isErased = originalKeyPair.privateKey.every(byte => byte === 0);
    const backupValid = toHex(backup.privateKey) === originalPrivateKey;
    
    console.log(`   Original key erased: ${isErased}`);
    console.log(`   Backup preserved: ${backupValid}`);

    // 9. Comprehensive validation
    console.log('\n9. Comprehensive Validation');
    console.log('-'.repeat(30));
    
    const validationTests = [
      await validateKeyPair(keyPair),
      await validateKeyPair(deterministicKeyPair1),
      await validateKeyPair(childKeyPair1),
      await validateKeyPair(passwordKeyPair1),
      await validateKeyPair(backup)
    ];
    
    const allValid = validationTests.every(test => test);
    console.log(`   All key pairs valid: ${allValid}`);
    console.log(`   Validation results: [${validationTests.join(', ')}]`);

    console.log('\n✅ Key pair management example completed successfully!');
    console.log('\nKey Features Demonstrated:');
    console.log('- ✓ Basic key generation');
    console.log('- ✓ Deterministic generation');
    console.log('- ✓ Batch operations');
    console.log('- ✓ Hierarchical derivation');
    console.log('- ✓ Password-based derivation');
    console.log('- ✓ Serialization/deserialization');
    console.log('- ✓ Secure memory management');

  } catch (error) {
    console.error('❌ Error in key pair example:', error);
    process.exit(1);
  }
}

// Run the example
if (require.main === module) {
  keyPairExample().catch(console.error);
}

export { keyPairExample };