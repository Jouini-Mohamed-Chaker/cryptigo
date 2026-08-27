> **Note:** .md files are ai generated to guide me because I'm a beginner

## **File Encrypter/Decrypter CLI** 

### Core Requirements & Features

**1. Symmetric Encryption (AES-256-GCM)**

* **AES-256:** Standard, fast, and virtually unbreakable when implemented correctly.
* **GCM (Galois/Counter Mode):** Provides **authenticated encryption**. This is crucial: it not only hides the data, but it also creates an authentication tag. If someone tries to tamper with the encrypted file (or inputs the wrong password), decryption will fail explicitly instead of producing corrupted garbage output.

**2. Key Derivation (PBKDF2 or Argon2)**

* Passwords typed by humans are low-entropy. You cannot use a raw password string as an AES key directly.
* Use a Key Derivation Function (KDF) to stretch the user's password into a secure 32-byte (256-bit) key.
* Incorporate a unique, randomly generated **Salt** per file to protect against rainbow table attacks.

**3. Stream / Chunked File I/O**

* Avoid reading an entire 4 GB video file into RAM all at once.
* Use Go's `io.Reader` and `io.Writer` interfaces to read source files in memory-efficient chunks (e.g., 64 KB buffers), encrypting and writing them progressively to the target file.

**4. Binary Layout / File Format Specification**
Your encrypted output file (`.enc`) needs a predictable header so your tool knows how to decrypt it later. A standard sequence stored at the start of the file includes:

* **Salt** (e.g., 16 bytes) – used for key derivation.
* **Nonce / Initialization Vector (IV)** (12 bytes for AES-GCM) – ensures identical plaintexts generate completely different ciphertexts every time.
* **Ciphertext Data** – the actual encrypted bytes.
* **Authentication Tag** – automatically managed by GCM to verify file integrity.

---

### Key Deliverables

* **The Binary CLI App:** A standalone executable compiled with `go build` (e.g., `filecrypt`).
* **Subcommand CLI Interface:**
* **Encryption mode:** `filecrypt encrypt -in secret.pdf -out secret.pdf.enc`
* **Decryption mode:** `filecrypt decrypt -in secret.pdf.enc -out secret.pdf`


* **Secure Password Prompting:** Standard terminal input masking (hiding password characters as the user types) so passwords aren't exposed on screen or recorded in terminal command history.
* **Overwriting / Destruction Flag (Optional):** A `-shred` or `-rm` flag that safely wipes and deletes the original unencrypted file after successful encryption.

---

### Command-Line Architecture & Subcommands

Organize your CLI options clean and intuitively using Go's `flag` package or standard subcommands:

| Subcommand | Flags | Description |
| --- | --- | --- |
| `encrypt` | `-in <file>`, `-out <file>` | Encrypts source file to output path. Prompts for password. |
| `decrypt` | `-in <file>`, `-out <file>` | Decrypts source file to output path. Prompts for password. |
| **Global / Optional** | `-overwrite` | Safely removes input file upon completion. |
|  | `-pass <string>` | Optional non-interactive flag for scripting (discouraged for manual use). |

---

### Operational & Execution Flow

#### Encryption Workflow

1. **Input Validation:** Ensure source file exists and destination path is writeable.
2. **Password Capture:** Prompt the user for a passphrase and ask them to confirm/re-enter it.
3. **Random Generation:** Generate a cryptographically secure 16-byte Salt and 12-byte Nonce using Go's `crypto/rand`.
4. **Key Derivation:** Run the password + Salt through PBKDF2 with high iterations (e.g., 100,000+ rounds) to produce the 256-bit key.
5. **Header Writing:** Write the Salt and Nonce directly to the beginning of the target file.
6. **Streaming Encryption:** Read the source file in chunks, encrypt using AES-GCM, and append ciphertext to the target file.
7. **Clean Exit:** Wipe the derived key from memory variables and print a success message with file size stats.

#### Decryption Workflow

1. **Input Reading:** Open the encrypted file and read the header (the fixed-length Salt and Nonce bytes).
2. **Password Capture:** Prompt the user once for the passphrase.
3. **Key Derivation:** Recompute the 256-bit key using the user's password and the Salt extracted from the header.
4. **Decryption & Verification:** Stream through the ciphertext chunk by chunk. AES-GCM will automatically check the integrity tag.
5. **Error Handling:** If the password is wrong or the file was modified, return a clear error (*"Decryption failed: Incorrect password or corrupted file"*) without producing a half-written output file.