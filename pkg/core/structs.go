package core

// State is a base status/error object, that can be embedded in any other type to proide status notifications
type State struct {
	Status string
	err    error
}

// Address is used as the key for searching for public keys (addresses also), scripts, transactions and blocks, generated using the hash160 function, which is a sha256 followed by ripemd160.
type Address string

// Hash is a 256 byte hash stored as a string for use with maps, used for block hashes, transaction hashes, message hashes, and other things
type Hash string

// MerkleTx is the merkle tree hash data for a transaction
type MerkleTx struct {
	HashBlock      []byte
	MerkleBranch   []byte
	Index          int64
	merkleVerified bool
}
