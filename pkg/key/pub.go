package key

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec"
	"github.com/parallelcointeam/duo/pkg/buf"
	"golang.org/x/crypto/ripemd160"
)

// NewPub creates a new public key
func NewPub() *Pub {
	return new(Pub)
}

// IsCompressed returns true if the key is compressed
func (r *Pub) IsCompressed() bool {
	return btcec.IsCompressedPubKey(*r.Bytes())
}

// Compress converts the key to compressed format if it is in anothter format
func (r *Pub) Compress() *Pub {
	compressed := r.AsCompressed()
	r.Copy(compressed.Bytes())
	return r
}

// Decompress converts the key to compressed format if it is in anothter format
func (r *Pub) Decompress() *Pub {
	decompressed := r.AsUncompressed()
	r.Copy(decompressed.Bytes())
	return r
}

// AsCompressed returns the compressed serialised form (33 bytes, prefix 2 or 3 depending on whether y is odd or even)
func (r *Pub) AsCompressed() (out *buf.Byte) {
	if r == nil {
		r = NewPub().SetStatus(er.NilRec).(*Pub)
	}
	priv, err := btcec.ParsePubKey(*r.Bytes(), btcec.S256())
	if r.SetStatusIf(err); err != nil {
		return buf.NewByte().SetStatus(r.Error()).(*buf.Byte)
	}
	p := priv.SerializeCompressed()
	out = buf.NewByte().Copy(&p).(*buf.Byte)
	return
}

// AsUncompressed returns the uncompressed serialised form (65 bytes with x and y with the prefix 4)
func (r *Pub) AsUncompressed() (out *buf.Byte) {
	if r == nil {
		r = NewPub().SetStatus(er.NilRec).(*Pub)
	}
	priv, err := btcec.ParsePubKey(*r.Bytes(), btcec.S256())
	if r.SetStatusIf(err); err != nil {
		return buf.NewByte().SetStatus(r.Error()).(*buf.Byte)
	}
	p := priv.SerializeUncompressed()
	out = buf.NewByte().Copy(&p).(*buf.Byte)
	return
}

// AsHybrid returns the uncompressed serialised form with the first byte taken from the first bit of the y coordinate, either 0 or 1, with both x and y coordinates (this is not really used)
func (r *Pub) AsHybrid() (out *buf.Byte) {
	if r == nil {
		r = NewPub().SetStatus(er.NilRec).(*Pub)
	}
	priv, err := btcec.ParsePubKey(*r.Bytes(), btcec.S256())
	if r.SetStatusIf(err); err != nil {
		return buf.NewByte().SetStatus(r.Error()).(*buf.Byte)
	}
	p := priv.SerializeHybrid()
	out = buf.NewByte().Copy(&p).(*buf.Byte)
	return
}

// AsEC returns the EC public key
func (r *Pub) AsEC() (out *ecdsa.PublicKey) {
	if r == nil {
		r = NewPub().SetStatus(er.NilRec).(*Pub)
		return new(ecdsa.PublicKey)
	}
	priv, err := btcec.ParsePubKey(*r.Bytes(), btcec.S256())
	if r.SetStatusIf(err); err != nil {
		return
	}
	return priv.ToECDSA()
}

// ID returns the ripemd160 hash of the public key
func (r *Pub) ID() (out *buf.Byte) {
	if r == nil {
		r = NewPub().SetStatus(er.NilRec).(*Pub)
	}
	if r.Val == nil {
		r.SetStatus(er.NilBuf)
		return
	}
	h := ripemd160.New()
	h.Write(*r.Bytes())
	o := h.Sum(nil)
	out.Copy(&o)
	return
}
