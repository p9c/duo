package wallet

import (
	"sync"
	"time"

	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
)

var (
	// Filename is the default filename being in the data directory for the wallet file
	Filename string
	// Locktime is the time delay after which the wallet will automatically lock
	Locktime = time.Minute * 15
	// Db is a shared wallet for the typical application using one
	Db DB
	// Prefix is loaded in init to contain KeyNames
	Prefix map[string][]byte
	// KeyNames is the list of key types stored in the wallet
	KeyNames = []string{"name", "tx", "acentry", "key", "wkey", "mkey", "ckey", "keymeta", "defaultkey", "pool", "version", "cscript", "orderposnext", "acc", "setting", "bestblock", "minversion"}
)

func init() {
	Prefix = make(map[string][]byte)
	for i := range KeyNames {
		Prefix[KeyNames[i]] = append([]byte{byte(len(KeyNames[i]))}, []byte(KeyNames[i])...)
	}
	Filename = *args.DataDir + "/" + *args.Wallet
}

// DB is an interface to a wallet.dat file
type DB struct {
	*bdb.Database
	Filename      string
	UnlockedUntil int64
	mutex         sync.Mutex
	updateCount uint64
}

type dB interface {
	Flush()
	SetFilename(string)
	Open() error
	Close() error
	Verify() error
	Encrypt() error
	Unlock() error
	KVDec([]byte, []byte) interface{}
	KVEnc(interface{}) *[2][]byte
	KVToString(*[2][]byte) (string, bool)
	StringToVars(string) interface{}
	Dump() (string, error)
	Version() int
	GetBalance() float64
	GetOldestKeyPoolTime() int64
	GetKeyPoolSize() int
	WriteName(string, string) error
	EraseName(string) error
	WriteTx(Uint.U256, *Tx)
	EraseTx(Uint.U256)
	WriteKey(*key.Pub, *key.Priv, *KeyMetadata) error
	WriteCryptedKey(*key.Pub, []byte, *KeyMetadata) error
	WriteMasterKey(uint, *crypto.MasterKey) error
	WriteScript(Uint.U160, *key.Script) error
	WriteBestBlock(*block.Locator) error
	ReadBestBlock(*block.Locator) error
	WriteOrderPosNext(int64) error
	WriteDefaultKey(*key.Pub) error
	ReadPool(int64, KeyPool) error
	WritePool(int64, KeyPool) error
	ErasePool(int64) error
	ReadSetting(string, interface{}) error
	WriteSetting(string, interface{}) error
	EraseSetting(string) error
	WriteMinVersion(int) error
	ReadAccount(string, *Account) error
	WriteAccount(string, *Account) error
	writeAccountingEntry(uint64, *AccountingEntry) error
	WriteAccountingEntry(*AccountingEntry) error
	GetAccountCreditDebit(string) int64
	ListAccountCreditDebit(*string, *[]AccountingEntry)
	ReorderTransactions(Wallet) error
	LoadWallet(Wallet) error
	RecoverOnlyKeys(*bdb.Environment, string) error
	Recover(*bdb.Environment, string) error
	Backup(*Wallet, string) error
	Read(*interface{}, *interface{}) bool
	Write(*interface{}, *interface{}, bool) bool
	Erase(*interface{}) bool
	Exists(*interface{})
	GetCursor() *bdb.Cursor
	GetupdateCount() uint64
}

// ScanState stores the state of a wallet
type ScanState struct {
	Keys, CKeys, KeyMeta      uint
	IsEncrypted, AnyUnordered bool
	FileVersion               int
	WalletUpgrade             []*Uint.U256
}
