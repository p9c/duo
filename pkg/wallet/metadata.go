package wallet
import (
)
// NewKeyMetadata makes a new KeyMetadata structure
func NewKeyMetadata(createTime int64) (M *KeyMetadata) {
	M.Version = CurrentVersion
	M.CreateTime = createTime
	return
}
