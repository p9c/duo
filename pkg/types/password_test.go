package types

import (
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {
	P := "testing"
	p := NewPassword().FromString(&P)
	fmt.Println(p.ToString())
}
