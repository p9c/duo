package buf

import (
	"encoding/json"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/lib/array"
	"gitlab.com/parallelcoin/duo/lib/coding"
	"gitlab.com/parallelcoin/duo/lib/status"
	"math/big"
)

type bsv []byte

func (r bsv) Addr() *[]byte {
	R := []byte(r)
	return &R
}

// Secure is a simple buffer for slices of bytes and can ingest almost any other kind of data if one were so inclined.
type Secure struct {
	Buf    *[]byte
	Status string
	Coding string
	lb     *memguard.LockedBuffer
}

// NewSecure makes a new Secure
func NewSecure() *Secure {
	b := new(Secure)
	b.Coding = "base64"
	return b
}

// Freeze writes the current data structure into a string that format as content of a JSON struct node
func (r *Secure) Freeze() (S string) {
	if r == nil {
		r = NewSecure()
		r.SetStatus("nil receiver")
	}
	coding := r.Coding
	r.SetCoding("base64")
	S = `{"Buf":"` + r.String() + `",` +
		`"Status":"` + r.Status + `",` +
		`"Coding":"` + coding + `"}`
	r.SetCoding(coding)
	return
}

// Thaw takes a frozen Secure buffer and returns the structure it represents
func (r *Secure) Thaw(s string) interface{} {
	var err error
	if r == nil {
		r = NewSecure()
	} else {
		r.Null().Free()
	}
	err = json.Unmarshal([]byte(s), r)
	if err != nil {
		return r.SetStatusIf(err).(*Secure)
	}
	if len(*r.Buf) == 0 {
		r.SetStatus("zero length buffer")
		r.Buf = nil
		if r.lb != nil {
			r.lb.Destroy()
		}
		r.lb = nil
	} else {
		r.lb, err = memguard.NewMutable(len(*r.Buf))
		r.SetStatusIf(err)
		r.lb.Move(*r.Buf)
		r.Buf = bsv(r.lb.Buffer()).Addr()
	}
	return r.SetStatusIf(err).(*Secure)
}

// Data returns the content of the buffer
func (r *Secure) Data(out ...interface{}) interface{} {
	r = r.UnsetStatus().(*Secure)
	if r.lb == nil {
		return &[]byte{}
	}
	return bsv(r.lb.Buffer()).Addr()
}

