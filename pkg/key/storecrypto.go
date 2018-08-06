package key
import (
	"crypto/aes"
	"errors"
	"gitlab.com/parallelcoin/duo/pkg/crypto"
)
type CryptedKeys struct {
	Pub    *Pub
	Secret []byte
}
// CryptedKeyMap stores encrypted private keys
type CryptedKeyMap map[*ID]*CryptedKeys
// StoreCrypto keeps also encrypted keys
type StoreCrypto struct {
	Store
	cryptedKeys CryptedKeyMap
	masterKey   crypto.KeyingMaterial
	encrypted   bool
}
type storeCrypto interface {
	SetCrypted() error
	Lock() error
	IsCrypted() bool
	IsLocked() bool
	Unlock(crypto.KeyingMaterial) error
	AddKeyPair(*Priv, *Pub) error
	AddEncryptedKey(*Pub, []byte) error
	GetPriv(*ID) *Priv
	GetPub(*ID) *Pub
	EncryptKeys(crypto.KeyingMaterial) error
}
// NewStoreCrypto creates a new key.StoreCrypto
func NewStoreCrypto() *StoreCrypto {
	return &StoreCrypto{}
}
// SetCrypted encrypts a wallet
func (s *StoreCrypto) SetCrypted() (err error) {
	if s.encrypted {
		return
	}
	if len(s.Keys) > 0 {
		return errors.New("There is unencrypted keys in the store")
	}
	s.Mutex.Lock()
	s.encrypted = true
	s.Mutex.Unlock()
	return
}
// Lock locks the wallet
func (s *StoreCrypto) Lock() (err error) {
	if err = s.SetCrypted(); err != nil {
		return
	}
	s.Mutex.Lock()
	s.masterKey.Clear()
	s.Mutex.Unlock()
	// notification needs to be sent at this point somewhere
	return
}
func (s *StoreCrypto) IsCrypted() bool {
	return s.encrypted
}
func (s *StoreCrypto) IsLocked() bool {
	return s.masterKey == nil
}
// Unlock unlocks an encrypted wallet using a passphrase
func (s *StoreCrypto) Unlock(passphrase crypto.KeyingMaterial) (err error) {
	if err = s.SetCrypted(); err != nil {
		return
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	block, err := NewCipher(passphrase)
	if err != nil {
		return
	}
	var cleartext, recovertext []byte
	for i := range s.cryptedKeys {
		block.Decrypt(cleartext, s.cryptedKeys[i].Secret)
		if len(cleartext) != 32 {
			return errors.New("Recovered key was incorrect length")
		}
		block.Encrypt(recovertext, cleartext)
		for k := range recovertext {
			if s.cryptedKeys[i].Secret[k] != recovertext[k] {
				return errors.New("Passphrase did not unlock one of the keys")
			}
		}
	}
	s.masterKey = cleartext
	return
}
// AddKeyPair adds a new key pair
func (s *StoreCrypto) AddKeyPair(priv *Priv, pub *Pub) (err error) {
	switch {
	case !s.IsCrypted():
		return s.Store.AddKeyPair(priv, pub)
	case s.IsLocked():
		return errors.New("Cannot add keys to locked wallet")
	}
	block, err := NewCipher(s.masterKey)
	var secret []byte
	block.Encrypt(secret, priv.Get())
	err = s.AddEncryptedKey(pub, secret)
	return
}
// AddEncryptedKey adds a new encrypted key
func (s *StoreCrypto) AddEncryptedKey(pub *Pub, secret []byte) (err error) {
	s.Mutex.Lock()
	s.cryptedKeys[pub.GetID()].Pub = pub
	s.cryptedKeys[pub.GetID()].Secret = secret
	s.Mutex.Unlock()
	return
}
// GetPriv gets an private key
func (s *StoreCrypto) GetPriv(id *ID) (priv *Priv, err error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	switch {
	case !s.IsCrypted():
		return s.Store.GetPriv(id)
	case s.IsLocked():
		return nil, errors.New("Wallet is locked")
	}
	if _, ok := s.cryptedKeys[id]; ok {
		block, err := aes.NewCipher(s.masterKey)
		if err != nil {
			return nil, err
		}
		var cleartext []byte
		block.Decrypt(cleartext, s.cryptedKeys[id].Secret)
		priv.Set(cleartext)
		return priv, nil
	}
	return nil, errors.New("Key ID not found in store")
}
// GetPub gets a public key
func (s *StoreCrypto) GetPub(id *ID) (pub *Pub) {
	if !s.IsCrypted() {
		return s.Store.GetPub(id)
	}
	pub = s.cryptedKeys[id].Pub
	return
}
// EncryptKeys encrypts keys with a specified passphrase
func (s *StoreCrypto) EncryptKeys(passphrase crypto.KeyingMaterial) (err error) {
	if len(s.cryptedKeys) < 1 {
		return errors.New("There is no keys to encrypt")
	}
	if s.encrypted {
		return errors.New("The keys are already encrypted")
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.encrypted = true
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return
	}
	for i := range s.Keys {
		var ciphertext, recovertext []byte
		block.Encrypt(ciphertext, s.Keys[i].Private.Get())
		block.Decrypt(recovertext, ciphertext)
		if len(s.cryptedKeys[i].Secret) != len(recovertext) {
			return errors.New("Passphrase did not encrypt one of the keys")
		}
		for j := range recovertext {
			if s.Keys[i].Private.Get()[j] != recovertext[j] {
				return errors.New("Passphrase did not encrypt one of the keys")
			}
		}
		s.cryptedKeys[i].Secret = ciphertext
	}
	return
}
