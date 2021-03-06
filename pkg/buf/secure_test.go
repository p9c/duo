package buf

import (
	"errors"
	"fmt"
	"testing"
)

func TestSecure(t *testing.T) {
	test := []byte("Test String")
	a := NewSecure()
	b := new(Secure)
	var c *Secure
	c.SetStatus("test")
	c.SetStatusIf(errors.New("test"))
	c.UnsetStatus()
	c.OK()
	c.SetElem(1, 1)
	c.SetCoding(*c.GetCoding())
	var byt *byte
	byt = c.GetElem(1).(*byte)
	c.Len()
	b.SetElem(1, 1)
	byt = b.GetElem(1).(*byte)
	b.Len()
	a.Copy(&test)
	b.Copy(&test)
	b.Copy(&[]byte{})
	out := a.Bytes()
	fmt.Println(*out, a.Error())
	out = b.Bytes()
	fmt.Println(*out, b.Error())
	out = c.Bytes()
	fmt.Println(*out, c.Error())
	c.Copy(&test)
	fmt.Println(*out, c.Error())
	c.Copy(nil)
	fmt.Println(*out, c.Error())
	d := make([]byte, 0)
	a.Copy(&d)
	fmt.Println(d, a.Error())
	c.Zero()
	fmt.Println(c.Error())
	a.Free()
	fmt.Println(a.Error())
	a.Zero()
	fmt.Println(a.Error())
	c.Free()
	fmt.Println(c.Error())
	codings := c.ListCodings()
	fmt.Println(codings)
	code := a.GetCoding()
	fmt.Println(*code)
	a.SetCoding("hex")
	code = a.GetCoding()
	fmt.Println(*code)
	a.SetCoding("b0rk")
	code = a.GetCoding()
	fmt.Println(*code)
	var err error
	a.SetStatusIf(err)
	fmt.Println(a.Status, a.OK())
	err = errors.New("testing status")
	a.SetStatusIf(err)
	fmt.Println(a.Status, a.OK())
	fmt.Println(a.String())
	fmt.Println(b.String())
	fmt.Println(c.String())
	testtypes := []string{"bytes", "string", "decimal", "hex", "base32", "base58check", "base64"}
	for i := range testtypes {
		b.SetCoding(testtypes[i])
		fmt.Println(testtypes[i], b.String())
	}
	out = b.Freeze()
	fmt.Println("not nil", string(*out))
	b.Thaw(out)
	fmt.Println(b.String())
	out = c.Freeze()
	fmt.Println("nil", string(*out))
	c.Thaw(out)
	fmt.Println(c.String())
	test = []byte("Test String")
	f := NewSecure()
	f.Copy(&test)
	_ = f.GetElem(24)
	f.SetElem(24, byt)
	_ = f.GetElem(1)
	f.Copy(&[]byte{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21})
	f.SetElem(1, byte(8))
	f.SetElem(f.Len()+5, byte(0))
	f.SetElem(1, &[]byte{})
	many := []interface{}{
		byte(100),
		int8(100),
		int(100),
		uint(100),
		uint16(100),
		int16(100),
		uint32(100),
		int32(100),
		uint64(100),
		int64(100),
	}
	for i := range many {
		f.SetElem(3, many[i])
	}
	fmt.Println(f.Free())
	fmt.Println(f.Bytes())
	fmt.Println(f.IsEqual(f.Bytes()))
	fmt.Println(a.IsEqual(a.Bytes()))
	fmt.Println(a.IsEqual(&[]byte{}))
	fmt.Println(a.IsEqual(nil))
	fmt.Println(c.IsEqual(nil))
	bt := &[]byte{100, 112, 134, 234, 22, 151}
	fmt.Println(a.IsEqual(bt))
	a.Copy(bt)
	fmt.Println(a.IsEqual(a.Bytes()))
	ct := &[]byte{100, 112, 134, 234, 22, 51}
	fmt.Println(a.IsEqual(ct))
	fmt.Println(c.Rand(23).Bytes())
	fmt.Println(c.Free().(*Secure).Rand(0))
	fmt.Println(a.Rand(23))
	fmt.Println(a.Rand(23))
	ij := []byte("invalid [] {}")
	fmt.Println(NewByte().Thaw(&ij).(*Byte).String())
}
