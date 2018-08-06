package state
import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/block"
)
type DiskBlockPos struct {
	File int
	Pos  uint
}
type DiskTxPos struct {
	DiskBlockPos
	TxOffset uint
}
// BlockIndexWorkComparator -
type BlockIndexWorkComparator struct{}
type CoinStats struct {
	Height                                                        int
	HashBlock                                                     Uint.U256
	Transactions, TransactionOutputs, SerializedSize, TotalAmount uint64
	HashSerialized                                                Uint.U256
}
type BlockTemplate struct {
	Block            block.Block
	TxFees, TxSigOps []int64
}
