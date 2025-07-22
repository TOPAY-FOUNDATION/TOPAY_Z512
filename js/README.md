# TOPAY-Z512 JavaScript/TypeScript Implementation

This is the JavaScript/TypeScript implementation of the TOPAY-Z512 cryptographic library, providing a 512-bit post-quantum cryptography solution with support for Key Encapsulation Mechanism (KEM) and cryptographic hashing.

## Features

- 512-bit cryptographic hash function based on SHA3-512
- Key Encapsulation Mechanism (KEM) based on Learning With Errors (LWE)
- Fragmentation support for better performance on resource-constrained devices
- TypeScript support with full type definitions

## Installation

```bash
npm install topayz512
# or
yarn add topayz512
```

## Usage

### Hash Function

```typescript
import { Hash, hash, hashCombine, newWithTime, hashWithTime } from 'topayz512';

// Create a new hash from data
const data = 'Hello, TOPAY-Z512!';
const hashValue = Hash.new(data);

// Get the hash bytes
const bytes = hashValue.getBytes();

// Convert hash to hex string
const hexString = hashValue.toHex();

// Create hash from hex string
const hashFromHex = Hash.fromHex(hexString);

// Combine two pieces of data into a single hash
const data1 = 'TOPAY';
const data2 = 'Z512';
const combinedHash = Hash.combine(data1, data2);

// Convenience functions
const hashBytes = hash(data);
const combinedHashBytes = hashCombine(data1, data2);

// Binary data
const binaryData = new Uint8Array([0, 1, 2, 3, 4, 5]);
const binaryHash = Hash.new(binaryData);

// Time-based hashing (useful for generating random-like hashes)
const timeHash = Hash.newWithTime();
const timeHashObj = newWithTime();
const timeHashBytes = hashWithTime();
```

### Examples

See the `examples` directory for more detailed examples:

```bash
# Compile TypeScript
npm run build

# Run example
node dist/examples/hash_example.js
```

## Development

### Building

```bash
npm run build
```

### Testing

```bash
npm test
```

### Linting

```bash
npm run lint
```

## Documentation

The API documentation is available in the `docs` directory at the root of the repository.

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](../LICENSE) file for details.
