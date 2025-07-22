//! Performance Benchmark for TOPAY-Z512 Optimizations
//! 
//! This example demonstrates the performance improvements achieved through optimization.

use topayz512::*;
use std::time::{Duration, Instant};

#[cfg(feature = "fragmentation")]
use topayz512::fragment::{FragmentEngine, Fragment};

fn main() {
    println!("=== TOPAY-Z512 Performance Benchmark ===\n");
    
    // Benchmark hash operations
    benchmark_hash_operations();
    
    // Benchmark key operations
    benchmark_key_operations();
    
    // Benchmark KEM operations
    benchmark_kem_operations();
    
    #[cfg(feature = "fragmentation")]
    benchmark_fragmentation_operations();
    
    // Feature detection
    benchmark_feature_detection();
    
    println!("\n=== Benchmark Complete ===");
}

fn benchmark_hash_operations() {
    println!("1. Hash Operations Benchmark:");
    
    let data_sizes = [64, 256, 1024, 4096, 16384];
    
    for &size in &data_sizes {
        let data = vec![42u8; size];
        
        let start = Instant::now();
        let iterations = 10000;
        
        for _ in 0..iterations {
            let _hash = Hash::new(&data);
        }
        
        let elapsed = start.elapsed();
        let per_op = elapsed / iterations;
        let throughput = (size as f64 * iterations as f64) / elapsed.as_secs_f64() / 1_000_000.0;
        
        println!("   {} bytes: {:?} per hash, {:.2} MB/s throughput", 
                 size, per_op, throughput);
    }
    
    // Test hex operations
    let data = vec![42u8; 1024];
    let hash = Hash::new(&data);
    
    let start = Instant::now();
    let iterations = 100000;
    for _ in 0..iterations {
        let _hex = hash.to_hex();
    }
    let hex_time = start.elapsed() / iterations;
    
    let hex_str = hash.to_hex();
    let start = Instant::now();
    for _ in 0..iterations {
        let _hash = Hash::from_hex(&hex_str).unwrap();
    }
    let from_hex_time = start.elapsed() / iterations;
    
    println!("   Hex encoding: {:?} per operation", hex_time);
    println!("   Hex decoding: {:?} per operation", from_hex_time);
    println!();
}

fn benchmark_key_operations() {
    println!("2. Key Operations Benchmark:");
    
    // Single key generation
    let start = Instant::now();
    let iterations = 10000;
    for _ in 0..iterations {
        let _keypair = KeyPair::generate();
    }
    let single_time = start.elapsed() / iterations;
    
    // Batch key generation
    let start = Instant::now();
    let batch_size = 100;
    let batches = 100;
    for _ in 0..batches {
        let _keypairs = KeyPair::batch_generate(batch_size);
    }
    let batch_time = start.elapsed() / (batches * batch_size) as u32;
    
    println!("   Single key generation: {:?} per keypair", single_time);
    println!("   Batch key generation: {:?} per keypair", batch_time);
    println!("   Batch improvement: {:.2}x faster", 
             single_time.as_nanos() as f64 / batch_time.as_nanos() as f64);
    
    // Public key derivation
    let private_key = PrivateKey::generate();
    let start = Instant::now();
    let iterations = 100000;
    for _ in 0..iterations {
        let _public_key = private_key.public_key();
    }
    let derivation_time = start.elapsed() / iterations;
    
    println!("   Public key derivation: {:?} per operation", derivation_time);
    println!();
}

fn benchmark_kem_operations() {
    println!("3. KEM Operations Benchmark:");
    
    // Key generation
    let start = Instant::now();
    let iterations = 1000;
    for _ in 0..iterations {
        let _keys = Kem::keygen();
    }
    let keygen_time = start.elapsed() / iterations;
    
    // Encapsulation
    let (public_key, _) = Kem::keygen();
    let start = Instant::now();
    for _ in 0..iterations {
        let _result = Kem::encapsulate(&public_key);
    }
    let encap_time = start.elapsed() / iterations;
    
    // Decapsulation
    let (public_key, secret_key) = Kem::keygen();
    let (ciphertext, _) = Kem::encapsulate(&public_key);
    let start = Instant::now();
    for _ in 0..iterations {
        let _shared_secret = Kem::decapsulate(&secret_key, &ciphertext);
    }
    let decap_time = start.elapsed() / iterations;
    
    // Batch operations
    let start = Instant::now();
    let batch_size = 10;
    let batches = 100;
    for _ in 0..batches {
        let _keys = Kem::batch_keygen(batch_size);
    }
    let batch_keygen_time = start.elapsed() / (batches * batch_size) as u32;
    
    println!("   KEM keygen: {:?} per operation", keygen_time);
    println!("   KEM encapsulate: {:?} per operation", encap_time);
    println!("   KEM decapsulate: {:?} per operation", decap_time);
    println!("   Batch keygen: {:?} per operation", batch_keygen_time);
    println!("   Batch improvement: {:.2}x faster", 
             keygen_time.as_nanos() as f64 / batch_keygen_time.as_nanos() as f64);
    println!();
}

