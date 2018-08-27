package crypt

import (
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/byteprint"
	. "gitlab.com/parallelcoin/duo/pkg/password"
	// "os"
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
	fmt.Println("fromstring")
	Print(d.Buf()).SP().Str(d.Buf()).SP().Quo("thisisastring").SP().Brc(&s).CR().Brc(',')
	NewCrypt().Password()
	var f *Crypt
	f.Password()
}
