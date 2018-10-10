package sync

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	fallocate "github.com/detailyang/go-fallocate"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/rpc"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
)

// NewNode creates a new blockchain sync node/server
func NewNode() (r *Node) {
	r.RPC = rpc.NewClient("127.0.0.1", 11048, "user", "pa55word", false)
	var path string
	var err error
	if path, err = homedir.Dir(); err != nil {
		fmt.Println("getting homedir", err.Error())
		os.Exit(1)
	}

	blockstoreBaseDir := path + "/" + db.DefaultBaseDir

	if _, err := os.Stat(blockstoreBaseDir + "/blocks"); !os.IsNotExist(err) {
		dbOptions := &badger.DefaultOptions
		dbOptions.Dir = blockstoreBaseDir + "/blockchain"
		dbOptions.ValueDir = dbOptions.Dir + "/values"
		db, err := badger.Open(*dbOptions)
		if err != nil {
			fmt.Println("opening db", err.Error())
			os.Exit(1)
		}
		r.DB = db
	}

	blockchain, err := os.Create(blockstoreBaseDir + "/blocks")
	if err != nil {
		fmt.Println("creating blockchain file", err.Error())
		os.Exit(1)
	}
	defer blockchain.Close()

	err = fallocate.Fallocate(blockchain, 0, 256*256*256*256)
	if err != nil {
		fmt.Println("allocating blockchain file", err.Error())
		os.Exit(1)
	}

	return
}

// Close shuts down the blockchain sync server
func (r *Node) Close() *Node {
	r.DB.Close()
	return r
}

// GetBestBlockHeight returns the newest consensus block height
func (r *Node) GetBestBlockHeight() uint64 {
	resp, err := r.RPC.Call("getinfo", nil)
	if err != nil {
		fmt.Println("rpc getinfo", err.Error())
		os.Exit(1)
	}
	str := new(rpc.GetInfo)
	err = json.Unmarshal(resp.Result, str)
	if err != nil {
		fmt.Println("unmarshalling info", err.Error())
		os.Exit(1)
	}
	return str.Blocks
}

// GetRawBlock gets the raw block given a block height
func (r *Node) GetRawBlock(height uint64) *[]byte {
	getHash, err := r.RPC.Call("getblockhash", []uint64{height})
	if err != nil {
	} else {
		unquoted, _ := strconv.Unquote(string(getHash.Result))
		bytes, _ := hex.DecodeString(string(unquoted))
		resp, err := r.RPC.Call("getblock", []interface{}{hex.EncodeToString(bytes), false})
		if err != nil {
		} else {
			unquoted, _ = strconv.Unquote(string(resp.Result))
			bytes, _ = hex.DecodeString(unquoted)
		}
		return &bytes
	}
	return &[]byte{}
}

// GetLatestSynced returns the newest block height stored in the database, updates it if it wasn't already stored
func (r *Node) GetLatestSynced() (latest uint64, latesthash []byte) {
	err := r.DB.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte("latest"))
		if item == nil {
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
				// var table byte
				var height uint64
				if key[0] == 1 {
					h := append(key[1:8], 0)
					core.BytesToInt(&height, &h)
					if height > latest {
						latest = height
						latesthash = value
					}
				}
			}
			err := r.DB.Update(func(txn *badger.Txn) error {
				k := append(*core.IntToBytes(latest), latesthash...)
				fmt.Println("latest", *core.IntToBytes(latest), latesthash)
				return txn.Set([]byte("latest"), k)
			})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			fmt.Println("updated until block height", latest, "hash", hex.EncodeToString(latesthash))
			return nil
		}
		v, err := item.Value()
		if err != nil {
			return err
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
	return
}

// Sync updates the blockchain to the latest current available block height
func (r *Node) Sync() *Node {

	return r
}
