/**
 * TOPAY-Z512 Cryptographic Library
 * 
 * A 512-bit post-quantum cryptography library based on LWE.
 * This library provides implementations for key encapsulation mechanisms (KEM)
 * and cryptographic hashing functions.
 */

// Export modules
export * from './hash';
export * from './keypair';
export * from './kem';

// Library version information
export const VERSION = '0.1.0';
export const NAME = 'topayz512';