package key

import (
	"sync"

	"github.com/parallelcointeam/duo/pkg/Uint"
)

// Various script flags
const (
	// SigHashAll is
	SigHashAll = 1
	// SigHashNone is
	SigHashNone = 2
	// SigHashSingle is
	SigHashSingle = 3
	// SigHashAnyoneCanPay is
	SigHashAnyoneCanPay = 0x80
	// VerifyNone is
	VerifyNone = 0
	// VerifyP2SH is
	VerifyP2SH = 1
	// VerifyStrictEnc is
	VerifyStrictEnc = 1 << 1
	// VerifyNocache is
	VerifyNocache = 1 << 2
	// TxNonstandard is a nonstandard transaction (eg just writing bytes)
	TxNonstandard = iota
	// TxPubKey is is a public key
	TxPubKey
	// TxPubKeyHash is the hash of a public key
	TxPubKeyHash
	// TxScriptHash is the hash of a script
	TxScriptHash
	// TxMultisig is a multisignature transaction
	TxMultisig
)

type ScriptID ID
type Script []byte
type script interface {
	GetID() *ScriptID
}

// GetID computes the RIPEMD160 hash of the script
func (s *Script) GetID() (sid *ScriptID) {
	sid.U160 = *Uint.RIPEMD160([]byte(*s))
	return
}

// ScriptCompressor -
type ScriptCompressor struct {
	specialScripts uint // 6 defined
	script         Script
}

// SigData -
type SigData struct {
	Hash   Uint.U256
	Sig    []byte
	Pubkey Pub
}

// SignatureCache -
type SignatureCache struct {
	ValidSet []bool
	Mutex    sync.RWMutex
}

// StoreIsMineVisitor -
type StoreIsMineVisitor struct {
	KeyStore *Store
}

// AffectedKeysVisitor -
type AffectedKeysVisitor struct {
	KeyStore *Store
	Keys     []ID
}

// ScriptVisitor -
type ScriptVisitor struct {
	script *Script
}
