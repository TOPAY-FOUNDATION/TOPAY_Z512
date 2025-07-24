# Key Pair Management Guide

## Overview

This guide covers best practices for key pair generation, storage, and management in TOPAY-Z512. Proper key management is crucial for maintaining the security of your cryptographic operations.

## Key Generation

### Secure Random Generation

TOPAY-Z512 uses cryptographically secure random number generators:

- **Rust**: `getrandom` crate with OS entropy
- **Go**: `crypto/rand` package
- **JavaScript**: `crypto.getRandomValues()` (Browser) or `crypto.randomBytes()` (Node.js)

### Generation Process

```rust
// Rust
let keypair = KeyPair::generate()?;
```

```go
// Go
keypair, err := GenerateKeyPair()
```

```typescript
// TypeScript
const keypair = await generateKeyPair();
```

### Performance Considerations

Key generation is computationally intensive:

- **Time**: 0.5-2ms depending on platform
- **Memory**: ~4KB working memory
- **Entropy**: 256 bits of system entropy required

## Key Storage

### In-Memory Storage

#### Secure Memory Handling

```rust
// Rust - Automatic secure cleanup
{
    let keypair = KeyPair::generate()?;
    // Private key automatically zeroed when dropped
}
```

```go
// Go - Manual secure cleanup
keypair, err := GenerateKeyPair()
defer keypair.PrivateKey.SecureZero()
```

```typescript
// TypeScript - Manual secure cleanup
const keypair = await generateKeyPair();
try {
    // Use keypair
} finally {
    keypair.privateKey.secureZero();
}
```

#### Memory Protection

- **Stack allocation**: Prefer stack over heap when possible
- **Secure pages**: Use mlock/VirtualLock for sensitive data
- **Swap protection**: Prevent swapping of key material

### Persistent Storage

#### File-Based Storage

```rust
// Rust - Encrypted storage
use std::fs;
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
struct EncryptedKeyPair {
    encrypted_private_key: Vec<u8>,
    public_key: Vec<u8>,
    salt: Vec<u8>,
    nonce: Vec<u8>,
}

fn save_keypair_encrypted(keypair: &KeyPair, password: &str, path: &str) -> Result<(), TopayError> {
    let salt = generate_random_bytes(32)?;
    let key = derive_key_from_password(password, &salt)?;
    let nonce = generate_random_bytes(12)?;
    
    let encrypted_private = encrypt_aes_gcm(&keypair.private_key.to_bytes(), &key, &nonce)?;
    
    let encrypted_keypair = EncryptedKeyPair {
        encrypted_private_key: encrypted_private,
        public_key: keypair.public_key.to_bytes(),
        salt,
        nonce,
    };
    
    let serialized = serde_json::to_string(&encrypted_keypair)?;
    fs::write(path, serialized)?;
    
    Ok(())
}
```

#### Database Storage

```sql
-- PostgreSQL schema for key storage
CREATE TABLE keypairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    public_key BYTEA NOT NULL,
    encrypted_private_key BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE
);

-- Index for efficient lookups
CREATE INDEX idx_keypairs_user_id ON keypairs(user_id);
CREATE INDEX idx_keypairs_public_key ON keypairs(public_key);
```

#### Hardware Security Modules (HSM)

```rust
// Example HSM integration
trait HsmKeyStore {
    fn generate_keypair(&self, label: &str) -> Result<KeyPairHandle, HsmError>;
    fn sign(&self, handle: &KeyPairHandle, data: &[u8]) -> Result<Vec<u8>, HsmError>;
    fn get_public_key(&self, handle: &KeyPairHandle) -> Result<PublicKey, HsmError>;
}

struct Pkcs11KeyStore {
    session: pkcs11::Session,
}

impl HsmKeyStore for Pkcs11KeyStore {
    fn generate_keypair(&self, label: &str) -> Result<KeyPairHandle, HsmError> {
        // PKCS#11 key generation
        let (public_handle, private_handle) = self.session.generate_key_pair(
            &pkcs11::Mechanism::TopayZ512KeyGen,
            &[
                pkcs11::Attribute::Label(label.as_bytes()),
                pkcs11::Attribute::Token(true),
                pkcs11::Attribute::Private(true),
            ],
            &[
                pkcs11::Attribute::Label(label.as_bytes()),
                pkcs11::Attribute::Token(true),
                pkcs11::Attribute::Private(true),
                pkcs11::Attribute::Sensitive(true),
                pkcs11::Attribute::Extractable(false),
            ],
        )?;
        
        Ok(KeyPairHandle {
            public_handle,
            private_handle,
        })
    }
}
```

## Key Serialization

### Binary Format

```rust
// Rust - Raw binary serialization
let public_bytes = keypair.public_key.to_bytes();
let private_bytes = keypair.private_key.to_bytes();

// Deserialization
let public_key = PublicKey::from_bytes(&public_bytes)?;
let private_key = PrivateKey::from_bytes(&private_bytes)?;
```

