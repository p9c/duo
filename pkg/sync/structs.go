package sync

import (
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/rpc"

	"github.com/dgraph-io/badger"
)

// Node is a sync client that updates by polling from full node, but will eventually implement as a standalone p2p client
//
// Records structures
//
// In the on-disk form of these records several fields are abbreviated by zero un-padding, height appears only alone as key and value, and trailing zeroes are removed, only 24 bits are used for length (16mb max block size) and the starting position is 64 bits stored with a byte length prefix and trailing zeroes removed, and lastly, the hash of blocks by design contains leading zeroes that indicate the difficulty, these are dispensed with and fill the last part of the value for a Block.
//
// For decoding these abbreviated storage formats, the proper full length is known and the bytes are padded first to restore orignial format and then converted into the format specified, block hash has its prefix zeroes readded, which are required to generate the correct hhash64
type Node struct {
	RPC        *rpc.Client
	DB         *badger.DB
	Latest     uint32
	LatestHash []byte
	Best       uint32
	BestTime   int64
	core.State
}

// Block links height to hash. This record type is identified by a prefix 1 byte
type Block struct {
	// key
	Height uint32
	// value
	Hash []byte
}

// Hash links the block hash to the height. This enables reverse lookup from hash to height. In the database the height value is pruned of its' trailing zeros to save space.
type Hash struct {
	// key
	HHash []byte
	// value
	Height uint32
}

// Address is a record that notes the places an address appears by block and transaction.
//
// This is a necessary index for searching for information about addresses, especially calculating their balance at a given height.
type Address struct {
	// key
	//     HighwayHash 64 bit hash of 160 bit address
	HHash []byte
	// value
	//     This is stored as a single value, 4 bytes for height, 2 bytes to specify the position of the address' appearance in the array of transactions.
	Locations []Location
}

// Location is a specification for a transaction in a block. 32 bits encodes up to 4 billion (more than enough) blocks and 16 bits for transaction number allows addressing up to 65536 transactions in the array inside a block, which is precisely the information required to gather the inputs and outputs and compute a balance.
type Location struct {
	Height uint32
	TxNum  uint16
}

// BalanceCache is a result cache that stores the results of previous queries of balances of an address with the contemporary best block height so subsequent queries don't have to make as many RPC queries to get the answer.
//
// We aren't storing the tx data, just making it much faster to find it.
type BalanceCache struct {
	// key
	//     Highway Hash 64 of 160 bit address
	HHash []byte

	// value
	// Balance is prefix-length-trailing-zero-trimmed
	Balance uint64
	// Height is trailing zero trimmed
	Height uint32
}
