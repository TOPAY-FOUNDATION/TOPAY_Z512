"use strict";
/**
 * TOPAY-Z512 Hash implementation
 *
 * This module provides a 512-bit cryptographic hash function implementation
 * that can be used as part of the TOPAY-Z512 cryptographic suite.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Hash = exports.HASH_SIZE_BYTES = void 0;
exports.hash = hash;
exports.hashCombine = hashCombine;
const js_sha3_1 = require("js-sha3");
/** The size of TOPAY-Z512 hash output in bytes (512 bits = 64 bytes) */
exports.HASH_SIZE_BYTES = 64;
/**
 * Represents a TOPAY-Z512 hash value (512 bits)
 */
class Hash {
    /**
     * Creates a new Hash instance
     * @param bytes - The 64-byte hash value
     */
    constructor(bytes) {
        if (bytes.length !== exports.HASH_SIZE_BYTES) {
            throw new Error(`Invalid hash length: ${bytes.length}, expected ${exports.HASH_SIZE_BYTES}`);
        }
        this.bytes = new Uint8Array(bytes);
    }
    /**
     * Creates a new hash from the given data
     * @param data - The data to hash
     * @returns A new Hash instance
     */
    static new(data) {
        const inputData = typeof data === 'string' ? new TextEncoder().encode(data) : data;
        const hashValue = js_sha3_1.sha3_512.arrayBuffer(inputData);
        return new Hash(new Uint8Array(hashValue));
    }
    /**
     * Creates a new hash by combining two input values
     * @param data1 - The first data to hash
     * @param data2 - The second data to hash
     * @returns A new Hash instance
     */
    static combine(data1, data2) {
        const inputData1 = typeof data1 === 'string' ? new TextEncoder().encode(data1) : data1;
        const inputData2 = typeof data2 === 'string' ? new TextEncoder().encode(data2) : data2;
        // First hash the individual inputs
        const hash1 = Hash.new(inputData1);
        const hash2 = Hash.new(inputData2);
        // Then combine the hashes with a separator to ensure it's different from concatenation
        const combinedBytes = new Uint8Array(hash1.getBytes().length + hash2.getBytes().length + 1);
        combinedBytes.set(hash1.getBytes(), 0);
        combinedBytes[hash1.getBytes().length] = 0xFF; // Add a separator byte
        combinedBytes.set(hash2.getBytes(), hash1.getBytes().length + 1);
        // Hash the combined result
        return Hash.new(combinedBytes);
    }
    /**
     * Returns the hash value as a byte array
     * @returns The hash bytes
     */
    getBytes() {
        return new Uint8Array(this.bytes);
    }
    /**
     * Converts the hash to a hexadecimal string
     * @returns The hash as a hex string
     */
    toHex() {
        return Array.from(this.bytes)
            .map(b => b.toString(16).padStart(2, '0'))
            .join('');
    }
    /**
     * Creates a hash from a hexadecimal string
     * @param hex - The hex string to convert
     * @returns A new Hash instance
     */
    static fromHex(hex) {
        if (hex.length !== exports.HASH_SIZE_BYTES * 2) {
            throw new Error(`Invalid hex length: ${hex.length}, expected ${exports.HASH_SIZE_BYTES * 2}`);
        }
        if (!/^[0-9a-fA-F]+$/.test(hex)) {
            throw new Error('Invalid hex string: contains non-hex characters');
        }
        const bytes = new Uint8Array(exports.HASH_SIZE_BYTES);
        for (let i = 0; i < exports.HASH_SIZE_BYTES; i++) {
            bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
        }
        return new Hash(bytes);
    }
    /**
     * Compares this hash with another hash for equality
     * @param other - The other hash to compare with
     * @returns True if the hashes are equal, false otherwise
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
    /**
     * Returns a string representation of the hash
     * @returns The hash as a hex string
     */
    toString() {
        return this.toHex();
    }
}
exports.Hash = Hash;
/**
 * A convenience function to hash data
 * @param data - The data to hash
 * @returns The hash bytes
 */
function hash(data) {
    return Hash.new(data).getBytes();
}
/**
 * A convenience function to hash two pieces of data together
 * @param data1 - The first data to hash
 * @param data2 - The second data to hash
 * @returns The hash bytes
 */
function hashCombine(data1, data2) {
    return Hash.combine(data1, data2).getBytes();
}
