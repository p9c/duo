package wallet

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/ec"
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
		num := BytesToUint64(numB)
		return []interface{}{id, string(acc), num}
	case "key", "wkey":
		pubB, _ := ser.GetPreLen(keyRem)
		if pub, err := ParsePub(pubB); err != nil {
			return nil
		} else {
			if string(id) == "key" {
				priv := SetPriv(v)
				return []interface{}{id, ToPub(pub), priv}
			} else {
				priv, keyRem := ser.GetPreLen(v)
				commentLen := keyRem[16]
				comment := string(keyRem[16 : commentLen+16])
				wkey := &Key{
					PrivKey: SetPriv(priv),
					TimeCreated: BytesToInt64(keyRem[:8]),
					TimeExpires: BytesToInt64(keyRem[8:16]),
					Comment: comment,
				}
				return []interface{}{id, ToPub(pub), wkey}
			}
		}
	case "mkey":
		mkeyID := BytesToInt64(keyRem[:8])
		mkey, keyRem := ser.GetPreLen(v)
		salt, keyRem := ser.GetPreLen(keyRem)
		method, iterations := BytesToUint32(keyRem[:4]), BytesToUint32(keyRem[4:8])
		other := keyRem[8:]
		masterkey := &crypto.MasterKey{
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
		if pub, err := ParsePub(pubB); err != nil {
			return err
		} else {
			return []interface{}{id, ToPub(pub), priv}
		}
	case "keymeta":
		keyB, _ := ser.GetPreLen(keyRem)
		if pub, err := ec.ParsePubKey(keyB, ec.S256()); err != nil {
			return err
		} else {
			version := BytesToUint32([]byte(string(v))[:4])
			createtime := time.Unix(BytesToInt64([]byte(string(v))[4:]), 0)
			return []interface{}{id, ToPub(pub), version, createtime}
		}
	case "defaultkey":
		pubB, _ := ser.GetPreLen(v)
		if pub, err := ParsePub(pubB); err != nil {
			return nil
		} else {
			return []interface{}{id, ToPub(pub)}
		}
	case "pool":
		index := BytesToUint64(keyRem)
		version := BytesToUint32([]byte(string(v))[:4])
		ptime := BytesToInt64([]byte(string(v))[4:12])
		pubB := []byte(string(v))[13:]
		if pub, err := ParsePub(pubB); err != nil {
			return nil
		} else {
			return []interface{}{id, index, version, ptime, ToPub(pub)}
		}
	case "version":
		version := BytesToUint32([]byte(string(v))[:4])
		return []interface{}{id, version}
	case "cscript":
		hID, _ := ser.GetPreLen(keyRem)
		hashID := Uint.Zero160().SetBytes(hID)
		script, _ := ser.GetPreLen(v)
		return []interface{}{id, hashID, script}
	case "orderposnext":
		orderposnext := BytesToInt64(v)
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
		minversion := BytesToUint32(v)
		return []interface{}{id, minversion}
	}
	return nil
}
