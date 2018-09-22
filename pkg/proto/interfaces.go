package proto

// Buffer is a generic interface for []byte buffers
type Buffer interface {
	Bytes() (out *[]byte)
	Copy(in *[]byte) Buffer
	Zero() Buffer
	Free() Buffer
	IsEqual(*[]byte) bool
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

// H160 is a 20 byte hash created with hash160.Sum that is used as a key for scripts and addresses
type H160 interface {
	GetID() ID
}

// ID is used as the key for searching for public keys (addresses also), scripts, transactions and blocks, generated using the hash160 function, which is a sha256 followed by ripemd160.
type ID string
