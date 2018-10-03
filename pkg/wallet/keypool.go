package wallet

import (
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// NewKeyPool creates a new pool of keys in reserve for generating transactions
func (r *Wallet) NewKeyPool() *Wallet {

	r.KeyPool = make(KeyPool)
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

// Destroy deletes an entire keypool
func (r *Wallet) Destroy() *Wallet {
	return r
}
