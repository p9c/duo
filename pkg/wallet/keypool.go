package wallet

import (
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// NewKeyPool creates a new pool of keys in reserve for generating transactions
func (r *Wallet) NewKeyPool() *Wallet {
	if r == nil {
		return New(nil)
	}
	if r.KeyPool != nil {
		r.EmptyKeyPool()
	}
	r.KeyPool = &KeyPool{
		High:     100,
		Low:      10,
		Lifespan: 90 * 24 * time.Hour,
	}
	r.KeyPool.Pool = make(map[int]*rec.Pool)
	for i := 0; i < r.KeyPool.High; i++ {
		var nk *key.Priv
		if r.DB.BC != nil {
			nk = key.NewPriv()
			nk.WithBC(r.DB.BC)
		} else {
			nk = key.NewPriv()
		}
		nk.Make()
		I := []byte(nk.GetID())
		idx := core.Hash64(&I)
		np := &rec.Pool{
			Address: buf.NewByte().Copy(&I).(*buf.Byte),
			Idx:     *idx,
			Seq:     i,
			Priv:    nk.Crypt,
			Pub:     nk.PubKey().(*buf.Byte),
			Created: time.Now().Unix(),
			Expires: time.Now().Add(r.KeyPool.Lifespan).Unix(),
		}
		r.DB.WritePool(np)
		r.KeyPool.Pool[i] = np
		r.KeyPool.Size++
	}
	return r
}

// LoadKeyPool loads the keypool into memory from the database. Note that this does not load the public and private keys, the keys of pool entries contain everything needed to decide which one to use
func (r *Wallet) LoadKeyPool() *Wallet {
	r.KeyPool = &KeyPool{
		High:     100,
		Low:      10,
		Lifespan: 90 * 24 * time.Hour,
	}
	r.KeyPool.Pool = make(PoolMap)
	opt := badger.DefaultIteratorOptions
	err := r.DB.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			k := item.Key()
			meta := item.UserMeta()
			table := string(k[:8])
			if table == rec.TS["Pool"] {
				idx := k[8:16]
				var seq int
				seqB := k[16:24]
				core.BytesToInt(&seq, &seqB)
				var cre, exp int64
				var creB, expB []byte
				var address []byte
				if meta&1 == 1 {
					if r.DB.BC != nil {
						address = k[24:60] // 36 bytes
						address = *r.DB.BC.Decrypt(&address)
						creB = k[60:84] // 24 bytes
						creB = *r.DB.BC.Decrypt(&creB)
						expB = k[84:] // 24 bytes
						expB = *r.DB.BC.Decrypt(&expB)
					} else {
						break
					}
				} else {
					address = k[16:36]
					creB = k[36:44]
					expB = k[44:52]
				}
				core.BytesToInt(&cre, &creB)
				core.BytesToInt(&exp, &expB)
				r.KeyPool.Pool[seq] = &rec.Pool{
					Idx:     idx,
					Address: buf.NewByte().Copy(&address).(*buf.Byte),
					Seq:     seq,
					Created: cre,
					Expires: exp,
				}
				r.KeyPool.Size++
			}
		}
		return nil
	})
	if !r.SetStatusIf(err).OK() {
		fmt.Println("\nERROR:", r.Error())
	}
	return r
}

// AddReserveKey -
// func (r *Wallet) AddReserveKey(kp *rec.Pool) *Wallet { return r }

// ReserveKeyFromKeyPool -
// func (r *Wallet) ReserveKeyFromKeyPool(int64, *rec.Pool) {}

