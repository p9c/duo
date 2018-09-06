// Package passbuf is a SecBuf with string conversion functions, which are not recommended to be used. However the distinct type adds other safety benefits for the programmer and can be treated exactly the same as the SecBuf otherwise.
package passbuf

import (
	"encoding/json"
	"gitlab.com/parallelcoin/duo/pkg/buf/sec"
)

// Password is just a SecBuf but recommended for use in storing password fields in other structures to signify what kind of data it is.
type Password struct {
	*secbuf.SecBuf
}

// guards against nil pointer receivers
func nilError(s string) *Password {
	r := new(Password)
	r.SecBuf = secbuf.New()
	r.SetError(s + " nil receiver")
	return r
}

// NewPassword creates a new Password
func New(r ...*Password) *Password {
	if len(r) == 0 {
		r = append(r, new(Password))
		r[0].SecBuf = secbuf.New()
	}
	if r[0] == nil {
		r[0] = new(Password)
		r[0].SecBuf = secbuf.New()
	}
	r[0].SetCoding("string")
	return r[0]
}

// FromString loads the Lockedbuffer with the bytes of a string. The string is immutable so it is not removed from memory except automatically.
func (r *Password) FromString(s string) *Password {
	if r == nil {
		r = New(r)
	}
	r.SecBuf = r.New(len(s)).(*secbuf.SecBuf)
	for i := range s {
		r.SetElem(i, s[i])
	}
	return r
}

// MarshalJSON marshalls the JSON for a Password
func (r *Password) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = nilError("MarshalJSON()")
	}
	if r.SecBuf == nil {
		r.SecBuf = secbuf.New()
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  string(*(r.SecBuf.Buf()).(*[]byte)),
		Coding: r.Coding(),
		Error:  r.Error(),
	})
}

// Stringer implementation

// String returns value encoded according to the current coding mode
func (r *Password) String() string {
	if nil == r {
		r = nilError("String()")
	}
	return r.SecBuf.String()
}
