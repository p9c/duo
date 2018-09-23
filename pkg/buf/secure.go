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

// Bytes returns a pointer to the buffer
func (r *Secure) Bytes() (out *[]byte) {
	switch {
	case r == nil:
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
		fallthrough
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
	switch {
	case r == nil:
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
		fallthrough
	case in == nil:
		r.SetStatus(er.NilParam)
	case len(*in) == 0:
		r.SetStatus(er.ZeroLen)
	case len(*in) > 0:
		b := make([]byte, len(*in))
		copy(b, *in)
		B, err := memguard.NewMutableFromBytes(b)
		if r.SetStatusIf(err).OK() {
			r.Val = B
		}
	}
	r.UnsetStatus()
	return r
}

// Zero is a
func (r *Secure) Zero() proto.Buffer {
	switch {
	case r == nil:
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
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
	switch {
	case r == nil:
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
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
	switch {
	case r == nil:
		r = NewSecure()
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

// Rand creates a secure buffer containing cryptographically secure random bytes
func (r *Secure) Rand(length int) *Secure {
	switch {
	case r == nil:
		r = NewSecure()
		r.SetStatus(er.NilRec)
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
	if r == nil {
		r = NewSecure()
		r.SetStatus(er.NilRec)
	}
	out = &r.Coding
	return
}

// SetCoding sets the encoding for the stringer
func (r *Secure) SetCoding(in string) proto.Coder {
	if r == nil {
		r = NewSecure()
		r.SetStatus(er.NilRec)
	}
	r.Coding = "hex"
	found := false
	for i := range proto.StringCodings {
		if in == proto.StringCodings[i] {
			found = true
			break
		}
	}
	if found {
		r.Coding = in
	}
	return r
}

// ListCodings returns the set of codings available
func (r *Secure) ListCodings() (out *[]string) {
	if r == nil {
		r = NewSecure()
		r.SetStatus(er.NilRec)
	}
	out = &proto.StringCodings
	return
}

// Freeze returns a json format struct of the data
func (r *Secure) Freeze() (out *[]byte) {
	if r == nil {
		r = NewSecure()
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
	r.UnsetStatus()
	return
}

// Thaw turns a json representation back into a variable
func (r *Secure) Thaw(in *[]byte) proto.Streamer {
	if r == nil {
		r = NewSecure()
		r.SetStatus(er.NilRec)
	}
	out := NewSecure()
	if err := json.Unmarshal(*in, out); !out.SetStatusIf(err).OK() {
		r.Zero()
		r = out
	}
	return r
}

// SetStatus sets the status of an object after an operation
func (r *Secure) SetStatus(s string) proto.Status {
	switch {
	case r == nil:
		r = NewSecure()
		r.SetStatus(er.NilRec)
	default:
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *Secure) SetStatusIf(err error) proto.Status {
	switch {
	case r == nil:
		r = NewSecure()
		r.SetStatus(er.NilRec)
	case err != nil:
		r.Status = err.Error()
	default:
		r.UnsetStatus()
	}
	return r
}

// UnsetStatus is a
func (r *Secure) UnsetStatus() proto.Status {
	switch {
	case r == nil:
		r = NewSecure()
		r.Status = er.NilRec
	default:
		r.Status = ""
	}
	return r
}

// OK returns true if there is no error
func (r *Secure) OK() bool {
	if r == nil {
		r = NewSecure()
		return r.SetStatus(er.NilRec).OK()
	}
	return r.Status == ""
}

// SetElem is a
func (r *Secure) SetElem(index int, in interface{}) proto.Array {
	switch {
	case r == nil:
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
	case index > r.Len():
		r.SetStatus(er.OutOfBounds)
	default:
		switch in.(type) {
		case *byte:
			b := r.Val.Buffer()
			b[index] = *in.(*byte)
		default:
			r.SetStatus(er.InvalidType)
		}
	}
	return r
}

// GetElem is a
func (r *Secure) GetElem(index int) (out interface{}) {
	o := byte(0)
	out = &o
	switch {
	case r == nil:
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
	case index > r.Len():
		r.SetStatus(er.OutOfBounds)
	default:
		out = &r.Val.Buffer()[index]
	}
	return
}

// Len is a
func (r *Secure) Len() (length int) {
	switch {
	case r == nil:
		r = NewSecure()
		r.SetStatus(er.NilRec)
		return -1
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
		return -1
	default:
		return r.Val.Size()
	}
}

// Error implements the Error interface
func (r *Secure) Error() string {
	if r == nil {
		r = NewSecure()
		r.SetStatus(er.NilRec)
	}
	return r.Status
}

// String implements the stringer, uses coding to determine how the string is contstructed
func (r *Secure) String() (s string) {
	switch {
	case r == nil:
		r = NewSecure()
		r.SetStatus(er.NilRec)
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
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
