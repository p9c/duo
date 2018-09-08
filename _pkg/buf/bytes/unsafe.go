package buf

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"math/big"
)

// Unsafe is a struct that stores and manages byte slices for security purposes, automatically wipes old data when new data is loaded.
type Unsafe struct {
	buf    *[]byte
	coding *def.StringCoding
	err    *def.ErrorStatus
}

// NewUnsafe makes a new Unsafe
func NewUnsafe() *Unsafe {
	r := new(Unsafe)
	r.coding = new(def.StringCoding)
	r.err = new(def.ErrorStatus)
	return r
}

// Buf returns a variable pointing to the value stored in a Unsafe.
func (r *Unsafe) Buf() (R interface{}) {
	switch {
	case nil == r:
		r = NewUnsafe()
		r.Status().Set("nil receiver")
		return &[]byte{}
	case nil == r.buf:
		r.Status().Set("buffer is nil")
		return &[]byte{}
	case len(*r.buf) == 0:
		r.Status().Set("buffer is zero length")
		fallthrough
	default:
		return r.buf
	}
}

// Copy duplicates the data from the buffer provided and zeroes and replaces its contents, clearing the error value.
func (r *Unsafe) Copy(b interface{}) def.Buffer {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set("nil interface")
	}
	if b == nil {
		r.Status().Set("nil interface parameter")
		return r
	}
	switch b.(type) {
	case *Unsafe:
		B := b.(def.Buffer)
		switch {
		case r == b:
			r.Status().Set("parameter is receiver")
			return r
		case B.Len() == 0:
			r.Status().Set("zero buffer")
			r = r.OfLen(0).(*Unsafe)
		case B.Len() > 0:
			B := b.(def.Array)
			bbuf := make([]byte, B.Len())
			r.buf = &bbuf
			for i := range bbuf {
				r.SetElem(i, B.Elem(i))
			}
		}
	case []byte:
		B := b.([]byte)
		bb := make([]byte, len(B))
		copy(bb, B)
		r.buf = &B
	case byte:
		r.Copy(&[]byte{b.(byte)})
	}
	return r
}

// Free dereferences the buffer
func (r *Unsafe) Free() def.Buffer {
	if r == nil {
		r = NewUnsafe()
		r.Status().Set(" nil receiver")
	} else {
		r.buf = nil
	}
	return r
}

// Len returns the length of the *[]byte if it has a value assigned, or -1
func (r *Unsafe) Len() int {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set("nil receiver")
		return -1
	}
	if nil == r.buf {
		r.Status().Set("nil buffer")
		return -1
	}
	return len(*r.buf)
}

// Link copies the link from a buffer after purging its contents, if any
func (r *Unsafe) Link(b def.Buffer) (R def.Buffer) {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set("nil receiver")
	}
	if b == nil {
		r.Status().Set("nil parameter")
		return r
	}
	switch b.(type) {
	case *Unsafe:
		switch {
		case r.buf == nil:
			r.Status().Set("nil buffer")
		default:
			r = r.Null().(*Unsafe)
			r.buf = b.(*Unsafe).buf
		}
	default:
		r.Status().Set("only accepts *Unsafe")
	}
	return r
}

// OfLen nulls the Unsafe and assigns an empty *[]byte with a specified size
func (r *Unsafe) OfLen(size int) (R def.Buffer) {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set("nil receiver")
	}
	if nil != r.buf {
		r.Null()
	}
	b := make([]byte, size)
	r.buf = &b
	return r
}

// Null wipes the value stored
func (r *Unsafe) Null() (R def.Buffer) {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set("nil receiver")
	} else {
		if nil == r.buf {
			r.Status().Set("nil buffer")
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
		r = NewUnsafe()
		r.Status().Set("nil receiver")
	}
	if len(size) > 0 {
		if size[0] < 0 {
			r.Status().Set("negative size")
			return r
		}
		r = r.OfLen(size[0]).(*Unsafe)
		rr := *r.Buf().(*[]byte)
		n, err := rand.Read(rr)
		if n != size[0] {
			r.Status().Set("Rand() only got", n, "of", size[0], "bytes")
		}
		if err != nil {
			r.Status().Set(err.Error())
		}
	} else {
		r.Status().Set("size parameter required")
	}
	return r
}

// Coding returns the string coding type handler
func (r *Unsafe) Coding() def.Coding {
	if r == nil {
		r = NewUnsafe()
	}
	return r.coding
}

// Status returns the status object
func (r *Unsafe) Status() def.Status {
	if r == nil {
		r = NewUnsafe()
	}
	if r.err == nil {
		r.err = new(def.ErrorStatus)
	}
	return r.err
}

// Array returns the array module of the Fenced buffer
func (r *Unsafe) Array() def.Array {
	if r == nil {
		r = NewUnsafe()
		r.Status().Set("nil receiver")
	}
	return r
}

// Elem returns the byte at a given index of the buffer
func (r *Unsafe) Elem(i int) (R interface{}) {
	switch {
	case nil == r:
		r = NewUnsafe()
		r.Status().Set("nil receiver")
		return byte(0)
	case nil == r.buf:
		r.Status().Set("nil buffer")
		return byte(0)
	case r.Len() == 0:
		r.Status().Set("array is zero elements")
		return byte(0)
	case i < 0:
		r.Status().Set("index less than zero")
		return byte(0)
	case r.Len() < i:
		r.Status().Set("index out of bounds")
		return byte(0)
	}
	return (*r.buf)[i]
}

// SetElem sets an element in the buffer
func (r *Unsafe) SetElem(i int, b interface{}) (R interface{}) {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set(" nil receiver")
	}
	if nil == r.buf {
		r.Status().Set("nil buffer")
		return byte(0)
	}
	switch b.(type) {
	case byte:
		if b != nil {
			switch {
			case r.Len() == 0:
				r.Status().Set("array is zero elements")
				return byte(0)
			case i < 0:
				r.Status().Set("negative index")
			case r.Len() < i:
				r.Status().Set("index out of bounds")
				return byte(0)
			default:
				(*r.buf)[i] = b.(byte)
			}
		}
	default:
		r.Status().Set("parameter not a byte")
		return r
	}
	return r
}

// String returns the Unsafe in the currently set coding format
func (r *Unsafe) String() (S string) {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set(" nil receiver")
	}
	if nil == r.buf {
		r.Status().Set("nil buffer")
		return "<nil>"
	}
	if len(*r.buf) == 0 {
		r.Status().Set("zero length buffer")
		return "{}"
	}
	switch def.StringCodingTypes[*r.Coding().(*def.StringCoding)] {
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
	r.Status().Unset()
	return
}

// MarshalJSON renders the data as JSON
func (r *Unsafe) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = NewUnsafe()
		r.Status().Set("nil receiver")
	}
	errstring := r.Status().Error()
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  r.String(),
		Coding: def.StringCodingTypes[*r.Coding().(*def.StringCoding)],
		Error:  errstring,
	})
}
