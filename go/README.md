# TOPAY-Z512 Go Implementation

This is the Go implementation of the TOPAY-Z512 cryptographic library, providing a 512-bit post-quantum cryptography solution with support for Key Encapsulation Mechanism (KEM) and cryptographic hashing.

## Features

- 512-bit cryptographic hash function based on SHA3-512
- Key Encapsulation Mechanism (KEM) based on Learning With Errors (LWE)
- Fragmentation support for better performance on resource-constrained devices

## Installation

```bash
go get github.com/topay-foundation/topayz512
```

## Usage

### Hash Function

```go
import (
    "fmt"
    "github.com/topay-foundation/topayz512/pkg/hash"
)

// Create a new hash from data
data := []byte("Hello, TOPAY-Z512!")
hashValue := hash.New(data)

// Get the hash bytes
bytes := hashValue.Bytes()

// Convert hash to hex string
hexString := hashValue.String()

// Create hash from hex string
hashFromHex, err := hash.FromHex(hexString)
if err != nil {
    fmt.Printf("Error: %v\n", err)
}

// Combine two pieces of data into a single hash
data1 := []byte("TOPAY")
data2 := []byte("Z512")
combinedHash := hash.Combine(data1, data2)

// Convenience functions
hashBytes := hash.Sum512(data)
combinedHashBytes := hash.SumCombine(data1, data2)

// Time-based hashing (useful for generating random-like hashes)
timeHash := hash.NewWithTime()
timeHashBytes := hash.Sum512WithTime()
```

## Running the Examples

The TOPAY-Z512 Go implementation includes several examples that demonstrate the functionality of the library. You can run these examples using the provided `main.go` file.

### Available Examples

1. **Hash Example** - Demonstrates the TOPAY-Z512 hashing functionality, including basic hashing, hash combination, and hex conversion.
2. **Key Pair Example** - Shows how to generate and use key pairs, including public key derivation and hex conversion.
3. **Private to Public Key Conversion Example** - Demonstrates how to convert a private key to a public key.

### Example Usage Instructions

You can run the examples in two ways:

#### Option 1: Interactive Menu

Run the following command to display an interactive menu that allows you to select which example to run:

```bash
go run main.go
```

#### Option 2: Direct Execution

You can also run a specific example directly by providing the example number as a command-line argument:

```bash
go run main.go 1  # Run the Hash Example
go run main.go 2  # Run the Key Pair Example
go run main.go 3  # Run the Private to Public Key Conversion Example
```

### Examples

See the `examples` directory for more detailed examples:

```bash
go run examples/hash_example.go
```

## Testing

Run the tests with:

```bash
go test ./...
```

## Benchmarking

Run the benchmarks with:

```bash
go test -bench=. ./...
```

## Documentation

Generate and view the documentation with:

```bash
go doc -all github.com/topay-foundation/topayz512/pkg/hash
```

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](../LICENSE) file for details.
