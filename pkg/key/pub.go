package key

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcec"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
)

// NewPub creates a new public key
func NewPub() *Pub {
	r := new(Pub)
	return r
}

// NewIf creates a new Pub if the receiver is nil
func (r *Pub) NewIf() *Pub {
	if r == nil {
		r = NewPub()
		r.SetStatus(er.NilRec)
	}
	return r
}

// Bytes returns the private key
func (r *Pub) Bytes() (out *[]byte) {
	r = r.NewIf()
	out = r.Byte.Bytes()
	return
}

// Copy loads the public key
func (r *Pub) Copy(in *[]byte) core.Buffer {
	r = r.NewIf()
	switch {
	case in == nil:
		r.SetStatus(er.NilParam)
	case len(*in) < 1:
		r.SetStatus(er.ZeroLen)
	default:
		out := make([]byte, len(*in))
		copy(out, *in)
		r.Byte.Val = &out
	}
	return r
}

// Zero wipes the key
func (r *Pub) Zero() core.Buffer {
	r = r.NewIf()
	switch {
	case r.Len() < 1:
		r.SetStatus(er.ZeroLenBuf)
	default:
		r.Byte.Zero()
	}
	return r
}

// Free deallocates the buffer of the key
func (r *Pub) Free() core.Buffer {
	r = r.NewIf()
	switch {
	case r.Len() < 1:
		r.SetStatus(er.ZeroLenBuf)
	default:
		r.Byte.Free()
	}
	return r
}

// IsCompressed returns true if the key is compressed
func (r *Pub) IsCompressed() bool {
	if r == nil {
		r = r.NewIf()
		return false
	}
	return btcec.IsCompressedPubKey(*r.Bytes())
}

// Compress converts the key to compressed format if it is in another format
func (r *Pub) Compress() *Pub {
	if r == nil {
		r = r.NewIf()
		return r
	}
	compressed := r.AsCompressed()
	r.Copy(compressed.Bytes())
	return r
}

// Decompress converts the key to compressed format if it is in anothter format
func (r *Pub) Decompress() *Pub {
	r = r.NewIf()
	switch {
	default:
		decompressed := r.AsUncompressed()
		r.Copy(decompressed.Bytes())
	}
	return r
}

// AsCompressed returns the compressed serialised form (33 bytes, prefix 2 or 3 depending on whether y is odd or even)
func (r *Pub) AsCompressed() (out *buf.Byte) {
	if r == nil {
		r = r.NewIf()
		return &buf.Byte{}
	}
	priv, err := btcec.ParsePubKey(*r.Bytes(), btcec.S256())
	if !r.SetStatusIf(err).OK() {
		out = buf.NewByte()
		out.SetStatus(r.Error())
		return
	}
	p := priv.SerializeCompressed()
	out = buf.NewByte().Copy(&p).(*buf.Byte)
	return
}

// AsUncompressed returns the uncompressed serialised form (65 bytes with x and y with the prefix 4)
func (r *Pub) AsUncompressed() (out *buf.Byte) {
	if r == nil {
		r = r.NewIf()
		return &buf.Byte{}
	}
	priv, err := btcec.ParsePubKey(*r.Bytes(), btcec.S256())
	if !r.SetStatusIf(err).OK() {
		out = buf.NewByte()
		out.SetStatus(r.Error())
		return
	}
	p := priv.SerializeUncompressed()
	out = buf.NewByte().Copy(&p).(*buf.Byte)
	return
}

// AsHybrid returns the uncompressed serialised form with the first byte taken from the first bit of the y coordinate, either 0 or 1, with both x and y coordinates (this is not really used)
func (r *Pub) AsHybrid() (out *buf.Byte) {
	if r == nil {
		r = r.NewIf()
		return &buf.Byte{}
	}
	priv, err := btcec.ParsePubKey(*r.AsUncompressed().Bytes(), btcec.S256())
	if !r.SetStatusIf(err).OK() {
		out = buf.NewByte()
		out.SetStatus(r.Error())
		return
	}
	p := priv.SerializeHybrid()
	out = buf.NewByte().Copy(&p).(*buf.Byte)
	return
}

// AsEC returns the EC public key
func (r *Pub) AsEC() (out *ecdsa.PublicKey) {
	if r == nil {
		r = r.NewIf()
		return new(ecdsa.PublicKey)
	}
	priv, err := btcec.ParsePubKey(*r.Bytes(), btcec.S256())
	if !r.SetStatusIf(err).OK() {
		return new(ecdsa.PublicKey)
	}
	return priv.ToECDSA()
}

// GetID returns the hash160 ID of the public key
func (r *Pub) GetID() core.Address {
	if r == nil {
		r = r.NewIf()
		return ""
	}
	return NewID(r.Bytes())
}
