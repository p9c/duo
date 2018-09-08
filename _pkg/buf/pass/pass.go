// Package pass is a Fenced with string conversion functions, which are not recommended to be used. However the distinct type adds other safety benefits for the programmer and can be treated exactly the same as the Fenced otherwise.
package pass

import (
	"encoding/json"
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
)

// Password is just a Fenced but recommended for use in storing password fields in other structures to signify what kind of data it is.
type Password struct {
	buf.Fenced
}

// New creates a new Password
func New() (R *Password) {
	R = new(Password)
	R.Fenced = *buf.NewFenced()
	R.Status().Set("string")
	return
}

// FromString loads the Lockedbuffer with the bytes of a string. The string is immutable so it is not removed from memory except automatically.
func (r *Password) FromString(s string) *Password {
	if r == nil {
		r = new(Password)
		r.Fenced = *buf.NewFenced()
		r.Fenced.Status().Set("FromString() nil receiver")
	}
	r.Fenced = *r.OfLen(len(s)).(*buf.Fenced)
	for i := range s {
		r.SetElem(i, s[i])
	}
	return r
}

// MarshalJSON marshalls the JSON for a Password
func (r *Password) MarshalJSON() ([]byte, error) {
	if nil == r {
		r = new(Password)
		r.Fenced = *buf.NewFenced()
		r.Status().Set("FromString() nil receiver")
	} else {
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  string(*(r.Fenced.Buf()).(*[]byte)),
		Coding: r.Coding().Get(),
		Error:  (*r.Status().Error()).Error(),
	})
}

// Stringer implementation

// String returns value encoded according to the current coding mode
func (r *Password) String() string {
	if nil == r {
		return "<nil receiver>"
	}
	return r.Fenced.String()
}
