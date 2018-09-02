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
	. "gitlab.com/parallelcoin/duo/pkg/byte"
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
	set    bool
	coding int
	err    error
}

// NewBytes empties an existing bytes or makes a new one
func NewBytes() (R *Bytes) {
	R = new(Bytes)
	R.buf = nil
	return
}

func (r *Bytes) ifnil() (R *Bytes) {
	if r == nil {
		return NewBytes().SetError("nil receiver").(*Bytes)
	}
	return r
}

/////////////////////////////////////////
// Buffer implementations
/////////////////////////////////////////

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() (R *[]byte) {
	r = r.ifnil()
	if r == nil || !r.IsSet() {
		return &[]byte{0}
	}
	return r.buf
}

// Copy duplicates the data from the buffer provided and zeroes and replaces its contents, clearing the error value.
func (r *Bytes) Copy(buf Buffer) (R Buffer) {
	r = r.ifnil()
	r.UnsetError()
	if buf == nil {
		return r.Free().SetError("nil parameter")
	}
	if r == buf {
		return r.SetError("parameter is receiver")
	}
	if buf.Size() == 0 {
		return r.Null().Load(buf.Buf()).SetError("empty buffer received")
	}
	r.New(buf.Size())
	r.ForEach(func(i int) {
		r.SetElem(i, buf.Elem(i))
	})
	r.Set()
	return r
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *Bytes) ForEach(f func(int)) (R Buffer) {
	r = r.ifnil()
	if r.buf != nil {
		for i := range *r.buf {
			f(i)
		}
	}
	return r
}

