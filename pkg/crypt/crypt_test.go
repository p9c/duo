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
	d.FromString(&s)
	fmt.Println("fromstring", d)
	Print(d.Buf()).SP().Str(d.Buf()).SP().Quo("thisisastring").SP().Brc(&s).CR().Brc(',')
	NewCrypt().Password()
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
	fmt.Println(x.SetIV(NewBytes().Rand(11)).Error())
	var z *Crypt
	fmt.Println(z.Error())
	fmt.Println(z.SetError("nothing"))
	var n *Crypt
	fmt.Println(n.SetIV(NewBytes().Rand(12)))
	var m *Crypt
	fmt.Println(m.SetIV(nil).IV().Error())
	fmt.Println(*m.SetRandomIV().IV().Buf())
}
