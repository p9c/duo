package buf

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"math/big"
)

// Unsafe is a struct that stores and manages byte slices for security purposes, automatically wipes old data when new data is loaded.
//
// The structure stores a boolean signifying whether its value is set to point at a valid slice or not, and an error value, which allows one to use the type in assignments without multiple return values, while still allowing one to check the error value of functions performed with it.
//
// To use it, simply new(Unsafe) to get pointer to a empty new structure, and then after that you can call the methods of the interface.
type Unsafe struct {
	buf    *[]byte
	coding int
	err    error
}

// New makes a new Unsafe
func NewUnsafe() (R *Unsafe) {
	R = new(Unsafe)
	return
}

// Buf returns a variable pointing to the value stored in a Unsafe.
func (r *Unsafe) Buf() (R interface{}) {
	switch {
	case nil == r:
		r = NewUnsafe().SetError("Buf() nil receiver").(*Unsafe)
		return &[]byte{}
	case nil == r.buf:
		r.SetError("Buf() buffer is nil")
		return &[]byte{}
	case len(*r.buf) == 0:
		r.SetError("Buf() buffer is zero length")
		fallthrough
	default:
		return r.buf
	}
}

// Copy duplicates the data from the buffer provided and zeroes and replaces its contents, clearing the error value.
func (r *Unsafe) Copy(b def.Buffer) def.Buffer {
	switch {
	case nil == r:
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
		fallthrough
	case nil == b:
		return r.SetError("Copy() nil parameter").(*Unsafe)
	case r == b.(*Unsafe):
		return r.SetError("Copy() parameter is receiver").(*Unsafe)
	case b.Size() > 0:
		bbuf := make([]byte, b.Size())
		r.buf = &bbuf
		for i := range bbuf {
			r.SetElem(i, b.Elem(i))
		}
	}
	return r
}

// Free dereferences the buffer
func (r *Unsafe) Free() def.Buffer {
	if r == nil {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	} else {
		r.buf = nil
	}
	return r
}

// Link nulls the Unsafe and copies the pointer in from another Unsafe. Whatever is done to one's []byte will also affect the other, but they keep separate error values. I
func (r *Unsafe) Link(b interface{}) (R def.Buffer) {
	switch b.(type) {
	case *Unsafe:
		switch {
		case nil == r:
			r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
			fallthrough
		case b == nil:
			return r.SetError("nil parameter").(*Unsafe)
		default:
			r.Null().(*Unsafe).buf = b.(*Unsafe).buf
		}
	default:
		r.SetError("Link() only accepts *Unsafe")
	}
	return r
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter
func (r *Unsafe) Load(b interface{}) (R def.Buffer) {
	if nil == b {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	}
	switch b.(type) {
	case *[]byte:
		switch {
		case nil == r:
			r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
		case nil != r.buf:
			r.Null()
		}
		r.buf = b.(*[]byte)
		return r.UnsetError().(*Unsafe)
	default:
		r.SetError("Load() only *[]byte accepted")
		return r
	}
}

// Move copies the *[]byte pointer into the Unsafe structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Unsafe) Move(buf def.Buffer) (R def.Buffer) {
	switch buf.(type) {
	case (*Unsafe):
		switch {
		case nil == r:
			r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
			fallthrough
		case nil == buf:
			r.SetError("Move() nil interface")
			return r
		case nil == buf.(*Unsafe):
			r.SetError("Move() nil parameter")
			return r
		default:
			r.buf = buf.(*Unsafe).buf
			buf.Free()
		}
	default:
		r.SetError("type interface not implemented")
	}
	return r
}

// New nulls the Unsafe and assigns an empty *[]byte with a specified size
func (r *Unsafe) New(size int) (R def.Buffer) {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	}
	x := make([]byte, size)
	r = r.Load(&x).(*Unsafe)
	return r
}

// Null wipes the value stored
func (r *Unsafe) Null() (R def.Buffer) {
	if nil == r {
		r = NewUnsafe().SetError("Null() nil receiver").(*Unsafe)
	} else {
		if nil == r.buf {
			r.SetError("Null() nil buffer")
		} else {
			for i := range *r.buf {
				r.SetElem(i, byte(0))
			}
		}
	}
	return r
}

// Rand loads a cryptographically random string of []byte of a specified size.
func (r *Unsafe) Rand(size ...int) (R def.Buffer) {
	if nil == r {
		r = NewUnsafe().SetError("Rand() nil receiver").(*Unsafe)
	}
	if len(size) > 0 {
		if size[0] < 0 {
			r.SetError("Rand() negative size")
			return r
		}
		r = r.New(size[0]).(*Unsafe)
		rr := *r.Buf().(*[]byte)
		_, r.err = rand.Read(rr)
	}
	return r
}

