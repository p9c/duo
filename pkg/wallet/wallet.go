package wallet

import (
	"time"

	"github.com/parallelcointeam/duo/pkg/wallet/db"
)

// New returns a new Wallet
func New(newWDB *db.DB) *Wallet {
	w := &Wallet{
		DB:           newWDB,
		version:      FeatureBase,
		maxVersion:   FeatureBase,
		FileBacked:   false,
		OrderPosNext: 0,
		KeyPool: &KeyPool{
			High:     100,
			Low:      10,
			Lifespan: 90 * 24 * time.Hour,
		},
	}
	return w
}
