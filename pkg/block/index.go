// Package block is a library for working with blocks of the Parallelcoin DUO token ledger
package block

import (
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/tx"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

const (
	// ModeValid means a block is valid
	ModeValid = iota
	// ModeInvalid means a block is invalid
	ModeInvalid
	// ModeError means there was an error in the block
	ModeError
)

var (
	// ChainIndex is a centralised index of a server's chain
	ChainIndex Index
)

// Contents is a
type Contents struct {
	Txs        []tx.Transaction
	MerkleTree []core.Hash
}

// Header is the data contained in a block header
type Header struct {
	CurrentVersion, Version       int
	HashPrevBlock, HashMerkleRoot core.Hash
	Time, Bits, Nonce             uint
	Contents
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

// DiskIndex is the on-disk index
type DiskIndex struct {
	Index
	HashPrev core.Hash
}

// ValidationState stores the state of validation of a block
type ValidationState struct {
	Mode int
	DoS  int
}

// Locator allows you to quickly find a block
type Locator struct {
	Have []core.Hash
}

// CoinStats stores the current state of the currency
type CoinStats struct {
	Height                         int
	HashBlock                      core.Hash
	Txs, TxOutputs, SerializedSize uint64
	HashSerialized                 core.Hash
	TotalAmount                    uint64
}

// CoinsView is
type CoinsView struct{}

// CoinsViewBacked is
type CoinsViewBacked struct {
	base *CoinsView
}

// CoinsViewCache is
type CoinsViewCache struct {
	IndexTip   Index
	CacheCoins map[*core.Hash]tx.Coins
}

// CoinsViewMemPool is
type CoinsViewMemPool struct {
	MemPool tx.MemPool
}

// Block is the data of a block
type Block struct {
	Transactions []tx.Transaction
	MerkleTree   []*core.Hash
}

// Template is a block template
type Template struct {
	Block          Block
	TxFees, SigOps []int64
}

// Undo is a collection of tx undos
type Undo struct {
	Txs tx.Undo
}

// ScriptCheck tracks the verification of scripts
type ScriptCheck struct {
	ScriptPubKey        rec.Script
	To                  tx.Transaction
	In, Flags, HashType uint
}

// MerkleTx is
type MerkleTx struct {
	tx.Transaction
	HashBlock    core.Hash
	MerkleBranch []*core.Hash
	Index        int
	Verified     bool
}

// PartialMerkleTree is
type PartialMerkleTree struct {
	Transactions uint
	Bits         []bool
	Hash         []core.Hash
	Bad          bool
}

// FileInfo stores details about the blockchain on disk
type FileInfo struct {
	Blocks, Size, UndoSize, HeightFirst, HeightLast uint
	TimeFirst, TimeLast                             uint64
}

// MerkleBlock is
type MerkleBlock struct {
	Header Header
}
