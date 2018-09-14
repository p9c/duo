package proto

// Buffer is a generic interface for []byte buffers
type Buffer interface {
	Bytes() *[]byte
	Copy(*[]byte) Buffer
	Zero() Buffer
	Free() Buffer
}

// Coder is an interface for encoding raw bytes in various base number formats
type Coder interface {
	GetCoding() string
	SetCoding(string) Coder
	ListCodings() []string
}

// Streamer is an interface for serialising data
type Streamer interface {
	Freeze() *[]byte
	Thaw(*[]byte) interface{}
}

// Status keeps track of errors on an ongoing basis and hooks into the logger which fills with snapshots of data state for debugging
type Status interface {
	SetStatus(string) Status
	SetStatusIf(error) Status
	UnsetStatus() Status
}

// Array is an interface to access elements of a buffer
type Array interface {
	SetElem(int, interface{}) Array
	GetElem(int) interface{}
	Len() int
}
