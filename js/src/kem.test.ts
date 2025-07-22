import { Ciphertext, SharedSecret, keygen, encapsulate, decapsulate } from './kem';
import { describe, test, expect } from '@jest/globals';

describe('KEM', () => {
  describe('Ciphertext', () => {
    test('should create a ciphertext and convert to/from bytes', () => {
      // Generate a key pair
      const keypair = keygen();
      
      // Encapsulate a shared secret
      const { ciphertext } = encapsulate(keypair.publicKey);
      
      // Convert to bytes and back
      const bytes = ciphertext.getBytes();
      const ciphertext2 = Ciphertext.fromBytes(bytes);
      
      // Compare the bytes
      expect(ciphertext2.getBytes()).toEqual(bytes);
    });
    
    test('should convert ciphertext to/from hex', () => {
      // Generate a key pair
      const keypair = keygen();
      
      // Encapsulate a shared secret
      const { ciphertext } = encapsulate(keypair.publicKey);
      
      // Convert to hex and back
      const hex = ciphertext.toHex();
      const ciphertext2 = Ciphertext.fromHex(hex);
      
      // Compare the bytes
      expect(ciphertext2.getBytes()).toEqual(ciphertext.getBytes());
    });
    
    test('should throw error for invalid bytes length', () => {
      expect(() => {
        Ciphertext.fromBytes(new Uint8Array(10));
      }).toThrow('Invalid ciphertext length');
    });
    
    test('should throw error for invalid hex length', () => {
      expect(() => {
        Ciphertext.fromHex('abcd');
      }).toThrow('Invalid hex string length');
    });
  });
  
  describe('SharedSecret', () => {
    test('should create a shared secret and convert to/from bytes', () => {
      // Generate a key pair
      const keypair = keygen();
      
      // Encapsulate a shared secret
      const { sharedSecret } = encapsulate(keypair.publicKey);
      
      // Convert to bytes and back
      const bytes = sharedSecret.getBytes();
      const sharedSecret2 = new SharedSecret(bytes);
      
      // Compare the bytes
      expect(sharedSecret2.getBytes()).toEqual(bytes);
    });
    
    test('should convert shared secret to/from hex', () => {
      // Generate a key pair
      const keypair = keygen();
      
      // Encapsulate a shared secret
      const { sharedSecret } = encapsulate(keypair.publicKey);
      
      // Convert to hex and back
      const hex = sharedSecret.toHex();
      const sharedSecret2 = SharedSecret.fromHex(hex);
      
      // Compare the bytes
      expect(sharedSecret2.getBytes()).toEqual(sharedSecret.getBytes());
    });
    
    test('should throw error for invalid bytes length', () => {
      expect(() => {
        new SharedSecret(new Uint8Array(10));
      }).toThrow('Invalid shared secret length');
    });
    
    test('should throw error for invalid hex length', () => {
      expect(() => {
        SharedSecret.fromHex('abcd');
      }).toThrow('Invalid hex string length');
    });
  });
  
  describe('KEM operations', () => {
    test('should encapsulate and decapsulate correctly', () => {
      // Generate a key pair
      const keypair = keygen();
      
      // Encapsulate a shared secret
      const { ciphertext, sharedSecret: encapsulatedSecret } = encapsulate(keypair.publicKey);
      
      // Decapsulate the shared secret
      const decapsulatedSecret = decapsulate(keypair.privateKey, ciphertext);
      
      // Compare the shared secrets
      expect(decapsulatedSecret.getBytes()).toEqual(encapsulatedSecret.getBytes());
    });
  });
});