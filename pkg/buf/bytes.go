package buf

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/base58check"
	"github.com/parallelcointeam/duo/pkg/proto"
	"math/big"
	"strings"
)

// Bytes is a simple single byte
type Bytes proto.Bytes

// NewBytes creates a new Bytes
func NewBytes() *Bytes {
	r := new(Bytes)
	r.Coding = "hex"
	return r
}

// Bytes returns a pointer to the buffer
func (r *Bytes) Bytes() (out *[]byte) {
	switch {
	case r == nil:
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
		fallthrough
	case r.Val == nil:
		r = NewBytes().SetStatus(er.NilBuf).(*Bytes)
		out = &[]byte{}
	default:
		out = r.Val
	}
	return
}

// Copy copies the bytes from a provided byte slice to a new buffer
func (r *Bytes) Copy(in *[]byte) proto.Buffer {
	switch {
	case r == nil:
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
		fallthrough
	case in == nil:
		r.SetStatus(er.NilParam)
	case len(*in) == 0:
		r.SetStatus(er.ZeroLen)
	case len(*in) > 0:
		b := make([]byte, len(*in))
		for i := range *in {
			b[i] = (*in)[i]
		}
		r.Val = &b
	}
	return r
}

// Zero writes zeroes to the byte slice
func (r *Bytes) Zero() proto.Buffer {
	switch {
	case r == nil:
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	case r.Val == nil:
		r = NewBytes().SetStatus(er.NilBuf).(*Bytes)
	default:
		b := *r.Val
		for i := range b {
			b[i] = 0
		}
	}
	return r
}

// Free is a
func (r *Bytes) Free() proto.Buffer {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	r.UnsetStatus()
	r.Val = nil
	return r
}

// GetCoding is a
func (r *Bytes) GetCoding() (out *string) {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	out = &r.Coding
	return
}

// SetCoding is a
func (r *Bytes) SetCoding(in string) proto.Coder {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	found := false
	for i := range proto.StringCodings {
		if in == proto.StringCodings[i] {
			found = true
			break
		}
	}
	if found != true {
		r.Coding = "hex"
	} else {
		r.Coding = in
	}
	return r
}

// ListCodings is a
func (r *Bytes) ListCodings() (out *[]string) {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	out = &proto.StringCodings
	return
}

// Freeze returns a json format struct of the data
func (r *Bytes) Freeze() (out *[]byte) {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	s := []string{
		`{"Val":`,
		`"` + r.String() + `",`,
		`"Status":`,
		`"` + r.Status + `",`,
		`"Coding":`,
		`"` + r.Coding + `"}`,
	}
	b := []byte(strings.Join(s, ""))
	out = &b
	return
}

// Thaw is a
func (r *Bytes) Thaw(in *[]byte) proto.Streamer {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	out := NewBytes()
	err := json.Unmarshal(*in, out)
	out.SetStatusIf(err)
	if r.Status != "" {
		return r
	}
	r.Zero()
	r = out
	return r
}

// SetStatus is a
func (r *Bytes) SetStatus(s string) proto.Status {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	} else {
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *Bytes) SetStatusIf(err error) proto.Status {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	} else {
		if err != nil {
			r.Status = err.Error()
		}
	}
	return r
}

// UnsetStatus is a
func (r *Bytes) UnsetStatus() proto.Status {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	} else {
		r.Status = ""
	}
	return r
}

// OK returns true if there is no error
func (r *Bytes) OK() bool {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
		return false
	}
	return r.Status == ""
}

// SetElem is a
func (r *Bytes) SetElem(index int, in interface{}) proto.Array {
	switch {
	case r == nil:
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
	case index > r.Len():
		r.SetStatus(er.OutOfBounds)
	default:
		switch in.(type) {
		case *byte:
			(*r.Val)[index] = *in.(*byte)
		default:
			r.SetStatus(er.InvalidType)
		}
	}
	return r
}

// GetElem is a
func (r *Bytes) GetElem(index int) (out interface{}) {
	var byt byte
	switch {
	case r == nil:
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
		out = &byt
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
		out = &byt
	case index > r.Len():
		r.SetStatus(er.OutOfBounds)
		out = &byt
	default:
		out = &(*r.Val)[index]
	}
	return
}

// Len is a
func (r *Bytes) Len() int {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
		return -1
	}
	if r.Val == nil {
		r.SetStatus(er.NilBuf)
		return -1
	}
	return len(*r.Val)
}

// Error implements the Error interface
func (r *Bytes) Error() string {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	return r.Status
}

// String implements the stringer, uses coding to determine how the string is contstructed
func (r *Bytes) String() string {
	if r == nil {
		r = NewBytes().SetStatus(er.NilRec).(*Bytes)
	}
	if r.Val == nil {
		return ""
	}
	r.UnsetStatus()
	switch r.Coding {
	case "bytes":
		return fmt.Sprint(*r.Val)
	case "string":
		return string(*r.Val)
	case "decimal":
		bi := big.NewInt(0)
		bi.SetBytes(*r.Val)
		return fmt.Sprint(bi)
	case "hex":
		return hex.EncodeToString(*r.Val)
	case "base32":
		return base32.StdEncoding.EncodeToString(*r.Val)
	case "base58check":
		b := *r.Val
		pre := hex.EncodeToString(b[0:0])
		body := hex.EncodeToString(b[1:])
		s, err := base58check.Encode(pre, body)
		r.SetStatusIf(err)
		return s
	case "base64":
		dst := make([]byte, len(*r.Val)*4)
		base64.StdEncoding.Encode(dst, *r.Val)
		return string(dst)
	default:
		r.SetStatus("unrecognised coding")
		r.SetCoding("hex")
		return hex.EncodeToString(*r.Val)
	}
}
