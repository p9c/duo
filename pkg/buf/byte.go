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

// NewIf creates a new Byte
func (r *Byte) NewIf() *Byte {
	if r == nil {
		r = NewByte()
		r.SetStatus(er.NilRec)
	}
	return r
}

// Bytes returns a pointer to the buffer
func (r *Byte) Bytes() (out *[]byte) {
	r = r.NewIf()
	out = &[]byte{}
	switch {
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
	r = r.NewIf()
	switch {
	case in == nil:
		r.SetStatus(er.NilParam)
	case len(*in) < 1:
		r.SetStatus(er.ZeroLen)
	default:
		v := make([]byte, len(*in))
		I := *in
		copy(v, I)
		r.Val = &v
	}
	return r
}

// Zero writes zeroes to the byte slice
func (r *Byte) Zero() proto.Buffer {
	r = r.NewIf()
	switch {
	case r.Val == nil:
		r = NewByte()
		r.SetStatus(er.NilBuf)
	default:
		proto.Zero(r.Val)
	}
	return r
}

// Free is a
func (r *Byte) Free() proto.Buffer {
	r = r.NewIf()
	switch {
	default:
		r.UnsetStatus()
		r.Val = nil
	}
	return r
}

// IsEqual returns true if a serialized public key matches this one, also in format (compressed is preferred in a distributed ledger due to size)
func (r *Byte) IsEqual(p *[]byte) (is bool) {
	r = r.NewIf()
	switch {
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
	r = r.NewIf()
	out = &r.Coding
	return
}

// SetCoding is a
func (r *Byte) SetCoding(in string) proto.Coder {
	r = r.NewIf()
	for i := range proto.StringCodings {
		if in == proto.StringCodings[i] {
			r.Coding = in
			break
		}
	}
	return r
}

// ListCodings is a
func (r *Byte) ListCodings() (out *[]string) {
	r = r.NewIf()
	out = &proto.StringCodings
	return
}

// Freeze returns a json format struct of the data
func (r *Byte) Freeze() (out *[]byte) {
	r = r.NewIf()
	var status string
	switch {
	case !r.OK():
		status = ""
	}
	s := []string{
		`{"Val":`,
		`"` + r.String() + `",`,
		`"Status":`,
		`"` + status + `",`,
		`"Coding":`,
		`"` + r.Coding + `"}`,
	}
	b := []byte(strings.Join(s, ""))
	out = &b
	return
}

// Thaw is a
func (r *Byte) Thaw(in *[]byte) proto.Streamer {
	r = r.NewIf()
	out := NewByte()
	if err := json.Unmarshal(*in, out); !out.SetStatusIf(err).OK() {
		r.Zero().Copy(out.Bytes())
	}
	return r
}

// SetStatus is a
func (r *Byte) SetStatus(s string) proto.Status {
	r = r.NewIf()
	switch {
	case s == "":
		r.State.SetStatus("empty status string given")
	default:
		r.State.SetStatus(s)
	}
	return r
}

// SetStatusIf is a
func (r *Byte) SetStatusIf(err error) proto.Status {
	r = r.NewIf()
	if err != nil {
		r.State.SetStatus(err.Error())
	} else {
		r.State.UnsetStatus()
	}
	return r
}

// UnsetStatus is a
func (r *Byte) UnsetStatus() proto.Status {
	r = r.NewIf()
	r.State.UnsetStatus()
	return r
}

// OK returns true if there is no error
func (r *Byte) OK() bool {
	if r == nil {
		r = r.NewIf()
	}
	return r.State.OK()
}

// SetElem is a
func (r *Byte) SetElem(index int, in interface{}) proto.Array {
	r = r.NewIf()
	switch {
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
		r = r.NewIf()
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
	if r == nil || r.Val == nil {
		return -1
	}
	return len(*r.Val)
}

// String implements the stringer, uses coding to determine how the string is contstructed
func (r *Byte) String() (s string) {
	r = r.NewIf()
	switch {
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
