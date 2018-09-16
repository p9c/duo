package proto

// Buffer is a generic interface for []byte buffers
type Buffer interface {
	Bytes() (out *[]byte)
	Copy(in *[]byte) Buffer
	Zero() Buffer
	Free() Buffer
}

// Coder is an interface for encoding raw bytes in various base number formats
type Coder interface {
	GetCoding() (out *string)
	SetCoding(in string) Coder
	ListCodings() (out *[]string)
}

// Streamer is an interface for serialising data
type Streamer interface {
	Freeze() (out *[]byte)
	Thaw(in *[]byte) Streamer
}

// Status keeps track of errors on an ongoing basis and hooks into the logger which fills with snapshots of data state for debugging
type Status interface {
	SetStatus(string) Status
	SetStatusIf(error) Status
	UnsetStatus() Status
	OK() bool
}

// Array is an interface to access elements of a buffer
type Array interface {
	SetElem(index int, in interface{}) Array
	GetElem(index int) (out interface{})
	Len() int
}
