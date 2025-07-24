/**
 * Comprehensive Test Suite for TOPAY-Z512 JavaScript/TypeScript Implementation
 *
 * This test suite covers all major functionality including:
 * - Hash operations
 * - Key pair generation and management
 * - KEM operations
 * - Data fragmentation
 * - Performance benchmarks
 * - Utility functions
 */

import {
  // Hash functions
  computeHash,
  computeHashWithSalt,
  computeHmac,
  batchHash,
  computeMerkleRoot,
  deriveKey,
  computeHashChain,
  verifyHashChain,

  // Key pair functions
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

  // KEM functions
  kemKeyGen,
  kemEncapsulate,
  kemDecapsulate,
  batchKEMKeyGen,
  batchKEMEncapsulate,
  batchKEMDecapsulate,
  validateKEMKeyPair,
  // testKEMOperations as kemOperationsTest, // Commented out unused import
  secureEraseKEMKeyPair,
  backupKEMKeyPair,
  serializeKEMKeyPair,
  deserializeKEMKeyPair,

  // Fragmentation functions
  fragmentData,
  reconstructData,
  parallelFragmentation,
  parallelReconstruction,
  validateFragment,
  estimateMobileLatency,
  getOptimalFragmentSize,
  serializeFragments,
  deserializeFragments,
  compressFragments,
  decompressFragments,

  // Performance functions
  benchmarkHash,
  benchmarkKeyPairGeneration,
  benchmarkKEM,
  benchmarkFragmentation,
  runBenchmarkSuite,
  monitorMemoryUsage,
  profileCPUUsage,
  estimateMobilePerformance,
  generatePerformanceReport,

  // Utility functions
  secureRandom,
  constantTimeEqual,
  secureZero,
  toHex,
  fromHex,
  xorBytes,
  validateSize,
  copyBytes,
  concatBytes,
  // timestamp, // Commented out unused import
  sleep,
  measureTime,
  hasWebCrypto,
  getSystemCapabilities,

  // Constants
  VERSION,
  PRIVATE_KEY_SIZE,
  PUBLIC_KEY_SIZE,
  HASH_SIZE,
  KEM_SECRET_KEY_SIZE,
  KEM_PUBLIC_KEY_SIZE,
  KEM_CIPHERTEXT_SIZE,
  KEM_SHARED_SECRET_SIZE,
  DEFAULT_FRAGMENT_SIZE,
  MAX_FRAGMENT_SIZE,
  MIN_FRAGMENT_SIZE
} from '../index';

interface TestResult {
  name: string;
  passed: boolean;
  error?: string;
  duration: number;
}

interface TestSuite {
  name: string;
  tests: TestResult[];
  totalTests: number;
  passedTests: number;
  failedTests: number;
  totalDuration: number;
}

class TestRunner {
  private suites: TestSuite[] = [];
  private currentSuite: TestSuite | null = null;

  startSuite(name: string): void {
    this.currentSuite = {
      name,
      tests: [],
      totalTests: 0,
      passedTests: 0,
      failedTests: 0,
      totalDuration: 0
    };
  }

  async runTest(name: string, testFn: () => Promise<void>): Promise<void> {
    if (!this.currentSuite) {
      throw new Error('No test suite started');
    }

    const startTime = Date.now();
    let passed = false;
    let error: string | undefined;

    try {
      await testFn();
      passed = true;
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    }

    const duration = Date.now() - startTime;

    const result: TestResult = {
      name,
      passed,
      duration,
      ...(error && { error })
    };

    this.currentSuite.tests.push(result);
    this.currentSuite.totalTests++;
    this.currentSuite.totalDuration += duration;

    if (passed) {
      this.currentSuite.passedTests++;
    } else {
      this.currentSuite.failedTests++;
    }

    const status = passed ? '✅' : '❌';
    console.log(`   ${status} ${name} (${duration}ms)`);
    if (error) {
      console.log(`      Error: ${error}`);
    }
  }

  endSuite(): void {
    if (!this.currentSuite) {
      throw new Error('No test suite to end');
    }

    this.suites.push(this.currentSuite);

    const suite = this.currentSuite;
    console.log(`\n📊 ${suite.name} Results:`);
    console.log(
      `   Total: ${suite.totalTests}, Passed: ${suite.passedTests}, Failed: ${suite.failedTests}`
    );
    console.log(`   Duration: ${suite.totalDuration}ms`);
    console.log(`   Success Rate: ${((suite.passedTests / suite.totalTests) * 100).toFixed(1)}%`);

    this.currentSuite = null;
  }

