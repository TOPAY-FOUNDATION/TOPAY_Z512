# TOPAY-Z512 Key Pair

## Overview

The TOPAY-Z512 Key Pair module provides functionality for generating and managing 512-bit cryptographic key pairs. These key pairs can be used for various cryptographic operations, such as digital signatures and key exchange protocols.

## Implementation Details

### Key Sizes

- **Private Key**: 512 bits (64 bytes)
- **Public Key**: 512 bits (64 bytes)

### Key Generation

The key generation process follows these steps:

1. Generate a random 512-bit private key using a cryptographically secure random number generator.
2. Derive the public key from the private key using the SHA3-512 hash function.

### Private to Public Key Conversion

The library provides convenient functions to derive a public key from a private key:

- Rust: `private_to_public(&private_key)`
- Go: `PrivateToPublic(privateKey)`
- JavaScript/TypeScript: `privateToPublic(privateKey)`

These functions use SHA3-512 to hash the private key, producing a deterministic public key.

### Security Considerations

- The private key should be kept secret and never shared.
- The public key can be freely distributed.
- The key pair generation uses cryptographically secure random number generators specific to each language implementation:
  - Rust: `rand::rngs::OsRng`
  - Go: `crypto/rand`
  - JavaScript/TypeScript: `crypto.getRandomValues`

## API Usage

Refer to the [API Reference](api_reference.md#key-pair-api) for detailed usage examples in Rust, Go, and JavaScript/TypeScript.

## Test Vectors

Test vectors for key pair operations are provided in the `test-vectors/003_keypair.json` file. These test vectors ensure consistent behavior across all language implementations.

## Cross-Language Compatibility

The key pair implementation is designed to be compatible across all supported languages (Rust, Go, and JavaScript/TypeScript). A private key generated in one language can be used to derive the same public key in another language, ensuring interoperability.
