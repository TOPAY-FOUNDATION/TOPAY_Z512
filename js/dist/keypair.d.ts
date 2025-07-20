/**
 * TOPAY-Z512 Key Pair implementation
 *
 * This module provides a 512-bit cryptographic key pair implementation
 * that can be used as part of the TOPAY-Z512 cryptographic suite.
 */
/** The size of TOPAY-Z512 keys in bytes (512 bits = 64 bytes) */
export declare const KEY_SIZE_BYTES = 64;
/**
 * Represents a TOPAY-Z512 private key (512 bits)
 */
export declare class PrivateKey {
    private readonly bytes;
    /**
     * Creates a new PrivateKey instance
     * @param bytes - The 64-byte private key value
     */
    constructor(bytes: Uint8Array);
    /**
     * Generates a new private key using secure random data
     * @returns A new PrivateKey instance
     */
    static generate(): PrivateKey;
    /**
     * Returns the private key as a byte array
     * @returns The private key bytes
     */
    getBytes(): Uint8Array;
    /**
     * Converts the private key to a hexadecimal string
     * @returns The private key as a hex string
     */
    toHex(): string;
    /**
     * Creates a private key from a hexadecimal string
     * @param hex - The hex string to convert
     * @returns A new PrivateKey instance
     */
    static fromHex(hex: string): PrivateKey;
}
/**
 * Represents a TOPAY-Z512 public key (512 bits)
 */
export declare class PublicKey {
    private readonly bytes;
    /**
     * Creates a new PublicKey instance
     * @param bytes - The 64-byte public key value
     */
    constructor(bytes: Uint8Array);
    /**
     * Derives a public key from a private key
     * @param privateKey - The private key to derive from
     * @returns A new PublicKey instance
     */
    static fromPrivateKey(privateKey: PrivateKey): PublicKey;
    /**
     * Returns the public key as a byte array
     * @returns The public key bytes
     */
    getBytes(): Uint8Array;
    /**
     * Converts the public key to a hexadecimal string
     * @returns The public key as a hex string
     */
    toHex(): string;
    /**
     * Creates a public key from a hexadecimal string
     * @param hex - The hex string to convert
     * @returns A new PublicKey instance
     */
    static fromHex(hex: string): PublicKey;
    /**
     * Compares this public key with another public key for equality
     * @param other - The other public key to compare with
     * @returns True if the public keys are equal, false otherwise
     */
    equals(other: PublicKey): boolean;
}
/**
 * Represents a TOPAY-Z512 key pair (private key and public key)
 */
export declare class KeyPair {
    readonly privateKey: PrivateKey;
    readonly publicKey: PublicKey;
    /**
     * Creates a new KeyPair instance
     * @param privateKey - The private key
     * @param publicKey - The public key
     */
    constructor(privateKey: PrivateKey, publicKey: PublicKey);
    /**
     * Generates a new key pair
     * @returns A new KeyPair instance
     */
    static generate(): KeyPair;
}
/**
 * Convenience function to generate a key pair
 * @returns A new KeyPair instance
 */
export declare function generateKeyPair(): KeyPair;
/**
 * Convenience function to derive a public key from a private key
 * @param privateKey - The private key to derive from
 * @returns A new PublicKey instance
 */
export declare function privateToPublic(privateKey: PrivateKey): PublicKey;
