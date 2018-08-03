package walletdat

import (
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
)

func init() {
	Prefix = make(map[string][]byte)
	for i := range KeyNames {
		Prefix[KeyNames[i]] = append([]byte{byte(len(KeyNames[i]))}, []byte(KeyNames[i])...)
	}
	Filename = *args.DataDir + "/" + *args.Wallet
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
