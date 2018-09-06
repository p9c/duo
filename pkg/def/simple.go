package def

import (
	"errors"
)

var me = "gitlab.com/parallelcoin/duo/pkg/def/const.go"

func debug(s ...string) {
	s = append([]string{me + " :"}, s...)
	Debug(s)
}

// StringCoding is a coding type for converting data to strings
type StringCoding int

// Get returns the current coding type
func (r *StringCoding) Get() string {
	switch {
	case int(*r) > len(StringCodingTypes):
		Debug("Get() code higher than maximum")
		*r = 0
	case int(*r) < 0:
		Debug("Get() negative coding")
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
func (r *ErrorStatus) Set(s string) interface{} {
	r.err = errors.New(s)
	return r
}

// Unset the status string/error
func (r *ErrorStatus) Unset() interface{} {
	r.err = nil
	return r
}