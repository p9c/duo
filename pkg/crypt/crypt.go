// Package crypt is a structure to define the key, password and IV of an encryption/decryption function to be embedded with data that is to be kept encrypted except when being used
package crypt

import (
	"crypto/cipher"
	"errors"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
	. "gitlab.com/parallelcoin/duo/pkg/password"
	. "gitlab.com/parallelcoin/duo/pkg/pipe"
)

// Crypt has a primary embed from a Bytes type that stores the encrypted data, so loading it is simple.
type Crypt struct {
	*Pipe
	crypt           *Bytes
	password        *Password
	ciphertext      *LockedBuffer
	iv              *Bytes
	iterations      int
	unlocked, armed bool
	gcm             *cipher.AEAD
	err             error
}
type crypt interface {
	Error() error
	SetError(string) *crypt
	Crypt() *Bytes
	Load(*Bytes) *Crypt
	IsLoaded() bool
	Null()
	Generate(*Password) *Crypt
	Password() *Password
	Ciphertext() *LockedBuffer
	IV() *Bytes
	SetIV(*Bytes) *Crypt
	SetRandomIV() *Crypt
	Unlock(*Password) *Crypt
	Lock() *Crypt
	IsUnlocked() bool
	Arm() *Crypt
	Disarm() *Crypt
	IsArmed() bool
	Encrypt(*LockedBuffer) *Bytes
	Decrypt(*Bytes) *LockedBuffer
}

// NewCrypt returns a new empty Crypt
func NewCrypt() *Crypt {
	return new(Crypt)
}

// Error returns the error stored in the crypt
func (r *Crypt) Error() error {
	if r == nil {
		return errors.New("receiver was nil")
	}
	return r.err
}

// SetError sets the error in the Crypt
func (r *Crypt) SetError(s string) *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	r.err = errors.New(s)
	return r
}

// Crypt returns the Bytes buffer crypt
func (r *Crypt) Crypt() *Bytes {
	if r == nil {
		return new(Bytes)
	}
	if r.crypt == nil {
		r.crypt = new(Bytes)
	}
	return r.crypt
}

// Load moves a bytes into the crypt
func (r *Crypt) Load(bytes *Bytes) *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	if r.crypt == nil {
		r.crypt = new(Bytes)
	}
	r.crypt.Move(bytes)
	return r
}

// IsLoaded returns true if the crypt contains data
func (r *Crypt) IsLoaded() bool {
	if r == nil || r.crypt == nil {
		return false
	}
	return true
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Crypt) Null() *Crypt {
	return NullCrypt(r).(*Crypt)
}

// NullCrypt empties a crypt or creates an empty crypt
func NullCrypt(R interface{}) interface{} {
	r := R.(*Crypt)
	if r == nil {
		r = new(Crypt)
	}
	if r.crypt != nil {
		r.crypt.Null()
	}
	r.crypt.Null()
	if r.password == nil {
		r.password = new(Password)
	} else {
		r.password.Null()
	}
	if r.ciphertext == nil {
		r.ciphertext = new(LockedBuffer)
	} else {
		r.ciphertext.Null()
	}
	if r.iv == nil {
		r.iv = new(Bytes)
	} else {
		r.iv.Null()
	}
	r.unlocked = false
	r.armed = false
	r.gcm = nil
	return r
}

// Generate creates a new crypt based on a password and a newly generated random ciphertext
func (r *Crypt) Generate(*Password) *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	if r == nil {
		r = new(Crypt)
	}
	r.ciphertext = NewLockedBuffer().Rand(36)
	r.SetRandomIV()
	return r
}

// Password returns the password stored in the Crypt
func (r *Crypt) Password() *Password {
	if r == nil {
		r = new(Crypt)
		r.SetError("receiver was nil")
	}
	if r.password == nil {
		r.password = NewPassword()
		r.password.SetError("password was nil")
	}
	return r.password
}

// Ciphertext returns the ciphertext stored in the crypt
func (r *Crypt) Ciphertext() *LockedBuffer {
	if r == nil {
		R := new(LockedBuffer)
		R.SetError("receiver was nil")
		return R
	}
	if r.ciphertext == nil {
		r.ciphertext = new(LockedBuffer)
		r.SetError("ciphertext was nil")
	}
	return r.ciphertext
}

// IV returns the initialisation vector stored in the crypt
func (r *Crypt) IV() *Bytes {
	if r == nil {
		return new(Bytes)
	}
	if r.iv == nil {
		r.iv = new(Bytes)
	}
	return r.iv
}

// SetIV loads the IV with a Bytes. It must be 12 bytes long.
func (r *Crypt) SetIV(b *Bytes) *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	if b == nil {
		b.SetError("nil Bytes")
	} else if b.Len() != 12 {
		b.SetError("must be 12 bytes")
	}
	if r.iv == nil {
		r.iv = NewBytes()
	}
	r.iv.Move(b)
	return r
}

// SetRandomIV loads the IV with a random 12 bytes.
func (r *Crypt) SetRandomIV() *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	if r.iv == nil {
		r.iv = new(Bytes)
	}
	r.iv.Rand(12)
	return r
}

// Unlock sets the password, runs the KDF and arms the
func (r *Crypt) Unlock(*Password) *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	r.unlocked = true
	return r
}

// Lock clears the password and disarms the crypt if it is armed
func (r *Crypt) Lock() *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	r.unlocked = false
	return r
}

// IsUnlocked returns whether the crypt is locked or not
func (r *Crypt) IsUnlocked() bool {
	if r == nil {
		return false
	}
	return r.unlocked
}

// Arm generates the ciphertext from the password, uses it to decrypt the crypt into the crypt's main cyphertext, and creates the AES-GCM cipher
func (r *Crypt) Arm() *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	r.armed = true
	return r
}

// Disarm clears the ciphertext
func (r *Crypt) Disarm() *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	r.armed = false
	return r
}

// IsArmed returns true if the crypt is armed
func (r *Crypt) IsArmed() bool {
	if r == nil {
		return false
	}
	return r.armed
}

// Encrypt encrypts a Lockedbuffer and returns the ciphertext as Bytes
func (r *Crypt) Encrypt(*LockedBuffer) *Bytes {
	if r == nil {
		return &Bytes{}
	}
	return nil
}

// Decrypt takes an encrypted Bytes and returns the decrypted data in a LockedBuffer
func (r *Crypt) Decrypt(*Bytes) *LockedBuffer {
	if r == nil {
		return &LockedBuffer{}
	}
	return nil
}
