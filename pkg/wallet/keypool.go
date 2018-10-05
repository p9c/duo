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
	if r.KeyPool != nil {
		r.EmptyKeyPool()
	}
	r.KeyPool = make(KeyPool)
	for i := 0; i < r.KeyPoolHigh; i++ {
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
			Expires: time.Now().Add(r.KeyPoolLifespan).Unix(),
		}
		r.DB.WritePool(np)
		r.KeyPool[i] = np
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
	opt := badger.DefaultIteratorOptions
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
