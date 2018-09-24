package buf

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/base58check"
	"github.com/awnumar/memguard"
	"github.com/parallelcointeam/duo/pkg/proto"
	"math/big"
	"strings"
)

// NewSecure creates a new Secure
func NewSecure() *Secure {
	r := new(Secure)
	r.Coding = "hex"
	return r
}

// NewIf creates a new Byte
func (r *Secure) NewIf() *Secure {
	if r == nil {
		r = NewSecure()
		r.SetStatus(er.NilRec)
	}
	return r
}

// Bytes returns a pointer to the buffer
func (r *Secure) Bytes() (out *[]byte) {
	r = r.NewIf()
	switch {
	case r.Val == nil:
		r = NewSecure().SetStatus(er.NilBuf).(*Secure)
		out = &[]byte{}
	default:
		b := r.Val.Buffer()
		out = &b
		r.UnsetStatus()
	}
	return
}

// Copy is a
func (r *Secure) Copy(in *[]byte) proto.Buffer {
	r = r.NewIf()
	switch {
	case in == nil:
		r.SetStatus(er.NilParam)
	case len(*in) == 0:
		r.SetStatus(er.ZeroLen)
	case len(*in) > 0:
		b := make([]byte, len(*in))
		copy(b, *in)
		B, err := memguard.NewMutableFromBytes(b)
		if r.SetStatusIf(err); B != nil {
			r.Val = B
		}
	}
	r.UnsetStatus()
	return r
}

// Zero is a
func (r *Secure) Zero() proto.Buffer {
	r = r.NewIf()
	switch {
	case r.Val == nil:
		r = NewSecure().SetStatus(er.NilBuf).(*Secure)
	default:
		r.Val.Wipe()
	}
	r.UnsetStatus()
	return r

}

// Free is a
func (r *Secure) Free() proto.Buffer {
	r = r.NewIf()
	switch {
	case r.Val != nil:
		r.Val.Destroy()
		fallthrough
	default:
		r.Val = nil
		r.UnsetStatus()
	}
	return r
}

// IsEqual returns true if a serialized public key matches this one, also in format (compressed is preferred in a distributed ledger due to size).
func (r *Secure) IsEqual(p *[]byte) (is bool) {
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

// Rand creates a secure buffer containing cryptographically secure random bytes
func (r *Secure) Rand(length int) *Secure {
	r = r.NewIf()
	switch {
	case r.Val != nil:
		r.Val.Destroy()
	default:
		var err error
		r.Val, err = memguard.NewMutableRandom(length)
		r.SetStatusIf(err)
	}
	return r
}

// GetCoding is a
func (r *Secure) GetCoding() (out *string) {
	r = r.NewIf()
	out = &r.Coding
	return
}

// SetCoding sets the encoding for the stringer
func (r *Secure) SetCoding(in string) proto.Coder {
	r = r.NewIf()
	r.Coding = "hex"
	for i := range proto.StringCodings {
		if in == proto.StringCodings[i] {
			r.Coding = in
			break
		}
	}
	return r
}

// ListCodings returns the set of codings available
func (r *Secure) ListCodings() (out *[]string) {
	r = r.NewIf()
	out = &proto.StringCodings
	return
}

// Freeze returns a json format struct of the data
func (r *Secure) Freeze() (out *[]byte) {
	r = r.NewIf()
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
	r.UnsetStatus()
	return
}

// Thaw turns a json representation back into a variable
func (r *Secure) Thaw(in *[]byte) proto.Streamer {
	r = r.NewIf()
	out := NewSecure()
	if err := json.Unmarshal(*in, out); !out.SetStatusIf(err).OK() {
		r.Zero()
		r = out
	}
	return r
}

// SetStatus sets the status of an object after an operation
func (r *Secure) SetStatus(s string) proto.Status {
	r = r.NewIf()
	r.State.SetStatus(s)
	return r
}

// SetStatusIf is a
func (r *Secure) SetStatusIf(err error) proto.Status {
	r = r.NewIf()
	switch {
	case err != nil:
		r.State.SetStatus(err.Error())
	default:
		r.UnsetStatus()
	}
	return r
}

// UnsetStatus is a
func (r *Secure) UnsetStatus() proto.Status {
	r = r.NewIf()
	r.State.SetStatus("")
	return r
}

// OK returns true if there is no error
func (r *Secure) OK() bool {
	r = r.NewIf()
	return r.Status == ""
}

// SetElem sets a byte in the buffer to a given value if it is in bounds
func (r *Secure) SetElem(index int, in interface{}) proto.Array {
	r = r.NewIf()
	var elem byte
	switch {
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
	case index > r.Len():
		r.SetStatus(er.OutOfBounds)
	default:
		switch in.(type) {
		case byte:
			elem = in.(byte)
		case int8:
			elem = byte(in.(int8))
		case int:
			elem = byte(in.(int))
		case uint:
			elem = byte(in.(uint))
		case uint16:
			elem = byte(in.(uint16))
		case int16:
			elem = byte(in.(int16))
		case uint32:
			elem = byte(in.(uint32))
		case int32:
			elem = byte(in.(int32))
		case uint64:
			elem = byte(in.(uint64))
		case int64:
			elem = byte(in.(int64))
		default:
			r.SetStatus(er.InvalidType)
		}
		b := r.Val.Buffer()
		b[index] = elem
	}
	return r
}

// GetElem returns the byte at a given position if it is in bounds
func (r *Secure) GetElem(index int) (out interface{}) {
	r = r.NewIf()
	o := byte(0)
	out = &o
	switch {
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
	case index > r.Len():
		r.SetStatus(er.OutOfBounds)
	default:
		out = &r.Val.Buffer()[index]
	}
	return
}

// Len returns the length of the buffer or -1 if it is not allocated
func (r *Secure) Len() (length int) {
	switch {
	case r == nil:
		return -1
	case r.Val == nil:
		return -1
	default:
		return r.Val.Size()
	}
}

// Error implements the Error interface
func (r *Secure) Error() string {
	r = r.NewIf()
	return r.Status
}

// String implements the stringer, uses coding to determine how the string is contstructed
func (r *Secure) String() (s string) {
	switch {
	case r == nil || r.Val == nil:
		s = "<nil>"
	default:
		switch r.Coding {
		case "bytes":
			s = fmt.Sprint(*r.Val)
		case "string":
			s = string(r.Val.Buffer())
		case "decimal":
			bi := big.NewInt(0)
			bi.SetBytes(r.Val.Buffer())
			s = fmt.Sprint(bi)
		case "hex":
			s = hex.EncodeToString(r.Val.Buffer())
		case "base32":
			s = base32.StdEncoding.EncodeToString(r.Val.Buffer())
		case "base58check":
			b := r.Val.Buffer()
			pre := hex.EncodeToString(b[0:0])
			body := hex.EncodeToString(b[1:])
			var err error
			s, err = base58check.Encode(pre, body)
			r.SetStatusIf(err)
		case "base64":
			dst := make([]byte, r.Val.Size()*4)
			base64.StdEncoding.Encode(dst, r.Val.Buffer())
			s = string(dst)
		default:
			r.SetStatus("unrecognised coding")
			r.SetCoding("decimal")
			s = fmt.Sprint(*r.Val)
		}
	}
	return s
}
