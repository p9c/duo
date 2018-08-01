package ser

import (
	"errors"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"
)

const (
	MaxSize = 0x02000000
	Network = 1
	Disk    = 2
	GetHash = 4
)

type StreamPlaceholder struct {
	Type, Version int
}
type Data [][]byte

type FlatData struct {
	Begin, End []byte
}

// DataStream is a stream of bytes
type DataStream struct {
	Data                 []byte
	ReadPos              uint
	State, ExceptionMask int16
	Type, Version        int
}

// AutoFile is an auto file
type AutoFile struct {
	File                 os.File
	State, ExceptionMask int16
	Type, Version        int
}

// BufferedFile is a buffered file
type BufferedFile struct {
	src                                os.File
	srcPos, readPos, readLimit, rewind uint64
	buf                                []byte
	state, exceptionMask               int16
	Type, Version                      int
}

// Serialize converts a variable into its binary representation, optionally adding a length prefix for C format strings and byte arrays. Strings and byte slices over 255 bytes in length do not get this prefix and are assumed to be stored with nothing further concatenated to the end (for further binary data parsing)
func Serialize(a interface{}) (b []byte, err error) {
	buf := new(bytes.Buffer)
	preLen := false
	switch a.(type) {
	case string:
		if len(a.(string)) < 256 {
			preLen = true
		}
	case  []byte:
		if len(a.([]byte)) < 256 {
			preLen = true
		}
	}
	err = binary.Write(buf, binary.LittleEndian, a)
	b = buf.Bytes()
	if preLen {
		b = append([]byte{byte(len(b))}, b...)
	}
	return
}

// Deserialize converts a binary representation back into it's in-memory form, trimming the prefix length byte from C strings and byte arrays
func Deserialize(b []byte, a interface{}) (keyRem []byte, err error) {
	bLen := len(b)
	buf := bytes.NewReader(b)
	preLen := false
	switch a.(type) {
	case string:
		if len(a.(string)) < 256 {
			preLen = true
		}
	case  []byte:
		if len(a.([]byte)) < 256 {
			preLen = true
		}
	}
	binary.Read(buf, binary.LittleEndian, a)
	if preLen {
		switch a.(type) {
		case string:
			plen := byte(a.(string)[0])
			a = a.(string)[1:plen]
			if bLen > int(plen)+1 {
				keyRem = b[plen:]
			}
		case []byte:
			plen := byte(a.([]byte)[0])
			a = a.([]byte)[1:plen]
			if bLen > int(plen)+1 {
				keyRem = b[plen:]
			}
		}
	} else {
		switch a.(type) {
		case int8:
			a = int8(b[0])
			if bLen > 1 {
				keyRem = b[1:]
			}
		case byte:
			a = b[0]
			if bLen > 1 {
				keyRem = b[1:]
			}
		case int16, uint16:
			B := bytes.NewBuffer(b[:2])
			binary.Read(B, binary.LittleEndian, &a)
			if bLen > 2 {
				keyRem = b[2:]
			}
		case int32, uint32:
			B := bytes.NewBuffer(b[:4])
			binary.Read(B, binary.LittleEndian, &a)
			if bLen > 4 {
				keyRem = b[4:]
			}
		case int64, uint64:
			B := bytes.NewBuffer(b[:8])
			binary.Read(B, binary.LittleEndian, &a)
			if bLen > 8 {
				keyRem = b[8:]
			}
		case int, uint:
			var I int
			iSize := unsafe.Sizeof(I)
			B := bytes.NewBuffer(b[:iSize])
			binary.Read(B, binary.LittleEndian, &a)
			if bLen > int(iSize) {
				keyRem = b[iSize:]
			}
		case Uint.U160:
			B := b[:20]
			a = Uint.Zero160().FromBytes(B)
			if bLen > 20 {
				keyRem = b[20:]
			}
		case Uint.U256:
			B := b[:32]
			a = Uint.Zero256().FromBytes(B)
			if bLen > 32 {
				keyRem = b[32:]
			}
		default:
			err = errors.New("Data type not handled")
		}
	}
	return
}
