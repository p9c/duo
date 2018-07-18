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
			kv := db.KVDec(rec[0], rec[1])
			result := kv.([]interface{})
			switch i {
			case "name":
				address := result[1].(string)
				name := result[2].(string)
				d += i + " " + address + " " + name + "\n"
			case "tx":
				txHash := hex.EncodeToString(result[1].(*Uint.U256).ToBytes())
				tx := hex.EncodeToString(result[2].([]byte))
				d += i + " " + txHash + " " + tx + " " + "\n"
			case "acentry":
				account := result[1].(string)
				amount := fmt.Sprint(result[2].(uint64))
				d += i + " " + account + " " + amount + " " + "\n"
			case "key":
				pub := hex.EncodeToString(result[1].(*key.Pub).GetPub().SerializeUncompressed())
				priv := hex.EncodeToString(result[2].(*key.Priv).Get())
				d += i + " " +
					"pubkey " + pub + " " + "privkey " + priv + " " + "\n"
			case "wkey":
				pub := hex.EncodeToString(result[1].(*key.Pub).GetPub().SerializeUncompressed())
				wkey := hex.EncodeToString(result[2].(*Key).PrivKey.Get())
				created := time.Unix(result[2].(*Key).TimeCreated, 0).UTC().String()
				expires := time.Unix(result[2].(*Key).TimeExpires, 0).UTC().String()
				comment := result[2].(*Key).Comment
				d += i + " " +
					"pub" + pub + " " + "priv" + wkey + " " +
					"created" + created + " " + "expires" + expires + " " +
					"comment '" + comment + "'\n"
			case "mkey":
				mkeyID := fmt.Sprint( result[1].(int64))
				masterkey := result[2].(*crypto.MasterKey)
				d += i + " ID " + mkeyID + " " +
					"key " + hex.EncodeToString(masterkey.EncryptedKey) + " " +
					"salt " + hex.EncodeToString(masterkey.Salt) + " " +
					"method " + fmt.Sprint(masterkey.DerivationMethod) + " " +
					"iterations " + fmt.Sprint(masterkey.DeriveIterations) + " " +
					"other " + hex.EncodeToString(masterkey.OtherDerivationParameters) +
					"\n"
			case "ckey":
				pubKey := hex.EncodeToString(result[1].(*key.Pub).GetPub().SerializeUncompressed())
				encrypted := hex.EncodeToString(result[2].([]byte))
				d += i + " " + pubKey + " " + encrypted + "\n"
			case "keymeta":
				key := hex.EncodeToString(result[1].(*key.Pub).GetPub().SerializeUncompressed())
				d += i + " " + key + " " +
					"version " + fmt.Sprint(result[2].(uint32)) + " " +
					"created '" + fmt.Sprint(result[3].(time.Time).UTC()) + "' " +
					"\n"
			case "defaultkey":
				d += i + " " + hex.EncodeToString(result[1].(*key.Pub).GetPub().SerializeCompressed()) + "\n"
			case "pool":
				index := fmt.Sprint(result[1].(uint64))
				version := fmt.Sprint(result[2].(uint32))
				t := time.Unix(result[3].(int64), 0).String()
				pub := hex.EncodeToString(result[4].(*key.Pub).GetPub().SerializeUncompressed())
				d += i + " " + "index " + index + " " + 
					"version " + version + " " + "time '" + t + "' " + "publickey " + pub + "\n"
			case "version":
				d += i + " " + fmt.Sprint(result[1].(uint32)) + "\n"
			case "cscript":
				hashID := hex.EncodeToString(result[1].(*Uint.U160).ToBytes())
				script := hex.EncodeToString(result[2].([]byte))
				d += i + " " + hashID + " " + script + "\n"
			case "orderposnext":
				d += i + " " + fmt.Sprint(result[1].(int64)) + "\n"
			case "account":
				d += i + " " + result[1].(string) + " " + "\n"
			case "setting":
				name := result[1].(string)
				value := result[2].([]byte)
				d += i + " " + name + " "
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
				d += i + " " + hex.EncodeToString(result[1].([]byte)) + "\n"
			case "minversion":
				minversion := fmt.Sprint(result[1].(uint32))
				d += i + " " + minversion + " " + "\n"
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
	id := s[0]
	switch id {
	case "name":
		address := s[1]
		name := s[2]
		return []interface{}{id, address, name}
	case "tx":
		if txHash, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if tx, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			return []interface{}{id, txHash, tx}
		}
	case "acentry":
		account := s[1]
		if amount, err := strconv.ParseUint(s[2], 10, 64); err != nil {
			return err
		} else {
			return []interface{}{id, account, amount}
		}
	case "key":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ParsePub(pubB); err != nil {
			return err
		} else if privB, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			var priv key.Priv
			priv.Set(privB)
			return []interface{}{id, pub, priv}
		}
	case "wkey":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ParsePub(pubB); err != nil {
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
				return []interface{}{id, pub, priv, created, expires, comment}
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
			return []interface{}{id, mkey, salt, method, iterations, other}
		}
	case "ckey":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ParsePub(pubB); err != nil {
			return err
		} else if encrypted, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			return []interface{}{id, pub, encrypted}
		}
	case "keymeta":
		if pubB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ParsePub(pubB); err != nil {
			return err
		} else if version, err := strconv.ParseUint(s[2], 10, 32); err != nil {
			return err
		} else if created, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", s[3]); err != nil {
			return err
		} else {
			return []interface{}{id, pub, uint32(version), created}
		}
	case "defaultkey":
		if keyB, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if pub, err := ec.ParsePubKey(keyB, ec.S256()); err != nil {
			return err
		} else {
			return []interface{}{id, pub}
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
		} else if pub, err := ParsePub(pubB); err != nil {
			return err
		} else {
			return []interface{}{id, index, uint32(version), t, pub}
		}
	case "version":
		if version, err := strconv.ParseUint(s[1], 10, 32); err != nil {
			return err
		} else {
			return []interface{}{id, uint32(version)}
		}
	case "cscript":
		if h, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else if script, err := hex.DecodeString(s[2]); err != nil {
			return err
		} else {
			hashID := Uint.Zero160().SetBytes(h)
			return []interface{}{id, hashID, script}
		}
	case "orderposnext":
		if orderposnext, err := strconv.ParseInt(s[1], 10, 64); err != nil {
			return err
		} else {
			return []interface{}{id, orderposnext}
		}
	case "account":
		return []interface{}{id, s[1]}
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
			return []interface{}{id, name, value}
		} else {
			return errors.New("Unknown setting key")
		}
	case "bestblock":
		if bytes, err := hex.DecodeString(s[1]); err != nil {
			return err
		} else {
			return []interface{}{id, bytes}
		}
	case "minversion":
		if minversion, err := strconv.ParseUint(s[1], 10, 32); err != nil {
			return err
		} else {
			return []interface{}{id, uint32(minversion)}
		}
	}
	return
}
