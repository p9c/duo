// A library for working with blocks of the Parallelcoin DUO token ledger
package block
import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/tx"
)
var (
	// ChainIndex is a centralised index of a server's chain
	ChainIndex Index
)
type Contents struct {
	Txs        []tx.Transaction
	MerkleTree []Uint.U256
}
type Header struct {
	CurrentVersion, Version       int
	HashPrevBlock, HashMerkleRoot Uint.U256
	Time, Bits, Nonce             uint
	Contents
}
type Index struct {
	HashBlock                     Uint.U256
	Prev                          *Index
	Height, File                  int
	DataPos, UndoPos              uint
	ChainWork                     Uint.U256
	TxCount, TxCumulative, Status uint
	Version                       int
	HashMerkleRoot                Uint.U256
	Time, Bits, Nonce             uint
}
type DiskIndex struct {
	Index
	HashPrev Uint.U256
}
const (
	// ModeValid means a block is valid
	ModeValid = iota
	// ModeInvalid means a block is invalid
	ModeInvalid
	// ModeError means there was an error in the block
	ModeError
)
// ValidationState stores the state of validation of a block
type ValidationState struct {
	Mode int
	DoS  int
}
// Locator allows you to quickly find a block
type Locator struct {
	Have []Uint.U256
}
// CoinStats stores the current state of the currency
type CoinStats struct {
	Height                         int
	HashBlock                      Uint.U256
	Txs, TxOutputs, SerializedSize uint64
	HashSerialized                 Uint.U256
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
	CacheCoins map[*Uint.U256]tx.Coins
}
// CoinsViewMemPool is
type CoinsViewMemPool struct {
	MemPool tx.MemPool
}
type Block struct {
	Transactions []tx.Transaction
	MerkleTree   []*Uint.U256
}
type Template struct {
	Block          Block
	TxFees, SigOps []int64
}
type Undo struct {
	Txs tx.Undo
}
// ScriptCheck tracks the verification of scripts
type ScriptCheck struct {
	ScriptPubKey        tx.Script
	To                  tx.Transaction
	In, Flags, HashType uint
}
// MerkleTx is
type MerkleTx struct {
	tx.Transaction
	HashBlock    Uint.U256
	MerkleBranch []*Uint.U256
	Index        int
	Verified     bool
}
// PartialMerkleTree is
type PartialMerkleTree struct {
	Transactions uint
	Bits         []bool
	Hash         []Uint.U256
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
