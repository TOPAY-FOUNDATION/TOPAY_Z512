/**
 * KEM (Key Encapsulation Mechanism) Example for TOPAY-Z512
 *
 * This example demonstrates advanced KEM operations including
 * batch processing, validation, and secure key exchange scenarios.
 */

import {
  kemKeyGen,
  kemEncapsulate,
  kemDecapsulate,
  batchKEMKeyGen,
  batchKEMEncapsulate,
  batchKEMDecapsulate,
  validateKEMKeyPair,
  testKEMOperations,
  secureEraseKEMKeyPair,
  backupKEMKeyPair,
  serializeKEMKeyPair,
  deserializeKEMKeyPair,
  toHex,
  measureTime,
  secureRandom
} from '../index';

async function kemExample(): Promise<void> {
  console.log('🔒 TOPAY-Z512 KEM Operations Example');
  console.log('='.repeat(50));

  try {
    // 1. Basic KEM operations
    console.log('\n1. Basic KEM Operations');
    console.log('-'.repeat(30));

    const { result: kemKeyPair, timeMs: keyGenTime } = await measureTime(async () => {
      return await kemKeyGen();
    });

    console.log(`   KEM key pair generated in ${keyGenTime.toFixed(3)}ms`);
    console.log(`   Secret key: ${toHex(kemKeyPair.secretKey).substring(0, 32)}...`);
    console.log(`   Public key: ${toHex(kemKeyPair.publicKey).substring(0, 32)}...`);

    // Validate the key pair
    const isValid = await validateKEMKeyPair(kemKeyPair);
    console.log(`   Key pair validation: ${isValid ? 'VALID' : 'INVALID'}`);

    // 2. Encapsulation
    console.log('\n2. Key Encapsulation');
    console.log('-'.repeat(30));

    const { result: encapResult, timeMs: encapTime } = await measureTime(async () => {
      return await kemEncapsulate(kemKeyPair.publicKey);
    });

    console.log(`   Encapsulation completed in ${encapTime.toFixed(3)}ms`);
    console.log(`   Ciphertext length: ${encapResult.ciphertext.length} bytes`);
    console.log(`   Shared secret length: ${encapResult.sharedSecret.length} bytes`);
    console.log(`   Ciphertext: ${toHex(encapResult.ciphertext).substring(0, 32)}...`);
    console.log(`   Shared secret: ${toHex(encapResult.sharedSecret).substring(0, 32)}...`);

    // 3. Decapsulation
    console.log('\n3. Key Decapsulation');
    console.log('-'.repeat(30));

    const { result: decapSecret, timeMs: decapTime } = await measureTime(async () => {
      return await kemDecapsulate(kemKeyPair.secretKey, encapResult.ciphertext);
    });

    console.log(`   Decapsulation completed in ${decapTime.toFixed(3)}ms`);
    console.log(`   Decapsulated secret: ${toHex(decapSecret).substring(0, 32)}...`);

    const secretsMatch = toHex(encapResult.sharedSecret) === toHex(decapSecret);
    console.log(`   Secret verification: ${secretsMatch ? 'SUCCESS' : 'FAILED'}`);

    // 4. Batch KEM operations
    console.log('\n4. Batch KEM Operations');
    console.log('-'.repeat(30));

    const batchSize = 10;
    console.log(`   Generating ${batchSize} KEM key pairs...`);

    const { result: batchKeyPairs, timeMs: batchKeyGenTime } = await measureTime(async () => {
      return await batchKEMKeyGen(batchSize);
    });

    console.log(`   Batch key generation: ${batchKeyGenTime.toFixed(3)}ms`);
    console.log(`   Average per key pair: ${(batchKeyGenTime / batchSize).toFixed(3)}ms`);

    // Extract public keys for batch encapsulation
    const publicKeys = batchKeyPairs.map(kp => kp.publicKey);

    console.log(`   Performing batch encapsulation...`);
    const { result: batchEncapResults, timeMs: batchEncapTime } = await measureTime(async () => {
      return await batchKEMEncapsulate(publicKeys);
    });

    console.log(`   Batch encapsulation: ${batchEncapTime.toFixed(3)}ms`);
    console.log(`   Average per encapsulation: ${(batchEncapTime / batchSize).toFixed(3)}ms`);

    // Batch decapsulation
    const secretKeys = batchKeyPairs.map(kp => kp.secretKey);
    const ciphertexts = batchEncapResults.map(result => result.ciphertext);

    console.log(`   Performing batch decapsulation...`);
    const { result: batchDecapSecrets, timeMs: batchDecapTime } = await measureTime(async () => {
      const decapInputs = secretKeys.map((secretKey, index) => ({
        secretKey,
        ciphertext: ciphertexts[index]!
      }));
      return await batchKEMDecapsulate(decapInputs);
    });

    console.log(`   Batch decapsulation: ${batchDecapTime.toFixed(3)}ms`);
    console.log(`   Average per decapsulation: ${(batchDecapTime / batchSize).toFixed(3)}ms`);

    // Verify all secrets match
    let matchCount = 0;
    for (let i = 0; i < batchSize; i++) {
      const originalSecret = batchEncapResults[i]!.sharedSecret;
      const decapsulatedSecret = batchDecapSecrets[i]!;
      if (toHex(originalSecret) === toHex(decapsulatedSecret)) {
        matchCount++;
      }
    }
    console.log(`   Secret verification: ${matchCount}/${batchSize} matches`);

    // 5. KEM operations testing
    console.log('\n5. KEM Operations Testing');
    console.log('-'.repeat(30));

    const testIterations = 10;
    let successCount = 0;
    const keyGenTimes = [];
    const encapTimes = [];
    const decapTimes = [];

    const startTime = Date.now();

    for (let i = 0; i < testIterations; i++) {
      const { result: testKeyPair, timeMs: keyGenTime } = await measureTime(async () => {
        return await kemKeyGen();
      });
      keyGenTimes.push(keyGenTime);

      const { result: encapResult, timeMs: encapTime } = await measureTime(async () => {
        return await kemEncapsulate(testKeyPair.publicKey);
      });
      encapTimes.push(encapTime);

      const { result: decapResult, timeMs: decapTime } = await measureTime(async () => {
        return await kemDecapsulate(testKeyPair.secretKey, encapResult.ciphertext);
      });
      decapTimes.push(decapTime);

      const testSuccess = await testKEMOperations(testKeyPair);
      if (testSuccess) successCount++;
    }

    const totalTime = Date.now() - startTime;
    const successRate = successCount / testIterations;
    const avgKeyGenTime = keyGenTimes.reduce((a, b) => a + b) / keyGenTimes.length;
    const avgEncapTime = encapTimes.reduce((a, b) => a + b) / encapTimes.length;
    const avgDecapTime = decapTimes.reduce((a, b) => a + b) / decapTimes.length;

    console.log(`   Test iterations: ${testIterations}`);
    console.log(`   Success rate: ${(successRate * 100).toFixed(1)}%`);
    console.log(`   Average key generation: ${avgKeyGenTime.toFixed(3)}ms`);
    console.log(`   Average encapsulation: ${avgEncapTime.toFixed(3)}ms`);
    console.log(`   Average decapsulation: ${avgDecapTime.toFixed(3)}ms`);
    console.log(`   Total test time: ${totalTime.toFixed(3)}ms`);

    // 6. Key exchange simulation
    console.log('\n6. Secure Key Exchange Simulation');
    console.log('-'.repeat(30));

    // Alice generates her KEM key pair
    const aliceKeyPair = await kemKeyGen();
    console.log('   Alice generated her KEM key pair');

    // Alice sends her public key to Bob (simulated)
    const alicePublicKey = aliceKeyPair.publicKey;
    console.log(`   Alice's public key: ${toHex(alicePublicKey).substring(0, 32)}...`);

    // Bob encapsulates a shared secret using Alice's public key
    const bobEncapResult = await kemEncapsulate(alicePublicKey);
    console.log('   Bob encapsulated shared secret');
    console.log(`   Bob's ciphertext: ${toHex(bobEncapResult.ciphertext).substring(0, 32)}...`);

    // Bob sends the ciphertext to Alice (simulated)
    // Alice decapsulates the shared secret
    const aliceSharedSecret = await kemDecapsulate(
      aliceKeyPair.secretKey,
      bobEncapResult.ciphertext
    );
    console.log('   Alice decapsulated shared secret');

    // Verify both parties have the same shared secret
    const keyExchangeSuccess = toHex(bobEncapResult.sharedSecret) === toHex(aliceSharedSecret);
    console.log(`   Key exchange: ${keyExchangeSuccess ? 'SUCCESS' : 'FAILED'}`);
    console.log(`   Shared secret: ${toHex(aliceSharedSecret).substring(0, 32)}...`);

    // 7. Multiple party key exchange
    console.log('\n7. Multiple Party Key Exchange');
    console.log('-'.repeat(30));

    const parties = ['Alice', 'Bob', 'Charlie', 'Diana'];
    const partyKeyPairs = [];
    const sharedSecrets = new Map();

    // Each party generates a key pair
    for (const party of parties) {
      const keyPair = await kemKeyGen();
      partyKeyPairs.push({ party, keyPair });
      console.log(`   ${party} generated key pair`);
    }

    // Each party establishes shared secrets with every other party
    for (let i = 0; i < parties.length; i++) {
      for (let j = i + 1; j < parties.length; j++) {
        const party1 = partyKeyPairs[i]!;
        const party2 = partyKeyPairs[j]!;

        // Party1 encapsulates using Party2's public key
        const encapResult = await kemEncapsulate(party2.keyPair.publicKey);
        const decapSecret = await kemDecapsulate(party2.keyPair.secretKey, encapResult.ciphertext);

        const secretKey = `${party1.party}-${party2.party}`;
        sharedSecrets.set(secretKey, {
          encapsulated: encapResult.sharedSecret,
          decapsulated: decapSecret,
          match: toHex(encapResult.sharedSecret) === toHex(decapSecret)
        });
      }
    }

    console.log(`   Established ${sharedSecrets.size} pairwise shared secrets`);
    const successfulExchanges = Array.from(sharedSecrets.values()).filter(s => s.match).length;
    console.log(`   Successful exchanges: ${successfulExchanges}/${sharedSecrets.size}`);

    // 8. KEM key pair serialization
    console.log('\n8. KEM Key Pair Serialization');
    console.log('-'.repeat(30));

    const testKemKeyPair = await kemKeyGen();
    const serialized = serializeKEMKeyPair(testKemKeyPair);
    console.log(`   Serialized length: ${serialized.length} characters`);

    const deserialized = deserializeKEMKeyPair(serialized);
    const serializationValid =
      toHex(testKemKeyPair.secretKey) === toHex(deserialized.secretKey) &&
      toHex(testKemKeyPair.publicKey) === toHex(deserialized.publicKey);
    console.log(`   Serialization round-trip: ${serializationValid ? 'SUCCESS' : 'FAILED'}`);

    // 9. Secure memory management
    console.log('\n9. Secure Memory Management');
    console.log('-'.repeat(30));

    const memTestKeyPair = await kemKeyGen();
    const originalSecretKey = toHex(memTestKeyPair.secretKey);

    // Create backup before erasure
    const backup = backupKEMKeyPair(memTestKeyPair);
    console.log('   Created backup of KEM key pair');

    // Securely erase the original
    secureEraseKEMKeyPair(memTestKeyPair);

    // Check if original is erased
    const isErased =
      memTestKeyPair.secretKey.every(byte => byte === 0) &&
      memTestKeyPair.publicKey.every(byte => byte === 0);
    const backupValid = toHex(backup.secretKey) === originalSecretKey;

    console.log(`   Original key pair erased: ${isErased}`);
    console.log(`   Backup integrity: ${backupValid ? 'PRESERVED' : 'CORRUPTED'}`);

    // 10. Performance comparison
    console.log('\n10. Performance Comparison');
    console.log('-'.repeat(30));

    const iterations = 100;
    console.log(`   Running ${iterations} iterations for each operation...`);

    // Key generation performance
    const perfKeyGenTimes = [];
    for (let i = 0; i < iterations; i++) {
      const { timeMs } = await measureTime(async () => {
        return await kemKeyGen();
      });
      perfKeyGenTimes.push(timeMs);
    }

    // Encapsulation performance
    const testKeyPair = await kemKeyGen();
    const perfEncapTimes = [];
    for (let i = 0; i < iterations; i++) {
      const { timeMs } = await measureTime(async () => {
        return await kemEncapsulate(testKeyPair.publicKey);
      });
      perfEncapTimes.push(timeMs);
    }

    // Decapsulation performance
    const { ciphertext } = await kemEncapsulate(testKeyPair.publicKey);
    const perfDecapTimes = [];
    for (let i = 0; i < iterations; i++) {
      const { timeMs } = await measureTime(async () => {
        return await kemDecapsulate(testKeyPair.secretKey, ciphertext);
      });
      perfDecapTimes.push(timeMs);
    }

    const avgKeyGen = perfKeyGenTimes.reduce((a, b) => a + b) / perfKeyGenTimes.length;
    const avgEncap = perfEncapTimes.reduce((a, b) => a + b) / perfEncapTimes.length;
    const avgDecap = perfDecapTimes.reduce((a, b) => a + b) / perfDecapTimes.length;

    console.log(`   Average key generation: ${avgKeyGen.toFixed(3)}ms`);
    console.log(`   Average encapsulation: ${avgEncap.toFixed(3)}ms`);
    console.log(`   Average decapsulation: ${avgDecap.toFixed(3)}ms`);
    console.log(`   Total cycle time: ${(avgKeyGen + avgEncap + avgDecap).toFixed(3)}ms`);

    console.log('\n✅ KEM example completed successfully!');
    console.log('\nKEM Features Demonstrated:');
    console.log('- ✓ Basic KEM operations (KeyGen, Encap, Decap)');
    console.log('- ✓ Batch processing capabilities');
    console.log('- ✓ Key pair validation');
    console.log('- ✓ Secure key exchange protocols');
    console.log('- ✓ Multiple party scenarios');
    console.log('- ✓ Serialization and deserialization');
    console.log('- ✓ Secure memory management');
    console.log('- ✓ Performance analysis');
  } catch (error) {
    console.error('❌ Error in KEM example:', error);
    process.exit(1);
  }
}

// Run the example
if (require.main === module) {
  kemExample().catch(console.error);
}

export { kemExample };
