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
		B, err := memguard.NewMutable(len(*in))
		r.SetStatusIf(err)
		if r != nil {
			b := B.Buffer()
			for i := range *in {
				b[i] = (*in)[i]
			}
			r.Val = B
		}
	}
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
	return r

}

// Free is a
func (r *Secure) Free() proto.Buffer {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	}
	r.UnsetStatus()
	if r.Val != nil {
		r.Val.Destroy()
	}
	r.Val = nil
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
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	}
	if r.Val != nil {
		r.Val.Destroy()
	}
	var err error
	r.Val, err = memguard.NewMutableRandom(length)
	if r.SetStatusIf(err); err != nil {
		return r
	}
	return r
}

// GetCoding is a
func (r *Secure) GetCoding() (out *string) {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	}
	out = &r.Coding
	return
}

// SetCoding is a
func (r *Secure) SetCoding(in string) proto.Coder {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
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
func (r *Secure) ListCodings() (out *[]string) {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	}
	out = &proto.StringCodings
	return
}

// Freeze returns a json format struct of the data
func (r *Secure) Freeze() (out *[]byte) {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
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
func (r *Secure) Thaw(in *[]byte) proto.Streamer {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	}
	out := NewSecure()
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
func (r *Secure) SetStatus(s string) proto.Status {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	} else {
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *Secure) SetStatusIf(err error) proto.Status {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	} else {
		if err != nil {
			r.Status = err.Error()
		}
	}
	return r
}

// UnsetStatus is a
func (r *Secure) UnsetStatus() proto.Status {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	} else {
		r.Status = ""
	}
	return r
}

// OK returns true if there is no error
func (r *Secure) OK() bool {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
		return false
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
	var byt byte
	switch {
	case r == nil:
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
		out = &byt
	case r.Val == nil:
		r.SetStatus(er.NilBuf)
		out = &byt
	case index > r.Len():
		r.SetStatus(er.OutOfBounds)
		out = &byt
	default:
		out = &r.Val.Buffer()[index]
	}
	return
}

// Len is a
func (r *Secure) Len() (length int) {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
		return -1
	}
	if r.Val == nil {
		r.SetStatus(er.NilBuf)
		return -1
	}
	return r.Val.Size()
}

// Error implements the Error interface
func (r *Secure) Error() string {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	}
	return r.Status
}

// String implements the stringer, uses coding to determine how the string is contstructed
func (r *Secure) String() string {
	if r == nil {
		r = NewSecure().SetStatus(er.NilRec).(*Secure)
	}
	if r.Val == nil {
		return ""
	}
	r.UnsetStatus()
	switch r.Coding {
	case "bytes":
		return fmt.Sprint(*r.Val)
	case "string":
		return string(r.Val.Buffer())
	case "decimal":
		bi := big.NewInt(0)
		bi.SetBytes(r.Val.Buffer())
		return fmt.Sprint(bi)
	case "hex":
		return hex.EncodeToString(r.Val.Buffer())
	case "base32":
		return base32.StdEncoding.EncodeToString(r.Val.Buffer())
	case "base58check":
		b := r.Val.Buffer()
		pre := hex.EncodeToString(b[0:0])
		body := hex.EncodeToString(b[1:])
		s, err := base58check.Encode(pre, body)
		r.SetStatusIf(err)
		return s
	case "base64":
		dst := make([]byte, r.Val.Size()*4)
		base64.StdEncoding.Encode(dst, r.Val.Buffer())
		return string(dst)
	default:
		r.SetStatus("unrecognised coding")
		r.SetCoding("decimal")
		return fmt.Sprint(*r.Val)
	}
}
