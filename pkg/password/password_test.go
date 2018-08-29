package password

import (
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {
	p := "this is a test!"
	P := NewPassword().FromString(p)
	fmt.Println("original", p, "From->ToString", *P.ToString())
	var n *Password
	fmt.Println(*n.ToString())
	fmt.Println(*n.FromString(p).ToString())
	fmt.Println(P.SetError("testing"))
	fmt.Println(P.String())
	fmt.Println(P.MarshalJSON())
	var pp *Password
	fmt.Println(pp.MarshalJSON())
	pp = new(Password)
	fmt.Println(pp.MarshalJSON())
}
