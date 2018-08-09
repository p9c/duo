package wallet
import (
	"encoding/binary"
	"gitlab.com/parallelcoin/duo/pkg/util"
	"github.com/mitchellh/go-homedir"
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"time"
)
// A collection of tables from a wallet.dat file
type Imports struct {
	Names []Name
	Metadata []Metadata
	Keys []Key
	WKeys []WKey
	MKeys []MKey
	CKeys []CKey
}
// Import reads an existing wallet.dat and returns all the keys and address data in it
func Import(filename ...string) (imp Imports, err error) {
	var db = &BDB{}
	if len(filename) == 0 {
		home, _ := homedir.Dir()
		db.SetFilename(home+"/.parallelcoin/wallet.dat")
	} else {
		db.SetFilename(filename[0])
	}
	if err = db.Open(); err != nil {
		return
	} else if cursor, err := db.Cursor(bdb.NoTransaction); err != nil {
		return Imports{}, err 
	} else {
		rec := [2][]byte{}
		if err = cursor.First(&rec); err != nil {
			return Imports{}, err
		} else {
			for {
				idLen := rec[0][0]+1
				rec[0] = []byte(string(rec[0]))
				rec[1] = []byte(string(rec[1]))
				id := string(rec[0][1:idLen])
				// logger.Debug(id)
				switch id {
				case "name":
					addrLen := rec[0][idLen]+1
					addr := string(rec[0][idLen+1:idLen+addrLen])
					nameLen := rec[1][0]+1
					name := string(rec[1][1:nameLen])
					// logger.Debug(id, "\""+addr+"\"", "\""+name+"\"")
					var e Name
					e.Addr = addr
					e.Name = name
					imp.Names = append(imp.Names, e)
				case "key":
					pubLen := rec[0][idLen]+1
					pubB := rec[0][1:pubLen]
					if pub, err := util.ParsePub(pubB); err != nil {
						return Imports{}, err
					} else {
						priv := util.SetPriv(rec[1])
						// logger.Debug(pub, priv)
						var e Key
						e.Pub = util.ToPub(pub)
						e.Priv = priv
						imp.Keys = append(imp.Keys, e)
					}
				case "keymeta":
					pubLen := rec[0][idLen]
					pubB := rec[0][idLen+1:pubLen+idLen+1]
					// logger.Debug(rec[0], pubLen, pubB)
					if pub, err := util.ParsePub(pubB); err != nil {
						return Imports{}, err
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
					pubLen := rec[0][idLen]+1
					pubB := rec[0][idLen:pubLen+idLen]
					if pub, err := util.ParsePub(pubB); err != nil {
						return Imports{}, err
					} else {
						pLen := rec[1][0]+1
						privB := rec[1][1:pLen]
						tc := rec[1][pLen:pLen+8]
						te := rec[1][pLen+8:pLen+16]
						timeCreated := time.Unix(int64(binary.LittleEndian.Uint64(tc)), 0)
						timeExpires := time.Unix(int64(binary.LittleEndian.Uint64(te)), 0)
						cLen := rec[1][pLen+16]
						comment := string(rec[1][pLen+16:pLen+cLen+16])
						// logger.Debug(id, *pub, timeCreated, timeExpires, "'"+comment+"'")
						var e WKey
						e.Pub = util.ToPub(pub)
						e.Priv = util.SetPriv(privB)
						e.TimeCreated = timeCreated
						e.TimeExpires = timeExpires
						e.Comment = comment
						imp.WKeys = append(imp.WKeys, e)
					}
				case "mkey":
					keyID := int64(binary.LittleEndian.Uint32(rec[0][idLen:idLen+4]))
					ekLen := rec[1][0]+1
					eKey := rec[1][1:ekLen]
					sLen := rec[1][ekLen]
					salt := rec[1][ekLen+1:sLen+ekLen+1]
					method := binary.LittleEndian.Uint32(rec[1][sLen+ekLen+1:sLen+ekLen+5])
					iterations := binary.LittleEndian.Uint32(rec[1][sLen+ekLen+5:sLen+ekLen+9])
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
					pubLen := rec[0][idLen]+1
					pubB := rec[0][idLen+1:pubLen+idLen]
					pub, err := util.ParsePub(pubB)
					if err != nil {
						return Imports{}, err
					}
					pLen := rec[1][0]+1
					priv := rec[1][1:pLen]
					var e CKey
					e.Pub = util.ToPub(pub)
					e.Priv = priv
					imp.CKeys = append(imp.CKeys, e)
				}
				if err = cursor.Next(&rec); err != nil {
					err = cursor.Close()
					break
				}
			}
		}
	}
	return
}
