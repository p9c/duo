package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"github.com/awnumar/memguard"
	"github.com/mitchellh/go-homedir"
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/util"
	"time"
)

// Stores an address book entry in a wallet.dat
type BName struct {
	Addr []byte
	Name []byte
}

// Stores key metadata in a wallet.dat
type BMetadata struct {
	Pub        []byte
	Version    uint32
	CreateTime time.Time
}

// Stores unencrypted keys in a wallet.dat
type BKey struct {
	Pub  []byte
	Priv []byte
}

// An unencrypted key pair with extra metadata for managing expiry in a wallet.dat
type BWKey struct {
	Pub         []byte
	Priv        []byte
	TimeCreated time.Time
	TimeExpires time.Time
	Comment     string
}

// A key pair with plaintext public key and AES-256-CBC encrypted private key
type BCKey struct {
	Pub  []byte
	Priv []byte
}

// Stores the default key that will appear in a wallet interface when creating a payment request
type BDefaultKey struct {
	Key []byte
}

// A collection of tables from a wallet.dat file with optional en/decryptors
type Imports struct {
	Serializable
	Names      []BName
	Metadata   []BMetadata
	Keys       []BKey
	WKeys      []BWKey
	CKeys      []BCKey
	DefaultKey BDefaultKey
}
type imports interface {
	ToEncryptedStore() (es EncryptedStore)
	EncryptData(dst *memguard.LockedBuffer, src []byte)
}

