import { Hash, newWithTime, hashWithTime } from '../src/hash';

// Example of using time-based hash functions
console.log('TOPAY-Z512 Time-based Hash Example\n');

// Create a hash using the current time with the static method
const timeHash = Hash.newWithTime();
console.log(`Time-based hash: ${timeHash.toHex()}`);

// Create another hash after a short delay
setTimeout(() => {
  const timeHash2 = Hash.newWithTime();
  console.log(`Time-based hash after 1 second: ${timeHash2.toHex()}`);
  
  // Demonstrate that the hashes are different
  console.log(`Hashes are different: ${!timeHash.equals(timeHash2)}`);
  
  // Using the convenience functions
  const timeHashObj = newWithTime();
  console.log(`Time-based hash object: ${timeHashObj.toHex()}`);
  
  const timeHashBytes = hashWithTime();
  console.log(`Time-based hash bytes length: ${timeHashBytes.length}`);
}, 1000);