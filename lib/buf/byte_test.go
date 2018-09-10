package buf

import (
	"fmt"
	"gitlab.com/parallelcoin/duo/lib/debug"
	"runtime"
	"testing"
)

func TestEverything(t *testing.T) {
	dbg.Init()
	defer func() {
		if recover() != nil {
			pc, fil, line, _ := runtime.Caller(0)
			fmt.Println(pc, fil, line)
			dbg.D.Close()
		}

	}()
	b := NewByte()
	b.Buf = 42
	b.Coding = "hex"
	b.Status = "testing error"
	want := b.Freeze()

	dbg.Append("", b)
	b.Thaw(want)
	dbg.Append("", b)
	var c *Byte
	dbg.Append("", c)
	if want != b.Freeze() {
		t.Error("Freeze/Thaw test failed: wanted", want, "got", b.Freeze())
	}
	want = `"Type":"*Byte","Value":{ "Buf": 0, "Status": "", "Coding": ""}`
	if want != c.Freeze() {
		t.Error("Freeze/Thaw test failed: wanted", want, "got", b.Freeze())
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
	dbg.Append("", b.Copy(13).(*Byte))
	b = NewByte()
	for i := range testvars {
		dbg.Append("", b.Copy(testvars[i]).(*Byte))
	}
	pByte := &[]byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	vByte := []byte{0, 1, 1, 2, 3, 5, 8, 11, 19}
	dbg.Append("", b.Copy(pByte).(*Byte))
	dbg.Append("", b.Copy(vByte).(*Byte))
	dbg.Append("", dbg.NewNote("Byte.Data() = "+fmt.Sprint(b.Data().(byte))))
	dbg.Append("", b.Free().(*Byte))
	dbg.Append("", struct{ *Byte }{}.Free().(*Byte))
	dbg.D.Close()
}
