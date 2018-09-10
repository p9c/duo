package buf

import (
	"fmt"
	// "gitlab.com/parallelcoin/duo/lib/array"
	"gitlab.com/parallelcoin/duo/lib/debug"
	"runtime"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	y := "*Bytes"
	dbg.Init()
	defer func() {
		if recover() != nil {
			for i := 0; i < 9; i++ {
				pc, fil, line, _ := runtime.Caller(i)
				fmt.Println(pc, fil+":"+fmt.Sprint(line))
			}
			dbg.D.Close()
		}
	}()
	b := NewBytes()
	b.Copy(int64(42))
	b.Coding = "hex"
	b.Status = "testing error"
	want := b.Freeze()
	dbg.Append("*Bytes", b)
	rethawed := b.Thaw(want).(*Bytes)
	got := rethawed.Freeze()
	dbg.Append("*Bytes", rethawed)
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
	dbg.Append("*Bytes", b.Copy(13).(*Bytes))
	b = NewBytes()
	for i := range testvars {
		dbg.Append("*Bytes", b.Copy(testvars[i]).(*Bytes))
	}
	pBytes := &[]byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	vBytes := []byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	c := NewBytes()
	dbg.Append(y, b.Copy(pBytes).(*Bytes))
	dbg.Append(y, c.Copy(vBytes).(*Bytes))
	dbg.Append(y, b.Copy(b).(*Bytes))
	dbg.Append(y, b.Copy(c).(*Bytes))
	dbg.Append("", dbg.NewNote("Bytes.Data() = "+fmt.Sprint(b.Data().(*[]byte))))
	dbg.Append(y, b.Free().(*Bytes))
	dbg.Append(y, b.Null().(*Bytes))
	dbg.Append(y, struct{ *Bytes }{}.Copy(byte(19)).(*Bytes))
	dbg.Append(y, struct{ *Bytes }{}.Free().(*Bytes))
	dbg.Append(y, struct{ *Bytes }{}.Null().(*Bytes))
	dbg.Append("", dbg.NewNote("Bytes.GetCoding() default = "+fmt.Sprint(b.GetCoding())))
	dbg.Append(y, struct{ *Bytes }{}.SetStatus("test status").(*Bytes))
	dbg.Append(y, struct{ *Bytes }{}.SetCoding("not a coding").(*Bytes))
	dbg.Append(y, b.SetCoding("base32").(*Bytes))
	dbg.Append(y, struct{ *Bytes }{})
	var e *Bytes
	dbg.Append(y, e.Thaw(e.Freeze()).(*Bytes))
	dbg.Append(y, b.Thaw(`{`).(*Bytes))
	b.Coding = ""
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Coding == "")))
	dbg.Append(y, b)
	dbg.Append("", dbg.NewNote(b.GetCoding()))
	dbg.Append("", dbg.NewNote(strings.Join(b.ListCodings(), ",")))
	b.Elem(8)
	dbg.Append("", dbg.NewNote(b.Error()))
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Len())))
	dbg.Append(y, b.SetElem(1, 1).(*Bytes))
	dbg.Append(y, b.SetElem(1, 0).(*Bytes))
	dbg.Append(y, b.SetElem(6, 1).(*Bytes))
	dbg.Append(y, b.SetElem(6, 0).(*Bytes))
	dbg.Append(y, b.Copy(pBytes).(*Bytes))
	dbg.Append("", dbg.NewNote(b.SetCoding("base58check").(*Bytes).String()))
	b.Buf = nil
	dbg.Append("", dbg.NewNote("Bytes.Data() "+fmt.Sprint(b.Data())))
	b.Coding = ""
	dbg.Append("", dbg.NewNote("Bytes.GetCoding() "+fmt.Sprint(b.GetCoding())))
	n := NewBytes()
	n.Copy(vBytes)
	dbg.Append("", dbg.NewNote("Elem() oob"+fmt.Sprint(n.Elem(n.Len()+2))))
	dbg.Append("", dbg.NewNote("SetElem() oob"+fmt.Sprint(n.SetElem(n.Len()+2, byte(144)))))
	dbg.Append("", dbg.NewNote("Elem() "+fmt.Sprint(n.Elem(3))))
	dbg.Append(y, n)
	dbg.Append(y, n.Null().(*Bytes))
	n.Buf = &[]byte{}
	dbg.Append("", dbg.NewNote("SetElem() oob"+fmt.Sprint(n.SetElem(n.Len()+2, byte(144)))))
	n.Buf = &[]byte{}
	dbg.Append("", dbg.NewNote("SetElem() zlb"+fmt.Sprint(n.SetElem(0, byte(144)))))
	dbg.Append("", dbg.NewNote("Elem() "+fmt.Sprint(n.Elem(0))))
	dbg.Append("", dbg.NewNote("Elem() "+fmt.Sprint(n.Elem(-1))))
	dbg.Append(y, n)
	dbg.Append(y, n.Copy([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}).(*Bytes))
	dbg.Append(y, n.Copy([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}).(*Bytes))
	dbg.D.Close()
}
