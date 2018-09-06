// Package cipher is a structure to define the key, password and IV of an encryption/decryption function to be embedded with data that is to be kept encrypted except when being used
package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
	"gitlab.com/parallelcoin/duo/pkg/buf/pass"
	"gitlab.com/parallelcoin/duo/pkg/buf/sec"
	"time"
)

// Cipher has a primary embed from a buf.Unsafe type that stores the encrypted data, so loading it is simple.
type Cipher struct {
	crypt           *buf.Unsafe
	password        *passbuf.Password
	ciphertext      *buf.Fenced
	iv              *buf.Unsafe
	iterations      int
	unlocked, armed bool
	gcm             *cipher.AEAD
	err             error
}

func nilError(s string) *Cipher {
	r := New()
	r.err = errors.New(s + " nil receiver")
	return r
}

// New empties a crypt or creates an empty crypt
func New(r ...*Cipher) *Cipher {
	if len(r) == 0 {
		r = append(r, new(Cipher))
	}
	if r[0] == nil {
		r[0] = new(Cipher)
	}
	if r[0].crypt == nil {
		r[0].crypt = new(buf.Unsafe)
	} else {
		r[0].crypt.Null()
	}
	if r[0].password == nil {
		r[0].password = new(passbuf.Password)
	} else {
		r[0].password.Null()
	}
	if r[0].ciphertext == nil {
		r[0].ciphertext = new(buf.Fenced)
	} else {
		r[0].ciphertext.Null()
	}
	if r[0].iv == nil {
		r[0].iv = new(buf.Unsafe)
	} else {
		r[0].iv.Null()
	}
	r[0].unlocked = false
	r[0].armed = false
	r[0].gcm = nil
	return r[0]
}

// Arm generates the ciphertext from the password, uses it to decrypt the crypt into the crypt's main cyphertext, and creates the AES-GCM cipher
func (r *Cipher) Arm() *Cipher {
	switch {
	case r == nil:
		r = nilError("Arm() nil receiver")
		return r
	case r.Password().Len() < 1:
		r.SetError("Arm() nil password")
		return r
	case r.Cipher().Len() < 1:
		r.SetError("Arm() nil crypt")
		return r
	case r.IV().Len() < 1:
		r.SetError("Arm() nil IV")
		return r
	default:
		var C *buf.Fenced
		var IV *buf.Unsafe
		C, IV, r.err = Gen(r.Password(), r.IV(), r.iterations)
		if r.err != nil {
			return r
		}
		var block cipher.Block
		block, r.err = aes.NewCipher(*C.Buf().(*[]byte))
		var blockmode cipher.AEAD
		blockmode, r.err = cipher.NewGCM(block)
		var c []byte
		c, r.err = blockmode.Open(nil, *IV.Buf().(*[]byte), *r.Cipher().Buf().(*[]byte), nil)
		if r.err == nil {
			return r
		}
		r.ciphertext = buf.New().Load(&c).(*buf.Fenced)
		block, r.err = aes.NewCipher(*r.ciphertext.Buf().(*[]byte))
		if r.err == nil {
			fmt.Println(r.err)
			return r
		}
		if block == nil {
			fmt.Println(r.err)
			return r
		}
		blockmode, r.err = cipher.NewGCM(block)
		if r.err == nil {
			fmt.Println(r.err)
			return r
		}
		r.gcm = &blockmode
		r.armed = true
	}
	return r
}

// Ciphertext returns the ciphertext stored in the crypt
func (r *Cipher) Ciphertext() *buf.Fenced {
	if r == nil {
		r := buf.New()
		r.SetError("Ciphertext() receiver was nil")
		return r
	}
	if r.ciphertext == nil {
		r.ciphertext = new(buf.Fenced)
		r.SetError("ciphertext was nil")
	}
	return r.ciphertext
}

// Cipher returns the buf.Unsafe buffer crypt
func (r *Cipher) Cipher() *buf.Unsafe {
	if r == nil {
		b := buf.New()
		b.SetError("Cipher() nil receiver")
		return b
	}
	if r.crypt == nil {
		r.crypt = buf.New()
		r.crypt.SetError("Cipher nil crypt")
	}
	return r.crypt
}

// Decrypt takes an encrypted buf.Unsafe and returns the decrypted data in a buf.Fenced
func (r *Cipher) Decrypt(b *buf.Unsafe) *buf.Fenced {
	switch {
	case r == nil:
		r.SetError("Decrypt() nil receiver")
	case !r.armed:
		r.SetError("Decrypt() not armed")
	case r.gcm == nil:
		r.SetError("Decrypt() nil gcm")
	default:
		var bb []byte
		bb, r.err = (*r.gcm).Open(nil, *r.IV().Buf().(*[]byte), *b.Buf().(*[]byte), nil)
		B := buf.New().Load(&bb).(*buf.Fenced)
		return B
	}
	B := buf.New()
	B.SetError(r.Error())
	return B
}

// Disarm clears the ciphertext
func (r *Cipher) Disarm() *Cipher {
	if r == nil {
		r = New()
		r.SetError("Disarm() nil receiver")
	}
	if r.gcm != nil {
		r.gcm = nil
	}
	r.ciphertext.Null().Free()
	r.ciphertext = nil
	r.armed = false
	return r
}

// Encrypt encrypts a Lockedbuffer and returns the ciphertext as buf.Unsafe
func (r *Cipher) Encrypt(lb *buf.Fenced) *buf.Unsafe {
	switch {
	case r == nil:
		r.SetError("Encrypt() nil receiver")
	case !r.armed:
		r.SetError("Encrypt() not armed")
	case r.gcm == nil:
		r.SetError("Encrypt() nil gcm")
	default:
		b := (*r.gcm).Seal(nil, *r.IV().Buf().(*[]byte), *lb.Buf().(*[]byte), nil)
		B := buf.New().Load(&b).(*buf.Unsafe)
		return B
	}
	b := buf.New()
	b.SetError(r.Error())
	return b
}

