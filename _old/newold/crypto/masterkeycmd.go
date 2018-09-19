package crypto

import (
	"github.com/awnumar/memguard"
)

// Sets the master key from a passphrase
func (mk *MasterKey) SetKeyFromPassphrase(keyData memguard.LockedBuffer, salt []byte, rounds, method uint64) (err error) {
   
   return
}
// Encrypt a block of data
func (mk *MasterKey) Encrypt(plainText memguard.LockedBuffer, cipherText []byte) (err error) {
   return
}
// Decrypts a block of data
func (mk *MasterKey) Decrypt(cipherText []byte, plainText memguard.LockedBuffer) (err error) {
   return
}
// Set the master key from a raw binary data and initialization vector
func (mk *MasterKey) SetKey(newKey memguard.LockedBuffer, newIV []byte) (err error) {
   return
}
// Clears the bytes where sensitive data was stored
func (mk *MasterKey) CleanKey() {

}