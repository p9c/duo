package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/ec"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"strconv"
	"strings"
	"time"
)

// KVToString converts a key/value pair a simple string format
func (db *DB) KVToString(rec [2][]byte) (d string) {
	for i := range Prefix {
		hasprefix := true
		for j := range Prefix[i] {
			if rec[0][j] != Prefix[i][j] {
				hasprefix = false
				break
			}
		}
		if hasprefix {
			kv := db.KVToVars(rec[0], rec[1])
			result := kv.([]interface{})
			switch i {
			case "name":
				address := result[1].(string)
				name := result[2].(string)
				d += i + " " + address + " " +
					name +
					"\n"
			case "tx":
				txHash := result[1].(Uint.U256)
				tx := result[2].([]byte)
				d += "found tx\n" +
					hex.EncodeToString(txHash.ToBytes()) + " " +
					hex.EncodeToString(tx) + " " +
					"\n"
			case "acentry":
				account := result[1].(string)
				amount := result[2].(uint64)
				d += i + " " +
					account + " " +
					fmt.Sprint(amount) + " " +
					"\n"
			case "key":
				pub := result[1].(*ec.PublicKey)
				priv := result[2].(*key.Priv)
				d += i + " " +
					"pubkey " + hex.EncodeToString(pub.SerializeCompressed()) + " " +
					"privkey " + hex.EncodeToString(priv.Get()) + " " +
					"\n"
			case "wkey":
				pub := result[1].(*ec.PublicKey)
				wkey := result[2].(*Key)
				d += i + " " +
					"pub" + hex.EncodeToString(pub.SerializeCompressed()) + " " +
					"priv" + hex.EncodeToString(wkey.PrivKey.Get()) + " " +
					"created" + time.Unix(wkey.TimeCreated, 0).UTC().String() + " " +
					"expires" + time.Unix(wkey.TimeExpires, 0).UTC().String() + " " +
					"comment '" + wkey.Comment + "'\n"
			case "mkey":
				mkeyID := result[1].(int64)
				masterkey := result[2].(crypto.MasterKey)
				d += i + " ID " + fmt.Sprint(mkeyID) + " " +
					"key " + hex.EncodeToString(masterkey.EncryptedKey) + " " +
					"salt " + hex.EncodeToString(masterkey.Salt) + " " +
					"method " + fmt.Sprint(masterkey.DerivationMethod) + " " +
					"iterations " + fmt.Sprint(masterkey.DeriveIterations) + " " +
					"other " + hex.EncodeToString(masterkey.OtherDerivationParameters) +
					"\n"
			case "ckey":
				pubKey := result[1].(*ec.PublicKey)
				encrypted := result[2].([]byte)
				d += i + " " +
					hex.EncodeToString(pubKey.SerializeCompressed()) + " " +
					hex.EncodeToString(encrypted) +
					"\n"
			case "keymeta":
				key := result[1].(*ec.PublicKey)
				d += i + " " + hex.EncodeToString(key.SerializeCompressed()) + " " +
					"version " + fmt.Sprint(result[2].(uint32)) + " " +
					"created '" + fmt.Sprint(result[3].(time.Time).UTC()) + "' " +
					"\n"
			case "defaultkey":
				d += i + " " + hex.EncodeToString(result[1].(*ec.PublicKey).SerializeCompressed()) + "\n"
			case "pool":
				index := result[1].(uint64)
				version := result[2].(uint32)
				t := time.Unix(result[3].(int64), 0)
				pub := result[4].(*ec.PublicKey)
				d += i + " " +
					"index " + fmt.Sprint(index) + " " +
					"version " + fmt.Sprint(version) + " " +
					"time '" + t.String() + "' " +
					"publickey " + hex.EncodeToString(pub.SerializeUncompressed()) +
					"\n"
			case "version":
				d += i + " " +
					fmt.Sprint(result[1].(uint32)) +
					"\n"
			case "cscript":
				hashID := result[1].(Uint.U160)
				script := result[2].([]byte)
				d += i + " " +
					hex.EncodeToString(hashID.ToBytes()) + " " +
					hex.EncodeToString(script) +
					"\n"
			case "orderposnext":
				d += i + " " +
					fmt.Sprint(result[1].(int64)) + "\n"
			case "account":
				d += i + " " +
					result[1].(string) + " " +
					"\n"
			case "setting":
				name := result[1].(string)
				value := result[2].([]byte)
				d += i + " " +
					name + " "
				if string(name) == "addrIncoming" {
					d += fmt.Sprint(value[0]) + "." +
						fmt.Sprint(value[1]) + "." +
						fmt.Sprint(value[2]) + "." +
						fmt.Sprint(value[3])
				} else {
					d += hex.EncodeToString(value)
				}
				d += "\n"
			case "bestblock":
				d += i + " " +
					hex.EncodeToString(result[1].([]byte)) + "\n"
			case "minversion":
				minversion := result[1].(uint32)
				d += i + " " +
					fmt.Sprint(minversion) + " " +
					"\n"
			default:
				return
			}
			break
		}
	}
	return
}

