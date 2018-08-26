package b

import (
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	a := new(Bytes)
	A := []byte("test")
	a.Assign(&A)
	fmt.Println("a", *a.Buf())
	b := new(Bytes)
	b.Copy(a)
	fmt.Println("copy a to b", *b.Buf())
	fmt.Println("before move a", *a, "b", *b)
	b.Move(a)
	fmt.Println("after move a", *a, "b", *b)
	a.Link(b)
	fmt.Println("link emptied b to a", *a.Buf(), *b.Buf())
	c := *a.val
	c[0] = 1
	fmt.Println("now both the same memory (changed byte zero of first only)", *a.Buf(), *b.Buf())
}
