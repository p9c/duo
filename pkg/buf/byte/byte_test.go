package byt

import (
	"testing"
)

func TestByte(t *testing.T) {
	NewByte(0).
		Rand().
		Free().
		Link(NewByte(1)).
		OfSize(1).
		Copy(byte(4)).
		Null().
		Copy(1).
		Buf()
	var v *Byte
	v.Len()
	v.Array()
	v.Status()
	v.Coding()
	v = v.Copy(v).(*Byte)
	v.err = nil
	v.Status()
}
