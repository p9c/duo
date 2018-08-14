package wallet

import (
	"crypto/cipher"
	"encoding/json"
	"time"
)

type Serializable struct{}
type serializable interface {
	ToJSON() (j []byte)
	FromJSON(j []byte)
}

type MasterKey struct {
	Serializable
	MKeyID       int64
	EncryptedKey []byte
	Salt         []byte
	Method       uint32
	Iterations   uint32
	Other        []byte
}

type AddressBook struct {
	Serializable
	Pub   []byte
	Label []byte
}
type Metadata struct {
	Serializable
	Pub        []byte
	Version    uint32
	CreateTime []byte
}
type Key struct {
	Serializable
	Pub  []byte
	Priv []byte
}
type Wdata struct {
	Serializable
	Pub     []byte
	Created []byte
	Expires []byte
	Comment []byte
}
type Tx struct {
	Serializable
	TxHash []byte
	TxData []byte
}
type Pool struct {
	Serializable
	Index   uint64
	Version uint32
	Time    []byte
	Pub     []byte
}
type Script struct {
	Serializable
	ID   []byte
	Data []byte
}
type Account struct {
	Serializable
	Account []byte
	Version int32
	Pub     []byte
}
type Setting struct {
	Serializable
	Name  string
	Value []byte
}
type EncryptedStore struct {
	Serializable
	LastLocked   time.Time
	en, de       cipher.BlockMode
	MasterKey    []MasterKey
	AddressBook  []AddressBook
	Metadata     []Metadata
	Key          []Key
	Wdata        []Wdata
	Tx           []Tx
	Pool         []Pool
	Script       []Script
	Account      []Account
	Setting      []Setting
	DefaultKey   []byte
	BestBlock    []byte
	OrderPosNext int64
	Version      uint32
	MinVersion   uint32
}

func (es *EncryptedStore) ToJSON() (j []byte) {
	j, _ = json.MarshalIndent(es, "", "    ")
	return
}
