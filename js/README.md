# TOPAY-Z512 JavaScript/TypeScript Implementation

A high-performance, quantum-resistant cryptographic library implementing the TOPAY-Z512 protocol in JavaScript/TypeScript.

## Features

- **Quantum-Resistant Cryptography**: Post-quantum secure key encapsulation mechanism (KEM)
- **High-Performance Hashing**: Optimized SHA-512 based operations
- **Advanced Key Management**: Hierarchical deterministic (HD) wallets and key derivation
- **Data Fragmentation**: Efficient data splitting and reconstruction with mobile optimization
- **Cross-Platform**: Works in Node.js, browsers, and mobile environments
- **TypeScript Support**: Full type safety and IntelliSense support
- **Comprehensive Testing**: Extensive test suite with performance benchmarks

## Installation

```bash
npm install @topay/topayz512
```

## Quick Start

```typescript
import {
  generateKeyPair,
  computeHash,
  kemKeyGen,
  kemEncapsulate,
  kemDecapsulate,
  fragmentData,
  reconstructData
} from '@topay/topayz512';

// Generate a key pair
const keyPair = await generateKeyPair();
console.log('Key pair generated:', {
  privateKey: keyPair.privateKey.length,
  publicKey: keyPair.publicKey.length
});

// Hash some data
const data = new TextEncoder().encode('Hello, TOPAY-Z512!');
const hash = await computeHash(data);
console.log('Hash:', hash);

// KEM operations
const kemKeyPair = await kemKeyGen();
const { ciphertext, sharedSecret } = await kemEncapsulate(kemKeyPair.publicKey);
const decapsulatedSecret = await kemDecapsulate(kemKeyPair.secretKey, ciphertext);
console.log('KEM successful:', sharedSecret.length === decapsulatedSecret.length);

// Data fragmentation
const largeData = new Uint8Array(2048);
crypto.getRandomValues(largeData);

const fragResult = await fragmentData(largeData);
console.log('Fragments created:', fragResult.fragments.length);

const reconResult = await reconstructData(fragResult.fragments);
console.log('Reconstruction successful:', reconResult.isComplete);
```

## API Reference

### Core Types

```typescript
interface KeyPair {
  privateKey: Uint8Array;
  publicKey: Uint8Array;
  timestamp: number;
}

interface KEMKeyPair {
  secretKey: Uint8Array;
  publicKey: Uint8Array;
  timestamp: number;
}

interface Fragment {
  id: number;
  data: Uint8Array;
  checksum: Uint8Array;
  metadata: FragmentMetadata;
}

interface FragmentationResult {
  fragments: Fragment[];
  metadata: FragmentMetadata;
}

interface ReconstructionResult {
  data: Uint8Array;
  isComplete: boolean;
  missingFragments: number[];
  integrityVerified: boolean;
}
```

### Hash Operations

```typescript
// Basic hashing
const hash = await computeHash(data);

// Salted hashing
const saltedHash = await computeHashWithSalt(data, salt);

// HMAC
const hmac = await computeHmac(data, key);

// Batch hashing
const hashes = await batchHash([data1, data2, data3]);

// Merkle root
const root = await computeMerkleRoot(leaves);

// Key derivation
const derivedKey = await deriveKey(password, salt, keyLength);

// Hash chains
const chain = await computeHashChain(initialValue, iterations);
const isValid = await verifyHashChain(chain);
```

### Key Pair Management

```typescript
// Generate random key pair
const keyPair = await generateKeyPair();

// Generate from seed (deterministic)
const seed = await secureRandom(32);
const deterministicKeyPair = await generateKeyPairFromSeed(seed);

// Batch generation
const keyPairs = await batchGenerateKeyPairs(10);

// Validate key pair
const isValid = await validateKeyPair(keyPair);

// Child key derivation
const childKeyPair = await deriveChildKeyPair(parentPrivateKey, index);

// HD wallet generation
const wallet = await generateHDWallet(seed, count);

// Password-based derivation
const passwordKeyPair = await deriveKeyPairFromPassword(password, salt);

// Serialization
const serialized = serializeKeyPair(keyPair);
const deserialized = deserializeKeyPair(serialized);

// Secure erasure
secureEraseKeyPair(keyPair);

// Backup
const backup = backupKeyPair(keyPair);
```

### KEM Operations

```typescript
// Generate KEM key pair
const kemKeyPair = await kemKeyGen();

// Encapsulation
const { ciphertext, sharedSecret } = await kemEncapsulate(publicKey);

// Decapsulation
const decapsulatedSecret = await kemDecapsulate(secretKey, ciphertext);

// Batch operations
const kemKeyPairs = await batchKEMKeyGen(5);
const encapResults = await batchKEMEncapsulate(publicKeys);
const decapSecrets = await batchKEMDecapsulate(secretKeys, ciphertexts);

// Validation
const isValid = await validateKEMKeyPair(kemKeyPair);

// Testing
const testResults = await testKEMOperations();

// Serialization
const serialized = serializeKEMKeyPair(kemKeyPair);
const deserialized = deserializeKEMKeyPair(serialized);

// Secure erasure
secureEraseKEMKeyPair(kemKeyPair);
```

### Data Fragmentation

