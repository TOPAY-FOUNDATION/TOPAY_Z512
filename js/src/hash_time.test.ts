import { describe, test, expect } from '@jest/globals';
import { Hash, HASH_SIZE_BYTES, newWithTime, hashWithTime } from './hash';

describe('Time-based Hash Functions', () => {
  test('Hash.newWithTime should create different hashes at different times', async () => {
    const hash1 = Hash.newWithTime();
    
    // Small delay to ensure different timestamps
    await new Promise(resolve => setTimeout(resolve, 10));
    
    const hash2 = Hash.newWithTime();
    
    // The hashes should be different due to different timestamps
    expect(hash1.equals(hash2)).toBe(false);
    expect(hash1.toHex()).not.toBe(hash2.toHex());
    
    // Verify the hash size is correct
    expect(hash1.getBytes().length).toBe(HASH_SIZE_BYTES);
  });
  
  test('newWithTime should create different hashes at different times', async () => {
    const hash1 = newWithTime();
    
    // Small delay to ensure different timestamps
    await new Promise(resolve => setTimeout(resolve, 10));
    
    const hash2 = newWithTime();
    
    // The hashes should be different due to different timestamps
    expect(hash1.equals(hash2)).toBe(false);
    expect(hash1.toHex()).not.toBe(hash2.toHex());
    
    // Verify the hash size is correct
    expect(hash1.getBytes().length).toBe(HASH_SIZE_BYTES);
  });
  
  test('hashWithTime should create different hash bytes at different times', async () => {
    const hashBytes1 = hashWithTime();
    
    // Small delay to ensure different timestamps
    await new Promise(resolve => setTimeout(resolve, 10));
    
    const hashBytes2 = hashWithTime();
    
    // The hash bytes should be different due to different timestamps
    expect(hashBytes1).not.toEqual(hashBytes2);
    
    // Verify the hash size is correct
    expect(hashBytes1.length).toBe(HASH_SIZE_BYTES);
  });
});