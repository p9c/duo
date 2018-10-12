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

		h := append(make([]byte, 32-len(v)), v...)

		return []interface{}{
			k[0],
			&Block{
				Height: height,
				Hash:   h,
			},
		}
	case 2:
		// hash reverse lookup
		// key is 8 bytes of HHash64
		hhash := k[1:]

		// The value only stores the non-zero trailing bytes
		hL := len(v)
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
		addr.HHash = k[1:]
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
	case 8:
		// cached balance query results
		var bal BalanceCache
		bal.HHash = k[:1]

		bL := v[0]
		b := v[1 : bL+1]
		balance := append(b, make([]byte, 8-bL)...)
		h := v[bL+1:]
		height := append(h, make([]byte, 4-len(h))...)
		core.BytesToInt(&bal.Balance, &balance)
		core.BytesToInt(&bal.Height, &height)
		return []interface{}{
			k[0],
			bal,
		}
	}
	return nil
}
