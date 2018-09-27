package walletdb

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/proto"
	"github.com/parallelcointeam/duo/pkg/walletdb/entries"
)

// ReadKey reads a key entry from the database
func (r *DB) ReadKey(address *[]byte) (out *key.Priv) {
	k := []byte(rec.Tables["Key"])
	idx := proto.Hash64(address)
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
		var priv, pub *[]byte
		switch {
		case r.BC != nil && meta&1 == 1:
			encpriv := V[:48]
			encpub := V[48:]
			out = key.NewPriv()
			out.WithBC(r.BC)
			priv = r.BC.Decrypt(&encpriv)
			pub = r.BC.Decrypt(&encpub)
		case meta&1 != 1:
			pr, pu := V[:32], V[32:]
			priv, pub = &pr, &pu
		default:
			r.SetStatus("record marked encrypted but no BC to decrypt with")
			fmt.Println(r.Error())
			out = new(key.Priv)
		}
		out.SetKey(priv, pub)
	}
	return
}

// WriteKey writes a key entry to the database
func (r *DB) WriteKey(priv *key.Priv) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if priv.Crypt.Len() < 1 {
		r.SetStatus("zero length crypt")
		return r
	}
	I := []byte(priv.GetID())
	address := &I
	idx := proto.Hash64(address)
	var pk, pp *[]byte
	var meta byte
	if r.BC != nil {
		meta = 1
		address = priv.BC.Encrypt(address)
		pk = priv.BC.Encrypt(priv.Bytes())
		pp = priv.BC.Encrypt(priv.PubKey().Bytes())
	} else {
		pk = priv.Bytes()
		pp = priv.PubKey().Bytes()
	}
	k := []byte(rec.Tables["Key"])
	k = append(k, *idx...)
	k = append(k, *address...)
	v := *pk
	v = append(v, *pp...)
	txn := r.DB.NewTransaction(true)
	err := txn.SetWithMeta(k, v, meta)
	if r.SetStatusIf(err).OK() {
		r.SetStatusIf(txn.Commit(nil))
	}
	return r
}

// EraseKey deletes a key entry
func (r *DB) EraseKey(address *[]byte) *DB {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	idx := proto.Hash64(address)
	search := append(rec.Tables["Key"], *idx...)
	ID := buf.NewSecure().Copy(address).(*buf.Secure)
	encid := r.Encrypt(ID)
	search = append(search, *encid.Bytes()...)
	txn := r.DB.NewTransaction(true)
	if r.SetStatusIf(txn.Delete(search)).OK() {
		txn.Commit(nil)
	}
	return r
}
