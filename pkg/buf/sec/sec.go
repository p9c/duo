// Package secbuf is a wrapper around the memguard SecBuf that automatically handles destroying data no longer needed and enables copy, link and move functions on the data contained inside the structure.
package secbuf

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"math/big"
)

// SecBuf is a struct that stores and manages memguard.LockedBuffer, ensuring that buffers are destroyed when no longer in use.
//
// Do not use struct literals and not assign them to a name and null() (deletes and zeroes struct) afterwards, or you could run out of memguard Secs
//
// All functions except for those exporting buffers will automatically allocate the struct of the receiver if it is nil. This permits the use of struct literals in assignment for one-liners that initialise values and call a function with data. It may introduce side effects in code if you did not intend to create a new variable.
//
// The maximum size of buffer is around 172500 bytes on a linux 4.18, it may be more may be less.
type SecBuf struct {
	buf    *memguard.LockedBuffer
	coding int
	err    error
}

// New clears the passed SecBuf or creates a new one if null
func New() (R *SecBuf) {
	R = new(SecBuf)
	return
}

func nilError(s string) *SecBuf {
	r := New()
	r.err = errors.New(s + " nil receiver")
	return r
}

// def.Buffer implementation

// Buf returns a pointer to the byte slice in the SecBuf.
func (r *SecBuf) Buf() interface{} {
	var b []byte
	if nil == r {
		r = nilError("Buf()")
		return []byte{}
	}
	if nil == r.buf {
		b = []byte{}
		r.SetError("Buf() nil buffer")
	} else {
		b = r.UnsetError().(*SecBuf).buf.Buffer()
	}
	return &b
}

// Copy duplicates the buffer from another SecBuf.
func (r *SecBuf) Copy(b def.Buffer) def.Buffer {
	if nil == r {
		r = nilError("Copy()")
	}
	if nil == b {
		r.SetError("Copy() nil interface")
		return r
	}
	if nil == b.(*SecBuf) {
		r.SetError("Copy() nil parameter")
		return r
	} else {
		if r == b.(*SecBuf) {
			r.SetError("Copy() parameter is receiver")
			return r
		}
		if b.(*SecBuf).buf == nil {
			r.SetError("Copy() nil buffer received")
			return r
		}
		r.New(b.Size())
		for i := range b.(*SecBuf).buf.Buffer() {
			r.SetElem(i, b.Elem(i))
		}
	}
	return r
}

// Free destroys the SecBuf and dereferences it
func (r *SecBuf) Free() def.Buffer {
	if nil == r {
		r = nilError("Free()")
		return r
	}
	r.buf = nil
	r.UnsetError()
	return r
}

// Link copies the pointer from another SecBuf's content, meaning what is written to one will also be visible in the other
func (r *SecBuf) Link(buf interface{}) def.Buffer {
	if nil == r {
		r = nilError("Link()")
	}
	if buf == nil {
		r.SetError("Link() nil interface")
		return r
	}
	switch buf.(type) {
	case *SecBuf:
		if buf.(*SecBuf) != nil {
			if nil != r.buf {
				r.buf.Destroy()
			}
		}
		r.buf = buf.(*SecBuf).buf
	default:
		r.SetError("Link() cannot link to other type of buffer")
	}
	return r
}

// Load moves the contents of a byte slice into the SecBuf, erasing the original copy.
func (r *SecBuf) Load(bytes interface{}) def.Buffer {
	if nil == r {
		r = nilError("Load()")
	}
	if nil == bytes {
		r.SetError("nil parameter")
		return r
	}
	if r.buf, r.err = memguard.NewMutableFromBytes(*bytes.(*[]byte)); r.err == nil {
		r.UnsetError()
	}
	return r
}

// Move copies the pointer to the buffer into the receiver and nulls the passed SecBuf
func (r *SecBuf) Move(buf def.Buffer) def.Buffer {
	if nil == r {
		r = nilError("Move()")
	}
	if nil == buf {
		r.SetError("Move() nil parameter")
		return r
	}
	r.buf = buf.(*SecBuf).buf
	r.UnsetError()
	buf.UnsetError()
	buf.(*SecBuf).buf = nil
	return r
}

// New destroys the old Lockedbuffer and assigns a new one with a given length
func (r *SecBuf) New(size int) def.Buffer {
	if nil == r {
		r = nilError("New()")
	}
	r.Null().(*SecBuf).buf, r.err = memguard.NewMutable(size)
	if r.err != nil {
		return r
	}
	r.UnsetError()
	return r
}

// Null zeroes out a SecBuf
func (r *SecBuf) Null() def.Buffer {
	if nil == r {
		r = nilError("Null()")
	}
	if r.buf != nil {
		r.UnsetError().(*SecBuf).buf.Wipe()
	} else {
		r.SetError("Null() nil .buf")
	}
	return r
}

