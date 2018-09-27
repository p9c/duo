package bc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/awnumar/memguard"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/kdf"
	"github.com/parallelcointeam/duo/pkg/proto"
	"time"
)

// New creates a new, empty BlockCrypt
func New() *BlockCrypt {
	r := new(BlockCrypt)
	return r
}

// Generate creates a new crypt based on a password and a newly generated random ciphertext.
//
// After this function the crypt is unlocked, the crypt contains the encrypted ciphertext and the ciphertext is destroyed and Arm() must be called to activate the cipher
func (r *BlockCrypt) Generate(p *buf.Secure) *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
		fallthrough
	case r.Ciphertext != nil:
		r.Ciphertext.Free()
		fallthrough
	case r.Password != nil:
		r.Password.Free()
		fallthrough
	case r.Crypt != nil:
		r.Crypt.Free()
		fallthrough
	case p == nil:
		r.SetStatus("no password given")
	default:
		r.Password = p
		var err error
		r.Ciphertext = buf.NewSecure()
		r.Ciphertext.Val, err = memguard.NewMutableRandom(32)
		if r.SetStatusIf(err); err == nil {
			bb := make([]byte, 12)
			n, err := rand.Read(bb)
			if r.SetStatusIf(err); n == 12 {
				if err == nil {
					r.IV = buf.NewByte()
					r.IV.Copy(&bb)
					r.Iterations = kdf.Bench(time.Second / 100)
					var C *buf.Secure
					C, err = kdf.Gen(r.Password, r.IV, r.Iterations)
					if r.SetStatusIf(err); err == nil {
						var block cipher.Block
						block, err = aes.NewCipher(C.Val.Buffer())
						if r.SetStatusIf(err); err == nil {
							var blockmode cipher.AEAD
							blockmode, err = cipher.NewGCM(block)
							if r.SetStatusIf(err); err == nil {
								c := blockmode.Seal(nil, *r.IV.Bytes(), *r.Ciphertext.Bytes(), nil)
								if r.SetStatusIf(err); err == nil {
									r.Crypt = buf.NewByte()
									r.Crypt.Copy(&c)
									r.Unlocked = true
									r.Ciphertext.Free()
									r.Armed = false
									r.UnsetStatus()
								}
							}
						}
					}
				}
			}
		}
	}
	return r
}

// LoadCrypt loads a crypt, IV and iterations as from the file where the encryption is used.
//
// The password still needs to be loaded to unlock the crypt and the crypt unlocked to arm the BlockCrypt. This function clears any existing data in the Blockcrypt.
func (r *BlockCrypt) LoadCrypt(crypt *[]byte, iv *[]byte, iterations int64) *BlockCrypt {
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
			fallthrough
		default:
			r.IV = buf.NewByte()
			r.IV.Copy(iv)
			r.Crypt = buf.NewByte()
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
		r = New()
		r.SetStatus(er.NilRec)
	case r.Password != nil:
		r.Password.Zero().Free()
		fallthrough
	default:
		r.Password = pass
		r.Unlocked = true
		r.decryptCrypt()
	}
	return r
}

// Lock clears the password and disarms the crypt
func (r *BlockCrypt) Lock() *BlockCrypt {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	case r.Password != nil:
		r.Password.Zero().Free()
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
	passCiphertext, err := kdf.Gen(r.Password, r.IV, r.Iterations)
	if r.SetStatusIf(err); err == nil {
		block, err := aes.NewCipher(*passCiphertext.Bytes())
		if r.SetStatusIf(err); err == nil {
			blockmode, err := cipher.NewGCM(block)
			if r.SetStatusIf(err); err == nil {
				c, err := blockmode.Open(nil, *r.IV.Bytes(), *r.Crypt.Bytes(), nil)
				if r.SetStatusIf(err); err == nil {
					if r.Ciphertext != nil {
						r.Ciphertext.Free()
					}
					r.Ciphertext = buf.NewSecure()
					r.Ciphertext.Copy(&c)
					proto.Zero(&c)
				}
			}
		}
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
		if r.SetStatusIf(err); err == nil {
			a, err := cipher.NewGCM(block)
			if r.SetStatusIf(err); err == nil {
				r.GCM = &a
				r.Armed = true
			}
		}
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
		o := (*r.GCM).Seal(nil, *r.IV.Bytes(), *buf, nil)
		out = &o
		r.UnsetStatus()
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
		o, err := (*r.GCM).Open(nil, *r.IV.Bytes(), *buf, nil)
		if r.SetStatusIf(err); err == nil {
			out = &o
		}
		r.UnsetStatus()
	}
	return
}

// Status implementation

// SetStatus sets the status of the crypt
func (r *BlockCrypt) SetStatus(s string) proto.Status {
	switch {
	case r == nil:
		r = New().SetStatus(er.NilRec).(*BlockCrypt)
	default:
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
	default:
		r.UnsetStatus()
	}
	return r
}

// UnsetStatus clears the error state
func (r *BlockCrypt) UnsetStatus() proto.Status {
	switch {
	case r == nil:
		r = New()
		r.SetStatus(er.NilRec)
	default:
		r.Status = ""
	}
	return r
}

// OK returns true if the status text is empty
func (r *BlockCrypt) OK() bool {
	if r == nil {
		r = New()
		r.SetStatus(er.NilRec)
		return false
	}
	return r.Status == ""
}

// Error implementation

// Error implemennts the Error() interface
func (r *BlockCrypt) Error() string {
	if r == nil {
		r = New()
		r.SetStatus(er.NilRec)
	}
	return r.Status
}
