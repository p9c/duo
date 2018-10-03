package walletdb

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/mitchellh/go-homedir"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/wallet/db/entries"
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
	} else {
		r.UnsetStatus()
	}
	return r
}

func (r *DB) dump() {
	fmt.Println("\nDUMP")
	counter := 0
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			counter++
			item := iter.Item()
			k := item.Key()
			v, err := item.Value()
			meta := item.UserMeta()
			t := rec.TS
			for i := range t {
				if bytes.Compare(k[:8], []byte(t[i])) == 0 {
					fmt.Printf("\t%s\n", i)
				}
			}
			fmt.Println("\t\tkey   ", hex.EncodeToString(k))
			fmt.Println("\t\tvalue ", hex.EncodeToString(v))
			fmt.Println("\t\terr   ", err, "\tmeta\t", meta)
		}
		return nil
	})
	if err != nil {
		fmt.Println("\tERROR:", err.Error())
	}
	itemS := "items"
	if counter == 1 {
		itemS = "item"
	}
	fmt.Println("\t", counter, itemS, "found")
}

func (r *DB) deleteAll() {
	// fmt.Print("\nDELETE ALL\t")
	counter := 0
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		found := false
		for iter.Rewind(); iter.Valid(); iter.Next() {
			counter++
			item := iter.Item()
			k := item.Key()
			t := rec.TS
			for i := range t {
				if bytes.Compare(k[:8], []byte(t[i])) == 0 {
					fmt.Print("\ndeleted item type ", i, " index ", hex.EncodeToString(k[8:16]))
				}
				err := txn.Delete(k)
				if err != nil {
					// fmt.Println("\nERROR", err.Error())
				} else {
					err := txn.Delete(k)
					if err != nil {
						// fmt.Println("\nERROR", err.Error())
					}
				}
			}
			if !found {
				if !r.SetStatusIf(txn.Delete(k)).OK() {
					// fmt.Println("\nERROR:", r.Status)
				}
			}
		}
		return nil
	})
	if err != nil {
		// fmt.Println("\nERROR:", err.Error())
	}
	// fmt.Println("\n", counter, "items deleted")
}

// WithBC attaches a BlockCrypt and thus enabling encryption of sensitive data in the wallet. Changes the encryption if already encrypted or enables it.
func (r *DB) WithBC(BC *bc.BlockCrypt) *DB {
	r = r.NewIf()
	if BC != nil {
		r.BC = BC
	} else {
		r.SetStatus(er.NilParam)
		return r
	}
	opt := badger.DefaultIteratorOptions
	err := r.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			k := item.Key()
			v, err := item.Value()
			if !r.SetStatusIf(err).OK() {
				return r
			}
			// meta := item.UserMeta()
			t := rec.TS
			switch string(k[:8]) {
			case t["Account"]:
				if r.SetStatusIf(txn.Delete(k)).OK() {
					addr := k[16:]
					r.WriteAccount(&addr, &v)
				}
			case t["Name"]:
				if r.SetStatusIf(txn.Delete(k)).OK() {
					addr := k[16:]
					r.WriteName(&addr, &v)
				}
			case t["Key"]:
				if r.SetStatusIf(txn.Delete(k)).OK() {
					priv := v[:32]
					pub := v[32:]
					pk := key.NewPriv()
					pk.SetKey(&priv, &pub)
					r.WriteKey(pk)
				}
			case t["Tx"]:
			case t["Seed"]:
			case t["Script"]:
			case t["Pool"]:
			case t["Setting"]:
			case t["Accounting"]:
			case t["CreditDebit"]:
			case t["BestBlock"]:
			case t["MinVersion"]:
			case t["DefaultKey"]:
			}
		}
		iter.Close()
		txn.Commit(nil)
		return nil
	})
	if err != nil {
		// fmt.Println("ERROR:", err.Error())
	}
	r.WriteMasterKey(BC)
	return r
}

// RemoveBC removes the BlockCrypt and decrypts all the records in the database.
func (r *DB) RemoveBC() *DB {
	if r.BC == nil {
		r.SetStatus("BC was already removed")
		return r
	}
	BC := r.BC
	var masterkeyIdx []byte
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			k := item.Key()
			v, err := item.Value()
			meta := item.UserMeta()
			if !r.SetStatusIf(err).OK() {
				return r
			}
			table := string(k[:8])
			t := rec.TS
			// K, V := hex.EncodeToString(k), hex.EncodeToString(v)
			switch table {
			case t["MasterKey"]:
				K := k[8:16]
				r.EraseMasterKey(&K)
			case t["Name"]:
				Naddress := k[16:]
				label := v
				if meta&1 == 1 {
					Naddress = *r.BC.Decrypt(&Naddress)
					label = *r.BC.Decrypt(&label)
				}
				r.EraseName(&Naddress)
				r.BC = nil
				r.WriteName(&Naddress, &label)
				r.BC = BC
			case t["Key"]:
				Kaddress := k[16:]
				priv := v[:48]
				pub := v[48:]
				if meta&1 == 1 {
					Kaddress = *r.BC.Decrypt(&Kaddress)
					priv = *r.BC.Decrypt(&priv)
					pub = *r.BC.Decrypt(&pub)
				}
				r.EraseKey(&Kaddress)
				r.BC = nil
				pk := key.NewPriv()
				pk.SetKey(&priv, &pub)
				r.WriteKey(pk)
				r.BC = BC
			case t["Account"]:
				Aaddress := k[16:]
				Aaddress = *r.BC.Decrypt(&Aaddress)
				r.EraseAccount(&Aaddress)
				pub := v
				if meta&1 == 1 {
					pub = *r.BC.Decrypt(&pub)
				}
				r.BC = nil
				r.WriteAccount(&Aaddress, &pub)
				r.BC = BC
			case t["Tx"]:
				// fmt.Println("\nTx         ", K, "\n           ", V)
			case t["Seed"]:
				// fmt.Println("\nSeed       ", K, "\n           ", V)
			case t["Script"]:
				// fmt.Println("\nScript     ", K, "\n           ", V)
			case t["Pool"]:
				// fmt.Println("\nPool       ", K, "\n           ", V)
			case t["Setting"]:
				// fmt.Println("\nSetting    ", K, "\n           ", V)
			case t["Accounting"]:
				// fmt.Println("\nAccounting ", K, "\n           ", V)
			case t["CreditDebit"]:
				// fmt.Println("\nCreditDebit", K, "\n           ", V)
			case t["BestBlock"]:
				// fmt.Println("\nBestBlock  ", K, "\n           ", V)
			case t["MinVersion"]:
				// fmt.Println("\nMinVersion ", K, "\n           ", V)
			case t["DefaultKey"]:
				// fmt.Println("\nDefaultKey ", K, "\n           ", V)
			}
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
		if masterkeyIdx != nil {
			r.EraseMasterKey(&masterkeyIdx)
		}
	}
	r.BC = nil
	return r
}

// Close shuts down a wallet database
func (r *DB) Close() {
	r.DB.Close()
}

// LoadWallet opens a wallet ready to use
func (r *DB) LoadWallet() {

}

// Recover attempts to recover as much data as possible from the database files by parsing their key and value tables as raw data
func (r *DB) Recover() {}
