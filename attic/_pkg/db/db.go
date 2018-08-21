// A library for interacting with various database types used in the duo blockchain client built on top of JVZC, built on the Badger key/value database from dgraph.io
package db
import (
	"sync"
)
const (
	// VerifyOK is the return for when the database verifies correctly
	VerifyOK = iota
	// RecoverOK means the database was broken and it was fixed
	RecoverOK
	// RecoverFail means the database was broken and could not be fixed
	RecoverFail
)
var (
	// WalletUpdated is a counter that we increment every time we update the database
	WalletUpdated uint
	// BitDB is the database environmennt for the main database
	BitDB Env
)
//
type Db struct {
}
type KV struct {
	Key, Value []byte
}
type Env struct {
	dbEnvInit, mockDB bool
	path              string
	Mutex             sync.RWMutex
	// DB                bolt.DB
	FileUseCount map[string]int
	DBMap        map[string]*Db
}
type DB struct {
	Db   *Db
	File string
	// ActiveTx bolt.Tx
	ReadOnly bool
}
type Addrs struct {
	Path string
}
