package ser

import (
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

// GetSerializeSize returns the length in bytes to serialize a variable
func GetSerializeSize(a interface{}) uint {
	return uint(unsafe.Sizeof(a))
}

// Serialize converts a variable into its binary representation, optionally adding a length prefix if requested (C format strings and byte arrays)
func Serialize(a interface{}) (b []byte, err error) {
	buf := new(bytes.Buffer)
	preLen := false
	switch a.(type) {
	case string, []byte:
		preLen = true
	}
	err = binary.Write(buf, binary.LittleEndian, a)
	b = buf.Bytes()
	if preLen {
		b = append([]byte{byte(len(b))}, b...)
	}
	return
}

// Deserialize converts a binary representation back into it's in-memory form, trimming the prefix length byte from C strings and byte arrays
func Deserialize(b []byte, a interface{}) interface{} {
	buf := bytes.NewReader(b)
	preLen := false
	switch a.(type) {
	case string, []byte:
		preLen = true
	}
	binary.Read(buf, binary.LittleEndian, a)
	if preLen {
		switch a.(type) {
		case string:
			plen := byte(a.(string)[0])
			a = a.(string)[2:plen]
		case []byte:
			plen := byte(a.([]byte)[0])
			a = a.([]byte)[2:plen]
		}
	}
	return a
}

// GetPreLen cuts a prefix length marked section of bytes, returns the value and the remainder slice
func GetPreLen(i []byte) (first, remainder []byte) {
	b := []byte(string(i))
	preLen := int(b[0]) + 1
	if len(b) > preLen {
		first = b[1:preLen]
		remainder = b[preLen:]
	} else {
		first = b[1:]
	}
	return
}

// GetPreLenString cuts a prefix length marked section of bytes, returns the value and the remainder slice
func GetPreLenString(i []byte) (first string, remainder []byte) {
	b := []byte(string(i))
	preLen := int(b[0]) + 1
	if len(b) > preLen {
		first = string(b[1:preLen])
		remainder = b[preLen:]
	} else {
		first = string(b[1:])
	}
	return
}

// GetInt extracts an integer of arbitrary type from the front of a byte slice
func GetInt(b []byte, i interface{}) (remainder []byte) {
	switch i.(type) {
	case int8:
		i = int8(b[0])
		return b[1:]
	case byte:
			i = b[0]
			return b[1:]
		case int16, uint16:
		B := bytes.NewBuffer(b[:2])
		binary.Read(B, binary.LittleEndian, &i)
		return b[2:]
	case int32, uint32:
		B := bytes.NewBuffer(b[:4])
		binary.Read(B, binary.LittleEndian, &i)
		return b[4:]
	case int64, uint64:
		B := bytes.NewBuffer(b[:8])
		binary.Read(B, binary.LittleEndian, &i)
		return b[8:]
	case int, uint:
		var I int
		iSize := unsafe.Sizeof(I)
		B := bytes.NewBuffer(b[:iSize])
		binary.Read(B, binary.LittleEndian, &i)
		return b[iSize:]
	case Uint.U160:
		B := b[:20]
		i = Uint.Zero160().FromBytes(B)
		return b[20:]
	case Uint.U256:
		B := b[:32]
		i = Uint.Zero256().FromBytes(B)
		return b[32:]
	}
	return
}