package sync

import (
	"encoding/binary"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
)

func removeTrailingZeroes(in []byte) []byte {
	length := 0
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] != 0 {
			length = i
			break
		}
	}
	return in[:length+1]
}

func removeLeadingZeroes(in []byte) []byte {
	nonzerostart := 0
	for i := range in {
		if in[i] != 0 {
			nonzerostart = i
			break
		}
	}
	return in[nonzerostart:]
}

func (r *Node) getLatest() (h uint32) {
	var latestB []byte

	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("latest"))
		if err == nil {
			latestB, err = item.Value()
		}
		return err
	}))
	if latestB != nil {
		heightB := latestB[:4]
		core.BytesToInt(&h, &heightB)
	}
	r.UnsetStatus()
	return
}

func decodeAddressRecord(addr []byte) (out []Location, length uint32) {
	cursor, length := 0, 0
	for cursor < len(addr) {
		var height, txnum uint64
		height, step := binary.Uvarint(addr[cursor:])
		cursor += step
		txnum, step = binary.Uvarint(addr[cursor:])
		cursor += step
		out = append(out, Location{Height: uint32(height), TxNum: uint16(txnum)})
		cursor++
	}
	return
}

func encodeAddressRecord(existing []byte, loc Location) (out []byte) {
	ex, _ := decodeAddressRecord(existing)
	ex = append(ex, loc)
	for i := range ex {
		if i == 0 {
			h := make([]byte, 5)
			l := binary.PutUvarint(h, uint64(ex[0].Height))
			h = h[:l]
			n := make([]byte, 5)
			l = binary.PutUvarint(n, uint64(ex[0].TxNum))
			n = n[:l]
			out = append(h, n...)
		} else {
			nh := ex[i].Height - ex[i-1].Height
			h := make([]byte, 5)
			l := binary.PutUvarint(h, uint64(nh))
			h = h[:l]
			n := make([]byte, 5)
			l = binary.PutUvarint(n, uint64(ex[0].TxNum))
			n = n[:l]
			out = append(out, append(h, n...)...)
		}
	}

	return
}
