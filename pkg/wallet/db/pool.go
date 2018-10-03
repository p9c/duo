package db

import (
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// ReadPool gets the oldest (lowest sequence number) key and moves it to the private key map, and triggers a keypool refill if it drops below KeyPoolLow
func (r *DB) ReadPool() *rec.Pool {
	return nil
}

// WritePool adds a new pool key to the wallet. The sequence number will be set to subsequent to the highest in the pool existing
func (r *DB) WritePool(item *rec.Pool) *DB {
	r = r.NewIf()
	if !r.OK() {
		return r
	}
	out := key.NewPriv()
	address := out.PubKey().Bytes()
	idx := core.Hash64(address)
	var meta byte
	if r.BC != nil {
		meta = 1
		address = r.BC.Encrypt(address)
	}
	k := []byte(rec.Tables["Pool"])
	k = append(k, *idx...)
	k = append(k, *address...)
	// fmt.Println("    ReadKey", hex.EncodeToString(k))
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	counter := 0

	err := r.DB.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			counter++
			item := iter.Item()
			k := item.Key()
			imeta := item.UserMeta()
			t := rec.TS
			if imeta != meta {

			}
			if string(k[:8]) == t["Pool"] {

			}
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
	}
	return r
}

// ErasePool removes a pool key
func (r *DB) ErasePool() {

}