  getSummary(): {
    totalTests: number;
    totalPassed: number;
    totalFailed: number;
    totalDuration: number;
  } {
    return this.suites.reduce(
      (acc, suite) => ({
        totalTests: acc.totalTests + suite.totalTests,
        totalPassed: acc.totalPassed + suite.passedTests,
        totalFailed: acc.totalFailed + suite.failedTests,
        totalDuration: acc.totalDuration + suite.totalDuration
      }),
      { totalTests: 0, totalPassed: 0, totalFailed: 0, totalDuration: 0 }
    );
  }
}

async function runAllTests(): Promise<void> {
  console.log('🧪 TOPAY-Z512 Test Suite');
  console.log('='.repeat(50));
  console.log(`Version: ${VERSION}`);
  console.log(`Node.js: ${process.version}`);
  console.log(`Platform: ${process.platform} ${process.arch}`);
  console.log('');

  const runner = new TestRunner();

  // Test Constants
  await testConstants(runner);

  // Test Utility Functions
  await testUtilities(runner);

  // Test Hash Operations
  await testHashOperations(runner);

  // Test Key Pair Operations
  await testKeyPairOperations(runner);

  // Test KEM Operations
  await testKEMOperations(runner);

  // Test Fragmentation
  await testFragmentation(runner);

  // Test Performance
  await testPerformance(runner);

  // Test Integration
  await testIntegration(runner);

  // Print final summary
  const summary = runner.getSummary();
  console.log('\n🎯 Final Test Summary');
  console.log('='.repeat(30));
  console.log(`Total Tests: ${summary.totalTests}`);
  console.log(`Passed: ${summary.totalPassed}`);
  console.log(`Failed: ${summary.totalFailed}`);
  console.log(`Success Rate: ${((summary.totalPassed / summary.totalTests) * 100).toFixed(1)}%`);
  console.log(`Total Duration: ${summary.totalDuration}ms`);

  if (summary.totalFailed > 0) {
    console.log('\n❌ Some tests failed!');
    throw new Error(`${summary.totalFailed} out of ${summary.totalTests} tests failed`);
  } else {
    console.log('\n✅ All tests passed!');
  }
}

async function testConstants(runner: TestRunner): Promise<void> {
  runner.startSuite('Constants');

  await runner.runTest('Version constant exists', async () => {
    if (typeof VERSION !== 'string' || VERSION.length === 0) {
      throw new Error('VERSION constant is invalid');
    }
  });

  await runner.runTest('Key size constants', async () => {
    if (PRIVATE_KEY_SIZE !== 64) throw new Error('Invalid PRIVATE_KEY_SIZE');
    if (PUBLIC_KEY_SIZE !== 64) throw new Error('Invalid PUBLIC_KEY_SIZE');
    if (HASH_SIZE !== 64) throw new Error('Invalid HASH_SIZE');
  });

  await runner.runTest('KEM size constants', async () => {
    if (KEM_SECRET_KEY_SIZE !== 64) throw new Error('Invalid KEM_SECRET_KEY_SIZE');
    if (KEM_PUBLIC_KEY_SIZE !== 64) throw new Error('Invalid KEM_PUBLIC_KEY_SIZE');
    if (KEM_CIPHERTEXT_SIZE !== 64) throw new Error('Invalid KEM_CIPHERTEXT_SIZE');
    if (KEM_SHARED_SECRET_SIZE !== 32) throw new Error('Invalid KEM_SHARED_SECRET_SIZE');
  });

  await runner.runTest('Fragment size constants', async () => {
    if (DEFAULT_FRAGMENT_SIZE !== 1024) throw new Error('Invalid DEFAULT_FRAGMENT_SIZE');
    if (MAX_FRAGMENT_SIZE !== 65536) throw new Error('Invalid MAX_FRAGMENT_SIZE');
    if (MIN_FRAGMENT_SIZE !== 256) throw new Error('Invalid MIN_FRAGMENT_SIZE');
  });

  runner.endSuite();
}

