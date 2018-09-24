package proto

import (
	"github.com/minio/highwayhash"
	"unsafe"
)

var (
	_int    int
	_uint   uint
	intlen  = unsafe.Sizeof(_int)
	uintlen = unsafe.Sizeof(_uint)
)

// Zero makes all the bytes in a slice zero
func Zero(b *[]byte) {
	B := *b
	for i := range B {
		B[i] = 0
	}
}

// IntToBytes converts any integer to a byte slice
func IntToBytes(u interface{}) (out *[]byte) {
	out = &[]byte{}
	switch u.(type) {
	case byte:
		out = &[]byte{u.(byte)}
	case int8:
		out = &[]byte{byte(u.(int8))}
	case uint16:
		t := u.(uint16)
		out = &[]byte{
			byte(t),
			byte(t >> 8),
		}
	case int16:
		t := u.(int16)
		out = &[]byte{
			byte(t),
			byte(t >> 8),
		}
	case uint32:
		t := u.(uint32)
		out = &[]byte{
			byte(t),
			byte(t >> 8),
			byte(t >> 16),
			byte(t >> 24),
		}
	case int32:
		t := u.(int32)
		out = &[]byte{
			byte(t),
			byte(t >> 8),
			byte(t >> 16),
			byte(t >> 24),
		}
	case uint64:
		t := u.(uint64)
		out = &[]byte{
			byte(t),
			byte(t >> 8),
			byte(t >> 16),
			byte(t >> 24),
			byte(t >> 32),
			byte(t >> 40),
			byte(t >> 48),
			byte(t >> 56),
		}
	case int64:
		t := u.(int64)
		out = &[]byte{
			byte(t),
			byte(t >> 8),
			byte(t >> 16),
			byte(t >> 24),
			byte(t >> 32),
			byte(t >> 40),
			byte(t >> 48),
			byte(t >> 56),
		}
	case int:
		switch intlen {
		case 2:
			out = IntToBytes(int16(u.(int)))
		case 4:
			out = IntToBytes(int32(u.(int)))
		case 8:
			out = IntToBytes(int64(u.(int)))
		}
	case uint:
		switch intlen {
		case 2:
			out = IntToBytes(uint16(u.(int)))
		case 4:
			out = IntToBytes(uint32(u.(int)))
		case 8:
			out = IntToBytes(uint64(u.(int)))
		}
	}
	return out
}

// BytesToInt takes up to 8 byte long byte slice and a pointer to the type of int wanted back
func BytesToInt(out interface{}, in *[]byte) {
	I := *in
	switch out.(type) {
	case *byte:
		*out.(*byte) = I[0]
	case *int8:
		i := int8(I[0])
		*out.(*int8) = i
	case *uint16:
		i := uint16(I[0])
		i += uint16(I[1]) << 8
		*out.(*uint16) = i
	case *int16:
		i := int16(I[0])
		i += int16(I[1]) << 8
		*out.(*int16) = i
	case *uint32:
		i := uint32(I[0])
		i += uint32(I[1]) << 8
		i += uint32(I[2]) << 16
		i += uint32(I[3]) << 24
		*out.(*uint32) = i
	case *int32:
		i := int32(I[0])
		i += int32(I[1]) << 8
		i += int32(I[2]) << 16
		i += int32(I[3]) << 24
		*out.(*int32) = i
	case *uint64:
		i := uint64(I[0])
		i += uint64(I[1]) << 8
		i += uint64(I[2]) << 16
		i += uint64(I[3]) << 24
		i += uint64(I[4]) << 32
		i += uint64(I[5]) << 40
		i += uint64(I[6]) << 48
		i += uint64(I[7]) << 56
		*out.(*uint64) = i
	case *int64:
		i := int64(I[0])
		i += int64(I[1]) << 8
		i += int64(I[2]) << 16
		i += int64(I[3]) << 24
		i += int64(I[4]) << 32
		i += int64(I[5]) << 40
		i += int64(I[6]) << 48
		i += int64(I[7]) << 56
		*out.(*int64) = i
	case *int:
		switch intlen {
		case 2:
			i := uint16(I[0])
			i += uint16(I[1]) << 8
			*out.(*uint16) = i
		case 4:
			i := int32(I[0])
			i += int32(I[1]) << 8
			i += int32(I[2]) << 16
			i += int32(I[3]) << 24
			*out.(*int32) = i
		case 8:
			i := int64(I[0])
			i += int64(I[1]) << 8
			i += int64(I[2]) << 16
			i += int64(I[3]) << 24
			i += int64(I[4]) << 32
			i += int64(I[5]) << 40
			i += int64(I[6]) << 48
			i += int64(I[7]) << 56
			*out.(*int64) = i
		}
	case *uint:
		switch intlen {
		case 2:
			i := uint16(I[0]) +
				uint16(I[1])<<8
			*out.(*uint16) = i
		case 4:
			i := uint32(I[0]) +
				uint32(I[1])<<8 +
				uint32(I[2])<<16 +
				uint32(I[3])<<24
			*out.(*uint32) = i
		case 8:
			i := uint64(I[0]) +
				uint64(I[1])<<8 +
				uint64(I[2])<<16 +
				uint64(I[3])<<24 +
				uint64(I[4])<<32 +
				uint64(I[5])<<40 +
				uint64(I[6])<<48 +
				uint64(I[7])<<56
			*out.(*uint64) = i
		}
	}

}

// Hash64 takes a byte slice and produces a 4 byte byte slice
func Hash64(in *[]byte) *[]byte {
	empty := make([]byte, highwayhash.Size)
	out := highwayhash.Sum64(*in, empty)
	return IntToBytes(out)
}
