package block

import (
	"github.com/parallelcointeam/duo/pkg/core"
)

func split(in []byte, pos int) (out []byte, piece []byte) {
	out = in[pos:]
	piece = in[:pos]
	return
}

// DecodeBlock reads a protocol serialised block and returns the raw block structure
func DecodeBlock(in []byte) (out Raw) {

	in, t := split(in, 4)
	core.BytesToInt(&out.Version, &t)

	in, t = split(in, 32)
	out.HashPrevBlock = *rev(t)

	in, t = split(in, 32)
	out.HashMerkleRoot = *rev(t)

	in, t = split(in, 4)
	core.BytesToInt(&out.Time, &t)

	in, t = split(in, 4)
	out.Bits = *rev(t)

	in, t = split(in, 4)
	core.BytesToInt(&out.Nonce, rev(t))

	var txCount uint64
	var txCountIface interface{}
	in, txCountIface = ExtractCompactInt(txCount, in)
	txCount = txCountIface.(uint64)

	txc := int(txCount)
	for txs := 0; txs < txc; txs++ {

		var tx Tx
		out.Transactions = append(out.Transactions, tx)

		var t []byte
		in, t = split(in, 4)
		core.BytesToInt(&out.Transactions[txs].Version, &t)

		var txV uint32
		var txVi uint64
		var txI interface{}
		in, txI = ExtractCompactInt(txVi, in)
		txiV := int(txI.(uint64))

		for txis := 0; txis < txiV; txis++ {
			var txIn TxIn

			in, txIn.PrevTxHash = split(in, 32)

			in, t = split(in, 4)
			core.BytesToInt(&txIn.PrevTxoutIndex, &t)

			var txsl interface{}
			in, txsl = ExtractCompactInt(txV, in)
			txV = txsl.(uint32)

			txIn.Script = make([]byte, txV)
			copy(txIn.Script, in[:txV])
			in = in[txV:]

			t = in[:4]
			core.BytesToInt(&txIn.Sequence, &t)
			in = in[4:]

			out.Transactions[txs].Ins = append(out.Transactions[txs].Ins, txIn)
		}

		in, txI = ExtractCompactInt(txV, in)
		txoV := int(txI.(uint32))

		for txos := 0; txos < txoV; txos++ {

			var txOut TxOut

			tx1val := in[:8]
			in = in[8:]
			core.BytesToInt(&txOut.Value, &tx1val)

			in, txI = ExtractCompactInt(txV, in)
			txV = txI.(uint32)

			txOut.Script = make([]byte, txV)
			copy(txOut.Script, in[:txV])
			in = in[txV:]

			out.Transactions[txs].Outs = append(out.Transactions[txs].Outs, txOut)
		}
		var lockb []byte
		in, lockb = split(in, 4)
		core.BytesToInt(out.Transactions[txs].Locktime, &lockb)
	}

	return
}

// EncodeBlock turns a completed block data structure into a serialised stream in protocol format
func EncodeBlock(in Raw) (out []byte) {

	return
}
