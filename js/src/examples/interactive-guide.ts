/**
 * Interactive Guide Example for TOPAY-Z512
 *
 * This example provides an interactive command-line interface to explore
 * all TOPAY-Z512 features with guided tutorials and real-time feedback.
 */

import * as readline from 'readline';
import {
  generateKeyPair,
  computeHash,
  kemKeyGen,
  kemEncapsulate,
  kemDecapsulate,
  fragmentData,
  reconstructData,
  benchmarkHash,
  validateKeyPair,
  toHex,
  fromHex,
  measureTime,
  generateKeyPairFromSeed,
  secureRandom,
  deriveChildKeyPair,
  serializeKeyPair,
  deserializeKeyPair
} from '../index';

interface InteractiveSession {
  rl: readline.Interface;
  currentKeyPair?: any;
  currentKEMKeyPair?: any;
  currentFragments?: any[];
  sessionData: Map<string, any>;
}

class InteractiveGuide {
  private session: InteractiveSession;

  constructor() {
    this.session = {
      rl: readline.createInterface({
        input: process.stdin,
        output: process.stdout
      }),
      sessionData: new Map()
    };
  }

  async start(): Promise<void> {
    console.log('🎯 TOPAY-Z512 Interactive Guide');
    console.log('='.repeat(50));
    console.log('Welcome to the interactive exploration of TOPAY-Z512!');
    console.log('This guide will walk you through all the features step by step.\n');

    await this.showMainMenu();
  }

  private async showMainMenu(): Promise<void> {
    console.log('\n📋 Main Menu');
    console.log('-'.repeat(20));
    console.log('1. 🔑 Key Pair Operations');
    console.log('2. 🔐 Hash Operations');
    console.log('3. 🔒 KEM Operations');
    console.log('4. 🧩 Data Fragmentation');
    console.log('5. ⚡ Performance Benchmarks');
    console.log('6. 📊 Session Summary');
    console.log('7. 🚪 Exit');

    const choice = await this.prompt('\nSelect an option (1-7): ');

    switch (choice.trim()) {
      case '1':
        await this.keyPairMenu();
        break;
      case '2':
        await this.hashMenu();
        break;
      case '3':
        await this.kemMenu();
        break;
      case '4':
        await this.fragmentationMenu();
        break;
      case '5':
        await this.benchmarkMenu();
        break;
      case '6':
        await this.showSessionSummary();
        break;
      case '7':
        await this.exit();
        return;
      default:
        console.log('❌ Invalid option. Please try again.');
        await this.showMainMenu();
    }
  }

  private async keyPairMenu(): Promise<void> {
    console.log('\n🔑 Key Pair Operations');
    console.log('-'.repeat(25));
    console.log('1. Generate new key pair');
    console.log('2. Generate from seed');
    console.log('3. Derive child key pair');
    console.log('4. Validate current key pair');
    console.log('5. Serialize/Deserialize key pair');
    console.log('6. Show current key pair');
    console.log('7. Back to main menu');

    const choice = await this.prompt('\nSelect an option (1-7): ');

    switch (choice.trim()) {
      case '1':
        await this.generateKeyPair();
        break;
      case '2':
        await this.generateKeyPairFromSeed();
        break;
      case '3':
        await this.deriveChildKeyPair();
        break;
      case '4':
        await this.validateCurrentKeyPair();
        break;
      case '5':
        await this.serializeKeyPair();
        break;
      case '6':
        await this.showCurrentKeyPair();
        break;
      case '7':
        await this.showMainMenu();
        return;
      default:
        console.log('❌ Invalid option. Please try again.');
    }

    await this.keyPairMenu();
  }

  private async generateKeyPair(): Promise<void> {
    console.log('\n🔄 Generating new key pair...');

    const { timeMs, result: keyPair } = await measureTime(async () => {
      return await generateKeyPair();
    });

    this.session.currentKeyPair = keyPair;
    this.session.sessionData.set(
      'keyPairsGenerated',
      (this.session.sessionData.get('keyPairsGenerated') || 0) + 1
    );

    console.log(`✅ Key pair generated in ${timeMs.toFixed(3)}ms`);
    console.log(`   Private key: ${toHex(keyPair.privateKey).substring(0, 32)}...`);
    console.log(`   Public key:  ${toHex(keyPair.publicKey).substring(0, 32)}...`);
  }

