/**
 * TOPAY-Z512 Hash implementation
 *
 * This module provides a 512-bit cryptographic hash function implementation
 * that can be used as part of the TOPAY-Z512 cryptographic suite.
 */
/** The size of TOPAY-Z512 hash output in bytes (512 bits = 64 bytes) */
export declare const HASH_SIZE_BYTES = 64;
/**
 * Represents a TOPAY-Z512 hash value (512 bits)
 */
export declare class Hash {
    private readonly bytes;
    /**
     * Creates a new Hash instance
     * @param bytes - The 64-byte hash value
     */
    constructor(bytes: Uint8Array);
    /**
     * Creates a new hash from the given data
     * @param data - The data to hash
     * @returns A new Hash instance
     */
    static new(data: Uint8Array | string): Hash;
    /**
     * Creates a new hash by combining two input values
     * @param data1 - The first data to hash
     * @param data2 - The second data to hash
     * @returns A new Hash instance
     */
    static combine(data1: Uint8Array | string, data2: Uint8Array | string): Hash;
    /**
     * Returns the hash value as a byte array
     * @returns The hash bytes
     */
    getBytes(): Uint8Array;
    /**
     * Converts the hash to a hexadecimal string
     * @returns The hash as a hex string
     */
    toHex(): string;
    /**
     * Creates a hash from a hexadecimal string
     * @param hex - The hex string to convert
     * @returns A new Hash instance
     */
    static fromHex(hex: string): Hash;
    /**
     * Compares this hash with another hash for equality
     * @param other - The other hash to compare with
     * @returns True if the hashes are equal, false otherwise
     */
    equals(other: Hash): boolean;
    /**
     * Returns a string representation of the hash
     * @returns The hash as a hex string
     */
    toString(): string;
}
/**
 * A convenience function to hash data
 * @param data - The data to hash
 * @returns The hash bytes
 */
export declare function hash(data: Uint8Array | string): Uint8Array;
/**
 * A convenience function to hash two pieces of data together
 * @param data1 - The first data to hash
 * @param data2 - The second data to hash
 * @returns The hash bytes
 */
export declare function hashCombine(data1: Uint8Array | string, data2: Uint8Array | string): Uint8Array;
