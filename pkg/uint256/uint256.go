// Package uint256 is a store for 32 byte values used for signatures and public/private keys
package uint256

import (
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/crypt"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
)

// Uint256 is a store  for 32 byte values either securely, using a crypt to keep data encrypted when not in use, or a regular, unsecure byteslice
type Uint256 struct {
	bytes        *Bytes
	lockedbuffer *LockedBuffer
	crypt        *Crypt
	secured      bool
	err          error
}

// NewUint256 creates a new Uint256
func NewUint256(r ...*Uint256) *Uint256 {
	if len(r) == 0 {
		r = append(r, new(Uint256))
	}
	if r[0] == nil {
		r[0] = new(Uint256)
	}
	if r[0].bytes == nil {
		r[0].bytes = NewBytes()
	}
	if r[0].lockedbuffer == nil {
		r[0].lockedbuffer = NewLockedBuffer()
	}
	if r[0].crypt == nil {
		r[0].crypt = NewCrypt()
	}
	r[0].secured = false
	return r[0]
}

type uint256 interface {
	Buf() []byte
	Copy(*Uint256) *Uint256
	Delete()
	Error() string
	IsSet() bool
	Len() int
	Link(*Uint256) *Uint256
	Load(*[]byte) *Uint256
	MarshalJSON() ([]byte, error)
	Move(*Uint256) *Uint256
	New(int) *Uint256
	Null() *Uint256
	Rand(int) *Uint256
	SetError(string) *Uint256
	String() string
	UseCrypt(*Crypt) *Uint256
	UseBytes() *Uint256
	IsUsingCrypt() bool
}

// Buf returns the contents of the buffer whether it's the secure or unsecured buffer
func (r *Uint256) Buf() []byte {
	return nil
}

// Copy copies in another Uint256
func (r *Uint256) Copy(*Uint256) *Uint256 {
	return r
}

// Delete wipes an insecure buffer or calls Destroy() on a LockedBuffer
func (r *Uint256) Delete() {
}

// Error returns the error stored in the Uint256
func (r *Uint256) Error() string {
	return ""
}

// IsSet returns true when the buffer is loaded
func (r *Uint256) IsSet() bool {
	return false
}

// Link copies the reference from one data to another. The link breaks if the security mode is changed.
func (r *Uint256) Link(*Uint256) *Uint256 {
	return r
}

// Load puts bytes into the buffer
func (r *Uint256) Load(*[]byte) *Uint256 {
	return r
}

// MarshalJSON renders the object into JSON format
func (r *Uint256) MarshalJSON() ([]byte, error) {
	return nil, nil
}

// Move references the buffer referred, and dereferences it from the source.
func (r *Uint256) Move(*Uint256) *Uint256 {
	return r
}

// New destroys the previous buffer and allocates a new one of a given size
func (r *Uint256) New(int) *Uint256 {
	return r
}

// Null destroys the buffer and sets the Uint256 into its Null state (all values default)
func (r *Uint256) Null() *Uint256 {
	return r
}

// Rand creates a new buffer from cryptographically secure random source, destroying the previous data
func (r *Uint256) Rand(int) *Uint256 {
	return r
}

// SetError sets the error stored in the Uint256
func (r *Uint256) SetError(string) *Uint256 {
	return r
}

// String converts the data to a string, in JSON format
func (r *Uint256) String() string {
	return ""
}

// UseCrypt enables the use of the crypt, referenced from an existing crypt.
func (r *Uint256) UseCrypt(c *Crypt) *Uint256 {
	return r
}

// UseBytes disables an armed crypt and decrypts the data
func (r *Uint256) UseBytes() *Uint256 {
	return r
}

// IsUsingCrypt returns true if the Uint256 is set to use a crypt
func (r *Uint256) IsUsingCrypt() bool {
	return false
}
