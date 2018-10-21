package sync

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestVarints(t *testing.T) {
	var bytes []byte
	slice := []interface{}{
		int(1),
		uint(10),
		int8(100),
		uint8(100),
		int16(10000),
		uint16(10000),
		int32(100000000),
		uint32(1000000000),
		int64(1000000000),
		uint64(1000000000),
	}
	for j := range slice {
		// fmt.Println(hex.EncodeToString(bytes))
		bytes = AppendVarint(bytes, slice[j])
		fmt.Println(hex.EncodeToString(bytes))
	}
	for i := range slice {
		var outint interface{}
		bytes, outint = ExtractVarint(slice[i], bytes)
		fmt.Println(outint, ", ", hex.EncodeToString(bytes))
	}
}
