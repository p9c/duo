// Package uint256
package uint256

import (
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/crypt"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
	. "gitlab.com/parallelcoin/duo/pkg/password"
)

//
type Uint256 struct {
	*Bytes
	*Crypt
	secured bool
}

//
func NewUint256(r ...*Uint256) *Uint256 {
	if len(r) == 0 {
		r = append(r, new(Uint256))
	}
	if r[0] == nil {
		r[0] = new(Uint256)
	}
	if r[0].Bytes == nil {
		r[0].Bytes = NewBytes()
	}
	if r[0].Crypt == nil {
		r[0].Crypt = NewCrypt()
	}
	r[0].secured = false
	return r[0]
}

type uint256 interface {
	Arm() *Uint256
	Buf() []byte
	Ciphertext() *LockedBuffer
	Copy(*Uint256) *Uint256
	Crypt() *Bytes
	Decrypt(*Bytes) *LockedBuffer
	Delete()
	Disarm() *Uint256
	Encrypt(*LockedBuffer) *Bytes
	Error() string
	Generate(*Password) *Uint256
	IV() *Bytes
	IsArmed() bool
	IsLoaded() bool
	IsSet() bool
	IsUTF8() bool
	IsUnlocked() bool
	Len() int
	Link(*Uint256) *Uint256
	Load(*Bytes) *Uint256
	Lock() *Uint256
	MarshalJSON() ([]byte, error)
	Move(*Uint256) *Uint256
	New(int) *Uint256
	Null() *Uint256
	Password() *Password
	Rand(int) *Uint256
	SetBin() *Uint256
	SetError(string) *Uint256
	SetIV(*Bytes) *Uint256
	SetRandomIV() *Uint256
	SetUTF8() *Uint256
	String() string
	Unlock(*Password) *Uint256
}

// Arm does...
func (r *Uint256) Arm() *Uint256 {
	return r
}

// Buf returns the contents of the buffer whether it's the secure or unsecured buffer
func (r *Uint256) Buf() []byte {
	return nil
}

// Ciphertext does...
func (r *Uint256) Ciphertext() *LockedBuffer {
	return r
}

// Copy does...
func (r *Uint256) Copy(*Uint256) *Uint256 {
	return r
}

// Crypt does...
func (r *Uint256) Crypt() *Bytes {
	return r
}

// Decrypt does...
func (r *Uint256) Decrypt(*Bytes) *LockedBuffer {
	return r
}

// Delete does...
func (r *Uint256) Delete() {
	return r
}

// Disarm does...
func (r *Uint256) Disarm() *Uint256 {
	return r
}

// Encrypt does...
func (r *Uint256) Encrypt(*LockedBuffer) *Bytes {
	return r
}

// Error does...
func (r *Uint256) Error() string {
	return r
}

// Generate does...
func (r *Uint256) Generate(*Password) *Uint256 {
	return r
}

// IV does...
func (r *Uint256) IV() *Bytes {
	return r
}

// IsArmed does...
func (r *Uint256) IsArmed() bool {
	return r
}

// IsLoaded does...
func (r *Uint256) IsLoaded() bool {
	return r
}

// IsSet does...
func (r *Uint256) IsSet() bool {
	return r
}

// IsUTF8 does...
func (r *Uint256) IsUTF8() bool {
	return r
}

// IsUnlocked does...
func (r *Uint256) IsUnlocked() bool {
	return r
}

// Len does...
func (r *Uint256) Len() int {
	return r
}

// Link does...
func (r *Uint256) Link(*Uint256) *Uint256 {
	return r
}

// Load does...
func (r *Uint256) Load(*Bytes) *Uint256 {
	return r
}

// Lock does...
func (r *Uint256) Lock() *Uint256 {
	return r
}

// MarshalJSON does...
func (r *Uint256) MarshalJSON() ([]byte, error) {
	return r
}

// Move does...
func (r *Uint256) Move(*Uint256) *Uint256 {
	return r
}

// New does...
func (r *Uint256) New(int) *Uint256 {
	return r
}

// Null does...
func (r *Uint256) Null() *Uint256 {
	return r
}

// Password does...
func (r *Uint256) Password() *Password {
	return r
}

// Rand does...
func (r *Uint256) Rand(int) *Uint256 {
	return r
}

// SetBin does...
func (r *Uint256) SetBin() *Uint256 {
	return r
}

// SetError does...
func (r *Uint256) SetError(string) *Uint256 {
	return r
}

// SetIV does...
func (r *Uint256) SetIV(*Bytes) *Uint256 {
	return r
}

// SetRandomIV does...
func (r *Uint256) SetRandomIV() *Uint256 {
	return r
}

// SetUTF8 does...
func (r *Uint256) SetUTF8() *Uint256 {
	return r
}

// String does...
func (r *Uint256) String() string {
	return r
}

// Unlock does...
func (r *Uint256) Unlock(*Password) *Uint256 {
	return r
}
