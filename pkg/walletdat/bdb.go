package walletdat

import (
	"encoding/hex"
	"os"
	"strconv"
	"sync"
	"time"

	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
)

var (
	// DBTypes is human-readable strings associated with bdb database types
	DBTypes = map[bdb.DatabaseType]string{
		bdb.BTree:    "btree",
		bdb.Hash:     "hash",
		bdb.Numbered: "numbered",
		bdb.Queue:    "queue",
		bdb.Unknown:  "unknown",
	}
	// Filename is the default filename being in the data directory for the wallet file
	Filename string
	// Db is a shared wallet for the typical application using one
	Db DB
	// Prefix is loaded in init to contain KeyNames
	Prefix map[string][]byte
	// KeyNames is the list of key types stored in the wallet
	KeyNames = []string{"name", "tx", "acentry", "key", "wkey", "mkey", "ckey", "keymeta", "defaultkey", "pool", "version", "cscript", "orderposnext", "account", "setting", "bestblock", "minversion"}
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
	updateCount   uint64
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
		db.UnlockedUntil = time.Now().Add(Locktime).Unix()
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

// Dump the set of keys and current stats of the chain in a string
func (db *DB) Dump() (dump string, err error) {
	// cursor, err := db.Cursor(bdb.NoTransaction)
	// if err != nil {
	// 	return "", err
	// }
	// rec := [2][]byte{}
	// err = cursor.First(&rec)
	// if err != nil {
	// 	return "", err
	// }
	// dbt, _ := db.Type()
	// dump += "databasetype " + DBTypes[dbt] + "\n"
	// for {
	// 	dump1 := db.KVToString(rec)
	// 	if dump1 != "" {
	// 		dump += dump1
	// 	} else {
	// 		dump += "key " + strconv.Itoa(len(rec[0])) + " " + hex.EncodeToString(rec[0]) +
	// 			" " + string(rec[0]) + "\n"
	// 		dump += "value " + strconv.Itoa(len(rec[1])) + " " + hex.EncodeToString(rec[1]) + "\n"
	// 	}
	// 	err = cursor.Next(&rec)
	// 	if err != nil {
	// 		err = cursor.Close()
	// 		break
	// 	}
	// }
	return
}
