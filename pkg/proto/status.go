package proto

import (
	"errors"
)

var er = Errors

// NewStatus creates a new Status object
func NewStatus() *State {
	r := new(State)
	return r
}

// NewIf creates a new Status object
func (r *State) NewIf() *State {
	if r == nil {
		r = NewStatus()
		r.SetStatus(er.NilRec)
	}
	return r
}

// SetStatus is a
func (r *State) SetStatus(s string) Status {
	r = r.NewIf()
	if r != nil {
		r.err = errors.New(s)
	}
	if s != "" {
		r.err = errors.New(s)
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *State) SetStatusIf(err error) Status {
	r = r.NewIf()
	if err != nil {
		r.err = err
		r.Status = r.err.Error()
	}
	return r
}

// UnsetStatus is a
func (r *State) UnsetStatus() Status {
	r = r.NewIf()
	r.Status, r.err = "", nil
	return r
}

// OK returns true if there is no error
func (r *State) OK() bool {
	if r == nil {
		r = r.NewIf()
		return false
	}
	return r.err == nil
}

// Error implements the error interface
func (r *State) Error() string {
	r = r.NewIf()
	if !r.OK() {
		return r.Status
	}
	return ""
}
