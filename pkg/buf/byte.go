package buf

import (
	"github.com/parallelcointeam/duo/pkg/proto"
)

// Byte is a simple single byte
type Byte proto.Byte

// Bytes is a
func (r *Byte) Bytes() *[]byte {
	panic("not implemented")
}

// Copy is a
func (r *Byte) Copy(*[]byte) proto.Buffer {
	panic("not implemented")
}

// Zero is a
func (r *Byte) Zero() proto.Buffer {
	panic("not implemented")
}

// Free is a
func (r *Byte) Free() proto.Buffer {
	panic("not implemented")
}

// GetCoding is a
func (r *Byte) GetCoding() string {
	panic("not implemented")
}

// SetCoding is a
func (r *Byte) SetCoding(string) proto.Coder {
	panic("not implemented")
}

// ListCodings is a
func (r *Byte) ListCodings() []string {
	panic("not implemented")
}

// Freeze is a
func (r *Byte) Freeze() *[]byte {
	panic("not implemented")
}

// Thaw is a
func (r *Byte) Thaw(*[]byte) interface{} {
	panic("not implemented")
}

// SetStatus is a
func (r *Byte) SetStatus(string) proto.Status {
	panic("not implemented")
}

// SetStatusIf is a
func (r *Byte) SetStatusIf(error) proto.Status {
	panic("not implemented")
}

// UnsetStatus is a
func (r *Byte) UnsetStatus() proto.Status {
	panic("not implemented")
}

// SetElem is a
func (r *Byte) SetElem(int, interface{}) proto.Array {
	panic("not implemented")
}

// GetElem is a
func (r *Byte) GetElem(int) interface{} {
	panic("not implemented")
}

// Len is a
func (r *Byte) Len() int {
	panic("not implemented")
}