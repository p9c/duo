package cipher

import (
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
	"gitlab.com/parallelcoin/duo/pkg/buf/pass"
	"gitlab.com/parallelcoin/duo/pkg/buf/sec"
	"gitlab.com/parallelcoin/duo/pkg/def"
)

// Crypt is a generic interface for a buffer that keeps data stored encrypted and decrypts it for read functions and encrypts it for write functions
type Crypt interface {
	def.Buffer
	Arm() *Crypt
	Ciphertext() *secbuf.SecBuf
	Disarm() *Crypt
	IV() *bytes.Bytes
	IsArmed() bool
	IsUnlocked() bool
	Lock() *Crypt
	Password() *passbuf.Password
	SetIV(b *bytes.Bytes) *Crypt
	SetRandomIV() *Crypt
	Unlock(p *passbuf.Password) *Crypt
}
