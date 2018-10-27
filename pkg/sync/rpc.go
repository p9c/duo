package sync

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/parallelcointeam/duo/pkg/rpc"
)

// LegacyGetBestBlockHeight returns the newest consensus block height and
func (r *Node) LegacyGetBestBlockHeight() uint32 {
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
	r.Best = str.Blocks
	r.BestTime = time.Now().Unix()
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

// LegacyGetBlockHash gets the block hash using an external parallelcoind full node
func (r *Node) LegacyGetBlockHash(height uint32) (blockHash []byte) {
	getHash, err := r.RPC.Call("getblockhash", []uint64{uint64(height)})
	if err != nil {
	} else {
		unquoted, _ := strconv.Unquote(string(getHash.Result))
		blockHash, _ = hex.DecodeString(unquoted)
	}
	return
}

// GetTxValue gets the transaction based on a given tx hash
func (r *Node) GetTxValue(txhash []byte) (out float64) {

	intVerbose := 1
	txS := hex.EncodeToString(txhash)
	// fmt.Println("getrawtransaction", txS, 1)
	txB, err := r.RPC.Call("getrawtransaction", []interface{}{txS, intVerbose})

	// vb := 1
	// txB, err := r.RPC.Call("getrawtransaction", []string{hex.EncodeToString(txhash)})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var I interface{}
		json.Unmarshal(txB.Result, &I)
		ii := I.(map[string]interface{})
		// fmt.Println(ii["vout"])
		iii := ii["vout"].([]interface{})[0].(map[string]interface{})
		// j, _ := json.MarshalIndent(iii["value"], "", "  ")
		// fmt.Println(uint64(iii["value"].(float64) * float64(core.COIN)))
		out = iii["value"].(float64)
		// .(map[string]interface{})["value"])
		// fmt.Println(string(txB.Result))
	}
	return
}
