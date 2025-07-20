# TOPAY-Z512 Rust Implementation

This is the Rust implementation of the TOPAY-Z512 cryptographic library, providing a 512-bit post-quantum cryptography solution with support for Key Encapsulation Mechanism (KEM) and cryptographic hashing.

## Features

- 512-bit cryptographic hash function based on SHA3-512
- Key Encapsulation Mechanism (KEM) based on Learning With Errors (LWE)
- Fragmentation support for better performance on resource-constrained devices

## Installation

Add this to your `Cargo.toml`:

```toml
[dependencies]
topayz512 = "0.1.0"
```

## Usage

### Hash Function

```rust
use topayz512::{Hash, hash, hash_combine};

// Create a new hash from data
let data = b"Hello, TOPAY-Z512!";
let hash_value = Hash::new(data);

// Get the hash bytes
let bytes = hash_value.as_bytes();

// Convert hash to hex string
let hex_string = hash_value.to_hex();

// Create hash from hex string
let hash_from_hex = Hash::from_hex(&hex_string).unwrap();

// Combine two pieces of data into a single hash
let data1 = b"TOPAY";
let data2 = b"Z512";
let combined_hash = Hash::combine(data1, data2);

// Convenience functions
let hash_bytes = hash(data);
let combined_hash_bytes = hash_combine(data1, data2);
```

### Examples

See the `examples` directory for more detailed examples:

```bash
cargo run --example hash_example
```

## Testing

Run the tests with:

```bash
cargo test
```

## Benchmarking

Run the benchmarks with:

```bash
cargo bench
```

## Documentation

Generate and view the documentation with:

```bash
cargo doc --open
```

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](../LICENSE) file for details.
