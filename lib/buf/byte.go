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
	buf    byte
	status string
	coding string
}

// NewByte makes a new Byte
func NewByte() *Byte {
	return new(Byte)
}

// Freeze writes the current data structure into a string that format as content of a JSON struct node
func (r *Byte) Freeze() (S string) {
	if r == nil {
		r = NewByte()
	}
	s := []string{
		`"Byte": {`,
		`"Buf":`,
		`"` + fmt.Sprint(r.buf) + `",`,
		`"Status":`,
		`"` + r.status + `",`,
		`"Coding":`,
		`"` + r.coding + `"`,
		"}",
	}
	S = strings.Join(s, " ")
	return
}

// Thaw is
func (r *Byte) Thaw(s string) interface{} {
	var out *Byte
	json.Unmarshal([]byte(s), out)
	return out
}

// Data is
func (r *Byte) Data() interface{} {
	return r.buf
}

// Copy accepts a parameter and copies the (first) byte in it into its buffer
func (r *Byte) Copy(b interface{}) Buf {
	if r == nil {
		r = NewByte()
		r.SetStatus("nil receiver")
	}
	switch b.(type) {
	case byte:
		r.buf = b.(byte)
	case int8:
		r.buf = byte(b.(int8))
	case uint16:
		r.buf = byte(b.(uint16))
	case int16:
		r.buf = byte(b.(int16))
	case uint32:
		r.buf = byte(b.(uint32))
	case int32:
		r.buf = byte(b.(int32))
	case uint64:
		r.buf = byte(b.(uint64))
	case int64:
		r.buf = byte(b.(int64))
	case []byte:
		if len(b.([]byte)) > 0 {
			r.buf = b.([]byte)[0]
		}
	case *[]byte:
		if len(*b.(*[]byte)) > 0 {
			r.buf = (*b.(*[]byte))[0]
		}
	default:
		r.SetStatus("parameter type not implemented")
	}
	return r
}

// Free doesn't really do anything but other buffers will need it
func (r *Byte) Free() Buf {
	if r == nil {
		r = NewByte()
		r.SetStatus("nil receiver")
	} else {
		r.UnsetStatus()
	}
	return r
}

// Null is
func (r *Byte) Null() Buf {
	if r == nil {
		r = NewByte()
		r.SetStatus("nil receiver")
	} else {
		r.UnsetStatus()
		r.buf = 0
		r.status = coding.Codings[0]
		r.coding = ""
	}
	return r
}

// SetStatus sets the error state
func (r *Byte) SetStatus(s string) status.Status {
	r.status = s
	return r
}

// UnsetStatus emptys the error state
func (r *Byte) UnsetStatus() status.Status {
	r.status = ""
	return r
}

// GetCoding is
func (r *Byte) GetCoding() string {
	return r.coding
}

// SetCoding is
func (r *Byte) SetCoding(s string) coding.Coding {
	found := false
	for i := range coding.Codings {
		if s == coding.Codings[i] {
			found = true
		}
	}
	if found {
		r.coding = s
	} else {
		r.coding = "hex"
	}
	return r
}

// ListCodings is
func (r *Byte) ListCodings() []string {
	return coding.Codings
}

// Elem is
func (r *Byte) Elem(int) interface{} {
	panic("not implemented")
}

// Len is
func (r *Byte) Len() int {
	panic("not implemented")
}

// SetElem is
func (r *Byte) SetElem(int, interface{}) array.Array {
	panic("not implemented")
}

// String is
func (r *Byte) String() string {
	panic("not implemented")
}

// Error is
func (r *Byte) Error() string {
	return r.status
}
