package walletdat

import (
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"os"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
)

var (
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
}

// DB is an interface to a wallet.dat file
type DB struct {
	*bdb.Database
	Filename      string
}

type dB interface {
	Open() error
	Close() error
	Verify() error
}


// Open a wallet.dat file
func (db *DB) Open() (err error) {
	dbenvconf := bdb.EnvironmentConfig{
		Create:        true,
		Recover:       true,
		Mode:          0600,
		Transactional: true,
	}
	dbenv, err := bdb.OpenEnvironment(*args.DataDir, &dbenvconf)
	if err != nil {
		return
	}
	dbconfig := bdb.DatabaseConfig{
		Create: false,
		Mode:   0600,
		Name:   "main",
	}
	db1, err := bdb.OpenDatabase(dbenv, bdb.NoTransaction, db.Filename, &dbconfig)
	if err == nil {
		db.Database = &db1
	} else {
		logger.Debug("Failed to open database", err)
		return
	}
	return
}

// Close an wallet.dat file
func (db *DB) Close() (err error) {
	err = db.Database.Close()
	return
}

// Verify the consistency of a wallet.dat database
func (db *DB) Verify() (err error) {
	if _, err = os.Stat(db.Filename); os.IsNotExist(err) {
		logger.Debug(err)
		return
	}
	if err = bdb.Verify(db.Filename); err != nil {
		logger.Debug(err)
		return
	}
	return
}

// SetFilename changes the name of the database we want to open
func (db *DB) SetFilename(filename string) {
	db.Filename = filename
}
