package buf

import (
	"encoding/json"
	"fmt"
	"gitlab.com/parallelcoin/duo/lib/array"
	"gitlab.com/parallelcoin/duo/lib/coding"
	"gitlab.com/parallelcoin/duo/lib/status"
	"strings"
)

// Byte is a simple buffer that just stores one byte
type Byte struct {
	Buf    byte
	Status string
	Coding string
}

// NewByte makes a new Byte
func NewByte() *Byte {
	b := new(Byte)
	b.Coding = "string"
	return b
}

// Freeze writes the current data structure into a string that format as content of a JSON struct node
func (r *Byte) Freeze() (S string) {
	if r == nil {
		r = NewByte()
		r.SetStatus("nil receiver")
	}
	s := []string{
		`{"Buf":`,
		`` + fmt.Sprint(r.Buf) + `,`,
		`"Status":`,
		`"` + r.Status + `",`,
		`"Coding":`,
		`"` + r.Coding + `"}`,
	}
	S = strings.Join(s, "")
	return
}

// Thaw is
func (r *Byte) Thaw(s string) interface{} {
	if r == nil {
		r = NewByte()
		r.SetStatus("nil receiver")
	}
	out := NewByte()
	err := json.Unmarshal([]byte(s), out)
	if err != nil {
		out.SetStatus(err.Error())
	}
	return out
}

// Data returns the content of the buffer
func (r *Byte) Data(out ...interface{}) interface{} {
	r = r.UnsetStatus().(*Byte)
	if len(out) > 0 {
		out[0] = r.Buf
	}
	return r.Buf
}

// Copy accepts a parameter and copies the (first) byte in it into its buffer
func (r *Byte) Copy(b interface{}) Buf {
	r = r.UnsetStatus().(*Byte)
	switch b.(type) {
	case int:
		r.Buf = byte(b.(int))
	case uint:
		r.Buf = byte(b.(uint))
	case byte:
		r.Buf = b.(byte)
	case int8:
		r.Buf = byte(b.(int8))
	case uint16:
		r.Buf = byte(b.(uint16))
	case int16:
		r.Buf = byte(b.(int16))
	case uint32:
		r.Buf = byte(b.(uint32))
	case int32:
		r.Buf = byte(b.(int32))
	case uint64:
		r.Buf = byte(b.(uint64))
	case int64:
		r.Buf = byte(b.(int64))
	case []byte:
		if len(b.([]byte)) > 0 {
			r.Buf = b.([]byte)[0]
		}
	case *[]byte:
		if len(*b.(*[]byte)) > 0 {
			r.Buf = (*b.(*[]byte))[0]
		}
	case *Byte:
		B := b.(*Byte)
		if r == B {
			return r.SetStatus("copy to self").(*Byte)
		}
		r.Buf, r.Status, r.Coding = B.Buf, B.Status, B.Coding
		return r
	default:
		r.SetStatus("parameter type not implemented")
		return r
	}
	return r
}

// Free doesn't really do anything but other buffers will need it
func (r *Byte) Free() Buf {
	r = r.UnsetStatus().(*Byte)
	return r
}

// Null is
func (r *Byte) Null() Buf {
	if r == nil {
		r = r.UnsetStatus().(*Byte)
	}
	r.Buf = 0
	r.UnsetStatus()
	r.SetCoding("")
	return r
}

// SetStatus sets the error state
func (r *Byte) SetStatus(s string) status.Status {
	r = r.UnsetStatus().(*Byte)
	r.Status = s
	return r
}

// SetStatusIf sets an error from a standard error interface variable if it is set
func (r *Byte) SetStatusIf(err error) status.Status {
	if err != nil {
		r.Status = err.Error()
	}
	return r
}

// UnsetStatus emptys the error state
func (r *Byte) UnsetStatus() status.Status {
	if r == nil {
		r = NewByte()
		r.SetStatus("nil receiver")
	} else {
		r.Status = ""
	}
	return r
}

// GetCoding is
func (r *Byte) GetCoding() string {
	r = r.UnsetStatus().(*Byte)
	if r.Coding == "" {
		r.Coding = "string"
	}
	r.SetCoding(r.Coding)
	return r.Coding
}

// SetCoding is
func (r *Byte) SetCoding(s string) coding.Coding {
	r = r.UnsetStatus().(*Byte)
	found := false
	for i := range coding.Codings {
		if s == coding.Codings[i] {
			found = true
		}
	}
	if found {
		r.Coding = s
	} else {
		r.Coding = "string"
	}
	return r
}

// ListCodings is
func (r *Byte) ListCodings() []string {
	r = r.UnsetStatus().(*Byte)
	return coding.Codings
}

// Elem returns a byte of 1 or 0 representing a bit
func (r *Byte) Elem(e int) interface{} {
	r = r.UnsetStatus().(*Byte)
	if e > 7 {
		r.SetStatus("only 8 bits in byte")
	}
	// compl := byte(7 - e)
	return r.Buf >> byte(e) & 1
}

// Len is 1, though we can also read the bits
func (r *Byte) Len() int {
	r = r.UnsetStatus().(*Byte)
	return 1
}

// SetElem is
func (r *Byte) SetElem(e int, val interface{}) arr.Array {
	if val == 0 {
		mask := ^(byte(1) << byte(e))
		r.Buf &= mask
	} else {
		r.Buf |= (byte(1) << byte(e))
	}
	return r
}

// String is
func (r *Byte) String() (S string) {
	c := r.GetCoding()
	return coding.Encode([]byte{r.Buf}, c)
}

// Error is
func (r *Byte) Error() string {
	return r.Status
}
