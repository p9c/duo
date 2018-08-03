package wallet

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/tx"
	"sync"
)

const (
	// FeatureBase is the base version number for a wallet
	FeatureBase = 10500
	// FeatureWalletCrypt indicates if the wallet enables encrypted keys
	FeatureWalletCrypt = 40000
	// FeatureCompressedPubKey indicates if the wallet enables compressed public keys
	FeatureCompressedPubKey = 60000
	// FeatureLatest is the newest version of the wallet
	FeatureLatest = 60000
)

// Wallet controls access to a wallet.db file containing keys and data relating to accounts and addresses
type Wallet struct {
	key.StoreCrypto
	dBEncryption        *DB
	version, maxVersion int
	Mutex               sync.RWMutex
	FileBacked          bool
	File                string
	KeyPoolSet          []int64
	KeyMetadataMap      map[*key.ID]*KeyMetadata
	MasterKeysMap       MasterKeyMap
	MasterKeyMaxID      uint
	WalletMap           map[*Uint.U256]Tx
	OrderPosNext        int64
	RequestCountMap     map[*Uint.U256]int
	AddressBookMap      map[key.TxDestination]string
	DefaultKey          key.Pub
	LockedCoinsSet      []tx.OutPoint
	TimeFirstKey        int64
}

// New returns a new Wallet
func New() *Wallet {
	return &Wallet{
		version:        FeatureBase,
		maxVersion:     FeatureBase,
		FileBacked:     false,
		MasterKeyMaxID: 0,
		OrderPosNext:   0,
	}
}

// NewFromFile makes a new wallet based on a wallet.dat file
func NewFromFile(filename string) *Wallet {
	return &Wallet{
		version:        FeatureBase,
		maxVersion:     FeatureBase,
		File:           filename,
		FileBacked:     true,
		MasterKeyMaxID: 0,
		OrderPosNext:   0,
	}
}

// TxPair is a transaction and an accounting entry
type TxPair map[*Tx]*AccountingEntry

// TxItems is an array of a map of transactions
type TxItems map[int64]TxPair

type wallet interface {
	CanSupportFeature(int) bool
	AvailableCoins([]Output, bool)
	SelectCoinsMinConf(int64, int, int, []Output) error
	IsLockedCoin(*Uint.U256, uint) bool
	LockCoin(*tx.OutPoint)
	UnlockCoin(*tx.OutPoint)
	UnlockAllCoins()
	ListLockedCoins([]tx.OutPoint)
	GenerateNewKey() *key.Pub
	AddKeyPair(*key.Priv, *key.Pub) bool
	LoadKey(*key.Priv, *key.Pub) bool
	LoadKeyMetadata(*key.Pub, *KeyMetadata) bool
	LoadMinVersion(int) bool
	AddCryptedKey(*key.Pub, *KeyMetadata) bool
	LoadCryptedKey(*key.Pub, []byte) bool
	AddScript(*key.Script) bool
	LoadScript(*key.Script) bool
	Unlock(string) bool
	ChangeWalletPassphrase(string, string) bool
	EncryptWallet(string)
	GetKeyBirthTimes(map[*key.ID]int64)
	IncOrderPosNext(*DB) int64
	OrderedTxItems([]AccountingEntry, string) *TxItems
	MarkDirty()
	AddToWallet(Tx) bool
	AddToWalletIfInvolvingMe(*Uint.U256, *tx.Transaction, *block.Block, bool, bool) bool
	EraseFromWallet(*Uint.U256) bool
	WalletUpdateSpent(*tx.Transaction)
	ScanForWalletTransactions(*block.Index, bool) int
	ReacceptWalletTransactions()
	ResendWalletTransactions()
	GetBalance() int64
	GetUnconfirmedBalance() int64
	GetImmatureBalance() int64
	CreateTransactions([]map[*key.Script]int64, *Tx, *ReserveKey, int64, string) bool
	CreateTransaction(*key.Script, int64, *Tx, *ReserveKey, int64, string) bool
	CommitTransaction(*Tx, *ReserveKey) bool
	SendMoney(*key.Script, int64, *Tx, bool) string
	SendMoneyToDestination(*key.TxDestination) string
	NewKeyPool() bool
	TopUpKeyPool() bool
	AddReserveKey(*KeyPool) int64
	ReserveKeyFromKeyPool(int64, *KeyPool)
	KeepKey(int64)
	ReturnKey(int64)
	GetKeyFromPool(*key.Pub, bool) bool
	GetOldestKeyPoolTime() int64
	GetAllReserveKeys() []key.ID
	GetAddressGroupings() []key.TxDestination
	GetAddressBalances() map[*key.TxDestination]int64
	IsMyTxIn(*tx.In) bool
	GetDebit(*tx.In) int64
	IsMyTxOut(*tx.Out) bool
	GetCredit(*tx.Out) int64
	IsChange(*tx.Out) bool
	GetChange(*tx.Out) int64
	IsMyTX(*tx.Transaction) bool
	IsFromMe(*tx.Transaction) bool
	GetTxDebit(*tx.Transaction) int64
	GetTxCredit(*tx.Transaction) int64
	GetTxChange(*tx.Transaction) int64
	SetBestChaing(*block.Locator)
	LoadWallet(bool) DBErrors
	SetAddressBookName(*key.TxDestination, string) bool
	DelAddressBookName(*key.TxDestination) bool
	UpdatedTransaction(*Uint.U256)
	PrintWallet(*block.Block)
	Inventory(*Uint.U256)
	GetKeyPoolSize() int
	GetTransaction(*Uint.U256, *Tx) bool
	SetDefaultKey(*key.Pub) bool
	SetMinVersion(int, *DB, bool) bool
	SetMaxVersion(int) bool
	GetVersion() int
	NotifyAddressBookChanged(*Wallet, *key.TxDestination, string, bool, int)
	NotifyTransactionChanged(*Wallet, *Uint.U256, int)
}
