package wallet
import (
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
)
type KeyMetadata struct {
	Version    uint32
	CreateTime int64
}
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
type Tx struct {
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
type Output struct {
	Tx       Tx
	I, Depth int
}
type Account struct {
	PubKey key.Pub
}
type AccountingEntry struct {
	Account               string
	CreditDebit, Time     int64
	OtherAccount, Comment string
	ValueMap              ValueMap
	OrderPos              int64
	EntryNo               uint64
}
type Key struct {
	PrivKey                  *key.Priv
	TimeCreated, TimeExpires int64
	Comment                  string
}
type ScanState struct {
	Keys, CKeys, KeyMeta      uint
	IsEncrypted, AnyUnordered bool
	FileVersion               int
	WalletUpgrade             []*Uint.U256
}
type TxPair map[*Tx]*AccountingEntry
type TxItems map[int64]TxPair
