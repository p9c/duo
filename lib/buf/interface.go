package buf

// Buf provides the basic services of storing, outputting, freeing and wiping buffers
type Buf interface {
	Data() interface{}
	Copy(interface{}) Buf
	Free() Buf
	Null() Buf
}
