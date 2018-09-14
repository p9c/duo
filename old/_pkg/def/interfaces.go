// Package def contains interfaces defining the basic types in duo and other global information
package def

// Buffer is a generic interface for structs that contain *[]byte elements
type Buffer interface {
	Buf() interface{}        // returns the pointer to the slice it represents
	Copy(interface{}) Buffer // duplicates a passed buffer
	Free() Buffer            // dereferences the buffer
	Len() int                // returns the byte length of the buffer
	Link(Buffer) Buffer      // copies the pointer of another Buffer
	OfLen(int) Buffer        // creates a new empty buffer of a given size
	Null() Buffer            // zeroes out the buffer
	Rand(...int) Buffer      // loads the buffer with random bytes
	Coding() Coding          // Defines  output encoding for the stringer
	Status() Status          // Embeds error handling
	Array() Array            // Allows the addressing of a buffer as a list of uniform sized subunits
	String() (S string)      // renders a string from the buffer
}

// Coding is an interface for data types that have multiple string output formats
type Coding interface {
	Get() string       // returns the current coding type
	Set(string) Coding // sets the coding type
	List() []string    // returns an array of strings representing available coding types
}

// Status is an extension of the Error interface that also adds getting and setting
type Status interface {
	Set(...interface{}) interface{} // sets the string that is returned by Error()
	Unset() interface{}             // clears error to represent the nominal state
	Error() string                  // returns the error field of a Status
}

// Array is an interface for working with data that can be split into multiple identical typed subunits.
type Array interface {
	Elem(int) interface{}                 // returns a buffer containing the numbered array element
	Len() int                             // returns the length of the array
	SetElem(int, interface{}) interface{} // loads the buffer at the given index
}
