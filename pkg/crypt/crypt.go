// Package crypt is the combination of a BlockCrypt and buf.Byte that stares encrypted and returns decrypted data if a BlockCrypt is loaded.
//
// This is used by the wallet to keep private keys from being copied inside the memory of the application when they are not being worked on.
//
// If there is no BlockCrypt the same code can be used but the data is not protected. This is for the case of a user who does not encrypt their wallet, or has other measures such as VM isolation to protect the wallet process's memory.
package crypt

import (
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/proto"
)

var er = proto.Errors

// Crypt is a generic structure for storing bytes encrypted and reading them to a another buffer, using a BlockCrypt AES-GCM 256 cipher
type Crypt struct {
	BC *blockcrypt.BlockCrypt
	buf.Byte
}

// New returns a new, empty Crypt
func New() *Crypt {
	return new(Crypt)
}

// New creates a new crypt out of the current one like a factory
func (r *Crypt) New() (R *Crypt) {
	R = New()
	if r.BC != nil {
		R.BC = r.BC
	}
	return
}

// WithBlockCrypt loads a crypter into the Crypt
func (r *Crypt) WithBlockCrypt(bc *blockcrypt.BlockCrypt) *Crypt {
	if r == nil {
		r = New()
	}
	if bc == nil {
		r.SetStatus(er.NilParam)
	} else {
		r.BC = bc
	}
	return r
}

// Get returns a secure buffer containing the decrypted data in the Crypt
func (r *Crypt) Get() (out proto.Buffer) {
	switch {
	case r == nil:
		r = r.SetStatus(er.NilRec).(*Crypt)
		out = buf.NewSecure()
	case r.BC == nil:
		out = buf.NewBytes().Copy(r.Byte.Byte())
	default:
		out = buf.NewSecure().Copy(r.BC.Decrypt(r.Byte()))
	}
	return out
}

// Put writes a secure buffer to the Crypt
func (r *Crypt) Put(in proto.Buffer) *Crypt {
	switch {
	case r == nil:
		r = r.SetStatus(er.NilRec).(*Crypt)
	case r.BlockCrypt == nil:
		r.Crypt.Copy(in.Byte())
	default:
		r.Crypt = in.Copy(r.Encrypt(in.Byte())).(*buf.Byte)
	}
	return r
}
