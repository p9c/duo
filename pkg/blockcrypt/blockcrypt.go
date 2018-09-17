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

// WithCiphertext loads a ciphertext as from the file where the encryption is used
func (r *BlockCrypt) WithCiphertext(ciphertext *[]byte) *BlockCrypt {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	}
	if r.Ciphertext != nil {
		r.Ciphertext.Val.Destroy()
	}
	if r.Ciphertext.Len() < 32 {
		r.SetStatus("ciphertext is too short")
	}
	var err error
	r.Ciphertext.Copy(ciphertext)
	r.SetStatusIf(err)
	if err != nil {
		return r
	}
	return r
}

// WithRandom loads a random ciphertext as from the file where the encryption is used
func (r *BlockCrypt) WithRandom() *BlockCrypt {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	}
	if r.Ciphertext != nil {
		r.Ciphertext.Val.Destroy()
	}
	var err error
	r.Ciphertext.Val, err = memguard.NewMutableRandom(32)
	r.SetStatusIf(err)
	if err != nil {
		return r
	}
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
	r.SetStatusIf(err)
	if err != nil {
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
	}
	r.Iterations = Bench(time.Second)
	var C *buf.Secure
	var IV *buf.Bytes
	C, IV, err = Gen(p, r.IV, r.Iterations)
	var block cipher.Block
	block, err = aes.NewCipher(*C.Bytes())
	var blockmode cipher.AEAD
	blockmode, err = cipher.NewGCM(block)
	c := blockmode.Seal(nil, *IV.Bytes(), *r.Ciphertext.Bytes(), nil)
	r.Crypt.Copy(&c)
	block, err = aes.NewCipher(*r.Ciphertext.Bytes())
	A := new(cipher.AEAD)
	a := *A
	a, err = cipher.NewGCM(block)
	r.GCM = &a
	r.Armed = true
	return r
}

// Arm generates the correct ciphertext that allows the encrypt/decrypt functions to operate
func (r *BlockCrypt) Arm() *BlockCrypt {
	if r == nil {
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
		return r
	}
	if r.Crypt == nil {
		r.SetStatus("no crypt is loaded")
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
