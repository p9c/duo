package blockcrypt

import (
	"crypto/cipher"
	"github.com/parallelcointeam/duo/pkg/buf"
)

// BlockCrypt stores the state of a GCM AES cipher for encrypting up to 4Gb of data
type BlockCrypt struct {
	Crypt           *buf.Bytes
	Password        *buf.Secure
	Ciphertext      *buf.Secure
	IV              *buf.Bytes
	Iterations      int
	GCM             *cipher.AEAD
	GCMIV           *buf.Bytes
	Unlocked, Armed bool
	Status          string
}
