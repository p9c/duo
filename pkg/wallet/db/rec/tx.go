package rec

import (
	"sync"

	"github.com/parallelcointeam/duo/pkg/core"
)

// Transaction -
type Transaction struct {
	MinTxFee, MinRelayTxFee int64
	CurrentVersion          int
	Version                 int
	Vin                     []TxIn
	Vout                    []TxOut
	LockTime                uint
}

// TxOutput is a transaction output
type TxOutput struct {
	Tx       *Tx
	I, Depth int
}

// TxOutPoint is one element in the collection of outputs of a transaction
type TxOutPoint struct {
	Hash core.Hash
	N    uint
}

// TxInPoint is one element in the collection of inputs of a transaction
type TxInPoint struct {
	Tx *Transaction
	N  uint
}

// TxIn is an input
type TxIn struct {
	PrevOut   TxOutPoint
	ScriptSig Script
	Sequence  uint
}

// TxOut is an output
type TxOut struct {
	Value        int64
	ScriptPubKey Script
}

// TxOutCompressor controls optimising a transaction
type TxOutCompressor struct {
	TxOut *TxOut
}

// TxInUndo is
type TxInUndo struct {
}

// TxUndo is
type TxUndo struct {
	PrevOut []TxInUndo
}

// Coins is
type Coins struct {
	Base    bool
	TxOut   []TxOut
	Height  int
	Version int
}

// MemPool stores the list of transactions received from the P2P network
type MemPool struct {
	Mutex sync.RWMutex
	Map   map[core.Hash]*Transaction
	Next  map[*TxOutPoint]*TxInPoint
}

// Orphan is a transaction that is not included in the current canonical chain, but older than the head
type Orphan struct {
	Tx                 Transaction
	DependsOn          []*core.Hash
	Priority, FeePerKB float64
}

// Pair is
type Pair map[*Transaction]*Accounting

// Items are
type Items map[int64]Pair

// TxDestination can be multiple types
type TxDestination interface{}

// NoDestination is nowhere
type NoDestination struct{}

// ScriptCompressor -
type ScriptCompressor struct {
	specialScripts uint // 6 defined
	script         Script
}
