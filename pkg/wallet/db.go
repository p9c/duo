package wallet
import (
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"sync"
	"time"
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
type DB struct {
	*bdb.Database
	Filename      string
	UnlockedUntil int64
	mutex         sync.Mutex
	updateCount uint64
}
type dB interface {
	Backup(*Wallet, string) error
	Close() error
	Dump() (string, error)
	Encrypt() error
	Erase(*interface{}) bool
	EraseName(string) error
	ErasePool(int64) error
	EraseSetting(string) error
	EraseTx(Uint.U256)
	Exists(*interface{})
	Flush()
	GetAccountCreditDebit(string) int64
	GetBalance() float64
	GetCursor() *bdb.Cursor
	GetKeyPoolSize() int
	GetOldestKeyPoolTime() int64
	GetupdateCount() uint64
	KVDec([]byte, []byte) interface{}
	KVEnc(interface{}) *[2][]byte
	KVToString(*[2][]byte) (string, bool)
	ListAccountCreditDebit(*string, *[]AccountingEntry)
	LoadWallet(Wallet) error
	Open() error
	Read(*interface{}, *interface{}) bool
	ReadAccount(string, *Account) error
	ReadBestBlock(*block.Locator) error
	ReadPool(int64, KeyPool) error
	ReadSetting(string, interface{}) error
	Recover(*bdb.Environment, string) error
	RecoverOnlyKeys(*bdb.Environment, string) error
	ReorderTransactions(Wallet) error
	SetFilename(string)
	StringToVars(string) interface{}
	Unlock() error
	Verify() error
	Version() int
	Write(*interface{}, *interface{}, bool) bool
	WriteAccount(string, *Account) error
	WriteAccountingEntry(*AccountingEntry) error
	writeAccountingEntry(uint64, *AccountingEntry) error
	WriteBestBlock(*block.Locator) error
	WriteCryptedKey(*key.Pub, []byte, *KeyMetadata) error
	WriteDefaultKey(*key.Pub) error
	WriteKey(*key.Pub, *key.Priv, *KeyMetadata) error
	WriteMasterKey(uint, *crypto.MasterKey) error
	WriteMinVersion(int) error
	WriteName(string, string) error
	WriteOrderPosNext(int64) error
	WritePool(int64, KeyPool) error
	WriteScript(Uint.U160, *key.Script) error
	WriteSetting(string, interface{}) error
	WriteTx(Uint.U256, *Tx)
}
