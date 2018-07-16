package crypto

const (
	// WalletKeySize is the size of keys (256 bits of course)
	WalletKeySize = 32
	// WalletSaltSize is the size of a salt used to initialise the master key
	WalletSaltSize = 8
)

// MasterKey is the key used to encrypt the contents of a wallet
type MasterKey struct {
	EncryptedKey                       []byte
	Salt                               []byte
	DerivationMethod, DeriveIterations uint32
	OtherDerivationParameters          []byte
}

// KeyingMaterial is the raw bytes used in the encryption, locked by the user password (this requires memguard)
type KeyingMaterial []byte

type keyingMaterial interface {
	Clear()
	Bytes() []byte
	FromBytes([]byte)
}

// Clear the buffer
func (k *KeyingMaterial) Clear() {
	km := *k
	for i := range km {
		km[i] = 0
	}
	k = &KeyingMaterial{}
}

// Bytes contained in the variable
func (k *KeyingMaterial) Bytes() (b []byte) {
	return []byte(*k)
}

// FromBytes puts bytes into the variable
func (k *KeyingMaterial) FromBytes(b []byte) {
	*k = KeyingMaterial(b)
}

// Encrypter is the data structure for managing encryption of the wallet keys
type Encrypter struct {
	Key      [WalletKeySize]byte
	IV       [WalletSaltSize]byte
	KeyIsSet bool
}
