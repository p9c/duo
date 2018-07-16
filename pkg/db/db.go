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

// KV is a key value structure of store elements in the database
type KV struct {
	Key, Value []byte
}

// Env is the database environment
type Env struct {
	dbEnvInit, mockDB bool
	path              string
	Mutex             sync.RWMutex
	// DB                bolt.DB
	FileUseCount map[string]int
	DBMap        map[string]*Db
}

// DB is a database for the token ledger
type DB struct {
	Db   *Db
	File string
	// ActiveTx bolt.Tx
	ReadOnly bool
}

// Addrs is ?
type Addrs struct {
	Path string
}
