package sync

import (
	"encoding/binary"
)

// EncodeKV decodes the record types used in the sync indices
func EncodeKV(in interface{}) (k, v []byte) {
	switch in.(type) {
	case Block:
		I := in.(Block)
		// 1 identifies a block record, trailing zeroes are removed from height value
		k = make([]byte, 5)
		l := binary.PutUvarint(k, uint64(I.Height))
		k = k[:l]
		k = append([]byte{1}, k...)

		// The full 32 byte block hash next has its 'difficulty' zero prefix bytes removed
		v = removeLeadingZeroes(I.Hash)

	case Hash:
		I := in.(Hash)

		// 2 identifies a hash to height index record. HHash is the 64 bit HighwayHash of the block hash, used as the search key
		k = append([]byte{2}, I.HHash...)

		// Height is stored with trailing zeroes omitted as in the Block record
		v = make([]byte, 5)
		l := binary.PutUvarint(v, uint64(I.Height))
		v = v[:l]

	case Address:
		I := in.(Address)

		// We also use the Highway Hash 64 to save space for the 20 byte address field. Prefix is 4 for no particular reason.
		k = append([]byte{4}, I.HHash...)

		for i := range I.Locations {
			height := make([]byte, 5)
			txnum := make([]byte, 3)
			l := binary.PutUvarint(height, uint64(I.Locations[i].Height))
			height = height[:l]
			l = binary.PutUvarint(txnum, uint64(I.Locations[i].TxNum))
			txnum = txnum[:l]
			v = append(append(v, height...), txnum...)
		}
	case BalanceCache:
		I := in.(BalanceCache)

		k = append([]byte{8}, I.HHash...)

		b := make([]byte, 9)
		l := binary.PutUvarint(b, I.Balance)
		b = b[:l]
		h := make([]byte, 5)
		l = binary.PutUvarint(h, uint64(I.Height))
		h = h[:l]
		v = append(b, h...)
	}
	return
}
