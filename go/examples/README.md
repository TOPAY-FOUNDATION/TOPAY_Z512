# TOPAY-Z512 Go Examples

This directory contains example code demonstrating how to use the TOPAY-Z512 library in Go.

## Examples

1. **Hash Example** - Demonstrates the hash functionality
2. **Time-based Hash Example** - Demonstrates time-based hash generation
3. **Key Pair Example** - Demonstrates key pair generation and usage
4. **Private to Public Key Conversion Example** - Demonstrates how to convert a private key to a public key

## Running the Examples

You can run the examples in two ways:

### Using the Menu

From the `go` directory, run:

```bash
go run cmd/main.go
```

This will display a menu where you can select which example to run.

### Running a Specific Example

You can also run a specific example by providing the example number as an argument:

```bash
go run cmd/main.go 1  # Run the Hash Example
go run cmd/main.go 2  # Run the Key Pair Example
go run cmd/main.go 3  # Run the Private to Public Key Conversion Example
```

## Example Structure

All examples are in the same package with different function names:

- `Run()` - Hash example
- `RunKeypair()` - Key pair example
- `RunPrivateToPublic()` - Private to public key conversion example

The `main.go` file serves as an entry point to run all the examples.

## Note on Package Structure

The examples are organized in a single package with different function names to avoid naming conflicts. In a real application, you would typically import the TOPAY-Z512 library and use it directly in your code.
