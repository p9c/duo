package wallet

import (
	"bytes"
	"encoding/binary"
	"time"

	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/ec"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/serialize"
)

// KVToVars reads a key/value pair from the wallet storage format
func (db *DB) KVToVars(k, v []byte) (result interface{}) {
	id, keyRem := ser.GetPreLen(k)
	switch string(id) {
	case "name":
		addr, _ := ser.GetPreLen(keyRem)
		name, _ := ser.GetPreLen(v)
		return []interface{}{id, string(addr), string(name)}
	case "tx":
		hashB, _ := ser.GetPreLen(keyRem)
		hashVal := Uint.Zero256().FromBytes(hashB)
		tx := v
		return []interface{}{id, &hashVal, tx}
	case "acentry":
		acc, numB := ser.GetPreLen(keyRem)
		var num uint64
		buf := bytes.NewBuffer([]byte(string(numB)))
		binary.Read(buf, binary.LittleEndian, &num)
		return []interface{}{id, string(acc), num}
	case "key", "wkey":
		pubB, _ := ser.GetPreLen(keyRem)
		if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return nil
		} else if string(id) == "key" {
			var priv key.Priv
			priv.Set(v)
			return []interface{}{id, pub, &priv}
		} else {
			priv, keyRem := ser.GetPreLen(v)
			created, _ := binary.Varint(keyRem[:8])
			expires, _ := binary.Varint(keyRem[8:16])
			commentLen := keyRem[17]
			comment := string(keyRem[17:commentLen])
			var wkey Key
			wkey.PrivKey.Set(priv)
			wkey.TimeCreated = created
			wkey.TimeExpires = expires
			wkey.Comment = comment
			return []interface{}{id, pub, &wkey}
		}
	case "mkey":
		var mkeyID int64
		mkeyIDB := bytes.NewBuffer(keyRem[:8])
		binary.Read(mkeyIDB, binary.LittleEndian, &mkeyID)
		mkey, keyRem := ser.GetPreLen(v)
		salt, keyRem := ser.GetPreLen(keyRem)
		var method, iterations uint32
		methodB := bytes.NewBuffer(keyRem[:4])
		binary.Read(methodB, binary.LittleEndian, &method)
		iterationsB := bytes.NewBuffer(keyRem[4:8])
		binary.Read(iterationsB, binary.LittleEndian, &iterations)
		other := keyRem[8:]
		masterkey := crypto.MasterKey{
			EncryptedKey:              mkey,
			Salt:                      salt,
			DerivationMethod:          method,
			DeriveIterations:          iterations,
			OtherDerivationParameters: other,
		}
		return []interface{}{id, mkeyID, masterkey}
	case "ckey":
		pubB, _ := ser.GetPreLen(keyRem)
		priv, _ := ser.GetPreLen(v)
		if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return err
		} else {
			return []interface{}{id, pub, priv}
		}
	case "keymeta":
		keyB, _ := ser.GetPreLen(keyRem)
		if pub, err := ec.ParsePubKey(keyB, ec.S256()); err != nil {
			return err
		} else {
			buf := bytes.NewBuffer([]byte(string(v))[:4])
			var version uint32
			binary.Read(buf, binary.LittleEndian, &version)
			buf = bytes.NewBuffer([]byte(string(v))[4:])
			var ct int64
			binary.Read(buf, binary.LittleEndian, &ct)
			createtime := time.Unix(ct, 0)
			return []interface{}{id, pub, version, createtime}
		}
	case "defaultkey":
		pubB, _ := ser.GetPreLen(v)
		if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return nil
		} else {
			return []interface{}{id, pub}
		}
	case "pool":
		index, _ := binary.Uvarint(keyRem)
		versionB := bytes.NewBuffer([]byte(string(v))[:4])
		var version uint32
		binary.Read(versionB, binary.LittleEndian, &version)
		ptimeB := bytes.NewBuffer([]byte(string(v))[4:12])
		var ptime int64
		binary.Read(ptimeB, binary.LittleEndian, &ptime)
		pubB := []byte(string(v))[13:]
		if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return nil
		} else {
			return []interface{}{id, index, version, ptime, pub}
		}
	case "version":
		buf := bytes.NewBuffer([]byte(string(v))[:4])
		var version uint32
		binary.Read(buf, binary.LittleEndian, &version)
		return []interface{}{id, version}
	case "cscript":
		var hashID Uint.U160
		hID, _ := ser.GetPreLen(keyRem)
		hashID.FromBytes(hID)
		script, _ := ser.GetPreLen(v)
		return []interface{}{id, hashID, script}
	case "orderposnext":
		buf := bytes.NewBuffer([]byte(string(v)))
		var orderposnext int64
		binary.Read(buf, binary.LittleEndian, &orderposnext)
		return []interface{}{id, orderposnext}
	case "account":
		// just guessing here it is a split key-only record
		account, _ := ser.GetPreLen(keyRem)
		return []interface{}{id, string(account)}
	case "setting":
		name, _ := ser.GetPreLen(keyRem)
		return []interface{}{id, string(name), v}
	case "bestblock":
		return []interface{}{id, v}
	case "minversion":
		buf := bytes.NewBuffer([]byte(string(v)))
		var minversion uint32
		binary.Read(buf, binary.LittleEndian, &minversion)
		return []interface{}{id, minversion}
	}
	return nil
}