// Size returns the length of the *[]byte if it has a value assigned, or -1
func (r *Unsafe) Size() int {
	if nil == r {
		r = NewUnsafe().SetError("Size() nil receiver").(*Unsafe)
		return -1
	}
	if nil == r.buf {
		r.SetError("Size() nil buffer")
		return -1
	}
	return len(*r.buf)
}

// Error implementation

// Error gets the error string
func (r *Unsafe) Error() string {
	if nil == r {
		r = NewUnsafe().SetError("Error() nil receiver").(*Unsafe)
	}
	if nil == r.err {
		return "<nil>"
	} else {
		return r.err.Error()
	}
}

//Array implementation

// Cap returns the amount of elements allocated (can be larger than the size), returns -1 if unallocated
func (r *Unsafe) Cap() int {
	if nil == r || r.buf == nil {
		return -1
	}
	return cap(*r.buf)
}

// Elem returns the byte at a given index of the buffer
func (r *Unsafe) Elem(i int) (R interface{}) {
	switch {
	case nil == r:
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
		return byte(0)
	case nil == r.buf:
		r.SetError("Elem() nil buffer")
	case r.Len() == 0:
		r.SetError("Elem() array is zero elements")
	case i < 0:
		r.SetError("Elem() index less than zero")
	case r.Len() < i:
		r.SetError("Elem() index out of bounds")
	}
	if r.err != nil {
		return byte(0)
	}
	return (*r.buf)[i]
}

// Len is just a synonym for size
func (r *Unsafe) Len() (i int) {
	return r.Size()
}

// SetElem sets an element in the buffer
func (r *Unsafe) SetElem(i int, b interface{}) (R interface{}) {
	switch b.(type) {
	case byte:
		if nil == r {
			r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
		}
		if b != nil {
			if r.Len() > i {
				(*r.buf)[i] = b.(byte)
			} else {
				r.SetError("SetElem() index out of bounds")
			}
		}
	default:
		return r.SetError("SetElem() parameter not a byte")
	}
	return r
}

// SetError sets the error string
func (r *Unsafe) SetError(s string) interface{} {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	}
	if s != "" {
		r.err = errors.New(s)
	}
	return r
}

// UnsetError sets the error to nil
func (r *Unsafe) UnsetError() interface{} {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	} else {
		r.err = nil
	}
	return r
}

// Coding implementation

// Coding returns the coding type to be used by the String function
func (r *Unsafe) Coding() string {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	}
	if r.coding > len(def.StringCodingTypes) || r.coding < 0 {
		r.SetError("Coding() invalid coding type")
		r.coding = 0
	}
	return def.StringCodingTypes[r.coding]
}

// SetCoding changes the encoding type
func (r *Unsafe) SetCoding(coding string) interface{} {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	}
	found := false
	for i := range def.StringCodingTypes {
		if coding == def.StringCodingTypes[i] {
			found = true
			r.coding = i
			break
		}
	}
	if !found {
		r.SetError("Coding() code type not found")
	}
	return r
}

// Codes returns a copy of the array of def.StringCodingTypes
func (r *Unsafe) Codes() (R []string) {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	}
	for i := range def.StringCodingTypes {
		R = append(R, def.StringCodingTypes[i])
	}
	return
}

// Stringer implementation

// String returns the Unsafe in the currently set coding format
func (r *Unsafe) String() (S string) {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
	}
	if r.coding > len(def.StringCodingTypes) || r.coding < 0 {
		r.SetError("String() invalid coding")
		r.SetCoding("hex")
		return "<invalid coding>"
	}
	if nil == r.buf {
		r.SetError("String() nil buffer")
		return "<nil>"
	}
	if len(*r.buf) == 0 {
		r.SetError("String() zero length buffer")
		return "{}"
	}
	switch def.StringCodingTypes[r.coding] {
	case "byte":
		S = fmt.Sprint(*r.Buf().(*[]byte))
	case "string":
		S = string(*r.Buf().(*[]byte))
	case "decimal":
		bi := big.NewInt(0)
		bi.SetBytes(*r.Buf().(*[]byte))
		S = fmt.Sprint(bi)
	case "hex":
		S = "0x" + hex.EncodeToString(*r.Buf().(*[]byte))
	}
	r.UnsetError()
	return
}

// JSON implementation

// MarshalJSON renders the data as JSON
func (r *Unsafe) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = NewUnsafe().SetError(" nil receiver").(*Unsafe)
		return []byte{}, r.err
	}
	errstring := ""
	if r.err != nil {
		errstring = r.err.Error()
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  r.String(),
		Coding: def.StringCodingTypes[r.coding],
		Error:  errstring,
	})
}
