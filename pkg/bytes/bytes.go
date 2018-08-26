// Bytes is a wrapper around the native byte slice  that automatically handles purging discarded data and enables copy, link and move functions on the data contained inside the structure.
package b

// To use it, simply new(Bytes) to get pointer to a empty new structure, and then after that you can call the methods of the interface.
type Bytes struct {
	val *[]byte
	set bool
	err error
}

type bytes interface {
	Len() int
	Null() *Bytes
	New(int) *Bytes
	Buf() []byte
	Assign(*[]byte) *Bytes
	Copy(*Bytes) *Bytes
	Link(*Bytes) *Bytes
	Move(*Bytes) *Bytes
}

// Len returns the length of the *[]byte if it has a value assigned.
func (r *Bytes) Len() int {
	if r.set {
		return len(*r.val)
	}
	return 0
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() *Bytes {
	if r.set {
		rr := *r.val
		for i := range rr {
			rr[i] = 0
		}
	}
	r.val = nil
	r.set = false
	r.err = nil
	return r
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size.
func (r *Bytes) New(size int) *Bytes {
	r.Null()
	b := make([]byte, size)
	r.Assign(&b)
	return r
}

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() *[]byte {
	return r.val
}

// Assign overwrites the pointer stored in a Bytes with a given byte slice, marks it set, and nulls the error value.
func (r *Bytes) Assign(bytes *[]byte) *Bytes {
	r.val = bytes
	r.set = true
	r.err = nil
	return r
}

// Copy duplicates the data from the *[]byte provided and zeroes and replaces its contents, clearing the error value.
func (r *Bytes) Copy(bytes *Bytes) *Bytes {
	r.Null()
	temp := make([]byte, bytes.Len())
	b := bytes.Buf()
	for i := range *b {
		temp[i] = (*b)[i]
	}
	r.Assign(&temp)
	return r
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(bytes *Bytes) *Bytes {
	r.Null()
	r.Assign(bytes.val)
	return r
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(bytes *Bytes) *Bytes {
	r.Assign(bytes.val)
	bytes.Null()
	return r
}
