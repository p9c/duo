package types

import (
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {
	P := "testing"
	p := NewPassword()
	p.FromString(&P)
	fmt.Println(p.Copy().ToString())
	q := p.Copy()
	fmt.Println(p, p.Copy().ToString(), q, q.Copy().ToString())
	fmt.Println(p.value.Buffer(), q.value.Buffer())
	r := p.Copy()
	fmt.Println(r)
	R := r.ToString()
	fmt.Println("The string", R, "The contents of the buffer after ToString()", r.value.Buffer(), "The original after copy", p.value.Buffer())
	s := p.ToPassword()
	fmt.Println("after ToPassword() original", p, "(should be nil)")
	fmt.Println("new:", s.value.Buffer())
	S := NewPassword().FromPassword(s)
	fmt.Println("after FromPassword() original", s, "(should be nil)")
	fmt.Println("new:", S.value.Buffer())

}
