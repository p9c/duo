package wallet

import (
	"github.com/1lann/cete"
	// "github.com/parallelcointeam/duo/pkg/Uint"
	// "github.com/parallelcointeam/duo/pkg/block"
	// "github.com/parallelcointeam/duo/pkg/crypto"
	// "github.com/parallelcointeam/duo/pkg/key"
	// "github.com/parallelcointeam/duo/pkg/server/args"
	"time"
)

const (
	Fname = iota
	Ftx
	Facentry
	Fkey
	Fwkey
	Fmkey
	Fckey
	Fkeymeta
	Fdefaultkey
	Fpool
	Fversion
	Fcscript
	Forderposnext
	Faccount
	Fbestblock
	Fminversion
	Flast
)

var (
	Filename string
	KeyNames = []string{"name", "tx", "acentry", "key", "wkey", "mkey", "ckey", "keymeta", "defaultkey", "pool", "version", "cscript", "orderposnext", "acc", "bestblock", "minversion"}
	K        = KeyNames
)

func init() {
	Filename = *args.DataDir + "/" + *args.Wallet
}

// DB is the structure for encryptable wallet database
type DB struct {
	*jvzc.DB
	UnlockedUntil int64
	updateCount   uint64
	Net           string
	Data          EncryptedStore
}
type dB interface {
	Backup(*Wallet, string) error
	Close() error
	Dump() string
	Encrypt() error
	EraseName(string) error
	EraseMasterKey(int64) error
	ErasePool(int64) error
	EraseTx(Uint.U256)
	Find(int, interface{}) (*jvzc.Range, error)
	Flush()
	GetAccountCreditDebit(string) int64
	GetBalance() float64
	GetKeyPoolSize() int
	GetOldestKeyPoolTime() int64
	GetUpdateCount() uint64
	ImportWalletDat(string) error
	ListAccountCreditDebit(*string, *[]AccountingEntry)
	LoadWallet(Wallet) error
	ReadAccount(string, *Account) error
	ReadBestBlock(*block.Locator) error
	ReadPool(int64, KeyPool) error
	Recover(string) error
	RecoverOnlyKeys(string) error
	ReorderTransactions(Wallet) error
	StringToVars(string) interface{}
	Unlock() error
	Verify() error
	Version() int
	WriteAccount(string, *Account) error
	WriteAccountingEntry(*AccountingEntry) error
	WriteBestBlock(*block.Locator) error
	WriteCryptedKey(*key.Pub, []byte, *KeyMetadata) error
	WriteDefaultKey(*key.Pub) error
	WriteKey(*key.Pub, *key.Priv, *KeyMetadata) error
	WriteMasterKey(*crypto.MasterKey) error
	WriteMinVersion(int) error
	WriteName(string, string) error
	WriteOrderPosNext(int64) error
	WritePool(int64, KeyPool) error
	WriteScript(Uint.U160, *key.Script) error
	WriteTx(Uint.U256, *Tx)
}

// NewDB creates a new DB and opens it - if it already exists it just opens it
func NewDB(path string) (db *DB, err error) {
	db = new(DB)
	db.DB, _ = jvzc.Open(path)
	db.UnlockedUntil = time.Now().Unix()
	return db, err
}
