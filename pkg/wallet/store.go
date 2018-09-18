package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"github.com/awnumar/memguard"
	"time"
)

// Abstract interface embeddable in other types that you wish to convert to other encodings such as JSON and potentially any type of desired encoding, as well as to enable encryption of sensitive data while it is not being utilised
type Serializable struct {
	armed     *MasterKey
	masterKey []*MasterKey
	safe      []*memguard.LockedBuffer
	ckey      *memguard.LockedBuffer
	iv        []byte
}

type serializable interface {
	Lock()
	Unlock(p *memguard.LockedBuffer)
	Copy() (S *Serializable)
}

// Locks the sensitive data that may have been created for an serializable object
func (s *Serializable) Lock() {
	if s.safe != nil {
		DeleteBuffers(s.safe...)
		DeleteBuffer(s.ckey)
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
	MKeyID       int64
	EncryptedKey []byte
	Salt         []byte
	Method       uint32
	Iterations   uint32
	Other        []byte
}

// Encrypts plaintext using the masterkey and a password
func (m *MasterKey) Encrypt(ckey *memguard.LockedBuffer, iv []byte, b ...*memguard.LockedBuffer) (r []*memguard.LockedBuffer, err error) {
	return encDec(true, ckey, iv, b...)
}

//Decrypts ciphertext using the masterkey and a password
func (m *MasterKey) Decrypt(ckey *memguard.LockedBuffer, iv []byte, b ...[]byte) (r []*memguard.LockedBuffer, err error) {
	B := make([]*memguard.LockedBuffer, len(b))
	for i := range b {
		B[i], err = NewBufferFromBytes(b[i])
		if err != nil {
			return
		}
	}
	r, err = encDec(false, ckey, iv, B...)
	DeleteBuffers(B...)
	return
}
func encDec(enc bool, ckey *memguard.LockedBuffer, iv []byte, b ...*memguard.LockedBuffer) (r []*memguard.LockedBuffer, err error) {
	if block, err := aes.NewCipher(ckey.Buffer()[:32]); err != nil {
		return nil, err
	} else {
		var blockmode cipher.BlockMode
		if enc {
			blockmode = cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])
		} else {
			blockmode = cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
		}
		for i := range b {
			if enc {
				b[i], err = PKCS7.Padding(b[i], blockmode.BlockSize())
				r = append(r, b[i])
			}
			blockmode.CryptBlocks(r[i].Buffer(), b[i].Buffer())
			if !enc {
				if r[i], err = PKCS7.Unpadding(r[i], blockmode.BlockSize()); err != nil {
					return nil, err
				}
			}
			// DeleteBuffers(b...)
		}
	}
	return
}

// Creates a new MasterKey
func NewMasterKey() (e MasterKey) {
	return
}

// Stores the details of human readable labels for own and others addresses
type AddressBook struct {
	*Serializable
	Pub   []byte
	Label []byte
}

type addressBook interface {
	Decrypt() (a *AddressBook)
	Encrypt() (a *AddressBook)
}

// Decrypt an addressbook record
func (a *AddressBook) Decrypt() (A AddressBook) {
	d := make([]*memguard.LockedBuffer, 2)
	d[0], _ = NewBuffer(len(a.Pub))
	d[1], _ = NewBuffer(len(a.Label))
	A.Serializable = a.Serializable
	A.safe = d
	r, _ := a.armed.Decrypt(a.ckey, a.iv, d[0].Buffer(), d[1].Buffer())
	A.Pub = r[0].Buffer()
	A.Label = r[1].Buffer()
	return
}

// Encrypt an addressbook record
func (a *AddressBook) Encrypt() (A AddressBook) {
	if a.armed != nil {
		p, _ := NewBufferFromBytes(a.Pub)
		l, _ := NewBufferFromBytes(a.Label)
		r, _ := a.armed.Encrypt(a.ckey, a.iv, p, l)
		r[0].Copy(A.Pub)
		r[1].Copy(A.Label)
	} else {
		A.Pub = a.Pub
		A.Label = a.Label
	}
	return
}

// Creates a new AddressBook
func NewAddressBook(s *Serializable) (a *AddressBook) {
	a = new(AddressBook)
	a.Serializable = s
	return
}

