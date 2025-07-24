/**
 * Data fragmentation and reconstruction for TOPAY-Z512
 */

import {
  FRAGMENT_SIZE,
  MIN_FRAGMENT_THRESHOLD,
  MAX_FRAGMENTS,
  Fragment,
  FragmentMetadata,
  FragmentationResult,
  ReconstructionResult,
  EmptyDataError,
  FragmentationFailedError,
  ReconstructionFailedError
  // InvalidFragmentCountError
} from './index';
import { computeHash } from './hash';
import { timestamp, /* copyBytes, */ constantTimeEqual } from './utils';

/**
 * Fragments data into smaller chunks for parallel processing
 * @param data - Data to fragment
 * @param fragmentSize - Size of each fragment (default: FRAGMENT_SIZE)
 * @returns Promise resolving to fragmentation result
 * @throws EmptyDataError if data is empty
 * @throws FragmentationFailedError if fragmentation fails
 */
export async function fragmentData(
  data: Uint8Array,
  fragmentSize: number = FRAGMENT_SIZE
): Promise<FragmentationResult> {
  if (data.length === 0) {
    throw new EmptyDataError();
  }

  if (fragmentSize <= 0) {
    throw new Error('Fragment size must be positive');
  }

  // Calculate number of fragments needed
  const fragmentCount = Math.ceil(data.length / fragmentSize);

  if (fragmentCount > MAX_FRAGMENTS) {
    throw new FragmentationFailedError();
  }

  // Create metadata
  const checksum = await computeHash(data);
  const metadata: FragmentMetadata = {
    originalSize: data.length,
    fragmentCount,
    timestamp: timestamp(),
    algorithm: 'TOPAY-Z512-FRAG-v1',
    checksum
  };

  // Create fragments
  const fragments: Fragment[] = [];

  for (let i = 0; i < fragmentCount; i++) {
    const start = i * fragmentSize;
    const end = Math.min(start + fragmentSize, data.length);
    const fragmentData = data.slice(start, end);

    fragments.push({
      index: i,
      data: fragmentData,
      metadata
    });
  }

  return {
    fragments,
    metadata
  };
}

/**
 * Reconstructs data from fragments
 * @param fragments - Array of fragments to reconstruct
 * @returns Promise resolving to reconstruction result
 * @throws ReconstructionFailedError if reconstruction fails
 */
export async function reconstructData(fragments: Fragment[]): Promise<ReconstructionResult> {
  if (fragments.length === 0) {
    throw new ReconstructionFailedError();
  }

  // Get metadata from first fragment
  const metadata = fragments[0]!.metadata;
  const expectedCount = metadata.fragmentCount;

  // Check if we have all fragments
  const receivedIndices = new Set(fragments.map(f => f.index));
  const missingCount = expectedCount - receivedIndices.size;
  const isComplete = missingCount === 0;

  if (!isComplete) {
    return {
      data: new Uint8Array(0),
      isComplete: false,
      missingCount,
      metadata
    };
  }

  // Sort fragments by index
  const sortedFragments = [...fragments].sort((a, b) => a.index - b.index);

  // Validate fragment sequence
  for (let i = 0; i < sortedFragments.length; i++) {
    if (sortedFragments[i]!.index !== i) {
      throw new ReconstructionFailedError();
    }
  }

  // Reconstruct data
  const reconstructedData = new Uint8Array(metadata.originalSize);
  let offset = 0;

  for (const fragment of sortedFragments) {
    const fragmentData = fragment.data;
    const copyLength = Math.min(fragmentData.length, metadata.originalSize - offset);

    reconstructedData.set(fragmentData.slice(0, copyLength), offset);
    offset += copyLength;
  }

  // Verify checksum
  const computedChecksum = await computeHash(reconstructedData);
  if (!constantTimeEqual(computedChecksum, metadata.checksum)) {
    throw new ReconstructionFailedError();
  }

  return {
    data: reconstructedData,
    isComplete: true,
    missingCount: 0,
    metadata
  };
}

/**
 * Performs parallel fragmentation of multiple data items
 * @param dataItems - Array of data to fragment
 * @param fragmentSize - Size of each fragment
 * @returns Promise resolving to array of fragmentation results
 */
export async function parallelFragmentation(
  dataItems: Uint8Array[],
  fragmentSize: number = FRAGMENT_SIZE
): Promise<FragmentationResult[]> {
  const promises = dataItems.map(data => fragmentData(data, fragmentSize));
  return Promise.all(promises);
}

/**
 * Performs parallel reconstruction of multiple fragment sets
 * @param fragmentSets - Array of fragment arrays to reconstruct
 * @returns Promise resolving to array of reconstruction results
 */
export async function parallelReconstruction(
  fragmentSets: Fragment[][]
): Promise<ReconstructionResult[]> {
  const promises = fragmentSets.map(fragments => reconstructData(fragments));
  return Promise.all(promises);
}

/**
 * Validates fragment integrity
 * @param fragment - Fragment to validate
 * @returns Promise resolving to true if fragment is valid
 */
export async function validateFragment(fragment: Fragment): Promise<boolean> {
  try {
    // Check basic structure
    if (!fragment.data || !fragment.metadata) {
      return false;
    }

    // Check index bounds
    if (fragment.index < 0 || fragment.index >= fragment.metadata.fragmentCount) {
      return false;
    }

    // Check algorithm
    if (fragment.metadata.algorithm !== 'TOPAY-Z512-FRAG-v1') {
      return false;
    }

    return true;
  } catch {
    return false;
  }
}

