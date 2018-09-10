package buf

import (
	"encoding/json"
	"fmt"
	"gitlab.com/parallelcoin/duo/lib/array"
	"gitlab.com/parallelcoin/duo/lib/coding"
	"gitlab.com/parallelcoin/duo/lib/status"
	"math/big"
	"strings"
)

// Bytes is a simple buffer for slices of bytes and can ingest almost any other kind of data if one were so inclined.
type Bytes struct {
	Buf    *[]byte
	Status string
	Coding string
}

// NewBytes makes a new Bytes
func NewBytes() *Bytes {
	b := new(Bytes)
	b.Coding = "string"
	return b
}

// Freeze writes the current data structure into a string that format as content of a JSON struct node
func (r *Bytes) Freeze() (S string) {
	if r == nil {
		r = NewBytes()
		r.SetStatus("nil receiver")
	}
	s := []string{
		`{"Buf":`,
		`"` + r.SetCoding("base64").(*Bytes).String() + `",`,
		`"Status":`,
		`"` + r.Status + `",`,
		`"Coding":`,
		`"` + r.Coding + `"}`,
	}
	S = strings.Join(s, "")
	return
}

// Thaw is
func (r *Bytes) Thaw(s string) interface{} {
	if r == nil {
		r = NewBytes()
		r.SetStatus("nil receiver")
	}
	out := NewBytes()
	err := json.Unmarshal([]byte(s), out)
	if err != nil {
		out.SetStatus(err.Error())
	}
	return out
}

// Data returns the content of the buffer
func (r *Bytes) Data() interface{} {
	r = r.UnsetStatus().(*Bytes)
	if r.Buf == nil {
		return &[]byte{}
	}
	return r.Buf
}

// Copy accepts a parameter and copies the (first) byte in it into its buffer
func (r *Bytes) Copy(b interface{}) Buf {
	r = r.UnsetStatus().(*Bytes)
	bi := big.NewInt(0)
	switch b.(type) {
	case int:
		bi.SetUint64(uint64(b.(int)))
		bb := bi.Bytes()
		r.Buf = &bb
	case uint:
		bi.SetUint64(uint64(b.(uint)))
		bb := bi.Bytes()
		r.Buf = &bb
	case byte:
		r.Buf = &[]byte{b.(byte)}
	case int8:
		r.Buf = &[]byte{byte(b.(int8))}
	case uint16:
		r.Buf = &[]byte{byte(b.(uint16) >> 8), byte(b.(uint16))}
	case int16:
		r.Buf = &[]byte{byte(uint16(b.(int16)) >> 8), byte(uint16(b.(int16)))}
	case uint32:
		r.Buf = &[]byte{
			byte(uint32(b.(uint32)) >> 24),
			byte(uint32(b.(uint32)) >> 16),
			byte(uint32(b.(uint32)) >> 8),
			byte(uint32(b.(uint32)))}
	case int32:
		r.Buf = &[]byte{
			byte(uint32(b.(int32)) >> 24),
			byte(uint32(b.(int32)) >> 16),
			byte(uint32(b.(int32)) >> 8),
			byte(uint32(b.(int32)))}
	case uint64:
		r.Buf = &[]byte{
			byte(b.(uint64) >> 56),
			byte(b.(uint64) >> 48),
			byte(b.(uint64) >> 40),
			byte(b.(uint64) >> 32),
			byte(b.(uint64) >> 24),
			byte(b.(uint64) >> 16),
			byte(b.(uint64) >> 8),
			byte(b.(uint64))}
	case int64:
		r.Buf = &[]byte{
			byte(b.(int64) >> 56),
			byte(b.(int64) >> 48),
			byte(b.(int64) >> 40),
			byte(b.(int64) >> 32),
			byte(b.(int64) >> 24),
			byte(b.(int64) >> 16),
			byte(b.(int64) >> 8),
			byte(b.(int64))}
	case []byte:
		bb := b.([]byte)
		if r.Buf != nil {
			rb := *r.Buf
			for i := range rb {
				rb[i] = 0
			}
		}
		fmt.Println(string(bb))
		rr := make([]byte, len(bb))
		copy(rr, bb)
		r.Buf = &rr
	case *[]byte:
		r.Buf = b.(*[]byte)
	case *Bytes:
		B := b.(*Bytes)
		if r == B {
			return r.SetStatus("copy to self").(*Bytes)
		}
		bb := *B.Buf
		rr := make([]byte, len(bb))
		copy(rr, *B.Buf)
		r.Buf, r.Status, r.Coding = B.Buf, B.Status, B.Coding
		return r
	default:
		r.SetStatus("parameter type not implemented")
		return r
	}
	return r
}

// Free doesn't really do anything but other buffers will need it
func (r *Bytes) Free() Buf {
	r = r.UnsetStatus().(*Bytes)
	r.Buf = nil
	return r
}

// Null is
func (r *Bytes) Null() Buf {
	if r == nil {
		r = r.UnsetStatus().(*Bytes)
		return r
	}
	if r.Buf == nil {
		r.SetStatus("nil buffer")
		return r
	}
	for i := range *r.Buf {
		r.SetElem(i, byte(0))
	}
	r.UnsetStatus()
	r.SetCoding("")
	return r
}

// SetStatus sets the error state
func (r *Bytes) SetStatus(s string) status.Status {
	r = r.UnsetStatus().(*Bytes)
	r.Status = s
	return r
}

// UnsetStatus emptys the error state
func (r *Bytes) UnsetStatus() status.Status {
	if r == nil {
		r = NewBytes()
		r.SetStatus("nil receiver")
	} else {
		r.Status = ""
	}
	return r
}

// GetCoding is
func (r *Bytes) GetCoding() string {
	r = r.UnsetStatus().(*Bytes)
	if r.Coding == "" {
		r.Coding = "string"
	}
	r.SetCoding(r.Coding)
	return r.Coding
}

// SetCoding is
func (r *Bytes) SetCoding(s string) coding.Coding {
	r = r.UnsetStatus().(*Bytes)
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
func (r *Bytes) ListCodings() []string {
	r = r.UnsetStatus().(*Bytes)
	return coding.Codings
}

// Elem returns a byte of 1 or 0 representing a bit
func (r *Bytes) Elem(e int) interface{} {
	r = r.UnsetStatus().(*Bytes)
	if r.Buf == nil {
		r.SetStatus("nil buffer")
		return byte(0)
	}
	if e > len(*r.Buf) {
		r.SetStatus("index out of bounds")
		return byte(0)
	}
	if len(*r.Buf) > 0 {
		return (*r.Buf)[e]
	}
	r.SetStatus("zero length buffer")
	return byte(0)
}

// Len is 1, though we can also read the bits
func (r *Bytes) Len() int {
	r = r.UnsetStatus().(*Bytes)
	if r.Buf == nil {
		return -1
	}
	return len(*r.Buf)
}

// SetElem is
func (r *Bytes) SetElem(e int, val interface{}) arr.Array {
	switch val.(type) {
	case byte:
		if e > len(*r.Buf) {
			r.SetStatus("index out of bounds")
			return r
		}
		if len(*r.Buf) > 0 {
			rr := (*r.Buf)
			rr[e] = val.(byte)
		} else {
			r.SetStatus("zero length buffer")
		}
	default:
		r.SetStatus("parameter not a byte")
	}
	return r
}

// String is
func (r *Bytes) String() (S string) {
	r = r.UnsetStatus().(*Bytes)
	if r.Buf == nil {
		return ""
	}
	return coding.Encode(*r.Buf, r.Coding)
}

// Error is
func (r *Bytes) Error() string {
	return r.Status
}
