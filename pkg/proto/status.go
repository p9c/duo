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

// SetStatus is a
func (r *State) SetStatus(s string) Status {
	if r == nil {
		r = NewStatus()
		r.err = errors.New(er.NilRec)
	} else {
		r.err = errors.New(s)
	}
	r.Status = r.err.Error()
	return r
}

// SetStatusIf is a
func (r *State) SetStatusIf(err error) Status {
	switch {
	case r == nil:
		r = NewStatus()
		r.err = errors.New(er.NilRec)
	case err != nil:
		r.err = err
	}
	r.Status = r.err.Error()
	return r
}

// UnsetStatus is a
func (r *State) UnsetStatus() Status {
	if r == nil {
		r = NewStatus()
		r.err = errors.New(er.NilRec)
		r.Status = ""
	} else {
		r.err = nil
		r.Status = ""
	}
	return r
}

// OK returns true if there is no error
func (r *State) OK() bool {
	if r == nil {
		r = NewStatus()
		r.SetStatus(er.NilRec)
		return false
	}
	return r.err == nil
}

// Error implements the error interface
func (r *State) Error() string {
	switch {
	case r == nil:
		r = NewStatus()
		r.SetStatus(er.NilRec)
	case r.err != nil:
		return r.err.Error()
	}
	return ""
}
