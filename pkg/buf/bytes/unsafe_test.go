package buf

import (
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"testing"
)

func TestUnsafe(t *testing.T) {
	al := &Unsafe{}
	an := NewUnsafe()
	var ay *Unsafe
	ab := make([]byte, 0)
	ap := &ab
	az := NewUnsafe().Load(ap).(*Unsafe)
	al1, _ := al.SetError("test").(*Unsafe).MarshalJSON()
	an1, _ := an.MarshalJSON()
	ay1, _ := ay.MarshalJSON()
	az1, _ := az.MarshalJSON()
	fmt.Printf("%s\n%s\n%s\n%s\n", al1, an1, ay1, az1)
	fmt.Println(ay.String(), ay.SetCoding("string").(*Unsafe).String())
	fmt.Println(ay.SetCoding("decimal").(*Unsafe).String(), ay.SetCoding("hex").(*Unsafe).String())
	fmt.Println(az.Coding(), ay.UnsetError(), ay.SetError("nil!"), ay.Error(), az.UnsetError().(*Unsafe).Error())
	fmt.Println(az.String(), ay.Codes(), az.Codes(), az.SetCoding("test"))
	fmt.Println(az.Coding(), ay.Coding(), az.Cap(), ay.Cap(), ay.Len(), az.Len())
	fmt.Println(ay.SetElem(1, byte(13)), az.SetElem(0, byte(23)), az.SetElem(1, 1))
	azt := []byte("this is a test")
	az = az.Load(&azt).(*Unsafe)
	fmt.Println(az.Len(), az.Elem(0), ay.Elem(0), az.Elem(-1), az.SetElem(9, byte(110)))
	fmt.Println(az.SetCoding("hex").(*Unsafe).String(), az.String())
	ab = make([]byte, 0)
	ap = &ab
	al.buf = ap
	fmt.Println(NewUnsafe().Len(), al.Len(), NewUnsafe().Elem(0), az.Elem(50), al.Elem(1))
	fmt.Println(NewUnsafe().Size(), az.Size(), ay.Size())
	fmt.Println(NewUnsafe().Rand(32), az.Rand(23), ay.Rand(13), az.Rand(-3))
	fmt.Println(az.SetCoding("hex").(*Unsafe).String(), ay.Null(), az.Null(), NewUnsafe().Null(), ay.New(32))
	str1 := "testing 123"
	ptr1 := []byte(str1)
	fmt.Println(az.Move(ay), ay.Move(az), az.Move(NewUnsafe().Load(&ptr1)))
	az.coding = -10
	fmt.Println(az.String(), az.Coding(), az.Error())
	az.coding = 100
	fmt.Println(az.Coding(), az.String(), az.Error())
	var iface interface{}
	fmt.Println(ay.Load(nil), az.Load([]byte(nil)), az.Load(iface), ay.Load(ap))
	fmt.Println(ay.Link(az), az.Link(an), az.Link(nil), ay.Free())
	fmt.Println(ay.Copy(az), az.Copy(an), az.Copy(az), az.Copy(NewUnsafe().Rand(32)))
	br := new(def.Buffer)
	fmt.Println(az.Move(*br))
	fmt.Println(az.SetCoding("string").(*Unsafe).String(), az.SetCoding("decimal").(*Unsafe).String())
	fmt.Println(NewUnsafe().SetCoding("hex").(*Unsafe).New(0).(*Unsafe).String())
	var vv *Unsafe
	fmt.Println(NewUnsafe().Buf(), vv.Buf(), NewUnsafe().New(0).Buf())
}