// Free dereferences the buffer
func (r *Bytes) Free() (R Buffer) {
	r = r.ifnil()
	if r != nil {
		r.Null()
		r.buf = nil
	}
	return r
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(bytes Buffer) (R Buffer) {
	r = r.ifnil()
	r.Null()
	r.buf, r.set = bytes.Buf(), bytes.IsSet()
	return r
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter.
func (r *Bytes) Load(bytes *[]byte) (R Buffer) {
	r = r.ifnil()
	if bytes == nil {
		r.SetError("nil parameter").Free().Unset()
	} else {
		r.buf = bytes
	}
	return r
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(bytes Buffer) (R Buffer) {
	r = r.ifnil()
	if bytes == nil {
		r.SetError("nil parameter")
	} else {
		r.Load(bytes.Buf())
		r.UnsetError().Set()
		bytes.Load(nil).UnsetError().Unset()
	}
	return r
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size.
func (r *Bytes) New(size int) (R Buffer) {
	r = r.ifnil()
	x := make([]byte, size)
	r.Load(&x)
	r.Set()
	return r
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() (R Buffer) {
	r = r.ifnil()
	R = r.ForEach(func(i int) {
		r.SetElem(i, NewByte().Load(&[]byte{0}))
	})
	return
}

// Rand loads a cryptographically random string of []byte of a specified size.
func (r *Bytes) Rand(size int) (R Buffer) {
	r = r.ifnil()
	r = r.New(size).(*Bytes)
	rand.Read(*r.Buf())
	r.set = true
	return r
}

// Size returns the length of the *[]byte if it has a value assigned, or -1
func (r *Bytes) Size() int {
	if r == nil {
		return -1
	}
	if r.IsSet() {
		if r.buf != nil {
			return len(*r.buf)
		}
	}
	return 0
}

/////////////////////////////////////////
// Error implementation
/////////////////////////////////////////

// Error gets the error string
func (r *Bytes) Error() string {
	r = r.ifnil()
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

/////////////////////////////////////////
// Array implementation
/////////////////////////////////////////

// Elem returns the byte at a given index of the buffer
func (r *Bytes) Elem(i int) (R Buffer) {
	r = r.ifnil()
	if r.buf == nil {
		r.SetError("nil buffer")
		return &Byte{}
	}
	R = NewByte().Load(r.buf)
	return
}

// Len is just a synonym for size, returns -1 if unallocated
func (r *Bytes) Len() int {
	if r == nil {
		return -1
	}
	return r.Size()
}

// Cap returns the amount of elements allocated (can be larger than the size), returns -1 if unallocated
func (r *Bytes) Cap() int {
	if r == nil {
		return -1
	}
	return cap(*r.buf)
}

// Purge zeroes out all of the buffers in the array
func (r *Bytes) Purge() (R Array) {
	r = r.ifnil()
	r.ForEach(func(i int) {
		(*r.buf)[i] = 0
	})
	return r
}

// SetElem sets an element in the buffer
func (r *Bytes) SetElem(i int, b Buffer) (R Array) {
	r = r.ifnil()
	if b.Len() < 1 {
		return r
	}
	if len(*r.buf) < i {
		return r.SetError("index out of bounds")
	}
	(*r.buf)[i] = (*b.Buf())[0]
	return r
}

/////////////////////////////////////////
// Toggle implementation
/////////////////////////////////////////

// IsSet returns true if the Bytes buffer has been loaded with a slice
func (r *Bytes) IsSet() bool {
	if r == nil {
		return false
	}
	return r.set
}

// Set mark that the value has been initialised/loaded
func (r *Bytes) Set() (R Toggle) {
	r = r.ifnil()
	r.set = true
	return r
}

// Unset changes the set flag in a Bytes to false and other functions will treat it as empty
func (r *Bytes) Unset() Toggle {
	r = r.ifnil()
	r.set = false
	return r
}

/////////////////////////////////////////
// Status implementation
/////////////////////////////////////////

// SetError sets the error string
func (r *Bytes) SetError(s string) (R Buffer) {
	r = r.ifnil()
	r.err = errors.New(s)
	return r
}

// UnsetError sets the error to nil
func (r *Bytes) UnsetError() (R Buffer) {
	r = r.ifnil()
	r.err = nil
	return r
}

/////////////////////////////////////////
// Coding implementation
/////////////////////////////////////////

// Coding returns the coding type to be used by the String function
func (r *Bytes) Coding() string {
	r = r.ifnil()
	if r.coding >= len(*r.buf) {
		r.coding = 0
		r.SetError("invalid coding type in Bytes")
	}
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *Bytes) SetCoding(coding string) (R Buffer) {
	r = r.ifnil()
	found := false
	for i := range CodeType {
		if coding == CodeType[i] {
			found = true
			r.coding = i
		}
	}
	if !found {
		r.SetError("code type not found")
	}
	return r
}

// Codes returns a copy of the array of CodeType
func (r *Bytes) Codes() (R []string) {
	r = r.ifnil()
	for i := range CodeType {
		R = append(R, CodeType[i])
	}
	return
}

/////////////////////////////////////////
// Stringer implementation
/////////////////////////////////////////

// String returns the Bytes in the currently set coding format
func (r *Bytes) String() (S string) {
	r = r.ifnil()
	switch CodeType[r.coding] {
	case "byte":
		return fmt.Sprint(*r.Buf())
	case "string":
		return string(*r.Buf())
	case "decimal":
		bi := big.NewInt(0)
		bi.SetBytes(*r.Buf())
		return fmt.Sprint(bi)
	case "hex":
		return "0x" + hex.EncodeToString(*r.Buf())
	default:
		return r.SetCoding("decimal").String()
	}
}

/////////////////////////////////////////
// JSON implementation
/////////////////////////////////////////

// MarshalJSON renders the data as JSON
func (r *Bytes) MarshalJSON() ([]byte, error) {
	r = r.ifnil()
	return json.Marshal(&struct {
		Value  string
		IsSet  bool
		Coding string
		Error  string
	}{
		Value:  r.String(),
		IsSet:  r.set,
		Coding: CodeType[r.coding],
		Error:  r.Error(),
	})
}
