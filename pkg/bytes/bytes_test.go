package bytes

import (
	"fmt"
	// . "gitlab.com/parallelcoin/duo/pkg/interfaces"
	"testing"
)

func TestBytes(t *testing.T) {
	al := &Bytes{}
	an := NewBytes()
	var ay *Bytes
	ab := make([]byte, 0)
	ap := &ab
	az := NewBytes().Load(ap).(*Bytes)
	al1, _ := al.SetError("test").(*Bytes).MarshalJSON()
	an1, _ := an.MarshalJSON()
	ay1, _ := ay.MarshalJSON()
	az1, _ := az.MarshalJSON()
	fmt.Printf("%s\n%s\n%s\n%s\n", al1, an1, ay1, az1)
	fmt.Println(ay.String(), ay.SetCoding("string").(*Bytes).String())
	fmt.Println(ay.SetCoding("decimal").(*Bytes).String(), ay.SetCoding("hex").(*Bytes).String())
	fmt.Println(az.Coding(), ay.UnsetError(), ay.SetError("nil!"), ay.Error(), az.UnsetError().(*Bytes).Error())
	fmt.Println(az.String(), ay.Codes(), az.Codes(), az.SetCoding("test"))
	fmt.Println(az.Coding(), ay.Coding(), az.Cap(), ay.Cap(), ay.Len(), az.Len())
	fmt.Println(ay.SetElem(1, byte(13)), az.SetElem(0, byte(23)), az.SetElem(1, 1))
	azt := []byte("this is a test")
	az = az.Load(&azt).(*Bytes)
	fmt.Println(az.Len(), az.Elem(0), ay.Elem(0), az.Elem(-1), az.SetElem(9, byte(110)))
	fmt.Println(ay.Purge(), al.Purge(), az.SetCoding("hex").(*Bytes).String(), az.Purge(), az.String())
	ab = make([]byte, 0)
	ap = &ab
	al.buf = ap
	fmt.Println(NewBytes().Len(), al.Len(), NewBytes().Elem(0), az.Elem(50), al.Elem(1))
	fmt.Println(NewBytes().Size(), az.Size(), ay.Size())
	fmt.Println(NewBytes().Rand(32), az.Rand(23), ay.Rand(13), az.Rand(-3))
	fmt.Println(az.SetCoding("hex").(*Bytes).String(), ay.Null(), az.Null(), NewBytes().Null(), ay.New(32))
	fmt.Println(az.Move(ay), ay.Move(az))
	az.coding = -10
	fmt.Println(az.String(), az.Coding(), az.Error())
	az.coding = 100
	fmt.Println(az.Coding(), az.String(), az.Error())
	var iface interface{}
	fmt.Println(ay.Load(nil), az.Load([]byte(nil)), az.Load(iface), ay.Load(ap))
	fmt.Println(ay.Link(az), az.Link(an), az.Link(nil), ay.Free())
	fmt.Println(ay.Copy(az), az.Copy(an), az.Copy(az), az.Copy(NewBytes().Rand(32)))
}
