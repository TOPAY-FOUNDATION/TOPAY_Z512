# pqcrypto512 Monorepo

**pqcrypto512** is a multi-language reference implementation of a post-quantum 512-bit encryption library. It includes idiomatic SDKs in Rust, Go, and JavaScript/TypeScript.

## Repository Structure

```base tree
pqcrypto512/                  # Root monorepo
├── LICENSE                    # Apache-2.0 License
├── README.md                  # This file
├── docs/                      # Shared design spec & API references
│   ├── design_spec.md
│   └── api_reference.md
├── ci/                        # CI workflows for each language
│   ├── rust.yml
│   ├── go.yml
│   └── js.yml
├── test-vectors/              # Canonical test vectors & KATs
│   └── basic.json
├── rust/                      # Rust implementation (Cargo project)
│   └── ...
├── go/                        # Go implementation (Go module)
│   └── ...
└── js/                        # JavaScript/TypeScript implementation (npm package)
    └── ...
```

## Getting Started

### Prerequisites

* **Rust** (stable toolchain) for the `rust/` directory
* **Go** (>=1.18) for the `go/` directory
* **Node.js** (>=14) and **npm**/yarn for the `js/` directory

### Build & Test All

From the project root, you can run each language’s CI script manually:

```bash
# Rust
cd rust && cargo test && cargo bench

# Go
cd go && go test ./... && go test -bench=.

# JS/TS
cd js && npm install && npm test && npm run bench
```

Or rely on GitHub Actions workflows under `.github/workflows/` which automate these steps on push/PR.

## Documentation

For detailed design rationale and API usage, see the shared docs:

* [Design Specification](docs/design_spec.md)
* [API Reference](docs/api_reference.md)

Each language folder also contains its own README with language-specific examples, installation, and usage instructions.

## Contributing

Please read `docs/contributing.md` for guidelines on:

* Coding style and linting
* Test vector updates (sync across languages)
* Pull request process and review criteria

## Licensing

This project is licensed under the Apache-2.0 License. See [LICENSE](LICENSE) for full terms.
