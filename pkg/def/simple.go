package def

import (
	"errors"
	"fmt"
	"github.com/jimlawless/whereami"
)

func dbg(s ...interface{}) {
	Debug(s)
}

// StringCoding is a coding type for converting data to strings
type StringCoding int

// Get returns the current coding type
func (r *StringCoding) Get() string {
	switch {
	case int(*r) > len(StringCodingTypes):
		dbg("Get() code higher than maximum")
		*r = 0
	case int(*r) < 0:
		dbg("Get() negative coding")
		*r = 0
		return StringCodingTypes[*r]
	}
	return StringCodingTypes[*r]
}

// Set sets the coding type
func (r *StringCoding) Set(s string) Coding {
	*r = 0
	for i := range StringCodingTypes {
		if s == StringCodingTypes[i] {
			*r = StringCoding(i)
			break
		}
	}
	return r
}

// List returns an array of strings representing available coding
func (r *StringCoding) List() (R []string) {
	for i := range StringCodingTypes {
		R = append(R, StringCodingTypes[i])
	}
	return
}

// ErrorStatus is a status getter/setter for controlling the error string
type ErrorStatus struct {
	err error
}

// Set the status string/error
func (r *ErrorStatus) Set(s interface{}) interface{} {
	if r == nil {
		r = new(ErrorStatus)
	}
	r.err = errors.New(whereami.WhereAmI() + fmt.Sprint(s))
	dbg(s)
	return r
}

// Unset the status string/error
func (r *ErrorStatus) Unset() interface{} {
	if r == nil {
		r = new(ErrorStatus)
	}
	r.err = nil
	return r
}

// Error returns the error field of a Status
func (r *ErrorStatus) Error() *error {
	return &r.err
}
