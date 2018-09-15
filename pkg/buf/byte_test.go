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
	a.Copy(&test)
	b.Copy(&test)
	out := new([]byte)
	a.Bytes(out)
	fmt.Println(*out, a.Error())
	b.Bytes(out)
	fmt.Println(*out, b.Error())
	c.Bytes(out)
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
	codings := new([]string)
	c.ListCodings(codings)
	fmt.Println(codings)
	code := new(string)
	a.GetCoding(code)
	fmt.Println(*code)
	a.SetCoding("hex")
	a.GetCoding(code)
	fmt.Println(*code)
	a.SetCoding("b0rk")
	a.GetCoding(code)
	fmt.Println(*code)
	var err error
	a.SetStatusIf(err)
	fmt.Println(a.Status, a.OK())
	err = errors.New("testing status")
	a.SetStatusIf(err)
	fmt.Println(a.Status, a.OK())
}