async function testUtilities(runner: TestRunner): Promise<void> {
  runner.startSuite('Utility Functions');

  await runner.runTest('secureRandom generates correct size', async () => {
    const random = await secureRandom(32);
    if (random.length !== 32) {
      throw new Error('secureRandom returned wrong size');
    }
  });

  await runner.runTest('constantTimeEqual works correctly', async () => {
    const a = new Uint8Array([1, 2, 3, 4]);
    const b = new Uint8Array([1, 2, 3, 4]);
    const c = new Uint8Array([1, 2, 3, 5]);

    if (!constantTimeEqual(a, b)) {
      throw new Error('constantTimeEqual failed for equal arrays');
    }
    if (constantTimeEqual(a, c)) {
      throw new Error('constantTimeEqual failed for different arrays');
    }
  });

  await runner.runTest('secureZero clears memory', async () => {
    const data = new Uint8Array([1, 2, 3, 4, 5]);
    secureZero(data);
    if (data.some(byte => byte !== 0)) {
      throw new Error('secureZero failed to clear memory');
    }
  });

  await runner.runTest('toHex and fromHex round trip', async () => {
    const original = new Uint8Array([0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0]);
    const hex = toHex(original);
    const restored = fromHex(hex);

    if (!constantTimeEqual(original, restored)) {
      throw new Error('toHex/fromHex round trip failed');
    }
  });

  await runner.runTest('xorBytes works correctly', async () => {
    const a = new Uint8Array([0xff, 0x00, 0xaa, 0x55]);
    const b = new Uint8Array([0x00, 0xff, 0x55, 0xaa]);
    const expected = new Uint8Array([0xff, 0xff, 0xff, 0xff]);

    const result = xorBytes(a, b);
    if (!constantTimeEqual(result, expected)) {
      throw new Error('xorBytes produced incorrect result');
    }
  });

  await runner.runTest('validateSize works correctly', async () => {
    const data = new Uint8Array(32);
    validateSize(data, 32, 'test data');

    try {
      validateSize(data, 64, 'test data');
      throw new Error('validateSize should have thrown');
    } catch (e) {
      if (!(e instanceof Error) || !e.message.includes('Invalid test data size')) {
        throw new Error('validateSize threw wrong error');
      }
    }
  });

  await runner.runTest('copyBytes works correctly', async () => {
    const source = new Uint8Array([1, 2, 3, 4, 5]);
    const copy = copyBytes(source);

    if (!constantTimeEqual(source, copy)) {
      throw new Error('copyBytes failed');
    }

    // Ensure it's a real copy
    source[0] = 99;
    if (copy[0] === 99) {
      throw new Error('copyBytes created reference, not copy');
    }
  });

  await runner.runTest('concatBytes works correctly', async () => {
    const a = new Uint8Array([1, 2]);
    const b = new Uint8Array([3, 4]);
    const expected = new Uint8Array([1, 2, 3, 4]);

    const result = concatBytes(a, b);
    if (!constantTimeEqual(result, expected)) {
      throw new Error('concatBytes failed');
    }
  });

  await runner.runTest('measureTime measures correctly', async () => {
    const { result, timeMs } = await measureTime(async () => {
      await sleep(10);
      return 'test';
    });

    if (timeMs < 8 || timeMs > 50) {
      throw new Error(`measureTime inaccurate: ${timeMs}ms`);
    }
    if (result !== 'test') {
      throw new Error('measureTime lost return value');
    }
  });

  await runner.runTest('hasWebCrypto returns boolean', async () => {
    const result = hasWebCrypto();
    if (typeof result !== 'boolean') {
      throw new Error('hasWebCrypto should return boolean');
    }
  });

  await runner.runTest('getSystemCapabilities returns object', async () => {
    const caps = getSystemCapabilities();
    if (typeof caps !== 'object' || !caps.hasWebCrypto || !caps.platform) {
      throw new Error('getSystemCapabilities returned invalid object');
    }
  });

  runner.endSuite();
}

