package wallet

import (
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
	for i := 0; i < r.KeyPoolTarget; i++ {
		nk := key.NewPriv().Make()
		idx := core.Hash64(nk.PubKey().Bytes())
		np := &rec.Pool{
			Idx: idx,
			Seq: int64(i),
			Key: nk,
		}
	}
	return r
}

// AddReserveKey -
func (r *Wallet) AddReserveKey(kp *rec.Pool) *Wallet { return w }

// GetKeyFromPool -
func (r *Wallet) GetKeyFromPool(*key.Pub, bool) *Wallet { return w }

// GetKeyPoolSize -
func (r *Wallet) GetKeyPoolSize() int { return 0 }

// GetOldestKeyPoolTime -
func (r *Wallet) GetOldestKeyPoolTime() int64 { return 0 }

// ReserveKeyFromKeyPool -
func (r *Wallet) ReserveKeyFromKeyPool(int64, *rec.Pool) {}

// TopUpKeyPool -
func (r *Wallet) TopUpKeyPool() *Wallet { return w }

// EmptyKeyPool deletes an entire keypool
func (r *Wallet) EmptyKeyPool() *Wallet {
	return r
}
