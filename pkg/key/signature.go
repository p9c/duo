package key

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/parallelcointeam/duo/pkg/buf"
)

// Sig is a bitcoin EC signature
type Sig struct {
	buf.Byte
}

// NewSig creates a new signature
func NewSig() (out *Sig) {
	return new(Sig)
}

// AsEC returns the signature in ecdsa format
func (r *Sig) AsEC() (out *btcec.Signature) {
	if r == nil {
		r = NewSig().SetStatus(er.NilRec).(*Sig)
		return
	}
	sig, err := btcec.ParseSignature(*r.Bytes(), btcec.S256())
	if r.SetStatusIf(err); err != nil {
		return new(btcec.Signature)
	}
	out = sig
	return
}

// Recover returns a public key with a buffer containing the public key if found and a status indicating if it was successful
func (r *Sig) Recover(h *[]byte) (out *Pub) {
	if r == nil {
		r = NewSig().SetStatus(er.NilRec).(*Sig)
		return NewPub().SetStatus(er.NilRec).(*Pub)
	}
	pub, found, err := btcec.RecoverCompact(btcec.S256(), *r.Bytes(), *h)
	if r.SetStatusIf(err); err != nil {
		return
	}
	if found {
		p := pub.SerializeCompressed()
		r.Copy(&p)
	}
	return
}
