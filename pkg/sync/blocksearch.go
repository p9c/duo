package sync

import (
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
)

// GetBlockHash returns the block hash from the chainsync database
func (r *Node) GetBlockHash(height uint32) (out []byte) {
	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
		k := append([]byte{1}, removeTrailingZeroes(*core.IntToBytes(height))...)
		item, err := txn.Get(k)
		v, err := item.Value()
		out = append(make([]byte, 32-len(v)), v...)
		return err
	}))
	return
}

// GetHeightFromHash returns the height of a block with a given hash
func (r *Node) GetHeightFromHash(h []byte) (out uint32) {
	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
		H := *core.Hash64(&h)
		item, err := txn.Get(append([]byte{2}, H...))
		if err != nil {
			v, err := item.Value()
			if item != nil && err == nil {
				core.BytesToInt(&out, &v)
			}
		}
		return err
	}))

	return ^uint32(0)
}
