package lockedbuffer

import (
	"fmt"
	"github.com/awnumar/memguard"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
	"testing"
)

func TestLockedBuffer(t *testing.T) {
	a := new(LockedBuffer)
	A := []byte("test")
	fmt.Println(a.SetCoding("string").(*LockedBuffer).Load(&A).String())
	fmt.Println(a.String())
	fmt.Println("a", a.Buf(), a.String(), a.Error())
	b := new(LockedBuffer)
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
	var d *LockedBuffer
	fmt.Println("copy to unallocated b", b, "d", d.Copy(b))
	fmt.Println("Should have an error, copying nil pointer: '" + b.Copy(d).Error() + "'")
	fmt.Println("Random bytes", b.Rand(32).Buf())
	fmt.Println("Chaining - New at end should mean empty at end", b.New(32).Rand(12).New(5).Buf())
	fmt.Println("Struct literal (VERY BAD! Dangles afterwards)", struct{ *LockedBuffer }{}.New(12))
	fmt.Println("Copying to ourself, should be an error '" + a.Copy(a).Error() + "'")
	fmt.Println("Copying to ourself, should be an error '" + a.Copy(&LockedBuffer{}).Error() + "'")
	fmt.Println("Getting length of struct literal should be 0:", struct{ *LockedBuffer }{}.Len())
	e := new(LockedBuffer)
	fmt.Println("Getting length on unallocated buffer (should be 0):", e.Len())
	fmt.Println("Nulling empty struct literal should be nil, false nil:", struct{ *LockedBuffer }{}.Null())
	fmt.Println("Getting buf fromempty struct literal should be empty:", struct{ *LockedBuffer }{}.Buf())
	fmt.Println("Rand on empty struct literal:", struct{ *LockedBuffer }{}.Rand(12).Buf())
	fmt.Println("Rand on empty struct literal with zero length: '"+struct{ *LockedBuffer }{}.Rand(0).Error(), "'")
	fmt.Println("New on empty struct literal with zero length: '" + struct{ *LockedBuffer }{}.New(0).Error() + "'")
	fmt.Println("Buf() on empty struct literal:", struct{ *LockedBuffer }{}.Null().Buf())
	A = []byte("testing")
	fmt.Println("Load() on empty struct literal:", struct{ *LockedBuffer }{}.Load(&A).Buf())
	fmt.Println("Link() on empty struct literal:", struct{ *LockedBuffer }{}.Link(b).Buf())
	fmt.Println("Move() on empty struct literal:", struct{ *LockedBuffer }{}.Move(b).Buf())
	f := []byte{}
	fmt.Println("Load([]byte{}) on empty struct literal:", struct{ *LockedBuffer }{}.Load(&f).Error())
	fmt.Println("Load(nil) on empty struct literal:", struct{ *LockedBuffer }{}.Load(nil))
	g := NewLockedBuffer().Rand(13)
	fmt.Println("NewLockedBuffer().Rand(13)", g, g.Buf())
	fmt.Println(NewLockedBuffer().New(13).Null())
	NewLockedBuffer().Free()
	var n *LockedBuffer
	fmt.Println("should be nil receiver:", n.Error())
	e.Rand(32)
	fmt.Println("Move(nil)", n.Move(nil))
	fmt.Println("nil SetError()", n.SetError("testing"))
	fmt.Println("JSON", e.String())
	var m *LockedBuffer
	m.UnsetError()
	m.SetElem(0, byte(100))
	m.Elem(0)
	fmt.Println("JSON", m.String())
	m.MarshalJSON()
	m.Load(NewBytes().Rand(32).Buf())
	m = new(LockedBuffer)
	m.SetCoding("string")
	m.buf, m.err = memguard.NewImmutableRandom(32)
	m.MarshalJSON()
	var oo *LockedBuffer
	oo.SetCoding("decimal")
	var pp *LockedBuffer
	pp.SetCoding("binary")
	NewLockedBuffer().SetElem(0, byte(100))
	NewLockedBuffer().Copy(nil).Elem(0)
	NewLockedBuffer().Load(&[]byte{'t', 'e', 's', 't'}).(*LockedBuffer).Free()
	yy := "test"
	YY := []byte(yy)
	xx := NewLockedBuffer().Load(&YY)
	fmt.Println(xx.String())
	xbuf, xerr := memguard.NewMutable(12)
	xx.(*LockedBuffer).buf = xbuf
	xx.(*LockedBuffer).err = xerr
	fmt.Println(xx.SetCoding("hex").(*LockedBuffer).String())
	xx.Free()
	fmt.Println(xx.String())
	fmt.Println(xx.Coding())
	fmt.Println(xx.Codes())
	xx.(*LockedBuffer).coding = len(CodeType) + 4
	xx.Coding()
	xx.Cap()
	xx.Purge().(*LockedBuffer).New(23).Cap()
	xx.Purge()
	fmt.Println(xx.Rand(23).SetCoding("decimal").(*LockedBuffer).String())
	xx.(*LockedBuffer).coding = len(CodeType) + 5
	fmt.Println(xx.Rand(23).String())
	fmt.Println(xx.Rand(23).SetCoding("base64").(*LockedBuffer).String())
	xx.SetElem(-1, byte(123))
	xx.Rand(23).SetElem(24, byte(23))
	xx.Rand(23).SetElem(24, nil)
	var xyz *LockedBuffer
	xyz.Purge()
	xyz.Coding()
	NewLockedBuffer().Free()
	nx := make([]byte, 0)
	Nx := &nx
	fmt.Println(Nx)
	xyz.Copy(NewLockedBuffer().Load(Nx))
	var abcd *LockedBuffer
	abcd.Free()
	abcd = abcd.Rand(23).(*LockedBuffer)
	var xxxx *LockedBuffer
	xxxx.Link(abcd)
	xxxx.Link(nil)
	xxxx.Link(NewBytes().Rand(23))
	NewLockedBuffer().New(24).Link(NewLockedBuffer().New(23))
}
