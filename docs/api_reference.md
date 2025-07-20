# TOPAY-Z512 API Reference

This document describes the API for the TOPAY-Z512 cryptographic library. The library provides implementations in Rust, Go, and JavaScript/TypeScript.

## Hash API

The Hash API provides a 512-bit cryptographic hash function based on SHA3-512 (Keccak).

### Rust API

```rust
// Create a new hash from data
let hash = topayz512::Hash::new(data);

// Get the hash bytes
let bytes: &[u8; 64] = hash.as_bytes();

// Convert hash to hex string
let hex_string = hash.to_hex();

// Create hash from hex string
let hash = topayz512::Hash::from_hex(hex_string).unwrap();

// Combine two pieces of data into a single hash
let combined_hash = topayz512::Hash::combine(data1, data2);

// Convenience functions
let hash_bytes = topayz512::hash(data);
let combined_hash_bytes = topayz512::hash_combine(data1, data2);
```

### Go API

```go
// Create a new hash from data
hash := hash.New(data)

// Get the hash bytes
bytes := hash.Bytes()

// Convert hash to hex string
hexString := hash.String()

// Create hash from hex string
hash, err := hash.FromHex(hexString)

// Combine two pieces of data into a single hash
combinedHash := hash.Combine(data1, data2)

// Convenience functions
hashBytes := hash.Sum512(data)
combinedHashBytes := hash.SumCombine(data1, data2)
```

### JavaScript/TypeScript API

```typescript
// Create a new hash from data
const hash = Hash.new(data);

// Get the hash bytes
const bytes = hash.getBytes();

// Convert hash to hex string
const hexString = hash.toHex();

// Create hash from hex string
const hash = Hash.fromHex(hexString);

// Combine two pieces of data into a single hash
const combinedHash = Hash.combine(data1, data2);

// Convenience functions
const hashBytes = hash(data);
const combinedHashBytes = hashCombine(data1, data2);
```

## Key Pair API

The Key Pair API provides functionality for generating and managing 512-bit cryptographic key pairs.

### Rust Key Pair API

```rust
// Generate a new key pair
let mut rng = rand::rngs::OsRng;
let keypair = topayz512::generate_keypair(&mut rng);

// Access the private and public keys
let private_key = &keypair.private_key;
let public_key = &keypair.public_key;

// Get the key bytes
let private_bytes: &[u8; 64] = private_key.as_bytes();
let public_bytes: &[u8; 64] = public_key.as_bytes();

// Convert keys to hex strings
let private_hex = private_key.to_hex();
let public_hex = public_key.to_hex();

// Create keys from hex strings
let private_key = topayz512::PrivateKey::from_hex(private_hex).unwrap();
let public_key = topayz512::PublicKey::from_hex(public_hex).unwrap();

// Derive public key from private key
let derived_public_key = topayz512::PublicKey::from_private_key(&private_key);

// Convenience function to derive public key from private key
let derived_public_key = topayz512::private_to_public(&private_key);
```

### Go Key Pair API

```go
// Generate a new key pair
kp, err := keypair.GenerateKeyPair()
if err != nil {
    // Handle error
}

// Access the private and public keys
privateKey := kp.PrivateKey
publicKey := kp.PublicKey

// Get the key bytes
privateBytes := privateKey.Bytes()
publicBytes := publicKey.Bytes()

// Convert keys to hex strings
privateHex := privateKey.String()
publicHex := publicKey.String()

// Create keys from hex strings
privateKey, err := keypair.PrivateKeyFromHex(privateHex)
publicKey, err := keypair.PublicKeyFromHex(publicHex)

// Derive public key from private key
derivedPublicKey := keypair.DerivePublicKey(privateKey)

// Convenience function to derive public key from private key
derivedPublicKey := keypair.PrivateToPublic(privateKey)
```

### JavaScript/TypeScript Key Pair API

```typescript
// Generate a new key pair
const keypair = generateKeyPair();

// Access the private and public keys
const privateKey = keypair.privateKey;
const publicKey = keypair.publicKey;

// Get the key bytes
const privateBytes = privateKey.getBytes();
const publicBytes = publicKey.getBytes();

// Convert keys to hex strings
const privateHex = privateKey.toHex();
const publicHex = publicKey.toHex();

// Create keys from hex strings
const privateKey = PrivateKey.fromHex(privateHex);
const publicKey = PublicKey.fromHex(publicHex);

// Derive public key from private key
const derivedPublicKey = PublicKey.fromPrivateKey(privateKey);

// Convenience function to derive public key from private key
const derivedPublicKey = privateToPublic(privateKey);
```

## Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `HASH_SIZE_BYTES` | 64 | Size of the hash output in bytes (512 bits) |
| `PRIVATE_KEY_SIZE_BYTES` | 64 | Size of the private key in bytes (512 bits) |
| `PUBLIC_KEY_SIZE_BYTES` | 64 | Size of the public key in bytes (512 bits) |

## Error Handling

### Rust

The Rust implementation returns `Result` types for operations that can fail, such as parsing from hex strings.

### Go

The Go implementation returns error values for operations that can fail.

### JavaScript/TypeScript

The JavaScript/TypeScript implementation throws exceptions for operations that can fail.
