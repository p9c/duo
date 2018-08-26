package lb

import (
	"fmt"
	"testing"
)

func TestLockedBuffer(t *testing.T) {
	a := new(LockedBuffer)
	A := []byte("test")
	a.Load(&A)
	fmt.Println("a", a.Buf())
	b := new(LockedBuffer)
	b.Copy(a)
	fmt.Println("copy a", a.Buf(), "to b", b.Buf())
	fmt.Println("before move a", *a, "b", *b)
	b.Move(a)
	fmt.Println("after move a", *a, "b", *b)
	a.Link(b)
	fmt.Println("link b to (empty) a", a.Buf(), b.Buf())
	c := a.Buf()
	C := *c
	C[0] = 1
	fmt.Println("now both the same memory (changed byte zero of a only)", a.Buf(), b.Buf())
}
