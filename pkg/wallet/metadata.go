package wallet

import (
	"encoding/binary"
	"errors"
)

// KeyMetadata is information about the key
type KeyMetadata struct {
	Version    int
	CreateTime int64
}

type keymetadata interface {
	FromBytes([]byte)
	ToBytes() []byte
}

// NewKeyMetadata makes a new KeyMetadata structure
func NewKeyMetadata(createTime int64) (M *KeyMetadata) {
	M.Version = CurrentVersion
	M.CreateTime = createTime
	return
}

// FromBytes takes a 12 element byte slice and copies it into the structure as stored in memory
func (m *KeyMetadata) FromBytes(b []byte) (err error) {
	lengthV := binary.Size(m.Version)
	lengthC := binary.Size(m.CreateTime)
	if len(b) != int(lengthC+lengthV) {
		err = errors.New("Bytes were wrong length")
	} else {
		version, _ := binary.Varint(b[:lengthV])
		m.Version = int(version)
		createtime, _ := binary.Varint(b[lengthV:])
		m.CreateTime = int64(createtime)
	}
	return
}

// ToBytes returns a byte slice containing the raw bytes of a Metadata structure in structure order
func (m *KeyMetadata) ToBytes() (b []byte) {
	lv, lc := binary.Size(m.Version), binary.Size(m.CreateTime)
	v, c := make([]byte, lv), make([]byte, lc)
	binary.PutVarint(v, int64(m.Version))
	binary.PutVarint(c, m.CreateTime)
	return append(v, c...)
}
