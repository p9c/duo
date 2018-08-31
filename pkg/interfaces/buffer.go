// Package interfaces
package interfaces

// Buffer is a generic interface for structs that contain *[]byte
type Buffer interface {
	Buf() *[]byte
	Copy(Buffer) Buffer
	Delete()
	Elem(int) byte
	Error() string
	IsSet() bool
	IsUTF8() bool
	Len() int
	Link(Buffer) Buffer
	Load(*[]byte) Buffer
	Move(Buffer) Buffer
	New(int) Buffer
	Null() Buffer
	Rand(int) Buffer
	SetBinary() Buffer
	SetElem(int, byte) Buffer
	SetError(string) Buffer
	SetUTF8() Buffer
	String() string
	Unset() Buffer
	UnsetError() Buffer
}
