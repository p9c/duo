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
	set    bool
	coding int
	err    error
}

// NewBytes empties an existing bytes or makes a new one
func NewBytes() (R *Bytes) {
	R = new(Bytes)
	buf := make([]byte, 0)
	R.buf = &buf
	return
}

func nilError() *Bytes {
	r := NewBytes()
	r.err = errors.New("nil receiver")
	return r
}

// Nil implementation

// Nil returns true if the receiver is nil
func (r *Bytes) Nil() bool {
	return r == nil
}

// Buffer implementation

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() (R interface{}) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(!r.set, func() {
				b := make([]byte, 1)
				r.buf = &b
			})
		})
	return r.buf
}

// Copy duplicates the data from the buffer provided and zeroes and replaces its contents, clearing the error value.
func (r *Bytes) Copy(b Buffer) Buffer {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(b,
				func() { r.Free().SetError("nil parameter") },
				func() {
					DoIf(r == b,
						func() { r.SetError("parameter is receiver") },
						func() {
							bb := b.(*Bytes)
							if bb.Size() > 0 {
								bbuf := make([]byte, bb.Size())
								r.buf = &bbuf
								ForEach(bbuf,
									func(i int) bool {
										return r.SetElem(i, bb.Elem(i)) == nil
									})
							}
						})
				})
		})
	return r
}

// Free dereferences the buffer
func (r *Bytes) Free() Buffer {
	DoIf(r == nil,
		func() { r = nilError() },
		func() { r.Null().(*Bytes).buf = nil })
	return r
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(b interface{}) (R Buffer) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(b == nil,
				func() {
					r.SetError("nil parameter")
				},
				func() {
					r.Null().(*Bytes).buf = b.(*Bytes).Buf().(*[]byte)
					r.set = b.(*Bytes).IsSet()
				})
		})
	return r
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter.
func (r *Bytes) Load(bytes interface{}) (R Buffer) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(bytes,
				func() { r.SetError("nil parameter").(Buffer).Free().Unset() },
				func() {
					r.buf = bytes.(*[]byte)
					r.Set()
				})
		})
	return r
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(buf Buffer) (R Buffer) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(buf,
				func() {
					DoIf(buf, func() { r.SetError("nil parameter") },
						func() {
							r.Load(buf.Buf()).UnsetError().(Buffer).Set()
							buf.Null().UnsetError().(Buffer).Unset().(Buffer).Free()
						})
				})
		})
	return r
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size.
func (r *Bytes) New(size int) (R Buffer) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			x := make([]byte, size)
			r = r.Load(&x).Set().(*Bytes)
		})
	return r
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() (R Buffer) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			r.ForEach(func(i int) bool {
				return r.SetElem(i, byte(0)) == nil
			})
		})
	return r
}

// Rand loads a cryptographically random string of []byte of a specified size.
func (r *Bytes) Rand(size ...int) (R Buffer) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(len(size) > 0, func() {
				r = r.New(size[0]).(*Bytes)
				rr := *r.Buf().(*[]byte)
				_, r.err = rand.Read(rr)
			})
		})
	return r
}

// Size returns the length of the *[]byte if it has a value assigned, or -1
func (r *Bytes) Size() (i int) {
	DoIf(r == nil,
		func() { i = -1 },
		func() {
			DoIf(r.IsSet(),
				func() {
					DoIf(r.buf != nil,
						func() { i = len(*r.buf) })
				})
		})
	return
}

// Error implementation

// Error gets the error string
func (r *Bytes) Error() (s string) {
	DoIf(r == nil,
		func() { s = "nil receiver" },
		func() { DoIf(r.err, func() {}, func() { s = r.err.Error() }) })
	return
}

// Array implementation

// Cap returns the amount of elements allocated (can be larger than the size), returns -1 if unallocated
func (r *Bytes) Cap() (i int) {
	DoIf(r,
		func() { i = -1 },
		func() { DoIf(i != -1, func() { i = cap(*r.buf) }) })
	return
}

