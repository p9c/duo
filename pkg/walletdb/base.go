package walletdb

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/mitchellh/go-homedir"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/key"
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

// WithBC attaches a BlockCrypt and thus enabling encryption of sensitive data in the wallet. Changes the encryption if already encrypted or enables it.
func (r *DB) WithBC(BC *bc.BlockCrypt) *DB {
	r = r.NewIf()
	if BC != nil {
		r.BC = BC
	}
	// TODO: have it read all entries and rewrite them encrypted and flush all old data
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
	err := r.DB.View(func(txn *badger.Txn) error {
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
			K, V := hex.EncodeToString(k), hex.EncodeToString(v)
			switch table {
			case t["MasterKey"]:
				fmt.Println("\nMasterKey  ", K, "\n           ", V)
				fmt.Println(">>> deleting...")
				masterkeyIdx = k[8:16]
			case t["Name"]:
				fmt.Println("\nName       ", K, "\n           ", V)
				if meta&1 == 1 {
					table := k[:8]
					idx := k[8:16]
					address := k[16:]
					label := v

					fmt.Println("\nENCRYPTED")
					fmt.Println("table  ", hex.EncodeToString(table))
					fmt.Println("idx    ", hex.EncodeToString(idx))
					fmt.Println("address", hex.EncodeToString(address))
					fmt.Println("label  ", hex.EncodeToString(label))

					address = *r.BC.Decrypt(&address)
					label = *r.BC.Decrypt(&label)
					r.EraseName(&address)

					fmt.Println("\nDECRYPTED")
					fmt.Println("address", hex.EncodeToString(address))
					fmt.Println("label  ", string(label))

					r.BC = nil
					r.WriteName(&address, &label)
					r.BC = BC
				}
			case t["Tx"]:
				fmt.Println("\nTx         ", K, "\n           ", V)
			case t["Seed"]:
				fmt.Println("\nSeed       ", K, "\n           ", V)
			case t["Key"]:
				fmt.Println("\nKey        ", K, "\n           ", V)
				if meta&1 == 1 {
					table := k[:8]
					idx := k[8:16]
					address := k[16:]
					priv := v[:48]
					pub := v[48:]

					fmt.Println("\nENCRYPTED")
					fmt.Println("table  ", hex.EncodeToString(table))
					fmt.Println("idx    ", hex.EncodeToString(idx))
					fmt.Println("address", hex.EncodeToString(address))
					fmt.Println("priv   ", hex.EncodeToString(priv))
					fmt.Println("pub    ", hex.EncodeToString(pub))

					address = *r.BC.Decrypt(&address)
					priv = *r.BC.Decrypt(&priv)
					pub = *r.BC.Decrypt(&pub)
					r.EraseKey(&address)

					fmt.Println("\nDECRYPTED")
					fmt.Println("address", hex.EncodeToString(address))
					fmt.Println("priv   ", hex.EncodeToString(priv))
					fmt.Println("pub    ", hex.EncodeToString(pub))

					r.BC = nil
					pk := key.NewPriv().SetKey(&priv, &pub)
					r.WriteKey(pk)
					r.BC = BC
				}
			case t["Script"]:
				fmt.Println("\nScript     ", K, "\n           ", V)
			case t["Pool"]:
				fmt.Println("\nPool       ", K, "\n           ", V)
			case t["Setting"]:
				fmt.Println("\nSetting    ", K, "\n           ", V)
			case t["Account"]:
				fmt.Println("\nAccount    ", K, "\n           ", V)
				if meta&1 == 1 {
					table := k[:8]
					idx := k[8:16]
					address := k[16:]
					pub := v

					fmt.Println("\nENCRYPTED")
					fmt.Println("table  ", hex.EncodeToString(table))
					fmt.Println("idx    ", hex.EncodeToString(idx))
					fmt.Println("address", hex.EncodeToString(address))
					fmt.Println("pub    ", hex.EncodeToString(pub))

					address = *r.BC.Decrypt(&address)
					r.EraseAccount(&address)
					pub = *r.BC.Decrypt(&pub)

					fmt.Println("\nDECRYPTED")
					fmt.Println("address", hex.EncodeToString(address))
					fmt.Println("pub    ", hex.EncodeToString(pub))

					r.BC = nil
					r.WriteAccount(&address, &pub)
					r.BC = BC
				}
			case t["Accounting"]:
				fmt.Println("\nAccounting ", K, "\n           ", V)
			case t["CreditDebit"]:
				fmt.Println("\nCreditDebit", K, "\n           ", V)
			case t["BestBlock"]:
				fmt.Println("\nBestBlock  ", K, "\n           ", V)
			case t["MinVersion"]:
				fmt.Println("\nMinVersion ", K, "\n           ", V)
			case t["DefaultKey"]:
				fmt.Println("\nDefaultKey ", K, "\n           ", V)
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
