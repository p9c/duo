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

// NewByte creates a new Byte
func NewByte() *Byte {
	r := new(Byte)
	r.Coding = "hex"
	return r
}

// Bytes returns a pointer to the buffer
func (r *Byte) Bytes() (out *[]byte) {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
		fallthrough
	case r.Val == nil:
		r = NewByte().SetStatus(er.NilBuf).(*Byte)
		out = &[]byte{}
	default:
		out = r.Val
	}
	return
}

// Copy copies the byte from a provided byte slice to a new buffer
func (r *Byte) Copy(in *[]byte) proto.Buffer {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
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
func (r *Byte) Zero() proto.Buffer {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	case r.Val == nil:
		r = NewByte().SetStatus(er.NilBuf).(*Byte)
	default:
		b := *r.Val
		for i := range b {
			b[i] = 0
		}
	}
	return r
}

// Free is a
func (r *Byte) Free() proto.Buffer {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	}
	r.UnsetStatus()
	r.Val = nil
	return r
}

// GetCoding is a
func (r *Byte) GetCoding() (out *string) {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	}
	out = &r.Coding
	return
}

// IsEqual returns true if a serialized public key matches this one, also in format (compressed is preferred in a distributed ledger due to size)
func (r *Byte) IsEqual(p *[]byte) (is bool) {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	case r.Len() != len(*p):
		r.SetStatus("buffers are different length")
	case r.Len() < 1:
		r.SetStatus(er.ZeroLenBuf)
		fallthrough
	case len(*p) < 1:
		r.SetStatus(er.ZeroLen)
		fallthrough
	default:
		is = true
		for i := range *p {
			if (*p)[i] != (*r.Bytes())[i] {
				is = false
				break
			}
		}
	}
	return
}

// SetCoding is a
func (r *Byte) SetCoding(in string) proto.Coder {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
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
func (r *Byte) ListCodings() (out *[]string) {
	out = &proto.StringCodings
	return
}

// Freeze returns a json format struct of the data
func (r *Byte) Freeze() (out *[]byte) {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
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
func (r *Byte) Thaw(in *[]byte) proto.Streamer {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	}
	out := NewByte()
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
func (r *Byte) SetStatus(s string) proto.Status {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	} else {
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *Byte) SetStatusIf(err error) proto.Status {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	} else {
		if err != nil {
			r.Status = err.Error()
		}
	}
	return r
}

// UnsetStatus is a
func (r *Byte) UnsetStatus() proto.Status {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	} else {
		r.Status = ""
	}
	return r
}

// OK returns true if there is no error
func (r *Byte) OK() bool {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
		return false
	}
	return r.Status == ""
}

// SetElem is a
func (r *Byte) SetElem(index int, in interface{}) proto.Array {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
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
func (r *Byte) GetElem(index int) (out interface{}) {
	var byt byte
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
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
func (r *Byte) Len() int {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
		return -1
	}
	if r.Val == nil {
		r.SetStatus(er.NilBuf)
		return -1
	}
	return len(*r.Val)
}

// Error implements the Error interface
func (r *Byte) Error() string {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	}
	return r.Status
}

// String implements the stringer, uses coding to determine how the string is contstructed
func (r *Byte) String() string {
	if r == nil {
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	}
	if r.Val == nil {
		return ""
	}
	r.UnsetStatus()
	switch r.Coding {
	case "byte":
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
