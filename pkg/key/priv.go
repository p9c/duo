package key

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/crypt"
	"github.com/parallelcointeam/duo/pkg/proto"
)

// NewPriv creates a new Priv
func NewPriv() (priv *Priv) {
	priv = new(Priv)
	priv.Crypt = crypt.New()
	priv.pub = NewPub()
	return
}

// WithBC copies in the reference to a BlockCrypt to enable encryption
func (r *Priv) WithBC(bc *blockcrypt.BlockCrypt) *Priv {
	switch {
	case r == nil:
		r = NewPriv()
		r.SetStatus(er.NilRec)
	default:
		r.Crypt = crypt.New().WithBC(bc)
	}
	return r
}

// IsValid returns true if the Priv is currently valid
func (r *Priv) IsValid() bool {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
	}
	return r.valid
}

// Invalidate zeroes the key and marks it invalid
func (r *Priv) Invalidate() *Priv {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
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
		r = NewPriv()
		r.SetStatus(er.NilRec)
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
	r.Crypt.Put(b)
	return r
}

// Zero zeroes the key and marks it invalid
func (r *Priv) Zero() proto.Buffer {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
	} else {
		r.Crypt.Free()
		if r.pub != nil {
			r.pub.Zero()
		}
	}
	r.valid = false
	return r
}

// Free frees the crypt inside the Priv and marks it invalid
func (r *Priv) Free() proto.Buffer {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
	} else {
		r.Crypt.Free()
		if r.pub != nil {
			r.pub.Free()
		}
	}
	r.valid = false
	return r
}

// SetKey loads a private key from raw bytes, and zeroes the input bytes of the private key
func (r *Priv) SetKey(priv *[]byte, pub *[]byte) *Priv {
	if r != nil {
		r.Zero().Free()
		r.pub.Zero().Free()
	}
	r.Copy(priv)
	for i := range *priv {
		(*priv)[i] = 0
	}
	r.pub.Copy(pub)
	r.valid = true
	return r
}

// Make generates a new private key from random bytes. By default it uses compressed format for the public key, to get another format append a further decompression or hybrid method invocation.
func (r *Priv) Make() *Priv {
	if r == nil {
		r = NewPriv()
	}
	priv, err := btcec.NewPrivateKey(btcec.S256())
	if r.SetStatusIf(err); err == nil {
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
		return new(ecdsa.PrivateKey)
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), *r.Bytes())
	ecpriv = priv.ToECDSA()
	return
}

// PubKey returns a copy of the public key
func (r *Priv) PubKey() proto.Buffer {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
		return buf.NewByte()
	}
	return buf.NewByte().Copy(r.pub.Bytes())
}

// Sign the hash of a message
func (r *Priv) Sign(h *[]byte) (out *Sig) {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
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
		r = NewPriv()
		r.SetStatus(er.NilRec)
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
func (r *Priv) GetID() proto.ID {
	if r == nil {
		r = NewPriv()
		r.SetStatus(er.NilRec)
		return ""
	}
	return NewID(r.pub.Bytes())
}