#[cfg(feature = "fragmentation")]
fn benchmark_fragmentation_operations() {
    println!("4. Fragmentation Operations Benchmark:");
    
    let data_sizes = [1024, 4096, 16384, 65536];
    
    for &size in &data_sizes {
        let data = vec![42u8; size];
        
        // Fragmentation
        let start = Instant::now();
        let iterations = 1000;
        for _ in 0..iterations {
            let _fragments = FragmentEngine::fragment_data(&data).unwrap();
        }
        let frag_time = start.elapsed() / iterations;
        
        // Reconstruction
        let fragments = FragmentEngine::fragment_data(&data).unwrap();
        let start = Instant::now();
        for _ in 0..iterations {
            let _reconstructed = FragmentEngine::reconstruct_data(&fragments).unwrap();
        }
        let recon_time = start.elapsed() / iterations;
        
        // Serialization
        let fragment = &fragments[0];
        let start = Instant::now();
        for _ in 0..iterations {
            let _bytes = fragment.to_bytes();
        }
        let ser_time = start.elapsed() / iterations;
        
        // Deserialization
        let bytes = fragment.to_bytes();
        let start = Instant::now();
        for _ in 0..iterations {
            let _fragment = Fragment::from_bytes(&bytes).unwrap();
        }
        let deser_time = start.elapsed() / iterations;
        
        let throughput = (size as f64 * iterations as f64) / (frag_time + recon_time).as_secs_f64() / 1_000_000.0;
        
        println!("   {} bytes: frag={:?}, recon={:?}, ser={:?}, deser={:?}, {:.2} MB/s", 
                 size, frag_time, recon_time, ser_time, deser_time, throughput);
    }
    
    // Mobile latency estimation
    println!("   Mobile latency estimates:");
    for &size in &[1000, 5000, 10000, 50000] {
        let latency = FragmentEngine::estimate_mobile_latency(size);
        println!("     {} bytes: {}ms", size, latency);
    }
    println!();
}

fn benchmark_feature_detection() {
    println!("5. System Capabilities:");
    
    let start = Instant::now();
    let has_simd = features::has_simd();
    let simd_time = start.elapsed();
    
    let start = Instant::now();
    let has_hw_rng = features::has_hardware_rng();
    let rng_time = start.elapsed();
    
    let start = Instant::now();
    let thread_count = features::optimal_thread_count();
    let thread_time = start.elapsed();
    
    println!("   SIMD support: {} (detected in {:?})", has_simd, simd_time);
    println!("   Hardware RNG: {} (detected in {:?})", has_hw_rng, rng_time);
    println!("   Optimal threads: {} (detected in {:?})", thread_count, thread_time);
    
    // Utility function benchmarks
    let data1 = [1, 2, 3, 4, 5];
    let data2 = [1, 2, 3, 4, 5];
    let data3 = [1, 2, 3, 4, 6];
    
    let start = Instant::now();
    let iterations = 1000000;
    for _ in 0..iterations {
        let _eq = utils::constant_time_eq(&data1, &data2);
    }
    let eq_time = start.elapsed() / iterations;
    
    let start = Instant::now();
    for _ in 0..iterations {
        let _eq = utils::constant_time_eq(&data1, &data3);
    }
    let neq_time = start.elapsed() / iterations;
    
    println!("   Constant-time equality (equal): {:?}", eq_time);
    println!("   Constant-time equality (different): {:?}", neq_time);
    
    // Hex encoding benchmark
    let data = [0x12, 0x34, 0xab, 0xcd, 0xef];
    let start = Instant::now();
    for _ in 0..iterations {
        let _hex = utils::fast_hex_encode(&data);
    }
    let hex_encode_time = start.elapsed() / iterations;
    
    let hex_str = "1234abcdef";
    let start = Instant::now();
    for _ in 0..iterations {
        let _data = utils::fast_hex_decode(hex_str).unwrap();
    }
    let hex_decode_time = start.elapsed() / iterations;
    
    println!("   Fast hex encode: {:?}", hex_encode_time);
    println!("   Fast hex decode: {:?}", hex_decode_time);
    println!();
}