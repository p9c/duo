package buf

import (
	"crypto/rand"
	"gitlab.com/parallelcoin/duo/pkg/def"
)

// Byte is a buffer that represents a byte
type Byte struct {
	buf    byte
	coding *def.StringCoding
	err    *def.ErrorStatus
}

// NewByte makes a new Byte
func NewByte(b byte) *Byte {
	r := new(Byte)
	r.buf = b
	r.coding = new(def.StringCoding)
	r.err = new(def.ErrorStatus)
	return r
}

// Buf returns the byte
func (r *Byte) Buf() interface{} {
	if r == nil {
		return byte(0)
	}
	return r.buf
}

// Copy copies a Byte or byte
func (r *Byte) Copy(b interface{}) def.Buffer {
	if r == nil {
		r = NewByte(0)
		r.Status().Set("nil interface")
	}
	if b == nil {
		r.Status().Set("parameter")
		return r
	}
	switch b.(type) {
	case *Byte:
		B := b.(*Byte)
		if B == nil {
			r = NewByte(0)
			r.Status().Set("nil parameter")
		}
		if B == nil {
			r.Status().Set("nil Byte")
		}
		r.buf = B.Buf().(byte)
		// r.coding = b.(*Byte).Coding().(*def.StringCoding)
		// r.err = b(*Byte).err
		r.Status().Unset()
	case byte:
		r.buf = b.(byte)
	default:
		r.Status().Set("buffer not Byte")
	}
	return r
}

// Free is no-op as the buffer part of the struct
func (r *Byte) Free() def.Buffer {
	r.Status().Set("no pointers, cannot deref")
	return r
}

// Len always returns 1 because this is a byte
func (r *Byte) Len() int {
	return 1
}

// Link does nothing because bytes are stored in the struct
func (r *Byte) Link(b def.Buffer) def.Buffer {
	r.Status().Set("not a pointer type")
	return r
}

// OfLen does nothing because this type only provides 1 byte
func (r *Byte) OfLen(int) def.Buffer {
	return NewByte(0)
}

// Null zeroes out the buffer
func (r *Byte) Null() def.Buffer {
	r.buf = 0
	r.Status().Unset()
	return r
}

// Rand gets a random byte and loads it into the buffer
func (r *Byte) Rand(...int) def.Buffer {
	b := make([]byte, 1)
	rand.Read(b)
	r.buf = b[0]
	r.Status().Unset()
	return r
}

// Coding returns the string coding type handler
func (r *Byte) Coding() def.Coding {
	if r == nil {
		r = NewByte(0)
	}
	return r.coding
}

// Status returns the status object
func (r *Byte) Status() def.Status {
	if r == nil {
		r = NewByte(0)
	}
	if r.err == nil {
		r.err = new(def.ErrorStatus)
	}
	return r.err
}

// Array is not implemented for Byte
func (r *Byte) Array() def.Array {
	if r == nil {
		r = NewByte(0)
	}
	r.Status().Set("Bytes does not implement array")
	return nil
}

func (r *Byte) String() (S string) {
	b := r.Buf().(byte)
	R := NewUnsafe().Copy(b)
	R.Coding().Set(r.Coding().Get())
	S = R.String()
	return
}
