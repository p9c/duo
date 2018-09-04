// Package interfaces contains interfaces defining the basic types in duo
package interfaces

// Nil is a helper interface to stop nil receiver panics. This interface has an abstract function that works with any type of concrete type, to use it simply implement the function and it can be directly used.
type Nil interface {
	// Nil returns true if the receiver is nil
	Nil() bool
}

func doif(b bool, f ...func()) {
	switch {
	case len(f) == 1:
		if f[0] != nil && b {
			f[0]()
		}
	case len(f) == 2:
		if f[0] != nil {
			if b {
				f[0]()
			}
			if f[1] != nil {
				if !b {
					f[1]()
				}
			}
		}

	}
}

func donil(b bool, f ...func()) {
	switch {
	case len(f) == 1:
		if b {
			if f[0] != nil {
				f[0]()
			}
		}
	case len(f) == 2:
		if b {
			if f[0] != nil {
				f[0]()
			}
			if f[0] != nil {
				f[1]()
			}
		} else {
			if f[0] != nil {
				f[1]()
			}
		}
	}
}

// DoIf runs first func if true, seecond if false, if is a pointer, runs both if the pointer is nil
func DoIf(is interface{}, do ...func()) (R interface{}) {
	switch is.(type) {
	case bool:
		b := is.(bool)
		doif(b, do...)
	case Nil, *[]byte:
		var b bool
		switch is.(type) {
		case Nil:
			b = is.(Nil).Nil()
		case *[]byte:
			b = is.(*[]byte) != nil
		}
		donil(b, do...)
	}
	return is
}

// Buffer is a generic interface for structs that contain *[]byte elements
type Buffer interface {
	Buf() interface{}        // returns the pointer to the slice it represents
	Copy(Buffer) Buffer      // duplicates a passed buffer
	Free() Buffer            // dereferences the buffer after any necessary preprocessing
	Link(interface{}) Buffer // copies the pointer of another Buffer
	Load(interface{}) Buffer // copies content into the buffer, and wipes the original
	Move(Buffer) Buffer      // copies the pointer to the underlying buffer and removes it from the source
	New(int) Buffer          // creates a new empty buffer of a given size
	Null() Buffer            // zeroes out the buffer
	Rand(...int) Buffer      // loads the buffer with some number of random bytes. Has variadic parameter to allow fixed length buffers and specify element sizes in a list
	Size() int               // returns the byte length of the buffer
	String() string          // returns a string representing the content according to the coding mode set
	Coding                   // Defines  output encoding for the stringer
	Status                   // Embeds error handling
	Toggle                   // Indicates whether initialisation has been performed
	Array                    // Allows the addressing of a buffer as a list of uniform sized subunits
}

// Coding is an interface for data types that have multiple string output formats
type Coding interface {
	Coding() string               // returns the current coding type
	SetCoding(string) interface{} // sets the coding type
	Codes() []string              // returns an array of strings representing available coding types
}

// Status is an extension of the Error interface that also adds getting and setting
type Status interface {
	SetError(string) interface{} // sets the string that is returned by Error()
	UnsetError() interface{}     // clears error to represent the nominal state
	Error() string               // returns the error string
}

// Toggle is an interface for a boolean value
type Toggle interface {
	IsSet() bool        // returns true if toggle is set
	Set() interface{}   // sets toggle to true
	Unset() interface{} // sets toggle to false
}

// Array is an interface for working with data that can be split into multiple identical typed subunits.
type Array interface {
	Cap() int                             // returns the current maximum capacity of the array (reserve allocation)
	Elem(int) interface{}                 // returns a buffer containing the numbered array element
	ForEach(func(int) bool) bool          // runs a loop over an array and calls a function that can be a closure, return value indicates escaped before full iteration
	Len() int                             // returns the length of the array
	Purge() interface{}                   // zeroes out all of the buffers in the array
	SetElem(int, interface{}) interface{} // loads the buffer at the given index
}

// ForEach calls a passed function for each element in a array or slice type, that returns a boolean to indicate success and break the loop.
func ForEach(iface interface{}, fn func(int) bool) bool {
	switch iface.(type) {
	case []string:
		for i := range iface.([]string) {
			if fn(i) {
				return true
			}
		}
	case []byte:
		for i := range iface.([]byte) {
			if fn(i) {
				return true
			}
		}
	}
	return false
}
