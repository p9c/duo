package buf

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/def"
	"math/big"
)

const (
	me = "gitlab.com/parallelcoin/duo/pkg/buf/bytes/sec.go"
)

// Fenced is a struct that stores and manages memguard.LockedBuffer, ensuring that buffers are destroyed when no longer in use.
type Fenced struct {
	buf    *memguard.LockedBuffer
	coding def.Coding
	status def.Status
}

// NewFenced creates a new secure buffer
func NewFenced() (R *Fenced) {
	R = new(Fenced)
	R.coding = new(def.StringCoding)
	R.status = new(def.ErrorStatus)
	return
}

// Buf returns a pointer to the byte slice in the Fenced.
func (r *Fenced) Buf() interface{} {
	var b []byte
	switch {
	case nil == r:
		r = NewFenced()
		r.Status().Set("receiver")
		b = []byte{}
	case nil == r.buf:
		r.Status().Set("buffer")
		b = []byte{}
	default:
		b = r.buf.Buffer()
		r.Status().Unset()
	}
	return &b
}

// Copy duplicates the buffer from another Fenced.
func (r *Fenced) Copy(b interface{}) def.Buffer {
	switch {
	case nil == r:
		r = NewFenced()
		r.Status().Set("receiver")
	case nil == b:
		r.Status().Set("interface")
		return r
	default:
		switch b.(type) {
		case *Fenced:
			switch {
			case nil == b.(*Fenced):
				r.Status().Set("parameter")
				return r
			case r == b.(*Fenced):
				r.Status().Set("is receiver")
				return r
			case b.(*Fenced).buf == nil:
				r.Status().Set("buffer received")
				return r
			}
			r = r.OfLen(b.(*Fenced).Len()).(*Fenced)
			for i := range b.(*Fenced).buf.Buffer() {
				r.SetElem(i, b.(*Fenced).Array().Elem(i))
			}
		case def.Buffer:
			switch {
			case nil == b.(def.Buffer):

			case r == b.(def.Buffer):
			case b.(def.Buffer).Buf() == nil:
			}
		case *[]byte:
		case []byte:
		case byte:
		default:
			return r.Status().Set(me + "Copy() parameter is wrong type").(*Fenced)
		}
	}
	return r
}

// Free destroys the Fenced and dereferences it
func (r *Fenced) Free() def.Buffer {
	if nil == r {
		r = NewFenced().Status().Set(me + "Buf() nil receiver").(*Fenced)
		return r
	}
	if r.buf != nil {
		r.buf.Destroy()
		r.buf = nil
	}
	r.Status().Unset()
	return r
}

// Len returns the length of the Fenced if it has been loaded, or -1 if not. Note that this method satisfies two interfaces, Buffer and Array
func (r *Fenced) Len() (i int) {
	if nil == r {
		r = NewFenced()
		r.Status().Set(me + "Size() nil receiver")
		return -1
	}
	if r.buf == nil {
		r.Status().Set(me + "Size() nil buffer")
		return 0
	}
	return r.buf.Size()
}

// Link copies the pointer from another Fenced's content, meaning what is written to one will also be visible in the other
func (r *Fenced) Link(buf def.Buffer) def.Buffer {
	if nil == r {
		r = NewFenced().Status().Set(me + "Link() nil receiver").(*Fenced)
	}
	if buf == nil {
		r.Status().Set("interface")
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
		r.Status().Set(me + "Link() cannot link to other type of buffer")
	}
	return r
}

// OfLen allocates a new buffer of a specified size and nulls and frees the existing, if any
func (r *Fenced) OfLen(size int) def.Buffer {
	if nil == r {
		r = NewFenced().Status().Set(me + "Buf() nil receiver").(*Fenced)
	}
	var err error
	r.Null().(*Fenced).buf, err = memguard.NewMutable(size)
	if err != nil {
		r.Status().Set(err.Error())
		return r
	}
	r.Status().Unset()
	return r
}

// Null zeroes out a Fenced buffer
func (r *Fenced) Null() def.Buffer {
	if nil == r {
		return NewFenced().Status().Set(me + "Buf() nil receiver").(*Fenced)
	}
	if r.buf != nil {
		r.Status().Unset().(*Fenced).buf.Wipe()
	} else {
		r.Status().Set(".buf")
	}
	return r
}

// Rand loads the Fenced with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *Fenced) Rand(size ...int) def.Buffer {
	if nil == r {
		r = NewFenced()
		r.Status().Set(me + "Buf() nil receiver")
	}
	if r.buf != nil {
		r.buf.Destroy()
	}
	if len(size) > 0 {
		var err error
		r.buf, err = memguard.NewMutableRandom(size[0])
		if err != nil {
			r.Status().Set(err.Error())
		}
	} else {
		r.Status().Set("length parameter must be > 0")
	}
	return r
}

// Coding returns the coding module of the Fenced buffer
func (r *Fenced) Coding() def.Coding {
	return r.coding
}

// Status returns the status module of the Fenced buffer
func (r *Fenced) Status() def.Status {
	return r.status
}

// Array returns the array module of the Fenced buffer
func (r *Fenced) Array() def.Array {
	return r
}

///////////////////////////////////////////////////////////////////////////////

//  Error returns the string in the err field
func (r *Fenced) Error() (err error) {
	if nil == r {
		r = NewFenced().Status().Set(me + "Error() nil receiver").(*Fenced)
	}
	if r.Status().Error() != nil {
		err = *r.status.Error()
	}
	return
}

///////////////////////////////////////////////////////////////////////////////

// Elem returns the byte at a given index of the buffer
func (r *Fenced) Elem(i int) (I interface{}) {
	if nil == r {
		r = NewFenced().Status().Set(me + "Size() nil receiver").(*Fenced)
	}
	if nil == r.buf {
		r.Status().Set("buffer")
		return byte(0)
	}
	return r.buf.Buffer()[i]
}

// SetElem sets an element in the buffer
func (r *Fenced) SetElem(i int, b interface{}) interface{} {
	switch b.(type) {
	case byte:
		if nil == r {
			r = NewFenced().Status().Set("receiver").(*Fenced)
		}
		if nil == r.buf {
			r.Status().Set("buffer")
			return r
		}
		if i < 0 {
			r.Status().Set("index")
			return r
		}
		if r.Len() > i {
			rr := r.buf.Buffer()
			rr[i] = b.(byte)
		} else {
			r.Status().Set("out of bounds")
		}
	default:
		return r.Status().Set("not a byte")
	}
	return r
}

///////////////////////////////////////////////////////////////////////////////

// MarshalJSON marshals the data of this object into JSON
func (r *Fenced) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = NewFenced().Status().Set(me + "Size() nil receiver").(*Fenced)
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  r.String(),
		Coding: r.Coding().Get(),
		Error:  (*r.Status().Error()).Error(),
	})
}

//	String returns the JSON representing the data in a Fenced
func (r *Fenced) String() (s string) {
	if nil == r {
		r = NewFenced().Status().Set(me + "Size() nil receiver").(*Fenced)
		return "<nil>"
	}
	if nil == r.buf {
		r.Status().Set("buffer")
		return "<nil>"
	}
	if *r.Coding().(*def.StringCoding) > def.StringCoding(len(def.StringCodingTypes)) {
		r.Status().Set("coding")
		r.Coding().Set("decimal")
	}
	switch def.StringCodingTypes[*r.Coding().(*def.StringCoding)] {
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
		r.Coding().Set("decimal")
		return r.String()
	}
}
