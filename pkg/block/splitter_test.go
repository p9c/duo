package block

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/anaskhan96/base58check"

	"github.com/parallelcointeam/duo/gocoin/btc"

	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/hash160"
	"github.com/parallelcointeam/duo/pkg/key"
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
	for i := 0; i < 20; i++ {
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
		// fmt.Print("Version ", Version)
		// switch Version {
		// case 2:
		// 	fmt.Println(" : SHA256D PoW")
		// case 514:
		// 	fmt.Println(" : Scrypt PoW")
		// }

		// HashPrevBlock := *rev(r[:32])
		r = r[32:]
		// fmt.Println("HashPrevBlock          ", hx(HashPrevBlock))

		// HashMerkleRoot := *rev(r[:32])
		r = r[32:]
		// fmt.Println("HashMerkleRoot         ", hx(HashMerkleRoot))

		ti := r[:4]
		r = r[4:]
		var t32 int32
		core.BytesToInt(&t32, &ti)
		// BlockTime := int64(t32)
		// fmt.Println("Unix timestamp         ", hx(ti))
		// fmt.Println("Time", time.Unix(BlockTime, 0))

		// Bits := *rev(r[:4])
		r = r[4:]
		// fmt.Println("Bits                   ", hx(Bits))
		// coeff := Bits[0]
		// base := Bits[1:]
		// tail := make([]byte, coeff-3)
		// tgt := append(base, tail...)
		// Target := append(make([]byte, 32-len(tgt)), tgt...)
		// fmt.Println("                 Target", hx(Target))

		nn := r[:4]
		nn = *rev(nn)
		r = r[4:]
		var Nonce uint32
		core.BytesToInt(&Nonce, &nn)
		// fmt.Println("Nonce                  ", Nonce)

		var txCount uint64
		var txCountIface interface{}
		r, txCountIface = ExtractCompactInt(txCount, r)
		txCount = txCountIface.(uint64)
		// fmt.Println("TxCount                 ", txCount)

		txc := int(txCount)
		for txs := 0; txs < txc; txs++ {

			fmt.Println("Transaction", txs)

			rV := r[:4]
			r = r[4:]
			var tx0v uint32
			core.BytesToInt(&tx0v, &rV)
			// fmt.Println("    tx version", tx0v)

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
			// fmt.Println("    in-counter", txiV)

			for txis := 0; txis < txiV; txis++ {

				tx1pth := *rev(r[:32])
				r = r[32:]
				// tx1txi := r[:4]
				r = r[4:]

				var txsl interface{}
				r, txsl = ExtractCompactInt(txV, r)
				txV = txsl.(uint64)
				// fmt.Println("     Txin script length", txV)
				tx1scr := r[:txV]
				r = r[txV:]

				// tx1seq := r[:4]
				r = r[4:]
				if bytes.Compare(make([]byte, 32), tx1pth) != 0 {
					// fmt.Println("    Txin", txis)
					// fmt.Println("             PrevTxHash", hx(tx1pth))

					// GET VALUE OF PREVTX
					value := uint64(node.GetTxValue(tx1pth) * float64(core.COIN))
					// fmt.Println("                 value", value)

					// fmt.Printf("       Prev Txout Index %08x\n", tx1txi)

					// fmt.Println("                 script", hx(tx1scr))

					var p1 string
					if len(tx1scr) > 64 {
						p1 = hx(tx1scr[1 : tx1scr[0]+1])
						// fmt.Println("           signature  ", p1)
						if len(tx1scr) > len(p1)/2 {
							rem := *rev(tx1scr[tx1scr[0]+2:])
							k, _ := base58check.Encode(key.B58prefixes["mainnet"]["pubkey"], hx(*hash160.Sum(&rem)))

							fmt.Printf("    %s -%012d\n", k, value)

							// fmt.Println("                    <<<", k)
							// fmt.Println("                    >>> Typical payment redemption")
						}
					}
					// fmt.Println("             seq number", hx(tx1seq))
				} else {
					// fmt.Println(">>> Generation transaction")
				}

			}
			r, txI = ExtractCompactInt(txV, r)
			txoV := int(txI.(uint64))
			// fmt.Println("    out-counter", txoV)

			for txos := 0; txos < txoV; txos++ {

				// fmt.Println("    Txout", txos)

				tx1val := r[:8]
				r = r[8:]
				var tx1V uint64
				core.BytesToInt(&tx1V, &tx1val)
				value := uint64(0)
				core.BytesToInt(&value, &tx1val)
				// fmt.Printf("                  value %4.7f\n", float64(tx1V)/core.COIN)

				r, txI = ExtractCompactInt(txV, r)
				txV = txI.(uint64)
				// fmt.Println("        Txout script length", txV)

				tx1scro := r[:txV]
				r = r[txV:]

				str, _ := btc.ScriptToText(tx1scro)
				fmt.Print("    ")
				// fmt.Println("                       ", str)
				if len(str) == 5 {
					k, _ := base58check.Encode(key.B58prefixes["mainnet"]["pubkey"], str[2])
					fmt.Printf("%s +%012d\n", k, value)

					// fmt.Println("                    >>>", k)
				} else {
					// fmt.Println(">>>>>>>", hx(tx1scro[1:tx1scro[0]]))
					k := tx1scro[1:tx1scro[0]]
					// fmt.Println(hx(k))
					K, _ := base58check.Encode(key.B58prefixes["mainnet"]["pubkey"], hx(*hash160.Sum(&k)))
					// fmt.Println(K, "+", value)
					fmt.Printf("%s +%012d\n", K, value)

					// fmt.Println("                    >>>", K)
					// fmt.Println("                 script", hx(tx1scro))
				}

			}
			lockb := r[:4]
			r = r[4:]
			var lock uint32
			core.BytesToInt(lock, &lockb)
			if lock != 0 {
				fmt.Println("    locktime", lock)
			}
		}
		// fmt.Println("Rest:\n", hx(r))
		// fmt.Println()
	}
}

func TestDecodeEncodeBlock(t *testing.T) {
	b := make([]byte, 4)
	node := sync.NewNode()
	best := node.LegacyGetBestBlockHeight()
	var B uint32
	for i := 0; i < 10; i++ {
		for {
			rand.Read(b)
			core.BytesToInt(&B, &b)
			if B < best {
				break
			}
		}
		rb := node.GetRawBlock(uint64(B))
		in := *rb
		fmt.Println("\nblock", B)
		out := Decode(in)
		// j, _ := json.MarshalIndent(out, "", "  ")
		// fmt.Println(string(j))
		re := Encode(out)
		fmt.Println(hx(re))
		fmt.Println(hx(in))
	}
}
