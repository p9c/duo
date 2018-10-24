package block

import "encoding/hex"

func rev(in []byte) (out *[]byte) {
	o := make([]byte, len(in))
	out = &o
	for i := range in {
		(*out)[len(in)-i-1] = in[i]
	}
	return
}

func hx(in []byte) string {
	return hex.EncodeToString(in)
}