// Import reads an existing wallet.dat and returns all the keys and address data in it. If a password is given, the private keys in the CKeys array are decrypted and the encrypter/decrypter functions are armed.
func Import(pass *memguard.LockedBuffer, filename ...string) (imp Imports) {
	var db = &BDB{}
	if len(filename) == 0 {
		home, _ := homedir.Dir()
		db.SetFilename(home + "/.parallelcoin/wallet.dat")
	} else {
		db.SetFilename(filename[0])
	}
	if err := db.Open(); err != nil {
		return
	} else if cursor, err := db.Cursor(bdb.NoTransaction); err != nil {
		return
	} else {
		rec := [2][]byte{}
		if err := cursor.First(&rec); err != nil {
			return
		} else {
			for {
				idLen := rec[0][0] + 1
				rec[0] = []byte(string(rec[0]))
				rec[1] = []byte(string(rec[1]))
				id := string(rec[0][1:idLen])
				switch id {
				case "name":
					addrLen := rec[0][idLen] + 1
					addr := string(rec[0][idLen+1 : idLen+addrLen])
					nameLen := rec[1][0] + 1
					name := string(rec[1][1:nameLen])
					// logger.Debug(id, "\""+addr+"\"", "\""+name+"\"")
					var e BName
					e.Addr = []byte(addr)
					e.Name = []byte(name)
					imp.Names = append(imp.Names, e)
				case "key":
					pubLen := rec[0][idLen] + 1
					pubB := rec[0][1:pubLen]
					if pub, err := util.ParsePub(pubB); err != nil {
						return
					} else {
						mg, _ := memguard.NewMutableFromBytes(rec[1])
						priv := key.NewPrivFromBytes(mg)
						// logger.Debug(pub, priv)
						var e BKey
						e.Pub = pub.SerializeCompressed()
						e.Priv = priv.Key()
						imp.Keys = append(imp.Keys, e)
					}
				case "keymeta":
					pubLen := rec[0][idLen]
					pubB := rec[0][idLen+1 : pubLen+idLen+1]
					// logger.Debug(rec[0], pubLen, pubB)
					if pub, err := util.ParsePub(pubB); err != nil {
						return
					} else {
						versionB := rec[1][:4]
						createtimeB := rec[1][4:12]
						// logger.Debug(id, *pub, binary.LittleEndian.Uint32(versionB), time.Unix(int64(binary.LittleEndian.Uint64(createtimeB)), 0))
						var e BMetadata
						e.Pub = pub.SerializeCompressed()
						e.Version = binary.LittleEndian.Uint32(versionB)
						e.CreateTime = time.Unix(int64(binary.LittleEndian.Uint64(createtimeB)), 0)
						imp.Metadata = append(imp.Metadata, e)
					}
				case "wkey":
					pubLen := rec[0][idLen] + 1
					pubB := rec[0][idLen : pubLen+idLen]
					if pub, err := util.ParsePub(pubB); err != nil {
						return
					} else {
						pLen := rec[1][0] + 1
						privB := rec[1][1:pLen]
						tc := rec[1][pLen : pLen+8]
						te := rec[1][pLen+8 : pLen+16]
						timeCreated := time.Unix(int64(binary.LittleEndian.Uint64(tc)), 0)
						timeExpires := time.Unix(int64(binary.LittleEndian.Uint64(te)), 0)
						cLen := rec[1][pLen+16]
						comment := string(rec[1][pLen+16 : pLen+cLen+16])
						// logger.Debug(id, *pub, timeCreated, timeExpires, "'"+comment+"'")
						var e BWKey
						e.Pub = pub.SerializeCompressed()
						mg, _ := memguard.NewMutableFromBytes(privB)
						priv := key.NewPrivFromBytes(mg)
						e.Priv = priv.Key()
						e.TimeCreated = timeCreated
						e.TimeExpires = timeExpires
						e.Comment = comment
						imp.WKeys = append(imp.WKeys, e)
					}
				case "mkey":
					keyID := int64(binary.LittleEndian.Uint32(rec[0][idLen : idLen+4]))
					ekLen := rec[1][0] + 1
					eKey := rec[1][1:ekLen]
					sLen := rec[1][ekLen]
					salt := rec[1][ekLen+1 : sLen+ekLen+1]
					method := binary.LittleEndian.Uint32(rec[1][sLen+ekLen+1 : sLen+ekLen+5])
					iterations := binary.LittleEndian.Uint32(rec[1][sLen+ekLen+5 : sLen+ekLen+9])
					other := rec[1][sLen+ekLen+9:]
					e := new(MasterKey)
					e.MKeyID = keyID
					e.EncryptedKey = eKey
					e.Salt = salt
					e.Method = method
					e.Iterations = iterations
					e.Other = other
					if imp.masterKey == nil {
						imp.masterKey = new([]*MasterKey)
					}
					mk := append(*imp.masterKey, e)
					(*imp.masterKey) = mk
				case "ckey":
					pubLen := rec[0][idLen] + 1
					pub := rec[0][idLen+1 : pubLen+idLen]
					privLen := rec[1][0] + 1
					priv := rec[1][1:privLen]
					var e BCKey
					e.Pub = pub
					e.Priv = priv
					imp.CKeys = append(imp.CKeys, e)
				case "defaultkey":
					klen := rec[1][0] + 1
					defaultKey := rec[1][1:klen]
					var e BDefaultKey
					e.Key = defaultKey
					imp.DefaultKey = e
				}
				if err = cursor.Next(&rec); err != nil {
					err = cursor.Close()
					break
				}
			}
		}
		if pass != nil {
			ckey, iv, _ := (*imp.masterKey)[0].DeriveCipher(pass)
			block, _ := aes.NewCipher(ckey.Buffer()[:32])
			de := cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
			en := cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])
			imp.de = &de
			imp.en = &en
			b := make([][]byte, len(imp.CKeys))
			for i := range b {
				b[i] = imp.CKeys[i].Priv
			}
			r, _ := (*imp.masterKey)[0].Decrypt(pass, b...)
			for i := range b {
				imp.CKeys[i].Priv = r[i]
			}
		} else {

		}
	}
	return
}