  private async generateKeyPairFromSeed(): Promise<void> {
    const seedInput = await this.prompt('\nEnter seed (hex) or press Enter for random: ');

    let seed: Uint8Array;
    if (seedInput.trim()) {
      try {
        seed = fromHex(seedInput.trim());
      } catch {
        console.log('❌ Invalid hex format. Using random seed.');
        seed = await secureRandom(32);
      }
    } else {
      seed = await secureRandom(32);
    }

    console.log(`🌱 Using seed: ${toHex(seed).substring(0, 32)}...`);

    const { timeMs, result: keyPair } = await measureTime(async () => {
      return await generateKeyPairFromSeed(seed);
    });

    this.session.currentKeyPair = keyPair;

    console.log(`✅ Deterministic key pair generated in ${timeMs.toFixed(3)}ms`);
    console.log(`   Private key: ${toHex(keyPair.privateKey).substring(0, 32)}...`);
    console.log(`   Public key:  ${toHex(keyPair.publicKey).substring(0, 32)}...`);
  }

  private async deriveChildKeyPair(): Promise<void> {
    if (!this.session.currentKeyPair) {
      console.log('❌ No current key pair. Generate one first.');
      return;
    }

    const indexInput = await this.prompt('\nEnter child index (0-4294967295): ');
    const index = parseInt(indexInput.trim()) || 0;

    console.log(`🌿 Deriving child key pair at index ${index}...`);

    const { timeMs, result: childKeyPair } = await measureTime(async () => {
      return await deriveChildKeyPair(this.session.currentKeyPair.privateKey, index);
    });

    console.log(`✅ Child key pair derived in ${timeMs.toFixed(3)}ms`);
    console.log(`   Child private key: ${toHex(childKeyPair.privateKey).substring(0, 32)}...`);
    console.log(`   Child public key:  ${toHex(childKeyPair.publicKey).substring(0, 32)}...`);

    const useChild = await this.prompt('\nUse this as current key pair? (y/n): ');
    if (useChild.toLowerCase() === 'y') {
      this.session.currentKeyPair = childKeyPair;
      console.log('✅ Child key pair is now the current key pair.');
    }
  }

  private async validateCurrentKeyPair(): Promise<void> {
    if (!this.session.currentKeyPair) {
      console.log('❌ No current key pair. Generate one first.');
      return;
    }

    console.log('\n🔍 Validating current key pair...');

    const { timeMs, result: isValid } = await measureTime(async () => {
      return await validateKeyPair(this.session.currentKeyPair);
    });

    console.log(`${isValid ? '✅' : '❌'} Key pair validation: ${isValid ? 'VALID' : 'INVALID'}`);
    console.log(`   Validation time: ${timeMs.toFixed(3)}ms`);
  }

  private async serializeKeyPair(): Promise<void> {
    if (!this.session.currentKeyPair) {
      console.log('❌ No current key pair. Generate one first.');
      return;
    }

    console.log('\n📦 Serializing key pair...');

    const serialized = serializeKeyPair(this.session.currentKeyPair);
    console.log(`✅ Serialized length: ${serialized.length} characters`);
    console.log(`   Serialized data: ${serialized.substring(0, 64)}...`);

    const testRoundTrip = await this.prompt('\nTest round-trip serialization? (y/n): ');
    if (testRoundTrip.toLowerCase() === 'y') {
      const deserialized = deserializeKeyPair(serialized);
      const matches =
        toHex(this.session.currentKeyPair.privateKey) === toHex(deserialized.privateKey);
      console.log(`${matches ? '✅' : '❌'} Round-trip test: ${matches ? 'SUCCESS' : 'FAILED'}`);
    }
  }

  private async showCurrentKeyPair(): Promise<void> {
    if (!this.session.currentKeyPair) {
      console.log('❌ No current key pair. Generate one first.');
      return;
    }

    console.log('\n🔑 Current Key Pair');
    console.log('-'.repeat(20));
    console.log(`Private key: ${toHex(this.session.currentKeyPair.privateKey)}`);
    console.log(`Public key:  ${toHex(this.session.currentKeyPair.publicKey)}`);
  }