// Copy accepts a parameter and copies the (first) byte in it into its buffer
func (r *Secure) Copy(b interface{}) Buf {
	var err error
	bi := big.NewInt(0)
	r = r.UnsetStatus().(*Secure)
	if b == nil {
		return r.SetStatus("nil interface").(*Secure)
	}
	if r.lb != nil {
		r.lb.Destroy()
	}
	switch b.(type) {
	case nil:
		return r.SetStatus("nil parameter").(*Secure)
	case int:
		bi.SetUint64(uint64(b.(int)))
		r.lb, err = memguard.NewMutable(len(bi.Bytes()))
		if err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		r.lb.Copy(bi.Bytes())
	case uint:
		bi.SetUint64(uint64(b.(uint)))
		r.lb, err = memguard.NewMutable(len(bi.Bytes()))
		if err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		r.lb.Copy(bi.Bytes())
	case byte:
		if r.lb, err = memguard.NewMutable(1); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		(*r.Buf)[0] = b.(byte)
	case int8:
		if r.lb, err = memguard.NewMutable(1); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		(*r.Buf)[0] = byte(b.(int8))
	case uint16:
		if r.lb, err = memguard.NewMutable(2); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		(*r.Buf)[0], (*r.Buf)[1] = byte(b.(uint16)), byte(b.(uint16)>>8)
	case int16:
		if r.lb, err = memguard.NewMutable(2); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		s, d := b.(int16), *r.Buf
		d[0], d[1] = byte(s), byte(s>>8)
	case uint32:
		if r.lb, err = memguard.NewMutable(4); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		s, d := b.(uint32), *r.Buf
		d[0] = byte(uint32(s))
		d[1] = byte(uint32(s) >> 8)
		d[2] = byte(uint32(s) >> 16)
		d[3] = byte(uint32(s) >> 24)
	case int32:
		if r.lb, err = memguard.NewMutable(4); err != nil {
			r.SetStatus(err.Error())
			return r
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		s, d := b.(int32), *r.Buf
		d[0] = byte(uint32(s))
		d[1] = byte(uint32(s) >> 8)
		d[2] = byte(uint32(s) >> 16)
		d[3] = byte(uint32(s) >> 24)
	case uint64:
		if r.lb, err = memguard.NewMutable(8); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		s, d := b.(uint64), *r.Buf
		d[0] = byte(s)
		d[1] = byte(s >> 8)
		d[2] = byte(s >> 16)
		d[3] = byte(s >> 24)
		d[4] = byte(s >> 32)
		d[5] = byte(s >> 40)
		d[6] = byte(s >> 48)
		d[7] = byte(s >> 56)
	case int64:
		if r.lb, err = memguard.NewMutable(8); err != nil {
			r.SetStatus(err.Error())
			return r
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		s, d := b.(int64), *r.Buf
		d[0] = byte(s)
		d[1] = byte(s >> 8)
		d[2] = byte(s >> 16)
		d[3] = byte(s >> 24)
		d[4] = byte(s >> 32)
		d[5] = byte(s >> 40)
		d[6] = byte(s >> 48)
		d[7] = byte(s >> 56)
	case []byte:
		bb := b.([]byte)
		if r.lb, err = memguard.NewMutable(len(bb)); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		r.lb.Copy(bb)
	case *[]byte:
		bv := *b.(*[]byte)
		if r.lb, err = memguard.NewMutable(len(bv)); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		r.lb.Copy(bv)
	case *Secure:
		B := b.(*Secure)
		if r == B {
			return r.SetStatus("copy to self").(*Secure)
		}
		if r.lb, err = memguard.NewMutableFromBytes(B.lb.Buffer()); err != nil {
			return r.SetStatus(err.Error()).(*Secure)
		}
		r.Buf = bsv(r.lb.Buffer()).Addr()
		r.Status, r.Coding = B.Status, B.Coding
	default:
		r.SetStatus("parameter type not implemented")
	}
	return r
}

// Free destroys the buffer, wiping it as well
func (r *Secure) Free() Buf {
	r = r.UnsetStatus().(*Secure)
	if r.lb != nil {
		r.lb.Destroy()
		r.lb = nil
	}
	r.Buf = nil
	return r
}

// Null sets all the bytes in the buffer to zero
func (r *Secure) Null() Buf {
	r = r.UnsetStatus().(*Secure)
	if r.Buf == nil {
		return r.SetStatus("nil buffer").(*Secure)
	}
	r.lb.Wipe()
	r.SetCoding("")
	return r
}

// SetStatus sets the error state
func (r *Secure) SetStatus(s string) status.Status {
	r = r.UnsetStatus().(*Secure)
	r.Status = s
	return r
}

// SetStatusIf sets an error from a standard error interface variable if it is set
func (r *Secure) SetStatusIf(err error) status.Status {
	if err != nil {
		r.Status = err.Error()
	} else {
		r.Status = ""
	}
	return r
}

// UnsetStatus emptys the error state
func (r *Secure) UnsetStatus() status.Status {
	if r == nil {
		r = NewSecure().SetStatus("nil receiver").(*Secure)
	} else {
		r.Status = ""
	}
	return r
}

// GetCoding is
func (r *Secure) GetCoding() string {
	r = r.UnsetStatus().(*Secure)
	return r.SetCoding(r.Coding).(*Secure).Coding
}

// SetCoding is
func (r *Secure) SetCoding(s string) coding.Coding {
	r = r.UnsetStatus().(*Secure)
	found := false
	for i := range coding.Codings {
		if s == coding.Codings[i] {
			found = true
			break
		}
	}
	if found {
		r.Coding = s
	} else {
		r.Coding = "base64"
	}
	return r
}

// ListCodings is
func (r *Secure) ListCodings() []string {
	r = r.UnsetStatus().(*Secure)
	return coding.Codings
}

// Elem returns a byte of 1 or 0 representing a bit
func (r *Secure) Elem(e int) interface{} {
	r = r.UnsetStatus().(*Secure)
	switch {
	case r.Buf == nil:
		r.SetStatus("nil buffer")
	case r.lb.Size() < 1:
		r.SetStatus("zero length buffer")
	case e > r.lb.Size():
		r.SetStatus("index out of bounds")
	case e < 1:
		r.SetStatus("zero length buffer")
	default:
		return (*r.Buf)[e]
	}
	return byte(0)
}

// Len is 1, though we can also read the bits
func (r *Secure) Len() int {
	r = r.UnsetStatus().(*Secure)
	if r.Buf == nil {
		return -1
	}
	return r.lb.Size()
}

// SetElem is
func (r *Secure) SetElem(e int, val interface{}) arr.Array {
	r = r.UnsetStatus().(*Secure)
	switch val.(type) {
	case byte:
		switch {
		case e > r.lb.Size():
			r.SetStatus("index out of bounds")
		case r.lb.Size() < 1:
			r.SetStatus("zero length buffer")
		case r.lb.Size() > 0:
			rr := *r.Buf
			rr[e] = val.(byte)
		}
	case int:
		switch {
		case e > r.lb.Size():
			r.SetStatus("index out of bounds")
		case r.lb.Size() < 1:
			r.SetStatus("zero length buffer")
		case r.lb.Size() > 0:
			rr := *r.Buf
			rr[e] = byte(val.(int))
		}
	default:
		r.SetStatus("parameter not a byte")
	}
	return r
}

// String is
func (r *Secure) String() (S string) {
	r = r.UnsetStatus().(*Secure)
	if r.lb == nil {
		return ""
	}
	return coding.Encode(r.lb.Buffer(), r.Coding)
}

// Error is
func (r *Secure) Error() string {
	return r.Status
}

// MarshalJSON renders the data as JSON
func (r *Secure) MarshalJSON() ([]byte, error) {
	var v, c, e string
	if r == nil {
		v, c, e = "", "base64", "nil receiver"
	} else {
		v = coding.Encode(r.lb.Buffer(), "base64")
		c = r.Coding
		e = r.Status
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  v,
		Coding: c,
		Error:  e,
	})
}