async function testHashOperations(runner: TestRunner): Promise<void> {
  runner.startSuite('Hash Operations');

  await runner.runTest('computeHash produces correct size', async () => {
    const data = new TextEncoder().encode('test');
    const hash = await computeHash(data);
    if (hash.length !== HASH_SIZE) {
      throw new Error(`Hash size incorrect: ${hash.length} !== ${HASH_SIZE}`);
    }
  });

  await runner.runTest('computeHash is deterministic', async () => {
    const data = new TextEncoder().encode('test');
    const hash1 = await computeHash(data);
    const hash2 = await computeHash(data);

    if (!constantTimeEqual(hash1, hash2)) {
      throw new Error('Hash is not deterministic');
    }
  });

  await runner.runTest('computeHashWithSalt works', async () => {
    const data = new TextEncoder().encode('test');
    const salt = await secureRandom(16);
    const hash = await computeHashWithSalt(data, salt);

    if (hash.length !== HASH_SIZE) {
      throw new Error('Salted hash size incorrect');
    }
  });

  await runner.runTest('computeHmac works', async () => {
    const data = new TextEncoder().encode('test');
    const key = await secureRandom(32);
    const hmac = await computeHmac(data, key);

    if (hmac.length !== HASH_SIZE) {
      throw new Error('HMAC size incorrect');
    }
  });

  await runner.runTest('batchHash processes multiple inputs', async () => {
    const inputs = [
      new TextEncoder().encode('test1'),
      new TextEncoder().encode('test2'),
      new TextEncoder().encode('test3')
    ];

    const hashes = await batchHash(inputs);
    if (hashes.length !== inputs.length) {
      throw new Error('Batch hash count mismatch');
    }

    for (const hash of hashes) {
      if (hash.length !== HASH_SIZE) {
        throw new Error('Batch hash size incorrect');
      }
    }
  });

  await runner.runTest('computeMerkleRoot works', async () => {
    const leaves = [
      await secureRandom(32),
      await secureRandom(32),
      await secureRandom(32),
      await secureRandom(32)
    ];

    const root = await computeMerkleRoot(leaves);
    if (root.length !== HASH_SIZE) {
      throw new Error('Merkle root size incorrect');
    }
  });

  await runner.runTest('deriveKey works', async () => {
    const password = new TextEncoder().encode('test password');
    const salt = await secureRandom(16);
    const key = await deriveKey(password, salt, 10000, 32);

    if (key.length !== 32) {
      throw new Error('Derived key size incorrect');
    }
  });

  await runner.runTest('hash chain verification works', async () => {
    const initialValue = await secureRandom(32);
    const chain = await computeHashChain(initialValue, 5);

    if (chain.length !== 6) {
      // initial + 5 iterations
      throw new Error('Hash chain length incorrect');
    }

    const isValid = await verifyHashChain(chain, initialValue);
    if (!isValid) {
      throw new Error('Hash chain verification failed');
    }
  });

  runner.endSuite();
}