  private async hashMenu(): Promise<void> {
    console.log('\n🔐 Hash Operations');
    console.log('-'.repeat(20));
    console.log('1. Hash text input');
    console.log('2. Hash random data');
    console.log('3. Hash file simulation');
    console.log('4. Compare hash performance');
    console.log('5. Back to main menu');

    const choice = await this.prompt('\nSelect an option (1-5): ');

    switch (choice.trim()) {
      case '1':
        await this.hashTextInput();
        break;
      case '2':
        await this.hashRandomData();
        break;
      case '3':
        await this.hashFileSimulation();
        break;
      case '4':
        await this.compareHashPerformance();
        break;
      case '5':
        await this.showMainMenu();
        return;
      default:
        console.log('❌ Invalid option. Please try again.');
    }

    await this.hashMenu();
  }

  private async hashTextInput(): Promise<void> {
    const text = await this.prompt('\nEnter text to hash: ');
    const data = new TextEncoder().encode(text);

    console.log(`\n📝 Hashing "${text}"...`);

    const { timeMs, result: hash } = await measureTime(async () => {
      return await computeHash(data);
    });

    console.log(`✅ Hash computed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Input size: ${data.length} bytes`);
    console.log(`   Hash: ${toHex(hash)}`);

    this.session.sessionData.set(
      'hashesComputed',
      (this.session.sessionData.get('hashesComputed') || 0) + 1
    );
  }

  private async hashRandomData(): Promise<void> {
    const sizeInput = await this.prompt('\nEnter data size in bytes (default 1024): ');
    const size = parseInt(sizeInput.trim()) || 1024;

    const data = new Uint8Array(size);
    crypto.getRandomValues(data);

    console.log(`\n🎲 Hashing ${size} bytes of random data...`);

    const { timeMs, result: hash } = await measureTime(async () => {
      return await computeHash(data);
    });

    console.log(`✅ Hash computed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Throughput: ${(size / (timeMs / 1000) / 1024 / 1024).toFixed(2)} MB/s`);
    console.log(`   Hash: ${toHex(hash).substring(0, 32)}...`);
  }

  private async hashFileSimulation(): Promise<void> {
    const sizes = [1024, 10240, 102400, 1048576]; // 1KB, 10KB, 100KB, 1MB

    console.log('\n📁 File Hash Simulation');
    console.log('-'.repeat(25));
    console.log('Size      | Time (ms) | Throughput (MB/s)');
    console.log('-'.repeat(40));

    for (const size of sizes) {
      const data = new Uint8Array(size);
      crypto.getRandomValues(data);

      const { timeMs, result: hash } = await measureTime(async () => {
        return await computeHash(data);
      });

      const throughput = size / (timeMs / 1000) / 1024 / 1024;
      const sizeStr =
        size < 1024
          ? `${size}B`
          : size < 1048576
            ? `${(size / 1024).toFixed(0)}KB`
            : `${(size / 1048576).toFixed(1)}MB`;

      console.log(
        `${sizeStr.padEnd(9)} | ${timeMs.toFixed(3).padStart(9)} | ${throughput.toFixed(2).padStart(16)}`
      );
    }
  }

  private async compareHashPerformance(): Promise<void> {
    console.log('\n⚡ Hash Performance Comparison');
    console.log('-'.repeat(30));

    const benchmark = await benchmarkHash();

    console.log(
      `Operations per second: ${benchmark.throughputOpsPerSec?.toLocaleString() || 'N/A'}`
    );
    console.log(`Execution time: ${benchmark.executionTimeMs.toFixed(3)}ms`);
    console.log(`Operation: ${benchmark.operation}`);
  }

  private async kemMenu(): Promise<void> {
    console.log('\n🔒 KEM Operations');
    console.log('-'.repeat(20));
    console.log('1. Generate KEM key pair');
    console.log('2. Encapsulate shared secret');
    console.log('3. Decapsulate shared secret');
    console.log('4. Full KEM cycle demo');
    console.log('5. Show current KEM key pair');
    console.log('6. Back to main menu');

    const choice = await this.prompt('\nSelect an option (1-6): ');

    switch (choice.trim()) {
      case '1':
        await this.generateKEMKeyPair();
        break;
      case '2':
        await this.kemEncapsulate();
        break;
      case '3':
        await this.kemDecapsulate();
        break;
      case '4':
        await this.fullKEMDemo();
        break;
      case '5':
        await this.showCurrentKEMKeyPair();
        break;
      case '6':
        await this.showMainMenu();
        return;
      default:
        console.log('❌ Invalid option. Please try again.');
    }

    await this.kemMenu();
  }

