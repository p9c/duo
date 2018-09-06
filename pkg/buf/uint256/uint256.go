// Package b32 is a store for 32 byte values used for signatures and public/private keys
package b32

import (
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
	"gitlab.com/parallelcoin/duo/pkg/buf/pass"
	"gitlab.com/parallelcoin/duo/pkg/buf/sec"
	"gitlab.com/parallelcoin/duo/pkg/cipher"
	"gitlab.com/parallelcoin/duo/pkg/def"
)

// B32 is a store  for 32 byte values either securely, using a cipher to keep data encrypted when not in use, or a regular, unsecure byteslice
type B32 struct {
	bytes        *buf.Unsafe
	lockedbuffer *buf.Fenced
	cipher       *cipher.Cipher
	secured      bool
	err          error
}

// New creates a new B32
func New() *B32 {
	r := new(B32)
	if r.bytes == nil {
		r.bytes = buf.New()
	}
	if r.lockedbuffer == nil {
		r.lockedbuffer = buf.New()
	}
	if r.cipher == nil {
		r.cipher = cipher.New()
	}
	r.secured = false
	return r
}

// Buffer implementation

// Buf is
func (r *B32) Buf() interface{} {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Copy is
func (r *B32) Copy(def.Buffer) def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Free is
func (r *B32) Free() def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Link is
func (r *B32) Link(interface{}) def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Load is
func (r *B32) Load(interface{}) def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Move is
func (r *B32) Move(def.Buffer) def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// New is
func (r *B32) New(int) def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Null is
func (r *B32) Null() def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Rand is
func (r *B32) Rand(...int) def.Buffer {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Size is
func (r *B32) Size() int {
	if nil == r {
		return -1
	}
	return 0
}

// String is
func (r *B32) String() string {
	if nil == r {
		return "nil receiver"
	}
	return ""
}

// Coding is
func (r *B32) Coding() string {
	if nil == r {
		return "nil receiver"
	}
	return ""
}

// SetCoding is
func (r *B32) SetCoding(string) interface{} {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Codes is
func (r *B32) Codes() []string {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return []string{}
}

// SetError is
func (r *B32) SetError(string) interface{} {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// UnsetError is
func (r *B32) UnsetError() interface{} {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Error is
func (r *B32) Error() string {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return ""
}

// Cap is
func (r *B32) Cap() int {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return 0
}

// Elem is
func (r *B32) Elem(int) interface{} {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Len is
func (r *B32) Len() int {
	if nil == r {
		return -1
	}
	return 0
}

// Purge is
func (r *B32) Purge() interface{} {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// SetElem is
func (r *B32) SetElem(int, interface{}) interface{} {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Cipher implementation

// Arm is
func (r *B32) Arm() cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Ciphertext is
func (r *B32) Ciphertext() *buf.Fenced {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return &buf.Fenced{}
}

// Disarm is
func (r *B32) Disarm() cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// IV is
func (r *B32) IV() *buf.Unsafe {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return &buf.Unsafe{}
}

// IsArmed is
func (r *B32) IsArmed() bool {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return false
}

// IsUnlocked is
func (r *B32) IsUnlocked() bool {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return false
}

// Lock is
func (r *B32) Lock() cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Password is
func (r *B32) Password() *passbuf.Password {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return &passbuf.Password{}
}

// SetIV is
func (r *B32) SetIV(b *buf.Unsafe) cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// SetRandomIV is
func (r *B32) SetRandomIV() cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Unlock is
func (r *B32) Unlock(p *passbuf.Password) cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// IsSecure is
func (r *B32) IsSecure() bool {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return false
}

// Secure is
func (r *B32) Secure(*buf.Fenced, *passbuf.Password, *buf.Unsafe) cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}

// Unsecure is
func (r *B32) Unsecure() cipher.Crypt {
	if nil == r {
		r = New().SetError("nil receiver").(*B32)
	}
	return r
}
