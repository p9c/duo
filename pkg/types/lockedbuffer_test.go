package types

import (
	"fmt"
	"testing"
)

func TestLockedBuffer(t *testing.T) {
	fmt.Println(NewLockedBuffer())
	fmt.Println(NewLockedBuffer().Len())
	fmt.Println(NewLockedBuffer().FromRandomBytes(23).ToBytes().ToByteSlice())
	a := "This is a test"
	A := &a
	fmt.Println(*NewLockedBuffer().FromBytes(NewBytes().FromString(A)).ToBytes().ToString())
	fmt.Println(NewLockedBuffer().ToBytes())
	fmt.Println(*NewLockedBuffer().FromRandomBytes(23).FromBytes(NewBytes().FromString(A)).ToBytes().ToString())
	fmt.Println(*NewLockedBuffer().FromLockedBuffer(NewLockedBuffer().FromRandomBytes(23).FromRandomBytes(33)).ToBytes().ToString())
	fmt.Println(NewLockedBuffer().FromRandomBytes(23).Len())
	fmt.Println(*NewLockedBuffer().WithSize(33).ToBytes().ToString())
	l := NewLockedBuffer().FromRandomBytes(22).ToLockedBuffer()
	fmt.Println("LockedBuffer #1", l)
	L := NewLockedBuffer().FromLockedBuffer(l)
	fmt.Println("After FromLockedBuffer() #1", l, "#2", L)
	m := []byte(a)
	M := NewLockedBuffer().FromByteSlice(&m)
	fmt.Println(*M.ToBytes().ToString(), "FromByteSlice()", m)
	fmt.Println("Should be 0s:", *M.ToByteSlice())
	n := NewLockedBuffer().FromRandomBytes(13)
	o := NewLockedBuffer().FromRandomBytes(23)
	fmt.Println("n", n, "o", o)
	p := o.FromLockedBuffer(n)
	fmt.Println("p = o.FromLockedBuffer(n) p", p, "o", o, "n", n)
	o.WithSize(19)
	fmt.Println("n after o.FromLockedBuffer(n)", n, "o after WithSize()", o)
	x := "final test"
	X := []byte(x)
	fmt.Println("before FromByteSlice()", p.value.Buffer())
	p.FromByteSlice(&X)
	fmt.Println("after FromByteSlice()", p.value.Buffer())
	s := NewLockedBuffer().WithSize(20)
	fmt.Println(s, s.value, s.value.Buffer())
	Sb := s.value.Buffer()
	fmt.Println(Sb)
	Sb[1] = 220
	fmt.Println(Sb)
	ss := s.Copy()
	fmt.Println(s, ss)
}