  private async generateKEMKeyPair(): Promise<void> {
    console.log('\n🔄 Generating KEM key pair...');

    const { timeMs, result: kemKeyPair } = await measureTime(async () => {
      return await kemKeyGen();
    });

    this.session.currentKEMKeyPair = kemKeyPair;

    console.log(`✅ KEM key pair generated in ${timeMs.toFixed(3)}ms`);
    console.log(`   Secret key: ${toHex(kemKeyPair.secretKey).substring(0, 32)}...`);
    console.log(`   Public key: ${toHex(kemKeyPair.publicKey).substring(0, 32)}...`);
  }

  private async kemEncapsulate(): Promise<void> {
    if (!this.session.currentKEMKeyPair) {
      console.log('❌ No current KEM key pair. Generate one first.');
      return;
    }

    console.log('\n🔐 Encapsulating shared secret...');

    const { timeMs, result } = await measureTime(async () => {
      return await kemEncapsulate(this.session.currentKEMKeyPair.publicKey);
    });

    this.session.sessionData.set('lastCiphertext', result.ciphertext);
    this.session.sessionData.set('lastSharedSecret', result.sharedSecret);

    console.log(`✅ Encapsulation completed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Ciphertext: ${toHex(result.ciphertext).substring(0, 32)}...`);
    console.log(`   Shared secret: ${toHex(result.sharedSecret).substring(0, 32)}...`);
  }

  private async kemDecapsulate(): Promise<void> {
    if (!this.session.currentKEMKeyPair) {
      console.log('❌ No current KEM key pair. Generate one first.');
      return;
    }

    const ciphertext = this.session.sessionData.get('lastCiphertext');
    if (!ciphertext) {
      console.log('❌ No ciphertext available. Perform encapsulation first.');
      return;
    }

    console.log('\n🔓 Decapsulating shared secret...');

    const { timeMs, result: sharedSecret } = await measureTime(async () => {
      return await kemDecapsulate(this.session.currentKEMKeyPair.secretKey, ciphertext);
    });

    const originalSecret = this.session.sessionData.get('lastSharedSecret');
    const matches = originalSecret && toHex(sharedSecret) === toHex(originalSecret);

    console.log(`✅ Decapsulation completed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Shared secret: ${toHex(sharedSecret).substring(0, 32)}...`);
    console.log(`${matches ? '✅' : '❌'} Secret verification: ${matches ? 'MATCH' : 'MISMATCH'}`);
  }

  private async fullKEMDemo(): Promise<void> {
    console.log('\n🎭 Full KEM Cycle Demonstration');
    console.log('-'.repeat(35));

    // Generate key pair
    console.log('1. Generating KEM key pair...');
    const { timeMs: keyGenTime, result: kemKeyPair } = await measureTime(async () => {
      return await kemKeyGen();
    });
    console.log(`   ✅ Generated in ${keyGenTime.toFixed(3)}ms`);

    // Encapsulate
    console.log('2. Encapsulating shared secret...');
    const { timeMs: encapTime, result: encapResult } = await measureTime(async () => {
      return await kemEncapsulate(kemKeyPair.publicKey);
    });
    console.log(`   ✅ Encapsulated in ${encapTime.toFixed(3)}ms`);

    // Decapsulate
    console.log('3. Decapsulating shared secret...');
    const { timeMs: decapTime, result: decapSecret } = await measureTime(async () => {
      return await kemDecapsulate(kemKeyPair.secretKey, encapResult.ciphertext);
    });
    console.log(`   ✅ Decapsulated in ${decapTime.toFixed(3)}ms`);

    // Verify
    const secretsMatch = toHex(encapResult.sharedSecret) === toHex(decapSecret);
    console.log(`4. Verifying shared secrets...`);
    console.log(
      `   ${secretsMatch ? '✅' : '❌'} Verification: ${secretsMatch ? 'SUCCESS' : 'FAILED'}`
    );

    const totalTime = keyGenTime + encapTime + decapTime;
    console.log(`\n📊 Total cycle time: ${totalTime.toFixed(3)}ms`);
  }

  private async showCurrentKEMKeyPair(): Promise<void> {
    if (!this.session.currentKEMKeyPair) {
      console.log('❌ No current KEM key pair. Generate one first.');
      return;
    }

    console.log('\n🔒 Current KEM Key Pair');
    console.log('-'.repeat(25));
    console.log(`Secret key: ${toHex(this.session.currentKEMKeyPair.secretKey)}`);
    console.log(`Public key: ${toHex(this.session.currentKEMKeyPair.publicKey)}`);
  }

