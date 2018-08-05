package wallet

import (
)

// KeyMetadata is information about the key
type KeyMetadata struct {
	Version    uint32
	CreateTime int64
}

type keymetadata interface {
	FromBytes([]byte)
	ToBytes() []byte
}

// NewKeyMetadata makes a new KeyMetadata structure
func NewKeyMetadata(createTime int64) (M *KeyMetadata) {
	M.Version = CurrentVersion
	M.CreateTime = createTime
	return
}
