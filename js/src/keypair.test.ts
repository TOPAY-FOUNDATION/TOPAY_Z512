/**
 * Tests for the TOPAY-Z512 Key Pair implementation
 */

import { KeyPair, PrivateKey, PublicKey, KEY_SIZE_BYTES, generateKeyPair } from './keypair';
import { describe, it, expect, jest } from '@jest/globals';

describe('KeyPair', () => {
  describe('generation', () => {
    it('should generate a valid key pair', () => {
      const keypair = KeyPair.generate();
      
      // Check that keys have correct size
      expect(keypair.privateKey.getBytes().length).toBe(KEY_SIZE_BYTES);
      expect(keypair.publicKey.getBytes().length).toBe(KEY_SIZE_BYTES);
    });
    
    it('should generate unique key pairs', () => {
      const keypair1 = KeyPair.generate();
      const keypair2 = KeyPair.generate();
      
      // Keys should be different
      expect(keypair1.privateKey.toHex()).not.toBe(keypair2.privateKey.toHex());
      expect(keypair1.publicKey.toHex()).not.toBe(keypair2.publicKey.toHex());
    });
    
    it('should provide a convenience function for generation', () => {
      const keypair = generateKeyPair();
      
      // Check that keys have correct size
      expect(keypair.privateKey.getBytes().length).toBe(KEY_SIZE_BYTES);
      expect(keypair.publicKey.getBytes().length).toBe(KEY_SIZE_BYTES);
    });
  });
});

describe('PrivateKey', () => {
  it('should reject invalid key sizes', () => {
    expect(() => {
      new PrivateKey(new Uint8Array(KEY_SIZE_BYTES - 1));
    }).toThrow();
    
    expect(() => {
      new PrivateKey(new Uint8Array(KEY_SIZE_BYTES + 1));
    }).toThrow();
  });
  
  it('should convert to and from hex correctly', () => {
    const privateKey = PrivateKey.generate();
    const hex = privateKey.toHex();
    
    // Hex string should have correct length
    expect(hex.length).toBe(KEY_SIZE_BYTES * 2);
    
    // Converting back should give the same key
    const recoveredKey = PrivateKey.fromHex(hex);
    expect(recoveredKey.toHex()).toBe(hex);
  });
  
  it('should reject invalid hex strings', () => {
    // Too short
    expect(() => {
      PrivateKey.fromHex('abcd');
    }).toThrow();
    
    // Invalid characters
    expect(() => {
      PrivateKey.fromHex('z'.repeat(KEY_SIZE_BYTES * 2));
    }).toThrow();
  });
});

describe('PublicKey', () => {
  it('should reject invalid key sizes', () => {
    expect(() => {
      new PublicKey(new Uint8Array(KEY_SIZE_BYTES - 1));
    }).toThrow();
    
    expect(() => {
      new PublicKey(new Uint8Array(KEY_SIZE_BYTES + 1));
    }).toThrow();
  });
  
  it('should derive consistently from private key', () => {
    const privateKey = PrivateKey.generate();
    const publicKey1 = PublicKey.fromPrivateKey(privateKey);
    const publicKey2 = PublicKey.fromPrivateKey(privateKey);
    
    // Deriving twice should give the same public key
    expect(publicKey1.toHex()).toBe(publicKey2.toHex());
  });
  
  it('should convert to and from hex correctly', () => {
    const privateKey = PrivateKey.generate();
    const publicKey = PublicKey.fromPrivateKey(privateKey);
    const hex = publicKey.toHex();
    
    // Hex string should have correct length
    expect(hex.length).toBe(KEY_SIZE_BYTES * 2);
    
    // Converting back should give the same key
    const recoveredKey = PublicKey.fromHex(hex);
    expect(recoveredKey.toHex()).toBe(hex);
  });
  
  it('should compare keys correctly', () => {
    const privateKey = PrivateKey.generate();
    const publicKey1 = PublicKey.fromPrivateKey(privateKey);
    const publicKey2 = PublicKey.fromPrivateKey(privateKey);
    const publicKey3 = PublicKey.fromPrivateKey(PrivateKey.generate());
    
    // Same keys should be equal
    expect(publicKey1.equals(publicKey2)).toBe(true);
    
    // Different keys should not be equal
    expect(publicKey1.equals(publicKey3)).toBe(false);
  });
});