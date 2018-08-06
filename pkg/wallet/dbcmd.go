package wallet
import (
	"errors"
	"os"
	"time"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
)
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
// Encrypt a wallet.dat database
func (db *DB) Encrypt() (err error) {
	return
}
// Unlock a wallet.dat database
func (db *DB) Unlock() (err error) {
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
// Version returns the version of the wallet
func (db *DB) Version() int {
	return FeatureBase
}
// GetBalance gets the balance of the wallet
func (db *DB) GetBalance() float64 {
	return 0.0
}
// GetOldestKeyPoolTime gets the oldest keypool time
func (db *DB) GetOldestKeyPoolTime() int64 {
	return 0
}
// GetKeyPoolSize gets the keypool size
func (db *DB) GetKeyPoolSize() int {
	return 0
}
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
// WriteName writes a new name to the database associated with an address
func (db *DB) WriteName(addr, name string) (err error) {
	r := db.KVEnc([]interface{}{"name", addr, name})
	if err = db.Put(bdb.NoTransaction, true, r); err != nil {
		return
	}
	db.updateCount++
	return
}
// EraseName deletes a name from the wallet
func (db *DB) EraseName(addr string) (err error) {
	r := db.KVEnc([]interface{}{"name", addr})
	if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
		return err
	}
	db.updateCount++
	return
}
// WriteTx writes a transaction to the wallet
func (db *DB) WriteTx(u *Uint.U256, t []byte) (err error) {
	r := db.KVEnc([]interface{}{"tx", u, t})
	if err = db.Put(bdb.NoTransaction, false, r); err != nil {
		return
	}
	db.updateCount++
	return
}
// EraseTx deletes a transaction from the wallet
func (db *DB) EraseTx(u *Uint.U256) (err error) {
	r := db.KVEnc([]interface{}{"tx", u})
	if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
		return err
	}
	db.updateCount++
	return
}
// WriteKey writes a new key to the wallet
func (db *DB) WriteKey(pub *key.Pub, priv *key.Priv, meta *KeyMetadata) (err error) {
	rKey := db.KVEnc([]interface{}{"key", pub, priv})
	rMeta := db.KVEnc([]interface{}{"keymeta", pub, meta.Version, meta.CreateTime})
	if err = db.Put(bdb.NoTransaction, false, rKey); err != nil {
		return
	} else if err = db.Put(bdb.NoTransaction, false, rMeta); err != nil {
		return
	}
	db.updateCount++
	return
}
func (db *DB) eraseKey(pub *key.Pub) (err error) {
	rKey := db.KVEnc([]interface{}{"key", pub})
	rMeta := db.KVEnc([]interface{}{"keymeta", pub})
	if err = db.Del(bdb.NoTransaction, rKey[0]); err != nil {
		return
	} else if err = db.Del(bdb.NoTransaction, rMeta[0]); err != nil {
		return
	}
	db.updateCount++
	return
}
// WriteCryptedKey writes an encrypted key to the wallet
func (db *DB) WriteCryptedKey(*key.Pub, []byte, *KeyMetadata) (err error) {
	db.updateCount++
	return
}
// WriteMasterKey writes a MasterKey to the wallet
func (db *DB) WriteMasterKey(id int64, mkey *crypto.MasterKey) (err error) {
	r := db.KVEnc([]interface{}{"mkey", id, mkey})
	if err = db.Put(bdb.NoTransaction, false, r); err != nil {
		return
	}
	db.updateCount++
	return
}
func (db *DB) eraseMasterKey(id int64) (err error) {
	r := db.KVEnc([]interface{}{"mkey", id})
	if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
		return
	}
	db.updateCount++
	return
}
// WriteScript writes a script to the wallet
func (db *DB) WriteScript(hashID *Uint.U160, script *key.Script) (err error) {
	r := db.KVEnc([]interface{}{"cscript", hashID, script})
	if err = db.Put(bdb.NoTransaction, false, r); err != nil {
		return
	}
	db.updateCount++
	return
}
func (db *DB) eraseScript(hashID *Uint.U160) (err error) {
	r := db.KVEnc([]interface{}{"cscript", hashID})
	if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
		return
	}
	db.updateCount++
	return
}
// WriteBestBlock writes the best block to the wallet
func (db *DB) WriteBestBlock(*block.Locator) (err error) {
	return
}
// ReadBestBlock returns the best block stored in the wallet
func (db *DB) ReadBestBlock(*block.Locator) (err error) {
	return
}
// WriteOrderPosNext moves the write position to the next
func (db *DB) WriteOrderPosNext(p int64) (err error) {
	r := db.KVEnc([]interface{}{"orderposnext", p})
	if err = db.Put(bdb.NoTransaction, true, r); err != nil {
		return
	}
	db.updateCount++
	return
}
func (db *DB) EraseOrderPosNext() (err error) {
	r := db.KVEnc([]interface{}{"orderposnext"})
	if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
		return
	}
	db.updateCount++
	return
}
// WriteDefaultKey writes the default key
func (db *DB) WriteDefaultKey(p *key.Pub) (err error) {
	r := db.KVEnc([]interface{}{"defaultkey", p})
	if err = db.Put(bdb.NoTransaction, true, r); err != nil {
		return
	}
	db.updateCount++
	return
}
func (db *DB) EraseDefaultKey() (err error) {
	r := db.KVEnc([]interface{}{"defaultkey"})
	if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
		return
	}
	db.updateCount++
	return
}
// ReadPool returns the KeyPool
func (db *DB) ReadPool(int64, KeyPool) (err error) {
	return
}
// WritePool writes to the KeyPool
func (db *DB) WritePool(int64, KeyPool) (err error) {
	return
}
// ErasePool erases a KeyPool
func (db *DB) ErasePool(int64) (err error) {
	return
}
// ReadSetting reads a setting (obsolete)
func (db *DB) ReadSetting(string, interface{}) (err error) {
	return
}
// WriteSetting writes a setting (obsolete)
func (db *DB) WriteSetting(string, interface{}) (err error) {
	return
}
// EraseSetting erases a setting (obsolete)
func (db *DB) EraseSetting(string) (err error) {
	return
}
// WriteMinVersion writes the MinVersion
func (db *DB) WriteMinVersion(int) (err error) {
	return
}
// ReadAccount returns the data of an Account
func (db *DB) ReadAccount(accname string, acc *Account) (err error) {
	return
}
// WriteAccount writes the data of an Account
func (db *DB) WriteAccount(string, *Account) (err error) {
	return
}
func (db *DB) writeAccountingEntry(uint64, *AccountingEntry) (err error) {
	return
}
// WriteAccountingEntry writes an AccountingEntry to the wallet
func (db *DB) WriteAccountingEntry(*AccountingEntry) (err error) {
	return
}
// GetAccountCreditDebit gets the Account credit/debit
func (db *DB) GetAccountCreditDebit(string) (err error) {
	return
}
// ListAccountCreditDebit gets the list off accounts and their credit/debits
func (db *DB) ListAccountCreditDebit(string, *[]AccountingEntry) (err error) {
	return
}
// ReorderTransactions reorders transactions in the wallet
func (db *DB) ReorderTransactions(*Wallet) (err error) {
	return
}
// LoadWallet loads the wallet
func (db *DB) LoadWallet(*Wallet) (err error) {
	return
}
// RecoverOnlyKeys recovers only the keys from the wallet
func (db *DB) RecoverOnlyKeys(*bdb.Environment, string) (err error) {
	return
}
func (db *DB) Recover(*bdb.Environment, string) (err error) {
	return
}
// Backup copies the current wallet to another location
func (db *DB) Backup(*Wallet, string) (err error) {
	return
}
// GetCursor returns a cursor to walk over the wallet database
func (db *DB) GetCursor() *bdb.Cursor {
	cursor, _ := db.Cursor(bdb.NoTransaction)
	return &cursor
}
func (db *DB) GetupdateCount() uint64 {
	return db.updateCount
}