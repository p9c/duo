package wallet

import (
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"time"
)

// Abstract interface embeddable in other types that you wish to convert to other encodings such as JSON and potentially any type of desired encoding, as well as to enable encryption of sensitive data while it is not being utilised
type Serializable struct {
	parent    interface{}
	masterKey []*MasterKey
	en        cipher.BlockMode
	de        cipher.BlockMode
}

type serializable interface {
	SetParent(p interface{}) interface{}
	ToJSON() (j []byte)
	FromJSON(j []byte)
}

// Sets the parent object so the json will be rendered
func (s *Serializable) SetParent(p interface{}) {
	s.parent = &p
}

// Converts the parent object into a json string
func (s *Serializable) ToJSON() (j string) {
	fmt.Println(s.parent)
	J, _ := json.MarshalIndent(s.parent, "", "    ")
	j = string(J)
	return
}

// Converts a JSON string into the data to fill the parent type
func (s *Serializable) FromJSON(j string) (err error) {
	return json.Unmarshal([]byte(j), s.parent)
}

// Stores the details of an encryption master key for a wallet
type MasterKey struct {
	Serializable
	MKeyID       int64
	EncryptedKey []byte
	Salt         []byte
	Method       uint32
	Iterations   uint32
	Other        []byte
}

// Stores the details of human readable labels for own and others addresses
type AddressBook struct {
	Serializable
	Pub   []byte
	Label []byte
}

// Stores version and create time for keys
type Metadata struct {
	Serializable
	Pub        []byte
	Version    uint32
	CreateTime []byte
}

// A public/private key pair corresponding to a wallet address
type Key struct {
	Serializable
	Pub  []byte
	Priv []byte
}

// Extra metadata for key expiry management
type Wdata struct {
	Serializable
	Pub     []byte
	Created []byte
	Expires []byte
	Comment []byte
}

// A raw transaction related to a key pair in the wallet
type Tx struct {
	Serializable
	TxHash []byte
	TxData []byte
}

// Keys in reserve for future change operations
type Pool struct {
	Serializable
	Index   uint64
	Version uint32
	Time    []byte
	Pub     []byte
}

// A payment script related to key pairs in the wallet
type Script struct {
	Serializable
	ID   []byte
	Data []byte
}

// I'm not sure what this is used for
type Account struct {
	Serializable
	Account []byte
	Version int32
	Pub     []byte
}

// Stores settings related to the user's wallet
type Setting struct {
	Serializable
	Name  string
	Value []byte
}

// A store of all the data related to a wallet with the ability to be encrypted and exported to other data formats
type EncryptedStore struct {
	Serializable
	MasterKey    []*MasterKey
	LastLocked   time.Time
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
