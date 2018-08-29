package crypt

import (
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/byteprint"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/password"
	"testing"
)

func TestCrypt(t *testing.T) {
	s := "testtest"
	a := NewCrypt()
	b := a.crypt.Rand(48).Buf()
	fmt.Println(b)
	c := a.crypt.Buf()
	fmt.Println(c)
	d := NewPassword()
	d.FromString(s)
	fmt.Println("fromstring", d)
	Print(d.Buf()).SP().Str(d.Buf()).SP().Quo("thisisastring").SP().Brc(&s).CR().Brc(',')
	fmt.Println("empty pw", NewCrypt().Password())
	var f *Crypt
	f.Password()
	f.Ciphertext()
	f.IV()
	f.Crypt()
	f.Null()
	g := new(Crypt)
	g.Password()
	g.Ciphertext()
	g.IV()
	g.Crypt()
	g.Null().crypt.Move(NewBytes().Rand(13)).Null()
	h := NewCrypt().Load(NewBytes().Rand(13))
	fmt.Println(h.crypt.Buf())
	fmt.Println(g.IsArmed(), f.IsArmed())
	fmt.Println(g.IsUnlocked(), f.IsUnlocked())
	fmt.Println(a.IsLoaded(), f.IsLoaded(), g.IsLoaded(), h.IsLoaded())
	x := f.Load(NewBytes().Rand(11))
	fmt.Println(*x.Crypt().Buf())
	fmt.Println(x.Crypt().IsSet())
	x.Crypt().Null().Null()
	fmt.Println(x)
	fmt.Println(*x.SetIV(NewBytes().Rand(12)).IV())
	fmt.Println(x.SetIV(NewBytes().Rand(11)))
	var z *Crypt
	fmt.Println(z)
	fmt.Println(z.SetError("nothing"))
	var n *Crypt
	fmt.Println("SetIV()", n.SetIV(NewBytes().Rand(12)).IV().Buf())
	var m *Crypt
	fmt.Println(m.SetIV(nil).IV())
	fmt.Println(*m.SetRandomIV().IV().Buf())
	// fmt.Println(NewCrypt().Generate(NewPassword().FromString("abcdef")).Password().Buf())
	// var v *Crypt
	// fmt.Println("nil receiver")
	// fmt.Println(v.Generate(NewPassword().FromString("ghijkl")).String())
	fmt.Println("valid receiver")
	fmt.Println("Generate()", NewCrypt().Unlock(NewPassword().FromString("abcdef")).SetRandomIV().Generate(NewPassword().FromString("abcdef")).String())
}
