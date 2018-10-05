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
			Created: time.Now().Unix(),
			Expires: time.Now().Add(r.KeyPool.Lifespan).Unix(),
		}
		r.DB.WritePool(np)
		r.KeyPool.Pool[i] = np
		r.KeyPool.Size++
	}
	return r
}

// LoadKeyPool loads the keypool into memory from the database
func (r *Wallet) LoadKeyPool() *Wallet {
	opt := badger.DefaultIteratorOptions
	err := r.DB.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			k := item.Key()
			table := string(k[:8])
			if table == rec.TS["Pool"] {

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
func (r *Wallet) GetKeyFromPool(*key.Pub, bool) *Wallet { return r }

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
