package wallet
import (
	"github.com/parallelcointeam/javazacdb"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"time"
)

const (
	// Identifiers matching the KeyNames array for convenience
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
	// The default filename being in the data directory for the wallet file
	Filename string
	// The time delay after which the wallet will automatically lock
	Locktime = time.Minute * 15
	// Db is a shared wallet for the typical application using one
	Db DB
	// The string identifiers of the various tables in a wallet database
	KeyNames = []string{"name", "tx", "acentry", "key", "wkey", "mkey", "ckey", "keymeta", "defaultkey", "pool", "version", "cscript", "orderposnext", "acc", "bestblock", "minversion"}
	K = KeyNames
)
func init() {
	Filename = *args.DataDir + "/" + *args.Wallet
}
type Name struct {
	Addr string
	Name string
}
type name interface {
	ToString() string
}
// Returns a string containing one line space separated values for a name record
func (n *Name) ToString() (r  string) {
	r = n.Addr + " " + n.Name
	return
}
type Metadata struct {
	Pub        *key.Pub
	Version    uint32
	CreateTime time.Time
}
type Key struct {
	Pub  *key.Pub
	Priv *key.Priv
}
type WKey struct {
	Pub         *key.Pub
	Priv     *key.Priv
	TimeCreated time.Time
	TimeExpires time.Time
	Comment     string
}
type MKey struct {
	MKeyID int64
	EncryptedKey              []byte
	Salt                      []byte
	Method          uint32
	Iterations          uint32
	Other []byte
}
type CKey struct {
	Pub  *key.Pub
	Priv []byte
}
// DB is the structure for encryptable wallet database
type DB struct {
	*jvzc.DB
	UnlockedUntil int64
	updateCount uint64
	Net string
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
	// Open() error
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
