// Package interfaces contains interfaces defining the basic types in duo
package interfaces

// Nil is a helper interface to stop nil receiver panics
type Nil interface {
	ifnil() interface{}
}

// Buffer is a generic interface for structs that contain *[]byte slice elements
type Buffer interface {
	// Returns the pointer to the slice of bytes the buffer represents
	Buf() *[]byte
	// Duplicates a passed buffer, nulling and freeing the old data
	Copy(Buffer) Buffer
	// Dereferences the buffer allowing it to be discarded
	Free() Buffer
	// Creates a link between the underlying byte slice in one buffer to another
	Link(Buffer) Buffer
	// Puts content into the buffer, wiping the source slice
	Load(*[]byte) Buffer
	// Copies the pointer to the underlying bytes and removes it from the source
	Move(Buffer) Buffer
	// Creates a new Buffer with an empty buffer of a given length
	New(int) Buffer
	// Zeroes out the buffer
	Null() Buffer
	// Loads the buffer with some number of random bytes
	Rand(int) Buffer
	// Range runs a loop and calls a function that can be a closure
	ForEach(func(int)) Buffer
	// Returns the byte length of the buffer
	Size() int
	// Returns the current coding type
	Coding() string
	// Sets the coding type
	SetCoding(string) Buffer
	// Returns an array of strings representing available coding types
	Codes() []string
	// Returns a string representing the content according to the coding mode set
	String() string
	// Return the string stored in the error. If implements status also implements errors
	Error() string
	// Sets the string that is returned by Error()
	SetError(string) Buffer
	// Clears error to represent the nominal state
	UnsetError() Buffer
	// Returns true if toggle is set
	IsSet() bool
	// Sets toggle to true
	Set() Buffer
	// Sets toggle to false
	Unset() Buffer
	Array
}

// Array is an abstract type for a memory ordered, dense list store that can be walked by increment/decrement of indices
type Array interface {
	// Returns a buffer containing the numbered array element
	Elem(int) Buffer
	// Returns the length of the array
	Len() int
	// Returns the current maximum capacity of the array (reserve allocation)
	Cap() int
	// Zeroes out all of the buffers in the array
	Purge() Array
	// Loads the buffer at the given index
	SetElem(int, Buffer) Array
}
