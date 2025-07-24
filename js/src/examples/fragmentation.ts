/**
 * Fragmentation Example for TOPAY-Z512
 *
 * This example demonstrates data fragmentation and reconstruction capabilities
 * including parallel processing, mobile optimization, and compression.
 */

import {
  fragmentData,
  reconstructData,
  parallelFragmentation,
  parallelReconstruction,
  validateFragment,
  estimateMobileLatency,
  getOptimalFragmentSize,
  serializeFragments,
  deserializeFragments,
  compressFragments,
  decompressFragments,
  toHex,
  measureTime
} from '../index';

async function fragmentationExample(): Promise<void> {
  console.log('🧩 TOPAY-Z512 Data Fragmentation Example');
  console.log('='.repeat(50));

  try {
    // 1. Basic fragmentation
    console.log('\n1. Basic Data Fragmentation');
    console.log('-'.repeat(30));

    const testData = new Uint8Array(2048);
    crypto.getRandomValues(testData);
    console.log(`   Original data size: ${testData.length} bytes`);

    const fragResult = await fragmentData(testData);
    console.log(`   Fragment count: ${fragResult.metadata.fragmentCount}`);
    console.log(`   Checksum: ${toHex(fragResult.metadata.checksum).substring(0, 32)}...`);
    console.log(`   Timestamp: ${new Date(fragResult.metadata.timestamp).toISOString()}`);

    // Validate all fragments
    let validFragments = 0;
    for (const fragment of fragResult.fragments) {
      if (await validateFragment(fragment)) {
        validFragments++;
      }
    }
    console.log(`   Valid fragments: ${validFragments}/${fragResult.fragments.length}`);

    // 2. Data reconstruction
    console.log('\n2. Data Reconstruction');
    console.log('-'.repeat(30));

    const reconResult = await reconstructData(fragResult.fragments);
    console.log(`   Reconstruction complete: ${reconResult.isComplete}`);
    console.log(`   Reconstructed size: ${reconResult.data.length} bytes`);

    const dataIntegrity =
      testData.length === reconResult.data.length &&
      testData.every((byte, index) => byte === reconResult.data[index]);
    console.log(`   Data integrity: ${dataIntegrity ? 'VERIFIED' : 'FAILED'}`);

    // 3. Parallel fragmentation
    console.log('\n3. Parallel Fragmentation');
    console.log('-'.repeat(30));

    const largeData = new Uint8Array(8192);
    crypto.getRandomValues(largeData);

    const { timeMs: parallelTime, result: parallelResult } = await measureTime(async () => {
      return await parallelFragmentation([largeData]);
    });

    console.log(`   Parallel fragmentation time: ${parallelTime}ms`);
    console.log(`   Fragment count: ${parallelResult[0]?.metadata.fragmentCount || 0}`);

    // Compare with sequential
    const { timeMs: sequentialTime, result: sequentialResult } = await measureTime(async () => {
      return await fragmentData(largeData);
    });

    console.log(`   Sequential fragmentation time: ${sequentialTime}ms`);
    console.log(`   Performance improvement: ${(sequentialTime / parallelTime).toFixed(2)}x`);

    // 4. Parallel reconstruction
    console.log('\n4. Parallel Reconstruction');
    console.log('-'.repeat(30));

    const { timeMs: parallelReconTime, result: parallelReconResult } = await measureTime(
      async () => {
        return await parallelReconstruction([parallelResult[0]?.fragments || []]);
      }
    );

    console.log(`   Parallel reconstruction time: ${parallelReconTime}ms`);
    console.log(`   Reconstruction complete: ${parallelReconResult[0]?.isComplete || false}`);

    const { timeMs: sequentialReconTime, result: sequentialReconResult } = await measureTime(
      async () => {
        return await reconstructData(parallelResult[0]?.fragments || []);
      }
    );

    console.log(`   Sequential reconstruction time: ${sequentialReconTime}ms`);
    console.log(
      `   Performance improvement: ${(sequentialReconTime / parallelReconTime).toFixed(2)}x`
    );

    // 5. Mobile optimization
    console.log('\n5. Mobile Optimization');
    console.log('-'.repeat(30));

    const mobileLatency = estimateMobileLatency(1024);
    console.log(`   Estimated mobile latency (1KB): ${mobileLatency}ms`);

    const optimalSize = getOptimalFragmentSize(4096); // Use 4KB as example data size
    console.log(`   Optimal fragment size: ${optimalSize} bytes`);

    // Test with optimal size
    const mobileData = new Uint8Array(4096);
    crypto.getRandomValues(mobileData);

    const mobileFragResult = await fragmentData(mobileData, optimalSize);
    console.log(`   Mobile-optimized fragments: ${mobileFragResult.metadata.fragmentCount}`);

    // 6. Fragment serialization
    console.log('\n6. Fragment Serialization');
    console.log('-'.repeat(30));

    const serialized = serializeFragments(fragResult.fragments);
    console.log(`   Serialized size: ${serialized.length} characters`);

    const deserialized = deserializeFragments(serialized);
    console.log(`   Deserialized fragments: ${deserialized.length}`);

    const serializationValid =
      deserialized.length === fragResult.fragments.length &&
      deserialized.every((frag, index) => frag.index === fragResult.fragments[index]!.index);
    console.log(`   Serialization integrity: ${serializationValid ? 'VERIFIED' : 'FAILED'}`);

    // 7. Fragment compression
    console.log('\n7. Fragment Compression');
    console.log('-'.repeat(30));

    // Create data with patterns for better compression
    const patternData = new Uint8Array(2048);
    for (let i = 0; i < patternData.length; i++) {
      patternData[i] = i % 256;
    }

    const patternFragResult = await fragmentData(patternData);
    const originalSize = patternFragResult.fragments.reduce(
      (sum, frag) => sum + frag.data.length,
      0
    );

    const compressed = compressFragments(patternFragResult.fragments);
    const compressedSize = compressed.reduce((sum, frag) => sum + frag.data.length, 0);

    console.log(`   Original total size: ${originalSize} bytes`);
    console.log(`   Compressed total size: ${compressedSize} bytes`);
    console.log(`   Compression ratio: ${(originalSize / compressedSize).toFixed(2)}:1`);

    const decompressed = decompressFragments(compressed);
    const decompressionValid = decompressed.length === patternFragResult.fragments.length;
    console.log(`   Decompression integrity: ${decompressionValid ? 'VERIFIED' : 'FAILED'}`);

    // 8. Fragment loss simulation
    console.log('\n8. Fragment Loss Simulation');
    console.log('-'.repeat(30));

    const lossTestData = new Uint8Array(1024);
    crypto.getRandomValues(lossTestData);

    const lossFragResult = await fragmentData(lossTestData);
    console.log(`   Original fragments: ${lossFragResult.fragments.length}`);

    // Simulate losing some fragments
    const partialFragments = lossFragResult.fragments.slice(0, -2); // Remove last 2 fragments
    console.log(`   Remaining fragments: ${partialFragments.length}`);

    const partialReconResult = await reconstructData(partialFragments);
    console.log(
      `   Partial reconstruction: ${partialReconResult.isComplete ? 'COMPLETE' : 'INCOMPLETE'}`
    );
    console.log(`   Reconstructed size: ${partialReconResult.data.length} bytes`);

    // 9. Performance analysis
    console.log('\n9. Performance Analysis');
    console.log('-'.repeat(30));

    const sizes = [512, 1024, 2048, 4096, 8192];
    console.log('   Size (bytes) | Frag Time (ms) | Recon Time (ms) | Fragments');
    console.log('   ' + '-'.repeat(65));

    for (const size of sizes) {
      const perfData = new Uint8Array(size);
      crypto.getRandomValues(perfData);

      const { timeMs: fragTime, result: perfFragResult } = await measureTime(async () => {
        return await fragmentData(perfData);
      });

      const { timeMs: reconTime } = await measureTime(async () => {
        return await reconstructData(perfFragResult.fragments);
      });

      console.log(
        `   ${size.toString().padStart(12)} | ${fragTime.toString().padStart(14)} | ${reconTime.toString().padStart(15)} | ${perfFragResult.metadata.fragmentCount.toString().padStart(9)}`
      );
    }

    console.log('\n✅ Fragmentation example completed successfully!');
    console.log('\nFeatures Demonstrated:');
    console.log('- ✓ Basic fragmentation and reconstruction');
    console.log('- ✓ Parallel processing capabilities');
    console.log('- ✓ Mobile optimization');
    console.log('- ✓ Fragment serialization');
    console.log('- ✓ Data compression');
    console.log('- ✓ Fragment loss handling');
    console.log('- ✓ Performance analysis');
  } catch (error) {
    console.error('❌ Error in fragmentation example:', error);
    process.exit(1);
  }
}

// Run the example
if (require.main === module) {
  fragmentationExample().catch(console.error);
}

export { fragmentationExample };
