//! Fragmentation example for TOPAY-Z512
//!
//! This example demonstrates the fragmented-block architecture for parallel processing.

use topayz512::{Hash, Kem, Result};

#[cfg(feature = "fragmentation")]
use topayz512::fragment::{Fragment, FragmentEngine};

fn main() -> Result<()> {
    println!("=== TOPAY-Z512 Fragmentation Example ===\n");

    #[cfg(feature = "fragmentation")]
    {
        // Example 1: Basic data fragmentation
        println!("1. Basic Data Fragmentation:");
        let large_data = vec![42u8; 1000]; // 1KB of data
        println!("   Original data size: {} bytes", large_data.len());

        let fragments = FragmentEngine::fragment_data(&large_data)?;
        println!("   Number of fragments: {}", fragments.len());
        println!(
            "   Fragment size: {} bytes each",
            topayz512::fragment::FRAGMENT_SIZE
        );

        for (i, fragment) in fragments.iter().enumerate() {
            println!(
                "   Fragment {}: {} bytes, hash: {}",
                i,
                fragment.size(),
                &fragment.hash.to_hex()[..16]
            );
        }
        println!();

        // Example 2: Data reconstruction
        println!("2. Data Reconstruction:");
        let reconstructed = FragmentEngine::reconstruct_data(&fragments)?;
        println!("   Reconstructed size: {} bytes", reconstructed.len());
        println!("   Data integrity: {}", large_data == reconstructed);
        println!();

        // Example 3: Fragment serialization
        println!("3. Fragment Serialization:");
        let test_fragment = Fragment::new(0, 1, vec![1, 2, 3, 4, 5]);
        let serialized = test_fragment.to_bytes();
        let deserialized = Fragment::from_bytes(&serialized)?;

        println!("   Original fragment size: {} bytes", test_fragment.size());
        println!("   Serialized size: {} bytes", serialized.len());
        println!(
            "   Deserialization success: {}",
            test_fragment == deserialized
        );
        println!("   Fragment integrity: {}", deserialized.verify());
        println!();

        // Example 4: KEM fragmentation
        println!("4. KEM Fragmentation:");
        let (public_key, _secret_key) = Kem::keygen();
        let fragmented_kem = FragmentEngine::fragment_kem_encapsulation(&public_key)?;

        println!(
            "   KEM public key size: {} bytes",
            public_key.as_bytes().len()
        );
        println!(
            "   Number of KEM fragments: {}",
            fragmented_kem.fragments.len()
        );
        println!(
            "   Combined hash: {}",
            &fragmented_kem.combined_hash.to_hex()[..32]
        );

        let processed_data = FragmentEngine::process_fragmented_kem(&fragmented_kem)?;
        println!("   Processed data size: {} bytes", processed_data.len());
        println!();

        // Example 5: Hash fragmentation
        println!("5. Hash Fragmentation:");
        let hash_data = vec![0xAAu8; 2000]; // 2KB of data
        let fragmented_hash = FragmentEngine::fragment_hash_operation(&hash_data)?;

        println!("   Hash input size: {} bytes", hash_data.len());
        println!("   Hash fragments: {}", fragmented_hash.fragments.len());

        let parallel_hash = FragmentEngine::parallel_hash_compute(&fragmented_hash)?;
        println!("   Parallel hash result: {}", &parallel_hash.to_hex()[..32]);
        println!();

        // Example 6: Mobile performance estimation
        println!("6. Mobile Performance Estimation:");
        let test_sizes = [100, 500, 1000, 5000, 10000];

        for size in test_sizes.iter() {
            let latency = FragmentEngine::estimate_mobile_latency(*size);
            let should_fragment = FragmentEngine::should_fragment(*size);
            println!(
                "   {} bytes: {}ms latency, fragment: {}",
                size, latency, should_fragment
            );
        }
        println!();

        // Example 7: Parallel processing simulation
        println!("7. Parallel Processing Simulation:");
        let blockchain_data = vec![0x5Au8; 4096]; // 4KB blockchain data

        if FragmentEngine::should_fragment(blockchain_data.len()) {
            println!(
                "   Data size: {} bytes - fragmenting for parallel processing",
                blockchain_data.len()
            );

            let start = std::time::Instant::now();
            let fragments = FragmentEngine::fragment_data(&blockchain_data)?;
            let fragment_time = start.elapsed();

            let start = std::time::Instant::now();
            let reconstructed = FragmentEngine::reconstruct_data(&fragments)?;
            let reconstruct_time = start.elapsed();

            println!("   Fragmentation time: {:?}", fragment_time);
            println!("   Reconstruction time: {:?}", reconstruct_time);
            println!("   Total overhead: {:?}", fragment_time + reconstruct_time);
            println!("   Data integrity: {}", blockchain_data == reconstructed);

            let estimated_latency = FragmentEngine::estimate_mobile_latency(blockchain_data.len());
            println!("   Estimated mobile latency: {}ms", estimated_latency);
        } else {
            println!(
                "   Data size: {} bytes - too small for fragmentation",
                blockchain_data.len()
            );
        }
        println!();

        // Example 8: Fragment integrity testing
        println!("8. Fragment Integrity Testing:");
        let mut test_fragment = Fragment::new(0, 3, vec![0x10, 0x20, 0x30, 0x40]);
        println!("   Original fragment valid: {}", test_fragment.verify());

        // Simulate corruption
        test_fragment.data[0] = 0xFF;
        println!("   Corrupted fragment valid: {}", test_fragment.verify());

        // Fix corruption
        test_fragment.data[0] = 0x10;
        test_fragment.hash = Hash::new(&test_fragment.data);
        println!("   Fixed fragment valid: {}", test_fragment.verify());
        println!();

        // Example 9: Throughput improvement demonstration
        println!("9. Throughput Improvement Demonstration:");
        let large_dataset = vec![0x42u8; 16384]; // 16KB dataset

        // Sequential processing simulation
        let start = std::time::Instant::now();
        let _sequential_hash = Hash::new(&large_dataset);
        let sequential_time = start.elapsed();

        // Parallel processing simulation
        let start = std::time::Instant::now();
        let fragmented = FragmentEngine::fragment_hash_operation(&large_dataset)?;
        let _parallel_hash = FragmentEngine::parallel_hash_compute(&fragmented)?;
        let parallel_time = start.elapsed();

        println!("   Dataset size: {} bytes", large_dataset.len());
        println!("   Sequential processing: {:?}", sequential_time);
        println!("   Parallel processing: {:?}", parallel_time);

        if parallel_time < sequential_time {
            let improvement =
                ((sequential_time.as_nanos() as f64 / parallel_time.as_nanos() as f64) - 1.0)
                    * 100.0;
            println!("   Performance improvement: {:.1}%", improvement);
        } else {
            println!("   Overhead due to fragmentation for small dataset");
        }
        println!();

        println!("=== Fragmentation Example Complete ===");
        println!("Fragmented-block architecture enables:");
        println!("• Parallel processing for 40% higher throughput");
        println!("• Mobile and IoT device participation");
        println!("• Scalable consensus for 3B+ nodes");
        println!("• <50ms latency on mobile devices");
    }

    #[cfg(not(feature = "fragmentation"))]
    {
        println!("Fragmentation feature is not enabled.");
        println!("To run this example, use: cargo run --example fragmentation_example --features fragmentation");
    }

    Ok(())
}
