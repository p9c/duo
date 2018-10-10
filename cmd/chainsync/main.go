package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/detailyang/go-fallocate"
	"github.com/dgraph-io/badger"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/rpc"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
)

func main() {
	C := rpc.NewClient("127.0.0.1", 11048, "user", "pa55word", false)
	resp, err := C.Call("getinfo", nil)
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

	var path string
	if path, err = homedir.Dir(); err != nil {
		fmt.Println("getting homedir", err.Error())
		os.Exit(1)
	}

	blockstoreBaseDir := path + "/" + db.DefaultBaseDir

	if _, err := os.Stat(blockstoreBaseDir + "/blocks"); !os.IsNotExist(err) {
		// Find current position and update from there
		dbOptions := &badger.DefaultOptions
		dbOptions.Dir = blockstoreBaseDir + "/blockchain"
		dbOptions.ValueDir = dbOptions.Dir + "/values"
		db, err := badger.Open(*dbOptions)
		if err != nil {
			fmt.Println("opening db", err.Error())
			os.Exit(1)
		}
		defer db.Close()
		// Key with content "latest" has value of the head block at last update
		var latest uint64
		var latesthash []byte
		var counter uint64
		err = db.View(func(txn *badger.Txn) error {
			// item, er := txn.Get([]byte("latest"))
			// if item == nil {
			// no latest key found, search the index for the latest and update it
			fmt.Println("counter:", counter)
			counter++
			opt := badger.DefaultIteratorOptions
			opt.PrefetchValues = false
			iter := txn.NewIterator(opt)
			defer iter.Close()
			var key []byte
			var hH, Hh, hp int
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
				switch key[0] {
				case 1:
					hH++
					// table = key[0]
					h := append(key[1:8], 0)
					core.BytesToInt(&height, &h)
					// fmt.Printf("height->hash height %09d hash %s\n",
					// height,
					// hex.EncodeToString(value))
					if height > latest {
						latest = height
						latesthash = value
					}
				case 2:
					Hh++
					// table = key[0]
				case 4:
					hp++
					// table = key[0]
					h := append(key[1:8], 0)
					core.BytesToInt(&height, &h)
					st := value[:8]
					le := value[8:]
					var start uint64
					var length uint32
					core.BytesToInt(&start, &st)
					core.BytesToInt(&length, &le)
					// fmt.Printf("height->position height %09d start %09d length %09d\n",
					// 	height,
					// 	start,
					// 	length)
				}
			}
			fmt.Println(hH, Hh, hp)
			err = db.Update(func(txn *badger.Txn) error {
				v := append(*core.IntToBytes(latest), latesthash...)
				fmt.Println("latest", *core.IntToBytes(latest), latesthash)
				return txn.Set([]byte("latest"), v)
			})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			// } else {

			// 	V, er := item.Value()
			// 	if er != nil {
			// 		fmt.Println(err.Error())
			// 		return er
			// 	}

			// 	latestB := V[:8]
			// 	core.BytesToInt(&latest, &latestB)
			// 	latesthash = V[8:]
			// }
			fmt.Println("updated until block height", latest, "hash", hex.EncodeToString(latesthash))
			return nil
		})
		if err != nil {
			fmt.Println("finding latest block", err.Error())
		}
	} else {
		// Sync from zero
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

		dbOptions := &badger.DefaultOptions
		dbOptions.Dir = blockstoreBaseDir + "/blockchain"
		dbOptions.ValueDir = dbOptions.Dir + "/values"
		db, err := badger.Open(*dbOptions)
		if err != nil {
			fmt.Println("opening db", err.Error())
			os.Exit(1)
		}
		defer db.Close()

		var hashes [][]byte
		var position uint64
		limit := str.Blocks
		// limit := uint64(1000)
		for i := uint64(0); i < limit; i++ {
			getHash, err := C.Call("getblockhash", []uint64{i})
			if err != nil {
				hashes = append(hashes, []byte{})
			} else {
				unquoted, _ := strconv.Unquote(string(getHash.Result))
				bytes, _ := hex.DecodeString(string(unquoted))
				hashes = append(hashes, bytes)
				resp, _ = C.Call("getblock", []interface{}{hex.EncodeToString(hashes[i]), false})
				unquoted, _ = strconv.Unquote(string(resp.Result))
				bytes, _ = hex.DecodeString(unquoted)
				start := position
				length := uint32(len(bytes))
				position += uint64(length)
				blockchain.Write(bytes)

				k1 := *core.IntToBytes(1 + i<<8)
				v1 := hashes[i]
				k2 := append([]byte{2}, hashes[i]...)
				v2 := *core.IntToBytes(i)
				k3 := *core.IntToBytes(4 + i<<8)
				v3 := append(*core.IntToBytes(start), *core.IntToBytes(length)...)
				err := db.Update(func(txn *badger.Txn) error {
					return txn.Set(k1, v1)
				})
				if err != nil {
					fmt.Println("writing kv 1", err.Error())
					os.Exit(1)
				}
				err = db.Update(func(txn *badger.Txn) error {
					return txn.Set(k2, v2)
				})
				if err != nil {
					fmt.Println("writing kv 2", err.Error())
					os.Exit(1)
				}
				err = db.Update(func(txn *badger.Txn) error {
					return txn.Set(k3, v3)
				})
				if err != nil {
					fmt.Println("writing kv 3", err.Error())
					os.Exit(1)
				}
			}
			if i%1000 == 0 {
				fmt.Print(i, "...")
			}
		}
	}
	fmt.Println()
}
