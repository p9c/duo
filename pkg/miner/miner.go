package miner

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/tx"
)

var (
	// SHA256InitState is the init state for SHA256
	SHA256InitState = []uint{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
	}
	// LastBlockTx is
	LastBlockTx uint64
	// LastBlockSize is the size of the last block
	LastBlockSize uint64
)

// TxPriority stores priority numbers relating to each transaction
type TxPriority struct {
	A, B float64
	Tx   *tx.Transaction
}

// PriorityCompare is a structure for ranking transactions by their fee value
type PriorityCompare struct {
	// Whether we will compare by fee or not
	ByFee bool
}

// Orphan is a transaction that does not link to one in the consensus chain
type Orphan struct {
	// Tx is the transaction data
	Tx tx.Transaction
	// DependsOn is the transaction that this transaction depends on
	DependsOn []*Uint.U256
	// Priority is the priority of this orphan transaction
	Priority float64
	// FeePerKb is the value of this Tx
	FeePerKb float64
}
