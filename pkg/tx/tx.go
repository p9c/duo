package tx

import (
	"sync"

	"gitlab.com/parallelcoin/duo/pkg/Uint"
)

// Transaction -
type Transaction struct {
	MinTxFee, MinRelayTxFee int64
	CurrentVersion          int
	Version                 int
	Vin                     []In
	Vout                    []Out
	LockTime                uint
}

// OutPoint is an amount out from an address
type OutPoint struct {
	Hash Uint.U256
	N    uint
}

// InPoint is an amount into an address
type InPoint struct {
	Tx *Transaction
	N  uint
}

// In is a list of InPoints
type In struct {
	PrevOut   OutPoint
	ScriptSig Script
	Sequence  uint
}

// Out is a list of OutPoints
type Out struct {
	Value        int64
	ScriptPubKey Script
}

// OutCompressor controls optimising a transaction
type OutCompressor struct {
	Out *Out
}

// InUndo is an undo buffer for inpoint lists
type InUndo struct {
}

// Undo is an undo buffer for outpoint lists
type Undo struct {
	PrevOut []InUndo
}

// Coins is a coinbase and its collection of movements
type Coins struct {
	Base    bool
	Out     []Out
	Height  int
	Version int
}

// Script is a transaction script
type Script struct{}

// MemPool stores the list of transactions received from the P2P network
type MemPool struct {
	Mutex sync.RWMutex
	Map   map[*Uint.U256]*Transaction
	Next  map[*OutPoint]*InPoint
}

// Orphan is a transaction that has fallen off the head chain
type Orphan struct {
	Tx                 Transaction
	DependsOn          []*Uint.U256
	Priority, FeePerKB float64
}
