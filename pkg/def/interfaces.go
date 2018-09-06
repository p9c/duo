// Package def contains interfaces defining the basic types in duo and other global information
package def

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

// Array is an interface for working with data that can be split into multiple identical typed subunits.
type Array interface {
	Cap() int                             // returns the current maximum capacity of the array (reserve allocation)
	Elem(int) interface{}                 // returns a buffer containing the numbered array element
	Len() int                             // returns the length of the array
	Purge() interface{}                   // zeroes out all of the buffers in the array
	SetElem(int, interface{}) interface{} // loads the buffer at the given index
}