// Elem returns the byte at a given index of the buffer
func (r *Bytes) Elem(i int) (R interface{}) {
	R = byte(0)
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(r.buf,
				func() { r.SetError("nil buffer") },
				func() {
					DoIf(r.Len() > i,
						func() {
							R = (*r.buf)[i]
						},
						func() {
							r.SetError("index out of bounds")
						})
				})
		})
	return R
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *Bytes) ForEach(f func(int) bool) (b bool) {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(r.buf, func() {
				if f != nil {
					b = ForEach(*r.buf,
						func(i int) bool { return f(i) })
				}
			})
		})
	return
}

// Len is just a synonym for size, returns -1 if unallocated
func (r *Bytes) Len() (i int) {
	return r.Size()
}

// Purge zeroes out all of the buffers in the array
func (r *Bytes) Purge() (R interface{}) {
	DoIf(nil == r,
		func() { r = nilError() },
		func() {
			r.ForEach(func(i int) bool {
				return r.SetElem(i, byte(0)).(*Bytes) == nil
			})
		})
	return r
}

// SetElem sets an element in the buffer
func (r *Bytes) SetElem(i int, b interface{}) (R interface{}) {
	switch b.(type) {
	case byte:
	default:
		return r.SetError("parameter not a byte")
	}
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(b != nil,
				func() {
					DoIf(r.Len() > i,
						func() { (*r.buf)[i] = b.(byte) },
						func() { r.SetError("index out of bounds") })
				})
		})
	return r
}

// Toggle implementation

// IsSet returns true if the Bytes buffer has been loaded with a slice
func (r *Bytes) IsSet() (b bool) {
	DoIf(r, func() { r = nilError() })
	return r.set
}

// Set mark that the value has been initialised/loaded
func (r *Bytes) Set() (R interface{}) {
	DoIf(r,
		func() { r = nilError() },
		func() { r.set = true; R = r })
	return r
}

// Unset changes the set flag in a Bytes to false and other functions will treat it as empty
func (r *Bytes) Unset() (R interface{}) {
	DoIf(r,
		func() { r = nilError() },
		func() { r.set = false; R = r })
	return r
}

// Status implementation

// SetError sets the error string
func (r *Bytes) SetError(s string) interface{} {
	DoIf(r,
		func() { r = nilError() },
		func() {
			DoIf(s != "",
				func() {
					r.err = errors.New(s)
				})
		})
	return r
}

// UnsetError sets the error to nil
func (r *Bytes) UnsetError() interface{} {
	DoIf(r, func() { r = nilError() },
		func() {
			r.err = nil
		})
	return r
}

// Coding implementation

// Coding returns the coding type to be used by the String function
func (r *Bytes) Coding() string {
	DoIf(r, func() { r = nilError() },
		func() {
			DoIf(r.coding >= len(*r.buf), func() {
				r.SetError("invalid coding type in Bytes")
				r.coding = 0
			})
		})
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *Bytes) SetCoding(coding string) interface{} {
	DoIf(r, func() { r = nilError() },
		func() {
			found := false
			for i := range CodeType {
				DoIf(coding == CodeType[i], func() {
					found = true
					r.coding = i
				})
			}
			if !found {
				r.SetError("code type not found")
			}
		})
	return r
}

// Codes returns a copy of the array of CodeType
func (r *Bytes) Codes() (R []string) {
	DoIf(r, func() { r = nilError() },
		func() {
			for i := range CodeType {
				R = append(R, CodeType[i])
			}
		})
	return R
}

// Stringer implementation

// String returns the Bytes in the currently set coding format
func (r *Bytes) String() (S string) {
	DoIf(r, func() { r = nilError() },
		func() {

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
			default:
				S = r.SetCoding("decimal").(Buffer).String()
			}
		})
	return
}

// JSON implementation

// MarshalJSON renders the data as JSON
func (r *Bytes) MarshalJSON() ([]byte, error) {
	DoIf(r, func() { r = nilError() })
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