  private async fragmentationMenu(): Promise<void> {
    console.log('\n🧩 Data Fragmentation');
    console.log('-'.repeat(25));
    console.log('1. Fragment text input');
    console.log('2. Fragment random data');
    console.log('3. Reconstruct fragments');
    console.log('4. Fragment loss simulation');
    console.log('5. Show current fragments');
    console.log('6. Back to main menu');

    const choice = await this.prompt('\nSelect an option (1-6): ');

    switch (choice.trim()) {
      case '1':
        await this.fragmentTextInput();
        break;
      case '2':
        await this.fragmentRandomData();
        break;
      case '3':
        await this.reconstructFragments();
        break;
      case '4':
        await this.fragmentLossSimulation();
        break;
      case '5':
        await this.showCurrentFragments();
        break;
      case '6':
        await this.showMainMenu();
        return;
      default:
        console.log('❌ Invalid option. Please try again.');
    }

    await this.fragmentationMenu();
  }

  private async fragmentTextInput(): Promise<void> {
    const text = await this.prompt('\nEnter text to fragment: ');
    const data = new TextEncoder().encode(text);

    console.log(`\n✂️ Fragmenting "${text}"...`);

    const { timeMs, result } = await measureTime(async () => {
      return await fragmentData(data);
    });

    this.session.currentFragments = result.fragments;

    console.log(`✅ Fragmentation completed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Original size: ${result.metadata.originalSize} bytes`);
    console.log(`   Fragment count: ${result.metadata.fragmentCount}`);
    console.log(`   Checksum: ${toHex(result.metadata.checksum).substring(0, 32)}...`);
  }

  private async fragmentRandomData(): Promise<void> {
    const sizeInput = await this.prompt('\nEnter data size in bytes (default 2048): ');
    const size = parseInt(sizeInput.trim()) || 2048;

    const data = new Uint8Array(size);
    crypto.getRandomValues(data);

    console.log(`\n✂️ Fragmenting ${size} bytes of random data...`);

    const { timeMs, result } = await measureTime(async () => {
      return await fragmentData(data);
    });

    this.session.currentFragments = result.fragments;

    console.log(`✅ Fragmentation completed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Fragment count: ${result.metadata.fragmentCount}`);
    console.log(`   Throughput: ${(size / (timeMs / 1000) / 1024 / 1024).toFixed(2)} MB/s`);
  }

  private async reconstructFragments(): Promise<void> {
    if (!this.session.currentFragments) {
      console.log('❌ No current fragments. Fragment some data first.');
      return;
    }

    console.log(`\n🔧 Reconstructing ${this.session.currentFragments.length} fragments...`);

    const { timeMs, result } = await measureTime(async () => {
      return await reconstructData(this.session.currentFragments!);
    });

    console.log(`✅ Reconstruction completed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Reconstruction complete: ${result.isComplete}`);
    console.log(`   Reconstructed size: ${result.data.length} bytes`);

    if (result.data.length < 200) {
      const text = new TextDecoder().decode(result.data);
      console.log(`   Reconstructed text: "${text}"`);
    }
  }

  private async fragmentLossSimulation(): Promise<void> {
    if (!this.session.currentFragments) {
      console.log('❌ No current fragments. Fragment some data first.');
      return;
    }

    const lossInput = await this.prompt('\nEnter number of fragments to lose (default 1): ');
    const lossCount = parseInt(lossInput.trim()) || 1;

    const originalCount = this.session.currentFragments.length;
    const partialFragments = this.session.currentFragments.slice(0, -lossCount);

    console.log(`\n💥 Simulating loss of ${lossCount} fragments...`);
    console.log(`   Original fragments: ${originalCount}`);
    console.log(`   Remaining fragments: ${partialFragments.length}`);

    const { timeMs, result } = await measureTime(async () => {
      return await reconstructData(partialFragments);
    });

    console.log(`✅ Partial reconstruction completed in ${timeMs.toFixed(3)}ms`);
    console.log(`   Reconstruction complete: ${result.isComplete}`);
    console.log(`   Reconstructed size: ${result.data.length} bytes`);
  }

