package cipher

import (
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
	"gitlab.com/parallelcoin/duo/pkg/buf/pass"
	"gitlab.com/parallelcoin/duo/pkg/def"
)

// Crypt is a generic interface for a buffer that keeps data stored encrypted and decrypts it for read functions and encrypts it for write functions
type Crypt interface {
	def.Buffer
	Arm() Crypt
	Ciphertext() *buf.Fenced
	Disarm() Crypt
	IV() *buf.Unsafe
	IsArmed() bool
	IsUnlocked() bool
	IsSecure() bool
	Lock() Crypt
	Password() *passbuf.Password
	Secure(*buf.Fenced, *passbuf.Password, *buf.Unsafe) Crypt
	SetIV(b *buf.Unsafe) Crypt
	SetRandomIV() Crypt
	Unlock(p *passbuf.Password) Crypt
	Unsecure() Crypt
}
