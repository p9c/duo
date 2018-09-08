package buf

import (
	"encoding/json"
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/def"
)

// Fenced is a struct that stores and manages memguard.LockedBuffer, ensuring that buffers are destroyed when no longer in use.
type Fenced struct {
	buf    *memguard.LockedBuffer
	coding def.StringCoding
	status def.ErrorStatus
}

// NewFenced creates a new secure buffer
func NewFenced() (R *Fenced) {
	R = new(Fenced)
	return
}

func (r *Fenced) Print() {
	if r == nil {
		fmt.Println("{}")
	} else {
		fmt.Print("{")
		if r.buf == nil {
			fmt.Print("buf: nil, ")
		} else {
			fmt.Print("buf:", r.buf, ",")
		}
		fmt.Print("coding: ", r.coding, ",")
		fmt.Print("status: ", r.status.Err, ",")
		fmt.Print("},")
	}
}

// Buf returns a pointer to the byte slice in the Fenced.
func (r *Fenced) Buf() interface{} {
	var b []byte
	switch {
	case nil == r:
		r = NewFenced()
		r.Status().Set("nil receiver")
		b = []byte{}
	case nil == r.buf:
		r.Status().Set("nil buffer")
		b = []byte{}
	default:
		b = r.buf.Buffer()
		r.Status().Unset()
	}
	return &b
}

// Copy duplicates the buffer from another Fenced.
func (r *Fenced) Copy(b interface{}) def.Buffer {
	// r.Print()
	switch {
	case nil == r:
		r = NewFenced()
		r.Status().Set("nil receiver")
		fallthrough
	default:
		switch b.(type) {
		case nil:
			r.Status().Set("nil interface")
			r.Print()
			return r
		case *Fenced:
			B := b.(*Fenced)
			switch {
			case r == B:
				r.Status().Set("parameter is receiver")
				r.Print()
				B.Print()
				return r
			case B.buf == nil:
				r.Status().Set("nil buffer received")
				r.Print()
				B.Print()
				return r
			}
			B.Print()
			r = r.OfLen(B.Len()).(*Fenced)
			for i := range B.buf.Buffer() {
				r.SetElem(i, B.Array().Elem(i))
			}
		case def.Buffer:
			switch {
			case r == b.(def.Buffer):
				r.Print()
			case b.(def.Buffer).Buf() == nil:
				r.Print()
			}
		case *[]byte:
			B := b.(*[]byte)
			if B == nil {
				r.Print()
				r.Status().Set("nil *[]byte")
				return r
			}
			var err error
			r.buf, err = memguard.NewMutable(len(*B))
			if err != nil {
				r.Status().Set(err.Error())
				return r
			}
			for i := range *B {
				r.SetElem(i, (*B)[i])
				(*B)[i] = 0
			}
		case []byte:
			r.Print()
			B := b.([]byte)
			if B == nil {
				r.Status().Set("nil []byte")
				return r
			}
			var err error
			r.buf, err = memguard.NewMutable(len(B))
			if err != nil {
				r.Status().Set(err.Error())
				return r
			}
			for i := range B {
				r.SetElem(i, B[i])
				B[i] = 0
			}
		case byte:
			r.Print()
			var err error
			r.buf, err = memguard.NewMutable(1)
			if err != nil {
				r.Status().Set(err.Error())
				return r
			}
			r.SetElem(0, b.(byte))
		default:
			r.Print()
			r.Status().Set("parameter is wrong type")
		}
	}
	return r
}

// Free destroys the Fenced and dereferences it
func (r *Fenced) Free() def.Buffer {
	if nil == r {
		r = NewFenced()
		r.Status().Set("nil receiver")
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
		r.Status().Set("nil receiver")
		return -1
	}
	if r.buf == nil {
		r.Status().Set("nil buffer")
		return 0
	}
	return r.buf.Size()
}

