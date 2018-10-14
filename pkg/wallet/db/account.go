package db

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// ReadAccount finds an account stored due to being a correspondent account
func (r *DB) ReadAccount(address *[]byte) (out *rec.Account) {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	out = new(rec.Account)
	k := []byte(rec.Tables["Account"])
	idx := core.Hash64(address)
	k = append(k, *idx...)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
	}
	k = append(k, *address...)
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	var V []byte
	var meta byte
	err := r.DB.View(func(txn *badger.Txn) error {
		item, er := txn.Get(k)
		if er != nil {
			return er
		}
		meta = item.UserMeta()
		V, er = item.Value()
		if er != nil {
			return er
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
		out.Idx = *idx
		switch {
		case r.BC != nil && meta&1 == 1:
			out.Address = r.BC.Decrypt(address)
			if len(V) > 1 {
				out.Pub = r.BC.Decrypt(&V)
			}
		case meta&1 != 1:
			out.Address = address
			out.Pub = &V
		default:
			r.SetStatus("record marked encrypted but no BC to decrypt with")
			fmt.Println(r.Error())
			out = new(rec.Account)
		}
	}
	return
}

// WriteAccount writes a new account entry
func (r *DB) WriteAccount(address, pub *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if address == nil {
		r.SetStatus(er.NilParam)
		return r
	}
	if pub == nil {
		pub = &[]byte{}
	}
	idx := core.Hash64(address)
	var k, v []byte
	var meta byte
	if r.BC != nil {
		meta = 1
		address = r.BC.Encrypt(address)
		pub = r.BC.Encrypt(pub)
	}
	k = []byte(rec.Tables["Account"])
	k = append(k, *idx...)
	k = append(k, *address...)
	if pub != nil {
		v = *pub
	} else {
		v = []byte{}
	}
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.SetWithMeta(k, v, meta)
	}))
	return r
}

// EraseAccount deletes an account from the wallet database
func (r *DB) EraseAccount(address *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	idx := core.Hash64(address)
	search := append(rec.Tables["Account"], *idx...)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
	}
	search = append(search, *address...)
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(search)
	}))
	return r
}

// GetAllAccountIDs returns an array containing all of the accounts in the database by their index
func (r *DB) GetAllAccountIDs() (out []rec.Idx) {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			k := item.KeyCopy(nil)
			table := string(k[:8])
			if table == rec.TS["Account"] {
				out = append(out, k[8:16])
			}
		}
		return nil
	}))
	return
}
