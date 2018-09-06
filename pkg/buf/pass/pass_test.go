package passbuf

import (
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {
	p := "this is a test!"
	P := New().FromString(p)
	fmt.Println("original", p, "From->String", P.String())
	var n *Password
	fmt.Println(n.String())
	fmt.Println(n.FromString(p).String())
	fmt.Println(P.SetError("testing"))
	fmt.Println(P.String())
	fmt.Println(P.MarshalJSON())
	var pp *Password
	fmt.Println(pp.MarshalJSON())
	pp = new(Password)
	fmt.Println(pp.MarshalJSON())
}