// StringToVars converts a key/value pair in the format used in a wallet.dat to the variables stored therein
func (db *DB) StringToVars(input string) (result interface{}) {
	s := strings.Split(input, " ")
	keyType := s[0]
	switch keyType {
	case "name":
		address := s[1]
		name := s[2]
		return []interface{}{keyType, address, name}
	case "tx":
		if txHash, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if tx, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			return []interface{}{keyType, txHash, tx}
		}
	case "acentry":
		account := s[1]
		if amount, err := strconv.ParseUint(s[2], 10, 64); err != nil {
			return err
		} else {
			return []interface{}{keyType, account, amount}
		}
	case "key":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return err
		} else if privB, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			var priv key.Priv
			priv.Set(privB)
			return []interface{}{keyType, pub, priv}
		}
	case "wkey":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return err
		} else if privB, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			var priv key.Priv
			priv.Set(privB)
			if created, err := strconv.ParseInt(s[3], 10, 64); err != nil {
				return err
			} else if expires, err := strconv.ParseInt(s[4], 10, 64); err != nil {
				return err
			} else {
				comment := s[5]
				return []interface{}{keyType, pub, priv, created, expires, comment}
			}
		}
	case "mkey":
		if mkey, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if salt, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else if method, err := strconv.ParseUint(s[3], 10, 64); err != nil {
			return err
		} else if iterations, err := strconv.ParseUint(s[4], 10, 64); err != nil {
			return err
		} else if other, err := hex.DecodeString(s[5]); err != nil {
			return err
		} else {
			return []interface{}{keyType, mkey, salt, method, iterations, other}
		}
	case "ckey":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return err
		} else if encrypted, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			return []interface{}{keyType, pub, encrypted}
		}
	case "keymeta":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return err
		} else if version, err := strconv.ParseUint(s[2], 10, 32); err != nil {
			return err
		} else if created, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", s[3]); err != nil {
			return err
		} else {
			return []interface{}{keyType, pub, uint32(version), created}
		}
	case "defaultkey":
		if keyB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ec.ParsePubKey(keyB, ec.S256()); err != nil {
			return err
		} else {
			return []interface{}{keyType, pub}
		}
	case "pool":
		if index, err := strconv.ParseUint(s[1], 10, 64); err != nil {
			return err
		} else if version, err := strconv.ParseUint(s[2], 10, 32); err != nil {
			return err
		} else if t, err := strconv.ParseInt(s[3], 10, 64); err != nil {
			return err
		} else if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ec.ParsePubKey(pubB, ec.S256()); err != nil {
			return err
		} else {
			return []interface{}{keyType, index, uint32(version), t, pub}
		}
	case "version":
		if version, err := strconv.ParseUint(s[1], 10, 32); err != nil {
			return err
		} else {
			return []interface{}{keyType, uint32(version)}
		}
	case "cscript":
		if h, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if script, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			hashID := Uint.Zero160().SetBytes(h)
			return []interface{}{keyType, hashID, script}
		}
	case "orderposnext":
		if orderposnext, err := strconv.ParseInt(s[1], 10, 64); err != nil {
			return err
		} else {
			return []interface{}{keyType, orderposnext}
		}
	case "account":
		return []interface{}{keyType, s[1]}
	case "setting":
		name := s[1]
		var value []byte
		if name == "addrIncoming" {
			ipaddr := strings.Split(s[2], ".")
			for i := range ipaddr {
				if point, err := strconv.Atoi(ipaddr[i]); err != nil {
					return err
				} else {
					value = append(value, byte(point))
				}
			}
			return []interface{}{keyType, name, value}
		} else {
			return errors.New("Unknown setting key")
		}
	case "bestblock":
		if bytes, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else {
			return []interface{}{keyType, bytes}
		}
	case "minversion":
		if minversion, err := strconv.ParseUint(s[1], 10, 32); err != nil {
			return err
		} else {
			return []interface{}{keyType, uint32(minversion)}
		}
	}
	return
}