// VarsToKV converts a of set variables to a key/value byte slice pair. Can be used to construct the key half only for search and delete operations, by only passing the fields that go into the key.
func (db *DB) VarsToKV(vars interface{}) (result *[2][]byte) {
	V := vars.([]interface{})
	vType := V[0].(string)
	result[0] = append([]byte{byte(len(vType))}, []byte(vType)...)
	switch vType {
	case "name":
		result[0] = append(result[0], []byte(V[1].(string))...)
		if len(V) > 2 {
			result[1] = append([]byte{byte(len(V[2].(string)))}, []byte(V[2].(string))...)
		}
	case "tx":
		txhash := V[1].(Uint.U256)
		txhashB := txhash.ToBytes()
		result[0] = append(result[0], append([]byte{byte(len(txhashB))}, txhashB...)...)
		if len(V) > 2 {
			result[1] = append([]byte{}, V[2].([]byte)...)
		}
	case "acentry":
		result[0] = append(result[0], []byte(V[1].(string))...)
		if len(V) > 2 {
			result[1] = make([]byte, 4)
			binary.LittleEndian.PutUint64(result[1], V[2].(uint64))
		}
	case "key":
		pub := V[1].(ec.PublicKey)
		pubB := pub.SerializeCompressed()
		result[0] = append(result[0], append([]byte{byte(len(pubB))}, pubB...)...)
		if len(V) > 2 {
			priv := V[2].(key.Priv)
			result[1] = append([]byte{byte(len(priv.Get()))}, priv.Get()...)
		}
	case "wkey":
		pub := V[1].(ec.PublicKey)
		pubB := pub.SerializeCompressed()
		result[0] = append(result[0], append([]byte{byte(len(pubB))}, pubB...)...)
		if len(V) > 2 {
			key := V[2].(Key)
			privb := key.PrivKey.Get()
			privB := append([]byte{byte(len(privb))}, privb...)
			createdB := make([]byte, 8)
			binary.LittleEndian.PutUint64(createdB, uint64(V[3].(int64)))
			expiresB := make([]byte, 8)
			binary.LittleEndian.PutUint64(expiresB, uint64(V[4].(int64)))
			commentb := []byte(V[5].(string))
			commentB := append([]byte{byte(len(commentb))}, commentb...)
			result[1] = append(privB, createdB...)
			result[1] = append(result[1], expiresB...)
			result[1] = append(result[1], commentB...)
		}
	case "mkey":
		id := V[1].(int64)
		idB := bytes.NewBuffer(make([]byte, 8))
		binary.Write(idB, binary.LittleEndian, id)
		result[0] = append(result[0], append([]byte{byte(len(idB.Bytes()))}, idB.Bytes()...)...)
		if len(V) > 2 {
			mk := V[2].(crypto.MasterKey)
			result[1] = append([]byte{byte(len(mk.EncryptedKey))}, mk.EncryptedKey...)
			result[1] = append(result[1], append([]byte{byte(len(mk.Salt))}, mk.Salt...)...)
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
		return nil
	}
	return
}
