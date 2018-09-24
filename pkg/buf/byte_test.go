package buf

import (
	"errors"
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	test := []byte("Test String")
	a := NewByte()
	b := new(Byte)
	var c *Byte
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
	b.SetStatus("")
	b.SetStatus("test")
	a.Copy(&test)
	b.Copy(&test)
	b.Copy(&[]byte{})
	out := a.Bytes()
	fmt.Println(*out, a.Error())
	fmt.Println(c.Bytes())
	c.Copy(&test)
	fmt.Println(c.Bytes())
	fmt.Println(c.Copy(nil))
	d := make([]byte, 0)
	a.Copy(&d)
	fmt.Println(d, a.Error())
	fmt.Println(c.Zero())
	a.Zero()
	fmt.Println(a.Error())
	a.Free()
	fmt.Println(a.Error())
	fmt.Println(c.Free())
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
	fmt.Println(a.State)
	fmt.Println(a.OK())
	err = errors.New("testing status")
	a.SetStatusIf(err)
	fmt.Println(a.State, a.OK())
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
	f := NewByte()
	f.Copy(&test)
	_ = f.GetElem(14)
	f.SetElem(14, byt)
	_ = f.GetElem(1)
	f.SetElem(1, byt)
	f.SetElem(1, &[]byte{})
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
	ij := []byte("invalid [] {}")
	fmt.Println(NewByte().Thaw(&ij).(*Byte).String())
}
