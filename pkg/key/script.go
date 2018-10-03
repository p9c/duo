package key

import (
	"sync"

	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
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

// SigData -
type SigData struct {
	Hash   core.Hash
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
	Keys     []core.Address
}

// ScriptVisitor -
type ScriptVisitor struct {
	script *rec.Script
}
