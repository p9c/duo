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

	blockchain, err := os.Create(path + "/" + db.DefaultBaseDir + "/blocks")
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
	dbOptions.Dir = path + "/" + db.DefaultBaseDir + "/blockchain"
	dbOptions.ValueDir = dbOptions.Dir + "/values"
	db, err := badger.Open(*dbOptions)
	if err != nil {
		fmt.Println("opening db", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	var hashes [][]byte
	var position uint64

	for i := uint64(0); i < str.Blocks; i++ {
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

			// fmt.Println("k 01", fmt.Sprintf("%06x", i), "v", hex.EncodeToString(hashes[i]))
			// fmt.Println("k 02", hex.EncodeToString(hashes[i]), "v", fmt.Sprintf("%08x", i))
			// fmt.Println("k 04", fmt.Sprintf("%06x v %08x %08x", i, start, end))
			// fmt.Println()

			k1 := *core.IntToBytes(0x01000000 & i)
			v1 := hashes[i]
			k2 := append([]byte{2}, hashes[i]...)
			v2 := *core.IntToBytes(i)
			k3 := *core.IntToBytes(0x04000000 & i)
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
	fmt.Println()
}
