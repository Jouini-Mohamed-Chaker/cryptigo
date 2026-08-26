# Cryptigo

A simple, fast CLI tool for encrypting and decrypting files using AES-256-GCM.

> **Note:** All source code (`.go` files) in this repository is 100% hand-written. Only Markdown (`.md`) files were written with AI assistance.

---

## Features

* **AES-256-GCM:** Authenticated symmetric encryption to ensure data secrecy and detect file tampering.
* **Key Derivation (PBKDF2):** Derives secure 256-bit keys from passphrases using unique per-file salts.
* **Stream I/O:** Memory-efficient chunk processing for handling large files safely.

---

## Project Structure


```
cryptigo/
├── Project-Plan.md       # Implementation notes & specs
├── cmd/
│   └── cryptigo/         # CLI entry point (main.go)
├── go.mod                # Go module definition
└── internal/
    └── crypt/            # Encrypt, decrypt, and validation logic
```

---

## Usage

### Build

```bash
go build -o cryptigo ./cmd/cryptigo
```

### Encryption

```bash
./cryptigo encrypt -in secret.pdf -out secret.pdf.enc
```

### Decryption

```bash
./cryptigo decrypt -in secret.pdf.enc -out secret.pdf
```