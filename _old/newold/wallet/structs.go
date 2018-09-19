package wallet

import (
	"github.com/parallelcointeam/duo/pkg/Uint"
	"github.com/parallelcointeam/duo/pkg/block"
	"github.com/parallelcointeam/duo/pkg/key"
)

type MasterKeyMap map[uint]*KeyMetadata
type KeyPool struct {
	Time   int64
	PubKey key.Pub
}
type ReserveKey struct {
	wallet *Wallet
	Index  int64
	PubKey key.Pub
}
type ValueMap map[string]string

// Orders are a key value pair
type Orders struct {
	Key, Value string
}
type Transaction struct {
	block.MerkleTx
	wallet                                        *Wallet
	Prev                                          []block.MerkleTx
	OrderForm                                     []Orders
	TimeReceivedIsTxTime, TimeReceived, TimeSmart uint
	FromMe                                        byte
	FromAccount                                   string
	Spent                                         []byte
	OrderPos                                      int64
	CachedFlags                                   struct {
		Debit, Credit, ImmatureCredit, AvailableCredit, Change bool
	}
	CachedValues struct {
		Debit, Credit, ImmatureCredit, AvailableCredit, Change int64
	}
}
type TxOutput struct {
	Tx       Tx
	I, Depth int
}
type AccountingEntry struct {
	Account               string
	CreditDebit, Time     int64
	OtherAccount, Comment string
	ValueMap              ValueMap
	OrderPos              int64
	EntryNo               uint64
}
type ScanState struct {
	Keys, CKeys, KeyMeta      uint
	IsEncrypted, AnyUnordered bool
	FileVersion               int
	WalletUpgrade             []*Uint.U256
}
type TxPair map[*Tx]*AccountingEntry
type TxItems map[int64]TxPair
