package types

// Crypt is a secure data store for sensitive data. It consists of two fenced memory buffers, one for the password, one for the actual symmetric ciphertext, and a regular buffer that stores the data encrypted. It is intended to be used to minimise the number of times sensitive data ends up being copied around in memory.
type Crypt struct {
	password      *Password
	data          *Bytes
	ciphertext    *LockedBuffer
	locked, armed bool
}
type crypt interface {
	Copy() *Password
	Crypt() *Bytes
	Ciphertext() *LockedBuffer
	Unlock(*Password) *Crypt
	Lock() *Crypt
	IsLocked() *Crypt
	Arm() *Crypt
	Disarm() *Crypt
	IsArmed()
	Encrypt(*LockedBuffer) *Bytes
	Decrypt(*Bytes) *LockedBuffer
}

// NewCrypt makes a new Crypt variable
func NewCrypt() *Crypt {
	return new(Crypt)
}

// Password returns a copy of the password stored in the Crypt
func (c *Crypt) Password() (P *Password) {
	return
}

// Crypt returns a copy of the encrypted data in the Crypt
func (c *Crypt) Crypt() (B *Bytes) {
	return c.data
}

// Ciphertext returns a copy of the unencrypted data stored in the Crypt
func (c *Crypt) Ciphertext() (LB *LockedBuffer) {
	return c.ciphertext
}

// Unlock loads the password into its LockedBuffer
func (c *Crypt) Unlock(*Password) *Crypt {
	return c
}

// Lock disarms the crypt if it is armed, and removes the password
func (c *Crypt) Lock() *Crypt {
	return c
}

// IsLocked returns true if the crypt is locked (does not have the password)
func (c *Crypt) IsLocked() *Crypt {
	return c
}

// Arm runs the KDF on the password and uses the resultant key to unlock the ciphertext enabling the encryption and decryption of data
func (c *Crypt) Arm() *Crypt {
	return c
}

// Disarm clears the unlocked ciphertext and disables encryption/decryption
func (c *Crypt) Disarm() *Crypt {
	return c
}

// IsArmed returns true if the ciphertext has been decrypted
func (c *Crypt) IsArmed() bool {
	return c.armed
}

// Encrypt encrypts the contents of a LockedBuffer and returns the encrypted form as Bytes
func (c *Crypt) Encrypt(*LockedBuffer) (B *Bytes) {
	return
}

// Decrypt decrypts the contents of a Bytes buffer and returns a LockedBuffer containing the decrypted plaintext
func (c *Crypt) Decrypt(*Bytes) (LB *LockedBuffer) {
	return
}
