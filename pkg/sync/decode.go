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
		hB = append(make([]byte, padL), hB...)
		var height uint32
		core.BytesToInt(&height, &hB)

		lB := append([]byte{0}, v[:2]...)
		var length uint32
		core.BytesToInt(&length, &lB)

		sL := v[2]
		sB := v[3 : 3+sL]
		sPadL := 8 - sL
		sB = append(sB, make([]byte, sPadL)...)
		var start uint64
		core.BytesToInt(&start, &sB)

		h := v[3+sL:]

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
		hhash := k

		hL := len(k)
		hPad := 4 - hL
		var height uint32
		hB := append(k[1:], make([]byte, hPad)...)
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

// EncodeKV decodes the record types used in the sync indices
func EncodeKV(in interface{}) (k, v []byte) {

	return
}
