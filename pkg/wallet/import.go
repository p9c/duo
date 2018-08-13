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
	en       cipher.BlockMode
	de       cipher.BlockMode
	Names    []Name
	Metadata []Metadata
	Keys     []Key
	WKeys    []WKey
	MKeys    []MKey
	CKeys    []CKey
}
type imports interface {
	ToBinaryFormatted() (bf BinaryFormatted, err error)
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
				// logger.Debug(id)
				switch id {
				case "name":
					addrLen := rec[0][idLen] + 1
					addr := string(rec[0][idLen+1 : idLen+addrLen])
					nameLen := rec[1][0] + 1
					name := string(rec[1][1:nameLen])
					// logger.Debug(id, "\""+addr+"\"", "\""+name+"\"")
					var e Name
					e.Addr = []byte(addr)
					e.Name = []byte(name)
					imp.Names = append(imp.Names, e)
				case "key":
					pubLen := rec[0][idLen] + 1
					pubB := rec[0][1:pubLen]
					if pub, err := util.ParsePub(pubB); err != nil {
						return Imports{}
					} else {
						mg, _ := memguard.NewMutableFromBytes(rec[1])
						priv := key.NewPrivFromBytes(mg)
						// logger.Debug(pub, priv)
						var e Key
						e.Pub = util.ToPub(pub)
						e.Priv = priv
						imp.Keys = append(imp.Keys, e)
					}
				case "keymeta":
					pubLen := rec[0][idLen]
					pubB := rec[0][idLen+1 : pubLen+idLen+1]
					// logger.Debug(rec[0], pubLen, pubB)
					if pub, err := util.ParsePub(pubB); err != nil {
						return Imports{}
					} else {
						versionB := rec[1][:4]
						createtimeB := rec[1][4:12]
						// logger.Debug(id, *pub, binary.LittleEndian.Uint32(versionB), time.Unix(int64(binary.LittleEndian.Uint64(createtimeB)), 0))
						var e Metadata
						e.Pub = util.ToPub(pub)
						e.Version = binary.LittleEndian.Uint32(versionB)
						e.CreateTime = time.Unix(int64(binary.LittleEndian.Uint64(createtimeB)), 0)
						imp.Metadata = append(imp.Metadata, e)
					}
				case "wkey":
					pubLen := rec[0][idLen] + 1
					pubB := rec[0][idLen : pubLen+idLen]
					if pub, err := util.ParsePub(pubB); err != nil {
						return Imports{}
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
						var e WKey
						e.Pub = util.ToPub(pub)
						mg, _ := memguard.NewMutableFromBytes(privB)
						priv := key.NewPrivFromBytes(mg)
						e.Priv = priv
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
					var e MKey
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
					var e CKey
					e.Pub = pub
					e.Priv = priv
					imp.CKeys = append(imp.CKeys, e)
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
func (imp *Imports) EncryptData(dst *memguard.LockedBuffer, src []byte) {
	if imp.en != nil {
		En := imp.en
		En.CryptBlocks(dst.Buffer(), src)
	} else {
		dst.Copy(src)
	}
}

// Converts the raw, unencrypted imports into the secure binary format which has all sensitive data encrypted, for writing the initial wallet when importing from a legacy wallet.dat
func (imp *Imports) ToBinaryFormatted() (bf *BinaryFormatted) {
	bf = new(BinaryFormatted)
	for i := range imp.Names {
		bf.AddressBook = append(bf.AddressBook, *new(BAddressBook))
		imp.EncryptData(bf.AddressBook[i].Pub, imp.Names[i].Addr)
		imp.EncryptData(bf.AddressBook[i].Label, imp.Names[i].Name)
	}
	for i := range imp.Keys {
		bf.Key = append(bf.Key, *new(BKey))
		imp.Keys[i].Pub.Compress()
		imp.EncryptData(bf.Key[i].Pub, imp.Keys[i].Pub.Key())
		imp.EncryptData(bf.Key[i].Priv, imp.Keys[i].Priv.Key())
	}
	for i := range imp.CKeys {
		bf.Key = append(bf.Key, *new(BKey))
		Len := len(bf.Key)
		imp.EncryptData(bf.Key[Len].Pub, imp.CKeys[i].Pub)
		imp.EncryptData(bf.Key[Len].Priv, imp.CKeys[i].Priv)
	}
	for i := range imp.WKeys {
		bf.Key = append(bf.Key, *new(BKey))
		Len := len(bf.Key)
		imp.WKeys[i].Pub.Compress()
		imp.EncryptData(bf.Key[Len].Pub, imp.WKeys[i].Pub.Key())
		imp.EncryptData(bf.Key[Len].Priv, imp.WKeys[i].Priv.Key())
		bf.Wdata = append(bf.Wdata, *new(BWdata))
	}
	return
}
