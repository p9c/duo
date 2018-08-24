package types

import (
	"crypto/rand"
	"errors"
)

// Crypt is a secure data store for sensitive data. It consists of two fenced memory buffers, one for the password, one for the actual symmetric ciphertext, and a regular buffer that stores the data encrypted. It is intended to be used to minimise the number of times sensitive data ends up being copied around in memory.
// The accessors for internal values all make copies into the same type of encapsulation rather than returning a reference to prevent a programmer accidentally messing up a Crypt variable's internals and other parts of an application using the values, so good memory hygiene dictates one should make sure to destroy especially LockedBuffers out of these functions.
type Crypt struct {
	password      *Password
	crypt         *Bytes
	iv            *Bytes
	ciphertext    *LockedBuffer
	cipherIV      *Bytes
	iterations    int
	locked, armed bool
	err           error
}
type crypt interface {
	Generate(*Password) *Crypt
	Load(*Bytes) *Crypt
	Copy() *Crypt
	Password() *Password
	Crypt() *Bytes
	Ciphertext() *LockedBuffer
	IV() *Bytes
	SetRandomIV() *Crypt
	SetIV(*Bytes) *Crypt
	SetIterations(int) *Crypt
	Unlock(*Password) *Crypt
	Lock() *Crypt
	IsLocked() bool
	Arm() *Crypt
	Disarm() *Crypt
	IsArmed() bool
	Encrypt(*LockedBuffer) *Bytes
	Decrypt(*Bytes) *LockedBuffer
	Delete()
}

// NewCrypt makes a new Crypt variable
func NewCrypt() (c *Crypt) {
	c = new(Crypt)
	return
}

// Generate creates a new crypt (password encrypted symmetric key), random initialisation vector, based on a password
func (c *Crypt) Generate(p *Password) {
	c.SetRandomIV()

}

// Load moves the encrypted data into the crypt
func (c *Crypt) Load(b *Bytes) *Crypt {
	c.crypt.FromBytes(b)
	return c
}

// Copy makes a copy of a crypt
func (c *Crypt) Copy() (C *Crypt) {
	C = new(Crypt)
	C.crypt = c.crypt.Copy()
	C.password = c.password.Copy()
	C.ciphertext = c.ciphertext.Copy()
	return
}

// Password returns a copy of the password stored in the Crypt
func (c *Crypt) Password() (P *Password) {
	P = c.password.Copy()
	return
}

// Crypt returns a copy of the encrypted data in the Crypt
func (c *Crypt) Crypt() (B *Bytes) {
	B = c.crypt.Copy()
	return
}

// Ciphertext returns a copy of the unencrypted data stored in the Crypt
func (c *Crypt) Ciphertext() (LB *LockedBuffer) {
	LB = c.ciphertext.Copy()
	return
}

// IV returns the initialisation vector used by the ciphertext, and concatenated with the password in each cycle of the hashchain iteration.
func (c *Crypt) IV() *Bytes {
	return c.iv.Copy()
}

// SetRandomIV returns a new 16 byte initialisation vector from a cryptographically secure random source
func (c *Crypt) SetRandomIV() *Crypt {
	i := *c.iv.Buffer()
	rand.Read(i)
	return c
}

// SetIV sets the initialisation vector from a Bytes
func (c *Crypt) SetIV(B *Bytes) *Crypt {
	c.iv.FromBytes(B)
	return c
}

// SetIterations sets the number of iterations to use in the KDF
func (c *Crypt) SetIterations(I int) *Crypt {
	c.iterations = I
	return c
}

// Unlock loads the password into its LockedBuffer and arms it by default if the crypt has been loaded with encrypted data
func (c *Crypt) Unlock(p *Password) *Crypt {
	if c.password != nil {
		c.password.Delete()
	}
	c.password = p
	c.Arm()
	return c
}

// Lock disarms the crypt if it is armed, and removes the password
func (c *Crypt) Lock() *Crypt {
	c.password.Delete()
	c.Disarm()
	return c
}

// IsLocked returns true if the crypt is locked (does not have the password)
func (c *Crypt) IsLocked() bool {
	return c.locked
}

// Arm runs the KDF on the password and uses the resultant key to unlock the ciphertext from the crypt enabling the encryption and decryption of data
func (c *Crypt) Arm() *Crypt {
	if c.iv != nil && c.crypt != nil && c.armed != true && c.iterations > 0 {
		c.ciphertext, c.cipherIV, c.err = KDF(c.password, c.iv, c.iterations)
	} else {
		c.err = errors.New("Crypt is not fully populated")
	}
	if c.err == nil {
		c.armed = true
	}
	return c
}

// Disarm clears the unlocked ciphertext and disables encryption/decryption
func (c *Crypt) Disarm() *Crypt {
	c.ciphertext.Delete()
	c.cipherIV.Zero()
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

// Delete destroys all the lockedbuffers in a Crypt
func (c *Crypt) Delete() {
	c.ciphertext.Delete()
	c.password.Delete()
	iv := *c.iv.Buffer()
	for i := range iv {
		iv[i] = 0
	}
	crypt := *c.crypt.Buffer()
	for i := range crypt {
		crypt[i] = 0
	}
	c.cipherIV.Zero()
}
