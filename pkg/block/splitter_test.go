package block

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/sync"
)

/*

Block Header Format:

Version: Little Endian uint32 - 2 = sha256d, 514 = scrypt
HashPrevBlock: Little Endian 32 bytes
HashMerkleRoot: Little Endian 32 bytes
Time: Big Endian 32 bits unix timestamp seconds since Jan 1 1970 00:00 UTC
Bits: Little Endian Current PoW target in compact format
Nonce: Little Endian random 32 bit number incremented for each PoW round

Transactions



*/

func TestGetRawBlock(t *testing.T) {
	node := sync.NewNode()
	best := node.LegacyGetBestBlockHeight()
	var B uint32
	b := make([]byte, 4)
	for i := 0; i < 10; i++ {
		for {
			rand.Read(b)
			core.BytesToInt(&B, &b)
			if B < best {
				break
			}
		}
		rb := node.GetRawBlock(uint64(B))
		r := *rb

		fmt.Println("BLOCK", B)

		var Version uint32
		v := r[:4]
		r = r[4:]
		core.BytesToInt(&Version, &v)
		fmt.Print("Version ", Version)
		switch Version {
		case 2:
			fmt.Println(" : SHA256D PoW")
		case 514:
			fmt.Println(" : Scrypt PoW")
		}

		HashPrevBlock := *rev(r[:32])
		r = r[32:]
		fmt.Println("HashPrevBlock           ", hx(HashPrevBlock))

		HashMerkleRoot := *rev(r[:32])
		r = r[32:]
		fmt.Println("HashMerkleRoot          ", hx(HashMerkleRoot))

		ti := r[:4]
		r = r[4:]
		var t32 int32
		core.BytesToInt(&t32, &ti)
		BlockTime := int64(t32)
		fmt.Println("Unix timestamp          ", hx(ti))
		fmt.Println("                    Time", time.Unix(BlockTime, 0))

		Bits := *rev(r[:4])
		r = r[4:]
		fmt.Println("Bits                    ", hx(Bits))
		coeff := Bits[0]
		base := Bits[1:]
		tail := make([]byte, coeff-3)
		tgt := append(base, tail...)
		Target := append(make([]byte, 32-len(tgt)), tgt...)
		fmt.Println("                  Target", hx(Target))

		nn := r[:4]
		nn = *rev(nn)
		r = r[4:]
		var Nonce uint32
		core.BytesToInt(&Nonce, &nn)
		fmt.Println("Nonce                   ", Nonce)

		var txCount uint64
		var txCountIface interface{}
		r, txCountIface = ExtractCompactInt(txCount, r)
		txCount = txCountIface.(uint64)
		fmt.Println("TxCount                 ", txCount)

		txc := int(txCount)
		for txs := 0; txs < txc; txs++ {

			fmt.Println("Transaction", txs)

			rV := r[:4]
			r = r[4:]
			var tx0v uint32
			core.BytesToInt(&tx0v, &rV)
			fmt.Println("    tx version          ", tx0v)

			// This is present when there is segwit, which there isn't
			// fl := r[:2]
			// fmt.Println(fl)
			// var flg uint16
			// core.BytesToInt(&flg, &fl)
			// fmt.Printf("    flag                 %04x\n", flg)
			// if flg == 1 {
			// 	r = r[2:]
			// }

			var txV uint64
			var txI interface{}
			r, txI = ExtractCompactInt(txV, r)
			txiV := int(txI.(uint64))
			fmt.Println("in-counter              ", txiV)

			for txis := 0; txis < txiV; txis++ {

				fmt.Println("Txin", txis)

				tx1pth := *rev(r[:32])
				r = r[32:]
				fmt.Println("    PrevTxHash          ", hx(tx1pth))

				tx1txi := *rev(r[:4])
				r = r[4:]
				fmt.Printf("    Prev Txout Index     %08x\n", tx1txi)

				var txsl interface{}
				r, txsl = ExtractCompactInt(txV, r)
				txV = txsl.(uint64)
				fmt.Println("    Txin script length  ", txV)

				tx1scr := *rev(r[:txV])
				r = r[txV:]
				fmt.Println("    script              ", hx(tx1scr))

				tx1seq := r[:4]
				r = r[4:]
				fmt.Println("    seq number          ", hx(tx1seq))
			}
			r, txI = ExtractCompactInt(txV, r)
			txoV := int(txI.(uint64))
			fmt.Println("out-counter             ", txoV)

			for txos := 0; txos < txoV; txos++ {

				fmt.Println("Txout", txos)

				tx1val := r[:8]
				r = r[8:]
				var tx1V uint64
				core.BytesToInt(&tx1V, &tx1val)
				fmt.Printf("    value                %4.7f\n", float64(tx1V)/core.COIN)

				r, txI = ExtractCompactInt(txV, r)
				txV = txI.(uint64)
				fmt.Println("    Txout script length ", txV)

				tx1scro := *rev(r[:txV])
				r = r[txV:]
				fmt.Println("    script              ", hx(tx1scro))

			}
			lockb := *rev(r[:4])
			r = r[4:]
			var lock uint32
			core.BytesToInt(lock, &lockb)
			fmt.Println("locktime                ", lock)
		}
		fmt.Println("\nRest:\n", hx(r))
		fmt.Println()
	}
}
