package key

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/parallelcointeam/duo/pkg/bc"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
)

// NewPriv creates a new Priv
func NewPriv() (priv *Priv) {
	priv = new(Priv)
	priv.pub = NewPub()
	return
}

// NewIf creates a new Pub if the receiver is nil
func (r *Priv) NewIf() *Priv {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
	}
	return r
}

// WithBC copies in the reference to a BlockCrypt to enable encryption
func (r *Priv) WithBC(bc *bc.BlockCrypt) *Priv {
	r = r.NewIf()
	r.Crypt.WithBC(bc)
	return r
}

// IsValid returns true if the Priv is currently valid
func (r *Priv) IsValid() bool {
	r = r.NewIf()
	return r.valid
}

// Invalidate zeroes the key and marks it invalid
func (r *Priv) Invalidate() *Priv {
	r = r.NewIf()
	switch {
	default:
		r.Zero().Free()
		r = NewPriv()
	}
	r.valid = false
	return r
}

// Bytes returns the buffer via the Get function of the Crypt
func (r *Priv) Bytes() (out *[]byte) {
	r = r.NewIf()
	out = new([]byte)
	switch {
	case !r.valid:
		r.SetStatus("key invalid")
	default:
		out = r.Get().Bytes()
	}
	return
}

// Hex returns the hex representation of the contennt of the crypt via the get function
func (r *Priv) Hex() (out string) {
	r = r.NewIf()
	switch {
	case !r.valid:
		r.SetStatus("key invalid")
	default:
		out = hex.EncodeToString(*r.Get().Bytes())
	}
	return
}

// Copy stores the input buffer using the Put function of the Crypt
func (r *Priv) Copy(in *[]byte) core.Buffer {
	r = r.NewIf()
	b := buf.NewByte().Copy(in)
	r.Crypt.Put(b)
	return r
}

// Zero zeroes the key and marks it invalid
func (r *Priv) Zero() core.Buffer {
	r = r.NewIf()
	switch {
	case r.pub != nil:
		r.pub.Zero()
		fallthrough
	default:
		r.Crypt.Free()
		r.valid = false
	}
	return r
}

// Free frees the crypt inside the Priv and marks it invalid
func (r *Priv) Free() core.Buffer {
	r = r.NewIf()
	switch {
	case r != nil:
		r.Crypt.Free()
		fallthrough
	case r.pub != nil:
		r.pub.Free()
	}
	r.valid = false
	return r
}

// SetKey loads a private key from raw bytes, and zeroes the input bytes of the private key
func (r *Priv) SetKey(priv *[]byte, pub *[]byte) (out *Priv) {
	out = new(Priv)
	r = r.NewIf()
	r.UnsetStatus()
	if r.Copy(priv).(*Priv).OK() {
		if r.pub.Copy(pub).(*Pub).OK() {
			r.valid = true
		}
	}
	return r
}

// Make generates a new private key from random bytes. By default it uses compressed format for the public key, to get another format append a further decompression or hybrid method invocation.
func (r *Priv) Make() *Priv {
	r = r.NewIf()
	if priv, err := btcec.NewPrivateKey(btcec.S256()); r.SetStatusIf(err).OK() {
		pr := priv.Serialize()
		r.Crypt.Put(buf.NewByte().Copy(&pr))
		pub := priv.PubKey().SerializeCompressed()
		r.pub = NewPub()
		r.pub.Copy(&pub)
		r.valid = true
	}
	return r
}

// AsEC returns the key in ecdsa.PrivateKey format
func (r *Priv) AsEC() (ecpriv *ecdsa.PrivateKey) {
	if r == nil {
		r = r.NewIf()
		return new(ecdsa.PrivateKey)
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), *r.Bytes())
	ecpriv = priv.ToECDSA()
	return
}

// PubKey returns a copy of the public key
func (r *Priv) PubKey() core.Buffer {
	if r == nil {
		r = r.NewIf()
		return buf.NewByte()
	}
	return buf.NewByte().Copy(r.pub.Bytes())
}

// Sign the hash of a message
func (r *Priv) Sign(h *[]byte) (out *Sig) {
	if r == nil {
		r = r.NewIf()
		return &Sig{}
	}
	priv, pub := btcec.PrivKeyFromBytes(btcec.S256(), *r.Bytes())
	sig, err := priv.Sign(*h)
	if r.SetStatusIf(err); err == nil {
		if sig != nil {
			s := sig.Serialize()
			out = NewSig()
			out.Copy(&s)
			out.mh = buf.NewByte().Copy(h).(*buf.Byte)
			p := pub.SerializeCompressed()
			pp := NewPub().Copy(&p).(*Pub)
			out.addr = pp.GetID()
		}
	}
	return
}

// SignCompact produces a compact signature for BTC type systems
func (r *Priv) SignCompact(h *[]byte) (out *Sig) {
	if r == nil {
		r = r.NewIf()
		return &Sig{}
	}
	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), *r.Bytes())
	sig, err := btcec.SignCompact(btcec.S256(), pk, *h, true)
	if r.SetStatusIf(err); err == nil {
		if sig != nil {
			out = NewSig()
			out.Copy(&sig)
			out.mh.Copy(h)
		}
	}
	return
}

// GetID returns the hash160 ID of the public key
func (r *Priv) GetID() core.Address {
	if r == nil {
		r = r.NewIf()
		return ""
	}
	return NewID(r.pub.Bytes())
}
