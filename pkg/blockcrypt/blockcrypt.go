package blockcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/awnumar/memguard"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/proto"
	"time"
)

var er = proto.Errors

// New creates a new, empty BlockCrypt
func New() *BlockCrypt {
	return new(BlockCrypt)
}

// WithCiphertext loads a ciphertext, IV and iterations as from the file where the encryption is used
func (r *BlockCrypt) WithCiphertext(ciphertext *[]byte, iv *[]byte, iterations int) *BlockCrypt {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	}
	if r.Ciphertext != nil {
		r.Ciphertext.Val.Destroy()
	}
	if r.Ciphertext.Len() != 32 {
		r.SetStatus("ciphertext is not 32 bytes")
	}
	var err error
	r.Ciphertext.Copy(ciphertext)
	r.SetStatusIf(err)
	if err != nil {
		return r
	}
	if len(*iv) != 12 {
		r.SetStatus("IV is not 12 bytes long")
		return r
	}
	r.IV.Copy(iv)
	if iterations < 0 {
		r.SetStatus("iterations is negative")
		return r
	}
	r.Iterations = iterations
	return r
}

// Generate creates a new crypt based on a password and a newly generated random ciphertext
func (r *BlockCrypt) Generate(p *buf.Secure) *BlockCrypt {
	if r == nil {
		r = New()
		r.SetStatus(er.NilRec)
	}
	if p == nil {
		r.SetStatus("no password given")
		return r
	}
	r.Password.Val = p.Val
	var err error
	r.Ciphertext.Val, err = memguard.NewMutableRandom(32)
	if r.SetStatusIf(err); err != nil {
		return r
	}
	bb := make([]byte, 12)
	r.IV.Copy(&bb)
	n, err := rand.Read(*r.IV.Val)
	r.SetStatusIf(err)
	if err != nil {
		return r
	}
	if n != 12 {
		r.SetStatus("did not get requested 12 random bytes")
		return r
	}
	r.Iterations = Bench(time.Second)

	return r
}

// Arm generates the correct ciphertext that allows the encrypt/decrypt functions to operate
func (r *BlockCrypt) Arm() *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
		return r
	case r.Crypt == nil:
		r.SetStatus("no crypt is loaded")
	case r.Password == nil:
		r.SetStatus("no password given")
	default:
		var err error
		if r.Ciphertext != nil {
			r.Ciphertext.Free()
		}
		r.Ciphertext, r.GCMIV, err = Gen(r.Password, r.IV, r.Iterations)
		if r.SetStatusIf(err); err != nil {
			return r
		}
		var block cipher.Block
		block, err = aes.NewCipher(*r.Ciphertext.Bytes())
		if r.SetStatusIf(err); err != nil {
			return r
		}
		var blockmode cipher.AEAD
		blockmode, err = cipher.NewGCM(block)
		if r.SetStatusIf(err); err != nil {
			return r
		}
		c, err := blockmode.Open(nil, *r.IV.Bytes(), *r.Crypt.Bytes(), nil)
		if r.SetStatusIf(err); err != nil {
			return r
		}
		r.Ciphertext.Copy(&c)
		for i := range c {
			c[i] = 0
		}
		block, err = aes.NewCipher(*r.Ciphertext.Bytes())
		r.SetStatusIf(err)
		if err != nil {
			return r
		}
		a, err := cipher.NewGCM(block)
		if r.SetStatusIf(err); err != nil {
			return r
		}
		r.GCM = &a
		r.Armed = true
	}
	return r
}

// Unlock loads the password, enabling arming of the crypt
func (r *BlockCrypt) Unlock(pass proto.Buffer) *BlockCrypt {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	}
	if r.Password != nil {
		r.Password.Zero().Free()
	}
	var err error
	r.Password.Copy(pass.Bytes())
	r.SetStatusIf(err)
	if err != nil {
		return r
	}
	// TODO
	return r
}

// Lock clears the password
func (r *BlockCrypt) Lock() *BlockCrypt {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	}
	if r.Password != nil {
		r.Password.Zero().Free()
	}
	return r
}

// Disarm re-encrypts the ciphertext, destroys the cipher and unmarks the armed flag
func (r *BlockCrypt) Disarm() *BlockCrypt {
	panic("not implemented")
}

// Encrypt uses the armed cipher to encrypt a buffer
func (r *BlockCrypt) Encrypt(buf proto.Buffer) (out proto.Buffer) {
	panic("not implemented")
}

// Decrypt uses the armed cipher to decrypt a buffer
func (r *BlockCrypt) Decrypt(buf proto.Buffer) (out proto.Buffer) {
	panic("not implemented")
}

// Status implementation

// SetStatus sets the status of the crypt
func (r *BlockCrypt) SetStatus(s string) proto.Status {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	} else {
		r.Status = s
	}
	return r
}

// SetStatusIf sets the status according to an error output
func (r *BlockCrypt) SetStatusIf(err error) proto.Status {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	} else {
		if err != nil {
			r.Status = err.Error()
		}
	}
	return r
}

// UnsetStatus clears the error state
func (r *BlockCrypt) UnsetStatus() proto.Status {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	} else {
		r.Status = ""
	}
	return r
}

// OK returns true if the status text is empty
func (r *BlockCrypt) OK() bool {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
		return false
	}
	return r.Status == ""
}

// Error implementation

// Error implemennts the Error() interface
func (r *BlockCrypt) Error() string {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	}
	return r.Status
}
