package buf

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/base58check"
	"github.com/parallelcointeam/duo/pkg/proto"
	"strings"
)

// Byte is a simple single byte
type Byte proto.Byte

var er = proto.Errors

// NewByte creates a new Byte
func NewByte() *Byte {
	r := new(Byte)
	b := byte(0)
	r.Val = &b
	r.Coding = "decimal"
	return r
}

// Bytes returns the byte as a slice with one
func (r *Byte) Bytes(out *[]byte) proto.Buffer {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
		fallthrough
	case r.Val == nil:
		r = NewByte().SetStatus(er.NilBuf).(*Byte)
		fallthrough
	default:
		*out = []byte{*r.Val}
	}
	return r
}

// Copy copies the first byte from a given pointer to a byte slice
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
		b := byte((*in)[0])
		val := &b
		r.Val = val
	}
	return r
}

// Zero makes the value zero
func (r *Byte) Zero() proto.Buffer {
	switch {
	case r == nil:
		r = NewByte().SetStatus(er.NilRec).(*Byte)
	case r.Val == nil:
		r = NewByte().SetStatus(er.NilBuf).(*Byte)
		fallthrough
	default:
		*r.Val = 0
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
func (r *Byte) GetCoding(out *string) proto.Coder {
	*out = r.Coding
	return r
}

// SetCoding is a
func (r *Byte) SetCoding(in string) proto.Coder {
	found := false
	for i := range proto.StringCodings {
		if in == proto.StringCodings[i] {
			found = true
			break
		}
	}
	if found != true {
		r.Coding = "decimal"
	} else {
		r.Coding = in
	}
	return r
}

// ListCodings is a
func (r *Byte) ListCodings(out *[]string) proto.Coder {
	*out = proto.StringCodings
	return r
}

// Freeze is a
func (r *Byte) Freeze(out *[]byte) proto.Streamer {
	if r == nil {
		r = NewByte()
		r.SetStatus("nil receiver")
	}
	s := []string{
		`{"Val":`,
		`` + fmt.Sprint(*r.Val) + `,`,
		`"Status":`,
		`"` + r.Status + `",`,
		`"Coding":`,
		`"` + r.Coding + `"}`,
	}
	b := []byte(strings.Join(s, ""))
	*out = b
	return r
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
	case "bytes":
		return fmt.Sprint([]byte{*r.Val})
	case "string":
		return string([]byte{*r.Val})
	case "decimal":
		return fmt.Sprint(*r.Val)
	case "hex":
		return hex.EncodeToString([]byte{*r.Val})
	case "base32":
		return base32.StdEncoding.EncodeToString([]byte{*r.Val})
	case "base58check":
		s, err := base58check.Encode("00", hex.EncodeToString([]byte{*r.Val}))
		r.SetStatusIf(err)
		return s
	case "base64":
		dst := make([]byte, 8)
		base64.StdEncoding.Encode(dst, []byte{*r.Val, 0, 0, 0})
		return string(dst)
	default:
		r.SetStatus("unrecognised coding")
		r.SetCoding("decimal")
		return fmt.Sprint(*r.Val)
	}
}