```typescript
// Basic fragmentation
const fragResult = await fragmentData(data);

// Parallel fragmentation (for large data)
const parallelResult = await parallelFragmentation(data);

// Reconstruction
const reconResult = await reconstructData(fragments);

// Parallel reconstruction
const parallelReconResult = await parallelReconstruction(fragments);

// Fragment validation
const isValid = await validateFragment(fragment);

// Mobile optimization
const latency = estimateMobileLatency(dataSize);
const optimalSize = getOptimalFragmentSize();

// Serialization
const serialized = serializeFragments(fragments);
const deserialized = deserializeFragments(serialized);

// Compression
const compressed = compressFragments(fragments);
const decompressed = decompressFragments(compressed);
```

### Performance & Benchmarking

```typescript
// Individual benchmarks
const hashBench = await benchmarkHashOperations();
const keyPairBench = await benchmarkKeyPairGeneration();
const kemBench = await benchmarkKEMOperations();
const fragBench = await benchmarkFragmentation();

// Full benchmark suite
const fullResults = await runBenchmarkSuite();

// Memory monitoring
const memResults = await monitorMemoryUsage(async () => {
  // Your operation here
});

// CPU profiling
const cpuResults = await profileCPUUsage(async () => {
  // Your operation here
});

// Mobile performance estimation
const mobilePerf = await estimateMobilePerformance();

// Performance report
const report = await generatePerformanceReport();
```

### Utility Functions

```typescript
// Secure random generation
const randomBytes = await secureRandom(32);

// Constant-time comparison
const isEqual = constantTimeEqual(array1, array2);

// Secure memory clearing
secureZero(sensitiveData);

// Hex encoding/decoding
const hex = toHex(bytes);
const bytes = fromHex(hex);

// XOR operation
const result = xorBytes(array1, array2);

// Size validation
validateSize(data, expectedSize, 'data name');

// Array operations
const copy = copyBytes(original);
const combined = concatBytes(array1, array2);

// Timing utilities
const timestamp = timestamp();
await sleep(milliseconds);
const [time, result] = await measureTime(asyncOperation);

// System capabilities
const hasWebCrypto = hasWebCrypto();
const capabilities = getSystemCapabilities();
```

## Examples

The library includes comprehensive examples:

- **Quick Start** (`examples/quick-start.ts`): Basic usage patterns
- **Key Pair Management** (`examples/keypair.ts`): Advanced key operations
- **KEM Operations** (`examples/kem.ts`): Key encapsulation examples
- **Data Fragmentation** (`examples/fragmentation.ts`): Data splitting and reconstruction
- **Performance Benchmarks** (`examples/benchmark.ts`): Performance analysis
- **Interactive Guide** (`examples/interactive-guide.ts`): CLI exploration tool

Run examples:

```bash
# Quick start
npm run example:quick-start

# Key pair operations
npm run example:keypair

# KEM operations
npm run example:kem

# Fragmentation
npm run example:fragmentation

# Benchmarks
npm run example:benchmark

# Interactive guide
npm run example:interactive
```

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run specific test categories
npm run test:hash
npm run test:keypair
npm run test:kem
npm run test:fragmentation
npm run test:performance
```

## Performance

TOPAY-Z512 is optimized for high performance across different environments:

### Typical Performance (Node.js on modern hardware)

- **Hash Operations**: ~50,000 ops/sec
- **Key Pair Generation**: ~1,000 ops/sec
- **KEM Encapsulation**: ~2,000 ops/sec
- **KEM Decapsulation**: ~2,000 ops/sec
- **Data Fragmentation**: ~100 MB/sec
- **Data Reconstruction**: ~150 MB/sec

### Mobile Optimization

The library includes mobile-specific optimizations:

- Adaptive fragment sizing based on device capabilities
- Memory-efficient processing for large datasets
- Battery-aware operation scheduling
- Network latency compensation

## Security Features

- **Quantum Resistance**: Post-quantum cryptographic algorithms
- **Secure Memory Management**: Automatic clearing of sensitive data
- **Constant-Time Operations**: Protection against timing attacks
- **Cryptographic Randomness**: Secure random number generation
- **Input Validation**: Comprehensive parameter validation
- **Error Handling**: Secure error handling without information leakage

## Browser Support

- **Modern Browsers**: Chrome 60+, Firefox 55+, Safari 11+, Edge 79+
- **Node.js**: Version 14.0.0 or higher
- **Mobile**: iOS Safari 11+, Chrome Mobile 60+
- **WebCrypto API**: Required for cryptographic operations

## TypeScript Support

Full TypeScript support with:

- Complete type definitions
- Generic type parameters
- Strict null checks
- IntelliSense support
- Compile-time error checking

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Security

For security issues, please email <security@topay.foundation> instead of using the issue tracker.

## Changelog

### v1.0.0

- Initial release
- Complete TOPAY-Z512 implementation
- Full TypeScript support
- Comprehensive test suite
- Performance optimizations
- Mobile support
- Browser compatibility

## Support

- Documentation: [https://docs.topay.foundation/topayz512](https://docs.topay.foundation/topayz512)
- Issues: [GitHub Issues](https://github.com/TOPAY-FOUNDATION/TOPAY_Z512/issues)
- Community: [Discord](https://discord.gg/topay)
- Email: <support@topayfoundation.com>
