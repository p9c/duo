package block

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/sync"
)

func TestGetRawBlock(t *testing.T) {
	node := sync.NewNode()
	best := node.LegacyGetBestBlockHeight()
	var B uint32
	b := make([]byte, 4)
	for i := 0; i < 10; i++ {
		for {
			rand.Read(b)
			core.BytesToInt(&B, &b)
			// fmt.Println(B, b)
			if B < best {
				break
			}
		}
		rb := node.GetRawBlock(uint64(B))
		r := *rb

		fmt.Println(B)
		fmt.Println(hex.EncodeToString(rev(r)))
		fmt.Println()
		fmt.Println("version  prevblockhash                                                    merklehashroot                                                   time     bits     nonce")

		fmt.Print(hex.EncodeToString(rev(r[:4])), " ")
		fmt.Print(hex.EncodeToString(rev(r[4:36])), " ")
		fmt.Print(hex.EncodeToString(rev(r[36:68])), " ")
		fmt.Print(hex.EncodeToString(rev(r[68:72])), " ")
		fmt.Print(hex.EncodeToString(rev(r[72:76])), " ")
		fmt.Println(hex.EncodeToString(rev(r[76:80])), " ")

		fmt.Println(r[81], "transactions in block")
		var txcount uint64
		txcount = 0
		bitlen := 1
		switch {
		case r[81] < 0xFD:
			txcount = uint64(r[81])
			bitlen = 1
			// fmt.Println("txcount", txcount, "bitlen", bitlen)
			r = r[82:]
		case r[81] == 0xFD:
			t := r[82:84]
			core.BytesToInt(&txcount, &t)
			// txcount += uint64(r[81])
			bitlen = 2
			r = r[84:]
		case r[81] == 0xFE:
			t := r[82:86]
			core.BytesToInt(&txcount, &t)
			// txcount += uint64(r[81])
			bitlen = 4
			r = r[86:]
		case r[81] == 0xFF:
			t := r[82:90]
			core.BytesToInt(&txcount, &t)
			// txcount += uint64(r[81])
			bitlen = 8
			r = r[90:]
		}
		_ = bitlen

		// fmt.Println("00000001 0000000000000000000000000000000000000000000000000000000000000000 ffffffff"
		fmt.Println("version  prevtxhash                                                       outindex")
		fmt.Print(hex.EncodeToString(r[:4]), " ")

		fmt.Print(hex.EncodeToString(r[4:36]), " ")

		fmt.Print(hex.EncodeToString(r[36:40]), " ")

		switch {
		case r[41] < 0xFD:
			txcount = uint64(r[41])
			bitlen = 1
			// fmt.Println("txcount", txcount, "bitlen", bitlen)
			r = r[42:]
		case r[41] == 0xFD:
			t := r[42:44]
			core.BytesToInt(&txcount, &t)
			// txcount += uint64(r[81])
			bitlen = 2
			r = r[44:]
		case r[41] == 0xFE:
			t := r[42:46]
			core.BytesToInt(&txcount, &t)
			// txcount += uint64(r[81])
			bitlen = 4
			r = r[46:]
		case r[41] == 0xFF:
			t := r[42:50]
			core.BytesToInt(&txcount, &t)
			// txcount += uint64(r[81])
			bitlen = 8
			r = r[50:]
		}
		fmt.Println()

		// fmt.Println("value")
		// fmt.Print(hex.EncodeToString(r[:8]), " ")
		// var x uint64
		// xx := r[:8]
		// core.BytesToInt(&x, &xx)
		// // fmt.Print(hex.EncodeToString(r), " ")
		// fmt.Println()

		var value interface{}
		r, value = sync.ExtractVarint(uint64(0), r)
		fmt.Println("value", value.(uint64))

		fmt.Println("\nRest:\n", hex.EncodeToString(r))
		fmt.Println()
	}
}

func rev(in []byte) (out []byte) {
	out = make([]byte, len(in))
	for i := range in {
		out[len(in)-i-1] = in[i]
	}
	return
}
