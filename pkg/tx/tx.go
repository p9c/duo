package tx

import (
	"github.com/parallelcointeam/duo/pkg/Uint"
	"sync"
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
type OutPoint struct {
	Hash Uint.U256
	N    uint
}
type InPoint struct {
	Tx *Transaction
	N  uint
}
type In struct {
	PrevOut   OutPoint
	ScriptSig Script
	Sequence  uint
}
type Out struct {
	Value        int64
	ScriptPubKey Script
}

// OutCompressor controls optimising a transaction
type OutCompressor struct {
	Out *Out
}
type InUndo struct {
}
type Undo struct {
	PrevOut []InUndo
}
type Coins struct {
	Base    bool
	Out     []Out
	Height  int
	Version int
}
type Script struct{}

// MemPool stores the list of transactions received from the P2P network
type MemPool struct {
	Mutex sync.RWMutex
	Map   map[*Uint.U256]*Transaction
	Next  map[*OutPoint]*InPoint
}
type Orphan struct {
	Tx                 Transaction
	DependsOn          []*Uint.U256
	Priority, FeePerKB float64
}
