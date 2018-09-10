package buf

import (
	"fmt"
	"gitlab.com/parallelcoin/duo/lib/array"
	"gitlab.com/parallelcoin/duo/lib/debug"
	"runtime"
	"strings"
	"testing"
)

func TestEverything(t *testing.T) {
	y := "*Byte"
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
	b := NewByte()
	b.Buf = 42
	b.Coding = "hex"
	b.Status = "testing error"
	want := b.Freeze()
	dbg.Append("*Byte", b)
	rethawed := b.Thaw(want).(*Byte)
	got := rethawed.Freeze()
	dbg.Append("*Byte", rethawed)
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
	dbg.Append("*Byte", b.Copy(13).(*Byte))
	b = NewByte()
	for i := range testvars {
		dbg.Append("*Byte", b.Copy(testvars[i]).(*Byte))
	}
	pByte := &[]byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	vByte := []byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	c := NewByte()

	dbg.Append(y, b.Copy(pByte).(*Byte))
	dbg.Append(y, c.Copy(vByte).(*Byte))
	dbg.Append(y, b.Copy(b).(*Byte))
	dbg.Append(y, b.Copy(c).(*Byte))
	dbg.Append("", dbg.NewNote("Byte.Data() = "+fmt.Sprint(b.Data().(byte))))
	dbg.Append(y, b.Free().(*Byte))
	dbg.Append(y, b.Null().(*Byte))
	dbg.Append(y, struct{ *Byte }{}.Copy(byte(19)).(*Byte))
	dbg.Append(y, struct{ *Byte }{}.Free().(*Byte))
	dbg.Append(y, struct{ *Byte }{}.Null().(*Byte))
	dbg.Append("", dbg.NewNote("Byte.GetCoding() default = "+fmt.Sprint(b.GetCoding())))
	dbg.Append(y, struct{ *Byte }{}.SetStatus("test status").(*Byte))
	dbg.Append(y, struct{ *Byte }{}.SetCoding("not a coding").(*Byte))
	dbg.Append(y, b.SetCoding("base32").(*Byte))
	dbg.Append(y, struct{ *Byte }{})
	var e *Byte
	dbg.Append(y, e.Thaw(e.Freeze()).(*Byte))
	dbg.Append(y, b.Thaw(`{`).(*Byte))
	b.Coding = ""
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Coding == "")))
	dbg.Append(y, b)
	dbg.Append("", dbg.NewNote(b.GetCoding()))
	dbg.Append("", dbg.NewNote(strings.Join(b.ListCodings(), ",")))
	var bits string
	for i := 0; i < 8; i++ {
		if i != 0 {
			bits += ", "
		}
		bits += fmt.Sprint(b.Copy(131).(arr.Array).Elem(i).(byte))
	}
	dbg.Append("", dbg.NewNote(bits)) // from the value 131
	b.Elem(8)
	dbg.Append("", dbg.NewNote(b.Error()))
	dbg.Append("", dbg.NewNote(fmt.Sprint(b.Len())))
	dbg.Append(y, b.SetElem(1, 1).(*Byte))
	dbg.Append(y, b.SetElem(1, 0).(*Byte))
	dbg.Append(y, b.SetElem(6, 1).(*Byte))
	dbg.Append(y, b.SetElem(6, 0).(*Byte))
	dbg.Append("", dbg.NewNote(b.SetCoding("base58check").(*Byte).String()))
	dbg.D.Close()
}
