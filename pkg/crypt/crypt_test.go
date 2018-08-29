package crypt

import (
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/byteprint"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
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
	pass := NewPassword().FromString("abcdef")
	C := NewCrypt().SetRandomIV().Generate(pass)
	C.Decrypt(NewBytes().Rand(64))
	fmt.Println("Generate()", C.String())
	C.Disarm()
	fmt.Println("Disarm()", C.String())
	C.Decrypt(NewBytes().Rand(64))
	C.Lock()
	fmt.Println("Lock()", C.String())
	C.Decrypt(NewBytes().Rand(64))
	pass = NewPassword().FromString("abcdef")
	C.Unlock(pass)
	fmt.Println("Unlock()", C.String())
	C.Decrypt(NewBytes().Rand(64))
	C.Arm()
	fmt.Println("Arm()", C.String())
	C.Decrypt(NewBytes().Rand(64))
	ss := []byte("this is a test")
	fmt.Println("plaintext '" + string(ss) + "'")
	fmt.Println("bytes ", ss)
	bb := C.Encrypt(NewLockedBuffer().Load(&ss))
	fmt.Println("ciphertext", *bb.Buf())
	BB := C.Decrypt(bb)
	fmt.Println("recovertext", *BB.Buf())
	fmt.Println("recover as string '" + string(*BB.Buf()) + "'")
	C.Lock()
	fmt.Println("Unlock()", C.String())
	var cc *Crypt
	cc.Arm()
	cc = new(Crypt)
	cc.Arm()
	cc.password = NewPassword().FromString("abcdef")
	cc.Arm()
	cc.Load(NewBytes().Rand(32))
	cc.Arm()
	cc.SetIV(&Bytes{})
	cc.Arm()
	var dd *Crypt
	dd.Lock()
	dd.Generate(NewPassword().FromString("testing"))
	dd.Generate(nil)
	dd.Encrypt(NewLockedBuffer().Rand(32))
	dd.Decrypt(NewBytes().Rand(32))
	dd = new(Crypt)
	dd.Encrypt(NewLockedBuffer().Rand(32))
	dd.Decrypt(NewBytes().Rand(32))
	var ee *Crypt
	ee = new(Crypt)
	ee.Lock()
	var ff *Crypt
	ff.Unlock(nil)
	ff.Disarm()
	ff = new(Crypt)
	ff.Unlock(nil)
	ff.Generate(NewPassword().FromString("abcdef"))
	ff.gcm = nil
	ff.Decrypt(NewBytes().Rand(32))
	ff.Encrypt(NewLockedBuffer().Rand(32))
}
