package wallet

import (
	"bytes"
	"encoding/binary"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/ec"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/serialize"
	"time"
)

// KVDec reads a key/value pair from the wallet storage format
func (db *DB) KVDec(k, v []byte) (result interface{}) {
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
		return []interface{}{id, hashVal, tx}
	case "acentry":
		acc, numB := ser.GetPreLen(keyRem)
		var num uint64
		buf := bytes.NewBuffer([]byte(string(numB)))
		binary.Read(buf, binary.LittleEndian, &num)
		return []interface{}{id, string(acc), num}
	case "key", "wkey":
		pubB, _ := ser.GetPreLen(keyRem)
		if pubEC, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return nil
		} else {
			pub := key.Pub{}
			pub.SetPub(pubEC)
			if string(id) == "key" {
				var priv key.Priv
				priv.Set(v)
				return []interface{}{id, &pub, &priv}
			} else {
				priv, keyRem := ser.GetPreLen(v)
				commentLen := keyRem[16]
				comment := string(keyRem[16 : commentLen+16])
				var wkey Key
				wkey.PrivKey.Set(priv)
				wkey.TimeCreated = BytesToInt64(keyRem[:8])
				wkey.TimeExpires = BytesToInt64(keyRem[8:16])
				wkey.Comment = comment
				return []interface{}{id, &pub, &wkey}
			}
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
