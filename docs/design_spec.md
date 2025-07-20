# TOPAY-Z512 Design Specification

This document outlines the design and specifications for the TOPAY-Z512 cryptographic library.

## Overview

TOPAY-Z512 is a multi-language, open-source 512-bit post-quantum cryptography library with implementations in Rust, Go, and JavaScript/TypeScript. The library provides a standardized Key Encapsulation Mechanism (KEM) based on Learning With Errors (LWE), as well as cryptographic hashing functionality.

## Hash Function Specification

### Algorithm

The TOPAY-Z512 hash function is based on SHA3-512 (Keccak), providing a 512-bit (64-byte) output. This provides a high level of security against collision attacks and preimage attacks.

### Security Properties

- **Collision Resistance**: Finding two different inputs that produce the same hash output is computationally infeasible.
- **Preimage Resistance**: Given a hash output, finding an input that produces that hash is computationally infeasible.
- **Second Preimage Resistance**: Given an input and its hash, finding a different input that produces the same hash is computationally infeasible.

### Implementation Details

- **Hash Size**: 512 bits (64 bytes)
- **Internal State**: SHA3-512 internal state
- **Endianness**: Big-endian byte order for all operations
- **Hex Representation**: Lowercase hexadecimal characters (0-9, a-f)

### Hash Combination

The library provides a method to combine two inputs into a single hash. This is different from simply concatenating the inputs and then hashing. The combination method works as follows:

1. Initialize the hash function state
2. Update the state with the first input
3. Update the state with the second input
4. Finalize the hash

This approach ensures that `Hash::combine(a, b)` is different from `Hash::new(concat(a, b))`.

## KEM Parameters

The Key Encapsulation Mechanism (KEM) is based on the Learning With Errors (LWE) problem with the following parameters:

- **N**: 1024 (lattice dimension)
- **Q**: 65537 (modulus)
- **σ**: 3.2 (standard deviation for error distribution)
- **Secret Length**: 64 bytes

These parameters are designed to provide at least 512-bit classical security (approximately 256-bit quantum resistance).

## Fragmentation

To support high throughput on desktops and acceptable latency on smartphones and embedded devices, the library implements a fragmentation mechanism that breaks large operations into smaller workloads.

### Fragmentation Strategy

1. Split large matrices and vectors into smaller chunks
2. Process each chunk independently
3. Combine the results

This approach allows for better utilization of available resources and can be parallelized on platforms that support it.

## Cross-Language Compatibility

All implementations (Rust, Go, JavaScript/TypeScript) are designed to be compatible with each other. The same input should produce the same output regardless of the implementation used.

Test vectors are provided to ensure compatibility across implementations.
