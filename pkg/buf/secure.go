package buf

import (
	"github.com/parallelcointeam/duo/pkg/proto"
)

// Secure is a simple single byte
type Secure proto.Secure

// Bytes is a
func (r *Secure) Bytes(b *[]byte) proto.Buffer {
	panic("not implemented")
}

// Copy is a
func (r *Secure) Copy(in *[]byte) proto.Buffer {
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
func (r *Secure) GetCoding(out string) proto.Coder {
	panic("not implemented")
}

// SetCoding is a
func (r *Secure) SetCoding(string) proto.Coder {
	panic("not implemented")
}

// ListCodings is a
func (r *Secure) ListCodings(out *[]string) proto.Coder {
	panic("not implemented")
}

// Freeze is a
func (r *Secure) Freeze(out *[]byte) proto.Streamer {
	panic("not implemented")
}

// Thaw is a
func (r *Secure) Thaw(in *[]byte) proto.Streamer {
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
func (r *Secure) SetElem(index int, in interface{}) proto.Array {
	panic("not implemented")
}

// GetElem is a
func (r *Secure) GetElem(index int, out interface{}) proto.Array {
	panic("not implemented")
}

// Len is a
func (r *Secure) Len(length *int) proto.Array {
	panic("not implemented")
}
