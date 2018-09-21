package key

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/hash160"
)

// Sig is a bitcoin EC signature
type Sig struct {
	buf.Byte
	mh *buf.Byte
}

// NewSig creates a new signature
func NewSig() (out *Sig) {
	return new(Sig)
}

// AsEC returns the signature in ecdsa format
func (r *Sig) AsEC() (out *btcec.Signature) {
	if r == nil {
		r = NewSig()
		r.SetStatus(er.NilRec)
		return &btcec.Signature{}
	}
	sig, err := btcec.ParseSignature(*r.Bytes(), btcec.S256())
	if r.SetStatusIf(err); err != nil {
		r.SetStatus("compact sig has no EC format")
		return new(btcec.Signature)
	}
	return sig
}

// Recover returns a public key with a buffer containing the public key if found and a status indicating if it was successful
func (r *Sig) Recover(h *[]byte, addr *[]byte) (out *Pub) {
	if r == nil {
		r = NewSig()
		r.SetStatus(er.NilRec)
		out = NewPub()
		out.SetStatus(er.NilRec)
		return
	}
	pub, _, err := btcec.RecoverCompact(btcec.S256(), *r.Bytes(), *h)
	if pub != nil {
		out = NewPub()
		p := pub.SerializeCompressed()
		out.Copy(&p)
		return
	}
	if r.SetStatusIf(err); err != nil {
		var btcsig []byte
		btcsig = (*r.Bytes())[4:37]
		var prefix byte
		if (*r.Bytes())[1] == 69 {
			btcsig = append(btcsig[1:], (*r.Bytes())[39:]...)
		} else {
			btcsig = append(btcsig[:len(btcsig)-1], (*r.Bytes())[38:]...)
		}
		btcsig = append([]byte{0}, btcsig...)
		for prefix = 27; prefix < 35; prefix++ {
			btcsig[0] = prefix
			var comp bool
			pub, comp, err = btcec.RecoverCompact(btcec.S256(), btcsig, *h)
			if pub != nil {
				var p []byte
				out = NewPub()
				if comp {
					p = pub.SerializeCompressed()
				} else {
					p = pub.SerializeUncompressed()
				}
				ar := buf.NewByte().Copy(hash160.Sum(&p))
				if ar.IsEqual(addr) {
					out.Copy(&p)
					return
				}
			}
		}
	}
	return
}
