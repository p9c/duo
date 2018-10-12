package sync

import "github.com/parallelcointeam/duo/pkg/core"

// EncodeKV decodes the record types used in the sync indices
func EncodeKV(in interface{}) (k, v []byte) {
	switch in.(type) {
	case Block:
		I := in.(Block)
		// 1 identifies a block record, trailing zeroes are removed from height value
		k = append([]byte{1}, removeTrailingZeroes(*core.IntToBytes(I.Height))...)

		// The full 32 byte block hash next has its 'difficulty' zero prefix bytes removed
		v = removeLeadingZeroes(I.Hash)

	case Hash:
		I := in.(Hash)

		// 2 identifies a hash to height index record. HHash is the 64 bit HighwayHash of the block hash, used as the search key
		k = append([]byte{2}, I.HHash...)

		// Height is stored with trailing zeroes omitted as in the Block record
		v = removeTrailingZeroes(*core.IntToBytes(I.Height))

	case Address:
		I := in.(Address)

		// We also use the Highway Hash 64 to save space for the 20 byte address field. Prefix is 4 for no particular reason.
		k = append([]byte{4}, I.HHash...)

		for i := range I.Locations {
			height := *core.IntToBytes(I.Locations[i].Height)
			txnum := *core.IntToBytes(I.Locations[i].TxNum)
			v = append(append(v, height...), txnum...)
		}
	case BalanceCache:
		I := in.(BalanceCache)

		k = append([]byte{8}, I.HHash...)

		b := removeTrailingZeroes(*core.IntToBytes(I.Balance))
		k = append([]byte{byte(len(b))}, b...)
		h := removeTrailingZeroes(*core.IntToBytes(I.Height))
		k = append(k, h...)
	}
	return
}
