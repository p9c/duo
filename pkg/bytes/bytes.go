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

func donil(r *Bytes, f1 func(), f2 func()) *Bytes {
	if r == nil {
		f1()
	}
	if f2 != nil {
		f2()
	}
	return r
}

func doif(b bool, fn func()) {
	if b {
		fn()
	}
}

// NewBytes empties an existing bytes or makes a new one
func NewBytes() (R *Bytes) {
	return new(Bytes).Load(&[]byte{}).(*Bytes)
}

/////////////////////////////////////////
// Buffer implementations
/////////////////////////////////////////

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() (R *[]byte) {
	donil(r, func() {
		r = NewBytes().SetError("nil receiver").(*Bytes)
	}, nil)
	doif(r.buf == nil, func() { r.buf = &[]byte{} })
	return r.buf
}

// Copy duplicates the data from the buffer provided and zeroes and replaces its contents, clearing the error value.
func (r *Bytes) Copy(buf Buffer) (R Buffer) {
	return donil(r, func() {
		r = NewBytes().SetError("nil receiver").(*Bytes)
	},
		func() {
			if buf != nil {
				doif(r == buf, func() { R = r.SetError("parameter is receiver") })
				doif(buf == nil, func() { R = r.Free().SetError("nil parameter") })
			}
		})
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *Bytes) ForEach(f func(int)) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			doif(r.buf != nil, func() {
				for i := range *r.buf {
					f(i)
				}
			})
		})
}

// Free dereferences the buffer
func (r *Bytes) Free() (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() { doif(r != nil, func() { r.buf = nil }) })
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(bytes Buffer) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			r.Null()
			r.buf, r.set = bytes.Buf(), bytes.IsSet()
		})
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter.
func (r *Bytes) Load(bytes *[]byte) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			doif(bytes == nil, func() { r.SetError("nil parameter").Free().Unset() })
			r.buf = bytes
		})
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(bytes Buffer) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			doif(bytes == nil, func() { r.SetError("nil parameter") })
			doif(bytes != nil, func() {
				r.Load(bytes.Buf()).UnsetError().Set()
				bytes.Load(nil).UnsetError().Unset()
			})
		})
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size.
func (r *Bytes) New(size int) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			x := make([]byte, size)
			r = r.Load(&x).Set().(*Bytes)
		})
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			R = r.ForEach(func(i int) { r.SetElem(i, NewByte().Load(&[]byte{0})) })
		})
}

// Rand loads a cryptographically random string of []byte of a specified size.
func (r *Bytes) Rand(size int) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			r = r.New(size).Set().(*Bytes)
			rand.Read(*r.Buf())
		})
}

// Size returns the length of the *[]byte if it has a value assigned, or -1
func (r *Bytes) Size() (i int) {
	doif(r == nil, func() { i = -1 })
	doif(r.IsSet(), func() { doif(r.buf != nil, func() { i = len(*r.buf) }) })
	return
}

/////////////////////////////////////////
// Error implementation
/////////////////////////////////////////

// Error gets the error string
func (r *Bytes) Error() (s string) {
	donil(r, func() {
		r = NewBytes()
		r.err = errors.New("nil receiver")
	}, func() {
		doif(s != "", func() { r.err = errors.New(s) })
	})
	return
}

/////////////////////////////////////////
// Array implementation
/////////////////////////////////////////

// Elem returns the byte at a given index of the buffer
func (r *Bytes) Elem(i int) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			if r.buf == nil {
				R = r.SetError("nil buffer").Load(NewByte().Buf())
			} else {
				R = NewByte().Load(r.buf)
			}
			R = r
		})
}

// Len is just a synonym for size, returns -1 if unallocated
func (r *Bytes) Len() (i int) {
	return r.Size()
}

// Cap returns the amount of elements allocated (can be larger than the size), returns -1 if unallocated
func (r *Bytes) Cap() (i int) {
	doif(r == nil, func() { i = -1 })
	doif(i != -1, func() { i = cap(*r.buf) })
	return
}

// Purge zeroes out all of the buffers in the array
func (r *Bytes) Purge() (R Array) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() { r.ForEach(func(i int) { (*r.buf)[i] = 0 }) })
}

// SetElem sets an element in the buffer
func (r *Bytes) SetElem(i int, b Buffer) (R Array) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			doif(len(*r.buf) < i, func() {
				doif(b.Len() < 1, func() { r.SetError("index out of bounds") })
				doif(r.buf != nil && b.IsSet(), func() { (*r.buf)[i] = (*b.Buf())[0] })
			})
		})
}

// IsSet returns true if the Bytes buffer has been loaded with a slice
func (r *Bytes) IsSet() (b bool) {
	donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) }, nil)
	return r.set
}

// Set mark that the value has been initialised/loaded
func (r *Bytes) Set() (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() { r.set = true; R = r })
}

// Unset changes the set flag in a Bytes to false and other functions will treat it as empty
func (r *Bytes) Unset() (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() { r.set = false; R = r })
}

/////////////////////////////////////////
// Status implementation
/////////////////////////////////////////

// SetError sets the error string
func (r *Bytes) SetError(s string) (R Buffer) {
	return donil(r, func() {
		r = NewBytes()
		r.err = errors.New("nil receiver")
	},
		func() {
			doif(s != "", func() { r.err = errors.New(s) })
		})
}

// UnsetError sets the error to nil
func (r *Bytes) UnsetError() (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			r.err = nil
		})
}

/////////////////////////////////////////
// Coding implementation
/////////////////////////////////////////

// Coding returns the coding type to be used by the String function
func (r *Bytes) Coding() string {
	donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			doif(r.coding >= len(*r.buf), func() {
				r.SetError("invalid coding type in Bytes")
				r.coding = 0
			})
		})
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *Bytes) SetCoding(coding string) (R Buffer) {
	return donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			found := false
			for i := range CodeType {
				doif(coding == CodeType[i], func() {
					found = true
					r.coding = i
				})
			}
			if !found {
				r.SetError("code type not found")
			}
		})
}

// Codes returns a copy of the array of CodeType
func (r *Bytes) Codes() (R []string) {
	donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			for i := range CodeType {
				R = append(R, CodeType[i])
			}
		})
	return R
}

/////////////////////////////////////////
// Stringer implementation
/////////////////////////////////////////

// String returns the Bytes in the currently set coding format
func (r *Bytes) String() (S string) {
	donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) },
		func() {
			switch CodeType[r.coding] {
			case "byte":
				S = fmt.Sprint(*r.Buf())
			case "string":
				S = string(*r.Buf())
			case "decimal":
				bi := big.NewInt(0)
				bi.SetBytes(*r.Buf())
				S = fmt.Sprint(bi)
			case "hex":
				S = "0x" + hex.EncodeToString(*r.Buf())
			default:
				S = r.SetCoding("decimal").String()
			}
		})
	return
}

/////////////////////////////////////////
// JSON implementation
/////////////////////////////////////////

// MarshalJSON renders the data as JSON
func (r *Bytes) MarshalJSON() ([]byte, error) {
	donil(r, func() { r = NewBytes().SetError("nil receiver").(*Bytes) }, nil)
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