// Link copies the pointer from another Fenced's content, meaning what is written to one will also be visible in the other
func (r *Fenced) Link(buf def.Buffer) def.Buffer {
	if nil == r {
		r = NewFenced()
		r.Status().Set("nil receiver")
	}
	if buf == nil {
		r.Status().Set("nil parameter")
		return r
	}
	switch buf.(type) {
	case *Fenced:
		B := buf.(*Fenced)
		if B == nil {
			r.Status().Set("nil parameter")
			return r
		}
		if B != nil {
			if nil != r.buf {
				r.buf.Destroy()
			}
		}
		r.buf = B.buf
	default:
		r.Status().Set("cannot link to other type of buffer")
	}
	return r
}

// OfLen allocates a new buffer of a specified size and nulls and frees the existing, if any
func (r *Fenced) OfLen(size int) def.Buffer {
	switch {
	case nil == r:
		r = NewFenced()
		r.Status().Set("nil receiver")
	case size < 0:
		r.Status().Set("negative size")
		return r
	case size < 1:
		r.Status().Set("zero length not allowed")
	}
	var err error
	r.Null().(*Fenced).buf, err = memguard.NewMutable(size)
	if err != nil {
		r.Status().Set(err.Error())
	}
	r.Status().Unset()
	return r
}

// Null zeroes out a Fenced buffer
func (r *Fenced) Null() def.Buffer {
	if nil == r {
		r = NewFenced()
		r.Status().Set("nil receiver")
		return r
	}
	if r.buf != nil {
		r.Status().Unset()
		r.buf.Wipe()
	} else {
		r.Status().Set("nil buffer")
	}
	return r
}

// Rand loads the Fenced with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *Fenced) Rand(size ...int) def.Buffer {
	switch {
	case nil == r:
		r = NewFenced()
		r.Status().Set("nil receiver")
	case r.buf != nil:
		r.buf.Destroy()
	case len(size) > 0:
		var err error
		r.buf, err = memguard.NewMutableRandom(size[0])
		if err != nil {
			r.Status().Set(err.Error())
		} else {
			r.Status().Unset()
		}
	default:
		r.Status().Set("length parameter must be > 0")
	}
	return r
}

// Coding returns the coding module of the Fenced buffer
func (r *Fenced) Coding() def.Coding {
	if nil == r {
		r = NewFenced()
		r.Status().Set("nil receiver")
	}
	return r.coding
}

// Status returns the status module of the Fenced buffer
func (r *Fenced) Status() def.Status {
	if nil == r {
		r = NewFenced()
		r.Status().Set("nil receiver")
	}
	return &r.status
}

// Array returns the array module of the Fenced buffer
func (r *Fenced) Array() def.Array {
	if nil == r {
		r = NewFenced()
		r.Status().Set("nil receiver")
	}
	return r
}

// Elem returns the byte at a given index of the buffer
func (r *Fenced) Elem(i int) interface{} {
	var b byte
	switch {
	case nil == r:
		r = NewFenced()
		r.Status().Set("nil receiver")
		b = 0
	case nil == r.buf:
		r.Status().Set("nil buffer")
		b = 0
	case i >= r.Len():
		r.Status().Set("out of bounds")
		b = 0
	default:
		fmt.Println(r.Len())
		b = r.buf.Buffer()[i]
	}
	return b
}

// SetElem sets an element in the buffer
func (r *Fenced) SetElem(i int, b interface{}) interface{} {
	switch {
	case i < 0:
		r.Status().Set("negative index")
		return r
	case nil == r:
		r = NewFenced()
		r.Status().Set("nil receiver")
	case nil == r.buf:
		r.Status().Set("nil buffer")
		return r
	}
	switch b.(type) {
	case byte:
		switch {
		case r.Len() > i:
			rr := r.buf.Buffer()
			rr[i] = b.(byte)
		default:
			r.Status().Set("out of bounds")
		}
	default:
		return r.Status().Set("not a byte")
	}
	return r
}

// MarshalJSON marshals the data of this object into JSON
func (r *Fenced) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = NewFenced()
		r.Status().Set("nil receiver")
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  r.String(),
		Coding: r.Coding().Get(),
		Error:  r.Status().Error(),
	})
}

//	String returns the JSON representing the data in a Fenced
func (r *Fenced) String() (S string) {
	b := r.Buf().(*[]uint8)
	R := NewUnsafe().Copy(b)
	R.Coding().Set(r.Coding().Get())
	S = R.String()
	R.Null()
	return
}
