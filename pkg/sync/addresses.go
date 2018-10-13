package sync

import (
	"encoding/hex"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
)

// UpdateAddresses updates the address reference index
func (r *Node) UpdateAddresses() *Node {
	var addressLatest uint32
	latest, _ := r.GetLatestSynced()
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("address"))
		if item != nil {
			a, err := item.Value()
			if err != nil {
				return err
			}
			core.BytesToInt(&addressLatest, &a)
		}
		return err
	}))
	for i := addressLatest; i < latest; i++ {
		h := r.GetBlockHash(i)
		fmt.Println(i, hex.EncodeToString(h))
	}
	return r
}
