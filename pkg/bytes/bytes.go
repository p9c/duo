// Package bytes is a wrapper around the native byte slice that automatically handles purging discarded data and enables copy, link and move functions on the data contained inside the structure.
//
// The purpose of this structure is to enable the chaining of pointer methods and eliminate the need for separate assignments by passing error value within the structure instead of as the last term in the return tuple. This structuring has a similar functionality to default parameters, without the compile-time complexity. The same pattern can be used to extend the type to be incorporated into a new aggregate type that can contain more similar structures or add them in addition to implemented methods.
package bytes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
	"math/big"
)

// Bytes is a struct that stores and manages byte slices for security purposes, automatically wipes old data when new data is loaded.
//
// The structure stores a boolean signifying whether its value is set to point at a valid slice or not, and an error value, which allows one to use the type in assignments without multiple return values, while still allowing one to check the error value of functions performed with it.
//
// To use it, simply new(Bytes) to get pointer to a empty new structure, and then after that you can call the methods of the interface.
type Bytes struct {
	buf    *[]byte
	coding int
	err    error
}

// NewBytes empties an existing bytes or makes a new one
func NewBytes() (R *Bytes) {
	R = new(Bytes)
	return
}

// Nil implementation

func nilError(s string) *Bytes {
	r := NewBytes()
	r.SetError(s + " nil receiver")
	return r
}

