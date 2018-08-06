package wallet
import (
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/tx"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
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
	DB        *DB
	version, maxVersion int
	FileBacked          bool
	File                string
	KeyPoolSet          []int64
	KeyMetadataMap      map[*key.ID]*KeyMetadata
	MasterKeysMap       MasterKeyMap
	MasterKeyMaxID      uint
	WalletMap           map[*Uint.U256]*Tx
	OrderPosNext        int64
	RequestCountMap     map[*Uint.U256]int
	AddressBookMap      map[*key.TxDestination]string
	DefaultKey          *key.Pub
	LockedCoinsSet      []*tx.OutPoint
	TimeFirstKey        int64
}
type wallet interface {
	AddCryptedKey(*key.Pub, *KeyMetadata) bool
	AddKeyPair(*key.Priv, *key.Pub) bool
	AddReserveKey(*KeyPool) int64
	AddScript(*key.Script) bool
	AddToWallet(Tx) bool
	AddToWalletIfInvolvingMe(*Uint.U256, *tx.Transaction, *block.Block, bool, bool) bool
	AvailableCoins([]Output, bool)
	CanSupportFeature(int) bool
	ChangeWalletPassphrase(string, string) bool
	CommitTransaction(*Tx, *ReserveKey) bool
	CreateTransaction(*key.Script, int64, *Tx, *ReserveKey, int64, string) bool
	CreateTransactions([]map[*key.Script]int64, *Tx, *ReserveKey, int64, string) bool
	DelAddressBookName(*key.TxDestination) bool
	EncryptWallet(string)
	EraseFromWallet(*Uint.U256) bool
	GenerateNewKey() *key.Pub
	GetAddressBalances() map[*key.TxDestination]int64
	GetAddressGroupings() []key.TxDestination
	GetAllReserveKeys() []key.ID
	GetBalance() int64
	GetChange(*tx.Out) int64
	GetCredit(*tx.Out) int64
	GetDebit(*tx.In) int64
	GetImmatureBalance() int64
	GetKeyBirthTimes(map[*key.ID]int64)
	GetKeyFromPool(*key.Pub, bool) bool
	GetKeyPoolSize() int
	GetOldestKeyPoolTime() int64
	GetTransaction(*Uint.U256, *Tx) bool
	GetTxChange(*tx.Transaction) int64
	GetTxCredit(*tx.Transaction) int64
	GetTxDebit(*tx.Transaction) int64
	GetUnconfirmedBalance() int64
	GetVersion() int
	IncOrderPosNext(*DB) int64
	Inventory(*Uint.U256)
	IsChange(*tx.Out) bool
	IsFromMe(*tx.Transaction) bool
	IsLockedCoin(*Uint.U256, uint) bool
	IsMyTX(*tx.Transaction) bool
	IsMyTxIn(*tx.In) bool
	IsMyTxOut(*tx.Out) bool
	KeepKey(int64)
	ListLockedCoins([]tx.OutPoint)
	LoadCryptedKey(*key.Pub, []byte) bool
	LoadKey(*key.Priv, *key.Pub) bool
	LoadKeyMetadata(*key.Pub, *KeyMetadata) bool
	LoadMinVersion(int) bool
	LoadScript(*key.Script) bool
	LoadWallet(bool) error
	LockCoin(*tx.OutPoint)
	MarkDirty()
	NewKeyPool() bool
	NotifyAddressBookChanged(*Wallet, *key.TxDestination, string, bool, int)
	NotifyTransactionChanged(*Wallet, *Uint.U256, int)
	OrderedTxItems([]AccountingEntry, string) *TxItems
	PrintWallet(*block.Block)
	ReacceptWalletTransactions()
	ResendWalletTransactions()
	ReserveKeyFromKeyPool(int64, *KeyPool)
	ReturnKey(int64)
	ScanForWalletTransactions(*block.Index, bool) int
	SelectCoinsMinConf(int64, int, int, []Output) error
	SendMoney(*key.Script, int64, *Tx, bool) string
	SendMoneyToDestination(*key.TxDestination) string
	SetAddressBookName(*key.TxDestination, string) bool
	SetBestChaing(*block.Locator)
	SetDefaultKey(*key.Pub) bool
	SetMaxVersion(int) bool
	SetMinVersion(int, *DB, bool) bool
	TopUpKeyPool() bool
	Unlock(string) bool
	UnlockAllCoins()
	UnlockCoin(*tx.OutPoint)
	UpdatedTransaction(*Uint.U256)
	WalletUpdateSpent(*tx.Transaction)
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
// NewFromFile makes a new wallet by importing a wallet.dat file
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