// Converts the raw, unencrypted imports into the secure binary format which has all sensitive data encrypted, for writing the initial wallet when importing from a legacy wallet.dat
func (imp *Imports) ToEncryptedStore() (es EncryptedStore) {
	es.en = imp.en
	es.de = imp.de
	es.masterKey = imp.masterKey
	es.LastLocked = time.Now()
	for i := range imp.Names {
		pub := make([]byte, 48)
		for j := range imp.Names[i].Addr {
			pub[j] = imp.Names[i].Addr[j]
		}
		for j := len(imp.Names[i].Addr); j < 48; j++ {
			pub[j] = byte(48 - len(imp.Names[i].Addr))
		}
		lmod := len(imp.Names[i].Name) % 16
		llen := len(imp.Names[i].Name) + 16 - lmod
		label := make([]byte, llen)
		for j := range imp.Names[i].Name {
			label[j] = imp.Names[i].Name[j]
		}
		for j := len(imp.Names[i].Name); j < llen; j++ {
			label[j] = byte(llen - len(imp.Names[i].Name))
		}
		es.AddressBook = append(es.AddressBook, NewAddressBook())
		es.AddressBook[i].Pub = make([]byte, 48)
		es.AddressBook[i].Label = make([]byte, llen)
		es.EncryptData(es.AddressBook[i].Pub, pub)
		es.EncryptData(es.AddressBook[i].Label, label)
		es.AddressBook[i].en = imp.en
		es.AddressBook[i].de = imp.de
		es.AddressBook[i].masterKey = imp.masterKey
	}
	for i := range imp.Keys {
		es.Key = append(es.Key, NewKey())
		pub := append(imp.Keys[i].Pub, make([]byte, 15)...)
		priv := imp.Keys[i].Priv
		es.EncryptData(es.Key[i].Pub, pub)
		es.EncryptData(es.Key[i].Priv, priv)
		es.Key[i].en = imp.en
		es.Key[i].de = imp.de
		es.Key[i].masterKey = imp.masterKey

	}
	for i := range imp.CKeys {
		Len := len(es.Key)
		es.Key = append(es.Key, NewKey())
		es.Key[Len].Pub = make([]byte, 48)
		es.Key[Len].Priv = make([]byte, 48)
		es.EncryptData(es.Key[Len].Pub, append(imp.CKeys[i].Pub, make([]byte, 15)...))
		es.EncryptData(es.Key[Len].Priv, imp.CKeys[i].Priv)
		es.Key[i].en = imp.en
		es.Key[i].de = imp.de
		es.Key[i].masterKey = imp.masterKey

	}
	for i := range imp.WKeys {
		es.Key = append(es.Key, NewKey())
		k := es.Key
		Len := len(k)
		es.EncryptData(k[Len].Pub, imp.WKeys[i].Pub)
		es.EncryptData(k[Len].Priv, imp.WKeys[i].Priv)
		es.Wdata = append(es.Wdata, *new(Wdata))
		es.EncryptData(es.Wdata[i].Pub, imp.WKeys[i].Pub)
		var tc, te []byte
		binary.LittleEndian.PutUint64(tc, uint64(imp.WKeys[i].TimeCreated.Unix()))
		binary.LittleEndian.PutUint64(te, uint64(imp.WKeys[i].TimeExpires.Unix()))
		es.EncryptData(es.Wdata[i].Created, tc)
		es.EncryptData(es.Wdata[i].Expires, te)
		es.EncryptData(es.Wdata[i].Comment, []byte(imp.WKeys[i].Comment))
		es.Key[i].en = imp.en
		es.Key[i].de = imp.de
		es.Key[i].masterKey = imp.masterKey
	}
	for i := range imp.Metadata {
		es.Metadata = append(es.Metadata, *new(Metadata))
		pub := imp.Metadata[i].Pub
		es.Metadata[i].Pub = make([]byte, 48)
		es.EncryptData(es.Metadata[i].Pub, append(pub, make([]byte, 15)...))
		es.Metadata[i].Version = imp.Metadata[i].Version
		ct := make([]byte, 8)
		binary.LittleEndian.PutUint64(ct, uint64(imp.Metadata[i].CreateTime.Unix()))
		es.Metadata[i].CreateTime = make([]byte, 16)
		es.EncryptData(es.Metadata[i].CreateTime, append(ct, make([]byte, 8)...))
		es.Metadata[i].en = imp.en
		es.Metadata[i].de = imp.de
		es.Metadata[i].masterKey = imp.masterKey
	}
	for i := range *imp.masterKey {
		mk := NewMasterKey()
		*es.masterKey = append(*es.masterKey, &mk)
		(*es.masterKey)[i] = (*imp.masterKey)[i]
	}
	es.DefaultKey = make([]byte, 48)
	es.EncryptData(es.DefaultKey, append(imp.DefaultKey.Key, make([]byte, 15)...))
	return
}
