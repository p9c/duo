package buf

import (
	"fmt"
	"gitlab.com/parallelcoin/duo/lib/debug"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
)

func TestSecure(t *testing.T) {
	y := "*Secure"
	dbg.Init()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
			debug.PrintStack()
			stack := strconv.Quote(string(debug.Stack()))
			stack = stack[1:]
			stack = stack[:len(stack)-2]
			stack = strings.Replace(string(stack), `\n`, ` `, -1)
			stack = strings.Replace(string(stack), `\t`, ``, -1)
			stack = strings.Replace(string(stack), `\"`, "`", -1)
			stack = strings.Replace(string(stack), `\`, ``, -1)
			dbg.Append("", dbg.NewNote(stack))
			dbg.D.Close()
		}
	}()
	b := NewSecure()
	b.Copy(42)
	b.Status = "testing error"
	b.Coding = "hex"
	want := b.Freeze()
	dbg.Append("*Secure", b)
	rethawed := b.Thaw(want).(*Secure)
	got := rethawed.Freeze()
	dbg.Append("*Secure", rethawed)
	if got != want {
		t.Error("Freeze/Thaw test failed:\nwanted", want, "\n   got", got)
	}
	testvars := []interface{}{
		nil,
		int(42),
		uint(42),
		byte(42),
		int8(42),
		uint16(42),
		int16(42),
		uint32(42),
		int32(42),
		uint64(42),
		int64(42),
	}
	dbg.Append("*Secure", b.Copy(13).(*Secure))
	b = NewSecure()
	for i := range testvars {
		dbg.Append("*Secure", b.Copy(testvars[i]).(*Secure))
	}
	pSecure := &[]byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	vSecure := []byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	c := NewSecure()
	dbg.Append(y, b.Copy(pSecure).(*Secure))
	dbg.Append(y, c.Copy(vSecure).(*Secure))
	dbg.Append(y, b.Copy(b).(*Secure))
	dbg.Append(y, b.Copy(c).(*Secure))
	dbg.Append("", dbg.NewNote("Secure.Data() = "+fmt.Sprint(b.Data().(*[]byte))))
	dbg.Append(y, b.Free().(*Secure))
	dbg.Append(y, b.Null().(*Secure))
	dbg.Append(y, struct{ *Secure }{}.Copy(byte(19)).(*Secure))
	dbg.Append(y, struct{ *Secure }{}.Free().(*Secure))
	dbg.Append(y, struct{ *Secure }{}.Null().(*Secure))
	dbg.Append("", dbg.NewNote("Secure.GetCoding() default = "+fmt.Sprint(b.GetCoding())))
	dbg.Append(y, struct{ *Secure }{}.SetStatus("test status").(*Secure))
	dbg.Append(y, struct{ *Secure }{}.SetCoding("not a coding").(*Secure))
	dbg.Append(y, b.SetCoding("base32").(*Secure))
	dbg.Append(y, struct{ *Secure }{})
	var e *Secure
	dbg.Append(y, e.Thaw(e.Freeze()).(*Secure))
	dbg.Append(y, b.Thaw(`{`).(*Secure))
	b.Coding = ""
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Coding == "")))
	dbg.Append(y, b)
	dbg.Append("", dbg.NewNote(b.GetCoding()))
	dbg.Append("", dbg.NewNote(strings.Join(b.ListCodings(), ",")))
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Elem(8))))
	dbg.Append("", dbg.NewNote(b.Error()))
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Len())))
	b = NewSecure().Copy([]byte("this is a test")).(*Secure).SetCoding("string").(*Secure)
	dbg.Append(y, b.SetElem(1, 1).(*Secure))
	// bb := new(byte)
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Buf)))
	dbg.Append(y, b.SetElem(1, 0).(*Secure))
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Buf)))
	dbg.Append(y, b.SetElem(6, 1).(*Secure))
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Buf)))
	dbg.Append(y, b.SetElem(6, 0).(*Secure))
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Buf)))
	dbg.Append(y, b.Copy(pSecure).(*Secure))
	dbg.Append("", dbg.NewNote(b.SetCoding("base58check").(*Secure).String()))
	b.Buf = nil
	dbg.Append("", dbg.NewNote("Secure.Data() "+fmt.Sprint(b.Data())))
	b.Coding = ""
	dbg.Append("", dbg.NewNote("Secure.GetCoding() "+fmt.Sprint(b.GetCoding())))
	n := NewSecure()
	n.Copy(vSecure)
	dbg.Append("", dbg.NewNote("Elem() oob"+fmt.Sprint(n.Elem(n.Len()+2))))
	dbg.Append("", dbg.NewNote("SetElem() oob"+fmt.Sprint(n.SetElem(n.Len()+2, byte(144)))))
	dbg.Append("", dbg.NewNote("Elem() "+fmt.Sprint(n.Elem(3))))
	dbg.Append(y, n)
	dbg.Append(y, n.Null().(*Secure))
	dbg.Append("", dbg.NewNote("SetElem() oob"+fmt.Sprint(n.SetElem(n.Len()+2, byte(144)))))
	dbg.Append("", dbg.NewNote("SetElem() zlb"+fmt.Sprint(n.SetElem(0, byte(144)))))
	dbg.Append("", dbg.NewNote("Elem() "+fmt.Sprint(n.Elem(0))))
	dbg.Append("", dbg.NewNote("Elem() "+fmt.Sprint(n.Elem(-1))))
	dbg.Append(y, n)
	dbg.Append(y, n.Copy([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}).(*Secure))
	dbg.Append(y, n.Copy([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}).(*Secure))
	dbg.D.Close()
}
