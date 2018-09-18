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

// Generate creates a new crypt based on a password and a newly generated random ciphertext.
//
// After this function the crypt is unlocked, the crypt contains the encrypted ciphertext and the ciphertext is destroyed and Arm() must be called to activate the cipher
func (r *BlockCrypt) Generate(p *buf.Secure) *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	case p == nil:
		r.SetStatus("no password given")
	case r.Password != nil:
		r.Password.Free()
		fallthrough
	default:
		r.Password = p
		var err error
		r.Ciphertext.Val, err = memguard.NewMutableRandom(32)
		if r.SetStatusIf(err); err != nil {
			return r
		}
		bb := make([]byte, 12)
		n, err := rand.Read(bb)
		switch {
		case err != nil:
			r.SetStatusIf(err)
		case n != 12:
			r.SetStatus("did not get requested 12 random bytes")
		default:
			r.IV.Copy(&bb)
			r.Iterations = Bench(time.Second)
			var C *buf.Secure
			C, err = Gen(r.Password, r.IV, r.Iterations)
			var block cipher.Block
			block, err = aes.NewCipher(C.Val.Buffer())
			if r.SetStatusIf(err); err != nil {
				return r
			}
			var blockmode cipher.AEAD
			blockmode, err = cipher.NewGCM(block)
			if r.SetStatusIf(err); err != nil {
				return r
			}
			c := blockmode.Seal(nil, *r.IV.Bytes(), *r.Ciphertext.Bytes(), nil)
			if r.SetStatusIf(err); err != nil {
				return r
			}
			r.Crypt.Copy(&c)
			r.Unlocked = true
			r.Ciphertext.Free()
			r.Armed = false
			r.UnsetStatus()
		}
	}
	return r
}

// LoadCrypt loads a crypt, IV and iterations as from the file where the encryption is used.
//
// The password still needs to be loaded to unlock the crypt and the crypt unlocked to arm the BlockCrypt. This function clears any existing data in the Blockcrypt.
func (r *BlockCrypt) LoadCrypt(crypt *[]byte, iv *[]byte, iterations int) *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
		fallthrough
	case crypt == nil:
		r.SetStatus("crypt parameter is nil")
	case iv == nil:
		r.SetStatus("IV parameter is nil")
	case len(*iv) != 12:
		r.SetStatus("IV is not 12 bytes long")
	case iterations < 0:
		r.SetStatus("iterations is negative")
	case iterations == 0:
		r.SetStatus("iterations is zero")
	case r.Ciphertext != nil:
		r.Ciphertext.Free()
		fallthrough
	default:
		switch {
		case r.Password != nil:
			r.Password.Free()
			r.Password = nil
			r.Unlocked = false
			fallthrough
		case r.Ciphertext != nil:
			r.Ciphertext.Free()
			r.Ciphertext = nil
			r.Armed = false
			fallthrough
		case r.GCM != nil:
			r.GCM = nil
		default:
			r.IV.Copy(iv)
			r.Crypt.Copy(crypt)
			r.Iterations = iterations
			r.UnsetStatus()
		}
	}
	return r
}

// Unlock loads the password, enabling arming of the crypt
func (r *BlockCrypt) Unlock(pass *buf.Secure) *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	case r.Password != nil:
		r.Password.Free()
		fallthrough
	default:
		r.Password = pass
		r.Unlocked = true
		r.UnsetStatus()
	}
	return r
}

// Lock clears the password and disarms the crypt
func (r *BlockCrypt) Lock() *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	case r.Password != nil:
		r.Password.Free()
		r.Password = nil
		r.Unlocked = false
		fallthrough
	default:
		r.Disarm()
	}
	return r
}

func (r *BlockCrypt) decryptCrypt() *BlockCrypt {
	if r.Ciphertext != nil {
		r.Ciphertext.Free()
	}
	passCiphertext, err := Gen(r.Password, r.IV, r.Iterations)
	if r.SetStatusIf(err); err != nil {
		return r
	}
	block, err := aes.NewCipher(*passCiphertext.Bytes())
	if r.SetStatusIf(err); err != nil {
		return r
	}
	blockmode, err := cipher.NewGCM(block)
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
	return r
}

// Arm generates the correct ciphertext that allows the encrypt/decrypt functions to operate
func (r *BlockCrypt) Arm() *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	case r.Iterations < 1:
		r.SetStatus("iterations are less than 1")
	case r.Crypt == nil:
		r.SetStatus("no crypt is loaded")
	case r.Password == nil:
		r.SetStatus("no password loaded")
	default:
		r.decryptCrypt()
		block, err := aes.NewCipher(*r.Ciphertext.Bytes())
		if r.SetStatusIf(err); err != nil {
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

// Disarm re-encrypts the ciphertext, destroys the cipher and unmarks the armed flag
func (r *BlockCrypt) Disarm() *BlockCrypt {
	r.Ciphertext.Free()
	r.Ciphertext = nil
	r.GCM = nil
	r.Armed = false
	r.UnsetStatus()
	return r
}

// Encrypt uses the armed cipher to encrypt a buffer
func (r *BlockCrypt) Encrypt(buf *[]byte) (out *[]byte) {
	out = &[]byte{}
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	case buf == nil:
		r.SetStatus(er.NilParam)
	case len(*buf) < 1:
		r.SetStatus(er.ZeroLen)
	case r.GCM == nil:
		r.SetStatus("cipher was not found")
	default:
		o := (*r.GCM).Seal(nil, *buf, *r.IV.Bytes(), nil)
		out = &o
	}
	return
}

// Decrypt uses the armed cipher to decrypt a buffer
func (r *BlockCrypt) Decrypt(buf *[]byte) (out *[]byte) {
	out = &[]byte{}
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	case buf == nil:
		r.SetStatus(er.NilParam)
	case len(*buf) < 1:
		r.SetStatus(er.ZeroLen)
	case r.GCM == nil:
		r.SetStatus("cipher was not found")
	default:
		o, err := (*r.GCM).Open(nil, *buf, *r.IV.Bytes(), nil)
		if r.SetStatusIf(err); err == nil {
			out = &o
		}
	}
	return
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
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
		fallthrough
	case err != nil:
		r.Status = err.Error()
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
