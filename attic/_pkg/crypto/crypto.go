// A library for managing encryption based on the wallet EC cryptography primitives
package crypto

import (
)
const (
	// WalletKeySize is the size of keys (256 bits of course)
	WalletKeySize = 32
	// WalletSaltSize is the size of a salt used to initialise the master key
	WalletSaltSize = 8
)
// Encrypter is the controlling structure for managing en/decryption of wallet data
type Encrypter struct {
	Key      [WalletKeySize]byte
	IV       [WalletSaltSize]byte
	KeyIsSet bool
}
