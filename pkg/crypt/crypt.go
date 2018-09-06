// Package crypt is a structure to define the key, password and IV of an encryption/decryption function to be embedded with data that is to be kept encrypted except when being used
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
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

func nilError(s string) *Crypt {
	r := NewCrypt()
	r.err = errors.New(s + " nil receiver")
	return r
}

// NewCrypt empties a crypt or creates an empty crypt
func NewCrypt(r ...*Crypt) *Crypt {
	if len(r) == 0 {
		r = append(r, new(Crypt))
	}
	if r[0] == nil {
		r[0] = new(Crypt)
	}
	if r[0].crypt == nil {
		r[0].crypt = new(Bytes)
	} else {
		r[0].crypt.Null()
	}
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

// Arm generates the ciphertext from the password, uses it to decrypt the crypt into the crypt's main cyphertext, and creates the AES-GCM cipher
func (r *Crypt) Arm() *Crypt {
	switch {
	case r == nil:
		r = nilError("Arm()")
		r.SetError("nil receiver")
	case r.Password().Len() == 0:
		r.SetError("nil password")
	case r.Crypt().Len() == 0:
		r.SetError("nil crypt")
	case r.IV().Len() == 0:
		r.SetError("nil IV")
	default:
		var C *LockedBuffer
		var IV *Bytes
		C, IV, r.err = kdf.Gen(r.Password(), r.IV(), r.iterations)
		if r.err != nil {
			return r
		}
		var block cipher.Block
		block, r.err = aes.NewCipher(*C.Buf().(*[]byte))
		var blockmode cipher.AEAD
		blockmode, r.err = cipher.NewGCM(block)
		var c []byte
		c, r.err = blockmode.Open(nil, *IV.Buf().(*[]byte), *r.Crypt().Buf().(*[]byte), nil)
		r.ciphertext = NewLockedBuffer().Load(&c).(*LockedBuffer)
		block, r.err = aes.NewCipher(*r.ciphertext.Buf().(*[]byte))
		blockmode, r.err = cipher.NewGCM(block)
		r.gcm = &blockmode
		r.armed = true
	}
	return r
}

// Ciphertext returns the ciphertext stored in the crypt
func (r *Crypt) Ciphertext() *LockedBuffer {
	if r == nil {
		r := NewLockedBuffer()
		r.SetError("receiver was nil")
		return r
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
		b := NewBytes()
		b.SetError("nil receiver")
		return b
	}
	if r.crypt == nil {
		r.crypt = NewBytes()
		r.crypt.SetError("nil crypt")
	}
	return r.crypt
}

// Decrypt takes an encrypted Bytes and returns the decrypted data in a LockedBuffer
func (r *Crypt) Decrypt(b *Bytes) *LockedBuffer {
	switch {
	case r == nil:
		r.SetError("nil receiver")
	case !r.armed:
		r.SetError("not armed")
	case r.gcm == nil:
		r.SetError("nil gcm")
	default:
		var bb []byte
		bb, r.err = (*r.gcm).Open(nil, *r.IV().Buf().(*[]byte), *b.Buf().(*[]byte), nil)
		B := NewLockedBuffer().Load(&bb).(*LockedBuffer)
		return B
	}
	B := NewLockedBuffer()
	B.SetError(r.Error())
	return B
}

// Disarm clears the ciphertext
func (r *Crypt) Disarm() *Crypt {
	if r == nil {
		r = NewCrypt()
		r.SetError("nil receiver")
	}
	if r.gcm != nil {
		r.gcm = nil
	}
	r.ciphertext.Null().Free()
	r.ciphertext = nil
	r.armed = false
	return r
}

// Encrypt encrypts a Lockedbuffer and returns the ciphertext as Bytes
func (r *Crypt) Encrypt(lb *LockedBuffer) *Bytes {
	switch {
	case r == nil:
		r.SetError("nil receiver")
	case !r.armed:
		r.SetError("not armed")
	case r.gcm == nil:
		r.SetError("nil gcm")
	default:
		b := (*r.gcm).Seal(nil, *r.IV().Buf().(*[]byte), *lb.Buf().(*[]byte), nil)
		B := NewBytes().Load(&b).(*Bytes)
		return B
	}
	b := NewBytes()
	b.SetError(r.Error())
	return b
}

// Error returns the error stored in the crypt
func (r *Crypt) Error() string {
	if r == nil {
		return "receiver was nil"
	}
	if r.err != nil {
		return r.err.Error()
	}
	return ""
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
		return r
	}
	r.password = p
	r.ciphertext = NewLockedBuffer().Rand(32).(*LockedBuffer)
	r.SetRandomIV()
	r.iterations = kdf.Bench(time.Second)
	var C *LockedBuffer
	var IV *Bytes
	C, IV, r.err = kdf.Gen(r.Password(), r.IV(), r.iterations)
	// if r.err != nil {
	// 	return r
	// }
	var block cipher.Block
	block, r.err = aes.NewCipher(*C.Buf().(*[]byte))
	var blockmode cipher.AEAD
	blockmode, r.err = cipher.NewGCM(block)
	c := blockmode.Seal(nil, *IV.Buf().(*[]byte), *r.Ciphertext().Buf().(*[]byte), nil)
	r.crypt = r.crypt.Load(&c).(*Bytes)
	block, r.err = aes.NewCipher(*r.Ciphertext().Buf().(*[]byte))
	// if r.err != nil {
	// 	return r
	// }
	A := new(cipher.AEAD)
	a := *A
	a, r.err = cipher.NewGCM(block)
	// if r.err != nil {
	// 	return r
	// }
	r.gcm = &a
	r.armed = true
	return r
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
		r = NewCrypt()
		r.SetError("nil receciver")
		return r
	}
	if r.password == nil {
		r.password = NewPassword()
		r.SetError("nil password")
		return r
	}
	r.password.Null()
	r.gcm = nil
	r.unlocked = false
	r.Disarm()
	return r
}

// MarshalJSON renders the struct as JSON
func (r *Crypt) MarshalJSON() ([]byte, error) {
	var crypt, ciphertext, iv, password string
	if r.Crypt() != nil && r.Crypt().Len() != 0 {
		crypt = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.Crypt().Buf().(*[]byte)))...))
	}
	if r.Ciphertext() != nil && r.Ciphertext().Len() != 0 {
		ciphertext = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.Ciphertext().Buf().(*[]byte)))...))
	}
	if r.Error() == "ciphertext was nil" {
		r.err = nil
	}
	if r.IV() != nil && r.IV().Len() != 0 {
		iv = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.IV().Buf().(*[]byte)))...))
	}
	if r.Password() != nil && r.Password().Len() != 0 {
		password = string(*r.Password().Buf().(*[]byte))
	}
	return json.Marshal(&struct {
		Crypt      string
		Password   string
		Ciphertext string
		IV         string
		Iterations int64
		Unlocked   bool
		Armed      bool
		HasGCM     bool
		Error      string
	}{
		Crypt:      crypt,
		Password:   password,
		Ciphertext: ciphertext,
		IV:         iv,
		Iterations: int64(r.iterations),
		Unlocked:   r.unlocked,
		Armed:      r.armed,
		HasGCM:     r.gcm != nil,
		Error:      r.Error(),
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
		return r
	}
	if r.password == nil {
		r.password = NewPassword(r.password)
		r.SetError("nil password")
		return r
	}
	r.password = p
	r.unlocked = true
	return r
}