async function testKeyPairOperations(runner: TestRunner): Promise<void> {
  runner.startSuite('Key Pair Operations');

  await runner.runTest('generateKeyPair produces correct sizes', async () => {
    const keyPair = await generateKeyPair();
    if (keyPair.privateKey.length !== PRIVATE_KEY_SIZE) {
      throw new Error('Private key size incorrect');
    }
    if (keyPair.publicKey.length !== PUBLIC_KEY_SIZE) {
      throw new Error('Public key size incorrect');
    }
  });

  await runner.runTest('generateKeyPairFromSeed is deterministic', async () => {
    const seed = await secureRandom(32);
    const keyPair1 = await generateKeyPairFromSeed(seed);
    const keyPair2 = await generateKeyPairFromSeed(seed);

    if (!constantTimeEqual(keyPair1.privateKey, keyPair2.privateKey)) {
      throw new Error('Seed-based generation not deterministic');
    }
  });

  await runner.runTest('validateKeyPair works', async () => {
    const keyPair = await generateKeyPair();
    const isValid = await validateKeyPair(keyPair);
    if (!isValid) {
      throw new Error('Valid key pair failed validation');
    }
  });

  await runner.runTest('batchGenerateKeyPairs works', async () => {
    const keyPairs = await batchGenerateKeyPairs(5);
    if (keyPairs.length !== 5) {
      throw new Error('Batch generation count incorrect');
    }

    for (const keyPair of keyPairs) {
      if (keyPair.privateKey.length !== PRIVATE_KEY_SIZE) {
        throw new Error('Batch key pair private key size incorrect');
      }
    }
  });

  await runner.runTest('deriveChildKeyPair works', async () => {
    const parentKeyPair = await generateKeyPair();
    const childKeyPair = await deriveChildKeyPair(parentKeyPair.privateKey, 0);

    if (childKeyPair.privateKey.length !== PRIVATE_KEY_SIZE) {
      throw new Error('Child key pair size incorrect');
    }

    const isValid = await validateKeyPair(childKeyPair);
    if (!isValid) {
      throw new Error('Child key pair validation failed');
    }
  });

  await runner.runTest('generateHDWallet works', async () => {
    const seed = await secureRandom(64);
    const wallet = await generateHDWallet(seed, 3);

    if (wallet.length !== 3) {
      throw new Error('HD wallet size incorrect');
    }

    for (const keyPair of wallet) {
      const isValid = await validateKeyPair(keyPair);
      if (!isValid) {
        throw new Error('HD wallet key pair validation failed');
      }
    }
  });

  await runner.runTest('deriveKeyPairFromPassword works', async () => {
    const password = 'test password';
    const salt = await secureRandom(32);
    const keyPair = await deriveKeyPairFromPassword(password, salt);

    if (keyPair.privateKey.length !== PRIVATE_KEY_SIZE) {
      throw new Error('Password-derived key pair size incorrect');
    }

    const isValid = await validateKeyPair(keyPair);
    if (!isValid) {
      throw new Error('Password-derived key pair validation failed');
    }
  });

  await runner.runTest('key pair serialization round trip', async () => {
    const keyPair = await generateKeyPair();
    const serialized = serializeKeyPair(keyPair);
    const deserialized = deserializeKeyPair(serialized);

    if (!constantTimeEqual(keyPair.privateKey, deserialized.privateKey)) {
      throw new Error('Key pair serialization round trip failed');
    }
  });

  await runner.runTest('secureEraseKeyPair works', async () => {
    const keyPair = await generateKeyPair();
    const backup = backupKeyPair(keyPair);

    secureEraseKeyPair(keyPair);

    const isErased =
      keyPair.privateKey.every(byte => byte === 0) && keyPair.publicKey.every(byte => byte === 0);
    if (!isErased) {
      throw new Error('Key pair not properly erased');
    }

    // Verify backup is intact
    const backupValid = await validateKeyPair(backup);
    if (!backupValid) {
      throw new Error('Backup key pair corrupted');
    }
  });

  runner.endSuite();
}

async function testKEMOperations(runner: TestRunner): Promise<void> {
  runner.startSuite('KEM Operations');

  await runner.runTest('kemKeyGen produces correct sizes', async () => {
    const kemKeyPair = await kemKeyGen();
    if (kemKeyPair.secretKey.length !== KEM_SECRET_KEY_SIZE) {
      throw new Error('KEM secret key size incorrect');
    }
    if (kemKeyPair.publicKey.length !== KEM_PUBLIC_KEY_SIZE) {
      throw new Error('KEM public key size incorrect');
    }
  });

  await runner.runTest('KEM encapsulation/decapsulation works', async () => {
    const kemKeyPair = await kemKeyGen();
    const { ciphertext, sharedSecret } = await kemEncapsulate(kemKeyPair.publicKey);

    if (ciphertext.length !== KEM_CIPHERTEXT_SIZE) {
      throw new Error('KEM ciphertext size incorrect');
    }
    if (sharedSecret.length !== KEM_SHARED_SECRET_SIZE) {
      throw new Error('KEM shared secret size incorrect');
    }

    const decapsulatedSecret = await kemDecapsulate(kemKeyPair.secretKey, ciphertext);
    if (!constantTimeEqual(sharedSecret, decapsulatedSecret)) {
      throw new Error('KEM shared secrets do not match');
    }
  });

  await runner.runTest('validateKEMKeyPair works', async () => {
    const kemKeyPair = await kemKeyGen();
    const isValid = await validateKEMKeyPair(kemKeyPair);
    if (!isValid) {
      throw new Error('Valid KEM key pair failed validation');
    }
  });

  await runner.runTest('batchKEMKeyGen works', async () => {
    const kemKeyPairs = await batchKEMKeyGen(3);
    if (kemKeyPairs.length !== 3) {
      throw new Error('Batch KEM generation count incorrect');
    }

    for (const kemKeyPair of kemKeyPairs) {
      const isValid = await validateKEMKeyPair(kemKeyPair);
      if (!isValid) {
        throw new Error('Batch KEM key pair validation failed');
      }
    }
  });

  await runner.runTest('batch KEM operations work', async () => {
    const kemKeyPairs = await batchKEMKeyGen(3);
    const publicKeys = kemKeyPairs.map(kp => kp.publicKey);
    const secretKeys = kemKeyPairs.map(kp => kp.secretKey);

    const encapResults = await batchKEMEncapsulate(publicKeys);
    if (encapResults.length !== 3) {
      throw new Error('Batch encapsulation count incorrect');
    }

    // Prepare operations for batch decapsulation
    const operations = secretKeys.map((secretKey, i) => ({
      secretKey,
      ciphertext: encapResults[i]!.ciphertext
    }));
    const decapSecrets = await batchKEMDecapsulate(operations);

    for (let i = 0; i < 3; i++) {
      if (!constantTimeEqual(encapResults[i]!.sharedSecret, decapSecrets[i]!)) {
        throw new Error(`Batch KEM operation ${i} failed`);
      }
    }
  });

  // Note: testKEMOperations is already called as part of the test suite
  // No need for a separate test here

  await runner.runTest('KEM key pair serialization works', async () => {
    const kemKeyPair = await kemKeyGen();
    const serialized = serializeKEMKeyPair(kemKeyPair);
    const deserialized = deserializeKEMKeyPair(serialized);

    if (!constantTimeEqual(kemKeyPair.secretKey, deserialized.secretKey)) {
      throw new Error('KEM key pair serialization failed');
    }
  });

  await runner.runTest('KEM secure erasure works', async () => {
    const kemKeyPair = await kemKeyGen();
    const backup = backupKEMKeyPair(kemKeyPair);

    secureEraseKEMKeyPair(kemKeyPair);

    const isErased =
      kemKeyPair.secretKey.every(byte => byte === 0) &&
      kemKeyPair.publicKey.every(byte => byte === 0);
    if (!isErased) {
      throw new Error('KEM key pair not properly erased');
    }

    const backupValid = await validateKEMKeyPair(backup);
    if (!backupValid) {
      throw new Error('KEM backup corrupted');
    }
  });

  runner.endSuite();
}