### Text Formats

#### Hexadecimal

```rust
// Rust - Hex encoding
let public_hex = keypair.public_key.to_hex();
let private_hex = keypair.private_key.to_hex();

// Decoding
let public_key = PublicKey::from_hex(&public_hex)?;
let private_key = PrivateKey::from_hex(&private_hex)?;
```

#### Base64

```typescript
// TypeScript - Base64 encoding
const publicB64 = btoa(String.fromCharCode(...keypair.publicKey.toBytes()));
const privateB64 = btoa(String.fromCharCode(...keypair.privateKey.toBytes()));

// Decoding
const publicBytes = new Uint8Array(atob(publicB64).split('').map(c => c.charCodeAt(0)));
const publicKey = PublicKey.fromBytes(publicBytes);
```

#### PEM Format

```rust
// Rust - PEM-like format
fn to_pem(key: &[u8], key_type: &str) -> String {
    let b64 = base64::encode(key);
    let mut pem = format!("-----BEGIN {} KEY-----\n", key_type);
    
    // Split into 64-character lines
    for chunk in b64.as_bytes().chunks(64) {
        pem.push_str(&String::from_utf8_lossy(chunk));
        pem.push('\n');
    }
    
    pem.push_str(&format!("-----END {} KEY-----\n", key_type));
    pem
}

// Usage
let public_pem = to_pem(&keypair.public_key.to_bytes(), "TOPAY-Z512 PUBLIC");
let private_pem = to_pem(&keypair.private_key.to_bytes(), "TOPAY-Z512 PRIVATE");
```

## Key Lifecycle Management

### Key Rotation

```rust
// Rust - Key rotation strategy
struct KeyManager {
    current_keypair: KeyPair,
    previous_keypairs: Vec<(KeyPair, SystemTime)>,
    rotation_interval: Duration,
}

impl KeyManager {
    fn rotate_keys(&mut self) -> Result<(), TopayError> {
        // Archive current key
        self.previous_keypairs.push((
            std::mem::replace(&mut self.current_keypair, KeyPair::generate()?),
            SystemTime::now(),
        ));
        
        // Clean up old keys (keep for 30 days)
        let cutoff = SystemTime::now() - Duration::from_secs(30 * 24 * 3600);
        self.previous_keypairs.retain(|(_, timestamp)| *timestamp > cutoff);
        
        Ok(())
    }
    
    fn should_rotate(&self) -> bool {
        // Rotate every 90 days
        self.current_keypair.created_at().elapsed().unwrap_or_default() > self.rotation_interval
    }
}
```

### Key Backup and Recovery

```rust
// Rust - Shamir's Secret Sharing for key backup
use shamir_secret_sharing::ShamirSecretSharing;

fn backup_private_key(private_key: &PrivateKey, threshold: u8, shares: u8) -> Result<Vec<Vec<u8>>, TopayError> {
    let sss = ShamirSecretSharing::new(threshold, shares)?;
    let key_bytes = private_key.to_bytes();
    let shares = sss.split(&key_bytes)?;
    Ok(shares)
}

fn recover_private_key(shares: &[Vec<u8>]) -> Result<PrivateKey, TopayError> {
    let sss = ShamirSecretSharing::new(shares.len() as u8, shares.len() as u8)?;
    let recovered_bytes = sss.recover(shares)?;
    PrivateKey::from_bytes(&recovered_bytes)
}
```

## Security Best Practices

### Key Generation Guidelines

1. **Use system entropy**: Always use OS-provided random sources
2. **Generate on secure hardware**: Use HSMs when available
3. **Verify randomness**: Test random number generator quality
4. **Avoid predictable seeds**: Never use timestamps or user input as seeds

### Secure Storage Guidelines

1. **Encrypt at rest**: Always encrypt private keys when stored
2. **Use strong passwords**: Derive encryption keys from strong passwords
3. **Separate storage**: Store public and private keys separately
4. **Access control**: Implement strict access controls
5. **Audit logging**: Log all key access and operations

### Key Usage

1. **Minimize exposure**: Keep private keys in memory only when needed
2. **Secure channels**: Use TLS for key transmission
3. **Rate limiting**: Implement rate limits on key operations
4. **Monitoring**: Monitor for unusual key usage patterns

### Key Destruction

1. **Secure deletion**: Overwrite key material multiple times
2. **Memory clearing**: Clear all copies from memory
3. **Backup cleanup**: Securely delete all backup copies
4. **Certificate revocation**: Revoke associated certificates

## Common Pitfalls

### Weak Random Number Generation

```rust
// ❌ BAD - Predictable seed
let mut rng = StdRng::seed_from_u64(SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs());

// ✅ GOOD - System entropy
let mut rng = StdRng::from_entropy();
```

### Insecure Key Storage

