package wallet

import (
	"bytes"
	"encoding/binary"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/key"
)

// KVEnc converts a of set variables to a key/value byte slice pair. Can be used to construct the key half only for search and delete operations, by only passing the fields that go into the key.
func (db *DB) KVEnc(vars interface{}) (result [2][]byte) {
	V := vars.([]interface{})
	vType := V[0].(string)
	result[0] = FormatString(vType)
	switch vType {
	case "name":
		Append(&result[0], FormatString(V[1].(string)))
		if len(V) > 2 {
			result[1] = FormatString(V[2].(string))
		}
	case "tx":
		Append(&result[0], V[1].(*Uint.U256).ToBytes())
		if len(V) > 2 {
			result[1] = V[2].([]byte)
		}
	case "acentry":
		Append(&result[0], FormatString(V[1].(string)), Uint64ToBytes(V[2].(uint64)))
	case "key":
		Append(&result[0], FormatBytes(V[1].(*key.Pub).GetPub().SerializeUncompressed()))
		if len(V) > 2 {
			result[1] = FormatBytes(V[2].(*key.Priv).Get())
		}
	case "wkey":
		Append(&result[0], FormatBytes(V[1].(*key.Pub).GetPub().SerializeUncompressed()))
		if len(V) > 2 {
			Append(&result[1],
				FormatBytes(V[2].(*Key).PrivKey.Get()),
				Int64ToBytes(V[3].(int64)),
				Int64ToBytes(V[4].(int64)),
				FormatString(V[5].(string)))
		}
	case "mkey":
		Append(&result[0], Int64ToBytes(V[1].(int64)))
		if len(V) > 2 {
			mk := V[2].(*crypto.MasterKey)
			Append(&result[1],
				FormatBytes(mk.EncryptedKey),
				FormatBytes(mk.Salt),
				Uint32ToBytes(mk.DerivationMethod),
				Uint32ToBytes(mk.DeriveIterations),
				mk.OtherDerivationParameters)
		}
	case "ckey":
		Append(&result[0], FormatBytes(V[1].(*key.Pub).GetPub().SerializeUncompressed()))
		if len(V) > 2 {
			result[1] = FormatBytes(V[2].([]byte))
		}
	case "keymeta":
		Append(&result[0], FormatBytes(V[1].(*key.Pub).GetPub().SerializeUncompressed()))
		if len(V) > 2 {
			Append(&result[1], Uint32ToBytes(V[2].(uint32)), Int64ToBytes(V[3].(int64)))
		}
	case "defaultkey":
		if len(V) > 2 {
			Append(&result[1], V[1].(*key.Pub).GetPub().SerializeUncompressed())
		}
	case "pool":
		Append(&result[0], Uint64ToBytes(V[1].(uint64)))
		if len(V) > 2 {
			Append(&result[1],
				Uint32ToBytes(V[2].(uint32)),
				Int64ToBytes(V[3].(int64)),
				V[4].(*key.Pub).GetPub().SerializeUncompressed())
		}
	case "version":
		if len(V) > 2 {
			result[1] = Uint32ToBytes(V[1].(uint32))
		}
	case "cscript":
		hashID := V[1].(Uint.U160)
		hashIDB := hashID.ToBytes()
		result[0] = append(result[0], append([]byte{byte(len(hashIDB))}, hashIDB...)...)
		if len(V) > 2 {
			script := V[2].([]byte)
			result[1] = append([]byte{byte(len(script))}, script...)
		}
	case "orderposnext":
		if len(V) > 2 {
			opn := bytes.NewBuffer(make([]byte, 8))
			binary.Write(opn, binary.LittleEndian, V[1].(int64))
			result[1] = opn.Bytes()
		}
	case "account":
		acct := V[1].(string)
		result[0] = append(result[0], append([]byte{byte(len(acct))}, []byte(acct)...)...)
	case "setting":
		name := V[1].(string)
		result[0] = append(result[0], append([]byte{byte(len(name))}, []byte(name)...)...)
		if len(V) > 2 {
			result[1] = V[2].([]byte)
		}
	case "bestblock":
		if len(V) > 2 {
			result[1] = V[1].([]byte)
		}
	case "minversion":
		if len(V) > 2 {
			minversionB := bytes.NewBuffer(make([]byte, 4))
			binary.Write(minversionB, binary.LittleEndian, V[1].(uint32))
			result[1] = minversionB.Bytes()
		}
	default:
		return
	}
	return
}
