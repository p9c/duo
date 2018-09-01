// Package bytes is a wrapper around the native byte slice that automatically handles purging discarded data and enables copy, link and move functions on the data contained inside the structure.
//
// The purpose of this structure is to enable the chaining of pointer methods and eliminate the need for separate assignments by passing error value within the structure instead of as the last term in the return tuple. This structuring has a similar functionality to default parameters, without the compile-time complexity. The same pattern can be used to extend the type to be incorporated into a new aggregate type that can contain more similar structures or add them in addition to implemented methods.
package bytes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	. "gitlab.com/parallelcoin/duo/pkg/byte"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
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

// Nil guards against nil pointer receivers
func ifnil(r *Bytes) *Bytes {
	if r == nil {
		r = new(Bytes)
		r.SetError("nil receiver")
	}
	return r
}

// NewBytes empties an existing bytes or makes a new one
func NewBytes(r ...*Bytes) *Bytes {
	if len(r) == 0 {
		r = append(r, new(Bytes))
	}
	if r[0] == nil {
		r[0] = ifnil(r[0])
		r[0].SetError("receiver was nil")
	}
	if r[0].buf != nil {
		rr := *r[0].buf
		if r[0].set {
			for i := range rr {
				rr[i] = 0
			}
		}
	}
	r[0].buf, r[0].set, r[0].err = nil, false, nil
	return r[0]
}

/////////////////////////////////////////
// Buffer implementations
/////////////////////////////////////////

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() *[]byte {
	if r == nil || r.buf == nil {
		return &[]byte{}
	}
	return r.buf
}

// Copy duplicates the data from the buffer provided and zeroes and replaces its contents, clearing the error value.
func (r *Bytes) Copy(buf Buffer) Buffer {
	r = ifnil(r)
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
	// for i := range *r.Buf() {
	// }
	r.Set()
	return r
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *Bytes) ForEach(f func(int)) Buffer {
	for i := range *r.buf {
		f(i)
	}
	return r
}

// Free zeroes the buffer and dereferences it
func (r *Bytes) Free() Buffer {
	r.Null()
	r.buf = nil
	return r
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(bytes Buffer) Buffer {
	if r == nil {
		r = NewBytes(nil)
	}
	r.Null()
	r.buf, r.set = bytes.Buf(), bytes.IsSet()
	return r
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter.
func (r *Bytes) Load(bytes *[]byte) Buffer {
	if r == nil {
		r = NewBytes()
		r.SetError("nil receiver")
	}
	if bytes == nil {
		r.SetError("nil parameter")
		r.Unset()
		r.Free()
	} else {
		r.Null()
		r.buf = bytes
		r.UnsetError()
		r.Set()
	}
	return r
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(bytes Buffer) Buffer {
	if r == nil {
		r = NewBytes()
	}
	if bytes == nil {
		r.err = errors.New("nil parameter")
	} else {
		r.buf, r.set, r.err = bytes.Buf(), true, nil
		bytes.Load(nil)
		bytes.Unset()
		bytes.SetError("")
	}
	return r
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size.
func (r *Bytes) New(size int) Buffer {
	if r == nil {
		r = NewBytes()
	}
	r.Null()
	rr := make([]byte, size)
	r.buf = &rr
	r.set = true
	return r
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() Buffer {
	return NewBytes(r)
}

// Rand loads a cryptographically random string of []byte of a specified size.
func (r *Bytes) Rand(size int) Buffer {
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
	if r == nil {
		return "nil receiver"
	}
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

/////////////////////////////////////////
// Array implementation
/////////////////////////////////////////

// Elem returns the byte at a given index of the buffer
func (r *Bytes) Elem(i int) Buffer {
	if r == nil {
		r.SetError("nil receiver")
		return &Byte{}
	}
	if r.buf == nil {
		r.SetError("nil buffer")
		return &Byte{}
	}
	R := NewByte().Load(&[]byte{(*r.buf)[i]})
	return R
}

// Len is just a synonym for size
func (r *Bytes) Len() int {
	return r.Size()
}

// Cap returns the amount of elements allocated (can be larger than the size)
func (r *Bytes) Cap() int {
	return cap(*r.buf)
}

// Purge zeroes out all of the buffers in the array
func (r *Bytes) Purge() Array {
	if r == nil {
		r = NewBytes()
		r.SetError("nil receiver")
	}
	if r.buf == nil {
		r.SetError("nil buffer")
		return r
	}
	for i := range *r.buf {
		(*r.buf)[i] = 0
	}
	return r
}

// SetElem sets an element in the buffer
func (r *Bytes) SetElem(i int, b Buffer) Array {
	if r == nil {
		R := NewBytes()
		R.SetError("nil receiver")
		return R
	}
	if r.buf == nil {
		r.SetError("nil value")
		return r
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
func (r *Bytes) Set() Toggle {
	r.set = true
	return r
}

// Unset changes the set flag in a Bytes to false and other functions will treat it as empty
func (r *Bytes) Unset() Toggle {
	r.set = false
	return r
}

/////////////////////////////////////////
// Status implementation
/////////////////////////////////////////

// SetError sets the error string
func (r *Bytes) SetError(s string) Buffer {
	if r == nil {
		r = NewBytes(nil)
		r.SetError("nil receiver")
	}
	r.err = errors.New(s)
	return r
}

// UnsetError sets the error to nil
func (r *Bytes) UnsetError() Buffer {
	if r == nil {
		r = NewBytes()
		r.SetError("nil receiver")
	}
	r.err = nil
	return r
}

/////////////////////////////////////////
// Coding implementation
/////////////////////////////////////////

// Coding returns the coding type to be used by the String function
func (r *Bytes) Coding() string {
	if r == nil {
		r = NewBytes()
		r.SetError("nil receiver")
	}
	if r.coding >= len(*r.buf) {
		r.coding = 0
		r.SetError("invalid coding type in Bytes")
	}
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *Bytes) SetCoding(coding string) Buffer {
	if r == nil {
		r = NewBytes()
		r.SetError("nil receiver")
	}
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
	copy(R, CodeType)
	return
}

/////////////////////////////////////////
// Stringer implementation
/////////////////////////////////////////

// String returns the Bytes in the currently set coding format
func (r *Bytes) String() string {
	switch CodeType[r.coding] {
	case "utf8":
		return string(*r.Buf())
	case "hex":
		return hex.EncodeToString(*r.Buf())
	default:
		r.SetError("coding type not implemented")
	}
	return ""
}

/////////////////////////////////////////
// JSON implementation
/////////////////////////////////////////

// MarshalJSON renders the data as JSON
func (r *Bytes) MarshalJSON() ([]byte, error) {
	var buf string
	if r.buf != nil {
		if r.coding >= len(CodeType) {
			r.SetError("invalid coding type set")
		}
		switch CodeType[r.coding] {
		case "utf8":
			buf = string(*r.buf)
		case "hex":
			buf = string(append([]byte("0x"), []byte(hex.EncodeToString(*r.buf))...))
		default:
			r.SetError("coding type not yet implemented")
		}
	}
	return json.Marshal(&struct {
		Value  string
		IsSet  bool
		Coding string
		Error  string
	}{
		Value:  buf,
		IsSet:  r.set,
		Coding: CodeType[r.coding],
		Error:  r.Error(),
	})
}