async function testFragmentation(runner: TestRunner): Promise<void> {
  runner.startSuite('Data Fragmentation');

  await runner.runTest('fragmentData works correctly', async () => {
    const data = new Uint8Array(2048);
    crypto.getRandomValues(data);

    const result = await fragmentData(data);
    if (result.metadata.originalSize !== data.length) {
      throw new Error('Fragment metadata incorrect');
    }
    if (result.fragments.length === 0) {
      throw new Error('No fragments created');
    }
  });

  await runner.runTest('reconstructData works correctly', async () => {
    const data = new Uint8Array(1024);
    crypto.getRandomValues(data);

    const fragResult = await fragmentData(data);
    const reconResult = await reconstructData(fragResult.fragments);

    if (!reconResult.isComplete) {
      throw new Error('Reconstruction not complete');
    }
    if (!constantTimeEqual(data, reconResult.data)) {
      throw new Error('Reconstructed data does not match original');
    }
  });

  await runner.runTest('validateFragment works', async () => {
    const data = new Uint8Array(512);
    crypto.getRandomValues(data);

    const fragResult = await fragmentData(data);
    for (const fragment of fragResult.fragments) {
      const isValid = await validateFragment(fragment);
      if (!isValid) {
        throw new Error('Valid fragment failed validation');
      }
    }
  });

  await runner.runTest('parallel fragmentation works', async () => {
    const data = new Uint8Array(4096);
    crypto.getRandomValues(data);

    const result = await parallelFragmentation([data]);
    const reconResult = await reconstructData(result[0]!.fragments);

    if (!constantTimeEqual(data, reconResult.data)) {
      throw new Error('Parallel fragmentation failed');
    }
  });

  await runner.runTest('parallel reconstruction works', async () => {
    const data = new Uint8Array(4096);
    crypto.getRandomValues(data);

    const fragResult = await fragmentData(data);
    const reconResult = await parallelReconstruction([fragResult.fragments]);

    if (!constantTimeEqual(data, reconResult[0]!.data)) {
      throw new Error('Parallel reconstruction failed');
    }
  });

  await runner.runTest('fragment serialization works', async () => {
    const data = new Uint8Array(1024);
    crypto.getRandomValues(data);

    const fragResult = await fragmentData(data);
    const serialized = serializeFragments(fragResult.fragments);
    const deserialized = deserializeFragments(serialized);

    if (deserialized.length !== fragResult.fragments.length) {
      throw new Error('Fragment serialization count mismatch');
    }
  });

  await runner.runTest('fragment compression works', async () => {
    // Create data with patterns for better compression
    const data = new Uint8Array(2048);
    for (let i = 0; i < data.length; i++) {
      data[i] = i % 256;
    }

    const fragResult = await fragmentData(data);
    const compressed = compressFragments(fragResult.fragments);
    const decompressed = decompressFragments(compressed);

    if (decompressed.length !== fragResult.fragments.length) {
      throw new Error('Fragment compression/decompression failed');
    }
  });

  await runner.runTest('mobile optimization functions work', async () => {
    const latency = estimateMobileLatency(1024);
    if (typeof latency !== 'number' || latency <= 0) {
      throw new Error('Mobile latency estimation failed');
    }

    const optimalSize = getOptimalFragmentSize(1024);
    if (typeof optimalSize !== 'number' || optimalSize <= 0) {
      throw new Error('Optimal fragment size calculation failed');
    }
  });

  runner.endSuite();
}

