package wallet

import (
	"github.com/parallelcointeam/duo/pkg/key"
)

// KeyMetadata is a data structure for storing metadata related to a key pair
type KeyMetadata struct {
	Pub        *key.Pub
	Version    uint32
	CreateTime int64
}

// NewKeyMetadata makes a new KeyMetadata structure
func NewKeyMetadata(createTime int64) (M *KeyMetadata) {
	M.Version = CurrentVersion
	M.CreateTime = createTime
	return
}
