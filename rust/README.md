# TOPAY-Z512 Rust Implementation

A quantum-safe cryptographic library implementing the TOPAY-Z512 protocol for the TOPAY Foundation's blockchain ecosystem.

## Overview

TOPAY-Z512 provides quantum-resistant cryptographic primitives designed for blockchain applications, featuring:

- **512-bit Security**: Enhanced security against quantum attacks
- **Fragmented Architecture**: Optimized for mobile and IoT devices
- **High Performance**: Efficient implementations for production use
- **Cross-Platform**: Compatible across different platforms and languages

## Features

- ✅ **Core Cryptography**: 512-bit hashing and key management
- ✅ **Key Encapsulation Mechanism (KEM)**: Quantum-safe key exchange
- ✅ **Fragmented Processing**: Parallel and distributed computation support
- ✅ **Serialization**: Hex and binary format support
- ✅ **Cross-Language Compatibility**: Consistent behavior across implementations

## Quick Start

Add to your `Cargo.toml`:

```toml
[dependencies]
topayz512 = "0.1.0"
```

### Basic Usage

```rust
use topayz512::{Hash, KeyPair, Kem, Result};

fn main() -> Result<()> {
    // Generate a key pair
    let keypair = KeyPair::generate()?;
    println!("Public key: {}", keypair.public_key().to_hex());
    
    // Create a hash
    let data = b"Hello, TOPAY Foundation!";
    let hash = Hash::new(data);
    println!("Hash: {}", hash.to_hex());
    
    // Key encapsulation
    let (kem_public, kem_secret) = Kem::keygen()?;
    let (ciphertext, shared_secret) = Kem::encapsulate(&kem_public)?;
    let decapsulated = Kem::decapsulate(&kem_secret, &ciphertext)?;
    
    assert_eq!(shared_secret.as_bytes(), decapsulated.as_bytes());
    println!("KEM operation successful!");
    
    Ok(())
}
```

### Fragmentation Support

Enable the fragmentation feature for advanced parallel processing:

```toml
[dependencies]
topayz512 = { version = "0.1.0", features = ["fragmentation"] }
```

```rust
use topayz512::fragment::FragmentEngine;

fn fragmentation_example() -> topayz512::Result<()> {
    let large_data = vec![0x42u8; 4096]; // 4KB of data
    
    // Fragment the data for parallel processing
    let fragments = FragmentEngine::fragment_data(&large_data)?;
    println!("Created {} fragments", fragments.len());
    
    // Reconstruct the original data
    let reconstructed = FragmentEngine::reconstruct_data(&fragments)?;
    assert_eq!(large_data, reconstructed);
    
    Ok(())
}
```

## API Reference

### Core Types

#### `Hash`

512-bit cryptographic hash using SHA3-512.

```rust
// Create from data
let hash = Hash::new(b"data");

// Create from hex string
let hash = Hash::from_hex("deadbeef...")?;

// Combine multiple inputs
let combined = Hash::combine(b"first", b"second");

// Concatenate hashes
let concat = Hash::concat(&[&hash1, &hash2, &hash3]);

// Convert to hex
let hex_string = hash.to_hex();
```

#### `KeyPair`, `PrivateKey`, `PublicKey`

512-bit key pair management.

```rust
// Generate new key pair
let keypair = KeyPair::generate()?;

// Access keys
let private_key = keypair.private_key();
let public_key = keypair.public_key();

// Create from existing private key
let private_key = PrivateKey::from_bytes([0u8; 64]);
let keypair = KeyPair::from_private_key(private_key);

// Serialization
let hex = private_key.to_hex();
let restored = PrivateKey::from_hex(&hex)?;
```

#### `Kem` (Key Encapsulation Mechanism)

Quantum-safe key exchange.

```rust
// Generate KEM key pair
let (public_key, secret_key) = Kem::keygen()?;

// Encapsulate (sender side)
let (ciphertext, shared_secret) = Kem::encapsulate(&public_key)?;

// Decapsulate (receiver side)
let shared_secret = Kem::decapsulate(&secret_key, &ciphertext)?;
```

### Fragmentation (Optional Feature)

#### `FragmentEngine`

Advanced fragmentation for parallel processing.