async function testPerformance(runner: TestRunner): Promise<void> {
  runner.startSuite('Performance Functions');

  await runner.runTest('benchmarkHash works', async () => {
    const results = await benchmarkHash();
    if (!results.throughputOpsPerSec || results.throughputOpsPerSec <= 0) {
      throw new Error('Hash benchmark failed');
    }
  });

  await runner.runTest('benchmarkKeyPairGeneration works', async () => {
    const results = await benchmarkKeyPairGeneration();
    if (!results.throughputOpsPerSec || results.throughputOpsPerSec <= 0) {
      throw new Error('Key pair benchmark failed');
    }
  });

  await runner.runTest('benchmarkKEM works', async () => {
    const results = await benchmarkKEM();
    if (!results.throughputOpsPerSec || results.throughputOpsPerSec <= 0) {
      throw new Error('KEM benchmark failed');
    }
  });

  await runner.runTest('benchmarkFragmentation works', async () => {
    const results = await benchmarkFragmentation();
    if (!results.throughputOpsPerSec || results.throughputOpsPerSec <= 0) {
      throw new Error('Fragmentation benchmark failed');
    }
  });

  await runner.runTest('runBenchmarkSuite works', async () => {
    const results = await runBenchmarkSuite();
    if (results.summary.totalTimeMs <= 0) {
      throw new Error('Benchmark suite failed');
    }
  });

  await runner.runTest('monitorMemoryUsage works', async () => {
    const results = await monitorMemoryUsage(async () => {
      const data = new Uint8Array(1024);
      return data;
    });

    if (typeof results.memoryMetrics.beforeMB !== 'number') {
      throw new Error('Memory monitoring failed');
    }
  });

  await runner.runTest('profileCPUUsage works', async () => {
    const results = await profileCPUUsage(async () => {
      await sleep(10);
    });

    if (typeof results.cpuMetrics.averageUsagePercent !== 'number') {
      throw new Error('CPU profiling failed');
    }
  });

  await runner.runTest('estimateMobilePerformance works', async () => {
    const testMetrics = {
      operation: 'test',
      executionTimeMs: 100,
      throughputOpsPerSec: 10
    };
    const results = estimateMobilePerformance(testMetrics);
    if (results.executionTimeMs <= testMetrics.executionTimeMs) {
      throw new Error('Mobile performance estimation failed');
    }
  });

  await runner.runTest('generatePerformanceReport works', async () => {
    const testResults = await runBenchmarkSuite();
    const report = generatePerformanceReport(testResults);
    if (!report.includes('TOPAY-Z512 Performance Report')) {
      throw new Error('Performance report generation failed');
    }
  });

  runner.endSuite();
}

