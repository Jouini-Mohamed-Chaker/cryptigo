> **Note:** .md files are ai generated to guide me because I'm a beginner

**1. The 3 Core Concepts**

* **AES-256-GCM (The Safe):** AES is the algorithm that scrambles your data. **GCM** adds a "tamper-evident seal" (called an Authentication Tag). If someone tries to guess your password or change even a single byte of the encrypted file, GCM immediately fails instead of giving you broken output.
* **Salt (The Password Strengthener):** Human passwords (like `Secret123!`) are short and easy to guess. A Salt is a batch of completely random bytes generated fresh for every file. It gets mixed into your password so two files encrypted with the same password end up with completely different keys.
* **Nonce / IV (The One-Time Number):** Nonce stands for **N**umber used **ONCE**. It ensures that if you encrypt the exact same file twice, the output looks 100% different every time.

---

**2. The Simplest Way to Encrypt a File in Go**

If you are just getting started and want the most straightforward implementation, **read the whole file into memory at once**. This removes the complexity of chunking while keeping your code secure.

**Encryption Workflow:**

1. Open the source file and read all bytes into memory.
2. Generate **16 random bytes** for the `salt` and **12 random bytes** for the `nonce` using `crypto/rand`.
3. Pass your password and `salt` into `pbkdf2` (with 100,000+ iterations) to turn your text password into a 32-byte secret key.
4. Pass that 32-byte key to `aes.NewCipher` and `cipher.NewGCM`.
5. Call `gcm.Seal(...)` on your file bytes. This scrambles the file and attaches the integrity check tag to the end.
6. Save everything into your final output file in this layout:
`[ 16-byte Salt ] [ 12-byte Nonce ] [ Encrypted Data + Tag ]`

---

**3. The Decryption Workflow**

To reverse the process, you unpack the file from left to right:

1. Open the encrypted file.
2. Read the **first 16 bytes**—that's your `salt`.
3. Read the **next 12 bytes**—that's your `nonce`.
4. Read the **rest of the file**—that's your scrambled data.
5. Re-run `pbkdf2` using the user-provided password and the extracted `salt` to recreate the exact 32-byte key.
6. Initialize `cipher.NewGCM` with the key, then call `gcm.Open(...)` with the `nonce` and scrambled data.
* If the password is right and the file wasn't altered: you get your original bytes back!
* If the password is wrong: `gcm.Open` throws an error instantly.