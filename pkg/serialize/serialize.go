package ser

import (
	"bytes"
	"encoding/binary"
	"os"
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

// DataStream is a[i] stream of bytes
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

// BufferedFile is a[i] buffered file
type BufferedFile struct {
	src                                os.File
	srcPos, readPos, readLimit, rewind uint64
	buf                                []byte
	state, exceptionMask               int16
	Type, Version                      int
}

// Serialize converts a[i] variable into its binary representation, optionally adding a[i] length prefix for C format strings and byte arrays. Strings and byte slices over 255 bytes in length do not get this prefix and are assumed to be stored with nothing further concatenated to the end (for further binary data parsing)
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
