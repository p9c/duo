package wallet

import (
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
)

const (
	CurrentVersion          = 1
)

var (
	AccountingEntryNumber = 0
)

// MasterKeyMap is a list of master keys
type MasterKeyMap map[uint]*KeyMetadata

// KeyPool is a list element of a list of keys
type KeyPool struct {
	Time   int64
	PubKey key.Pub
}

// ReserveKey is a reserve key type
type ReserveKey struct {
	wallet *Wallet
	Index  int64
	PubKey key.Pub
}

// ValueMap is map of values
type ValueMap map[string]string

// Orders are a key value pair
type Orders struct {
	Key, Value string
}

// Tx is a transaction
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

// Output is a transaction output
type Output struct {
	Tx       Tx
	I, Depth int
}

// Account is a public key
type Account struct {
	PubKey key.Pub
}

// AccountingEntry is a record of a transaction
type AccountingEntry struct {
	Account               string
	CreditDebit, Time     int64
	OtherAccount, Comment string
	ValueMap              ValueMap
	OrderPos              int64
	EntryNo               uint64
}

// Key is a private key with extra data
type Key struct {
	PrivKey                  *key.Priv
	TimeCreated, TimeExpires int64
	Comment                  string
}
