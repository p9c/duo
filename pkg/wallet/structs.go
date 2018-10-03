package wallet

import (
	"time"

	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/tx"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// KeyPool is a collection of available addresses for constructing transactions
type KeyPool map[int64]*rec.Pool

// KeyMetadata is
type KeyMetadata map[core.Address]*KeyMetadata

// Transactions is a map of transactions in the wallet
type Transactions map[core.Hash]*rec.Tx

// AddressBook is a collection of correspondent addresses
type AddressBook map[core.Address]key.Account

// Wallet controls access to a wallet.db file containing keys and data relating to accounts and addresses
type Wallet struct {
	KeyStore            key.Store
	DB                  *db.DB
	version, maxVersion int
	FileBacked          bool
	File                string
	KeyPool             KeyPool
	KeyPoolHigh         int64
	KeyPoolLow          int64
	KeyPoolLifespan     time.Duration
	KeyMetadata         KeyMetadata
	MasterKeys          key.MasterKeys
	Transactions        Transactions
	OrderPosNext        int64
	RequestCountMap     map[core.Hash]int
	AddressBook         AddressBook
	DefaultKey          *key.Pub
	LockedCoinsSet      []*tx.OutPoint
	TimeFirstKey        int64
	core.State
}

// ReserveKey is
type ReserveKey struct {
	wallet *Wallet
	Index  int64
	PubKey key.Pub
}

// MasterKeyMap is the collection of masterkeys in the wallet
type MasterKeyMap map[uint64]*KeyMetadata

// ValueMap is
type ValueMap map[string]string

// Orders are a key value pair
type Orders struct {
	Key, Value string
}

// ScanState is
type ScanState struct {
	Keys, CKeys, KeyMeta      uint
	IsEncrypted, AnyUnordered bool
	FileVersion               int
	WalletUpgrade             []*core.Hash
}
