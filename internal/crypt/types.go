package crypt

const (
	SALT_SIZE_IN_BYTES         = 16
	NONCE_SIZE_IN_BYTES        = 12
	GCM_AUTH_TAG_SIZE_IN_BYTES = 16
)

// Params defines the config for Argon2id (ie the params passed
// to IDKey function, exluding the password and salt)
type Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// cfg is declared as a global variable so encrypt and decrypt
//
//	functions have the same config
var cfg = Params{
	Time:    3,
	Memory:  64 * 1024, // 64 MB
	Threads: 4,
	KeyLen:  32,
}