  private async showCurrentFragments(): Promise<void> {
    if (!this.session.currentFragments) {
      console.log('❌ No current fragments. Fragment some data first.');
      return;
    }

    console.log('\n🧩 Current Fragments');
    console.log('-'.repeat(20));
    console.log(`Fragment count: ${this.session.currentFragments.length}`);

    for (let i = 0; i < Math.min(5, this.session.currentFragments.length); i++) {
      const fragment = this.session.currentFragments[i]!;
      console.log(`   Fragment ${fragment.index}: ${fragment.data.length} bytes`);
    }

    if (this.session.currentFragments.length > 5) {
      console.log(`   ... and ${this.session.currentFragments.length - 5} more fragments`);
    }
  }

  private async benchmarkMenu(): Promise<void> {
    console.log('\n⚡ Performance Benchmarks');
    console.log('-'.repeat(30));
    console.log('1. Hash operations benchmark');
    console.log('2. Key pair generation benchmark');
    console.log('3. KEM operations benchmark');
    console.log('4. Fragmentation benchmark');
    console.log('5. Custom benchmark');
    console.log('6. Back to main menu');

    const choice = await this.prompt('\nSelect an option (1-6): ');

    switch (choice.trim()) {
      case '1':
        await this.runHashBenchmark();
        break;
      case '2':
        await this.runKeyPairBenchmark();
        break;
      case '3':
        await this.runKEMBenchmark();
        break;
      case '4':
        await this.runFragmentationBenchmark();
        break;
      case '5':
        await this.runCustomBenchmark();
        break;
      case '6':
        await this.showMainMenu();
        return;
      default:
        console.log('❌ Invalid option. Please try again.');
    }

    await this.benchmarkMenu();
  }

  private async runHashBenchmark(): Promise<void> {
    console.log('\n⚡ Running hash operations benchmark...');

    const benchmark = await benchmarkHash();

    console.log('✅ Hash Benchmark Results:');
    console.log(
      `   Operations per second: ${benchmark.throughputOpsPerSec?.toLocaleString() || 'N/A'}`
    );
    console.log(`   Execution time: ${benchmark.executionTimeMs.toFixed(3)}ms`);
    console.log(`   Operation: ${benchmark.operation}`);
  }

  private async runKeyPairBenchmark(): Promise<void> {
    console.log('\n⚡ Running key pair generation benchmark...');

    const startTime = Date.now();
    const keyPairs = [];

    for (let i = 0; i < 10; i++) {
      keyPairs.push(await generateKeyPair());
    }

    const endTime = Date.now();
    const totalTime = endTime - startTime;

    console.log('✅ Key Pair Benchmark Results:');
    console.log(`   Generated ${keyPairs.length} key pairs`);
    console.log(`   Total time: ${totalTime}ms`);
    console.log(`   Average time per key pair: ${(totalTime / keyPairs.length).toFixed(3)}ms`);
    console.log(`   Key pairs per second: ${(keyPairs.length / (totalTime / 1000)).toFixed(1)}`);
  }

  private async runKEMBenchmark(): Promise<void> {
    console.log('\n⚡ Running KEM operations benchmark...');

    const kemKeyPair = await kemKeyGen();
    const operations = 10;

    const startTime = Date.now();

    for (let i = 0; i < operations; i++) {
      const { ciphertext, sharedSecret } = await kemEncapsulate(kemKeyPair.publicKey);
      await kemDecapsulate(kemKeyPair.secretKey, ciphertext);
    }

    const endTime = Date.now();
    const totalTime = endTime - startTime;

    console.log('✅ KEM Benchmark Results:');
    console.log(`   Completed ${operations} full KEM cycles`);
    console.log(`   Total time: ${totalTime}ms`);
    console.log(`   Average time per cycle: ${(totalTime / operations).toFixed(3)}ms`);
    console.log(`   Cycles per second: ${(operations / (totalTime / 1000)).toFixed(1)}`);
  }

