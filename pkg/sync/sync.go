package sync

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/rpc"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
)

// NewNode creates a new blockchain sync node/server
func NewNode() (r *Node) {
	r = new(Node)
	r.RPC = rpc.NewClient("127.0.0.1", 11048, "user", "pa55word", false)

	var path string
	var err error
	if path, err = homedir.Dir(); err != nil {
		fmt.Println("getting homedir", err.Error())
		os.Exit(1)
	}

	blockstoreBaseDir := path + "/" + db.DefaultBaseDir

	dbOptions := &badger.DefaultOptions
	dbOptions.Dir = blockstoreBaseDir + "/index"
	dbOptions.ValueDir = dbOptions.Dir + "/values"
	db, err := badger.Open(*dbOptions)
	if err != nil {
		fmt.Println("opening db", err.Error())
		os.Exit(1)
	}
	r.DB = db

	return
}

// Close shuts down the blockchain sync server
func (r *Node) Close() *Node {
	r.DB.Close()
	return r
}

// Sync updates the blockchain to the latest current available block height
func (r *Node) Sync() *Node {
	var startHeight uint32
	var latestB []byte

	// var startBlockHash []byte

	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("latest"))
		if err == nil {
			latestB, err = item.Value()
		}
		return err
	}))
	if latestB != nil {
		heightB := latestB[:4]
		core.BytesToInt(&startHeight, &heightB)

		// we probably don't need the hash but it's there for other reasons
		// hashB := latestB[4:]
		// startBlockHash = append(make([]byte, 32-len(hashB)), hashB...)
	}

	// If we got a latest height we are assuming that the database is consistent up to this point. If we find errors or just want to recheck we can just delete the latest key and run this function and it will start from zero

	bestBlockHeight := r.LegacyGetBestBlockHeight()

	var v []byte
	k := []byte("latest")
	var lastBlockUpdated uint32
	// bestBlockHeight = 50
	for i := startHeight; i <= bestBlockHeight; i++ {
		k1, v1 := EncodeKV(Block{Height: i, Hash: r.LegacyGetBlockHash(i)})
		h := r.LegacyGetBlockHash(i)
		k2, v2 := EncodeKV(Hash{
			HHash:  *core.Hash64(&h),
			Height: i,
		})
		// fmt.Println(k2, v2)
		r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
			err := txn.Set(k1, v1)
			if err != nil {
				return err
			}
			err = txn.Set(k2, v2)
			return err
		}))
		lastBlockUpdated = i
		if i%1000 == 0 {
			// update the latest
			r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
				v = append(*core.IntToBytes(lastBlockUpdated),
					r.LegacyGetBlockHash(lastBlockUpdated)...)
				err := txn.Set(k, v)
				return err
			}))
			fmt.Print(lastBlockUpdated, "...")
		}
		if !r.OK() {
			break
		}
	}
	fmt.Println("done")

	// update the latest
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		v = append(*core.IntToBytes(lastBlockUpdated),
			r.LegacyGetBlockHash(lastBlockUpdated)...)
		err := txn.Set(k, v)
		return err
	}))
	if r.OK() {
		r.Latest = lastBlockUpdated
		hS := string(r.LegacyGetBlockHash(lastBlockUpdated))
		hS, _ = strconv.Unquote(hS)
		r.LatestHash, _ = hex.DecodeString(hS)
	}
	return r
}

// GetLatestSynced returns the newest block height stored in the database, updates it if it wasn't already stored
func (r *Node) GetLatestSynced() (latest uint32, latesthash []byte) {
	if r.Latest != 0 && r.LatestHash != nil {
		return r.Latest, r.LatestHash
	}
	return
}

// RemoveOldVersions walks the database and removes old versions of records
func (r *Node) RemoveOldVersions() *Node {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			item.DiscardEarlierVersions()
			// if result {
			// 	fmt.Println(item.Value())
			// }
		}
		return nil
	})
	if !r.SetStatusIf(err).OK() {
		fmt.Println("\nERROR:", r.Error())
	}
	return r
}
