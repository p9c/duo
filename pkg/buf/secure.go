package buf

import (
	"github.com/parallelcointeam/duo/pkg/proto"
)

// Secure is a simple single byte
type Secure proto.Secure

// Bytes is a
func (r *Secure) Bytes() *[]byte {
	panic("not implemented")
}

// Copy is a
func (r *Secure) Copy(*[]byte) proto.Buffer {
	panic("not implemented")
}

// Zero is a
func (r *Secure) Zero() proto.Buffer {
	panic("not implemented")
}

// Free is a
func (r *Secure) Free() proto.Buffer {
	panic("not implemented")
}

// GetCoding is a
func (r *Secure) GetCoding() string {
	panic("not implemented")
}

// SetCoding is a
func (r *Secure) SetCoding(string) proto.Coder {
	panic("not implemented")
}

// ListCodings is a
func (r *Secure) ListCodings() []string {
	panic("not implemented")
}

// Freeze is a
func (r *Secure) Freeze() *[]byte {
	panic("not implemented")
}

// Thaw is a
func (r *Secure) Thaw(*[]byte) interface{} {
	panic("not implemented")
}

// SetStatus is a
func (r *Secure) SetStatus(string) proto.Status {
	panic("not implemented")
}

// SetStatusIf is a
func (r *Secure) SetStatusIf(error) proto.Status {
	panic("not implemented")
}

// UnsetStatus is a
func (r *Secure) UnsetStatus() proto.Status {
	panic("not implemented")
}

// SetElem is a
func (r *Secure) SetElem(int, interface{}) proto.Array {
	panic("not implemented")
}

// GetElem is a
func (r *Secure) GetElem(int) interface{} {
	panic("not implemented")
}

// Len is a
func (r *Secure) Len() int {
	panic("not implemented")
}
