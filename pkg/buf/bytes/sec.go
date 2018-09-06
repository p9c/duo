package buf

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"math/big"
)

const (
	me = "gitlab.com/parallelcoin/duo/pkg/buf/bytes/sec.go"
)

// Fenced is a struct that stores and manages memguard.LockedBuffer, ensuring that buffers are destroyed when no longer in use.
//
// Do not use struct literals and not assign them to a name and null() (deletes and zeroes struct) afterwards, or you could run out of memguard Secs
//
// All functions except for those exporting buffers will automatically allocate the struct of the receiver if it is nil. This permits the use of struct literals in assignment for one-liners that initialise values and call a function with data. It may introduce side effects in code if you did not intend to create a new variable.
//
// The maximum size of buffer is around 172500 bytes on a linux 4.18, it may be more may be less.
type Fenced struct {
	buf    *memguard.LockedBuffer
	coding int
	err    error
}

// NewFenced creates a new secure buffer
func NewFenced() (R *Fenced) {
	R = new(Fenced)
	return
}

// BufferImplementation - implements def.Buffer
var BufferImplementation = true

// Buf returns a pointer to the byte slice in the Fenced.
func (r *Fenced) Buf() interface{} {
	var b []byte
	switch {
	case nil == r:
		r = NewFenced().SetError(me + "Buf() nil receiver").(*Fenced)
		return []byte{}
	case nil == r.buf:
		r.SetError(me + "Buf() nil buffer")
		b = []byte{}
	default:
		b = r.UnsetError().(*Fenced).buf.Buffer()
	}
	return &b
}

// Copy duplicates the buffer from another Fenced.
func (r *Fenced) Copy(b def.Buffer) def.Buffer {
	switch {
	case nil == r:
		r = NewFenced().SetError(me + "Copy() nil receiver").(*Fenced)
	case nil == b:
		return r.SetError("Copy() nil interface").(*Fenced)
	default:
		switch b.(type) {
		case *Fenced, :
			switch {
			case nil == b.(*Fenced):
				return r.SetError("Copy() nil parameter").(*Fenced)
			case r == b.(*Fenced):
				return r.SetError("Copy() parameter is receiver").(*Fenced)
			case b.(*Fenced).buf == nil:
				return r.SetError("Copy() nil buffer received").(*Fenced)
			}
			r.New(b.Size())
			for i := range b.(*Fenced).buf.Buffer() {
				r.SetElem(i, b.Elem(i))
			}
		case def.Buffer:

		default:
			return r.SetError(me + "Copy() parameter is wrong type").(*Fenced)
		}
	}
	return r
}

// Free destroys the Fenced and dereferences it
func (r *Fenced) Free() def.Buffer {
	if nil == r {
		r = NewFenced().SetError(me + "Buf() nil receiver").(*Fenced)
		return r
	}
	r.buf = nil
	r.UnsetError()
	return r
}

// Link copies the pointer from another Fenced's content, meaning what is written to one will also be visible in the other
func (r *Fenced) Link(buf interface{}) def.Buffer {
	if nil == r {
		r = NewFenced().SetError(me + "Link()() nil receiver").(*Fenced)
	}
	if buf == nil {
		r.SetError("me + Link() nil interface")
		return r
	}
	switch buf.(type) {
	case *Fenced:
		if buf.(*Fenced) != nil {
			if nil != r.buf {
				r.buf.Destroy()
			}
		}
		r.buf = buf.(*Fenced).buf
	default:
		r.SetError(me + "Link() cannot link to other type of buffer")
	}
	return r
}

// Load moves the contents of a byte slice into the Fenced, erasing the original copy.
func (r *Fenced) Load(b interface{}) def.Buffer {
	if nil == r {
		r = NewFenced().SetError(me + "Load() nil receiver").(*Fenced)
	}
	if nil == b {
		r.SetError(me + "Load() nil parameter")
		return r
	}
	if r.buf, r.err = memguard.NewMutableFromBytes(*b.(*[]byte)); r.err == nil {
		r.UnsetError()
	}
	return r
}


// New destroys the old Lockedbuffer and assigns a new one with a given length
func (r *Fenced) New(size int) def.Buffer {
	if nil == r {
		r = NewFenced().SetError(me + "Buf() nil receiver").(*Fenced)
	}
	r.Null().(*Fenced).buf, r.err = memguard.NewMutable(size)
	if r.err != nil {
		return r
	}
	r.UnsetError()
	return r
}

