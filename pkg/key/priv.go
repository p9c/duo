package key
import (
	"github.com/awnumar/memguard"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/anaskhan96/base58check"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/ec"
)
// In memory store for a private key (includes the public key)
type Priv struct {
	*Pub
	priv   []byte
	cipher *memguard.LockedBuffer
}
// Creates a new private key from random bytes
func NewPriv(cipher ...*memguard.LockedBuffer) (p Priv, err error) {
	var k *memguard.LockedBuffer
	k, err = memguard.NewMutableRandom(32)
	priv, pub := ec.PrivKeyFromBytes(ec.S256(), k.Buffer())
	p.SetPub(*pub)
	p.Pub.Compress()
	pr, _ := memguard.NewMutableFromBytes(priv.Serialize())
	p.Set(pr)
	return
}
// Creates a new private key from a raw secret in a LockedBuffer
func NewPrivFromBytes(key *memguard.LockedBuffer, cipher ...*memguard.LockedBuffer) (p *Priv) {
	p.cipher = cipher[0]
	p.Set(key)
	return
}
type priv interface {
	Clear()
	Bytes() memguard.LockedBuffer
	Size() int
	Invalidate()
	IsValid() bool
	Set(*Priv) bool
	GetPub() *Pub
	Sign(Uint.U256, []byte) bool
	SignCompact(Uint.U256, []byte) bool
	Verify(hash Uint.U256, S []byte)
	ToBase58Check() string
}
// Zeroes out all of the contents of the priv and detaches from the unlock cipher
func (p *Priv) Clear() {
	for i := range p.priv {
		p.priv[i] = 0
	}
	p.pub.X.SetBytes(make([]byte, 32))
	p.pub.Y.SetBytes(make([]byte, 32))
	p.cipher = nil
}
// Returns the full private key as a byte slice
func (p *Priv) Bytes() (B memguard.LockedBuffer) {
	decrypted, _ := memguard.NewMutableFromBytes(p.priv)
	if p.cipher != nil {
	}
	B.Copy(decrypted.Buffer())
	return 
}
// Returns the size of the key
func (p *Priv) Size() int {
	if p.invalid {
		return 0
	}
	return 65
}
// Marks the invalid flag true and wipes all the data
func (p *Priv) Invalidate() {
	p.Clear()
	p.invalid = true
}
// Returns true if key is valid
func (p *Priv) IsValid() bool {
	return !p.invalid
}
// Sets the private key of a Priv from the unencrypted bytes
func (p *Priv) Set(b *memguard.LockedBuffer) {
	// TODO: encrypt key and generate Pub from unencrypted source
	encrypted := b.Buffer()
	p.priv = encrypted
}
// GetPub returns a copy of the public key
func (p *Priv) GetPub() (P *Pub) {
	if !p.invalid {
		P = new(Pub)
		P.SetPub(p.pub)
	}
	return
}
// Sign a 256 bit hash
func (p *Priv) Sign(hash Uint.U256) (b []byte, err error) {
	// TODO: need to decrypt private key before signing
	// if sig, err := p.priv.Sign(hash.ToBytes()); err == nil {
	// 	return sig.Serialize(), err
	// }
	return
}
// SignCompact makes a compact signature on a 256 bit hash
func (p *Priv) SignCompact(hash Uint.U256) (sig Uint.U256, err error) {
	return //ec.SignCompact(ec.S256(), p.priv, hash.ToBytes(), p.compressed)
}
// Verify a signature on a hash
func (p *Priv) Verify(hash Uint.U256, S []byte) (key *Pub, err error) {
	var sig *ec.Signature
	if sig, err = ec.ParseSignature(S, ec.S256()); err != nil {
		return
	} else {
		var keyEC *ec.PublicKey
		keyEC, _, err = p.Recover(hash, S)
		if !ecdsa.Verify(keyEC.ToECDSA(), hash.Bytes(), sig.R, sig.S) {
			key = nil
		}
		key.SetPub(*keyEC)
		return
	}
}
// Recover public key from a signature, identify if it was compressed
func (p *Priv) Recover(hash Uint.U256, S []byte) (key *ec.PublicKey, compressed bool, err error) {
	key, compressed, err = ec.RecoverCompact(ec.S256(), S, hash.Bytes())
	return
}
// ToBase58Check returns a private key encoded in base58check with the network specified prefix
func (p *Priv) ToBase58Check(net string) (b58 string, err error) {
	// TODO: decrypt key first
	decrypted := p.priv
	h := hex.EncodeToString(decrypted)
	b58, err = base58check.Encode(B58prefixes[net]["privkey"], h)
	return 
}
