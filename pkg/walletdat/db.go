package walletdat

import (
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"sync"
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"errors"
	"os"
	"time"

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
		// DBTypes is human-readable strings associated with bdb database types
	DBTypes = map[bdb.DatabaseType]string{
		bdb.BTree:    "btree",
		bdb.Hash:     "hash",
		bdb.Numbered: "numbered",
		bdb.Queue:    "queue",
		bdb.Unknown:  "unknown",
	}
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

// NewDB creates a new database file
func NewDB(opts ...string) (db *DB, err error) {
	db = &DB{}
	var filename, environment string
	create := false
	switch {
	case len(opts) == 0:
		filename = *args.DataDir + "/" + *args.Wallet
		environment = *args.DataDir
	case len(opts) == 1:
		filename = opts[0]
		environment = *args.DataDir
	case len(opts) == 2:
		filename = opts[0]
		environment = opts[1]
	default:
		return &DB{}, errors.New("Excess arguments")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		create = true
	}
	dbenvconf := bdb.EnvironmentConfig{
		Create:        true,
		Mode:          0600,
		Transactional: true,
	}
	dbenv, err := bdb.OpenEnvironment(environment, &dbenvconf)
	dbconfig := bdb.DatabaseConfig{
		Create: create,
		Mode:   0600,
		Name:   "main",
		Type:   bdb.BTree,
	}
	db1, err := bdb.OpenDatabase(dbenv, bdb.NoTransaction, filename, &dbconfig)
	if err == nil {
		db.Database = &db1
		db.Filename = filename
		db.UnlockedUntil = time.Now().Add(Locktime).Unix()
	}
	return
}

// SetFilename changes the name of the database we want to open
func (db *DB) SetFilename(filename string) {
	db.Filename = filename
}

// Find returns a list of key/value pairs according to the key label and first part of the content of the key (it must be a complete first several values including length prefixes for strings and slices). This is iterative and will take on average the cost of visiting half the database. Fortunately they are generally small but this is not ideal.
func (db *DB) Find(label string, content []byte) (result [][2][]byte, err error) {
	labelB := append([]byte{byte(len(label))}, []byte(label)...)
	contentB := append([]byte{byte(len(content))}, content...)
	var cursor bdb.Cursor
	if cursor, err = db.Cursor(bdb.NoTransaction); err != nil {
		return
	}
	var rec [2][]byte
	if err = cursor.First(&rec); err != nil {
		return
	}
	matchfail := false
	for {
		for i := 0; i < len(labelB) && !matchfail; i++ {
			if rec[0][i] != labelB[i] {
				matchfail = true
			}
		}
		if content != nil {
			rem := []byte(string(rec[0]))[len(labelB):]
			for i := 0; i < len(rem) && !matchfail; i++ {
				if rem[i] != contentB[i] {
					matchfail = true
				}
			}
		}
		if matchfail {
			break
		} else {
			result = append(result, rec)
			if err = cursor.Next(&rec); err != nil {
				return result, nil
			}
		}
	}
	return
}