// Error returns the error stored in the crypt
func (r *Cipher) Error() string {
	if r == nil {
		return "receiver was nil"
	}
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

// Generate creates a new crypt based on a password and a newly generated random ciphertext
func (r *Cipher) Generate(p *passbuf.Password) *Cipher {
	if r == nil {
		r = New()
		r.SetError("nil receiver")
	}
	if p == nil {
		r.password = new(passbuf.Password)
		r.SetError("no password given")
		return r
	}
	r.password = p
	r.ciphertext = buf.New().Rand(32).(*buf.Fenced)
	r.SetRandomIV()
	r.iterations = Bench(time.Second)
	var C *buf.Fenced
	var IV *buf.Unsafe
	C, IV, r.err = Gen(r.Password(), r.IV(), r.iterations)
	var block cipher.Block
	block, r.err = aes.NewCipher(*C.Buf().(*[]byte))
	var blockmode cipher.AEAD
	blockmode, r.err = cipher.NewGCM(block)
	c := blockmode.Seal(nil, *IV.Buf().(*[]byte), *r.Ciphertext().Buf().(*[]byte), nil)
	r.crypt = r.crypt.Load(&c).(*buf.Unsafe)
	block, r.err = aes.NewCipher(*r.Ciphertext().Buf().(*[]byte))
	A := new(cipher.AEAD)
	a := *A
	a, r.err = cipher.NewGCM(block)
	r.gcm = &a
	r.armed = true
	return r
}

// IV returns the initialisation vector stored in the crypt
func (r *Cipher) IV() *buf.Unsafe {
	if r == nil {
		return new(buf.Unsafe)
	}
	if r.iv == nil {
		r.iv = new(buf.Unsafe)
	}
	return r.iv
}

// IsArmed returns true if the crypt is armed
func (r *Cipher) IsArmed() bool {
	if r == nil {
		return false
	}
	return r.armed
}

// IsLoaded returns true if the crypt contains data
func (r *Cipher) IsLoaded() bool {
	if r == nil || r.crypt == nil {
		return false
	}
	return true
}

// IsUnlocked returns whether the crypt is locked or not
func (r *Cipher) IsUnlocked() bool {
	if r == nil {
		return false
	}
	return r.unlocked
}

// Load moves a bytes into the crypt
func (r *Cipher) Load(b *buf.Unsafe) *Cipher {
	if r == nil {
		r = new(Cipher)
	}
	if r.crypt == nil {
		r.crypt = buf.New()
	}
	r.crypt.Move(b)
	return r
}

// Lock clears the password and disarms the crypt if it is armed
func (r *Cipher) Lock() *Cipher {
	if r == nil {
		r = New()
		r.SetError("nil receciver")
		return r
	}
	if r.password == nil {
		r.password = passbuf.New()
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
func (r *Cipher) MarshalJSON() ([]byte, error) {
	var crypt, ciphertext, iv, password string
	if r.Cipher() != nil && r.Cipher().Len() != 0 {
		crypt = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.Cipher().Buf().(*[]byte)))...))
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
		Cipher     string
		Password   string
		Ciphertext string
		IV         string
		Iterations int64
		Unlocked   bool
		Armed      bool
		HasGCM     bool
		Error      string
	}{
		Cipher:     crypt,
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

// Null wipes the value stored, and restores the buf.Unsafe to the same state as a newly created one (with a nil *[]byte).
func (r *Cipher) Null() *Cipher {
	return New(r)
}

// Password returns the password stored in the Cipher
func (r *Cipher) Password() *passbuf.Password {
	if r == nil {
		r = new(Cipher)
		r.SetError("receiver was nil")
	}
	if r.password == nil {
		r.password = passbuf.New()
		r.password.SetError("password was nil")
	}
	return r.password
}

// SetError sets the error in the Cipher
func (r *Cipher) SetError(s string) *Cipher {
	if r == nil {
		r = new(Cipher)
	}
	r.err = errors.New(s)
	return r
}

// SetIV loads the IV with a buf.Unsafe. It must be 12 bytes long.
func (r *Cipher) SetIV(b *buf.Unsafe) *Cipher {
	if r == nil {
		r = new(Cipher)
	}
	if b == nil {
		b.SetError("nil buf.Unsafe")
	} else if b.Len() != 12 {
		b.SetError("must be 12 bytes")
	}
	if r.iv == nil {
		r.iv = buf.New()
	}
	r.iv.Move(b)
	return r
}

// SetRandomIV loads the IV with a random 12 buf.
func (r *Cipher) SetRandomIV() *Cipher {
	if r == nil {
		r = new(Cipher)
	}
	if r.iv == nil {
		r.iv = new(buf.Unsafe)
	}
	r.iv.Rand(12)
	return r
}

// String prints the JSON representation of the data and structure
func (r *Cipher) String() string {
	s, _ := json.MarshalIndent(r, "", "    ")
	return string(s)
}

// Unlock sets the password, runs the KDF and arms the
func (r *Cipher) Unlock(p *passbuf.Password) *Cipher {
	if r == nil {
		r = new(Cipher)
		r.SetError("nil receiver")
		return r
	}
	if r.password == nil {
		r.password = passbuf.New(r.password)
		r.SetError("nil password")
		return r
	}
	r.password = p
	r.unlocked = true
	return r
}