  private async runFragmentationBenchmark(): Promise<void> {
    console.log('\n⚡ Running fragmentation benchmark...');

    const testData = new Uint8Array(4096);
    crypto.getRandomValues(testData);
    const operations = 10;

    const startTime = Date.now();

    for (let i = 0; i < operations; i++) {
      const fragResult = await fragmentData(testData);
      await reconstructData(fragResult.fragments);
    }

    const endTime = Date.now();
    const totalTime = endTime - startTime;

    console.log('✅ Fragmentation Benchmark Results:');
    console.log(`   Completed ${operations} fragmentation cycles`);
    console.log(`   Data size: ${testData.length} bytes per operation`);
    console.log(`   Total time: ${totalTime}ms`);
    console.log(`   Average time per cycle: ${(totalTime / operations).toFixed(3)}ms`);
    console.log(
      `   Throughput: ${((testData.length * operations) / (totalTime / 1000) / 1024 / 1024).toFixed(2)} MB/s`
    );
  }

  private async runCustomBenchmark(): Promise<void> {
    console.log('\n🎯 Custom Benchmark Configuration');
    console.log('-'.repeat(35));

    const operation = await this.prompt('Select operation (hash/keygen/kem/fragment): ');
    const iterations = parseInt(await this.prompt('Number of iterations (default 100): ')) || 100;

    console.log(`\n⚡ Running custom benchmark: ${operation} x${iterations}...`);

    const startTime = Date.now();

    switch (operation.toLowerCase()) {
      case 'hash':
        const data = new Uint8Array(1024);
        for (let i = 0; i < iterations; i++) {
          crypto.getRandomValues(data);
          await computeHash(data);
        }
        break;

      case 'keygen':
        for (let i = 0; i < iterations; i++) {
          await generateKeyPair();
        }
        break;

      case 'kem':
        const kemKeyPair = await kemKeyGen();
        for (let i = 0; i < iterations; i++) {
          const { ciphertext } = await kemEncapsulate(kemKeyPair.publicKey);
          await kemDecapsulate(kemKeyPair.secretKey, ciphertext);
        }
        break;

      case 'fragment':
        const fragData = new Uint8Array(2048);
        for (let i = 0; i < iterations; i++) {
          crypto.getRandomValues(fragData);
          const fragResult = await fragmentData(fragData);
          await reconstructData(fragResult.fragments);
        }
        break;

      default:
        console.log('❌ Unknown operation. Skipping benchmark.');
        return;
    }

    const endTime = Date.now();
    const totalTime = endTime - startTime;

    console.log('✅ Custom Benchmark Results:');
    console.log(`   Operation: ${operation}`);
    console.log(`   Iterations: ${iterations.toLocaleString()}`);
    console.log(`   Total time: ${totalTime}ms`);
    console.log(`   Average time: ${(totalTime / iterations).toFixed(3)}ms`);
    console.log(`   Operations per second: ${(iterations / (totalTime / 1000)).toFixed(1)}`);
  }

  private async showSessionSummary(): Promise<void> {
    console.log('\n📊 Session Summary');
    console.log('-'.repeat(20));

    const keyPairsGenerated = this.session.sessionData.get('keyPairsGenerated') || 0;
    const hashesComputed = this.session.sessionData.get('hashesComputed') || 0;

    console.log(`Key pairs generated: ${keyPairsGenerated}`);
    console.log(`Hashes computed: ${hashesComputed}`);
    console.log(`Current key pair: ${this.session.currentKeyPair ? 'Yes' : 'No'}`);
    console.log(`Current KEM key pair: ${this.session.currentKEMKeyPair ? 'Yes' : 'No'}`);
    console.log(
      `Current fragments: ${this.session.currentFragments ? this.session.currentFragments.length : 0}`
    );

    if (this.session.currentKeyPair) {
      console.log(`\nCurrent key pair details:`);
      console.log(
        `   Private: ${toHex(this.session.currentKeyPair.privateKey).substring(0, 32)}...`
      );
      console.log(
        `   Public:  ${toHex(this.session.currentKeyPair.publicKey).substring(0, 32)}...`
      );
    }

    await this.prompt('\nPress Enter to continue...');
    await this.showMainMenu();
  }

  private async exit(): Promise<void> {
    console.log('\n👋 Thank you for exploring TOPAY-Z512!');
    console.log('Visit our documentation for more advanced features.');
    this.session.rl.close();
  }

  private async prompt(question: string): Promise<string> {
    return new Promise(resolve => {
      this.session.rl.question(question, resolve);
    });
  }
}

async function interactiveGuideExample(): Promise<void> {
  const guide = new InteractiveGuide();
  await guide.start();
}

// Run the example
if (require.main === module) {
  interactiveGuideExample().catch(console.error);
}

export { interactiveGuideExample };
