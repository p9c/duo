// Package def contains interfaces defining the basic types in duo and other global information
package def

// Buffer is a generic interface for structs that contain *[]byte elements
type Buffer interface {
	Buf() interface{}        // returns the pointer to the slice it represents
	Copy(Buffer) Buffer      // duplicates a passed buffer
	Free() Buffer            // dereferences the buffer after any necessary preprocessing
	Link(interface{}) Buffer // copies the pointer of another Buffer
	Load(interface{}) Buffer // copies content into the buffer, and wipes the original
	OfSize(int) Buffer       // creates a new empty buffer of a given size
	Null() Buffer            // zeroes out the buffer
	Rand(...int) Buffer      // loads the buffer with some number of random buf. Has variadic parameter to allow fixed length buffers and specify element sizes in a list
	Size() int               // returns the byte length of the buffer
	Coding() Coding          // Defines  output encoding for the stringer
	Status() Status          // Embeds error handling
	Array() Array            // Allows the addressing of a buffer as a list of uniform sized subunits
}

// Coding is an interface for data types that have multiple string output formats
type Coding interface {
	Get() string       // returns the current coding type
	Set(string) Coding // sets the coding type
	List() []string    // returns an array of strings representing available coding types
}

// Status is an extension of the Error interface that also adds getting and setting
type Status interface {
	Set(string) interface{} // sets the string that is returned by Error()
	Unset() interface{}     // clears error to represent the nominal state
}

// Array is an interface for working with data that can be split into multiple identical typed subunits.
type Array interface {
	Elem(int) Buffer            // returns a buffer containing the numbered array element
	Len() int                   // returns the length of the array
	SetElem(int, Buffer) Buffer // loads the buffer at the given index
}
