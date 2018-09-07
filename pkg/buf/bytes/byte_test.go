package buf

import (
	"testing"
)

func TestByte(t *testing.T) {
	NewByte(0).
		Rand().
		Free().
		Link(NewByte(1)).
		OfLen(1).
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
	v.Copy(nil)
	v.err = nil
	v.Status()
}