// Null zeroes out a Fenced
func (r *Fenced) Null() def.Buffer {
	if nil == r {
		return NewFenced().SetError(me + "Buf() nil receiver").(*Fenced)
	}
	if r.buf != nil {
		r.UnsetError().(*Fenced).buf.Wipe()
	} else {
		r.SetError("Null() nil .buf")
	}
	return r
}

// Rand loads the Fenced with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *Fenced) Rand(size ...int) def.Buffer {
	if nil == r {
		r = NewFenced().SetError(me + "Buf() nil receiver").(*Fenced)
	}
	if len(size) > 0 {
		r.Null()
	}
	r.buf, r.err = memguard.NewMutableRandom(size[0])
	return r
}

// Size returns the length of the Fenced if it has been loaded, or -1 if not
func (r *Fenced) Size() (i int) {
	if nil == r {
		r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
		return -1
	}
	if r.buf == nil {
		r.SetError(me + "Size() nil buffer")
		return 0
	}
	return r.buf.Size()
}

var coding_implementation = true

// Coding returns the coding type to be used by the String function
func (r *Fenced) Coding() string {
	if nil == r {
		r = NewFenced().SetError(me + "Buf() nil receiver").(*Fenced)
	}
	if r.coding >= len(def.StringCodingTypes) {
		r.coding = 0
		r.SetError("Coding() invalid coding type")
	}
	return def.StringCodingTypes[r.coding]
}

// SetCoding changes the encoding type
func (r *Fenced) SetCoding(coding string) interface{} {
	if nil == r {
		r = NewFenced().SetError(me + "SetCoding() nil receiver").(*Fenced)
	}
	found := false
	for i := range def.StringCodingTypes {
		if coding == def.StringCodingTypes[i] {
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

//Codes returns a copy of the array of def.StringCodingTypes
func (r *Fenced) Codes() (R []string) {
	for i := range def.StringCodingTypes {
		R = append(R, def.StringCodingTypes[i])
	}
	return
}

var error_implementation = true

//  Error returns the string in the err field
func (r *Fenced) Error() (s string) {
	if nil == r {
		r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
	}
	if r.err != nil {
		s = r.err.Error()
	}
	return s
}

// SetError sets the string of the error in the err field
func (r *Fenced) SetError(s string) interface{} {
	if nil == r {
		r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
	}
	r.err = errors.New(s)
	fmt.Println("SetError() [", s, "]")
	return r
}

// UnsetError sets the error to nil
func (r *Fenced) UnsetError() interface{} {
	if nil == r {
		r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
	} else {
		r.err = nil
	}
	return r
}

var array_implementation = true

// Elem returns the byte at a given index of the buffer
func (r *Fenced) Elem(i int) (I interface{}) {
	if nil == r {
		r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
	}
	if nil == r.buf {
		r.SetError("Elem() nil buffer")
		return byte(0)
	}
	return r.buf.Buffer()[i]
}

// Len returns the length of the array
func (r *Fenced) Len() int {
	return r.Size()
}

// Cap returns the amount of elements allocated (can be larger than the size)
func (r *Fenced) Cap() (i int) {
	if nil == r || r.buf == nil {
		i = 0
	}
	i = cap(*(r.Buf().(*[]byte)))
	return i
}

// SetElem sets an element in the buffer
func (r *Fenced) SetElem(i int, b interface{}) interface{} {
	switch b.(type) {
	case byte:
		if nil == r {
			r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
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

// MarshalJSON marshals the data of this object into JSON
func (r *Fenced) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  r.String(),
		Coding: def.StringCodingTypes[r.coding],
		Error:  r.Error(),
	})
}

var stringer_implementation = true

//	String returns the JSON representing the data in a Fenced
func (r *Fenced) String() (s string) {
	if nil == r {
		r = NewFenced().SetError(me + "Size() nil receiver").(*Fenced)
		return "<nil>"
	}
	if nil == r.buf {
		r.SetError("String() nil buffer")
		return "<nil>"
	}
	if r.coding > len(def.StringCodingTypes) {
		r.SetError("invalid coding")
		r.SetCoding("decimal")
	}
	switch def.StringCodingTypes[r.coding] {
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
