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
	bytes        *bytes.Bytes
	lockedbuffer *secbuf.SecBuf
	cipher       *cipher.Cipher
	secured      bool
	err          error
}

// New creates a new B32
func New(r ...*B32) *B32 {
	if len(r) == 0 {
		r = append(r, new(B32))
	}
	if r[0] == nil {
		r[0] = new(B32)
	}
	if r[0].bytes == nil {
		r[0].bytes = bytes.New()
	}
	if r[0].lockedbuffer == nil {
		r[0].lockedbuffer = secbuf.New()
	}
	if r[0].cipher == nil {
		r[0].cipher = cipher.New()
	}
	r[0].secured = false
	return r[0]
}

// Buffer implementation

// Buf is
func (r *B32) Buf() interface{} {
	panic("not implemented")
}

// Copy is
func (r *B32) Copy(def.Buffer) def.Buffer {
	panic("not implemented")
}

// Free is
func (r *B32) Free() def.Buffer {
	panic("not implemented")
}

// Link is
func (r *B32) Link(interface{}) def.Buffer {
	panic("not implemented")
}

// Load is
func (r *B32) Load(interface{}) def.Buffer {
	panic("not implemented")
}

// Move is
func (r *B32) Move(def.Buffer) def.Buffer {
	panic("not implemented")
}

// New is
func (r *B32) New(int) def.Buffer {
	panic("not implemented")
}

// Null is
func (r *B32) Null() def.Buffer {
	panic("not implemented")
}

// Rand is
func (r *B32) Rand(...int) def.Buffer {
	panic("not implemented")
}

// Size is
func (r *B32) Size() int {
	panic("not implemented")
}

// String is
func (r *B32) String() string {
	panic("not implemented")
}

// Coding is
func (r *B32) Coding() string {
	panic("not implemented")
}

// SetCoding is
func (r *B32) SetCoding(string) interface{} {
	panic("not implemented")
}

// Codes is
func (r *B32) Codes() []string {
	panic("not implemented")
}

// SetError is
func (r *B32) SetError(string) interface{} {
	panic("not implemented")
}

// UnsetError is
func (r *B32) UnsetError() interface{} {
	panic("not implemented")
}

// Error is
func (r *B32) Error() string {
	panic("not implemented")
}

// Cap is
func (r *B32) Cap() int {
	panic("not implemented")
}

// Elem is
func (r *B32) Elem(int) interface{} {
	panic("not implemented")
}

// Len is
func (r *B32) Len() int {
	panic("not implemented")
}

// Purge is
func (r *B32) Purge() interface{} {
	panic("not implemented")
}

// SetElem is
func (r *B32) SetElem(int, interface{}) interface{} {
	panic("not implemented")
}

// Cipher implementation

// Arm is
func (r *B32) Arm() *cipher.Crypt {
	panic("not implemented")
}

// Ciphertext is
func (r *B32) Ciphertext() *secbuf.SecBuf {
	panic("not implemented")
}

// Disarm is
func (r *B32) Disarm() *cipher.Crypt {
	panic("not implemented")
}

// IV is
func (r *B32) IV() *bytes.Bytes {
	panic("not implemented")
}

// IsArmed is
func (r *B32) IsArmed() bool {
	panic("not implemented")
}

// IsUnlocked is
func (r *B32) IsUnlocked() bool {
	panic("not implemented")
}

// Lock is
func (r *B32) Lock() *cipher.Crypt {
	panic("not implemented")
}

// Password is
func (r *B32) Password() *passbuf.Password {
	panic("not implemented")
}

// SetIV is
func (r *B32) SetIV(b *bytes.Bytes) *cipher.Crypt {
	panic("not implemented")
}

// SetRandomIV is
func (r *B32) SetRandomIV() *cipher.Crypt {
	panic("not implemented")
}

// Unlock is
func (r *B32) Unlock(p *passbuf.Password) *cipher.Crypt {
	panic("not implemented")
}
