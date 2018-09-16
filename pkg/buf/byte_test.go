package buf

import (
	"errors"
	"fmt"
	"testing"
)

func TestByte(t *testing.T) {
	test := []byte("Test String")
	a := NewByte()
	b := new(Byte)
	var c *Byte
	c.SetStatus("test")
	c.SetStatusIf(errors.New("test"))
	c.UnsetStatus()
	c.OK()
	a.Copy(&test)
	b.Copy(&test)
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
	testtypes := []string{"bytes", "string", "hex", "base32", "base58check", "base64"}
	for i := range testtypes {
		b.SetCoding(testtypes[i])
		fmt.Println(testtypes[i], b.String())
	}
	b.Freeze(out)
	fmt.Println("not nil", string(*out))
	b.Thaw(out)
	fmt.Println(b.String())
	c.Freeze(out)
	fmt.Println("nil", string(*out))
	c.Thaw(out)
	fmt.Println(c.String())
}
