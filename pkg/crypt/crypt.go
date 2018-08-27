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
	*Bytes
	password      *Password
	ciphertext    *LockedBuffer
	iv            *Bytes
	iterations    int
	locked, armed bool
	gcm           *cipher.AEAD
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
	IsLocked() bool
	Arm() *Crypt
	Disarm() *Crypt
	IsArmed() bool
	Encrypt(*LockedBuffer) *Bytes
	Decrypt(*Bytes) *LockedBuffer
}

// Generate creates a new crypt based on a password and a newly generated random ciphertext
func (r *Crypt) Generate(*Password) *Crypt {
	return r
}

// Password returns the password stored in the Crypt
func (r *Crypt) Password() *Password {
	return nil
}

// Ciphertext returns the ciphertext stored in the crypt
func (r *Crypt) Ciphertext() *LockedBuffer {
	return nil
}

// IV returns the initialisation vector stored in the crypt
func (r *Crypt) IV() *Bytes {
	return nil
}

// SetIV loads the IV with a Bytes. It must be 12 bytes long.
func (r *Crypt) SetIV(*Bytes) *Crypt {
	return nil
}

// SetRandomIV loads the IV with a random 12 bytes.
func (r *Crypt) SetRandomIV() *Crypt {
	return r
}

// Unlock sets the password, runs the KDF and arms the
func (r *Crypt) Unlock(*Password) *Crypt {
	return r
}

// Lock clears the password and disarms the crypt if it is armed
func (r *Crypt) Lock() *Crypt {
	return r
}

// IsLocked returns whether the crypt is locked or not
func (r *Crypt) IsLocked() bool {
	return r.locked
}

// Arm generates the ciphertext from the password, uses it to decrypt the crypt into the crypt's main cyphertext, and creates the AES-GCM cipher
func (r *Crypt) Arm() *Crypt {
	return r
}

// Disarm clears the ciphertext
func (r *Crypt) Disarm() *Crypt {
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
