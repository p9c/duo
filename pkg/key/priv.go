package key

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/proto"
)

// NewPriv creates a new Priv
func NewPriv() *Priv {
	return new(Priv)
}

// IsValid returns true if the Priv is currently valid
func (r *Priv) IsValid() bool {
	return r.valid
}

// Invalidate zeroes the key and marks it invalid
func (r *Priv) Invalidate() *Priv {
	if r == nil {
		r = NewPriv().SetStatus(er.NilRec).(*Priv)
	} else {
		r.Zero().Free()
	}
	r.valid = false
	return r
}

// Bytes returns the buffer via the Get function of the Crypt
func (r *Priv) Bytes() (out *[]byte) {
	out = new([]byte)
	switch {
	case r == nil:
		r = NewPriv().SetStatus(er.NilRec).(*Priv)
	case !r.valid:
		r.SetStatus("key invalid")
	default:
		out = r.Get().Bytes()
	}
	return
}

// Copy stores the input buffer using the Put function of the Crypt
func (r *Priv) Copy(in *[]byte) proto.Buffer {
	b := buf.NewByte().Copy(in)
	r.Put(b)
	b.Zero().Free()
	return r
}

// Zero zeroes the key and marks it invalid
func (r *Priv) Zero() proto.Buffer {
	if r == nil {
		return NewPriv().SetStatus(er.NilRec).(*Priv)
	}
	r.valid = false
	r.Zero()
	r.pub.Zero()
	return r
}

// Free frees the crypt inside the Priv and marks it invalid
func (r *Priv) Free() proto.Buffer {
	r.valid = false
	return r.Crypt.Free()
}

// SetKey loads a private key from raw bytes, and zeroes the input bytes of the private key
func (r *Priv) SetKey(priv *[]byte, pub *[]byte) *Priv {
	if r != nil {
		r.Zero().Free()
		r.pub.Zero().Free()
	}
	p := buf.NewByte().Copy(priv)
	r.Put(p)
	p.Zero().Free()
	for i := range *priv {
		(*priv)[i] = 0
	}
	r.pub.Copy(pub)
	return r
}

// Make generates a new private key from random bytes. By default it uses compressed format for the public key, to get another format append a further decompression or hybrid method invocation.
func (r *Priv) Make() *Priv {
	if r != nil {
		r.Zero().Free()
	}
	priv, err := btcec.NewPrivateKey(btcec.S256())
	if r.SetStatusIf(err); err != nil {
		return r
	}
	pr := priv.Serialize()
	r.Copy(&pr)
	pub := priv.PubKey().SerializeCompressed()
	r.pub.Copy(&pub)
	return r
}

// AsEC returns the key in ecdsa.PrivateKey format
func (r *Priv) AsEC() (ecpriv *ecdsa.PrivateKey) {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), *r.Bytes())
	ecpriv = priv.ToECDSA()
	return
}

// PubKey returns a copy of the public key
func (r *Priv) PubKey() proto.Buffer {
	return buf.NewByte().Copy(r.pub.Bytes())
}

// Sign the hash of a message
func (r *Priv) Sign(h *[]byte) (out *Sig) {
	if r == nil {
		r = NewPriv().SetStatus(er.NilRec).(*Priv)
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), *r.Bytes())
	sig, err := priv.Sign(*h)
	if r.SetStatusIf(err); err != nil {
		return
	}
	s := sig.Serialize()
	r.Copy(&s)
	return
}

// SignCompact produces a compact signature for BTC type systems
func (r *Priv) SignCompact(h *[]byte) (out *Sig) {
	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), *r.Bytes())
	sig, err := btcec.SignCompact(btcec.S256(), pk, *h, true)
	if r.SetStatusIf(err); err != nil {
		return
	}
	out.Copy(&sig)
	return
}
