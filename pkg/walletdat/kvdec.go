package walletdat

import (
	"gitlab.com/parallelcoin/duo/pkg/wallet"
	"gitlab.com/parallelcoin/duo/pkg/util"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
	"gitlab.com/parallelcoin/duo/pkg/ec"
	"gitlab.com/parallelcoin/duo/pkg/serialize"
	"time"
)

// KVDec reads a key/value pair from the wallet storage format.
// Only the database entries containing key pairs and related data are imported, as the rest of the data can be computed using the blockchain and derived indices and metadata
func KVDec(k, v []byte) (result interface{}) {
	var id string
	keyRem, _ := ser.Deserialize(k, id)
	switch string(id) {
	case "name":
		// name entries are human readable labels for public keys, both keys controlled by the user as well as correspondent accounts in the address book
		var addr, name string
		ser.Deserialize(keyRem, addr)
		ser.Deserialize(v, name)
		return []interface{}{id, addr, name}
	case "keymeta":
		// keymeta entries store metadata about key pairs in the wallet
		var keyB []byte
		ser.Deserialize(keyRem, keyB)
		if pub, err := ec.ParsePubKey(keyB, ec.S256()); err != nil {
			return err
		} else {
			var version uint32
			ser.Deserialize([]byte(string(v))[:4], version)
			var ct int64
			ser.Deserialize([]byte(string(v))[4:], ct)
			createtime := time.Unix(ct, 0)
			return []interface{}{id, util.ToPub(pub), version, createtime}
		}
	case "key", "wkey":
		// wkey entries are key entries with metadata about creation, expiration and comments
		var pubB []byte
		ser.Deserialize(keyRem, pubB)
		if pub, err := util.ParsePub(pubB); err != nil {
			return err
		} else {
			if string(id) == "key" {
				// key entries are a simple unencrypted public/private key pair
				priv := util.SetPriv(v)
				return []interface{}{id, priv, util.ToPub(pub)}
			} else {
				var priv []byte
				keyRem, _ := ser.Deserialize(v, priv)
				var tc, te int64
				keyRem, _ = ser.Deserialize(keyRem, tc)
				keyRem, _ = ser.Deserialize(keyRem, te)
				var comment string
				keyRem, _ = ser.Deserialize(keyRem, comment)
				wkey := &wallet.Key{
					PrivKey:     util.SetPriv(priv),
					TimeCreated: tc,
					TimeExpires: te,
					Comment:     comment,
				}
				return []interface{}{id, wkey, util.ToPub(pub)}
			}
		}
	// mkey entries are a passphrase encrypted master key that encrypts ckey entries' private keys
	case "mkey":
		var mkeyID int64
		ser.Deserialize(keyRem[:8], mkeyID)
		var mkey, salt []byte
		keyRem, _ = ser.Deserialize(v, mkey)
		keyRem, _ = ser.Deserialize(keyRem, salt)
		var method, iterations uint32
		keyRem, _ = ser.Deserialize(keyRem, method)
		other, _ := ser.Deserialize(keyRem, iterations)
		masterkey := &crypto.MasterKey{
			EncryptedKey:              mkey,
			Salt:                      salt,
			DerivationMethod:          method,
			DeriveIterations:          iterations,
			OtherDerivationParameters: other,
		}
		return []interface{}{id, mkeyID, masterkey}
	// ckey entries are a public key and the corresponding encrypted private key
	case "ckey":
		var pubB []byte
		ser.Deserialize(keyRem, pubB)
		var priv []byte
		ser.Deserialize(v, priv)
		if pub, err := util.ParsePub(pubB); err != nil {
			return err
		} else {
			return []interface{}{id, util.ToPub(pub), priv}
		}

	}
	return nil
}
