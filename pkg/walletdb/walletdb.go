package walletdb

import (
	"bytes"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/mitchellh/go-homedir"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/proto"
	"github.com/parallelcointeam/duo/pkg/walletdb/entries"
)

// NewWalletDB creates a new walletdb.DB. Path, BaseDir, ValueDir the order of how the variadic options will be processed to override thte defaults
func NewWalletDB(params ...string) (db *DB) {
	var err error
	db = &DB{
		BaseDir:  DefaultBaseDir,
		ValueDir: DefaultValueDir,
	}
	if db.Path, err = homedir.Dir(); err != nil {
		db.SetStatus(err.Error())
		return
	}
	db.Options = &badger.DefaultOptions
	l := len(params)
	if l > 0 {
		switch {
		case l >= 1:
			db.Path = params[0]
		case l >= 2:
			db.BaseDir = params[1]
		case l >= 3:
			db.ValueDir = params[3]
		}
	}
	db.Options.Dir = db.Path + "/" + db.BaseDir
	db.Options.ValueDir = db.Path + "/" + db.BaseDir + "/" + db.ValueDir
	if db.DB, err = badger.Open(*db.Options); !db.SetStatusIf(err).OK() {
		return
	}
	return
}

// NewIf creates a new WalletDB
func (r *DB) NewIf() *DB {
	if r == nil {
		r = NewWalletDB()
		r.SetStatus(er.NilRec)
	}
	return r
}

// Close shuts down a wallet database
func (r *DB) Close() {
	r.DB.Close()
}

// WriteName writes a name entry to the database
func (r *DB) WriteName() {}

// EraseName removes a name entry from the database
func (r *DB) EraseName() {}

// ReadTx reads a transaction entry from the database
func (r *DB) ReadTx() {}

// WriteTx writes a transaction entry from the database
func (r *DB) WriteTx() {}

// EraseTx deletes a transaction entry from the database
func (r *DB) EraseTx() {}

// WriteKey writes a key entry to the database
func (r *DB) WriteKey() {}

// ReadMasterKeys returns all of the masterkey entries in the database
func (r *DB) ReadMasterKeys() (BC []*bc.BlockCrypt) {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := item.Key()
			fmt.Println("key", key)
			value, _ := item.Value()
			fmt.Println("value", value)
			table := key[:8]
			if bytes.Compare(table, rec.Tables["MasterKey"]) == 0 {
				idx := key[8:16]
				crypt := value[:48]
				cryptHash := proto.Hash64(&crypt)
				if bytes.Compare(idx, *cryptHash) != 0 {
					r.SetStatus("index of crypt was incorrect")
				}
				iv := value[48:60]
				iterations := value[60:68]
				var it int64
				proto.BytesToInt(&it, &iterations)
				BC = append(BC,
					&bc.BlockCrypt{
						Crypt:      buf.NewByte().Copy(&crypt).(*buf.Byte),
						IV:         buf.NewByte().Copy(&iv).(*buf.Byte),
						Iterations: it,
					})
			}
		}
		return nil
	})
	r.SetStatusIf(err)
	return
}

// WriteMasterKey adds a master key entry to the database
func (r *DB) WriteMasterKey(BC *bc.BlockCrypt) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if BC.Crypt.Len() < 1 {
		r.SetStatus("zero length crypt")
		return r
	}
	out := proto.Hash64(BC.Crypt.Bytes())
	key := append(rec.Tables["MasterKey"], *out...)
	value := *BC.Crypt.Bytes()
	value = append(value, *BC.IV.Bytes()...)
	value = append(value, *proto.IntToBytes(BC.Iterations)...)
	txn := r.DB.NewTransaction(true)
	err := txn.Set(key, value)
	r.SetStatusIf(err)
	txn.Commit(nil)
	return r
}

// EraseMasterKey deletes a masterkey entry from the database
func (r *DB) EraseMasterKey(idx *[]byte) *DB {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	search := append(rec.Tables["MasterKey"], *idx...)
	txn := r.DB.NewTransaction(true)
	if !r.SetStatusIf(txn.Delete(search)).OK() {
		fmt.Println(r.Error())
	}
	txn.Commit(nil)
	return r
}

// WriteScript writes a script entry to the database
func (r *DB) WriteScript() {}

// EraseScript deletes a script entry from the database
func (r *DB) EraseScript() {}

// WriteBestBlock updates the best block entry to the current head
func (r *DB) WriteBestBlock() {}

// WriteDefaultKey updates the default key used by interfaces when receiving payments
func (r *DB) WriteDefaultKey() {}

// ReadDefaultKey returns the current set default key
func (r *DB) ReadDefaultKey() {}

// ReadBestBlock gets the current best block entry
func (r *DB) ReadBestBlock() {}

// ReadPool gets the oldest available pool entry and refreshes the pool after addresses are used
func (r *DB) ReadPool() {}

// WritePool adds a new pool key to the wallet
func (r *DB) WritePool() {}

// ReadMinVersion returns the minimum version required to read this database
func (r *DB) ReadMinVersion() {}

// WriteMinVersion updates the minimum version
func (r *DB) WriteMinVersion() {}

// ReadAccount finds an account stored due to being a correspondent account
func (r *DB) ReadAccount() {}

// WriteAccount writes a new account entry
func (r *DB) WriteAccount() {}

// EraseAccount deletes an account from the wallet database
func (r *DB) EraseAccount() {}

// WriteAccountingEntry writes an accounting entry based on a transaction
func (r *DB) WriteAccountingEntry() {}

// GetAccountCreditDebit finds entries in the credit/debit records written related to each input transaction from a list of indexes of accounts of interest
func (r *DB) GetAccountCreditDebit() {}

// LoadWallet opens a wallet ready to use
func (r *DB) LoadWallet() {}

// Recover attempts to recover as much data as possible from the database files by parsing their key and value tables as raw data
func (r *DB) Recover() {}
