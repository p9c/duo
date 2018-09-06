package secbuf

import (
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"testing"
)

func TestLockedBuffer(t *testing.T) {
	a := new(SecBuf)
	A := []byte("test")
	fmt.Println(a.SetCoding("string").(*SecBuf).Load(&A).String())
	fmt.Println(a.String())
	fmt.Println("a", a.Buf(), a.String(), a.Error())
	b := new(SecBuf)
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
	var d *SecBuf
	fmt.Println("copy to unallocated b", b, "d", d.Copy(b))
	fmt.Println("Should have an error, copying nil pointer: '" + b.Copy(d).Error() + "'")
	fmt.Println("Random bytes", b.Rand(32).Buf())
	fmt.Println("Chaining - New at end should mean empty at end", b.New(32).Rand(12).New(5).Buf())
	fmt.Println("Struct literal (VERY BAD! Dangles afterwards)", struct{ *SecBuf }{}.New(12))
	fmt.Println("Copying to ourself, should be an error '" + a.Copy(a).Error() + "'")
	fmt.Println("Copying to ourself, should be an error '" + a.Copy(&SecBuf{}).Error() + "'")
	fmt.Println("Getting length of struct literal should be 0:", struct{ *SecBuf }{}.Len())
	e := new(SecBuf)
	fmt.Println("Getting length on unallocated buffer (should be 0):", e.Len())
	fmt.Println("Nulling empty struct literal should be nil, false nil:", struct{ *SecBuf }{}.Null())
	fmt.Println("Getting buf fromempty struct literal should be empty:", struct{ *SecBuf }{}.Buf())
	fmt.Println("Rand on empty struct literal:", struct{ *SecBuf }{}.Rand(12).Buf())
	fmt.Println("Rand on empty struct literal with zero length: '"+struct{ *SecBuf }{}.Rand(0).Error(), "'")
	fmt.Println("New on empty struct literal with zero length: '" + struct{ *SecBuf }{}.New(0).Error() + "'")
	fmt.Println("Buf() on empty struct literal:", struct{ *SecBuf }{}.Null().Buf())
	A = []byte("testing")
	fmt.Println("Load() on empty struct literal:", struct{ *SecBuf }{}.Load(&A).Buf())
	fmt.Println("Link() on empty struct literal:", struct{ *SecBuf }{}.Link(b).Buf())
	fmt.Println("Move() on empty struct literal:", struct{ *SecBuf }{}.Move(b).Buf())
	f := []byte{}
	fmt.Println("Load([]byte{}) on empty struct literal:", struct{ *SecBuf }{}.Load(&f).Error())
	fmt.Println("Load(nil) on empty struct literal:", struct{ *SecBuf }{}.Load(nil))
	g := New().Rand(13)
	fmt.Println("New().Rand(13)", g, g.Buf())
	fmt.Println(New().New(13).Null())
	New().Free()
	var n *SecBuf
	fmt.Println("should be nil receiver:", n.Error())
	e.Rand(32)
	fmt.Println("Move(nil)", n.Move(nil))
	fmt.Println("nil SetError()", n.SetError("testing"))
	fmt.Println("JSON", e.String())
	var m *SecBuf
	m.UnsetError()
	m.SetElem(0, byte(100))
	m.Elem(0)
	fmt.Println("JSON", m.String())
	m.MarshalJSON()
	m.Load(bytes.New().Rand(32).Buf())
	m = new(SecBuf)
	m.SetCoding("string")
	m.buf, m.err = memguard.NewImmutableRandom(32)
	m.MarshalJSON()
	var oo *SecBuf
	oo.SetCoding("decimal")
	var pp *SecBuf
	pp.SetCoding("binary")
	New().SetElem(0, byte(100))
	New().Copy(nil).Elem(0)
	New().Load(&[]byte{'t', 'e', 's', 't'}).(*SecBuf).Free()
	yy := "test"
	YY := []byte(yy)
	xx := New().Load(&YY)
	fmt.Println(xx.String())
	xbuf, xerr := memguard.NewMutable(12)
	xx.(*SecBuf).buf = xbuf
	xx.(*SecBuf).err = xerr
	fmt.Println(xx.SetCoding("hex").(*SecBuf).String())
	xx.Free()
	fmt.Println(xx.String())
	fmt.Println(xx.Coding())
	fmt.Println(xx.Codes())
	xx.(*SecBuf).coding = len(def.CodingTypes) + 4
	xx.Coding()
	xx.Cap()
	xx.Purge().(*SecBuf).New(23).Cap()
	xx.Purge()
	fmt.Println(xx.Rand(23).SetCoding("decimal").(*SecBuf).String())
	xx.(*SecBuf).coding = len(def.CodingTypes) + 5
	fmt.Println(xx.Rand(23).String())
	fmt.Println(xx.Rand(23).SetCoding("base64").(*SecBuf).String())
	xx.SetElem(-1, byte(123))
	xx.Rand(23).SetElem(24, byte(23))
	xx.Rand(23).SetElem(24, nil)
	var xyz *SecBuf
	xyz.Purge()
	xyz.Coding()
	New().Free()
	nx := make([]byte, 0)
	Nx := &nx
	fmt.Println(Nx)
	xyz.Copy(New().Load(Nx))
	var abcd *SecBuf
	abcd.Free()
	abcd = abcd.Rand(23).(*SecBuf)
	var xxxx *SecBuf
	xxxx.Link(abcd)
	xxxx.Link(nil)
	xxxx.Link(bytes.New().Rand(23))
	New().New(24).Link(New().New(23))
}
