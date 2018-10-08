package db

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/bc"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

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
			value, _ := item.Value()
			table := key[:8]
			if bytes.Compare(table, rec.Tables["MasterKey"]) == 0 {
				idx := key[8:16]
				crypt := value[:48]
				iv := value[48:60]
				iterations := value[60:68]
				var it int64
				core.BytesToInt(&it, &iterations)
				BC = append(BC,
					&bc.BlockCrypt{
						Idx:        &idx,
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
	o := make([]byte, 8)
	rand.Read(o)
	out := &o
	BC.Idx = out
	key := append(rec.Tables["MasterKey"], *out...)
	value := *BC.Crypt.Val
	value = append(value, *BC.IV.Bytes()...)
	value = append(value, *core.IntToBytes(BC.Iterations)...)
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.SetWithMeta(key, value, 0)
	}))
	return r
}

// EraseMasterKey deletes a masterkey entry from the database
func (r *DB) EraseMasterKey(idx *[]byte) *DB {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	search := append(rec.Tables["MasterKey"], *idx...)
	txn := r.DB.NewTransaction(true)
	_, err := txn.Get(search)
	if !r.SetStatusIf(err).OK() {
		return r
	}
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(search)
	}))
	return r
}

// EraseAllMasterKeys deletes all of the masterkeys
func (r *DB) EraseAllMasterKeys() *DB {
	counter := 0
	opt := badger.DefaultIteratorOptions
	err := r.DB.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			counter++
			item := iter.Item()
			if string(item.Key()[:8]) == rec.TS["MasterKey"] {
				r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
					return txn.Delete(item.Key())
				}))
			}
		}
		return nil
	})
	if !r.SetStatusIf(err).OK() {
		fmt.Println("\nERROR:", r.Error())
	}
	fmt.Println(counter, "master keys deleted")
	return r
}
