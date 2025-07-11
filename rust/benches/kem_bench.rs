use criterion::{black_box, criterion_group, criterion_main, Criterion};
use topayz512::{decapsulate, encapsulate, keygen};

#[cfg(feature = "fragmentation")]
use topayz512::fragmented;

pub fn kem_benchmark(c: &mut Criterion) {
    // Benchmark key generation
    c.bench_function("keygen", |b| {
        b.iter(|| {
            black_box(keygen().unwrap());
        });
    });

    // Generate a key pair for encapsulation and decapsulation benchmarks
    let (pk, sk) = keygen().unwrap();

    // Benchmark encapsulation
    c.bench_function("encapsulate", |b| {
        b.iter(|| {
            black_box(encapsulate(black_box(&pk)).unwrap());
        });
    });

    // Generate a ciphertext for decapsulation benchmark
    let (ct, _) = encapsulate(&pk).unwrap();

    // Benchmark decapsulation
    c.bench_function("decapsulate", |b| {
        b.iter(|| {
            black_box(decapsulate(black_box(&sk), black_box(&ct)).unwrap());
        });
    });

    // Benchmark fragmented operations if the feature is enabled
    #[cfg(feature = "fragmentation")]
    {
        // Number of fragments for benchmarking
        let num_fragments = 4;

        // Benchmark fragmented key generation
        c.bench_function("fragmented_keygen", |b| {
            b.iter(|| {
                black_box(fragmented::keygen(black_box(num_fragments)).unwrap());
            });
        });

        // Generate a fragmented key pair for encapsulation and decapsulation benchmarks
        let (fpk, fsk) = fragmented::keygen(num_fragments).unwrap();

        // Benchmark fragmented encapsulation
        c.bench_function("fragmented_encapsulate", |b| {
            b.iter(|| {
                black_box(fragmented::encapsulate(black_box(&fpk)).unwrap());
            });
        });

        // Generate a fragmented ciphertext for decapsulation benchmark
        let (fct, _) = fragmented::encapsulate(&fpk).unwrap();

        // Benchmark fragmented decapsulation
        c.bench_function("fragmented_decapsulate", |b| {
            b.iter(|| {
                black_box(fragmented::decapsulate(black_box(&fsk), black_box(&fct)).unwrap());
            });
        });
    }
}

criterion_group!(benches, kem_benchmark);
criterion_main!(benches);