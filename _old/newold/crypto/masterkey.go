package crypto

import (
	"github.com/awnumar/memguard"
	"github.com/parallelcointeam/duo/pkg/Uint"
)

// Stores the data used to encrypt sensitive data in the wallet when it is set to encrypt
type MasterKey struct {
	EncryptedKey                       []byte
	Salt                               []byte
	DerivationMethod, DeriveIterations uint64
	OtherDerivationParameters          []byte
}
type masterkey interface {
	SetKeyFromPassphrase(keyData, salt []byte, rounds, method uint64) (err error)
	Encrypt(plainText memguard.LockedBuffer, cipherText []byte) (err error)
	Decrypt(cipherText []byte, plainText memguard.LockedBuffer) (err error)
	SetKey(newKey memguard.LockedBuffer, newIV []byte) (err error)
	CleanKey()
}

// Creates a new MasterKey
func NewMasterKey() (mk *MasterKey) {
	mk = new(MasterKey)
	mk.DeriveIterations = 25000
	return
}

func EncryptSecret(masterKey, plainText memguard.LockedBuffer, iv Uint.U256, cipherText []byte) (err error) {
	return
}

func DecryptSecret(masterKey memguard.LockedBuffer, cipherText []byte, iv Uint.U256, plainText memguard.LockedBuffer) (err error) {
	return
}
