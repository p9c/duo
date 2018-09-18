package key

import (
	"github.com/anaskhan96/base58check"
	"github.com/parallelcointeam/duo/pkg/Uint"
	"github.com/parallelcointeam/duo/pkg/ec"
)

type Pub struct {
	pub                 ec.PublicKey
	compressed, invalid bool
}
type pub interface {
	SetPub(*ec.PublicKey)
	GetPub() *Pub
	Size() int
	Key() []byte
	GetID() ID
	GetHash() Uint.U256
	Invalidate()
	IsValid() bool
	Compress() bool
	Decompress() bool
	IsCompressed() bool
	Check([]byte) bool
	MakeNewKey(bool) *Pub
	ToBase58Check(string) string
}

// SetPub sets the public key of a Priv
func (p *Pub) SetPub(P ec.PublicKey) {
	p.pub = P
}

// GetPub returns the key stored in a Priv
func (p *Pub) GetPub() (P *ec.PublicKey) {
	return &p.pub
}

// Size returns the size of the key if it were requested as bytes
func (p *Pub) Size() int {
	if p.compressed {
		return 33
	}
	return 65
}

// Key returns the raw bytes according to the compression setting
func (p *Pub) Key() []byte {
	if p.compressed {
		return p.pub.SerializeCompressed()
	}
	return p.pub.SerializeUncompressed()
}

// GetID returns the RIPEMD160 hash of the 256 byte public key
func (p *Pub) GetID() (id *ID) {
	id.U160 = *Uint.RIPEMD160(p.pub.SerializeUncompressed())
	return
}

// GetHash returns the SHA256 hash of the public key
func (p *Pub) GetHash() (u Uint.U256) {
	return *Uint.SHA256(p.pub.SerializeUncompressed())
}

// Invalidate marks the invalid flag true
func (p *Pub) Invalidate() {
	p.invalid = true
}

// IsValid returns whether the key has been invalidated or not
func (p *Pub) IsValid() bool {
	return !p.invalid
}

// Decompress doesn't really decompress, it just flags that outputs will be not compressed if they were before
func (p *Pub) Decompress() bool {
	if p.compressed {
		p.compressed = false
		return true
	}
	return false
}

// Compress doesn't really compress, it just flags that outputs will be compressed if it wasn't before
func (p *Pub) Compress() bool {
	if !p.compressed {
		p.compressed = true
		return true
	}
	return false
}

// IsCompressed returns whether the serializing commands return compressed format
func (p *Pub) IsCompressed() bool {
	return p.compressed
}

// ToBase58Check returns a private key encoded in base58check with the network specified prefix
func (p *Pub) ToBase58Check(net string) string {
	b58, err := base58check.Encode(B58prefixes[net]["pubkey"], string(p.Key()))
	if err != nil {
		return "Base58check encoding failure"
	}
	return b58
}
