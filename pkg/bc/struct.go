package bc

import (
	"crypto/cipher"

	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
)

var er = core.Errors

// BlockCrypt stores the state of a GCM AES cipher for encrypting up to 4Gb of data
type BlockCrypt struct {
	Crypt           *buf.Byte
	Password        *buf.Secure
	Ciphertext      *buf.Secure
	IV              *buf.Byte
	Iterations      int64
	GCM             *cipher.AEAD
	Idx             *[]byte
	Unlocked, Armed bool
	core.State
}
