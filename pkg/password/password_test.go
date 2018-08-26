package password

import (
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {
	p := "this is a test!"
	P := NewPassword().FromString(&p)
	fmt.Println("original", p, "From->ToString", *P.ToString())
	var n *Password
	fmt.Println(n.ToString())
	fmt.Println(*n.FromString(&p).ToString())
}
