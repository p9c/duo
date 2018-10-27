package block

import "github.com/parallelcointeam/duo/pkg/core"

// Encode turns a completed block data structure into a serialised stream in protocol format
func Encode(in Raw) (out []byte) {
	out = append(out, *core.IntToBytes(in.Version)...)
	out = append(out, *rev(in.HashPrevBlock)...)
	out = append(out, *rev(in.HashMerkleRoot)...)
	out = append(out, *core.IntToBytes(in.Time)...)
	out = append(out, *rev(in.Bits)...)
	out = append(out, *core.IntToBytes(in.Nonce)...)
	out = AppendCompactInt(out, len(in.Transactions))
	for i := range in.Transactions {
		out = append(out, *core.IntToBytes(in.Transactions[i].Version)...)
		T := in.Transactions[i]
		out = AppendCompactInt(out, len(T.Ins))
		for ins := range T.Ins {
			out = append(out, T.Ins[ins].PrevTxHash...)
			out = append(out, *core.IntToBytes(T.Ins[ins].PrevTxoutIndex)...)
			out = AppendCompactInt(out, len(T.Ins[ins].Script))
			out = append(out, T.Ins[ins].Script...)
			out = append(out, *core.IntToBytes(T.Ins[ins].Sequence)...)
		}
		out = AppendCompactInt(out, len(T.Outs))
		for outs := range T.Outs {
			out = append(out, *core.IntToBytes(T.Outs[outs].Value)...)
			out = AppendCompactInt(out, len(T.Outs[outs].Script))
			out = append(out, T.Outs[outs].Script...)
		}
		out = append(out, *core.IntToBytes(uint32(in.Transactions[i].Locktime))...)
	}
	return
}
