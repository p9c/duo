package buf

import (
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"testing"
)

func TestLockedBuffer(t *testing.T) {
	a := new(Fenced)
	A := []byte("test")
	fmt.Println(a.SetCoding("string").(*Fenced).Load(&A).String())
	fmt.Println(a.String())
	fmt.Println("a", a.Buf(), a.String(), a.Error())
	b := new(Fenced)
	fmt.Println("copy a", a.Buf(), "to b", b.Buf())
	fmt.Println(b.Copy(a))
	fmt.Println("copy a", a.Buf(), "to b", b.Buf())
	fmt.Println("before move a", *a, "b", *b)
	fmt.Println(b.Move(a))
	fmt.Println("after move a", *a, "b", *b)
	a.Link(b)
	fmt.Println("link b to (empty) a", a.Buf(), b.Buf())
	a.SetElem(0, byte(0))
	fmt.Println("now both the same memory (changed byte zero of a only)", a.Buf(), b.Buf())
	var d *Fenced
	fmt.Println("copy to unallocated b", b, "d", d.Copy(b))
	fmt.Println("Should have an error, copying nil pointer: '" + b.Copy(d).Error() + "'")
	fmt.Println("Random bytes", b.Rand(32).Buf())
	fmt.Println("Chaining - New at end should mean empty at end", b.New(32).Rand(12).New(5).Buf())
	fmt.Println("Struct literal (VERY BAD! Dangles afterwards)", struct{ *Fenced }{}.New(12))
	fmt.Println("Copying to ourself, should be an error '" + a.Copy(a).Error() + "'")
	fmt.Println("Copying to ourself, should be an error '" + a.Copy(&Fenced{}).Error() + "'")
	fmt.Println("Getting length of struct literal should be 0:", struct{ *Fenced }{}.Len())
	e := new(Fenced)
	fmt.Println("Getting length on unallocated buffer (should be 0):", e.Len())
	fmt.Println("Nulling empty struct literal should be nil, false nil:", struct{ *Fenced }{}.Null())
	fmt.Println("Getting buf fromempty struct literal should be empty:", struct{ *Fenced }{}.Buf())
	fmt.Println("Rand on empty struct literal:", struct{ *Fenced }{}.Rand(12).Buf())
	fmt.Println("Rand on empty struct literal with zero length: '"+struct{ *Fenced }{}.Rand(0).Error(), "'")
	fmt.Println("New on empty struct literal with zero length: '" + struct{ *Fenced }{}.New(0).Error() + "'")
	fmt.Println("Buf() on empty struct literal:", struct{ *Fenced }{}.Null().Buf())
	A = []byte("testing")
	fmt.Println("Load() on empty struct literal:", struct{ *Fenced }{}.Load(&A).Buf())
	fmt.Println("Link() on empty struct literal:", struct{ *Fenced }{}.Link(b).Buf())
	fmt.Println("Move() on empty struct literal:", struct{ *Fenced }{}.Move(b).Buf())
	f := []byte{}
	fmt.Println("Load([]byte{}) on empty struct literal:", struct{ *Fenced }{}.Load(&f).Error())
	fmt.Println("Load(nil) on empty struct literal:", struct{ *Fenced }{}.Load(nil))
	g := NewFenced().Rand(13)
	fmt.Println("NewFenced().Rand(13)", g, g.Buf())
	fmt.Println(NewFenced().New(13).Null())
	NewFenced().Free()
	var n *Fenced
	fmt.Println("should be nil receiver:", n.Error())
	e.Rand(32)
	fmt.Println("Move(nil)", n.Move(nil))
	fmt.Println("nil SetError()", n.SetError("testing"))
	fmt.Println("JSON", e.String())
	var m *Fenced
	m.UnsetError()
	m.SetElem(0, byte(100))
	m.Elem(0)
	fmt.Println("JSON", m.String())
	m.MarshalJSON()
	m.Load(NewUnsafe().Rand(32).Buf())
	m = new(Fenced)
	m.SetCoding("string")
	m.buf, m.err = memguard.NewImmutableRandom(32)
	m.MarshalJSON()
	var oo *Fenced
	oo.SetCoding("decimal")
	var pp *Fenced
	pp.SetCoding("binary")
	NewFenced().SetElem(0, byte(100))
	NewFenced().Copy(nil).Elem(0)
	NewFenced().Load(&[]byte{'t', 'e', 's', 't'}).(*Fenced).Free()
	yy := "test"
	YY := []byte(yy)
	xx := NewFenced().Load(&YY)
	fmt.Println(xx.String())
	xbuf, xerr := memguard.NewMutable(12)
	xx.(*Fenced).buf = xbuf
	xx.(*Fenced).err = xerr
	fmt.Println(xx.SetCoding("hex").(*Fenced).String())
	xx.Free()
	fmt.Println(xx.String())
	fmt.Println(xx.Coding())
	fmt.Println(xx.Codes())
	xx.(*Fenced).coding = len(def.StringCodingTypes) + 4
	xx.Coding()
	xx.Cap()
	xx.(*Fenced).New(23).Cap()
	fmt.Println(xx.Rand(23).SetCoding("decimal").(*Fenced).String())
	xx.(*Fenced).coding = len(def.StringCodingTypes) + 5
	fmt.Println(xx.Rand(23).String())
	fmt.Println(xx.Rand(23).SetCoding("base64").(*Fenced).String())
	xx.SetElem(-1, byte(123))
	xx.Rand(23).SetElem(24, byte(23))
	xx.Rand(23).SetElem(24, nil)
	var xyz *Fenced
	xyz.Coding()
	NewFenced().Free()
	nx := make([]byte, 0)
	Nx := &nx
	fmt.Println(Nx)
	xyz.Copy(NewFenced().Load(Nx))
	var abcd *Fenced
	abcd.Free()
	abcd = abcd.Rand(23).(*Fenced)
	var xxxx *Fenced
	xxxx.Link(abcd)
	xxxx.Link(nil)
	xxxx.Link(NewUnsafe().Rand(23))
	NewFenced().New(24).Link(NewFenced().New(23))
}
