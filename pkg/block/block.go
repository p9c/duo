package block

import (
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/tx"
)

// Block is the data stored in a block
type Block struct {
	Header       *[]byte
	Transactions []*tx.Transaction
	MerkleTree   []*core.Hash
}

// Index is an index of blocks
type Index struct {
	HashBlock                     core.Hash
	Prev                          *Index
	Height, File                  int
	DataPos, UndoPos              uint
	ChainWork                     core.Hash
	TxCount, TxCumulative, Status uint
	Version                       int
	HashMerkleRoot                core.Hash
	Time, Bits, Nonce             uint
}

// DiskIndex is the on-disk augmented version of the Index
type DiskIndex struct {
	Index
	HashPrev core.Hash
}

// Locator allows you to quickly find a block
type Locator struct {
	Have []*core.Hash
}
