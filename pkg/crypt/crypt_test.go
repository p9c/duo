package crypt

import (
	"fmt"
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
	h := NewCrypt().Load(NewBytes().Rand(13).(*Bytes))
	fmt.Println(h.crypt.Buf())
	fmt.Println(g.IsArmed(), f.IsArmed())
	fmt.Println(g.IsUnlocked(), f.IsUnlocked())
	fmt.Println(a.IsLoaded(), f.IsLoaded(), g.IsLoaded(), h.IsLoaded())
	x := f.Load(NewBytes().Rand(11).(*Bytes))
	fmt.Println(*x.Crypt().Buf().(*[]byte))
	fmt.Println(x.Crypt())
	x.Crypt().Null().Null()
	fmt.Println(x)
	fmt.Println(*x.SetIV(NewBytes().Rand(12).(*Bytes)).IV())
	fmt.Println(x.SetIV(NewBytes().Rand(11).(*Bytes)))
	var z *Crypt
	fmt.Println(z)
	fmt.Println(z.SetError("nothing"))
	var n *Crypt
	fmt.Println("SetIV()", n.SetIV(NewBytes().Rand(12).(*Bytes)).IV().Buf())
	var m *Crypt
	fmt.Println(m.SetIV(nil).IV())
	fmt.Println(*m.SetRandomIV().IV().Buf().(*[]byte))
	fmt.Println("valid receiver")
	pass := NewPassword().FromString("abcdef")
	C := NewCrypt().SetRandomIV().Generate(pass)
	C.Decrypt(NewBytes().Rand(64).(*Bytes))
	fmt.Println("Generate()", C.String())
	C.Disarm()
	fmt.Println("Disarm()", C.String())
	C.Decrypt(NewBytes().Rand(64).(*Bytes))
	C.Lock()
	fmt.Println("Lock()", C.String())
	C.Decrypt(NewBytes().Rand(64).(*Bytes))
	pass = NewPassword().FromString("abcdef")
	C.Unlock(pass)
	fmt.Println("Unlock()", C.String())
	C.Decrypt(NewBytes().Rand(64).(*Bytes))
	C.Arm()
	fmt.Println("Arm()", C.String())
	C.Decrypt(NewBytes().Rand(64).(*Bytes))
	ss := []byte("this is a test")
	fmt.Println("plaintext '" + string(ss) + "'")
	fmt.Println("bytes ", ss)
	bb := C.Encrypt(NewLockedBuffer().Load(&ss).(*LockedBuffer))
	fmt.Println("ciphertext", *bb.Buf().(*[]byte))
	BB := C.Decrypt(bb)
	fmt.Println("recovertext", *BB.Buf().(*[]byte))
	fmt.Println("recover as string '" + string(*BB.Buf().(*[]byte)) + "'")
	C.Lock()
	fmt.Println("Unlock()", C.String())
	var cc *Crypt
	cc.Arm()
	cc = new(Crypt)
	cc.Arm()
	cc.password = NewPassword().FromString("abcdef")
	cc.Arm()
	cc.Load(NewBytes().Rand(32).(*Bytes))
	cc.Arm()
	cc.SetIV(&Bytes{})
	cc.Arm()
	var dd *Crypt
	dd.Lock()
	dd.Generate(NewPassword().FromString("testing"))
	dd.Generate(nil)
	dd.Encrypt(NewLockedBuffer().Rand(32).(*LockedBuffer))
	dd.Decrypt(NewBytes().Rand(32).(*Bytes))
	dd = new(Crypt)
	dd.Encrypt(NewLockedBuffer().Rand(32).(*LockedBuffer))
	dd.Decrypt(NewBytes().Rand(32).(*Bytes))
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
	ff.Decrypt(NewBytes().Rand(32).(*Bytes))
	ff.Encrypt(NewLockedBuffer().Rand(32).(*LockedBuffer))
	fffff := []byte{}
	ff.Load(NewBytes().Load(fffff).(*Bytes))
	ff.Arm()

}