```rust
// Data fragmentation
let fragments = FragmentEngine::fragment_data(&data)?;
let reconstructed = FragmentEngine::reconstruct_data(&fragments)?;

// KEM fragmentation
let fragmented_kem = FragmentEngine::fragment_kem_encapsulation(&public_key)?;
let results = FragmentEngine::process_fragmented_kem(&fragmented_kem)?;

// Parallel hash computation
let fragmented_hash = FragmentEngine::fragment_hash_operation(&data)?;
let hash = FragmentEngine::parallel_hash_compute(&fragmented_hash)?;

// Mobile performance estimation
let latency = FragmentEngine::estimate_mobile_latency(&fragments);
```

## Examples

The `examples/` directory contains comprehensive examples:

- `hash_example.rs` - Hash operations and combinations
- `keypair_example.rs` - Key generation and management
- `kem_example.rs` - Key encapsulation mechanism
- `fragmentation_example.rs` - Fragmented processing (requires feature)

Run examples with:

```bash
cargo run --example hash_example
cargo run --example kem_example --features fragmentation
```

## Benchmarks

Performance benchmarks are available:

```bash
cargo bench
```

This will run benchmarks for:

- Hash operations
- Key pair generation
- KEM operations
- Fragmentation (if enabled)

## Testing

Run the test suite:

```bash
# Unit tests
cargo test

# Integration tests
cargo test --test integration_tests

# All tests with fragmentation
cargo test --features fragmentation

# Test with release optimizations
cargo test --release
```

## Core Features and Capabilities

### Default Features

- `std` - Standard library support (enabled by default)

### Optional Features

- `fragmentation` - Enable fragmented-block architecture support

### Feature Combinations

```toml
# Minimal (no_std compatible)
topayz512 = { version = "0.1.0", default-features = false }

# With fragmentation
topayz512 = { version = "0.1.0", features = ["fragmentation"] }

# All features
topayz512 = { version = "0.1.0", features = ["std", "fragmentation"] }
```

## Security Considerations

### Quantum Resistance

TOPAY-Z512 is designed to be quantum-safe:

- 512-bit security level provides protection against quantum attacks
- Hash-based constructions resist quantum cryptanalysis
- Key sizes are chosen for post-quantum security

### Implementation Security

- Constant-time operations where possible
- Secure random number generation
- Memory safety through Rust's ownership system
- No unsafe code in core cryptographic operations

### Best Practices

- Always use `KeyPair::generate()` for new keys
- Verify fragment integrity when using fragmentation
- Use secure channels for key exchange
- Regularly update to latest versions

## Performance

### Typical Performance (Release Mode)

- Key pair generation: ~100μs
- Hash operation (1KB): ~10μs
- KEM encapsulation: ~200μs
- KEM decapsulation: ~150μs
- Fragmentation (4KB): ~50μs

Performance varies by platform and hardware. Run benchmarks for your specific environment.

## Cross-Language Compatibility

TOPAY-Z512 maintains compatibility across language implementations:

- Consistent hash outputs for identical inputs
- Compatible key derivation algorithms
- Standardized serialization formats
- Shared test vectors

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/topay-foundation/topay-z512.git
cd topay-z512/rust

# Run tests
cargo test --all-features

# Run benchmarks
cargo bench

# Check formatting
cargo fmt --check

# Run clippy
cargo clippy --all-features
```

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## TOPAY Foundation

TOPAY-Z512 is part of the TOPAY Foundation's mission to create a quantum-safe, human-centered blockchain ecosystem. Learn more at [topay.foundation](https://topay.foundation).

### Related Projects

- TOPAY Blockchain Core
- TOPAY Mobile SDK
- TOPAY Validator Network
- TOPAY Developer Tools

## Changelog

### v0.1.0 (Initial Release)

- Core cryptographic primitives
- Key pair generation and management
- Hash operations with SHA3-512
- Key Encapsulation Mechanism
- Fragmentation support
- Cross-platform compatibility
- Comprehensive test suite
- Performance benchmarks

## Support

For questions, issues, or contributions:

- GitHub Issues: [Report bugs or request features](https://github.com/topay-foundation/topay-z512/issues)
- Documentation: [Full API documentation](https://docs.rs/topayz512)
- Community: [TOPAY Foundation Discord](https://discord.gg/topay)

---

### Built with ❤️ by the TOPAY Foundation team
