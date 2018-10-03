package tx

import (
	"sync"

	"github.com/parallelcointeam/duo/pkg/proto"
	"github.com/parallelcointeam/duo/pkg/wallet/db/entries"
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

// OutPoint is one element in the collection of outputs of a transaction
type OutPoint struct {
	Hash proto.Hash
	N    uint
}

// InPoint is one element in the collection of inputs of a transaction
type InPoint struct {
	Tx *Transaction
	N  uint
}

// In is an input
type In struct {
	PrevOut   OutPoint
	ScriptSig rec.Script
	Sequence  uint
}

// Out is an output
type Out struct {
	Value        int64
	ScriptPubKey rec.Script
}

// OutCompressor controls optimising a transaction
type OutCompressor struct {
	Out *Out
}

// InUndo is
type InUndo struct {
}

// Undo is
type Undo struct {
	PrevOut []InUndo
}

// Coins is
type Coins struct {
	Base    bool
	Out     []Out
	Height  int
	Version int
}

// MemPool stores the list of transactions received from the P2P network
type MemPool struct {
	Mutex sync.RWMutex
	Map   map[proto.Hash]*Transaction
	Next  map[*OutPoint]*InPoint
}

// Orphan is a transaction that is not included in the current canonical chain, but older than the head
type Orphan struct {
	Tx                 Transaction
	DependsOn          []*proto.Hash
	Priority, FeePerKB float64
}
