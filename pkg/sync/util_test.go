package sync

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestVarints(t *testing.T) {
	var bytes []byte
	slice := []interface{}{
		int(5),
		uint(13),
		int8(100),
		uint8(140),
		int16(10035),
		uint16(10250),
		int32(10020000),
		uint32(1003230000),
		int64(10000006932),
		uint64(1002000302),
	}
	for j := range slice {
		bytes = AppendVarint(bytes, slice[j])
		fmt.Println(slice[j], hex.EncodeToString(bytes))
	}
	for i := range slice {
		var outint interface{}
		bytes, outint = ExtractVarint(slice[i], bytes)
		fmt.Println(outint, hex.EncodeToString(bytes))
	}
}
