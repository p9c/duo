// Package crypt is a structure to define the key, password and IV of an encryption/decryption function to be embedded with data that is to be kept encrypted except when being used
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	"gitlab.com/parallelcoin/duo/pkg/kdf"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
	. "gitlab.com/parallelcoin/duo/pkg/password"
	"time"
)

// Crypt has a primary embed from a Bytes type that stores the encrypted data, so loading it is simple.
type Crypt struct {
	crypt           *Bytes
	password        *Password
	ciphertext      *LockedBuffer
	iv              *Bytes
	iterations      int
	unlocked, armed bool
	gcm             *cipher.AEAD
	err             error
}

// NewCrypt empties a crypt or creates an empty crypt
func NewCrypt(r ...*Crypt) *Crypt {
	if len(r) == 0 {
		r = append(r, new(Crypt))
	}
	if r[0] == nil {
		r[0] = new(Crypt)
	}
	r[0].crypt.Null()
	if r[0].password == nil {
		r[0].password = new(Password)
	} else {
		r[0].password.Null()
	}
	if r[0].ciphertext == nil {
		r[0].ciphertext = new(LockedBuffer)
	} else {
		r[0].ciphertext.Null()
	}
	if r[0].iv == nil {
		r[0].iv = new(Bytes)
	} else {
		r[0].iv.Null()
	}
	r[0].unlocked = false
	r[0].armed = false
	r[0].gcm = nil
	return r[0]
}

type crypt interface {
	Arm() *Crypt
	Ciphertext() *LockedBuffer
	Crypt() *Bytes
	Decrypt(*Bytes) *LockedBuffer
	Disarm() *Crypt
	Encrypt(*LockedBuffer) *Bytes
	Error() string
	Generate(*Password) *Crypt
	IsArmed() bool
	IsLoaded() bool
	IsUnlocked() bool
	IV() *Bytes
	Load(*Bytes) *Crypt
	Lock() *Crypt
	MarshalJSON() ([]byte, error)
	Null()
	Password() *Password
	SetError(string) *crypt
	SetIV(*Bytes) *Crypt
	SetRandomIV() *Crypt
	String() string
	Unlock(*Password) *Crypt
}

// Arm generates the ciphertext from the password, uses it to decrypt the crypt into the crypt's main cyphertext, and creates the AES-GCM cipher
func (r *Crypt) Arm() *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	r.armed = true
	return r
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

// Decrypt takes an encrypted Bytes and returns the decrypted data in a LockedBuffer
func (r *Crypt) Decrypt(*Bytes) *LockedBuffer {
	if r == nil {
		return &LockedBuffer{}
	}
	return nil
}

// Disarm clears the ciphertext
func (r *Crypt) Disarm() *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	r.armed = false
	return r
}

// Encrypt encrypts a Lockedbuffer and returns the ciphertext as Bytes
func (r *Crypt) Encrypt(*LockedBuffer) *Bytes {
	if r == nil {
		return &Bytes{}
	}
	return nil
}

// Error returns the error stored in the crypt
func (r *Crypt) Error() string {
	if r == nil {
		return "receiver was nil"
	}
	if r.err != nil {
		return r.err.Error()
	}
	return "<nil>"
}

// Generate creates a new crypt based on a password and a newly generated random ciphertext
func (r *Crypt) Generate(p *Password) *Crypt {
	if r == nil {
		r = NewCrypt()
		r.SetError("nil receiver")
	}
	if p == nil {
		r.password = new(Password)
		r.SetError("no password given")
	} else {
		r.password = p
	}
	r.ciphertext = NewLockedBuffer().Rand(36)
	r.SetRandomIV()
	r.iterations = kdf.Bench(time.Second)
	fmt.Println("before", r.String())
	C, IV, _ := kdf.Gen(r.Password(), r.IV(), r.iterations)
	var block cipher.Block
	block, r.err = aes.NewCipher(*C.Buf())
	var blockmode cipher.AEAD
	blockmode, r.err = cipher.NewGCM(block)
	c := blockmode.Seal(nil, *IV.Buf(), *r.Ciphertext().Buf(), nil)
	r.crypt.Load(&c)
	r.armed = true
	return r
}

// IsArmed returns true if the crypt is armed
func (r *Crypt) IsArmed() bool {
	if r == nil {
		return false
	}
	return r.armed
}

// IsLoaded returns true if the crypt contains data
func (r *Crypt) IsLoaded() bool {
	if r == nil || r.crypt == nil {
		return false
	}
	return true
}

// IsUnlocked returns whether the crypt is locked or not
func (r *Crypt) IsUnlocked() bool {
	if r == nil {
		return false
	}
	return r.unlocked
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

// Lock clears the password and disarms the crypt if it is armed
func (r *Crypt) Lock() *Crypt {
	if r == nil {
		// r = new(Crypt).NilGuard(r, null).(*Crypt)
	}
	r.unlocked = false
	return r
}

// MarshalJSON renders the struct as JSON
func (r *Crypt) MarshalJSON() ([]byte, error) {
	var crypt, ciphertext, iv string
	if r.Crypt() != nil && r.Crypt().Len() != 0 {
		crypt = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.Crypt().Buf()))...))
	}
	if r.Ciphertext() != nil && r.Ciphertext().Len() != 0 {
		ciphertext = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.Ciphertext().Buf()))...))
	}
	if r.IV() != nil && r.IV().Len() != 0 {
		iv = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.IV().Buf()))...))
	}
	return json.Marshal(&struct {
		Crypt      string
		Password   string
		Ciphertext string
		IV         string
		Iterations int64
		Unlocked   bool
		Armed      bool
	}{
		Crypt:      crypt,
		Password:   string(*r.Password().Buf()),
		Ciphertext: ciphertext,
		IV:         iv,
		Iterations: int64(r.iterations),
		Unlocked:   r.unlocked,
		Armed:      r.armed,
	})
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Crypt) Null() *Crypt {
	return NewCrypt(r)
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

// SetError sets the error in the Crypt
func (r *Crypt) SetError(s string) *Crypt {
	if r == nil {
		r = new(Crypt)
	}
	r.err = errors.New(s)
	return r
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

// String prints the JSON representation of the data and structure
func (r *Crypt) String() string {
	s, _ := json.MarshalIndent(r, "", "    ")
	return string(s)
}

// Unlock sets the password, runs the KDF and arms the
func (r *Crypt) Unlock(p *Password) *Crypt {
	if r == nil {
		r = new(Crypt)
		r.SetError("nil receiver")
	}
	if r.password != nil {
		r.password = NewPassword(r.password)
	}
	r.password = p
	r.unlocked = true
	return r
}
