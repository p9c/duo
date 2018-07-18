package wallet

import (
	"bytes"
	"encoding/binary"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/ec"
	"gitlab.com/parallelcoin/duo/pkg/key"
)

// KVEnc converts a of set variables to a key/value byte slice pair. Can be used to construct the key half only for search and delete operations, by only passing the fields that go into the key.
func (db *DB) KVEnc(vars interface{}) (result [2][]byte) {
	V := vars.([]interface{})
	vType := V[0].(string)
	result[0] = PreLenString(vType)
	switch vType {
	case "name":
		result[0] = append(result[0], PreLenString(V[1].(string))...)
		if len(V) > 2 {
			result[1] = PreLenString(V[2].(string))
		}
	case "tx":
		result[0] = append(result[0], V[1].(*Uint.U256).ToBytes()...)
		if len(V) > 2 {
			result[1] = append([]byte{}, V[2].([]byte)...)
		}
	case "acentry":
		result[0] = append(result[0], []byte(PreLenString(V[1].(string)))...)
		amount := make([]byte, 4)
		binary.LittleEndian.PutUint64(amount, V[2].(uint64))
		result[0] = append(result[0], amount...)
	case "key":
		pub := V[1].(*key.Pub).GetPub().SerializeUncompressed()
		result[0] = append(result[0], PreLenBytes(pub)...)
		if len(V) > 2 {
			priv := V[2].(key.Priv)
			result[1] = PreLenBytes(priv.Get())
		}
	case "wkey":
		pubB := V[1].(*key.Pub).GetPub().SerializeUncompressed()
		result[0] = append(result[0], PreLenBytes(pubB)...)
		if len(V) > 2 {
			privB := PreLenBytes(V[2].(*Key).PrivKey.Get())
			createdB := Int64ToBytes(V[3].(int64))
			expiresB := Int64ToBytes(V[4].(int64))
			commentB := PreLenBytes([]byte(V[5].(string)))
			result[1] = append(privB, append(createdB, append(expiresB, commentB...)...)...)
		}
	case "mkey":
		id := V[1].(int64)
		idB := bytes.NewBuffer(make([]byte, 8))
		binary.Write(idB, binary.LittleEndian, id)
		result[0] = append(result[0], PreLenBytes(idB.Bytes())...)
		if len(V) > 2 {
			mk := V[2].(crypto.MasterKey)
			result[1] = append(PreLenBytes(mk.EncryptedKey), PreLenBytes(mk.Salt)...)
			methodB := bytes.NewBuffer(make([]byte, 4))
			binary.Write(methodB, binary.LittleEndian, mk.DerivationMethod)
			result[1] = append(result[1], methodB.Bytes()...)
			iterationsB := bytes.NewBuffer(make([]byte, 4))
			binary.Write(iterationsB, binary.LittleEndian, mk.DeriveIterations)
			result[1] = append(result[1], iterationsB.Bytes()...)
			result[1] = append(result[1], mk.OtherDerivationParameters...)
		}
	case "ckey":
		pub := V[1].(*ec.PublicKey)
		pubB := pub.SerializeCompressed()
		result[0] = append(result[0], append([]byte{byte(len(pubB))}, pubB...)...)
		if len(V) > 2 {
			priv := V[2].([]byte)
			result[1] = append([]byte{byte(len(priv))}, priv...)
		}
	case "keymeta":
		pub := V[1].(*ec.PublicKey)
		pubB := pub.SerializeCompressed()
		result[0] = append(result[0], append([]byte{byte(len(pubB))}, pubB...)...)
		if len(V) > 2 {
			versionB := bytes.NewBuffer(make([]byte, 4))
			binary.Write(versionB, binary.LittleEndian, V[2].(uint32))
			createtimeB := bytes.NewBuffer(make([]byte, 8))
			binary.Write(createtimeB, binary.LittleEndian, V[3].(int64))
			result[1] = append(versionB.Bytes(), createtimeB.Bytes()...)
		}
	case "defaultkey":
		if len(V) > 2 {
			pub := V[1].(*ec.PublicKey)
			pubB := pub.SerializeCompressed()
			result[1] = append([]byte{byte(len(pubB))}, pubB...)
		}
	case "pool":
		index := V[1].(uint64)
		indexB := bytes.NewBuffer(make([]byte, 8))
		binary.Write(indexB, binary.LittleEndian, index)
		result[0] = append(result[0], indexB.Bytes()...)
		if len(V) > 2 {
			versionB := bytes.NewBuffer(make([]byte, 4))
			binary.Write(versionB, binary.LittleEndian, V[2].(uint32))
			ptime := bytes.NewBuffer(make([]byte, 8))
			binary.Write(ptime, binary.LittleEndian, V[3].(int64))
			pubB := append([]byte{byte(len(V[4].([]byte)))}, V[4].([]byte)...)
			result[1] = append(versionB.Bytes(), ptime.Bytes()...)
			result[1] = append(result[1], pubB...)
		}
	case "version":
		if len(V) > 2 {
			versionB := bytes.NewBuffer(make([]byte, 4))
			binary.Write(versionB, binary.LittleEndian, V[1].(uint32))
			result[1] = versionB.Bytes()
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
