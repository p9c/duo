package crypt

import (
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/password"
	"testing"
)

func TestCrypt(t *testing.T) {
	s := "testtest"
	a := NewCrypt()
	b := a.Rand(48).Buf()
	fmt.Println(b)
	c := a.Buf()
	fmt.Println(c)
	d := NewPassword()
	d.FromString(&s)
	fmt.Println("fromstring", string(*d.Buf()))
}
