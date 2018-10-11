package sync

import (
	"encoding/hex"
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
	r = new(Node)
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

	err = fallocate.Fallocate(blockchain, 0, 256*256*256*256)
	if err != nil {
		fmt.Println("allocating blockchain file", err.Error())
		os.Exit(1)
	}
	r.Chain = blockchain
	return
}

// Close shuts down the blockchain sync server
func (r *Node) Close() *Node {
	r.DB.Close()
	r.Chain.Close()
	return r
}

// Sync updates the blockchain to the latest current available block height
func (r *Node) Sync() *Node {
	var hashes [][]byte
	var position uint64
	limit := r.GetBestBlockHeight()
	start, _ := r.GetLatestSynced()
	var value []byte
	if start != 0 {
		err := r.DB.View(func(txn *badger.Txn) error {
			var err error
			item, _ := txn.Get(*core.IntToBytes(start))
			value, err = item.Value()
			if !r.SetStatusIf(err).OK() {
				fmt.Println(err.Error())
				return err
			}
			return nil
		})
		if !r.SetStatusIf(err).OK() {
			fmt.Println(err.Error())
			return r
		} else {
			var start int64
			var length uint32
			sB := append(value[32:38], []byte{0, 0}...)
			core.BytesToInt(&start, &sB)
			lB := append(value[38:41], []byte{0}...)
			core.BytesToInt(&length, &lB)
			start += int64(length)
			r.Chain.Seek(start, 0)
			fmt.Println("seeking to end of last record at position", start)
		}

	}

	for i := start; i < limit; i++ {
		getHash, err := r.RPC.Call("getblockhash", []uint32{i})
		if err != nil {
			hashes = append(hashes, []byte{})
		} else {
			unquoted, _ := strconv.Unquote(string(getHash.Result))
			Hash, _ := hex.DecodeString(string(unquoted))
			hashes = append(hashes, Hash)
			resp, _ := r.RPC.Call("getblock", []interface{}{hex.EncodeToString(hashes[i]), false})
			unquoted, _ = strconv.Unquote(string(resp.Result))
			bytes, _ := hex.DecodeString(unquoted)
			begin := position
			length := uint32(len(bytes))
			position += uint64(length)
			r.Chain.Write(bytes)

			k1 := removeTrailingZeroes(*core.IntToBytes(i))
			// fmt.Print("K height ", hex.EncodeToString(k1))
			v1 := (*core.IntToBytes(length))[:3]
			// fmt.Print(" V length ", hex.EncodeToString(v1))
			b := removeTrailingZeroes((*core.IntToBytes(begin))[:6])
			b = append([]byte{byte(len(b))}, b...)
			// fmt.Print(" begin ", hex.EncodeToString(b))
			v1 = append(v1, 0)
			v1 = append(v1, b...)
			H := removeLeadingZeroes(Hash)
			// fmt.Print(" hash ", hex.EncodeToString(H))
			v1 = append(v1, H...)

			// fmt.Println()
			h := core.Hash64(&Hash)
			// fmt.Print("K hash ", hex.EncodeToString(*h))
			k2 := *h
			v2 := removeTrailingZeroes(*core.IntToBytes(i))
			// fmt.Println(" V height", hex.EncodeToString(v2))
			err := r.DB.Update(func(txn *badger.Txn) error {
				return txn.Set(k1, v1)
			})
			if err != nil {
				fmt.Println("writing kv 1", err.Error())
				os.Exit(1)
			}
			err = r.DB.View(func(txn *badger.Txn) error {
				item, _ := txn.Get(k2)
				if item != nil {
					v, err := item.Value()
					if err != nil {
						fmt.Println("reading colliding hash", err.Error())
						os.Exit(1)
					}
					fmt.Println("COLLISION ON HASH64 OF BLOCKHASH", k2)
					v2 = append(v, v2...)
				}
				return nil
			})

			err = r.DB.Update(func(txn *badger.Txn) error {
				return txn.Set(k2, v2)
			})
			if err != nil {
				fmt.Println("writing kv 2", err.Error())
				os.Exit(1)
			}
			// fmt.Println()
		}
		if i%1000 == 0 {
			fmt.Print(i, "...")
		}
	}
	return r
}
