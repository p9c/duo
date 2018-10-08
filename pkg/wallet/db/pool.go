package db

import (
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// ReadPool gets the oldest (lowest sequence number) key and moves it to the private key map, and triggers a keypool refill if it drops below KeyPoolLow
func (r *DB) ReadPool() *rec.Pool {
	return nil
}

// WritePool blindly writes a pool record, assuming its indices do not conflict (used by the NewKeyPool function)
func (r *DB) WritePool(newPool *rec.Pool) *DB {
	r = r.NewIf()
	if !r.OK() {
		return r
	}
	t := rec.TS
	address := newPool.Address
	pub := newPool.Pub
	priv := newPool.Priv
	idx := core.Hash64(address.Bytes())
	creB := core.IntToBytes(newPool.Created)
	expB := core.IntToBytes(newPool.Expires)
	var meta byte
	if r.BC != nil {
		meta = 1
		address.Copy(r.BC.Encrypt(address.Bytes()))
		priv.Copy(r.BC.Encrypt(priv.Bytes()))
		pub.Copy(r.BC.Encrypt(pub.Bytes()))
		creB = r.BC.Encrypt(creB)
		expB = r.BC.Encrypt(expB)
	}
	k := []byte(t["Pool"])
	seqB := core.IntToBytes(newPool.Seq)

	k = append(k, *idx...)
	k = append(k, *seqB...)
	k = append(k, *address.Bytes()...)
	k = append(k, *creB...)
	k = append(k, *expB...)

	v := *priv.Bytes()
	v = append(v, *pub.Bytes()...)

	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.SetWithMeta(k, v, meta)
	}))
	return r
}

// ErasePool removes a pool key
func (r *DB) ErasePool(pool *rec.Pool) *DB {
	k := []byte(rec.TS["Pool"])
	k = append(k, pool.Idx...)
	k = append(k, *core.IntToBytes(pool.Seq)...)
	creB := core.IntToBytes(pool.Created)
	expB := core.IntToBytes(pool.Expires)
	if r.BC != nil {
		k = append(k, *r.BC.Encrypt(pool.Address.Bytes())...)
		creB = r.BC.Encrypt(creB)
		expB = r.BC.Encrypt(expB)
	} else {
		k = append(k, *pool.Address.Bytes()...)
	}
	k = append(k, *creB...)
	k = append(k, *expB...)
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(k)
	}))
	return r
}
