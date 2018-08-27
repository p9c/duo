// Package pipe is a small struct for managing the construction of libraries with functions that pass the receiver through and allow chaining, which causes issues with nil pointers amongst other things
package pipe

import (
	"reflect"
)

// Pipe is a convenience device to prevent nil pointer panics for pass-through function design pattern
type Pipe struct {
}

// NilGuard checks if a receiver is empty, receives an object to use, and a null function that allocates the struct and sets default values, that it executes if it is given
func (r *Pipe) NilGuard(R interface{}, null func(interface{}) interface{}) interface{} {
	if r == nil {
		r = new(Pipe)
	}
	if R == nil {
		R = reflect.New(reflect.TypeOf(R))
	}
	if null != nil {
		R = null(R)
	}
	return R
}
