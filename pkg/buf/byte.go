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
	out = &[]byte{}
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
	case r.Val == nil:
		r = NewByte()
		r.SetStatus(er.NilBuf)
	default:
		out = r.Val
	}
	return
}

// Copy copies the byte from a provided byte slice to a new buffer
func (r *Byte) Copy(in *[]byte) proto.Buffer {
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
		fallthrough
	case in == nil:
		r.SetStatus(er.NilParam)
	case len(*in) < 1:
		r.SetStatus(er.ZeroLen)
	default:
		v := make([]byte, len(*in))
		copy(v, *in)
		r.Val = &v
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
		proto.Zero(r.Val)
	}
	return r
}

// Free is a
func (r *Byte) Free() proto.Buffer {
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
	default:
		r.UnsetStatus()
		r.Val = nil
	}
	return r
}

// IsEqual returns true if a serialized public key matches this one, also in format (compressed is preferred in a distributed ledger due to size)
func (r *Byte) IsEqual(p *[]byte) (is bool) {
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
	case p == nil:
		r.SetStatus(er.NilParam)
	case len(*p) < 1:
		r.SetStatus(er.ZeroLen)
	case r.Len() != len(*p):
		r.SetStatus("buffers are different length")
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

// GetCoding is a
func (r *Byte) GetCoding() (out *string) {
	if r == nil {
		r = NewByte()
		r.SetStatus(er.NilRec)
	}
	out = &r.Coding
	return
}

// SetCoding is a
func (r *Byte) SetCoding(in string) proto.Coder {
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
	default:
		found := false
		for i := range proto.StringCodings {
			if in == proto.StringCodings[i] {
				found = true
				break
			}
		}
		switch {
		case found != true:
			r.Coding = "hex"
		default:
			r.Coding = in
		}
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
		r = NewByte()
		r.SetStatus(er.NilRec)
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
		r = NewByte()
		r.SetStatus(er.NilRec)
	}
	out := NewByte()
	if err := json.Unmarshal(*in, out); !out.SetStatusIf(err).OK() {
		r.Zero().Copy(out.Bytes())
	}
	return r
}

// SetStatus is a
func (r *Byte) SetStatus(s string) proto.Status {
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
	default:
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *Byte) SetStatusIf(err error) proto.Status {
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
	case err != nil:
		r.Status = err.Error()
	}
	return r
}

// UnsetStatus is a
func (r *Byte) UnsetStatus() proto.Status {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	default:
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
		r = NewByte()
		r.SetStatus(er.NilRec)
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
		r = NewByte()
		r.SetStatus(er.NilRec)
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
	switch {
	case r == nil:
		r = NewByte()
		r.SetStatus(er.NilRec)
		return -1
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
		return -1
	}
	return len(*r.Val)
}

// Error implements the Error interface
func (r *Byte) Error() string {
	if r == nil {
		r = NewByte()
		r.SetStatus(er.NilRec)
	}
	return r.Status
}

// String implements the stringer, uses coding to determine how the string is contstructed
func (r *Byte) String() (s string) {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	case r.Val == nil:
		return ""
	default:
		switch r.Coding {
		case "bytes":
			s = fmt.Sprint(*r.Val)
		case "string":
			s = string(*r.Val)
		case "decimal":
			bi := big.NewInt(0).SetBytes(*r.Val)
			s = fmt.Sprint(bi)
		case "hex":
			s = hex.EncodeToString(*r.Val)
		case "base32":
			s = base32.StdEncoding.EncodeToString(*r.Val)
		case "base58check":
			b := *r.Val
			pre := hex.EncodeToString(b[0:0])
			body := hex.EncodeToString(b[1:])
			var err error
			s, err = base58check.Encode(pre, body)
			r.SetStatusIf(err)
		case "base64":
			dst := make([]byte, len(*r.Val)*4)
			base64.StdEncoding.Encode(dst, *r.Val)
			s = string(dst)
		default:
			r.SetStatus("unrecognised coding")
			r.SetCoding("hex")
			s = hex.EncodeToString(*r.Val)
		}
	}
	return
}
