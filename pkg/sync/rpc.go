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
