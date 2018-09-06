// Package password is a LockedBuffer with string conversion functions, which are not recommended to be used. However the distinct type adds other safety benefits for the programmer and can be treated exactly the same as the LockedBuffer otherwise.
package password

import (
	"encoding/json"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
)

// Password is just a LockedBuffer but recommended for use in storing password fields in other structures to signify what kind of data it is.
type Password struct {
	*LockedBuffer
}

// guards against nil pointer receivers
func nilError(s string) *Password {
	r := new(Password)
	r.LockedBuffer = new(LockedBuffer)
	r.SetError(s + " nil receiver")
	return r
}

// NewPassword creates a new Password
func NewPassword(r ...*Password) *Password {
	if len(r) == 0 {
		r = append(r, new(Password))
		r[0].LockedBuffer = NewLockedBuffer()
	}
	if r[0] == nil {
		r[0] = new(Password)
		r[0].LockedBuffer = NewLockedBuffer()
	}
	r[0].SetCoding("string")
	return r[0]
}

// FromString loads the Lockedbuffer with the bytes of a string. The string is immutable so it is not removed from memory except automatically.
func (r *Password) FromString(s string) *Password {
	if r == nil {
		r = NewPassword(r)
	}
	r.LockedBuffer = r.New(len(s)).(*LockedBuffer)
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
	if r.LockedBuffer == nil {
		r.LockedBuffer = NewLockedBuffer()
	}
	return json.Marshal(&struct {
		Value  string
		Coding string
		Error  string
	}{
		Value:  string(*(r.LockedBuffer.Buf()).(*[]byte)),
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
	return r.LockedBuffer.String()
}
