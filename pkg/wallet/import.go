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

// A collection of tables from a wallet.dat file with optional en/decryptors
type Imports struct {
	en         cipher.BlockMode
	de         cipher.BlockMode
	Names      []BName
	Metadata   []BMetadata
	Keys       []BKey
	WKeys      []BWKey
	MKeys      []MasterKey
	CKeys      []BCKey
	DefaultKey BDefaultKey
}
type imports interface {
	ToEncryptedStore() (bf EncryptedStore)
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
					var e MasterKey
					e.MKeyID = keyID
					e.EncryptedKey = eKey
					e.Salt = salt
					e.Method = method
					e.Iterations = iterations
					e.Other = other
					imp.MKeys = append(imp.MKeys, e)
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
			b := make([][]byte, len(imp.CKeys))
			for i := range b {
				b[i] = imp.CKeys[i].Priv
			}
			r, _ := imp.MKeys[0].Decrypt(pass, b...)
			for i := range b {
				imp.CKeys[i].Priv = r[i]
			}
			ckey, iv, _ := imp.MKeys[0].DeriveCipher(pass)
			block, _ := aes.NewCipher(ckey.Buffer()[:32])
			imp.de = cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
			imp.en = cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])
		}
	}
	return
}

// Encrypts data from one variable into another if the encrypter is armed
func (imp *Imports) EncryptData(dst, src []byte) {
	if imp.en != nil {
		imp.en.CryptBlocks(dst, src)
	} else {
		dst = src
	}
}

// Converts the raw, unencrypted imports into the secure binary format which has all sensitive data encrypted, for writing the initial wallet when importing from a legacy wallet.dat
func (imp *Imports) ToEncryptedStore() (bf *EncryptedStore) {
	bf = new(EncryptedStore)
	bf.LastLocked = time.Now()
	for i := range imp.Names {
		pub := make([]byte, 48)
		lmod := len(imp.Names[i].Name) % 16
		llen := len(imp.Names[i].Name) + 16 - lmod
		label := make([]byte, llen)
		for j := range imp.Names[i].Addr {
			pub[j] = imp.Names[i].Addr[j]
		}
		for j := range imp.Names[i].Name {
			label[j] = imp.Names[i].Name[j]
		}
		bf.AddressBook = append(bf.AddressBook, *new(AddressBook))
		bf.AddressBook[i].Pub = make([]byte, 48)
		bf.AddressBook[i].Label = make([]byte, llen)
		imp.EncryptData(bf.AddressBook[i].Pub, pub)
		imp.EncryptData(bf.AddressBook[i].Label, label)
	}
	for i := range imp.Keys {
		pub := append(imp.Keys[i].Pub, make([]byte, 15)...)
		priv := imp.Keys[i].Priv
		imp.EncryptData(bf.Key[i].Pub, pub)
		imp.EncryptData(bf.Key[i].Priv, priv)
	}
	for i := range imp.CKeys {
		Len := len(bf.Key)
		bf.Key = append(bf.Key, *new(Key))
		bf.Key[Len].Pub = make([]byte, 48)
		bf.Key[Len].Priv = make([]byte, 48)
		imp.EncryptData(bf.Key[Len].Pub, append(imp.CKeys[i].Pub, make([]byte, 15)...))
		imp.EncryptData(bf.Key[Len].Priv, imp.CKeys[i].Priv)
	}
	for i := range imp.WKeys {
		bf.Key = append(bf.Key, *new(Key))
		k := bf.Key
		Len := len(k)
		imp.EncryptData(k[Len].Pub, imp.WKeys[i].Pub)
		imp.EncryptData(k[Len].Priv, imp.WKeys[i].Priv)
		bf.Wdata = append(bf.Wdata, *new(Wdata))
		imp.EncryptData(bf.Wdata[i].Pub, imp.WKeys[i].Pub)
		var tc, te []byte
		binary.LittleEndian.PutUint64(tc, uint64(imp.WKeys[i].TimeCreated.Unix()))
		binary.LittleEndian.PutUint64(te, uint64(imp.WKeys[i].TimeExpires.Unix()))
		imp.EncryptData(bf.Wdata[i].Created, tc)
		imp.EncryptData(bf.Wdata[i].Expires, te)
		imp.EncryptData(bf.Wdata[i].Comment, []byte(imp.WKeys[i].Comment))
	}
	for i := range imp.Metadata {
		bf.Metadata = append(bf.Metadata, *new(Metadata))
		pub := imp.Metadata[i].Pub
		bf.Metadata[i].Pub = make([]byte, 48)
		imp.EncryptData(bf.Metadata[i].Pub, append(pub, make([]byte, 15)...))
		bf.Metadata[i].Version = imp.Metadata[i].Version
		ct := make([]byte, 8)
		binary.LittleEndian.PutUint64(ct, uint64(imp.Metadata[i].CreateTime.Unix()))
		bf.Metadata[i].CreateTime = make([]byte, 16)
		imp.EncryptData(bf.Metadata[i].CreateTime, append(ct, make([]byte, 8)...))
	}
	for i := range imp.MKeys {
		bf.MasterKey = append(bf.MasterKey, *new(MasterKey))
		bf.MasterKey[i] = imp.MKeys[i]
	}
	bf.DefaultKey = make([]byte, 48)
	imp.EncryptData(bf.DefaultKey, append(imp.DefaultKey.Key, make([]byte, 15)...))
	return
}
