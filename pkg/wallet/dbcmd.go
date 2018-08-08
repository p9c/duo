package wallet
import (
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"github.com/parallelcointeam/javazacdb"
)
// Backup copies the current wallet to another location
func (db *DB) Backup(*Wallet, string) (err error) {
	return
}
// Close a wallet.jvzc file
func (db *DB) Close() (err error) {
	// err = db.Database.Close()
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
// Encrypt a wallet.dat database
func (db *DB) Encrypt() (err error) {
	return
}
// EraseDefaultKey erases the default key
func (db *DB) EraseDefaultKey() (err error) {
	// r := db.KVEnc([]interface{}{"defaultkey"})
	// if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
	// 	return
	// }
	// db.updateCount++
	return
}
// EraseMasterKey erases a master key record from a wallet database
func (db *DB) EraseMasterKey(id int64) (err error) {
	// r := db.KVEnc([]interface{}{"mkey", id})
	// if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
	// 	return
	// }
	db.updateCount++
	return
}
// EraseName deletes a name from the wallet
func (db *DB) EraseName(addr string) (err error) {
	// r := db.KVEnc([]interface{}{"name", addr})
	// if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
	// 	return err
	// }
	// db.updateCount++
	return
}

// ErasePool erases a KeyPool
func (db *DB) ErasePool(int64) (err error) {
	return
}
// EraseTx deletes a transaction from the wallet
func (db *DB) EraseTx(u *Uint.U256) (err error) {
	// r := db.KVEnc([]interface{}{"tx", u})
	// if err = db.Del(bdb.NoTransaction, r[0]); err != nil {
	// 	return err
	// }
	// db.updateCount++
	return
}
// Find searches the wallet database for an item in a specified table with a given key
func (db *DB) Find(t int, k string) (r *jvzc.Range, err error) {

	return
}
// Flush forces writing to disk of dirty cache
func (db *DB) Flush() {
	return
}
// GetAccountCreditDebit gets the Account credit/debit
func (db *DB) GetAccountCreditDebit(string) (err error) {
	return
}
// GetBalance gets the balance of the wallet
func (db *DB) GetBalance() float64 {
	return 0.0
}
// GetKeyPoolSize gets the keypool size
func (db *DB) GetKeyPoolSize() int {
	return 0
}
// GetOldestKeyPoolTime gets the oldest keypool time
func (db *DB) GetOldestKeyPoolTime() int64 {
	return 0
}
// GetUpdateCount returns the number of times the wallet has been modified in the current session
func (db *DB) GetUpdateCount() uint64 {
	return db.updateCount
}
// ImportWalletDat imports an old style BDB wallet.dat file
func (db *DB) ImportWalletDat(string) (err error) {
	return
}
// ListAccountCreditDebit gets the list off accounts and their credit/debits
func (db *DB) ListAccountCreditDebit(string, []*AccountingEntry) (err error) {
	return
}
// LoadWallet loads the wallet
func (db *DB) LoadWallet(*Wallet) (err error) {
	return
}
// ReadAccount returns the data of an Account
func (db *DB) ReadAccount(accname string, acc *Account) (err error) {
	return
}
// ReadBestBlock returns the best block stored in the wallet
func (db *DB) ReadBestBlock(*block.Locator) (err error) {
	return
}
// ReadPool returns the KeyPool
func (db *DB) ReadPool(int64, KeyPool) (err error) {
	return
}
// Recover recovers everything from the wallet that it can
func (db *DB) Recover(string) (err error) {
	return
}
// RecoverOnlyKeys recovers only the keys from the wallet
func (db *DB) RecoverOnlyKeys(string) (err error) {
	return
}
// ReorderTransactions reorders transactions in the wallet
func (db *DB) ReorderTransactions(*Wallet) (err error) {
	return
}
// Unlock a wallet.dat database
func (db *DB) Unlock() (err error) {
	return
}
// Verify the consistency of a wallet database
func (db *DB) Verify() (err error) {
	// if _, err = os.Stat(db.Filename); os.IsNotExist(err) {
	// 	logger.Debug(err)
	// 	return
	// }
	// if err = bdb.Verify(db.Filename); err != nil {
	// 	logger.Debug(err)
	// 	return	
	// }
	return
}
// Version returns the version of the wallet
func (db *DB) Version() int {
	return FeatureBase
}
// WriteAccount writes the data of an Account
func (db *DB) WriteAccount(string, *Account) (err error) {
	return
}
// WriteAccountingEntry writes an AccountingEntry to the wallet
func (db *DB) WriteAccountingEntry(*AccountingEntry) (err error) {
	return
}
// WriteBestBlock writes the best block to the wallet
func (db *DB) WriteBestBlock(*block.Locator) (err error) {
	return
}
// WriteCryptedKey writes an encrypted key to the wallet
func (db *DB) WriteCryptedKey(pub *key.Pub, priv []byte, meta *KeyMetadata) (err error) {
	db.Table(KN[Fckey]).Set(pub.ToBase58Check(db.Net), priv)
	db.Table(KN[Fkeymeta]).Set(pub.ToBase58Check(db.Net), []interface{}{meta.Version, meta.CreateTime})
	db.updateCount++
	return
}
// WriteDefaultKey writes the default key
func (db *DB) WriteDefaultKey(p *key.Pub) (err error) {
	// r := db.KVEnc([]interface{}{"defaultkey", p})
	// if err = db.Put(bdb.NoTransaction, true, r); err != nil {
	// 	return
	// }
	// db.updateCount++
	return
}
// WriteKey writes a new key to the wallet
func (db *DB) WriteKey(pub *key.Pub, priv *key.Priv, meta *KeyMetadata) (err error) {
	db.Table(KN[Fkey]).Set(pub.ToBase58Check(db.Net), priv.Get())
	db.Table(KN[Fkeymeta]).Set(pub.ToBase58Check(db.Net), []interface{}{meta.Version, meta.CreateTime})
	db.updateCount++
	return
}
// WriteMasterKey writes a MasterKey to the wallet
func (db *DB) WriteMasterKey(mkey MKey) (err error) {
	db.Table(KN[Fmkey]).Set(fmt.Sprint(mkey.MKeyID), []interface{}{
		mkey.EncryptedKey,
		mkey.Salt,
		mkey.Method,
		mkey.Iterations,
		mkey.Other})
	db.updateCount++
	return
}
// WriteMinVersion writes the MinVersion
func (db *DB) WriteMinVersion(int) (err error) {
	return
}
// WriteName writes a new name to the database associated with an address
func (db *DB) WriteName(addr, name string) (err error) {
	db.Table(KN[Fname]).Set(addr, name)
	db.updateCount++
	return
}
// WriteOrderPosNext moves the write position to the next
func (db *DB) WriteOrderPosNext(p int64) (err error) {
	// r := db.KVEnc([]interface{}{"orderposnext", p})
	// if err = db.Put(bdb.NoTransaction, true, r); err != nil {
	// 	return
	// }
	// db.updateCount++
	return
}
// WritePool writes to the KeyPool
func (db *DB) WritePool(int64, KeyPool) (err error) {
	return
}
// WriteScript writes a script to the wallet
func (db *DB) WriteScript(hashID *Uint.U160, script *key.Script) (err error) {
	// r := db.KVEnc([]interface{}{"cscript", hashID, script})
	// if err = db.Put(bdb.NoTransaction, false, r); err != nil {
	// 	return
	// }
	// db.updateCount++
	return
}
// WriteTx writes a transaction to the wallet
func (db *DB) WriteTx(u *Uint.U256, t []byte) (err error) {
	// r := db.KVEnc([]interface{}{"tx", u, t})
	// if err = db.Put(bdb.NoTransaction, false, r); err != nil {
	// 	return
	// }
	// db.updateCount++
	return
}

// WriteWalletKey writes a 'wkey', which contains extra metadata
func (db *DB) WriteWalletKey(wkey *WKey) (err error) {
	db.Table(KN[Fwkey]).Set(wkey.Pub.ToBase58Check(db.Net), []interface{}{
		wkey.Priv.Get(),
		wkey.TimeCreated.Unix(),
		wkey.TimeExpires.Unix(),
		wkey.Comment})
	return
}