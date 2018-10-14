package sync

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/anaskhan96/base58check"
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/rpc"
)

// UpdateAddresses updates the address reference index
func (r *Node) UpdateAddresses() *Node {
	var addressLatest uint32
	latest, _ := r.GetLatestSynced()
	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
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
		h := hex.EncodeToString(r.GetBlockHash(i))
		resp, err := r.RPC.Call("getblock", []interface{}{h, true})
		var blk rpc.GetBlock
		err = json.Unmarshal(resp.Result, &blk)
		if err != nil {
			fmt.Println("ERROR", err.Error())
		}
		for j := range blk.Tx {
			resp, err := r.RPC.Call("getrawtransaction", []interface{}{blk.Tx[j], 1})
			if err == nil {
				var tx rpc.RawTransaction
				err = json.Unmarshal(resp.Result, &tx)
				if err != nil {
					fmt.Println(err)
				} else {
					for k := range tx.Vout {
						for l := range tx.Vout[k].ScriptPubKey.Addresses {
							addr := tx.Vout[k].ScriptPubKey.Addresses[l]
							id, _ := base58check.Decode(addr)
							I, _ := hex.DecodeString(id[2:])
							hhash := *core.Hash64(&I)
							height := *core.IntToBytes(i)
							posI := uint16(k)
							pos := *core.IntToBytes(posI)

							k := append([]byte{16}, hhash...)
							v := append(pos, height...)

							r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
								item, err := txn.Get(append([]byte{16}, hhash...))
								if item != nil {
									existing, err := item.Value()
									if err == nil {
										fmt.Println("reappearance", k, v)
										v = append(existing, v...)
										fmt.Println("new", v)
										return nil
									}
								}
								return err
							}))
							r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
								v = pruneToHeight(v, i)
								fmt.Println("k", k, "v", v)
								err = txn.Set(k, v)
								return err
							}))
						}
					}
				}
			}
		}
		if i%1000 == 0 {
			r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
				return txn.Set([]byte("address"), *core.IntToBytes(i))
			}))
		}
	}
	return r
}

func pruneToHeight(refs []byte, height uint32) []byte {
	if len(refs) < 3 {
		return []byte{}
	}
	var i uint32
	for {
		var t uint32
		n := refs[i+2 : i+6]
		core.BytesToInt(&t, &n)
		if t <= height && len(refs) > (int(i)+1)*6 {
			i++
			continue
		}
		return refs[:6+i*6]
	}
}
