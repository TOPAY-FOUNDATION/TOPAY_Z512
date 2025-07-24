# TOPAY-Z512 Go Implementation

This is the Go implementation of the TOPAY-Z512 cryptographic library, providing quantum-safe cryptographic primitives for the TOPAY Foundation blockchain ecosystem.

## Features

- **Quantum-Safe Security**: ≥512-bit classical security (~256-bit quantum resistance)
- **High Performance**: Optimized implementations with batch operations
- **Fragmented Architecture**: Support for parallel processing and mobile devices
- **Cross-Platform**: Compatible with all major operating systems
- **Developer Friendly**: Simple, consistent API design

## Installation

```bash
go get github.com/TOPAY-FOUNDATION/TOPAY_Z512/go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/TOPAY-FOUNDATION/TOPAY_Z512/go"
)

func main() {
    // Generate a key pair
    privateKey, publicKey, err := topayz512.GenerateKeyPair()
    if err != nil {
        panic(err)
    }
    
    // Hash some data
    data := []byte("Hello, TOPAY-Z512!")
    hash := topayz512.ComputeHash(data)
    
    // KEM operations
    kemPublic, kemSecret, err := topayz512.KEMKeyGen()
    if err != nil {
        panic(err)
    }
    
    ciphertext, sharedSecret, err := topayz512.KEMEncapsulate(kemPublic)
    if err != nil {
        panic(err)
    }
    
    decapsulatedSecret, err := topayz512.KEMDecapsulate(kemSecret, ciphertext)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Secrets match: %v\n", bytes.Equal(sharedSecret, decapsulatedSecret))
}
```

## API Reference

### Key Pair Operations

- `GenerateKeyPair() (PrivateKey, PublicKey, error)`
- `DerivePublicKey(privateKey PrivateKey) PublicKey`
- `BatchGenerateKeyPairs(count int) ([]PrivateKey, []PublicKey, error)`

### Hash Operations

- `ComputeHash(data []byte) Hash`
- `HashFromHex(hex string) (Hash, error)`
- `CombineHashes(hashes ...Hash) Hash`

### KEM Operations

- `KEMKeyGen() (KEMPublicKey, KEMSecretKey, error)`
- `KEMEncapsulate(publicKey KEMPublicKey) (Ciphertext, SharedSecret, error)`
- `KEMDecapsulate(secretKey KEMSecretKey, ciphertext Ciphertext) (SharedSecret, error)`
- `BatchKEMKeyGen(count int) ([]KEMPublicKey, []KEMSecretKey, error)`

### Fragmentation Operations (with `fragmentation` build tag)

- `FragmentData(data []byte) ([]Fragment, error)`
- `ReconstructData(fragments []Fragment) ([]byte, error)`
- `EstimateMobileLatency(dataSize int) time.Duration`

## Build Tags

- `fragmentation`: Enables fragmentation support for parallel processing

```bash
go build -tags fragmentation
```

## Performance

The Go implementation is optimized for performance:

- **Hash Operations**: ~6 GB/s throughput
- **Key Generation**: ~50ns per keypair (batch mode)
- **KEM Operations**: ~100ns per operation
- **Fragmentation**: ~1.5 GB/s throughput

## Testing

```bash
# Run all tests
go test ./...

# Run tests with fragmentation
go test -tags fragmentation ./...

# Run benchmarks
go test -bench=. ./...

# Run examples
go run examples/quick_start/main.go
```

## Examples

See the `examples/` directory for comprehensive usage examples:

- `quick_start/` - Basic usage guide
- `hash_example/` - Hash operations
- `kem_example/` - Key encapsulation
- `keypair_example/` - Key pair management
- `fragmentation_example/` - Parallel processing (requires fragmentation tag)
- `performance_benchmark/` - Performance testing

## Security

TOPAY-Z512 provides post-quantum security based on lattice-based cryptography:

- **Classical Security**: ≥512 bits
- **Quantum Security**: ~256 bits
- **Constant-Time Operations**: Protection against timing attacks
- **Secure Memory**: Automatic cleanup of sensitive data

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## TOPAY Foundation

TOPAY-Z512 is part of the TOPAY Foundation's quantum-safe blockchain ecosystem.

- Website: <https://www.topayfoundation.com>
- Email: <contact@topayfoundation.com>

Building the future of quantum-safe finance! 🚀
