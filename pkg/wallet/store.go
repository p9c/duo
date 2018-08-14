package wallet

import (
	"crypto/cipher"
	"encoding/json"
	"github.com/awnumar/memguard"
	"time"
)

// Abstract interface embeddable in other types that you wish to convert to other encodings such as JSON and potentially any type of desired encoding, as well as to enable encryption of sensitive data while it is not being utilised
type Serializable struct {
	masterKey *[]*MasterKey
	en        *cipher.BlockMode
	de        *cipher.BlockMode
	safe      *[]*memguard.LockedBuffer
}

type serializable interface {
	EncryptData()
	DecryptData()
	Wipe()
}

// Encrypts data from one variable into another if the encrypter is armed
func (s *Serializable) EncryptData(dst, src []byte) {
	if s.en != nil {
		(*s.en).CryptBlocks(dst, src)
	} else {
		dst = src
	}
}

// Decrypts data from one variable into another if the encrypter is armed
func (s *Serializable) DecryptData(dst, src []byte) {
	if s.en != nil {
		(*s.en).CryptBlocks(dst, src)
	} else {
		dst = src
	}
}

// Destroys the sensitive data that may have been created for an AddressBook entry
func (s *Serializable) Wipe() {
	for i := range *s.safe {
		(*s.safe)[i].Destroy()
	}
}

// Converts an object into a JSON string
func ToJSON(p interface{}) (j string) {
	J, _ := json.MarshalIndent(p, "", "    ")
	j = string(J)
	return
}

// Converts a JSON string into the data to fill the parent type
func (s *Serializable) FromJSON(j string) (J string, err error) {
	err = json.Unmarshal([]byte(j), J)
	return J, err
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

// Creates a new MasterKey
func NewMasterKey() (e MasterKey) {
	return
}

// Stores the details of human readable labels for own and others addresses
type AddressBook struct {
	Serializable
	Pub   []byte
	Label []byte
}

type addressBook interface {
	Decrypt() (a *AddressBook, err error)
}

// Decrypt an addressbook record
func (a *AddressBook) Decrypt() (A AddressBook) {
	d := make([]*memguard.LockedBuffer, 2)
	d[0], _ = memguard.NewMutable(len(a.Pub))
	d[1], _ = memguard.NewMutable(len(a.Label))
	a.safe = &d
	a.DecryptData((*a.safe)[0].Buffer(), a.Pub)
	a.DecryptData((*a.safe)[1].Buffer(), a.Label)
	A.Pub = (*a.safe)[0].Buffer()
	A.Label = (*a.safe)[1].Buffer()
	return
}

// Creates a new AddressBook
func NewAddressBook() (e AddressBook) {
	return
}

// Stores version and create time for keys
type Metadata struct {
	Serializable
	Pub        []byte
	Version    uint32
	CreateTime []byte
}

// Creates a new
func NewMetadata() (e Metadata) {
	return
}

// A public/private key pair corresponding to a wallet address
type Key struct {
	Serializable
	Pub  []byte
	Priv []byte
}

// Creates a new Key
func NewKey() (e Key) {
	return
}

// Extra metadata for key expiry management
type Wdata struct {
	Serializable
	Pub     []byte
	Created []byte
	Expires []byte
	Comment []byte
}

// Creates a new Wdata
func NewWdata() (e Wdata) {
	return
}

// A raw transaction related to a key pair in the wallet
type Tx struct {
	Serializable
	TxHash []byte
	TxData []byte
}

// Creates a new Tx
func NewTx() (e Tx) {
	return
}

// Keys in reserve for future change operations
type Pool struct {
	Serializable
	Index   uint64
	Version uint32
	Time    []byte
	Pub     []byte
}

// Creates a new Pool
func NewPool() (e Pool) {
	return
}

// A payment script related to key pairs in the wallet
type Script struct {
	Serializable
	ID   []byte
	Data []byte
}

// Creates a new Script
func NewScript() (e Script) {
	return
}

// I'm not sure what this is used for
type Account struct {
	Serializable
	Account []byte
	Version int32
	Pub     []byte
}

// Creates a new Account
func NewAccount() (e Account) {
	return
}

// Stores settings related to the user's wallet
type Setting struct {
	Serializable
	Name  string
	Value []byte
}

// Creates a new Setting
func NewSetting() (e Setting) {
	return
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

type encryptedStore interface {
}

// Creates a new EncryptedStore
func NewEncryptedStore() (e EncryptedStore) {
	return
}
