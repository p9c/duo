package block

/*
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

*/

type Tx struct {
	Version  uint32
	Ins      []TxIn
	Outs     []TxOut
	Locktime uint32
}

type TxIn struct {
	PrevTxHash     []byte
	PrevTxoutIndex int32
	Script         []byte
	Sequence       uint32
}

type TxOut struct {
	Value  uint64
	Script []byte
}

type Raw struct {
	Version        uint32
	HashPrevBlock  []byte
	HashMerkleRoot []byte
	Time           uint32
	Bits           []byte
	Nonce          uint32
	Transactions   []Tx
}