// Stores version and create time for keys
type Metadata struct {
	*Serializable
	Pub        []byte
	Version    uint32
	CreateTime []byte
}

type metadata interface {
	Decrypt() (a Metadata)
}

// Decrypt an metadata record
func (m *Metadata) Decrypt() (M *Metadata) {
	d := make([]*memguard.LockedBuffer, 2)
	d[0], _ = NewBuffer(len(m.Pub))
	d[1], _ = NewBuffer(len(m.CreateTime))
	M.safe = d
	d, _ = m.armed.Decrypt(m.ckey, m.iv, m.Pub, m.CreateTime)
	M.Pub = d[0].Buffer()
	M.CreateTime = d[1].Buffer()
	return
}

// Creates a new
func NewMetadata(s *Serializable) (m *Metadata) {
	m = new(Metadata)
	m.Serializable = s
	return
}

// A public/private key pair corresponding to a wallet address
type Key struct {
	*Serializable
	Pub  []byte
	Priv []byte
}

type ikey interface {
	Decrypt() (a Metadata)
}

// Decrypt an Key record
func (k *Key) Decrypt() (K Key) {
	d := make([]*memguard.LockedBuffer, 2)
	d[0], _ = NewBuffer(len(k.Pub))
	d[1], _ = NewBuffer(len(k.Priv))
	K.safe = d
	d, _ = k.armed.Decrypt(k.ckey, k.iv, k.Pub, k.Priv)
	K.Pub = d[0].Buffer()
	K.Priv = d[1].Buffer()
	return
}

// Creates a new Key
func NewKey(s *Serializable) (k *Key) {
	k = new(Key)
	k.Serializable = s
	return
}

// Extra metadata for key expiry management
type Wdata struct {
	*Serializable
	Pub     []byte
	Created []byte
	Expires []byte
	Comment []byte
}

// Creates a new Wdata
func NewWdata(s *Serializable) (w *Wdata) {
	w.Serializable = s
	return
}

// A raw transaction related to a key pair in the wallet
type Tx struct {
	*Serializable
	TxHash []byte
	TxData []byte
}

// Creates a new Tx
func NewTx(s *Serializable) (t *Tx) {
	t = new(Tx)
	t.Serializable = s
	return
}

// Keys in reserve for future change operations
type Pool struct {
	*Serializable
	Index   uint64
	Version uint32
	Time    []byte
	Pub     []byte
}

// Creates a new Pool
func NewPool(s *Serializable) (p *Pool) {
	p = new(Pool)
	p.Serializable = s
	return
}

// A payment script related to key pairs in the wallet
type Script struct {
	*Serializable
	ID   []byte
	Data []byte
}

// Creates a new Script
func NewScript(s *Serializable) (S *Script) {
	S = new(Script)
	S.Serializable = s
	return
}

// I'm not sure what this is used for
type Account struct {
	*Serializable
	Account []byte
	Version int32
	Pub     []byte
}

// Creates a new Account
func NewAccount(s *Serializable) (a *Account) {
	a.Serializable = s
	return
}

// Stores settings related to the user's wallet
type Setting struct {
	*Serializable
	Name  string
	Value []byte
}

// Creates a new Setting
func NewSetting(s *Serializable) (S *Setting) {
	S = new(Setting)
	S.Serializable = s
	return
}

// A store of all the data related to a wallet with the ability to be encrypted and exported to other data formats
type EncryptedStore struct {
	*Serializable
	MasterKey    []*MasterKey
	LastLocked   time.Time
	AddressBook  []*AddressBook
	Metadata     []*Metadata
	Key          []*Key
	Wdata        []*Wdata
	Tx           []*Tx
	Pool         []*Pool
	Script       []*Script
	Account      []*Account
	Setting      []*Setting
	DefaultKey   []byte
	BestBlock    []byte
	OrderPosNext int64
	Version      uint32
	MinVersion   uint32
}

type encryptedStore interface {
}

// Creates a new EncryptedStore
func NewEncryptedStore() (e *EncryptedStore) {
	e = new(EncryptedStore)
	e.Serializable = new(Serializable)
	return
}
