/**
 * Test setup file for TOPAY-Z512 JavaScript/TypeScript implementation
 */

// Polyfill for crypto in test environment
import { webcrypto } from 'crypto';

// Make crypto available globally for tests
if (typeof globalThis.crypto === 'undefined') {
  (globalThis as any).crypto = webcrypto;
}

// Increase timeout for performance tests
jest.setTimeout(30000);

// Mock console methods for cleaner test output
const originalConsoleLog = console.log;
const originalConsoleError = console.error;

beforeAll(() => {
  // Suppress console output during tests unless explicitly needed
  console.log = jest.fn();
  console.error = jest.fn();
});

afterAll(() => {
  // Restore console methods
  console.log = originalConsoleLog;
  console.error = originalConsoleError;
});

// Global test utilities
(global as any).testUtils = {
  // Helper to create test data
  createTestData: (size: number): Uint8Array => {
    const data = new Uint8Array(size);
    for (let i = 0; i < size; i++) {
      data[i] = i % 256;
    }
    return data;
  },

  // Helper to create random test data
  createRandomTestData: (size: number): Uint8Array => {
    const data = new Uint8Array(size);
    crypto.getRandomValues(data);
    return data;
  },

  // Helper to compare arrays
  arraysEqual: (a: Uint8Array, b: Uint8Array): boolean => {
    if (a.length !== b.length) return false;
    for (let i = 0; i < a.length; i++) {
      if (a[i] !== b[i]) return false;
    }
    return true;
  },

  // Helper to measure execution time
  measureTime: async <T>(fn: () => Promise<T>): Promise<[number, T]> => {
    const start = Date.now();
    const result = await fn();
    const time = Date.now() - start;
    return [time, result];
  }
};