// Buffer implementation

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() (R interface{}) {
	switch {
	case nil == r:
		r = nilError("Buf()")
		fallthrough
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
func (r *Bytes) Copy(b Buffer) Buffer {
	switch {
	case nil == r:
		r = nilError("Copy()")
		fallthrough
	case nil == b:
		return r.SetError("Copy() nil parameter").(*Bytes)
	case r == b.(*Bytes):
		return r.SetError("Copy() parameter is receiver").(*Bytes)
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
func (r *Bytes) Free() Buffer {
	if r == nil {
		r = nilError("Free()")
	}
	r.Null().(*Bytes).buf = nil
	return r
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(b interface{}) (R Buffer) {
	switch b.(type) {
	case *Bytes:
		switch {
		case nil == r:
			r = nilError("Link()")
			fallthrough
		case b == nil:
			return r.SetError("nil parameter").(*Bytes)
		default:
			r.Null().(*Bytes).buf = b.(*Bytes).buf
		}
	default:
		r.SetError("Link() only accepts *Bytes")
	}
	return r
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter
func (r *Bytes) Load(bytes interface{}) (R Buffer) {
	if nil == bytes {
		r = nilError("Load()")
	}
	switch bytes.(type) {
	case *[]byte:
		switch {
		case nil == r:
			r = nilError("Load()")
		}
		r.buf = bytes.(*[]byte)
		return r.UnsetError().(*Bytes)
	default:
		r.SetError("Load() only *[]byte accepted")
		return r
	}
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(buf Buffer) (R Buffer) {
	switch {
	case nil == r:
		r = nilError("Move()")
		fallthrough
	case nil == buf:
		r.SetError("nil parameter")
		return r
	default:
		r.Load(buf.Buf()).UnsetError()
		buf.Null().UnsetError().(*Bytes).Free()
	}
	return r
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size
func (r *Bytes) New(size int) (R Buffer) {
	if nil == r {
		r = nilError("New()")
	}
	x := make([]byte, size)
	r = r.Load(&x).(*Bytes)
	return r
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() (R Buffer) {
	if nil == r {
		r = nilError("Null()")
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
func (r *Bytes) Rand(size ...int) (R Buffer) {
	if nil == r {
		r = nilError("Rand()")
	}
	if len(size) > 0 {
		if size[0] < 0 {
			r.SetError("Rand() negative size")
			return r
		}
		r = r.New(size[0]).(*Bytes)
		rr := *r.Buf().(*[]byte)
		_, r.err = rand.Read(rr)
	}
	return r
}

// Size returns the length of the *[]byte if it has a value assigned, or -1
func (r *Bytes) Size() int {
	if nil == r {
		return -1
	}
	if r.buf != nil {
		return len(*r.buf)
	}
	return 0
}

// Error implementation

// Error gets the error string
func (r *Bytes) Error() string {
	if nil == r {
		return "Error() nil receiver"
	}
	if nil == r.err {
		return "<nil>"
	}
	return r.err.Error()
}

//Array implementation

// Cap returns the amount of elements allocated (can be larger than the size), returns -1 if unallocated
func (r *Bytes) Cap() int {
	if nil == r || r.buf == nil {
		return -1
	}
	return cap(*r.buf)
}

// Elem returns the byte at a given index of the buffer
func (r *Bytes) Elem(i int) (R interface{}) {
	switch {
	case nil == r:
		r = nilError("Elem()")
		return byte(0)
	case nil == r.buf:
		r.SetError("Elem() nil buffer")
		return byte(0)
	}
	if r.Len() == 0 {
		r.SetError("Elem() array is zero elements")
		return byte(0)
	}
	if i < 0 {
		r.SetError("Elem() index less than zero")
		return byte(0)
	}
	if r.Len() < i {
		r.SetError("Elem() index out of bounds")
		return byte(0)
	}
	return (*r.buf)[i]
}

// Len is just a synonym for size, returns -1 if unallocated
func (r *Bytes) Len() (i int) {
	if nil == r || r.buf == nil {
		return -1
	}
	return len(*r.buf)
}

// Purge zeroes out all of the buffers in the array
func (r *Bytes) Purge() interface{} {
	switch {
	case nil == r:
		r = nilError("Purge()")
	case nil == r.buf:
		r.SetError("Purge() nil buffer")
	default:
		for i := range *r.buf {
			r.SetElem(i, byte(0))
		}
	}
	return r
}

// SetElem sets an element in the buffer
func (r *Bytes) SetElem(i int, b interface{}) (R interface{}) {
	switch b.(type) {
	case byte:
		if nil == r {
			r = nilError("SetElem()")
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

// Status implementation

// SetError sets the error string
func (r *Bytes) SetError(s string) interface{} {
	if nil == r {
		r = nilError("SetError()")
	}
	if s != "" {
		r.err = errors.New(s)
	}
	return r
}

// UnsetError sets the error to nil
func (r *Bytes) UnsetError() interface{} {
	if nil == r {
		r = nilError("UnsetError()")
	} else {
		r.err = nil
	}
	return r
}

// Coding implementation

// Coding returns the coding type to be used by the String function
func (r *Bytes) Coding() string {
	if nil == r {
		r = nilError("Coding()")
	}
	if r.coding > len(CodeType) || r.coding < 0 {
		r.SetError("Coding() invalid coding type")
		r.coding = 0
	}
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *Bytes) SetCoding(coding string) interface{} {
	if nil == r {
		r = nilError("SetCoding()")
	}
	found := false
	for i := range CodeType {
		if coding == CodeType[i] {
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

// Codes returns a copy of the array of CodeType
func (r *Bytes) Codes() (R []string) {
	if nil == r {
		r = nilError("Codes()")
	}
	for i := range CodeType {
		R = append(R, CodeType[i])
	}
	return
}

// Stringer implementation

// String returns the Bytes in the currently set coding format
func (r *Bytes) String() (S string) {
	if nil == r {
		r = nilError("String()")
	}
	if r.coding > len(CodeType) || r.coding < 0 {
		r.SetError("String() invalid coding")
		r.SetCoding("hex")
		return "<invalid coding>"
	}
	switch CodeType[r.coding] {
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
	return
}

// JSON implementation

// MarshalJSON renders the data as JSON
func (r *Bytes) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = nilError("MarshalJSON()")
		return []byte{}, r.err
	}
	var errstring string
	if r.err != nil {
		errstring = r.err.Error()
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  r.String(),
		Coding: CodeType[r.coding],
		Error:  errstring,
	})
}