// GetKeyFromPool gets the oldest key form the keypool, if it is not allowed to be reused it is taken out of the pool and stored as a main wallet private/public key pair
func (r *Wallet) GetKeyFromPool(allowReuse bool) (out *key.Priv) {
	if r == nil {
		r.SetStatus(er.NilRec)
		out = key.NewPriv()
		out.SetStatus(er.NilRec)
		return out
	}
	if len(r.KeyPool.Pool) < 1 {
		r.NewKeyPool()
	}
	var sorted []int
	for i := range r.KeyPool.Pool {
		sorted = append(sorted, r.KeyPool.Pool[i].Seq)
	}
	sort.Ints(sorted)
	lowest := sorted[0]
	outKeyPool := r.KeyPool.Pool[lowest]
	k := []byte(rec.TS["Pool"])
	k = append(k, outKeyPool.Idx...)
	k = append(k, *core.IntToBytes(outKeyPool.Seq)...)
	if r.DB.BC != nil {
		k = append(k, *r.DB.BC.Encrypt(outKeyPool.Address.Bytes())...)
		k = append(k, *r.DB.BC.Encrypt(core.IntToBytes(outKeyPool.Created))...)
		k = append(k, *r.DB.BC.Encrypt(core.IntToBytes(outKeyPool.Expires))...)
	} else {
		k = append(k, *outKeyPool.Address.Bytes()...)
		k = append(k, *core.IntToBytes(outKeyPool.Created)...)
		k = append(k, *core.IntToBytes(outKeyPool.Expires)...)
	}
	txn := r.DB.DB.NewTransaction(true)
	defer txn.Commit(nil)
	item, err := txn.Get(k)
	if r.SetStatusIf(err).OK() {
		v, err := item.Value()
		meta := item.UserMeta()
		if r.SetStatusIf(err).OK() {
			if meta&1 == 1 {
				if r.DB.BC == nil {
					r.SetStatus("no crypt for crypted record")
					out = &key.Priv{}
					out.SetStatus(r.Error())
					return
				}
				privB := v[:64]
				pubB := v[64:]
				out = key.NewPriv()
				out.SetKey(r.DB.BC.Decrypt(&privB), r.DB.BC.Decrypt(&pubB))
				out.WithBC(r.DB.BC)
			} else {
				privB := v[:32]
				pubB := v[32:]
				out = key.NewPriv()
				out.SetKey(&privB, &pubB)
			}
		}
	}
	if !r.OK() {
		out = &key.Priv{}
		out.SetStatus(r.Error())
	} else {
		if !allowReuse {
			if !r.SetStatusIf(txn.Delete(k)).OK() {
				fmt.Println("ERROR", r.Error())
			}
			r.DB.WriteKey(out)
			delete(r.KeyPool.Pool, lowest)
		}
	}
	r.KeyPool.Size--
	return
}

// GetKeyPoolSize returns the number of keys in the keypool
func (r *Wallet) GetKeyPoolSize() int { return r.KeyPool.Size }

// GetOldestKeyPoolTime -
func (r *Wallet) GetOldestKeyPoolTime() int64 { return 0 }

// TopUpKeyPool -
func (r *Wallet) TopUpKeyPool() *Wallet {
	if r.KeyPool.Size < r.KeyPool.Low {
		toAdd := r.KeyPool.High - r.KeyPool.Size
		for i := 0; i < toAdd; i++ {
			var nk *key.Priv
			if r.DB.BC != nil {
				nk = key.NewPriv()
				nk.WithBC(r.DB.BC)
			} else {
				nk = key.NewPriv()
			}
			nk.Make()
			I := []byte(nk.GetID())
			idx := core.Hash64(&I)
			np := &rec.Pool{
				Address: buf.NewByte().Copy(&I).(*buf.Byte),
				Idx:     *idx,
				Seq:     i,
				Priv:    nk.Crypt,
				Pub:     nk.PubKey().(*buf.Byte),
				Created: time.Now().Unix(),
				Expires: time.Now().Add(r.KeyPool.Lifespan).Unix(),
			}
			r.DB.WritePool(np)
			r.KeyPool.Pool[i] = np
			r.KeyPool.Size++
		}
	}
	return r
}

// EmptyKeyPool deletes an entire keypool
func (r *Wallet) EmptyKeyPool() *Wallet {
	if r == nil {
		r = New(nil)
		r.SetStatus(er.NilRec)
		return r
	}
	if r.KeyPool == nil {
		r.KeyPool = &KeyPool{
			High:     100,
			Low:      10,
			Lifespan: 90 * 24 * time.Hour,
		}
		r.SetStatus("keypool not initialised")
		return r
	}
	if r.KeyPool.Size < 1 {
		return r
	}
	opt := badger.DefaultIteratorOptions
	// for i := range r.KeyPool.Pool {
	// 	r.DB.ErasePool(r.KeyPool.Pool[i])
	// 	delete(r.KeyPool.Pool, i)
	// 	r.KeyPool.Size--
	// }
	// And for what was not in memory...
	err := r.DB.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			k := item.Key()
			table := string(k[:8])
			if table == rec.TS["Pool"] {
				fmt.Println("del", hex.EncodeToString(k[8:16]))
				r.SetStatusIf(r.DB.DB.Update(func(txn *badger.Txn) error {
					return txn.Delete(item.Key())
				}))
			}
		}
		return nil
	})
	if !r.SetStatusIf(err).OK() {
		fmt.Println("\nERROR:", r.Error())
	}
	return r
}
