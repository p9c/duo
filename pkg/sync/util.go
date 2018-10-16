package sync

import (
	"encoding/binary"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/golang/snappy"
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
	var dec []byte
	var err error
	dec, err = snappy.Decode(nil, addr)
	if err == nil {
		addr = dec
	}
	cursor, length := 0, 0
	for cursor < len(addr) {
		var h, n uint64
		var height uint32
		var txnum uint16
		h, step := binary.Uvarint(addr[cursor:])
		cursor += step
		n, step = binary.Uvarint(addr[cursor:])
		cursor += step
		l := len(out)
		if l > 1 {
			height = out[l-1].Height + uint32(h)
		} else {
			height, txnum = uint32(h), uint16(n)
		}
		out = append(out, Location{Height: uint32(height), TxNum: uint16(txnum)})
		length++
	}
	return
}

func encodeAddressRecord(existing []byte, loc Location) (out []byte) {
	ex, _ := decodeAddressRecord(existing)
	// fmt.Println("ex", ex)
	ex = append(ex, loc)
	var h, n []byte
	for i := range ex {
		if i == 0 {
			h = make([]byte, 5)
			l := binary.PutUvarint(h, uint64(ex[0].Height))
			h = h[:l]
			n = make([]byte, 5)
			l = binary.PutUvarint(n, uint64(ex[0].TxNum))
			n = n[:l]
			out = append(h, n...)
		} else {
			nh := ex[i].Height - ex[i-1].Height
			// fmt.Println("\n", ex[i].Height, ex[i-1].Height, nh)
			h = make([]byte, 5)
			l := binary.PutUvarint(h, uint64(nh))
			h = h[:l]
			n = make([]byte, 5)
			l = binary.PutUvarint(n, uint64(ex[i].TxNum))
			n = n[:l]
			out = append(out, append(h, n...)...)
		}
		// fmt.Println("\n", uint64(ex[i].Height), h, uint64(ex[0].TxNum), n)
	}
	enc := snappy.Encode(nil, out)
	out = enc
	if len(enc) > len(out) {
		fmt.Println("\ncompressor expanded!")
	}
	// fmt.Print(len(out)) //, hex.EncodeToString(out))
	return
}
