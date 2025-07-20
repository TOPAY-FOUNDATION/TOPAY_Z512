"use strict";
/**
 * TOPAY-Z512 Key Pair implementation
 *
 * This module provides a 512-bit cryptographic key pair implementation
 * that can be used as part of the TOPAY-Z512 cryptographic suite.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.KeyPair = exports.PublicKey = exports.PrivateKey = exports.KEY_SIZE_BYTES = void 0;
exports.generateKeyPair = generateKeyPair;
exports.privateToPublic = privateToPublic;
const js_sha3_1 = require("js-sha3");
const hash_1 = require("./hash");
/** The size of TOPAY-Z512 keys in bytes (512 bits = 64 bytes) */
exports.KEY_SIZE_BYTES = hash_1.HASH_SIZE_BYTES;
/**
 * Represents a TOPAY-Z512 private key (512 bits)
 */
class PrivateKey {
    /**
     * Creates a new PrivateKey instance
     * @param bytes - The 64-byte private key value
     */
    constructor(bytes) {
        if (bytes.length !== exports.KEY_SIZE_BYTES) {
            throw new Error(`Invalid private key length: ${bytes.length}, expected ${exports.KEY_SIZE_BYTES}`);
        }
        this.bytes = new Uint8Array(bytes);
    }
    /**
     * Generates a new private key using secure random data
     * @returns A new PrivateKey instance
     */
    static generate() {
        const bytes = new Uint8Array(exports.KEY_SIZE_BYTES);
        crypto.getRandomValues(bytes);
        return new PrivateKey(bytes);
    }
    /**
     * Returns the private key as a byte array
     * @returns The private key bytes
     */
    getBytes() {
        return new Uint8Array(this.bytes);
    }
    /**
     * Converts the private key to a hexadecimal string
     * @returns The private key as a hex string
     */
    toHex() {
        return Array.from(this.bytes)
            .map(b => b.toString(16).padStart(2, '0'))
            .join('');
    }
    /**
     * Creates a private key from a hexadecimal string
     * @param hex - The hex string to convert
     * @returns A new PrivateKey instance
     */
    static fromHex(hex) {
        if (hex.length !== exports.KEY_SIZE_BYTES * 2) {
            throw new Error(`Invalid hex length: ${hex.length}, expected ${exports.KEY_SIZE_BYTES * 2}`);
        }
        if (!/^[0-9a-fA-F]+$/.test(hex)) {
            throw new Error('Invalid hex string: contains non-hex characters');
        }
        const bytes = new Uint8Array(exports.KEY_SIZE_BYTES);
        for (let i = 0; i < exports.KEY_SIZE_BYTES; i++) {
            bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
        }
        return new PrivateKey(bytes);
    }
}
exports.PrivateKey = PrivateKey;
/**
 * Represents a TOPAY-Z512 public key (512 bits)
 */
class PublicKey {
    /**
     * Creates a new PublicKey instance
     * @param bytes - The 64-byte public key value
     */
    constructor(bytes) {
        if (bytes.length !== exports.KEY_SIZE_BYTES) {
            throw new Error(`Invalid public key length: ${bytes.length}, expected ${exports.KEY_SIZE_BYTES}`);
        }
        this.bytes = new Uint8Array(bytes);
    }
    /**
     * Derives a public key from a private key
     * @param privateKey - The private key to derive from
     * @returns A new PublicKey instance
     */
    static fromPrivateKey(privateKey) {
        const privateKeyBytes = privateKey.getBytes();
        const hashValue = js_sha3_1.sha3_512.arrayBuffer(privateKeyBytes);
        return new PublicKey(new Uint8Array(hashValue));
    }
    /**
     * Returns the public key as a byte array
     * @returns The public key bytes
     */
    getBytes() {
        return new Uint8Array(this.bytes);
    }
    /**
     * Converts the public key to a hexadecimal string
     * @returns The public key as a hex string
     */
    toHex() {
        return Array.from(this.bytes)
            .map(b => b.toString(16).padStart(2, '0'))
            .join('');
    }
    /**
     * Creates a public key from a hexadecimal string
     * @param hex - The hex string to convert
     * @returns A new PublicKey instance
     */
    static fromHex(hex) {
        if (hex.length !== exports.KEY_SIZE_BYTES * 2) {
            throw new Error(`Invalid hex length: ${hex.length}, expected ${exports.KEY_SIZE_BYTES * 2}`);
        }
        if (!/^[0-9a-fA-F]+$/.test(hex)) {
            throw new Error('Invalid hex string: contains non-hex characters');
        }
        const bytes = new Uint8Array(exports.KEY_SIZE_BYTES);
        for (let i = 0; i < exports.KEY_SIZE_BYTES; i++) {
            bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
        }
        return new PublicKey(bytes);
    }
    /**
     * Compares this public key with another public key for equality
     * @param other - The other public key to compare with
     * @returns True if the public keys are equal, false otherwise
     */
    equals(other) {
        if (this.bytes.length !== other.bytes.length) {
            return false;
        }
        for (let i = 0; i < this.bytes.length; i++) {
            if (this.bytes[i] !== other.bytes[i]) {
                return false;
            }
        }
        return true;
    }
}
exports.PublicKey = PublicKey;
/**
 * Represents a TOPAY-Z512 key pair (private key and public key)
 */
class KeyPair {
    /**
     * Creates a new KeyPair instance
     * @param privateKey - The private key
     * @param publicKey - The public key
     */
    constructor(privateKey, publicKey) {
        this.privateKey = privateKey;
        this.publicKey = publicKey;
    }
    /**
     * Generates a new key pair
     * @returns A new KeyPair instance
     */
    static generate() {
        const privateKey = PrivateKey.generate();
        const publicKey = PublicKey.fromPrivateKey(privateKey);
        return new KeyPair(privateKey, publicKey);
    }
}
exports.KeyPair = KeyPair;
/**
 * Convenience function to generate a key pair
 * @returns A new KeyPair instance
 */
function generateKeyPair() {
    return KeyPair.generate();
}
/**
 * Convenience function to derive a public key from a private key
 * @param privateKey - The private key to derive from
 * @returns A new PublicKey instance
 */
function privateToPublic(privateKey) {
    return PublicKey.fromPrivateKey(privateKey);
}
