package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"os"

	"golang.org/x/crypto/argon2"
)

func Encrypt(inFile, outFile, password string) error {
	// functions params are already validated

	salt := make([]byte, SALT_SIZE_IN_BYTES)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	nonce := make([]byte, NONCE_SIZE_IN_BYTES)
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	secretKey := argon2.IDKey([]byte(password), salt, cfg.Time, cfg.Memory,
		cfg.Threads, cfg.KeyLen)

	// zero out secretKey from memory when Encrypt returns
	// for extra secruity
	defer func() {
		for i:= range secretKey {
			secretKey[i] = 0
		}
	}()

	// Initialize raw AES-256 block engine
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	// Wrap the block engine in GCM mode so it can handle any file size
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// TODO: replace reading the whole file into memory to chunking
	file, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	// Append salt and nonce to the file
	var encryptedFileBytes []byte
	encryptedFileBytes = append(encryptedFileBytes, salt...)
	encryptedFileBytes = append(encryptedFileBytes, nonce...)

	// gcm.Seal() scrambles the file and attaches the integrity check tag to the end
	// First param is the destination slice
	encryptedFileBytes = gcm.Seal(encryptedFileBytes, nonce, file, nil)

	// encrypted file format: [ 16-byte Salt ] [ 12-byte Nonce ] [ Encrypted Data + Tag ]
	// Save to disk
	err = os.WriteFile(outFile, encryptedFileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
