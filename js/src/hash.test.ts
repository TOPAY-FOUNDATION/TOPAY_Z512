import { Hash, hash, hashCombine, HASH_SIZE_BYTES } from './hash';
import { describe, test, expect } from '@jest/globals';

describe('Hash', () => {
  test('hash size should be correct', () => {
    const data = 'TOPAY-Z512 test data';
    const hashValue = Hash.new(data);
    expect(hashValue.getBytes().length).toBe(HASH_SIZE_BYTES);
  });

  test('hash should be deterministic', () => {
    const data = 'TOPAY-Z512 test data';
    const hash1 = Hash.new(data);
    const hash2 = Hash.new(data);
    expect(hash1.equals(hash2)).toBe(true);
    expect(hash1.toHex()).toBe(hash2.toHex());
  });

  test('different inputs should produce different hashes', () => {
    const data1 = 'TOPAY-Z512 test data 1';
    const data2 = 'TOPAY-Z512 test data 2';
    const hash1 = Hash.new(data1);
    const hash2 = Hash.new(data2);
    expect(hash1.equals(hash2)).toBe(false);
    expect(hash1.toHex()).not.toBe(hash2.toHex());
  });

  test('combine should be different from concatenation', () => {
    const data1 = 'TOPAY-Z512';
    const data2 = 'test data';
    
    // Combined hash
    const combined = Hash.combine(data1, data2);
    
    // Concatenated hash
    const concatenated = Hash.new(data1 + data2);
    
    expect(combined.equals(concatenated)).toBe(false);
    expect(combined.toHex()).not.toBe(concatenated.toHex());
  });

  test('hex conversion should be reversible', () => {
    const data = 'TOPAY-Z512 hex conversion test';
    const hashValue = Hash.new(data);
    const hex = hashValue.toHex();
    const hashFromHex = Hash.fromHex(hex);
    
    expect(hashValue.equals(hashFromHex)).toBe(true);
    expect(hashValue.toHex()).toBe(hashFromHex.toHex());
  });

  test('should reject invalid hex strings', () => {
    // Invalid length
    expect(() => {
      Hash.fromHex('invalid');
    }).toThrow(/Invalid hex length/);
    
    // Invalid characters
    expect(() => {
      Hash.fromHex('zz' + '0'.repeat(HASH_SIZE_BYTES * 2 - 2));
    }).toThrow(/Invalid hex string/);
  });

  test('convenience functions should work correctly', () => {
    const data = 'TOPAY-Z512 test data';
    const data1 = 'TOPAY-Z512';
    const data2 = 'test data';
    
    // Test hash function
    const hashBytes = hash(data);
    expect(hashBytes.length).toBe(HASH_SIZE_BYTES);
    expect(hashBytes).toEqual(Hash.new(data).getBytes());
    
    // Test hashCombine function
    const combineBytes = hashCombine(data1, data2);
    expect(combineBytes.length).toBe(HASH_SIZE_BYTES);
    expect(combineBytes).toEqual(Hash.combine(data1, data2).getBytes());
  });

  test('should work with Uint8Array input', () => {
    const textData = 'TOPAY-Z512 test data';
    const binaryData = new TextEncoder().encode(textData);
    
    const hashFromText = Hash.new(textData);
    const hashFromBinary = Hash.new(binaryData);
    
    expect(hashFromText.equals(hashFromBinary)).toBe(true);
  });
});