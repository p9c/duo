package state

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/block"
)

// DiskBlockPos is the position of a cursor on a file
type DiskBlockPos struct {
	File int
	Pos  uint
}

// DiskTxPos is the sequence of a transaction in a file
type DiskTxPos struct {
	DiskBlockPos
	TxOffset uint
}

// BlockIndexWorkComparator -
type BlockIndexWorkComparator struct{}

// CoinStats is statistics of the blockchain
type CoinStats struct {
	Height                                                        int
	HashBlock                                                     Uint.U256
	Transactions, TransactionOutputs, SerializedSize, TotalAmount uint64
	HashSerialized                                                Uint.U256
}

// BlockTemplate is the template for a block
type BlockTemplate struct {
	Block            block.Block
	TxFees, TxSigOps []int64
}
