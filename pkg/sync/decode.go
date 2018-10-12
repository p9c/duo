package sync

import "github.com/parallelcointeam/duo/pkg/core"

// DecodeKV decodes the record types used in the sync indices
func DecodeKV(k, v []byte) interface{} {
	switch k[0] {
	case 1:
		// height to hash and disk location
		hB := k[1:]
		hL := len(hB)
		padL := 4 - hL
		hB = append(hB, make([]byte, padL)...)
		var height uint32
		core.BytesToInt(&height, &hB)

		lB := append(v[:3], []byte{0}...)
		var length uint32
		core.BytesToInt(&length, &lB)

		sL := v[3]
		sB := v[4 : 4+sL]
		sPadL := 8 - sL
		sB = append(sB, make([]byte, sPadL)...)
		var start uint64
		core.BytesToInt(&start, &sB)

		h := v[4+sL:]

		return []interface{}{
			k[0],
			&Block{
				Height: height,
				Length: length,
				Start:  start,
				Hash:   h,
			},
		}
	case 2:
		// hash reverse lookup
		// key is 8 bytes of HHash64
		hhash := k[1:]

		// value is trailing zero byte truncated 32 bit block height
		hL := len(v)
		hPad := 4 - hL
		var height uint32
		hB := append(make([]byte, hPad), k[1:]...)
		core.BytesToInt(&height, &hB)
		return []interface{}{
			k[0],
			&Hash{
				HHash:  hhash,
				Height: height,
			},
		}
	case 4:
		// address record
		var addr Address
		addr.HHash = k
		numRefs := len(v) / 6
		for i := 0; i < numRefs; i++ {
			h := v[6*i : 6*i+4]
			t := v[6*i+4 : 6*i+4+2]
			var height uint32
			var txnum uint16
			core.BytesToInt(&height, &h)
			core.BytesToInt(&txnum, &t)
			addr.Locations = append(addr.Locations,
				Location{
					Height: height,
					TxNum:  txnum,
				})
		}
		return []interface{}{
			k[0],
			addr,
		}
	}
	return nil
}
