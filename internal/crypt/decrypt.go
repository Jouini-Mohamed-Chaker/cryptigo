package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
)

func Decrypt(inFile, outFile, password string) error {
	inFileBytes, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	// To avoid trying to access slice members that don't exist
	if len(inFileBytes) < SALT_SIZE_IN_BYTES+NONCE_SIZE_IN_BYTES+GCM_AUTH_TAG_SIZE_IN_BYTES {
		return fmt.Errorf("malformed file")
	}
	// First 16 bytes are the salt
	salt := inFileBytes[:SALT_SIZE_IN_BYTES]
	// 12 Bytes after that are the nonce
	nonce := inFileBytes[SALT_SIZE_IN_BYTES : SALT_SIZE_IN_BYTES+NONCE_SIZE_IN_BYTES]
	// The rest is cyphertext + auth tag
	cypherText := inFileBytes[SALT_SIZE_IN_BYTES+NONCE_SIZE_IN_BYTES:]

	// Derive argon's secret key using the data we have
	secretKey := argon2.IDKey([]byte(password), salt, cfg.Time,
		cfg.Memory, cfg.Threads, cfg.KeyLen)

	// zero out secretKey from memory when Decrypt returns
	// for extra secruity
	defer func() {
		for i:= range secretKey {
			secretKey[i] = 0
		}
	}()

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Now we try decrypting the file bytes loaded in memory
	decryptedFileBytes, err := gcm.Open(nil, nonce, cypherText, nil)
	if err != nil {
		return err
	}

	err = os.WriteFile(outFile, decryptedFileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
