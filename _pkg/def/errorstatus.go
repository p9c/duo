package def

import (
	"errors"
	"fmt"
)

// ErrorStatus is a status getter/setter for controlling the error string
type ErrorStatus struct {
	Err error
}

// Set the status string/error
func (r *ErrorStatus) Set(s ...interface{}) interface{} {
	if r == nil {
		r = new(ErrorStatus)
	}
	r.Err = errors.New(fmt.Sprint(s...))
	Dbg.Print(s...)
	return r
}

// Unset the status string/error
func (r *ErrorStatus) Unset() interface{} {
	if r == nil {
		r = new(ErrorStatus)
	}
	r.Err = nil
	return r
}

// Error returns the error field of a Status
func (r *ErrorStatus) Error() string {
	if r == nil {
		r = new(ErrorStatus)
	}
	if r.Err != nil {
		return r.Err.Error()
	}
	return ""
}
