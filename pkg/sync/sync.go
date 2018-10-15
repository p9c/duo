package sync

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/anaskhan96/base58check"
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
	// dbOptions.ValueThreshold = 40960
	// dbOptions.NumVersionsToKeep = 1
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
	startHeight := r.getLatest()
	// If we got a latest height we are assuming that the database is consistent up to this point. If we find errors or just want to recheck we can just delete the latest key and run this function and it will start from zero

	bestBlockHeight := r.LegacyGetBestBlockHeight()

	var v []byte
	var lastBlockUpdated uint32
	for i := startHeight; i <= bestBlockHeight; i++ {
		foundtx := false
		foundrepeat := false
		h := r.LegacyGetBlockHash(i)
		k1, v1 := EncodeKV(Block{Height: i, Hash: h})
		k2, v2 := EncodeKV(Hash{
			HHash:  *core.Hash64(&h),
			Height: i,
		})

		// Set the keys for the forward and reverse height/hash lookup table
		r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
			err := txn.SetWithDiscard(k1, v1, 0)
			if err != nil {
				return err
			}
			err = txn.SetWithDiscard(k2, v2, 0)
			return err
		}))

		// Update the address index, append the new reference if the record already exists
		resp, err := r.RPC.Call("getblock", []interface{}{hex.EncodeToString(h), true})
		if !r.SetStatusIf(err).OK() {
			fmt.Println("getting block", r.Error())
		}
		var blk rpc.GetBlock
		if !r.SetStatusIf(json.Unmarshal(resp.Result, &blk)).OK() {
			fmt.Println("unmarshalling block json", r.Error())
			return r
		}

		hugeaddr := false
		for j := range blk.Tx {
			resp, err := r.RPC.Call("getrawtransaction", []interface{}{blk.Tx[j], 1})
			if err == nil {
				foundtx = true
				// fmt.Println("tx", j)
				var tx rpc.RawTransaction
				err = json.Unmarshal(resp.Result, &tx)
				if err != nil {
					fmt.Println("unmarshalling transaction", err)
				} else {
					for k := range tx.Vout {
						for l := range tx.Vout[k].ScriptPubKey.Addresses {
							addr := tx.Vout[k].ScriptPubKey.Addresses[l]
							id, _ := base58check.Decode(addr)
							I, _ := hex.DecodeString(id[2:])
							hhash := *core.Hash64(&I)
							height := make([]byte, 5)
							l := binary.PutUvarint(height, uint64(i))
							height = height[:l]

							pos := make([]byte, 3)
							l = binary.PutUvarint(pos, uint64(k))
							pos = pos[:l]

							k3 := append([]byte{16}, hhash...)
							v3 := append(pos, height...)

							r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
								item, err := txn.Get(append([]byte{16}, hhash...))
								if item != nil {
									existing, err := item.ValueCopy(nil)
									if err == nil {
										foundrepeat = true

										V := append(existing, v3...)
										v3 = V
										if len(v3) > 8192 {
											hugeaddr = true
										}
										// fmt.Println("\nrepeat",)
										return nil
									}
								}
								return err
							}))
							r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
								err = txn.Delete(k3)
								return err
							}))
							r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
								// v = pruneToHeight(v3, i)
								err = txn.SetWithDiscard(k3, v3, 0)
								return err
							}))
						}
					}
				}
			}
		}

		lastBlockUpdated = i
		// update the latest
		v = append(*core.IntToBytes(lastBlockUpdated),
			r.LegacyGetBlockHash(lastBlockUpdated)...)
		r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
			err = txn.SetWithDiscard([]byte("latest"), v, 0)
			return err
		}))
		if i&16384 == 0 {
			r.DB.RunValueLogGC(0.5)
		}
		if i%288 == 0 {
			fmt.Println()
		}
		if i%72 == 0 {
			fmt.Printf("\n%7d ", i)
		} else {
			switch {
			case hugeaddr:
				fmt.Print("!")
			case !foundtx && !foundrepeat:
				fmt.Print(".")
			case foundtx && !foundrepeat:
				fmt.Print("*")
			case foundtx && foundrepeat:
				fmt.Print("+")
			}
		}
		if !r.OK() {
			break
		}
	}
	fmt.Println("done")

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
	fmt.Println("\nRemoving old versions of records")
	var counter int
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	err := r.DB.Update(func(txn *badger.Txn) error {
		iter := txn.NewIterator(opt)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			if item.DiscardEarlierVersions() {
				fmt.Print("%")
				counter++
			}
		}
		fmt.Printf("\n%d old records flushed\n", counter)
		return nil
	})
	r.DB.RunValueLogGC(0.5)
	if !r.SetStatusIf(err).OK() {
		fmt.Println("\nERROR:", r.Error())
	}
	return r
}
