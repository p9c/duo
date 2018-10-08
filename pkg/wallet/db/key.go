package db

import (
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// ReadKey reads a key entry from the database
func (r *DB) ReadKey(address *[]byte) (out *key.Priv) {
	r = r.NewIf()
	if !r.OK() {
		return &key.Priv{}
	}
	out = key.NewPriv()
	idx := core.Hash64(address)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
	}
	k := []byte(rec.Tables["Key"])
	k = append(k, *idx...)
	k = append(k, *address...)
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
	idx := core.Hash64(address)
	var pk, pp *[]byte
	var meta byte
	if r.BC != nil {
		meta = 1
		address = r.BC.Encrypt(address)
		pk = r.BC.Encrypt(priv.Bytes())
		pp = r.BC.Encrypt(priv.PubKey().Bytes())
	} else {
		pk = priv.Bytes()
		pp = priv.PubKey().Bytes()
	}
	k := []byte(rec.Tables["Key"])
	k = append(k, *idx...)
	k = append(k, *address...)
	// fmt.Println("    WriteKey", hex.EncodeToString(k))
	v := *pk
	v = append(v, *pp...)
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.SetWithMeta(k, v, meta)
	}))
	return r
}

// EraseKey deletes a key entry
func (r *DB) EraseKey(address *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	idx := core.Hash64(address)
	search := append(rec.Tables["Key"], *idx...)
	if r.BC != nil {
		ID := buf.NewSecure().Copy(address).(*buf.Secure)
		addr := r.Encrypt(ID)
		search = append(search, *addr.Bytes()...)
	} else {
		search = append(search, *address...)
	}
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(search)
	}))
	return r
}
