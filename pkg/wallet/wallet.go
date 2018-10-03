package wallet

import (
	"time"

	"github.com/parallelcointeam/duo/pkg/wallet/db"
)

// New returns a new Wallet
func New(newWDB *db.DB) *Wallet {
	w := &Wallet{
		DB:              newWDB,
		version:         FeatureBase,
		maxVersion:      FeatureBase,
		FileBacked:      false,
		OrderPosNext:    0,
		KeyPoolHigh:     100,
		KeyPoolLow:      10,
		KeyPoolLifespan: 90 * 24 * time.Hour,
	}
	return w
}
