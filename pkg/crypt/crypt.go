// Package crypt is the combination of a BlockCrypt and buf.Byte that stares encrypted and returns decrypted data if a BlockCrypt is loaded.
//
// This is used by the wallet to keep private keys from being copied inside the memory of the application when they are not being worked on.
//
// If there is no BlockCrypt the same code can be used but the data is not protected. This is for the case of a user who does not encrypt their wallet, or has other measures such as VM isolation to protect the wallet process's memory.
package crypt

import (
	"github.com/parallelcointeam/duo/pkg/bc"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
)

var er = core.Errors

// Crypt is a generic structure for storing bytes encrypted and reading them to a another buffer, using a BlockCrypt AES-GCM 256 cipher
type Crypt struct {
	buf.Byte
	BC *bc.BlockCrypt
}

// New returns a new, empty Crypt
func New() *Crypt {
	r := new(Crypt)
	return r
}

// CopyBC creates a new crypt out of the current one like a factory
func (r *Crypt) CopyBC() (R *Crypt) {
	R = New()
	if r.BC != nil {
		R.BC = r.BC
	}
	return
}

// WithBC loads a crypter into the Crypt
func (r *Crypt) WithBC(bc *bc.BlockCrypt) *Crypt {
	switch {
	case r == nil:
		r = New()
		r.SetStatus(er.NilRec)
	case bc == nil:
		r.SetStatus(er.NilParam)
	default:
		r.BC = bc
	}
	return r
}

// Get returns a secure buffer containing the decrypted data in the Crypt
func (r *Crypt) Get() (out core.Buffer) {
	switch {
	case r == nil:
		r = New()
		r.SetStatus(er.NilRec)
		out = buf.NewSecure()
	case r.BC == nil:
		out = buf.NewByte().Copy(r.Byte.Bytes())
	default:
		out = buf.NewSecure().Copy(r.BC.Decrypt(r.Bytes()))
	}
	return out
}

// Put writes a secure buffer to the Crypt
func (r *Crypt) Put(in core.Buffer) *Crypt {
	switch {
	case r == nil:
		r = New()
		r.SetStatus(er.NilRec)
		fallthrough
	case r.BC == nil:
		r.Copy(in.Bytes())
	default:
		r.Copy(r.BC.Encrypt(in.Bytes()))
	}
	return r
}
