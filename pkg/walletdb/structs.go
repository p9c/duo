package walletdb

import (
	"github.com/dgraph-io/badger"
)

// DB is the central data repository for the wallet database
type DB struct {
	Path     string
	Dir      string
	ValueDir string
	Options  *badger.Options
	DB       *badger.DB
	Status   string
}
