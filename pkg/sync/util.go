package sync

import (
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
)

func removeTrailingZeroes(in []byte) []byte {
	length := 0
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] != 0 {
			length = i
			break
		}
	}
	return in[:length+1]
}

func removeLeadingZeroes(in []byte) []byte {
	nonzerostart := 0
	for i := range in {
		if in[i] != 0 {
			nonzerostart = i
			break
		}
	}
	return in[nonzerostart:]
}

func (r *Node) getLatest() (h uint32) {
	var latestB []byte

	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("latest"))
		if err == nil {
			latestB, err = item.Value()
		}
		return err
	}))
	if latestB != nil {
		heightB := latestB[:4]
		core.BytesToInt(&h, &heightB)
	}
	r.UnsetStatus()
	return
}
