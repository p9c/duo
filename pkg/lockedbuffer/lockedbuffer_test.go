package lockedbuffer

import (
	"fmt"
	"testing"
)

func TestLockedBuffer(t *testing.T) {
	a := new(LockedBuffer)
	A := []byte("test")
	a.Load(&A)
	fmt.Println("a", a.Buf())
	b := new(LockedBuffer)
	b.Copy(a)
	fmt.Println("copy a", a.Buf(), "to b", b.Buf())
	fmt.Println("before move a", *a, "b", *b)
	b.Move(a)
	fmt.Println("after move a", *a, "b", *b)
	a.Link(b)
	fmt.Println("link b to (empty) a", a.Buf(), b.Buf())
	c := a.Buf()
	C := *c
	C[0] = 1
	fmt.Println("now both the same memory (changed byte zero of a only)", a.Buf(), b.Buf())
	var d *LockedBuffer
	fmt.Println("copy to unallocated b", b, "d", d.Copy(b))
	fmt.Println("Should have an error, copying nil pointer: '" + b.Copy(d).Error() + "'")
	fmt.Println("Random bytes", b.Rand(32).Buf())
	fmt.Println("Chaining - New at end should mean empty at end", b.New(32).Rand(12).New(5).Buf())
	fmt.Println("Struct literal (VERY BAD! Dangles afterwards)", struct{ *LockedBuffer }{}.New(12))
	fmt.Println("Copying to ourself, should be an error '" + a.Copy(a).Error() + "'")
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
}
