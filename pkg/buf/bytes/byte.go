package buf

import (
	"gitlab.com/parallelcoin/duo/pkg/def"
)

var me = "gitlab.com/parallelcoin/duo/pkg/def"

// Byte is a buffer that represennts a byte
type Byte struct {
	buf    byte
	coding *StringCoding
	err    *ErrorStatus
}

func NewByte(b byte) *Byte {
	r := new(Byte)
	r.coding = new(StringCoding)
	r.err = new(ErrorStatus)
	return r
}

var BufferImplementation = true

// Buf returns the byte
func (r *Byte) Buf() interface{} {
	return r.buf
}

// Copy copies a Byte
func (r *Byte) Copy(b def.Buffer) def.Buffer {
	switch b.(type) {
	case *Byte:
		if b.(*Byte) == nil {
			r = NewByte().Status().Set("Copy() nil receiver")
		}
		r.buf = b.(*Byte).buf
		r.coding = b.(*Byte).coding
		r.err = b.(*Byte).err
		r.Status().Unset()
	default:
		r.Status().Set("Copy() buffer not Byte")
	}
	return r
}

func (r *Byte) Free() def.Buffer {
	r.Status().Set("Free() no pointers, cannot deref")
	return r
}

func (r *Byte) Link(interface{}) def.Buffer {
	r.Status().Set("Link() not a pointer type")
	return r
}

func (r *Byte) Load(b interface{}) def.Buffer {
	if r == nil {
		r = NewByte(0)
		r.Status().Set("Load() nil receiver")
	}
	switch b.(type) {
	case *Byte:
		return r.Copy(b)
	case byte:
		r.buf = b
		r.Status().Unset()
		return r
	default:
		r.Status().Set("Copy() incorrect parameter type")
	}
	return r
}

func (r *Byte) Null() def.Buffer {
	r.buf = 0
	r.Status().Unset()
	return r
}

func (r *Byte) Rand(...int) def.Buffer {
	rand.Read(r.buf)
	r.Status().Unset()
	return r
}

func (r *Byte) Size() int {
	return 1
}

func (r *Byte) Coding() Coding {
	return r.coding
}

func (r *Byte) Status() Status {
	return r.err
}

func (r *Byte) Array() Array {
	r.Status().Set("Array() Bytes does not implement array")
	return nil
}
