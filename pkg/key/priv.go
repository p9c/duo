package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/ec"
)

// Priv is an encapsulated private key
type Priv struct {
	priv                *ec.PrivateKey
	pub                 *ec.PublicKey
	invalid, compressed bool
}

// New makes a new key pair
func (p *Priv) New() {
	bytes := make([]byte, 32)
	for {
		rand.Read(bytes)
		if Check(bytes) {
			break
		}
	}
	p.Set(bytes)
	return
}

type priv interface {
	New()
	Get() []byte
	Set([]byte) *Priv
	Size() int
	Invalidate()
	IsValid() bool
	SetPriv(*Priv) bool
	GetPriv() *Priv
	SetPub(*Pub)
	GetPub() *Pub
	Sign(Uint.U256, []byte) bool
	SignCompact(Uint.U256, []byte) bool
	Verify(hash Uint.U256, S []byte)
}

// Get - gets the full private key as a byte slice
func (p *Priv) Get() []byte {
	return p.priv.Serialize()
}

// Set changes the key to a new key if it is valid, or invalidates it
func (p *Priv) Set(b []byte) (P *Priv) {
	if Check(b) && len(b) == 32 {
		p.priv, p.pub = ec.PrivKeyFromBytes(elliptic.P256(), b)
	} else {
		p.Invalidate()
	}
	return p
}

// Size returns the size of the key
func (p *Priv) Size() int {
	if p.invalid {
		return 0
	}
	return 65
}

// Invalidate marks the invalid flag true
func (p *Priv) Invalidate() {
	p.invalid = true
}

// IsValid returns if the key is tagged as valid
func (p *Priv) IsValid() bool {
	return !p.invalid
}

// SetPriv sets the private key of a Priv
func (p *Priv) SetPriv(P *Priv) {
	p.priv = P.priv
	p.pub = P.pub
	p.compressed = P.compressed
	p.invalid = P.invalid
}

// GetPriv returns a copy of private key
func (p *Priv) GetPriv() *Priv {
	return &Priv{p.priv, p.pub, p.compressed, p.invalid}
}

// SetPub sets the public key of a Priv
func (p *Priv) SetPub(P *Pub) {
	p.pub = P.pub
	p.compressed = P.compressed
	p.invalid = P.invalid
}

// GetPub returns a copy of the public key
func (p *Priv) GetPub() (P *Pub) {
	if !p.invalid {
		P = &Pub{p.pub, p.compressed, p.invalid}
	}
	return
}

// Sign a 256 bit hash
func (p *Priv) Sign(hash Uint.U256) (b []byte, err error) {
	if sig, err := p.priv.Sign(hash.ToBytes()); err == nil {
		return sig.Serialize(), err
	}
	return
}

// SignCompact makes a compact signature on a 256 bit hash
func (p *Priv) SignCompact(hash Uint.U256) ([]byte, error) {
	return ec.SignCompact(ec.S256(), p.priv, hash.ToBytes(), p.compressed)
}

// Verify a signature on a hash
func (p *Priv) Verify(hash Uint.U256, S []byte) (key *ec.PublicKey, err error) {
	sig, err := ec.ParseSignature(S, ec.S256())
	if err == nil {
		key, _, err = p.Recover(hash, S)
	}
	if !ecdsa.Verify(key.ToECDSA(), hash.Bytes(), sig.R, sig.S) {
		key = nil
	}
	return
}

// Recover public key from a signature, identify if it was compressed
func (p *Priv) Recover(hash Uint.U256, S []byte) (key *ec.PublicKey, compressed bool, err error) {
	key, compressed, err = ec.RecoverCompact(ec.S256(), S, hash.Bytes())
	return
}