```rust
// ❌ BAD - Plaintext storage
std::fs::write("private_key.bin", keypair.private_key.to_bytes())?;

// ✅ GOOD - Encrypted storage
save_keypair_encrypted(&keypair, &password, "keypair.enc")?;
```

### Memory Leaks

```rust
// ❌ BAD - Key material in heap
let private_key_vec = keypair.private_key.to_bytes();
// Vector may not be securely cleared

// ✅ GOOD - Secure memory handling
{
    let private_key_bytes = keypair.private_key.to_bytes();
    // Use key...
    secure_zero(&mut private_key_bytes);
}
```

### Insufficient Access Control

```rust
// ❌ BAD - No access control
fn get_private_key() -> PrivateKey {
    // Anyone can call this
}

// ✅ GOOD - Proper access control
fn get_private_key(user: &AuthenticatedUser, key_id: &str) -> Result<PrivateKey, AuthError> {
    if !user.has_permission(Permission::AccessPrivateKey(key_id)) {
        return Err(AuthError::Forbidden);
    }
    // Return key only if authorized
}
```

## Integration Examples

### Web Application

```typescript
// TypeScript - Browser key management
class BrowserKeyManager {
    private static readonly STORAGE_KEY = 'topay_keypair';
    
    static async generateAndStore(password: string): Promise<KeyPair> {
        const keypair = await generateKeyPair();
        await this.storeEncrypted(keypair, password);
        return keypair;
    }
    
    static async loadFromStorage(password: string): Promise<KeyPair | null> {
        const encrypted = localStorage.getItem(this.STORAGE_KEY);
        if (!encrypted) return null;
        
        return await this.decryptKeypair(encrypted, password);
    }
    
    private static async storeEncrypted(keypair: KeyPair, password: string): Promise<void> {
        const salt = crypto.getRandomValues(new Uint8Array(32));
        const key = await this.deriveKey(password, salt);
        const iv = crypto.getRandomValues(new Uint8Array(12));
        
        const encrypted = await crypto.subtle.encrypt(
            { name: 'AES-GCM', iv },
            key,
            keypair.privateKey.toBytes()
        );
        
        const data = {
            encrypted: Array.from(new Uint8Array(encrypted)),
            salt: Array.from(salt),
            iv: Array.from(iv),
            publicKey: Array.from(keypair.publicKey.toBytes())
        };
        
        localStorage.setItem(this.STORAGE_KEY, JSON.stringify(data));
    }
}
```

### Server Application

```go
// Go - Server key management
type ServerKeyManager struct {
    db     *sql.DB
    hsm    HSMInterface
    logger *log.Logger
}

func (km *ServerKeyManager) GenerateUserKeypair(userID string) (*KeyPair, error) {
    // Generate keypair in HSM
    keypair, err := km.hsm.GenerateKeypair(fmt.Sprintf("user_%s", userID))
    if err != nil {
        return nil, fmt.Errorf("HSM key generation failed: %w", err)
    }
    
    // Store public key in database
    _, err = km.db.Exec(`
        INSERT INTO user_keypairs (user_id, public_key, hsm_handle, created_at)
        VALUES ($1, $2, $3, NOW())
    `, userID, keypair.PublicKey.ToBytes(), keypair.HSMHandle)
    
    if err != nil {
        return nil, fmt.Errorf("database storage failed: %w", err)
    }
    
    km.logger.Printf("Generated keypair for user %s", userID)
    return keypair, nil
}
```

## Compliance and Auditing

### FIPS 140-2 Compliance

- Use FIPS-approved random number generators
- Implement proper key zeroization
- Maintain audit logs of all key operations
- Use certified cryptographic modules

### SOC 2 Compliance

- Implement proper access controls
- Maintain detailed audit logs
- Regular security assessments
- Incident response procedures

### GDPR Compliance

- Implement right to erasure for key material
- Maintain data processing records
- Implement privacy by design
- Regular data protection impact assessments

## Monitoring and Alerting

```rust
// Rust - Key operation monitoring
struct KeyOperationMetrics {
    generations: Counter,
    usage_count: Counter,
    errors: Counter,
    response_time: Histogram,
}

impl KeyOperationMetrics {
    fn record_generation(&self, duration: Duration) {
        self.generations.inc();
        self.response_time.observe(duration.as_secs_f64());
    }
    
    fn record_error(&self, error_type: &str) {
        self.errors.with_label_values(&[error_type]).inc();
    }
}

// Alert on suspicious patterns
fn check_for_anomalies(metrics: &KeyOperationMetrics) {
    if metrics.generations.get() > 1000.0 {
        alert("High key generation rate detected");
    }
    
    if metrics.errors.get() > 10.0 {
        alert("High error rate in key operations");
    }
}
```

This comprehensive guide covers all aspects of key pair management in TOPAY-Z512, from generation to secure destruction, ensuring that developers can implement robust and secure key management practices.
