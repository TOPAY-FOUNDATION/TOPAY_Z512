/**
 * TOPAY-Z512 Hash Example
 *
 * This example demonstrates how to use the TOPAY-Z512 hash functionality.
 * It uses the current time as input to show dynamic hashing.
 */

import { Hash, hash, hashCombine } from '../src';

console.log('TOPAY-Z512 Hash Example with Time-based Input\n');

// Get current time as input
const now = new Date();
const timeStr = `TOPAY-Z512 time: ${Math.floor(now.getTime() / 1000)} seconds`;
const data = timeStr;
const hashValue = Hash.new(data);

console.log(`Input: ${data}`);
console.log(`Hash: ${hashValue}`);
console.log(`Hash size: ${hashValue.getBytes().length} bytes\n`);

// Hash combination with time components
const timeMillis = now.getTime().toString();
const data1 = `TOPAY-${timeMillis}`;
const data2 = `Z512-${now.getTime() % 1000000}`;

const combinedHash = Hash.combine(data1, data2);

console.log(`Input 1: ${data1}`);
console.log(`Input 2: ${data2}`);
console.log(`Combined Hash: ${combinedHash}`);

// Concatenated hash (different from combined hash)
const concatenated = data1 + data2;
const concatHash = Hash.new(concatenated);

console.log(`Concatenated Hash: ${concatHash}`);
console.log(`Are they equal? ${combinedHash.equals(concatHash)}\n`);

// Convenience functions
const hashBytes = hash(data);
console.log(`Hash bytes (first 8): ${Array.from(hashBytes.slice(0, 8))}`);

const combinedBytes = hashCombine(data1, data2);
console.log(`Combined hash bytes (first 8): ${Array.from(combinedBytes.slice(0, 8))}`);

// Hex conversion
const hex = hashValue.toHex();
console.log(`\nHex string: ${hex}`);

try {
  const hashFromHex = Hash.fromHex(hex);
  console.log(`Converted back from hex: ${hashFromHex}`);
  console.log(`Equal to original: ${hashValue.equals(hashFromHex)}`);
} catch (e) {
  console.log(`Error: ${e}`);
}

// Binary data example
const binaryData = new Uint8Array([0, 1, 2, 3, 4, 5, 6, 7, 8, 9]);
const binaryHash = Hash.new(binaryData);
console.log(`\nBinary data hash: ${binaryHash}`);