async function testIntegration(runner: TestRunner): Promise<void> {
  runner.startSuite('Integration Tests');

  await runner.runTest('Full workflow: key generation to fragmentation', async () => {
    // Generate key pair
    // const keyPair = await generateKeyPair(); // Commented out unused variable

    // Hash some data
    const data = new TextEncoder().encode('Integration test data');
    // const hash = await computeHash(data); // Commented out unused variable

    // Generate KEM key pair and perform operations
    const kemKeyPair = await kemKeyGen();
    const { ciphertext, sharedSecret } = await kemEncapsulate(kemKeyPair.publicKey);
    const decapSecret = await kemDecapsulate(kemKeyPair.secretKey, ciphertext);

    if (!constantTimeEqual(sharedSecret, decapSecret)) {
      throw new Error('KEM operation failed in integration test');
    }

    // Fragment the data
    const fragResult = await fragmentData(data);
    const reconResult = await reconstructData(fragResult.fragments);

    if (!constantTimeEqual(data, reconResult.data)) {
      throw new Error('Fragmentation failed in integration test');
    }
  });

  await runner.runTest('Secure key exchange simulation', async () => {
    // Alice and Bob generate key pairs
    // const aliceKeyPair = await kemKeyGen(); // Commented out unused variable
    const bobKeyPair = await kemKeyGen();

    // Alice encapsulates using Bob's public key
    const aliceEncap = await kemEncapsulate(bobKeyPair.publicKey);

    // Bob decapsulates using his secret key
    const bobSecret = await kemDecapsulate(bobKeyPair.secretKey, aliceEncap.ciphertext);

    if (!constantTimeEqual(aliceEncap.sharedSecret, bobSecret)) {
      throw new Error('Key exchange simulation failed');
    }

    // Use shared secret to derive encryption key
    const keyMaterial = new TextEncoder().encode('shared');
    const encryptionKey = await deriveKey(keyMaterial, aliceEncap.sharedSecret, 10000, 32);

    if (encryptionKey.length !== 32) {
      throw new Error('Encryption key derivation failed');
    }
  });

  await runner.runTest('Large data processing pipeline', async () => {
    // Create large test data
    const largeData = new Uint8Array(10240);
    crypto.getRandomValues(largeData);

    // Hash the data
    const dataHash = await computeHash(largeData);

    // Fragment the data
    const fragResult = await parallelFragmentation([largeData]);

    if (!fragResult || fragResult.length === 0 || !fragResult[0]) {
      throw new Error('Fragmentation failed');
    }

    // Compress fragments
    const compressed = compressFragments(fragResult[0].fragments);

    // Serialize fragments
    const serialized = serializeFragments(compressed);

    // Deserialize fragments
    const deserialized = deserializeFragments(serialized);

    // Decompress fragments
    const decompressed = decompressFragments(deserialized);

    // Reconstruct data
    const reconResult = await parallelReconstruction([decompressed]);

    if (!reconResult || reconResult.length === 0 || !reconResult[0]) {
      throw new Error('Reconstruction failed');
    }

    if (!reconResult[0].isComplete) {
      throw new Error('Large data pipeline reconstruction failed');
    }

    // Verify data integrity
    const reconHash = await computeHash(reconResult[0].data);
    if (!constantTimeEqual(dataHash, reconHash)) {
      throw new Error('Large data pipeline integrity check failed');
    }
  });

  await runner.runTest('Batch operations integration', async () => {
    // Generate multiple key pairs
    // const keyPairs = await batchGenerateKeyPairs(5); // Commented out unused variable

    // Generate multiple KEM key pairs
    const kemKeyPairs = await batchKEMKeyGen(5);

    // Perform batch KEM operations
    const publicKeys = kemKeyPairs.map(kp => kp.publicKey);
    const encapResults = await batchKEMEncapsulate(publicKeys);

    // Prepare operations for batch decapsulation
    const operations = kemKeyPairs.map((kp, i) => ({
      secretKey: kp.secretKey,
      ciphertext: encapResults[i]!.ciphertext
    }));
    const decapSecrets = await batchKEMDecapsulate(operations);

    // Verify all operations succeeded
    for (let i = 0; i < 5; i++) {
      if (!constantTimeEqual(encapResults[i]!.sharedSecret, decapSecrets[i]!)) {
        throw new Error(`Batch operation ${i} failed`);
      }
    }

    // Batch hash multiple data sets
    const dataSets = Array.from({ length: 5 }, (_, i) =>
      new TextEncoder().encode(`test data ${i}`)
    );
    const hashes = await batchHash(dataSets);

    if (hashes.length !== 5) {
      throw new Error('Batch hash operation failed');
    }
  });

  runner.endSuite();
}

// Run all tests
if (require.main === module) {
  runAllTests().catch(console.error);
}

export {
  runAllTests,
  testConstants,
  testUtilities,
  testHashOperations,
  testKeyPairOperations,
  testKEMOperations,
  testFragmentation,
  testPerformance,
  testIntegration,
  TestRunner
};
