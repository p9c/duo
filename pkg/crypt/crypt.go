// Package crypt is a structure to define the key, password and IV of an encryption/decryption function to be embedded with data that is to be kept encrypted except when being used
package crypt

import (
	"crypto/cipher"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
	. "gitlab.com/parallelcoin/duo/pkg/password"
)

// Crypt has a primary embed from a Bytes type that stores the encrypted data, so loading it is simple.
type Crypt struct {
	Bytes
	password        *Password
	ciphertext      *LockedBuffer
	iv              *Bytes
	iterations      int
	unlocked, armed bool
	gcm             *cipher.AEAD
}
type crypt interface {
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

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Crypt) Null() *Crypt {
	return null(r).(*Crypt)
}
func null(R interface{}) interface{} {
	r := R.(*Crypt)
	if r == nil {
		r = new(Crypt)
	}
	if r.Buf() != nil {
		rr := *r.Buf()
		if r.IsSet() {
			for i := range rr {
				rr[i] = 0
			}
		}
	}
	r.val = nil
	r.set = false
	r.err = nil
	return r
}

// Generate creates a new crypt based on a password and a newly generated random ciphertext
func (r *Crypt) Generate(*Password) *Crypt {
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
	}
	return nil
}

// Ciphertext returns the ciphertext stored in the crypt
func (r *Crypt) Ciphertext() *LockedBuffer {
	if r == nil {
		r = new(Crypt)
	}
	return nil
}

// IV returns the initialisation vector stored in the crypt
func (r *Crypt) IV() *Bytes {
	if r == nil {
		r = new(Crypt)
	}
	return nil
}

// SetIV loads the IV with a Bytes. It must be 12 bytes long.
func (r *Crypt) SetIV(b *Bytes) *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	if b.Len() == 0 {
		r.SetError("nil Bytes")
		return r
	}
	if b.Len() != 12 {
		r.SetError("must be 12 bytes")
	}
	r.iv.Move(b)
	return nil
}

// SetRandomIV loads the IV with a random 12 bytes.
func (r *Crypt) SetRandomIV() *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	r.iv.Move(NewBytes().Rand(12))
	return r
}

// Unlock sets the password, runs the KDF and arms the
func (r *Crypt) Unlock(*Password) *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	r.unlocked = true
	return r
}

// Lock clears the password and disarms the crypt if it is armed
func (r *Crypt) Lock() *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	r.unlocked = false
	return r
}

// IsUnlocked returns whether the crypt is locked or not
func (r *Crypt) IsUnlocked() bool {
	if r == nil {
		r = new(Crypt)
	}
	return r.unlocked
}

// Arm generates the ciphertext from the password, uses it to decrypt the crypt into the crypt's main cyphertext, and creates the AES-GCM cipher
func (r *Crypt) Arm() *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	r.armed = true
	return r
}

// Disarm clears the ciphertext
func (r *Crypt) Disarm() *Crypt {
	r.armed = false
	return r
}

// IsArmed returns true if the crypt is armed
func (r *Crypt) IsArmed() bool {
	return r.armed
}

// Encrypt encrypts a Lockedbuffer and returns the ciphertext as Bytes
func (r *Crypt) Encrypt(*LockedBuffer) *Bytes {
	return nil
}

// Decrypt takes an encrypted Bytes and returns the decrypted data in a LockedBuffer
func (r *Crypt) Decrypt(*Bytes) *LockedBuffer {
	return nil
}
