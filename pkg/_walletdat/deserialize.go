package walletdat

import (
	"errors"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"unsafe"
	"encoding/binary"
	"bytes"
)

var (
	String = ""
	Bytes = []byte{0}
	Uint64 = uint64(0)
	Int64 = int64(0)
	Uint32 = uint32(0)
	Int32 = int32(0)
	Int8 = int8(0)
	Byte = byte(0)
)

// Deserialize converts a binary representation back into it's in-memory form, trimming the prefix length byte from C strings and byte arrays
func Deserialize(b []byte, a ...interface{}) (o *[]interface{}, keyRem []byte, err error) {
	o = new([]interface{})
	b = []byte(string(b))
	bLen := len(b)
	buf := bytes.NewReader(b)
	for i := range a {
		switch a[i].(type) {
		case string:
			plen := b[0]+1
			*o = append(*o, string(b[1:plen]))
			if bLen > int(plen) {
				keyRem = b[plen:]
				break
			}
		case  []byte:
			plen := b[0]+1
			*o = append(*o, b[1:plen])
			if bLen > int(plen) {
				keyRem = b[plen:]
				break
			}
		default:
			binary.Read(buf, binary.LittleEndian, a[i])
			switch a[i].(type) {
			case int8:
				*o = append(*o, int8(b[0]))
				if bLen > 1 {
					keyRem = b[1:]
				}
			case byte:
				*o = append(*o, b[0])
				if bLen > 1 {
					keyRem = b[1:]
				}
			case int16, uint16:
				B := bytes.NewBuffer(b[:2])
				binary.Read(B, binary.LittleEndian, a[i])
				*o = append(*o, a[i])
				if bLen > 2 {
					keyRem = b[2:]
				}
			case int32, uint32:
				B := bytes.NewBuffer(b[:4])
				binary.Read(B, binary.LittleEndian, a[i])
				*o = append(*o, a[i])
				if bLen > 4 {
					keyRem = b[4:]
				}
			case int64, uint64:
				B := bytes.NewBuffer(b[:8])
				binary.Read(B, binary.LittleEndian, a[i])
				*o = append(*o, a[i])
				if bLen > 8 {
					keyRem = b[8:]
				}
			case int, uint:
				var I int
				iSize := unsafe.Sizeof(I)
				B := bytes.NewBuffer(b[:iSize])
				binary.Read(B, binary.LittleEndian, a[i])
				*o = append(*o, a[i])
				if bLen > int(iSize) {
					keyRem = b[iSize:]
				}
			case Uint.U160:
				B := b[:20]
				a[i] = Uint.Zero160().FromBytes(B)
				*o = append(*o, a[i])
				if bLen > 20 {
					keyRem = b[20:]
				}
			case Uint.U256:
				B := b[:32]
				a[i] = Uint.Zero256().FromBytes(B)
				*o = append(*o, a[i])
				if bLen > 32 {
					keyRem = b[32:]
				}
			default:
				err = errors.New("Data type not handled")
			}
		}
	}
	return
}
