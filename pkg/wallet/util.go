package wallet

import (
	"bytes"
	"encoding/binary"
)

func PreLenString(s string) (result []byte) {
	result = append([]byte{byte(len(s))}, s...)
	return
}

func PreLenBytes(b []byte) (result []byte) {
	result = append([]byte{byte(len(b))}, b...)
	return
}

func BytesToInt64(in []byte) (result int64) {
	resultB := bytes.NewBuffer(in)
	binary.Read(resultB, binary.LittleEndian, &result)
	return
}

func Uint64ToBytes(in uint64) (result []byte) {
	result = make([]byte, 8)
	binary.LittleEndian.PutUint64(result, in)
	return
}

func Int64ToBytes(in int64) (result []byte) {
	result = make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(in))
	return
}
