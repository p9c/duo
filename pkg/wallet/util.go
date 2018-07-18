package wallet

import (
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/ec"
	"bytes"
	"encoding/binary"
)

// FormatString prepends a byte with the length of a string for wire/storage formatting
func FormatString(s string) (result []byte) {
	result = append([]byte{byte(len(s))}, s...)
	return
}

// FormatBytes prepends a byte with the length of a byte slice for wire/storage formatting
func FormatBytes(b []byte) (result []byte) {
	result = append([]byte{byte(len(b))}, b...)
	return
}

func BytesToUint32(in []byte) (result uint32) {
	resultB := bytes.NewBuffer(in)
	binary.Read(resultB, binary.LittleEndian, &result)
	return
}


func Uint32ToBytes(in uint32) (result []byte) {
	resultB := bytes.NewBuffer([]byte{})
	binary.Write(resultB, binary.LittleEndian, &in)
	return resultB.Bytes()
}

func BytesToInt64(in []byte) (result int64) {
	resultB := bytes.NewBuffer(in)
	binary.Read(resultB, binary.LittleEndian, &result)
	return
}

func Int64ToBytes(in int64) (result []byte) {
	result = make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(in))
	return
}
func BytesToUint64(in []byte) (result uint64) {
	resultB := bytes.NewBuffer(in)
	binary.Read(resultB, binary.LittleEndian, &result)
	return
}


func Uint64ToBytes(in uint64) (result []byte) {
	result = make([]byte, 8)
	binary.LittleEndian.PutUint64(result, in)
	return
}


// Append a byte slice to a byte slice in the caller's scope
func Append(b *[]byte, B ...[]byte) {
	for i := range *b {
		*b = append(*b, B[i]...)
	}
}

func ToPub(pubEC *ec.PublicKey) (pub *key.Pub) {
	pub = &key.Pub{}
	pub.SetPub(pubEC)
	return
}

func ParsePub(pub []byte) (key *ec.PublicKey, err error) {
	return ec.ParsePubKey(pub, ec.S256())
}

func SetPriv(priv []byte) (result *key.Priv) {
	result = &key.Priv{}
	result.Set(priv)
	return
}