/**
 * Estimates mobile device latency for fragment processing
 * @param dataSize - Size of data to process
 * @param fragmentSize - Size of each fragment
 * @returns Estimated latency in milliseconds
 */
export function estimateMobileLatency(
  dataSize: number,
  fragmentSize: number = FRAGMENT_SIZE
): number {
  const fragmentCount = Math.ceil(dataSize / fragmentSize);

  // Base latency estimates for mobile devices (in ms)
  const baseFragmentationLatency = 2; // Per fragment
  const baseReconstructionLatency = 1.5; // Per fragment
  const networkLatency = 50; // Network overhead

  const fragmentationTime = fragmentCount * baseFragmentationLatency;
  const reconstructionTime = fragmentCount * baseReconstructionLatency;

  return fragmentationTime + reconstructionTime + networkLatency;
}

/**
 * Determines optimal fragment size for given data and constraints
 * @param dataSize - Size of data to fragment
 * @param maxLatency - Maximum acceptable latency in milliseconds
 * @param deviceType - Type of device ('mobile', 'desktop', 'iot')
 * @returns Optimal fragment size
 */
export function getOptimalFragmentSize(
  dataSize: number,
  maxLatency: number = 100,
  deviceType: 'mobile' | 'desktop' | 'iot' = 'mobile'
): number {
  // Device-specific constraints
  const constraints = {
    mobile: { minSize: 128, maxSize: 512, processingPower: 1.0 },
    desktop: { minSize: 256, maxSize: 2048, processingPower: 3.0 },
    iot: { minSize: 64, maxSize: 256, processingPower: 0.5 }
  };

  const constraint = constraints[deviceType];

  // Start with default fragment size
  let fragmentSize = FRAGMENT_SIZE;

  // Adjust based on data size
  if (dataSize < MIN_FRAGMENT_THRESHOLD) {
    return Math.min(dataSize, constraint.maxSize);
  }

  // Adjust based on latency requirements
  let estimatedLatency = estimateMobileLatency(dataSize, fragmentSize);

  while (estimatedLatency > maxLatency && fragmentSize > constraint.minSize) {
    fragmentSize = Math.max(fragmentSize * 0.8, constraint.minSize);
    estimatedLatency = estimateMobileLatency(dataSize, fragmentSize);
  }

  // Ensure within device constraints
  fragmentSize = Math.max(constraint.minSize, Math.min(fragmentSize, constraint.maxSize));

  return Math.floor(fragmentSize);
}

/**
 * Serializes fragments to a portable format
 * @param fragments - Fragments to serialize
 * @returns Serialized fragments as JSON string
 */
export function serializeFragments(fragments: Fragment[]): string {
  const serializable = fragments.map(fragment => ({
    index: fragment.index,
    data: Array.from(fragment.data),
    metadata: {
      ...fragment.metadata,
      checksum: Array.from(fragment.metadata.checksum)
    }
  }));

  return JSON.stringify(serializable);
}

/**
 * Deserializes fragments from a portable format
 * @param serialized - Serialized fragments JSON string
 * @returns Deserialized fragments
 * @throws Error if deserialization fails
 */
export function deserializeFragments(serialized: string): Fragment[] {
  try {
    const data = JSON.parse(serialized);

    if (!Array.isArray(data)) {
      throw new Error('Invalid fragment data format');
    }

    return data.map(item => ({
      index: item.index,
      data: new Uint8Array(item.data),
      metadata: {
        ...item.metadata,
        checksum: new Uint8Array(item.metadata.checksum)
      }
    }));
  } catch (error) {
    throw new Error(`Failed to deserialize fragments: ${error}`);
  }
}

/**
 * Compresses fragments using simple run-length encoding
 * @param fragments - Fragments to compress
 * @returns Compressed fragments
 */
export function compressFragments(fragments: Fragment[]): Fragment[] {
  return fragments.map(fragment => ({
    ...fragment,
    data: runLengthEncode(fragment.data)
  }));
}

/**
 * Decompresses fragments using run-length decoding
 * @param fragments - Compressed fragments to decompress
 * @returns Decompressed fragments
 */
export function decompressFragments(fragments: Fragment[]): Fragment[] {
  return fragments.map(fragment => ({
    ...fragment,
    data: runLengthDecode(fragment.data)
  }));
}

// Helper functions for compression

function runLengthEncode(data: Uint8Array): Uint8Array {
  if (data.length === 0) return data;

  const encoded: number[] = [];
  let current = data[0]!;
  let count = 1;

  for (let i = 1; i < data.length; i++) {
    if (data[i] === current && count < 255) {
      count++;
    } else {
      encoded.push(count, current);
      current = data[i]!;
      count = 1;
    }
  }

  encoded.push(count, current);
  return new Uint8Array(encoded);
}

function runLengthDecode(encoded: Uint8Array): Uint8Array {
  if (encoded.length % 2 !== 0) {
    throw new Error('Invalid run-length encoded data');
  }

  const decoded: number[] = [];

  for (let i = 0; i < encoded.length; i += 2) {
    const count = encoded[i]!;
    const value = encoded[i + 1]!;

    for (let j = 0; j < count; j++) {
      decoded.push(value);
    }
  }

  return new Uint8Array(decoded);
}
