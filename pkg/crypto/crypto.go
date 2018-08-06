// A library for managing encryption based on the wallet EC cryptography primitives
package crypto
const (
	// WalletKeySize is the size of keys (256 bits of course)
	WalletKeySize = 32
	// WalletSaltSize is the size of a salt used to initialise the master key
	WalletSaltSize = 8
)
type MasterKey struct {
	EncryptedKey                       []byte
	Salt                               []byte
	DerivationMethod, DeriveIterations uint32
	OtherDerivationParameters          []byte
}
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
// Encrypter is the controlling structure for managing en/decryption of wallet data
type Encrypter struct {
	Key      [WalletKeySize]byte
	IV       [WalletSaltSize]byte
	KeyIsSet bool
}