// Rand loads the SecBuf with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *SecBuf) Rand(size ...int) def.Buffer {
	if nil == r {
		r = nilError("Rand()")
	}
	if len(size) > 0 {
		r.Null()
	}
	r.buf, r.err = memguard.NewMutableRandom(size[0])
	return r
}

// Size returns the length of the SecBuf if it has been loaded, or -1 if not
func (r *SecBuf) Size() (i int) {
	if nil == r {
		return -1
	}
	if r.buf == nil {
		return 0
	}
	return r.buf.Size()
}

// Coding implementation

// Coding returns the coding type to be used by the String function
func (r *SecBuf) Coding() string {
	if nil == r {
		r = nilError("Coding()")
	}
	if r.coding >= len(def.CodingTypes) {
		r.coding = 0
		r.SetError("Coding() invalid coding type")
	}
	return def.CodingTypes[r.coding]
}

// SetCoding changes the encoding type
func (r *SecBuf) SetCoding(coding string) interface{} {
	if nil == r {
		r = nilError("SetCoding()")
	}
	found := false
	for i := range def.CodingTypes {
		if coding == def.CodingTypes[i] {
			r.coding = i
			found = true
			break
		}
	}
	if !found {
		r.SetError("SetCoding() code type not found")
	}
	r.UnsetError()
	return r
}

//Codes returns a copy of the array of def.CodingTypes
func (r *SecBuf) Codes() (R []string) {
	for i := range def.CodingTypes {
		R = append(R, def.CodingTypes[i])
	}
	return
}

// Status implementation

//  Error returns the string in the err field
func (r *SecBuf) Error() (s string) {
	if nil == r {
		r = nilError("Error()")
	}
	if r.err != nil {
		s = r.err.Error()
	}
	return s
}

// SetError sets the string of the error in the err field
func (r *SecBuf) SetError(s string) interface{} {
	if nil == r {
		r = nilError("SetError()")
	}
	r.err = errors.New(s)
	fmt.Println("SetError() [", s, "]")
	return r
}

// UnsetError sets the error to nil
func (r *SecBuf) UnsetError() interface{} {
	if nil == r {
		r = nilError("UnsetError()")
	} else {
		r.err = nil
	}
	return r
}

// Array implementation

// Elem returns the byte at a given index of the buffer
func (r *SecBuf) Elem(i int) (I interface{}) {
	if nil == r {
		r = nilError("Elem()")
	}
	if nil == r.buf {
		r.SetError("Elem() nil buffer")
		return byte(0)
	}
	return r.buf.Buffer()[i]
}

// Len returns the length of the array
func (r *SecBuf) Len() int {
	return r.Size()
}

// Cap returns the amount of elements allocated (can be larger than the size)
func (r *SecBuf) Cap() (i int) {
	if nil == r || r.buf == nil {
		i = 0
	}
	i = cap(*(r.Buf().(*[]byte)))
	return i
}

// SetElem sets an element in the buffer
func (r *SecBuf) SetElem(i int, b interface{}) interface{} {
	switch b.(type) {
	case byte:
		if nil == r {
			r = nilError("SetElem()")
		}
		if nil == r.buf {
			r.SetError("SetElem() nil buffer")
			return r
		}
		if i < 0 {
			r.SetError("SetElem() negative index")
			return r
		}
		if r.Len() > i {
			rr := r.buf.Buffer()
			rr[i] = b.(byte)
		} else {
			r.SetError("index out of bounds")
		}
	default:
		return r.SetError("parameter not a byte")
	}
	return r
}

// Purge zeroes out all of the buffers in the array
func (r *SecBuf) Purge() interface{} {
	if r == nil {
		r = nilError("Purge()")
	} else {
		if nil == r.buf {
			r.SetError("Purge() nil buffer")
		} else {
			r.buf.Wipe()
		}
	}
	return r
}

// JSON

// MarshalJSON marshals the data of this object into JSON
func (r *SecBuf) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = nilError("MarshalJSON()")
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  r.String(),
		Coding: def.CodingTypes[r.coding],
		Error:  r.Error(),
	})
}

// Stringer implementation

//	String returns the JSON representing the data in a SecBuf
func (r *SecBuf) String() (s string) {
	if nil == r {
		r = nilError("String()")
		return "<nil>"
	}
	if nil == r.buf {
		r.SetError("String() nil buffer")
		return "<nil>"
	}
	if r.coding > len(def.CodingTypes) {
		r.SetError("invalid coding")
		r.SetCoding("decimal")
	}
	switch def.CodingTypes[r.coding] {
	case "byte":
		return fmt.Sprint(*(r.Buf().(*[]byte)))
	case "string":
		return string(*r.Buf().(*[]byte))
	case "decimal":
		bi := big.NewInt(0)
		bi.SetBytes(*r.Buf().(*[]byte))
		return fmt.Sprint(bi)
	case "hex":
		return "0x" + hex.EncodeToString(*r.Buf().(*[]byte))
	default:
		return r.SetCoding("decimal").(def.Buffer).String()
	}
}
