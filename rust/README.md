# topayz512 - Rust Implementation

This is the Rust implementation of the topayz512 library, a 512-bit post-quantum cryptography library based on Learning With Errors (LWE). The library provides a Key Encapsulation Mechanism (KEM) with ≥512-bit classical security and ~256-bit quantum resistance.

## Features

- **High Security**: Achieves ≥512-bit classical security (~256-bit quantum resistance) using lattice-based LWE parameters.
- **Performance**: Optimized implementation with optional fragmentation support for mobile and embedded devices.
- **No-std Support**: Can be used in environments without the standard library.
- **Fragmentation**: Optional feature to split large operations into smaller workloads.

## Installation

Add the following to your `Cargo.toml`:

```toml
[dependencies]
topayz512 = "0.1.0"
```

To enable the fragmentation feature:

```toml
[dependencies]
topayz512 = { version = "0.1.0", features = ["fragmentation"] }
```

## Usage

### Basic Key Exchange

```rust
use topayz512::{keygen, encapsulate, decapsulate};

// Alice generates a key pair
let (alice_pk, alice_sk) = keygen()?;

// Bob encapsulates a shared secret using Alice's public key
let (ciphertext, bob_shared_secret) = encapsulate(&alice_pk)?;

// Alice decapsulates the shared secret using her secret key and Bob's ciphertext
let alice_shared_secret = decapsulate(&alice_sk, &ciphertext)?;

// Both parties now have the same shared secret
assert_eq!(alice_shared_secret, bob_shared_secret);
```

### Fragmented Key Exchange

When the fragmentation feature is enabled:

```rust
use topayz512::fragmented::{keygen, encapsulate, decapsulate};

// Number of fragments to use
let num_fragments = 4;

// Alice generates a fragmented key pair
let (alice_pk, alice_sk) = keygen(num_fragments)?;

// Bob encapsulates a shared secret using Alice's fragmented public key
let (ciphertext, bob_shared_secret) = encapsulate(&alice_pk)?;

// Alice decapsulates the shared secret using her fragmented secret key and Bob's ciphertext
let alice_shared_secret = decapsulate(&alice_sk, &ciphertext)?;

// Both parties now have the same shared secret
assert_eq!(alice_shared_secret, bob_shared_secret);
```

## Examples

The library includes several examples:

- `key_exchange.rs`: Demonstrates basic key exchange.
- `fragmented_key_exchange.rs`: Demonstrates fragmented key exchange (requires the `fragmentation` feature).

To run an example:

```bash
cargo run --example key_exchange
```

Or with the fragmentation feature:

```bash
cargo run --features fragmentation --example fragmented_key_exchange
```

## Benchmarks

To run the benchmarks:

```bash
cargo bench
```

To run the benchmarks with the fragmentation feature:

```bash
cargo bench --features fragmentation
```

## Security Parameters

The library uses the following security parameters:

- **N**: 1024 (dimension of the LWE problem)
- **Q**: 65537 (modulus for the LWE problem, 2^16 + 1)
- **σ**: 3.2 (standard deviation for the error distribution)

These parameters are designed to provide ≥512-bit classical security and ~256-bit quantum resistance.

## License

This library is licensed under the MIT License - see the LICENSE file for details.
