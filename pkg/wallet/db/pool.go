package db

import (
	"sort"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// ReadPool gets the oldest (lowest sequence number) key and moves it to the private key map, and triggers a keypool refill if it drops below KeyPoolLow
func (r *DB) ReadPool() *rec.Pool {
	return nil
}

// AddToPool adds a new pool key to the wallet. The sequence number will be set to subsequent to the highest in the pool existing
func (r *DB) AddToPool(newPool *rec.Pool) *DB {
	r = r.NewIf()
	if !r.OK() {
		return r
	}
	// out := key.NewPriv()
	address := buf.NewByte().Copy(newPool.Address.Bytes())
	priv := &newPool.Priv
	idx := core.Hash64(address.Bytes())
	var meta byte
	if r.BC != nil {
		meta = 1
		address.Copy(r.BC.Encrypt(address.Bytes()))
		priv.Copy(r.BC.Encrypt(priv.Bytes()))
	}
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	counter := 0
	pool := make(map[int]core.Address)
	t := rec.TS
	err := r.DB.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			counter++
			item := iter.Item()
			k := item.Key()
			imeta := item.UserMeta()
			if string(k[:8]) == t["Pool"] {
				seq := k[16:24]
				addr := k[24:]
				address.Copy(&addr)
				if imeta&1 == meta&1 {
					address = r.Decrypt(buf.NewByte().
						Copy(address.Bytes()).(*buf.Byte))
				}
				var seqI int
				core.BytesToInt(&seqI, &seq)
				pool[seqI] = core.Address(*address.Bytes())
			}
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
		var poolSorted []int
		for i := range pool {
			poolSorted = append(poolSorted, i)
		}
		if len(poolSorted) < 1 {
			sort.Ints(poolSorted)
			next := poolSorted[len(poolSorted)-1] + 1
			newPool.Seq = next
		}
		k := []byte(t["Pool"])
		seqB := core.IntToBytes(newPool.Seq)
		k = append(k, *idx...)
		k = append(k, *seqB...)
		k = append(k, *address.Bytes()...)
		txn := r.DB.NewTransaction(true)
		err := txn.SetWithMeta(k, *priv.Bytes(), meta)
		if r.SetStatusIf(err).OK() {
			r.SetStatusIf(txn.Commit(nil))
		}
	}
	return r
}

// WritePool blindly writes a pool record, assuming its indices do not conflict (used by the NewKeyPool function)
func (r *DB) WritePool(newPool *rec.Pool) *DB {
	r = r.NewIf()
	if !r.OK() {
		return r
	}
	t := rec.TS
	address := &newPool.Address
	pub := &newPool.Pub
	priv := &newPool.Priv
	idx := core.Hash64(address.Bytes())
	var meta byte
	if r.BC != nil {
		meta = 1
		address.Copy(r.BC.Encrypt(address.Bytes()))
		priv.Copy(r.BC.Encrypt(priv.Bytes()))
		pub.Copy(r.BC.Encrypt(pub.Bytes()))
	}
	k := []byte(t["Pool"])
	seqB := core.IntToBytes(newPool.Seq)
	k = append(k, *idx...)
	k = append(k, *seqB...)
	k = append(k, *address.Bytes()...)
	txn := r.DB.NewTransaction(true)
	v := *priv.Bytes()
	v = append(v, *pub.Bytes()...)
	err := txn.SetWithMeta(k, v, meta)
	if r.SetStatusIf(err).OK() {
		r.SetStatusIf(txn.Commit(nil))
	}
	return r
}

// ErasePool removes a pool key
func (r *DB) ErasePool() {

}
