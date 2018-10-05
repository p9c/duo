package wallet

import (
	"fmt"
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
		nk := key.NewPriv()
		if r.DB.BC != nil {
			nk.WithBC(r.DB.BC)
		}
		nk.Make()
		idx := core.Hash64(nk.PubKey().Bytes())
		np := &rec.Pool{
			Address: buf.NewByte().Copy(nk.PubKey().Bytes()).(*buf.Byte),
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
	r.KeyPool = new(KeyPool)
	r.KeyPool.Pool = make(PoolMap)
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
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
						address = k[16:65] // 49 bytes
						address = *r.DB.BC.Decrypt(&address)
						creB = k[65:89] // 24 bytes
						creB = *r.DB.BC.Decrypt(&creB)
						expB = k[89:] // 24 bytes
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
			}
		}
		iter.Close()
		return nil
	})
	if !r.SetStatusIf(err).OK() {
		fmt.Println("\nERROR:", r.Error())
	}
	return r
}

// AddReserveKey -
func (r *Wallet) AddReserveKey(kp *rec.Pool) *Wallet { return r }

// GetKeyFromPool -
func (r *Wallet) GetKeyFromPool(*key.Pub, bool) *Wallet {
	return r
}

// GetKeyPoolSize -
func (r *Wallet) GetKeyPoolSize() int { return 0 }

// GetOldestKeyPoolTime -
func (r *Wallet) GetOldestKeyPoolTime() int64 { return 0 }

// ReserveKeyFromKeyPool -
func (r *Wallet) ReserveKeyFromKeyPool(int64, *rec.Pool) {}

// TopUpKeyPool -
func (r *Wallet) TopUpKeyPool() *Wallet { return r }

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
	for i := range r.KeyPool.Pool {
		r.DB.ErasePool(r.KeyPool.Pool[i])
	}
	// And for what was not in memory...
	err := r.DB.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			k := item.Key()
			table := string(k[:8])
			if table == rec.TS["Pool"] {
				if !r.SetStatusIf(txn.Delete(k)).OK() {
					fmt.Println("\nERROR", r.Error())
				}
			}
		}
		iter.Close()
		return nil
	})
	if !r.SetStatusIf(err).OK() {
		fmt.Println("\nERROR:", r.Error())
	}
	return r
}
