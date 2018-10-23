package block

import "github.com/parallelcointeam/duo/pkg/core"

// AppendVarint takes any type of integer and returns the Varint. This is for the bitcoin protocol format varint that stores up to FC in 1 byte, FD means two bytes with uint16 after, FE 4 bytes after, FF 8 bytes after. The bytes are in little-endian, with the MSB first
func AppendVarint(to []byte, in interface{}) (out []byte) {
	var outint uint64
	switch in.(type) {
	case uint:
		outint = uint64(in.(uint))
	case byte:
		outint = uint64(in.(byte))
	case uint16:
		outint = uint64(in.(uint16))
	case uint32:
		outint = uint64(in.(uint32))
	case uint64:
		outint = uint64(in.(uint64))
	case int:
		outint = uint64(in.(int))
	case int8:
		outint = uint64(in.(int8))
	case int16:
		outint = uint64(in.(int16))
	case int32:
		outint = uint64(in.(int32))
	case int64:
		outint = uint64(in.(int64))
	default:
		return to
	}
	out = to
	// Bytes are appended from lowest to highest, ie little-endian
	switch {
	case outint < 0xFD:
		out = append(out, byte(outint))
	case outint >= 0xFD && outint < uint64(^uint16(0)):
		out = []byte{0xFE}
		t := uint16(outint)
		out = append(out, byte(t<<8>>8))
		out = append(out, byte(t>>8))
	case outint < uint64(^uint32(0)):
		t := uint32(outint)
		out = append(out, byte(t>>24))
		out = append(out, byte(t<<8>>24))
		out = append(out, byte(t<<16>>24))
		out = append(out, byte(t<<24>>24))
	default:
		t := uint64(outint)
		out = append(out, byte(t>>56))
		out = append(out, byte(t<<8>>56))
		out = append(out, byte(t<<16>>56))
		out = append(out, byte(t<<24>>56))
		out = append(out, byte(t<<32>>56))
		out = append(out, byte(t<<40>>56))
		out = append(out, byte(t<<48>>56))
		out = append(out, byte(t<<56>>56))
	}
	return
}

// ExtractVarint reads the first varint contained in a given byte slice and returns the value according to the type of the typ parameter, and slices the input bytes removing the extracted integer
func ExtractVarint(typ interface{}, in []byte) (outbytes []byte, outint interface{}) {
	switch {
	case in[0] < 0xFD:
		outbytes = in[1:]
		outint = outbytes
	case in[0] == 0xFD:
		var out uint16
		t := in[1:3]
		core.BytesToInt(&out, &t)
		outbytes = in[3:]
		outint = out
	case in[0] == 0xFE:
		var out uint32
		t := in[1:5]
		core.BytesToInt(&out, &t)
		outbytes = in[5:]
		outint = out
	case in[0] == 0xFF:
		var out uint64
		t := in[1:9]
		core.BytesToInt(&out, &t)
		outbytes = in[9:]
		outint = out
	}
	return
}
