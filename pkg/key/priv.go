package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/crypt"
	"github.com/parallelcointeam/duo/pkg/proto"
	"math/big"
)

// Priv is a private key, stored in a Crypt
type Priv struct {
	crypt.Crypt
	Pub        *buf.Byte
	valid      bool
	compressed bool
}

// NewPriv creates a new Priv
func NewPriv() *Priv {
	return new(Priv)
}

// IsValid returns true if the Priv is currently valid
func (r *Priv) IsValid() bool {
	return r.valid
}

// IsCompressed returns true if the Priv is currently compressed
func (r *Priv) IsCompressed() bool {
	return r.compressed
}

// Invalidate zeroes the key and marks it invalid
func (r *Priv) Invalidate() *Priv {
	if r == nil {
		r = NewPriv()
	} else {
		r.Zero().Free()
	}
	r.valid = false
	return r
}

// SetCompress sets the key to generate compressed public keys
func (r *Priv) SetCompress() *Priv {
	r.compressed = true
	return r
}

// UnsetCompress sets the key to generate uncompressed public keys
func (r *Priv) UnsetCompress() *Priv {
	r.compressed = false
	return r
}

// Bytes returns the buffer via the Get function of the Crypt
func (r *Priv) Bytes() (out *[]byte) {
	out = r.Crypt.Get().Bytes()
	return
}

// Copy stores the input buffer using the Put function of the Crypt
func (r *Priv) Copy(in *[]byte) proto.Buffer {
	b := buf.NewSecure().Copy(in)
	r.Crypt.Put(b)
	b.Free()
	return r
}

// Zero zeroes the key and marks it invalid
func (r *Priv) Zero() proto.Buffer {
	r.valid = false
	r.Crypt.Zero()
	return r
}

// Free frees the crypt inside the Priv and marks it invalid
func (r *Priv) Free() proto.Buffer {
	r.valid = false
	return r.Crypt.Free()
}

// SetPrivKey loads a private key from raw bytes, and zeroes the input bytes
func (r *Priv) SetPrivKey(priv *[]byte, compressed bool) *Priv {
	r.Copy(priv)
	for i := range *priv {
		(*priv)[i] = 0
	}
	r.compressed = compressed
	return r
}

// Make generates a new private key from random bytes
func (r *Priv) Make() *Priv {
	if r != nil {
		r.Zero().Free()
	}
	priv, x, y, err := elliptic.GenerateKey(
		elliptic.P256(),
		rand.Reader)
	if r.SetStatusIf(err); err != nil {
		return r
	}
	keyBytes := elliptic.Marshal(elliptic.P256(), x, y)
	r.Pub.Copy(&keyBytes)
	r.Copy(&priv)
	for i := range priv {
		priv[i] = 0
	}
	return r
}

// GetEC returns the key in ecdsa.PrivateKey format
func (r *Priv) GetEC() (priv *ecdsa.PrivateKey) {
	x, y := elliptic.Unmarshal(elliptic.P256(), *r.Crypt.Bytes())
	bi := big.NewInt(0)
	bi.SetBytes(*r.Bytes())
	r.Bytes()
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(), X: x, Y: y},
		D: bi,
	}
}

// GetPubKey returns the public key
func (r *Priv) GetPubKey() proto.Buffer {
	return buf.NewByte().Copy(r.Pub.Bytes())
}
