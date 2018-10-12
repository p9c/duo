package sync

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
)

// CheckConsistency walks through the database and checks that blocks make the same hash as the index and when it fails it moves the 'latest' tag to the previous to failed
func (r *Node) CheckConsistency() *Node {
	var latest uint32
	var latesthash, v []byte // k
	// var end uint64
	err := r.DB.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.PrefetchValues = false
		iter := txn.NewIterator(opt)
		defer iter.Close()
		var key []byte
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key = item.Key()
			value, err := item.Value()
			if err != nil {
				fmt.Println("reading key/value pairs", err.Error())
				os.Exit(1)
			}
			var height uint32
			if key[0] == 1 {
				h := append(key[1:8], 0)
				core.BytesToInt(&height, &h)
				if height > latest {
					latest = height
					latesthash = value
				}
			}
			err = r.DB.Update(func(txn *badger.Txn) error {
				k := append(*core.IntToBytes(latest), latesthash...)
				fmt.Println("latest", *core.IntToBytes(latest), latesthash)
				return txn.Set([]byte("latest"), k)
			})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			fmt.Println("updated until block height", latest, "hash", hex.EncodeToString(latesthash))
			// k = item.Key()[1:]
			// v, err := item.Value()
			return nil
		}
		if v != nil {
			heightB := v[:8]
			core.BytesToInt(&latest, &heightB)
			latesthash = v[8:]
		}
		return nil
	})
	if err != nil {
		fmt.Println("finding latest block", err.Error())
	}
	r.Latest = latest
	r.LatestHash = latesthash
	return r
}
