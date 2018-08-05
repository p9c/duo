package walletdat

import (
	"gitlab.com/parallelcoin/duo/pkg/util"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/ec"
	"time"
)

// KVDec reads a key/value pair from the wallet storage format.
// Only the database entries containing key pairs and related data are imported, as the rest of the data can be computed using the blockchain and derived indices and metadata
func KVDec(kv [2][]byte) (result interface{}) {
	k, v := kv[0], kv[1]
	res, keyRem, _ := Deserialize(k, String)
	id := *res[0].(string)
	switch id {
	case "name":
		// name entries are human readable labels for public keys, both keys controlled by the user as well as correspondent accounts in the address book
		addrs, _, _ := Deserialize(keyRem, String)
		names, _, _ := Deserialize(v, String)
		return []interface{}{id, addrs[0].(string), names[0].(string)}
	case "keymeta":
		// keymeta entries store metadata about key pairs in the wallet
		pubs, _, _ := Deserialize(keyRem, Bytes)
		res, _, _ := Deserialize(v, Uint32, Int64)
		createtime := time.Unix(res[2].(int64), 0)
		if pub, err := ec.ParsePubKey(pubs[0].([]byte), ec.S256()); err != nil {
			return err
		} else {
			return []interface{}{id, util.ToPub(pub), res[0].(uint32), createtime}
		}
	case "key":
		pubs, _, _ := Deserialize(keyRem, Bytes)
		if pub, err := util.ParsePub(pubs[0].([]byte)); err != nil {
			return err
		} else {
			priv := util.SetPriv(v)
			return []interface{}{id, util.ToPub(pub), priv}
		}
	case "wkey":
		// wkey entries are key entries with metadata about creation, expiration and comments
		pubs, _, _ := Deserialize(keyRem, Bytes)
		if pub, err := util.ParsePub(pubs[0].([]byte)); err != nil {
			return err
		} else {
			res, _, _ := Deserialize(v, Bytes, Int64, Int64, String)
			wkey := &WKey{
				Priv:     util.SetPriv(res[0].([]byte)),
				TimeCreated: res[1].(int64),
				TimeExpires: res[2].(int64),
				Comment:     res[3].(string),
			}
			return []interface{}{id, wkey, util.ToPub(pub)}
		}
	// mkey entries are a passphrase encrypted master key that encrypts ckey entries' private keys
	case "mkey":
		mkeyID, _, _ := Deserialize(keyRem, Uint32)
		res, keyRem, _ = Deserialize(v, Bytes, Bytes, Uint32, Uint32)
		masterkey := &crypto.MasterKey{
			EncryptedKey:              res[0].([]byte),
			Salt:                      res[1].([]byte),
			DerivationMethod:          res[2].(uint32),
			DeriveIterations:          res[3].(uint32),
			OtherDerivationParameters: keyRem,
		}
		return []interface{}{id, mkeyID[0].(uint32), masterkey}
	// ckey entries are a public key and the corresponding encrypted private key
	case "ckey":
		pubs, _, _ := Deserialize(keyRem, Bytes)
		privs, _, _ := Deserialize(v, Bytes)
		if pub, err := util.ParsePub(pubs[0].([]byte)); err != nil {
			return err
		} else {
			return []interface{}{id, util.ToPub(pub), privs[0].([]byte)}
		}

	}
	return nil
}
