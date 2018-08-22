package types

import (
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	S := "this is a test"
	Z := "test string number 2"
	fmt.Println(NewBytes())
	s := []byte(S)
	z := []byte(Z)
	b := NewBytes().FromByteSlice(&s)
	fmt.Println(string(*b.value))
	c := NewBytes().FromBytes(b)
	S = *c.ToString()
	fmt.Println(S)
	f := NewBytes().FromString(&S)
	fmt.Println(*f.ToString())
	d := f.ToBytes()
	fmt.Println(*d.ToString())
	d.FromByteSlice(&z)
	e := NewBytes().FromBytes(d)
	fmt.Println(*e.ToString())
	E := NewBytes().WithSize(23)
	fmt.Println(E.Len())
	g := E.WithSize(32)
	fmt.Println(g.Len())
	fmt.Println(E.Len())
	fmt.Println(*E.FromBytes(f).ToString())
	fmt.Println(c.ToByteSlice())
	fmt.Println(NewBytes().Len())
}
